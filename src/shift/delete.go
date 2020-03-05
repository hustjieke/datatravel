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
)

func (h *EventHandler) DeleteRow(e *canal.RowsEvent) {
	var conn *client.Conn

	h.wg.Add(1)
	executeFunc := func(conn *client.Conn) {
		defer h.wg.Done()
		var keep = true

		pks := e.Table.PKColumns
		for i, row := range e.Rows {
			var values []string

			// keep connection in the loop, just put conn to pool when execute the last row
			if (i + 1) == len(e.Rows) {
				keep = false
			}

			// We have pk columns.
			if len(pks) > 0 {
				for _, pk := range pks {
					v := row[pk]
					values = append(values, fmt.Sprintf("`%s`=%s", e.Table.Columns[pk].Name, h.ParseValue(e, pk, v)))
				}
			} else {
				for j, v := range row {
					if v == nil {
						continue
					}
					values = append(values, fmt.Sprintf("`%s`=%s", e.Table.Columns[j].Name, h.ParseValue(e, j, v)))
				}
			}

			query := &Query{
				sql:       fmt.Sprintf("delete from `%s`.`%s` where %s", e.Table.Schema, e.Table.Name, strings.Join(values, " and ")),
				typ:       QueryType_DELETE,
				skipError: false,
			}
			h.execute(conn, keep, query)
		}
	}

	if h.xaConn != nil {
		conn = h.xaConn
	} else {
		if conn = h.shift.toPool.Get(); conn == nil {
			h.shift.panicMe("shift.delete.get.to.conn.nil.error")
		}
	}

	// executeFunc(conn)
	tables, ok := h.shift.cfg.DBTablesMaps[e.Table.Schema]
	if ok {
		// 过滤
		for _, tbl := range tables {
			if e.Table.Name == tbl {
				executeFunc(conn)
			}
		}
	}
}
