/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package shift

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"xlog"

	"github.com/ant0ine/go-json-rest/rest"
)

var (
	restfulPort = 8181

	// Config for normal shift.
	mockCfg = &Config{
		From:         "127.0.0.1:3306",
		FromUser:     "root",
		FromDatabase: "shift_test_from",
		FromTable:    "t1",

		To:         "127.0.0.1:3306",
		ToUser:     "root",
		ToDatabase: "shift_test_to",
		ToTable:    "t1",

		Cleanup:   true,
		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		RadonURL:  fmt.Sprintf("http://127.0.0.1:%d", restfulPort),
		Checksum:  true,
	}

	// Config for system (mysql) shift.
	mockCfgMysql = &Config{
		From:         "127.0.0.1:3306",
		FromUser:     "root",
		FromDatabase: "mysql",
		FromTable:    "user",

		To:         "127.0.0.1:3306",
		ToUser:     "root",
		ToDatabase: "mysql",
		ToTable:    "userx",

		Cleanup:   true,
		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		RadonURL:  fmt.Sprintf("http://127.0.0.1:%d", restfulPort),
		Checksum:  false,
	}

	// Config for xa shift.
	mockCfgXa = &Config{
		From:         "127.0.0.1:3306",
		FromUser:     "root",
		FromDatabase: "shift_test_from",
		FromTable:    "t1",

		To:         "127.0.0.1:3307",
		ToUser:     "root",
		ToDatabase: "shift_test_to",
		ToTable:    "t1",

		Cleanup:   true,
		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		RadonURL:  fmt.Sprintf("http://127.0.0.1:%d", restfulPort),
		Checksum:  true,
	}

	// Config for ddl shift.
	mockCfgDDL = &Config{
		From:         "127.0.0.1:3306",
		FromUser:     "root",
		FromDatabase: "shift_test_from",
		FromTable:    "t1",

		To:         "127.0.0.1:3306",
		ToUser:     "root",
		ToDatabase: "shift_test_to",
		ToTable:    "t1",

		Cleanup:   true,
		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		RadonURL:  fmt.Sprintf("http://127.0.0.1:%d", restfulPort),
		Checksum:  false,
	}
)

