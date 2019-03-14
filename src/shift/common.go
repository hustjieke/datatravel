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
