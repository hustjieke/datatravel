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

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/schema"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/common"
)

func (h *EventHandler) ParseValue(e *canal.RowsEvent, idx int, v interface{}) string {
	if v == nil {
		return fmt.Sprintf("NULL")
	}

	if _, ok := v.([]byte); ok {
		return fmt.Sprintf("%q", v)
	} else {
		switch {
		case e.Table.Columns[idx].Type == schema.TYPE_NUMBER:
			return fmt.Sprintf("%d", v)
		case e.Table.Columns[idx].Type == schema.TYPE_BIT:
			// Here we should add prefix "0x" for hex
			return fmt.Sprintf("0x%x", v)
		default:
			switch e.Table.Columns[idx].RawType {
			case "tinyblob", "blob", "mediumblob", "longblob":
				// Here we should add prefix "0x" for hex
				str := fmt.Sprintf("0x%x", v)
				// If str is empty, we`ll got "0x"
				if str == "0x" {
					return "\"\""
				}
				return str
			default:
				s := fmt.Sprintf("%v", v)
				return fmt.Sprintf("\"%s\"", EscapeBytes(common.StringToBytes(s)))
			}
		}
	}
}
