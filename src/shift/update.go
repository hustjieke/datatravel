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

func (h *EventHandler) UpdateRow(e *canal.RowsEvent) {
	var conn *client.Conn

	h.wg.Add(1)
	executeFunc := func(conn *client.Conn) {
		defer h.wg.Done()
		var keep = true

		rows := e.Rows
		pks := e.Table.PKColumns
		for i := 0; i < len(rows); i += 2 {
			var values []string
			var wheres []string

			// keep connection in the loop, just put conn to pool when execute the last row
			if (i + 2) == len(e.Rows) {
				keep = false
			}

			// Old image.
			v1Row := rows[i]
			// New image.
			v2Row := rows[i+1]

			// We have pk columns.
			if len(pks) > 0 {
				for _, pk := range pks {
					v := v1Row[pk]
					wheres = append(wheres, fmt.Sprintf("`%s`=%s", e.Table.Columns[pk].Name, h.ParseValue(e, pk, v)))
				}
			}

			for i := range v2Row {
				v2 := v2Row[i]
				if v2 != nil {
					values = append(values, fmt.Sprintf("`%s`=%s", e.Table.Columns[i].Name, h.ParseValue(e, i, v2)))
				}

				if len(pks) == 0 {
					v1 := v1Row[i]
					if v1 != nil {
						wheres = append(wheres, fmt.Sprintf("`%s`=%s", e.Table.Columns[i].Name, h.ParseValue(e, i, v1)))
					}
				}
			}
			query := &Query{
				sql:       fmt.Sprintf("update `%s`.`%s` set %s where %s", e.Table.Schema, e.Table.Name, strings.Join(values, ","), strings.Join(wheres, " and ")),
				typ:       QueryType_UPDATE,
				skipError: false,
			}
			h.execute(conn, keep, query)
		}
	}

	if h.xaConn != nil {
		conn = h.xaConn
	} else {
		if conn = h.shift.toPool.Get(); conn == nil {
			h.shift.panicMe("shift.update.get.to.conn.nil.error")
		}
	}

	executeFunc(conn)
	//tables, ok := h.shift.cfg.DBTablesMaps[e.Table.Schema]
	//if ok {
	//	// 过滤
	//	for _, tbl := range tables {
	//		if e.Table.Name == tbl {
	//			executeFunc(conn)
	//			break
	//		}
	//	}
	//}
}
