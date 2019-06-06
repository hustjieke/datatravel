/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package shift

import (
	"strings"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/client"
)

// XAQuery execute xa statement (xa start|end|prepare|commit|rollback)
func (h *EventHandler) XAQuery(e *canal.XAEvent) {
	var XaType QueryType
	var conn *client.Conn

	if strings.Contains(string(e.Query), "XA START") {
		XaType = QueryType_XA_START
	} else if strings.Contains(string(e.Query), "XA END") {
		XaType = QueryType_XA_END
	} else if strings.Contains(string(e.Query), "XA PREPARE") {
		XaType = QueryType_XA_PREPARE
	} else if strings.Contains(string(e.Query), "XA COMMIT") {
		XaType = QueryType_XA_COMMIT
	} else if strings.Contains(string(e.Query), "XA ROLLBACK") {
		XaType = QueryType_XA_ROLLBACK
	} else {
		h.shift.panicMe("shift.handler.unsupported.XAQueryEvent[%+v]", e)
	}

	query := &Query{
		sql:       strings.Split(string(e.Query), ",")[0],
		typ:       XaType,
		skipError: false,
	}

	if h.xaConn != nil {
		conn = h.xaConn
	} else {
		conn = h.shift.toPool.Get()
	}

	if XaType == QueryType_XA_ROLLBACK || XaType == QueryType_XA_COMMIT {
		h.execute(conn, false, query)
	} else {
		h.execute(conn, true, query)
	}
}
