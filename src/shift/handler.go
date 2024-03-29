/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package shift

import (
	"config"
	"strings"
	"sync"
	"xlog"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/client"
)

type QueryType int

const (
	QueryType_INSERT      QueryType = 0
	QueryType_DELETE      QueryType = 1
	QueryType_UPDATE      QueryType = 2
	QueryType_XA_START    QueryType = 3
	QueryType_XA_END      QueryType = 4
	QueryType_XA_PREPARE  QueryType = 5
	QueryType_XA_COMMIT   QueryType = 6
	QueryType_XA_ROLLBACK QueryType = 7
)

type Query struct {
	sql       string
	typ       QueryType
	skipError bool
}

type EventHandler struct {
	wg    sync.WaitGroup
	log   *xlog.Log
	shift *Shift
	canal.DummyEventHandler
	xaConn *client.Conn
}

func NewEventHandler(log *xlog.Log, shift *Shift) *EventHandler {
	return &EventHandler{
		log:    log,
		shift:  shift,
		xaConn: nil,
	}
}

// OnRow used to handle the Insert/Delete/Update events.
func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
	switch e.Action {
	case canal.InsertAction:
		_, isSystem := sysDatabases[strings.ToLower(e.Table.Schema)]
		if h.shift.cfg.ToFlavor == config.ToMySQLFlavor ||
			h.shift.cfg.ToFlavor == config.ToMariaDBFlavor {
			h.InsertMySQLRow(e, isSystem)
		} else {
			h.InsertRadonDBRow(e, isSystem, h.shift.cfg.IsNotFisrtTime)
		}
	case canal.DeleteAction:
		h.DeleteRow(e)
	case canal.UpdateAction:
		h.UpdateRow(e)
	default:
		h.shift.panicMe("shift.handler.unsupported.event[%+v]", e)
	}
	return nil
}

// OnTableChanged used to handle the QueryEvent and XAEvent.
func (h *EventHandler) OnTableChanged(schema string, table string) error {
	cfg := h.shift.cfg

	if cfg.FromDatabase == schema && cfg.FromTable == table {
		h.shift.panicMe("datatravel.cant.do.ddl[%v, %v].during.shifting...", schema, table)
	}
	return nil
}

func (h *EventHandler) OnXA(e *canal.XAEvent) error {
	h.XAQuery(e)
	return nil
}

func (h *EventHandler) WaitWorkerDone() {
	h.wg.Wait()
}

func (h *EventHandler) execute(conn *client.Conn, keep bool, query *Query) {
	sql := query.sql
	log := h.log
	shift := h.shift
	pool := h.shift.toPool
	cfg := h.shift.cfg

	switch query.typ {
	case QueryType_INSERT, QueryType_DELETE, QueryType_UPDATE:
		{
			execFn := func() {
				if _, err := conn.Execute(sql); err != nil {
					// if query.skipError {
					if true {
						log.Error("shift.execute.sql[%s].error:%+v", sql, err)
					} else {
						log.Error("shift.execute.sql[%s].error:%+v", sql, err)
						shift.panicMe("shift.execute.sql[%s].error:%+v", sql, err)
					}
				}
			}

			if h.xaConn != nil {
				// Transactional query.
				execFn()
			} else {
				// Not transactional query.
				execFn()
				if !keep {
					pool.Put(conn)
				}
			}
		}
	case QueryType_XA_START:
		{
			// Prepare a xaConn for xa transaction
			if _, err := conn.Execute(sql); err != nil {
				shift.panicMe("shift.execute.sql[%s].error:%+v", sql, err)
			}
			h.xaConn = conn
		}
	case QueryType_XA_END:
		{
			if conn == nil {
				shift.panicMe("shift.xa.end.can't.find.xa.connection.sql[%s]", sql)
			}

			if _, err := conn.Execute(sql); err != nil {
				shift.panicMe("shift.execute.sql[%s].error:%+v", sql, err)
			}
		}
	case QueryType_XA_PREPARE:
		{
			if conn == nil {
				shift.panicMe("shift.xa.prepare.can't.find.xa.connection.sql[%s]", sql)
			}

			if _, err := conn.Execute(sql); err != nil {
				shift.panicMe("shift.execute.sql[%s].error:%+v", sql, err)
			}

			// Close xa connection and set xaConn to nil
			if err := conn.Close(); err != nil {
				shift.panicMe("shift.xa.prepare.xa.connection.close.error:%+v", err)
			}
			h.xaConn = nil

			// Put a new connection to ToPool
			Newconn, err := client.Connect(cfg.To, cfg.ToUser, cfg.ToPassword, "")
			if err != nil {
				shift.panicMe("shift.xa.prepare.new.connection.error:%+v", err)
			}
			pool.Put(Newconn)
		}
	case QueryType_XA_COMMIT, QueryType_XA_ROLLBACK:
		{
			if _, err := conn.Execute(sql); err != nil {
				log.Error("shift.execute.sql[%s].error:%+v", sql, err)
			}
			if !keep {
				pool.Put(conn)
			}
		}
	}
}
