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
	"strconv"
	"strings"

	"github.com/imroc/biu"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/schema"
)

func (h *EventHandler) InsertMySQLRow(e *canal.RowsEvent, systemTable bool) {
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
					switch {
					case e.Table.Columns[idx].Type == schema.TYPE_NUMBER:
						values = append(values, fmt.Sprintf("%d", v))
					case e.Table.Columns[idx].Type == schema.TYPE_BIT:
						// Here no need to add prefix "0x" for hexstr
						hexstr := fmt.Sprintf("%x", v)
						log.Debug("bit hexstr:", hexstr)
						if num64, err := strconv.ParseUint(hexstr, 16, 64); err != nil {
							panic(err)
						} else {
							num64bit := biu.ToBinaryString(num64)
							num64bit = strings.Replace(num64bit, " ", "", -1)
							num64bit = strings.TrimLeft(num64bit, "[")
							num64bit = strings.TrimRight(num64bit, "]")
							values = append(values, fmt.Sprintf("B'%s'", num64bit))
						}
					default:
						switch e.Table.Columns[idx].RawType {
						case "tinyblob", "blob", "mediumblob", "longblob":
							// Here we should add prefix "0x" for hex
							values = append(values, fmt.Sprintf("0x%x", v))
						default:
							log.Debug("insert table type and raw type:", e.Table.Name, e.Table.Columns[idx].Type, e.Table.Columns[idx].RawType)
							values = append(values, fmt.Sprintf("%#v", v))
						}
					}
				}
			}

			query := &Query{
				sql:       fmt.Sprintf("insert into `%s`.`%s` values (%s)", e.Table.Schema, e.Table.Name, strings.Join(values, ",")),
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
