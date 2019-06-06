/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package shift

import (
	"config"
	"fmt"
	"strings"
	"time"
	"xlog"
)

var (

	// Config for normal shift.
	mockCfg = &config.Config{
		SetGlobalReadLock:  true,
		MetaDir:            "/tmp",
		FkCheck:            false,
		MaxAllowedPacketMB: 16,
		From:               "127.0.0.1:3306",
		FromUser:           "root",
		FromDatabase:       "shift_test_from",
		FromTable:          "t1",

		To:         "127.0.0.1:3307",
		ToUser:     "root",
		ToDatabase: "shift_test_to",
		ToTable:    "t1",

		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		Checksum:  true,
	}

	// Config for system (mysql) shift.
	mockCfgMysql = &config.Config{
		SetGlobalReadLock:  true,
		MetaDir:            "/tmp",
		FkCheck:            false,
		MaxAllowedPacketMB: 16,
		From:               "127.0.0.1:3306",
		FromUser:           "root",
		FromDatabase:       "mysql",
		FromTable:          "user",

		To:         "127.0.0.1:3307",
		ToUser:     "root",
		ToDatabase: "mysql",
		ToTable:    "userx",

		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		Checksum:  false,
	}

	// Config for xa shift.
	mockCfgXa = &config.Config{
		SetGlobalReadLock:  true,
		MetaDir:            "/tmp",
		FkCheck:            false,
		MaxAllowedPacketMB: 16,
		From:               "127.0.0.1:3306",
		FromUser:           "root",
		FromDatabase:       "shift_test_from",
		FromTable:          "t1",

		To:         "127.0.0.1:3307",
		ToUser:     "root",
		ToDatabase: "shift_test_to",
		ToTable:    "t1",

		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		Checksum:  true,
	}

	// Config for ddl shift.
	mockCfgDDL = &config.Config{
		SetGlobalReadLock:  true,
		MetaDir:            "/tmp",
		FkCheck:            false,
		MaxAllowedPacketMB: 16,
		From:               "127.0.0.1:3306",
		FromUser:           "root",
		FromDatabase:       "shift_test_from",
		FromTable:          "t1",

		To:         "127.0.0.1:3307",
		ToUser:     "root",
		ToDatabase: "shift_test_to",
		ToTable:    "t1",

		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		Checksum:  false,
	}
)

