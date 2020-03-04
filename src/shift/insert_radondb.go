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
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/common"
)

func (h *EventHandler) InsertRadonDBRow(e *canal.RowsEvent, systemTable, isNotFisrtTime bool) {
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

			for idx, v := range row {
				if v == nil {
					values = append(values, fmt.Sprintf("NULL"))
					continue
				}

				if _, ok := v.([]byte); ok {
					values = append(values, fmt.Sprintf("%q", v))
				} else {
					switch {
					case e.Table.Columns[idx].Type == schema.TYPE_NUMBER:
						// In case dbs ---> db, db1.tbl and db2.tbl`s auto_increment may be conflicted
						// if e.Table.Columns[idx].IsAuto == true && isNotFisrtTime {
						// values = append(values, "0")
						if e.Table.Columns[idx].IsAuto {
							continue
						} else {
							values = append(values, fmt.Sprintf("%d", v))
						}
					case e.Table.Columns[idx].Type == schema.TYPE_BIT:
						// Here we should add prefix "0x" for hex
						values = append(values, fmt.Sprintf("0x%x", v))
					default:
						switch e.Table.Columns[idx].RawType {
						case "tinyblob", "blob", "mediumblob", "longblob":
							// Here we should add prefix "0x" for hex
							values = append(values, fmt.Sprintf("0x%x", v))
						default:
							s := fmt.Sprintf("%v", v)
							values = append(values, fmt.Sprintf("\"%s\"", EscapeBytes(common.StringToBytes(s))))
						}
					}
				}
			}

			cols := common.NewBuffer(256)
			len := len(e.Table.Columns)
			for idx, col := range e.Table.Columns {
				// Skip auto_increment col
				if col.IsAuto {
					continue
				}
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
			h.execute(conn, keep, query)
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
		tables, ok := h.shift.cfg.DBTablesMaps[e.Table.Schema]
		if ok {
			// 过滤
			for _, tbl := range tables {
				if e.Table.Name == tbl {
					executeFunc(conn)
				}
			}
		}
	} else {
		// canal.DUMPDATA, Backend worker for mysqldump.
		go func(conn *client.Conn) {
			executeFunc(conn)
		}(conn)
	}
}