func mockShift(log *xlog.Log, cfg *Config, hasPK bool, initData bool, readonlyHanler mockHandler, shardshiftHandler mockHandler, throttleHandler mockHandler) (*Shift, func()) {
	h := mockHttp(log, restfulPort, readonlyHanler, shardshiftHandler, throttleHandler)
	shift := NewShift(log, cfg)

	// Prepare connections.
	{
		if err := shift.prepareConnection(); err != nil {
			log.Panicf("mock.shift.prepare.connection.error:%+v", err)
		}
	}

	// Prepare the from database and table.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)
		toConn := shift.toPool.Get()
		defer shift.toPool.Put(toConn)

		// Cleanup To table first.
		{
			sql := fmt.Sprintf("drop table if exists `%s`.`%s`", cfg.ToDatabase, cfg.ToTable)
			if _, err := toConn.Execute(sql); err != nil {
				log.Panicf("mock.shift.prepare.table.error:%+v", err)
			}
		}

		if _, isSystem := sysDatabases[strings.ToLower(cfg.FromDatabase)]; !isSystem {
			// Cleanup To table first.
			{
				sql := fmt.Sprintf("drop table if exists `%s`.`%s`", cfg.FromDatabase, cfg.FromTable)
				if _, err := fromConn.Execute(sql); err != nil {
					log.Panicf("mock.shift.prepare.database.error:%+v", err)
				}
			}

			// Create database on from.
			sql := fmt.Sprintf("create database if not exists `%s`", cfg.FromDatabase)
			if _, err := fromConn.Execute(sql); err != nil {
				log.Panicf("mock.shift.prepare.database.error:%+v", err)
			}

			// Create table on from.
			if hasPK {
				sql = fmt.Sprintf("create table `%s`.`%s`(a int primary key, b int, c varchar(200), d DOUBLE NULL DEFAULT NULL, e json DEFAULT NULL, f INT UNSIGNED DEFAULT NULL, g BIGINT DEFAULT NULL, h BIGINT UNSIGNED DEFAULT NULL, i TINYINT NULL, j TINYINT UNSIGNED DEFAULT NULL, k SMALLINT DEFAULT NULL, l SMALLINT UNSIGNED DEFAULT NULL, m MEDIUMINT DEFAULT NULL, n MEDIUMINT UNSIGNED DEFAULT NULL)", cfg.FromDatabase, cfg.FromTable)
			} else {
				sql = fmt.Sprintf("create table `%s`.`%s`(a int, b int, c varchar(200), d DOUBLE NULL DEFAULT NULL, e json DEFAULT NULL, f INT UNSIGNED DEFAULT NULL, g BIGINT DEFAULT NULL, h BIGINT UNSIGNED DEFAULT NULL, i TINYINT NULL, j TINYINT UNSIGNED DEFAULT NULL, k SMALLINT DEFAULT NULL, l SMALLINT UNSIGNED DEFAULT NULL, m MEDIUMINT DEFAULT NULL, n MEDIUMINT UNSIGNED DEFAULT NULL)", cfg.FromDatabase, cfg.FromTable)
			}
			if _, err := fromConn.Execute(sql); err != nil {
				log.Panicf("mock.shift.prepare.database.error:%+v", err)
			}

			if initData {
				for i := 100; i < 108; i++ {
					sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
					if _, err := fromConn.Execute(sql); err != nil {
						log.Panicf("mock.shift.prepare.datas.error:%+v", err)
					}
				}
			}
		} else {
			// Prepare mysql.userx(fakes for mysql.user) table on TO.
			sql := fmt.Sprintf("show create table `%s`.`%s`", cfg.FromDatabase, cfg.FromTable)
			r, err := fromConn.Execute(sql)
			if err != nil {
				log.Panicf("mock.prepare.mysql.userx.error:%+v", err)
			}
			sql, _ = r.GetString(0, 1)
			sql = strings.Replace(sql, fmt.Sprintf("CREATE TABLE `%s`", cfg.FromTable), fmt.Sprintf("CREATE TABLE `%s`.`%s`", cfg.ToDatabase, cfg.ToTable), 1)
			if _, err = toConn.Execute(sql); err != nil {
				log.Panicf("mock.prepare.mysql.userx.error:%+v", err)
			}

			if initData {
				for i := 100; i < 108; i++ {
					sql := fmt.Sprintf(`insert into %s.%s values("%d", "%d","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","N","","","","",0,0,0,0,"mysql_native_password","*THISISNOTAVALIDPASSWORDTHATCANBEUSEDHERE","N","2017-06-22 17:37:18",NULL,"Y")`, shift.cfg.ToDatabase, shift.cfg.ToTable, i, i)
					if _, err := toConn.Execute(sql); err != nil {
						log.Panicf("mock.shift.prepare.datas.error:%+v", err)
					}
				}
			}
		}
	}

	// Prepare tables.
	{
		if err := shift.prepareTable(); err != nil {
			log.Panicf("mock.shift.prepare.table.error:%+v", err)
		}
	}

	// Prepare canal.
	{
		if err := shift.prepareCanal(); err != nil {
			log.Panicf("mock.shift.prepare.canal.error:%+v", err)
		}
		time.Sleep(time.Millisecond * 100)
	}

	// Prepare nearcheck.
	{
		if err := shift.behindsCheckStart(); err != nil {
			log.Panicf("mock.shift.behinds.check.error:%+v", err)
		}
	}
	return shift, func() {
		shift.Close()
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		h.Shutdown(ctx)
		time.Sleep(time.Millisecond * 100)
	}
}

func MockShift(log *xlog.Log, hasPK bool) (*Shift, func()) {
	return mockShift(log, mockCfg, hasPK, false, mockRadonReadonly, mockRadonShift, mockRadonThrottle)
}

func MockShiftWithCleanup(log *xlog.Log, hasPK bool) (*Shift, func()) {
	mockCfg.Cleanup = true
	return mockShift(log, mockCfg, hasPK, false, mockRadonReadonly, mockRadonShift, mockRadonThrottle)
}

func MockShiftWithData(log *xlog.Log, hasPK bool) (*Shift, func()) {
	return mockShift(log, mockCfg, hasPK, true, mockRadonReadonly, mockRadonShift, mockRadonThrottle)
}

func MockShiftXa(log *xlog.Log, hasPK bool) (*Shift, func()) {
	return mockShift(log, mockCfgXa, hasPK, false, mockRadonReadonly, mockRadonShift, mockRadonThrottle)
}

