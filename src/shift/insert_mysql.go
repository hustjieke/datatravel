/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package shift

import (
	"fmt"
	"strings"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/schema"
)

func (h *EventHandler) InsertMySQLRow(e *canal.RowsEvent, systemTable bool) {
	var conn *client.Conn
	h.wg.Add(1)

	executeFunc := func(conn *client.Conn) {
		defer h.wg.Done()
		var keep = true

		for i, row := range e.Rows {
			var values []string

			// keep connection in the loop, just put conn to pool when execute the last row
			if (i + 1) == len(e.Rows) {
				keep = false
			}

			var isEnum bool
			for idx, v := range row {
				values = append(values, h.ParseValue(e, idx, v))
				if e.Table.Columns[idx].Type == schema.TYPE_ENUM {
					isEnum = true
				}
			}

			if isEnum {
				query := &Query{
					sql:       fmt.Sprintf("insert ignore into `%s`.`%s` values (%s)", e.Table.Schema, e.Table.Name, strings.Join(values, ",")),
					typ:       QueryType_INSERT,
					skipError: systemTable,
				}
				h.execute(conn, keep, query)
			} else {
				query := &Query{
					sql:       fmt.Sprintf("insert into `%s`.`%s` values (%s)", e.Table.Schema, e.Table.Name, strings.Join(values, ",")),
					typ:       QueryType_INSERT,
					skipError: systemTable,
				}
				h.execute(conn, keep, query)
			}
		}
	}

	if h.xaConn != nil {
		conn = h.xaConn
	} else {
		if conn = h.shift.toPool.Get(); conn == nil {
			h.shift.panicMe("shift.insert.get.to.conn.nil.error")
		}
	}

	// if e.DataType == canal.BINLOGDATA {
	// Binlog sync
	if e.Header != nil {
		executeFunc(conn)
	} else {
		// canal.DUMPDATA, Backend worker for mysqldump.
		go func(conn *client.Conn) {
			executeFunc(conn)
		}(conn)
	}
}
