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
	"github.com/xelabs/go-mysqlstack/common"
)

func (h *EventHandler) InsertRadonDBRow(e *canal.RowsEvent, systemTable bool) {
	var conn *client.Conn
	log := h.log
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

			for idx, v := range row {
				if v == nil {
					values = append(values, fmt.Sprintf("NULL"))
					continue
				}

				if _, ok := v.([]byte); ok {
					values = append(values, fmt.Sprintf("%q", v))
				} else {
					if e.Table.Columns[idx].Type == schema.TYPE_NUMBER {
						values = append(values, fmt.Sprintf("%d", v))
					} else {
						values = append(values, fmt.Sprintf("%#v", v))
					}
				}
			}

			// TODO if not travel RadonDB, we don`t need to add columns
			cols := common.NewBuffer(256)
			len := len(e.Table.Columns)
			for idx, col := range e.Table.Columns {
				cols.WriteString(col.Name)
				if idx != (len - 1) {
					cols.WriteString(",")
				}
			}
			columns, _ := cols.ReadStringNUL()

			query := &Query{
				sql:       fmt.Sprintf("insert into `%s`.`%s`(%s) values (%s)", e.Table.Schema, e.Table.Name, columns, strings.Join(values, ",")),
				typ:       QueryType_INSERT,
				skipError: systemTable,
			}
			log.Debug("----no:%d, query:%+v", i, query)
			h.execute(conn, keep, query)
		}
	}

	if h.xaConn != nil {
		conn = h.xaConn
	} else {
		conn = h.shift.toPool.Get()
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