func MockShiftDDL(log *xlog.Log, hasPK bool) (*Shift, func()) {
	return mockShift(log, mockCfgDDL, hasPK, false, mockRadonReadonly, mockRadonShift, mockRadonThrottle)
}

func MockShiftMysqlTable(log *xlog.Log, hasPK bool) (*Shift, func()) {
	return mockShift(log, mockCfgMysql, hasPK, false, mockRadonReadonly, mockRadonShift, mockRadonThrottle)
}

func MockShiftMysqlTableWithData(log *xlog.Log, hasPK bool) (*Shift, func()) {
	return mockShift(log, mockCfgMysql, hasPK, true, mockRadonReadonly, mockRadonShift, mockRadonThrottle)
}

func MockShiftWithRadonReadonlyError(log *xlog.Log, hasPK bool) (*Shift, func()) {
	return mockShift(log, mockCfg, false, false, mockRadonReadonlyError, mockRadonShift, mockRadonThrottle)
}

func MockShiftWithRadonShardRuleError(log *xlog.Log, hasPK bool) (*Shift, func()) {
	return mockShift(log, mockCfg, false, false, mockRadonReadonly, mockRadonShiftError, mockRadonThrottle)
}

// RESTful api.
type mockHandler func(log *xlog.Log) rest.HandlerFunc

func mockHttp(log *xlog.Log, port int, readonly mockHandler, shardshift mockHandler, throttle mockHandler) *http.Server {
	httpAddr := fmt.Sprintf(":%d", port)
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Put("/v1/radon/readonly", readonly(log)),
		rest.Put("/v1/radon/throttle", throttle(log)),
		rest.Post("/v1/shard/shift", shardshift(log)),
	)
	if err != nil {
		log.Panicf("mock.shift.rest.make.router.error:%+v", err)
	}
	api.SetApp(router)
	handlers := api.MakeHandler()
	h := &http.Server{Addr: httpAddr, Handler: handlers}
	go func() {
		if err := h.ListenAndServe(); err != nil {
			log.Error("mock.shift.rest.error:%+v", err)
			return
		}
	}()
	time.Sleep(time.Millisecond * 100)
	return h
}

var readonlyLast bool

type readonlyParams struct {
	ReadOnly bool `json:"readonly"`
}

func mockRadonReadonly(log *xlog.Log) rest.HandlerFunc {
	f := func(w rest.ResponseWriter, r *rest.Request) {
		p := readonlyParams{}
		r.DecodeJsonPayload(&p)
		readonlyLast = p.ReadOnly
		log.Info("mock.api.radon.readonly.call.req:%+v", p)
	}
	return f
}

var throttleLast int

type throttleParams struct {
	Limits int `json:"limits"`
}

func mockRadonThrottle(log *xlog.Log) rest.HandlerFunc {
	f := func(w rest.ResponseWriter, r *rest.Request) {
		p := throttleParams{}
		r.DecodeJsonPayload(&p)
		throttleLast = p.Limits
		log.Info("mock.api.radon.throttle.call.req:%+v", p)
	}
	return f
}

func mockRadonShift(log *xlog.Log) rest.HandlerFunc {
	f := func(w rest.ResponseWriter, r *rest.Request) {
		log.Info("mock.api.radon.rule.call")
	}
	return f
}

func mockRadonReadonlyError(log *xlog.Log) rest.HandlerFunc {
	f := func(w rest.ResponseWriter, r *rest.Request) {
		log.Info("mock.api.readonly.error.call")
		readonlyLast = false
		rest.Error(w, "mock.api.readonly.error", http.StatusInternalServerError)
	}
	return f
}

func mockRadonShiftError(log *xlog.Log) rest.HandlerFunc {
	f := func(w rest.ResponseWriter, r *rest.Request) {
		log.Info("mock.api.shift.error.call")
		rest.Error(w, "mock.api.shift.error", http.StatusInternalServerError)
	}
	return f
}

func mockPanicMe(log *xlog.Log, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log.Info("mock.panicme.fired, msg:%s", msg)
	panic(1)
}

func mockRecoverPanicMe(log *xlog.Log, format string, v ...interface{}) {
	defer func() {
		if x := recover(); x != nil {
			msg := fmt.Sprintf(format, v...)
			log.Info("mock.panicme.fired, msg:%s", msg)
		}
	}()
	panic(1)
}