func mockShift(log *xlog.Log, cfg *config.Config, hasPK bool, initData bool) (*Shift, func()) {
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
		if fromConn == nil {
			panic("shift.mock.get.from.conn.nil.error")
		}
		toConn := shift.toPool.Get()
		defer shift.toPool.Put(toConn)
		if toConn == nil {
			panic("shift.mock.get.to.conn.nil.error")
		}

		if _, isSystem := sysDatabases[strings.ToLower(cfg.FromDatabase)]; !isSystem {
			// Cleanup To database first, datatravel shift FromDatabase to To.
			{
				sql := fmt.Sprintf("drop database if exists `%s`", cfg.FromDatabase)
				if _, err := toConn.Execute(sql); err != nil {
					log.Panicf("mock.shift.drop.to.table.error:%+v", err)
				}
			}

			// Cleanup From database.
			{
				sql := fmt.Sprintf("drop database if exists `%s`", cfg.FromDatabase)
				if _, err := fromConn.Execute(sql); err != nil {
					log.Panicf("mock.shift.drop.from.database.error:%+v", err)
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
				sql = fmt.Sprintf("create table `%s`.`%s`(a int, b int, c varchar(200), d DOUBLE NULL DEFAULT NULL, e json DEFAULT NULL, f INT UNSIGNED DEFAULT NULL, g BIGINT DEFAULT NULL, h BIGINT UNSIGNED DEFAULT NULL, i TINYINT NULL, j TINYINT UNSIGNED DEFAULT NULL, k SMALLINT DEFAULT NULL, l SMALLINT UNSIGNED DEFAULT NULL, m MEDIUMINT DEFAULT NULL, n INT UNSIGNED DEFAULT NULL)", cfg.FromDatabase, cfg.FromTable)
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
			// Cleanup To table first
			{
				sql := fmt.Sprintf("drop table if exists `%s`.`%s`", cfg.FromDatabase, cfg.ToTable)
				if _, err := toConn.Execute(sql); err != nil {
					log.Panicf("mock.shift.prepare.table.error:%+v", err)
				}
			}
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
		// Unlock global read lock tables
		shift.unLockTables()

		// Cleanup shift_test_from before end
		toConn := shift.toPool.Get()
		defer shift.toPool.Put(toConn)
		{
			sql := fmt.Sprintf("drop database if exists `%s`", mockCfg.FromDatabase)
			if _, err := toConn.Execute(sql); err != nil {
				log.Panicf("mock.shift.prepare.table.error:%+v", err)
			}
		}
		shift.Close()
		time.Sleep(time.Millisecond * 100)
	}
}

func MockShift(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfg.ToFlavor = flavor
	mockCfg.DBTablesMaps = make(map[string][]string) // init map
	mockCfg.Databases = make([]string, 0, 0)         // init dbs
	mockCfg.FromRows = 0
	mockCfg.ToRows = 0
	return mockShift(log, mockCfg, hasPK, false)
}

func MockShiftWithCleanup(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfg.ToFlavor = flavor
	mockCfg.DBTablesMaps = make(map[string][]string) // init map
	mockCfg.Databases = make([]string, 0, 0)         // init dbs
	mockCfg.FromRows = 0
	mockCfg.ToRows = 0
	return mockShift(log, mockCfg, hasPK, false)
}

func MockShiftWithData(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfg.ToFlavor = flavor
	mockCfg.DBTablesMaps = make(map[string][]string) // init map
	mockCfg.Databases = make([]string, 0, 0)         // init dbs
	mockCfg.FromRows = 0
	mockCfg.ToRows = 0
	return mockShift(log, mockCfg, hasPK, true)
}

func MockShiftXa(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfgXa.ToFlavor = flavor
	mockCfgXa.DBTablesMaps = make(map[string][]string) // init map
	mockCfgXa.Databases = make([]string, 0, 0)         // init dbs
	mockCfgXa.FromRows = 0
	mockCfgXa.ToRows = 0
	return mockShift(log, mockCfgXa, hasPK, false)
}

func MockShiftDDL(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfgDDL.ToFlavor = flavor
	mockCfgDDL.DBTablesMaps = make(map[string][]string) // init map
	mockCfgDDL.Databases = make([]string, 0, 0)         // init dbs
	mockCfgDDL.FromRows = 0
	mockCfgDDL.ToRows = 0
	return mockShift(log, mockCfgDDL, hasPK, false)
}

func MockShiftMysqlTable(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfgMysql.ToFlavor = flavor
	mockCfgMysql.DBTablesMaps = make(map[string][]string) // init map
	mockCfgMysql.Databases = make([]string, 0, 0)         // init dbs
	mockCfgMysql.FromRows = 0
	mockCfgMysql.ToRows = 0
	return mockShift(log, mockCfgMysql, hasPK, false)
}

func MockShiftMysqlTableWithData(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfgMysql.ToFlavor = flavor
	mockCfgMysql.DBTablesMaps = make(map[string][]string) // init map
	mockCfgMysql.Databases = make([]string, 0, 0)         // init dbs
	mockCfgMysql.FromRows = 0
	mockCfgMysql.ToRows = 0
	return mockShift(log, mockCfgMysql, hasPK, true)
}

func MockShiftWithRadonReadonlyError(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfg.ToFlavor = flavor
	mockCfg.DBTablesMaps = make(map[string][]string) // init map
	mockCfg.Databases = make([]string, 0, 0)         // init dbs
	mockCfg.FromRows = 0
	mockCfg.ToRows = 0
	return mockShift(log, mockCfg, false, false)
}

func MockShiftWithRadonShardRuleError(log *xlog.Log, hasPK bool, flavor string) (*Shift, func()) {
	mockCfg.ToFlavor = flavor
	mockCfg.DBTablesMaps = make(map[string][]string) // init map
	mockCfg.Databases = make([]string, 0, 0)         // init dbs
	mockCfg.FromRows = 0
	mockCfg.ToRows = 0
	return mockShift(log, mockCfg, false, false)
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
