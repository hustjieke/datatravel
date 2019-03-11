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

	"github.com/siddontang/go-mysql/client"
)

// Cleanup used to clean up the tables on the from who has shifted,
// Or cleanup the to tables who half shifted.
// This func must be called after canal closed, otherwise it maybe replicated by canal.
func (shift *Shift) Cleanup() {
	log := shift.log

	// Set throttle to unlimits.
	if err := shift.setRadonThrottle(0); err != nil {
		log.Error("shift.cleanup.set.radon.throttle.error:%+v", err)
	}

	// Set readonly to false.
	if err := shift.setRadonReadOnly(false); err != nil {
		log.Error("shift.cleanup.set.radon.readonly.error:%+v", err)
	}

	// Cleanup.
	if shift.cfg.Cleanup {
		if shift.allDone {
			// Func cleanupFrom is use for shift big table between
			// diffrent backends in RadonDB, it may be used in the future.
			// shift.cleanupFrom()

			// Now for safety of src data, we do not cleanup
			log.Info("datatravel.all.done")
		} else {
			switch shift.cfg.ToFlavor {
			case config.ToMySQLFlavor:
			case config.ToMariaDBFlavor:
				// For Mysql and MariaDB, we drop all tables
				shift.cleanupToByDrop()
			case config.ToRadonDBFlavor:
				// For RadonDB, we just truncate the tables as the tables
				// in RadonDB are all created by users before migration
				shift.cleanupToByTruncate()
			}
		}
	}
}

// cleanupFrom used to cleanup the table on from.
// This func was called after shift succuess with
// cfg.Cleanup=true and used when migration table in RadonDB.
func (shift *Shift) cleanupFrom() {
	log := shift.log
	cfg := shift.cfg

	log.Info("shift.cleanup.from.table[%s.%s]...", cfg.FromDatabase, cfg.FromTable)
	if _, isSystem := sysDatabases[strings.ToLower(cfg.FromDatabase)]; !isSystem {
		from, err := client.Connect(cfg.From, cfg.FromUser, cfg.FromPassword, "")
		if err != nil {
			shift.panicMe("shift.cleanup.connection.error:%+v", err)
		}
		defer from.Close()

		sql := fmt.Sprintf("drop table `%s`.`%s`", cfg.FromDatabase, cfg.FromTable)
		if _, err := from.Execute(sql); err != nil {
			shift.panicMe("shift.execute.sql[%s].error:%+v", sql, err)
		}
	} else {
		log.Info("shift.table.is.system.cleanup.skip...")
	}
	log.Info("shift.cleanup.from.table.done...")
}

// cleanupToByDrop used to cleanup the tables on to.
// This func was called when travel data failed.
// used for MySQL/MariaDB-->MySQL/MariaDB
func (shift *Shift) cleanupToByDrop() {
	log := shift.log
	cfg := shift.cfg

	log.Info("shift.cleanup.to[%s/%s]...", cfg.ToDatabase, cfg.ToTable)
	if _, isSystem := sysDatabases[strings.ToLower(cfg.FromDatabase)]; !isSystem {
		to, err := client.Connect(cfg.To, cfg.ToUser, cfg.ToPassword, "")
		if err != nil {
			log.Error("shift.cleanup.to.connect.error:%+v", err)
			return
		}
		defer to.Close()

		sql := fmt.Sprintf("drop table `%s`.`%s`", cfg.ToDatabase, cfg.ToTable)
		if _, err := to.Execute(sql); err != nil {
			log.Error("shift.cleanup.to.execute[%s].error:%+v", sql, err)
			return
		}
	} else {
		log.Info("shift.table.is.system.cleanup.skip...")
	}
	log.Info("shift.cleanup.to.done...")
}

// cleanupToByTruncate used to truncate tables in
// RadonDB when we travel data failed.
// used for MySQL/MariaDB-->RadonDB
func (shift *Shift) cleanupToByTruncate() {
	log := shift.log
	cfg := shift.cfg

	for db, tbls := range cfg.DBTablesMaps {
		for _, tbl := range tbls {
			log.Info("datatravel.cleanto.by.truncate[%s/%s]...", db, tbl)
			to, err := client.Connect(cfg.To, cfg.ToUser, cfg.ToPassword, "")
			if err != nil {
				log.Error("datatravel.cleanup.to.connect.error:%+v", err)
				return
			}
			defer to.Close()

			// e.g. truncate sbtest1.benchyou0;
			sql := fmt.Sprintf("truncate `%s`.`%s`", db, tbl)
			if _, err := to.Execute(sql); err != nil {
				log.Error("datatravel.truncate.to.execute[%s].error:%+v", sql, err)
				return
			}
		}
	}
	log.Info("datatravel.truncate.to.done...")
}
