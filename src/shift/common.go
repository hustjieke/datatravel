/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package shift

import (
	"xlog"

	"github.com/xelabs/go-mysqlstack/sqlparser/depends/common"
)

var sysDatabases = map[string]bool{
	"sys":                true,
	"mysql":              true,
	"information_schema": true,
	"performance_schema": true,
}

func logPanicHandler(log *xlog.Log, format string, v ...interface{}) {
	log.Fatal(format, v...)
}

// EscapeBytes used to escape the literal byte.
func EscapeBytes(bytes []byte) []byte {
	buffer := common.NewBuffer(128)
	for _, b := range bytes {
		// See https://dev.mysql.com/doc/refman/5.7/en/string-literals.html
		// for more information on how to escape string literals in MySQL.
		switch b {
		case 0:
			buffer.WriteString(`\0`)
		case '\'':
			buffer.WriteString(`\'`)
		case '"':
			buffer.WriteString(`\"`)
		case '\b':
			buffer.WriteString(`\b`)
		case '\n':
			buffer.WriteString(`\n`)
		case '\r':
			buffer.WriteString(`\r`)
		case '\t':
			buffer.WriteString(`\t`)
		case 0x1A:
			buffer.WriteString(`\Z`)
		case '\\':
			buffer.WriteString(`\\`)
		default:
			buffer.WriteU8(b)
		}
	}
	return buffer.Datas()
}
