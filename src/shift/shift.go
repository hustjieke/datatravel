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
	"sync"
	"time"

	"config"
	"xlog"

	"github.com/pingcap/errors"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
)

type Shift struct {
	log           *xlog.Log
	cfg           *config.Config
	toPool        *Pool
	fromPool      *Pool
	canal         *canal.Canal
	behindsTicker *time.Ticker
	done          chan bool
	handler       *EventHandler
	allDone       bool
	panicHandler  func(log *xlog.Log, format string, v ...interface{})

	// wg used for check when travel data done
	wg sync.WaitGroup
	// used for set flush tables with read lock and unlock tables;
	readLockConn *client.Conn
	// true: datatravel running normal; false: get some exception
	canalStatus bool
}

func NewShift(log *xlog.Log, cfg *config.Config) *Shift {
	log.Info("shift.cfg:%#v", cfg)
	return &Shift{
		log:           log,
		cfg:           cfg,
		done:          make(chan bool),
		behindsTicker: time.NewTicker(time.Duration(5000) * time.Millisecond),
		panicHandler:  logPanicHandler,
		canalStatus:   true,
	}
}

func (shift *Shift) CheckCanalStatus() bool {
	return shift.canalStatus
}

func (shift *Shift) SetCanalStatus(b bool) {
	shift.canalStatus = b
}

func (shift *Shift) prepareConnection() error {
	log := shift.log
	cfg := shift.cfg

	fromPool, err := NewPool(log, 16, cfg.From, cfg.FromUser, cfg.FromPassword, cfg.FkCheck)
	if err != nil {
		log.Error("shift.start.from.connection.pool.error:%+v", err)
		return err
	}
	shift.fromPool = fromPool
	log.Info("shift.[%s].connection.done...", cfg.From)

	toPool, err := NewPool(log, cfg.Threads, cfg.To, cfg.ToUser, cfg.ToPassword, cfg.FkCheck)
	if err != nil {
		log.Error("shift.start.to.connection.pool.error:%+v", err)
		return err
	}
	shift.toPool = toPool
	log.Info("shift.[%s].connection.done...", cfg.To)
	log.Info("shift.prepare.connections.done...")
	return nil
}

// Check if tbl is in tables
func (shift *Shift) isIncludeTbl(tbls []string, tbl string) bool {
	for _, v := range tbls {
		if v == tbl {
			return true
		}
	}
	return false
}

func (shift *Shift) prepareTable() error {
	log := shift.log
	cfg := shift.cfg

	// From connection.
	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)
	if fromConn == nil {
		shift.panicMe("shift.get.from.conn.nil.error")
	}

	// To connection.
	toConn := shift.toPool.Get()
	defer shift.toPool.Put(toConn)
	if toConn == nil {
		shift.panicMe("shift.get.to.conn.nil.error")
	}

	// Get databases if len(dbs) == 0
	if len(cfg.Databases) == 0 {
		if cfg.TableDB != "" {
			cfg.Databases = append(cfg.Databases, cfg.TableDB)
		} else {
			log.Info("shift.get.database...")
			sql := "show databases;"
			r, err := fromConn.Execute(sql)
			if err != nil {
				log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
				return err
			}
			for i := 0; i < r.RowNumber(); i++ {
				str, _ := r.GetString(i, 0)
				if _, isSystem := sysDatabases[strings.ToLower(str)]; !isSystem {
					cfg.Databases = append(cfg.Databases, str)
				}
			}
			if len(cfg.Databases) == 0 {
				return errors.New("no.database.to.shift")
			}
		}
	}

	for _, db := range cfg.Databases {
		// Prepare database, check the database is not system database and create them.
		log.Info("shift.prepare.database[%s]...", db)
		sql := fmt.Sprintf("select * from information_schema.tables where table_schema = '%s' limit 1", db)
		r, err := toConn.Execute(sql)
		if err != nil {
			log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
			return err
		}

		if r.RowNumber() == 0 {
			sql := fmt.Sprintf("create database if not exists `%s`", db)
			if _, err := toConn.Execute(sql); err != nil {
				log.Error("shift.create.database.sql[%s].error:%+v", sql, err)
				return err
			}
			log.Info("shift.prepare.database.done...")
		} else {
			log.Info("shift.database.exists...")
		}

		// Get tables
		sql = fmt.Sprintf("use `%s`", db)
		r, err = fromConn.Execute(sql)
		if err != nil {
			log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
			return err
		}

		sql = fmt.Sprintf("show table status")
		r, err = fromConn.Execute(sql)
		if err != nil {
			log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
			return err
		}

		var tables []string
		tblName := "Name"
		tblRows := "Rows"
		for i := 0; i < r.RowNumber(); i++ {
			tbl, _ := r.GetStringByName(i, tblName)
			// In case --table-db, --tables
			if len(cfg.Tables) != 0 {
				if shift.isIncludeTbl(cfg.Tables, tbl) {
					tables = append(tables, tbl)
					rows, _ := r.GetUintByName(i, tblRows)
					cfg.FromRows += rows
					log.Info("shift.show.[%s.%s].rows:%+v", db, tbl, rows)
				} else {
					continue
				}
			} else {
				// In case --databases or --table-db
				tables = append(tables, tbl)
				rows, _ := r.GetUintByName(i, tblRows)
				cfg.FromRows += rows
				log.Info("shift.show.[%s.%s].rows:%+v", db, tbl, rows)
			}
		}
		cfg.DBTablesMaps[db] = tables
		if len(tables) == 0 {
			log.Error("shift.check.database.[%+v].no.tables", db)
			continue // don`t need return err
		}

		// prepare tables
		for _, tbl := range tables {
			log.Info("shift.prepare.table[%s/%s]...", db, tbl)
			sql = fmt.Sprintf("show create table `%s`.`%s`", db, tbl)
			r, err = fromConn.Execute(sql)
			if err != nil {
				log.Error("shift.show.[%s].create.table.sql[%s].error:%+v", cfg.From, sql, err)
				return err
			}
			sql, err = r.GetString(0, 1)
			if err != nil {
				log.Error("shift.show.[%s].create.table.get.error:%+v", cfg.From, err)
				return err
			}
			sql = strings.Replace(sql, fmt.Sprintf("CREATE TABLE `%s`", tbl), fmt.Sprintf("CREATE TABLE `%s`.`%s`", db, tbl), 1)
			if _, err := toConn.Execute(sql); err != nil {
				log.Error("shift.create.[%s].table.sql[%s].error:%+v", cfg.To, sql, err)
				return err
			}
			log.Info("shift.prepare.table.done...")
		}
	}

	return nil
}

// Used for flavor RadonDB
func (shift *Shift) checkTableExistForRadonDB() error {
	log := shift.log
	cfg := shift.cfg

	// From connection.
	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)

	// To connection.
	toConn := shift.toPool.Get()
	defer shift.toPool.Put(toConn)

	// Get databases if len(dbs) == 0
	if len(cfg.Databases) == 0 {
		if cfg.TableDB != "" {
			cfg.Databases = append(cfg.Databases, cfg.TableDB)
		} else {
			// Get all dbs except system dbs
			log.Info("shift.get.database...")
			sql := "show databases;"
			r, err := fromConn.Execute(sql)
			if err != nil {
				log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
				return err
			}
			for i := 0; i < r.RowNumber(); i++ {
				str, _ := r.GetString(i, 0)
				if _, isSystem := sysDatabases[strings.ToLower(str)]; !isSystem {
					cfg.Databases = append(cfg.Databases, str)
				}
			}
			if len(cfg.Databases) == 0 {
				return errors.New("no.database.to.shift")
			}
		}
	}

	// Get databases
	/*
		log.Info("shift.get.database...")
		sql := "show databases;"
		r, err := fromConn.Execute(sql)
		if err != nil {
			log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
			return err
		}
		for i := 0; i < r.RowNumber(); i++ {
			str, _ := r.GetString(i, 0)
			if _, isSystem := sysDatabases[strings.ToLower(str)]; !isSystem {
				cfg.Databases = append(cfg.Databases, str)
			}
		}
		if len(cfg.Databases) == 0 {
			return errors.New("no.database.to.shift")
		}
	*/

	// Check if db.table exist
	for _, db := range cfg.Databases {
		// Prepare database, check the database is not system database and create them.
		log.Info("shift.check.if.database[%s].exist...", db)
		sql := fmt.Sprintf("use %s", db)
		r, err := toConn.Execute(sql)
		if err != nil {
			log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
			return err
		} else {
			log.Info("shift.check.database.ok[%s]", db)
		}

		// Get from tables
		sql = fmt.Sprintf("use `%s`", db)
		r, err = fromConn.Execute(sql)
		if err != nil {
			log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
			return err
		}

		sql = fmt.Sprintf("show table status")
		r, err = fromConn.Execute(sql)
		if err != nil {
			log.Error("shift.check.database.sql[%s].error:%+v", sql, err)
			return err
		}

		var tables []string
		tblName := "Name"
		tblRows := "Rows"
		for i := 0; i < r.RowNumber(); i++ {
			tbl, _ := r.GetStringByName(i, tblName)
			// In case --table-db, --tables
			if len(cfg.Tables) != 0 {
				if shift.isIncludeTbl(cfg.Tables, tbl) {
					tables = append(tables, tbl)
					rows, _ := r.GetUintByName(i, tblRows)
					cfg.FromRows += rows
					log.Info("shift.show.[%s.%s].rows:%+v", db, tbl, rows)
				} else {
					continue
				}
			} else {
				// In case --databases or --table-db
				tables = append(tables, tbl)
				rows, _ := r.GetUintByName(i, tblRows)
				cfg.FromRows += rows
				log.Info("shift.show.[%s.%s].rows:%+v", db, tbl, rows)
			}
		}
		cfg.DBTablesMaps[db] = tables
		if len(tables) == 0 {
			log.Error("shift.check.database.[%+v].no.tables", db)
			continue // don`t need return err
		}

		// Check if RadonDB`s tables and fields are consistent with from
		for _, tbl := range tables {
			log.Info("shift.check.if.radondb.table[%s/%s].exist...", db, tbl)
			sql = fmt.Sprintf("show create table `%s`.`%s`", db, tbl)
			_, errFrom := fromConn.Execute(sql)
			if errFrom != nil {
				log.Error("shift.show.[%s].sql.[%s].error:%+v", cfg.From, sql, errFrom)
				return errFrom
			}
		}
		// TODO show columns from db.tbl resp error
		/*
		   for _, tbl := range tables {
		       log.Info("shift.check.radondb.table[%s/%s]...", db, tbl)
		       sql = fmt.Sprintf("show columns from `%s`.`%s`", db, tbl)
		       rFrom, errFrom := fromConn.Execute(sql)
		       if errFrom != nil {
		           log.Error("shift.show.[%s].sql.[%s].error:%+v", cfg.From, sql, errFrom)
		           return errFrom
		       }
		       rTo, errTo := toConn.Execute(sql)
		       if errTo != nil {
		           log.Error("shift.show.[%s].sql.[%s].error:%+v", cfg.From, sql, errTo)
		           return errTo
		       }

		       if rFrom.ColumnNumber() == rTo.ColumnNumber() {
		           for idx, field := range rFrom.Fields {
		               if (string(field.Name) != string(rTo.Fields[idx].Name)) ||
		                   (field.Type != rTo.Fields[idx].Type) {
		                   errmsg := fmt.Sprintf("shift.check.db[%s],table[%s].columns.fields.not.consistent.with.from", db, tbl)
		                   return errors.New(errmsg)
		               }
		           }
		       } else {
		           errmsg := fmt.Sprintf("shift.check.db[%s],table[%s].columns.size.not.consistent.with.from", db, tbl)
		           return errors.New(errmsg)
		       }
		   }
		*/
	}
	return nil
}

// TODO(gry), check for mysql 5.6--datetime, some bug exist in checksum type datetime
// Now this func only used for travel mysql-->RadonDB
func (shift *Shift) getAutoIncTable() error {
	log := shift.log
	cfg := shift.cfg

	// From connection.
	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)
	if fromConn == nil {
		shift.panicMe("shift.get.from.conn.nil.error")
	}

	log.Info("shift.get.auto.inc.tables")
	for db, tbls := range cfg.DBTablesMaps {
		for _, tbl := range tbls {
			sql := fmt.Sprintf("show create table `%s`.`%s`", db, tbl)
			r, err := fromConn.Execute(sql)
			if err != nil {
				log.Error("shift.show.[%s].create.table.sql[%s].error:%+v", cfg.From, sql, err)
				return err
			}
			createSql, err := r.GetString(0, 1)
			if strings.Contains(createSql, "AUTO_INCREMENT") {
				dbTbl := fmt.Sprintf("`%s`.`%s`", db, tbl)
				cfg.AutoIncTable[dbTbl] = true
			}
		}
	}

	return nil
}

// Only used for mysql-->radondb
func (shift *Shift) getToChecksumInfos() error {
	log := shift.log
	cfg := shift.cfg

	// To connection.
	toConn := shift.toPool.Get()
	defer shift.toPool.Put(toConn)
	if toConn == nil {
		shift.panicMe("shift.get.to.conn.nil.error")
	}

	log.Info("shift.get.checksum.and.count.infos.from.to")
	for db, tbls := range cfg.DBTablesMaps {
		for _, tbl := range tbls {
			dbTbl := fmt.Sprintf("`%s`.`%s`", db, tbl)
			if cfg.AutoIncTable[dbTbl] {
				sql := fmt.Sprintf("use `%s`", db)
				r, err := toConn.Execute(sql)
				if err != nil {
					log.Error("shift.select.[%s].count.table.sql[%s].error:%+v", cfg.From, sql, err)
					return err
				}
				sql = fmt.Sprintf("select count(*) from `%s`.`%s`", db, tbl)
				r, err = toConn.Execute(sql)
				if err != nil {
					log.Error("shift.select.[%s].count.table.sql[%s].error:%+v", cfg.From, sql, err)
					return err
				}

				v, err := r.GetUint(0, 0)
				if err != nil {
					shift.panicMe("shift.get.%s.table[%s.%s].count.error:%+v", "to", db, tbl, err)
				}
				cfg.ToTblsRowsBefore[dbTbl] = v
			} else {
				sql := fmt.Sprintf("checksum table `%s`.`%s`", db, tbl)
				r, err := toConn.Execute(sql)
				if err != nil {
					log.Error("shift.checksum.[%s].table.sql[%s].error:%+v", cfg.From, sql, err)
					return err
				}
				v, err := r.GetUint(0, 1)
				if err != nil {
					shift.panicMe("shift.get.%s.table[%s.%s].checksum.error:%+v", "to", db, tbl, err)
				}
				log.Info("gry----v:", v)
				cfg.ToTblsChecksumBefore[dbTbl] = uint32(v)
			}
		}
	}

	return nil
}

// Check if db.table has auto_increment column.
func (shift *Shift) containAutoIncCol(db, tbl string) bool {
	dbTbl := fmt.Sprintf("`%s`.`%s`", db, tbl)
	return shift.cfg.AutoIncTable[dbTbl]
}

// Only used for mysql-->radondb
// If we fisrt do db1 db2 ... to one db, the ToRows should be 0.
func (shift *Shift) setIsNotFirstTime() {
	if shift.cfg.ToRows != 0 {
		shift.cfg.IsNotFisrtTime = true
	}
}

func (shift *Shift) prepareCanal() error {
	log := shift.log
	conf := shift.cfg
	log.Info("datatravel.cfg.after.prepare.table:%+v", conf)
	cfg := canal.NewDefaultConfig()
	cfg.Addr = conf.From
	cfg.User = conf.FromUser
	cfg.Password = conf.FromPassword
	cfg.Dump.ExecutionPath = conf.MySQLDump
	cfg.Dump.DiscardErr = false
	if len(conf.Tables) != 0 && conf.TableDB != "" {
		cfg.Dump.TableDB = conf.TableDB
		cfg.Dump.Tables = conf.Tables
	} else if len(conf.Tables) == 0 && conf.TableDB != "" {
		// In this case conf.Databases is eq to conf.TableDB
		cfg.Dump.Databases = conf.Databases
	} else {
		// like db1 [db2 db3 ....]
		cfg.Dump.Databases = conf.Databases
	}
	cfg.Dump.MaxAllowedPacketMB = conf.MaxAllowedPacketMB

	// canal
	canal, err := canal.NewCanal(cfg)
	if err != nil {
		log.Error("shift.canal.new.error:%+v", err)
		return err
	}

	handler := NewEventHandler(log, shift)
	canal.SetEventHandler(handler)
	shift.handler = handler
	shift.canal = canal
	go func() {
		if err := canal.Run(); err != nil {
			log.Error("shift.canal.run.error:%+v", err)
			shift.SetCanalStatus(false)
		}
	}()
	log.Info("shift.prepare.canal.done...")
	return nil
}

/*
	mysql> checksum table sbtest.sbtest1;
	+----------------+-----------+
	| Table          | Checksum  |
	+----------------+-----------+
	| sbtest.sbtest1 | 410139351 |
	+----------------+-----------+
*/
// ChecksumTable ensure that FromTable and ToTable are consistent
func (shift *Shift) ChecksumTables() error {
	for db, tbls := range shift.cfg.DBTablesMaps {
		shift.wg.Add(1)
		go shift.checksumTables(db, tbls)
	}
	shift.wg.Wait()
	return nil
}

// TODO(gry) 后面最好还是分离成checksum和checkcount两个函数，清晰的很
func (shift *Shift) checksumTables(db string, tbls []string) error {
	log := shift.log
	var fromchecksum, tochecksum uint32
	var fromcount, tocount uint64

	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)
	if fromConn == nil {
		shift.panicMe("shift.get.from.conn.nil.error")
	}

	toConn := shift.toPool.Get()
	defer shift.toPool.Put(toConn)
	if toConn == nil {
		shift.panicMe("shift.get.to.conn.nil.error")
	}

	for _, tbl := range tbls {
		checksumFunc := func(t string, Conn *client.Conn, Database string, Table string, c chan uint64) {
			sql := fmt.Sprintf("checksum table %s.%s", Database, Table)
			r, err := Conn.Execute(sql)
			if err != nil {
				shift.panicMe("shift.checksum.%s.table[%s.%s].error:%+v", t, Database, Table, err)
			}

			v, err := r.GetUint(0, 1)
			if err != nil {
				shift.panicMe("shift.get.%s.table[%s.%s].checksum.error:%+v", Database, Table, err)
			}
			c <- v
		}

		// Only used for mysql-->radondb
		countTblFunc := func(t string, Conn *client.Conn, Database string, Table string, c chan uint64) {
			sql := fmt.Sprintf("select count(*) from %s.%s", Database, Table)
			r, err := Conn.Execute(sql)
			if err != nil {
				shift.panicMe("shift.count.%s.table[%s.%s].error:%+v", t, Database, Table, err)
			}

			v, err := r.GetUint(0, 0)
			if err != nil {
				shift.panicMe("shift.get.%s.table[%s.%s].count.error:%+v", t, Database, Table, err)
			}
			c <- v
		}

		fromchan := make(chan uint64, 1)
		tochan := make(chan uint64, 1)

		if shift.cfg.ToFlavor == config.ToRadonDBFlavor && shift.containAutoIncCol(db, tbl) &&
			shift.cfg.IsNotFisrtTime {
			// execute count func
			{
				go countTblFunc("from", fromConn, db, tbl, fromchan)
				go countTblFunc("to", toConn, db, tbl, tochan)
			}
			fromcount = <-fromchan
			tocount = <-tochan
			dbTbl := fmt.Sprintf("`%s`.`%s`", db, tbl)
			toRealCount := tocount - shift.cfg.ToTblsRowsBefore[dbTbl]
			if fromcount != toRealCount {
				err := fmt.Errorf("count not equivalent: from-table[%v.%v] count is %v, to-table[%v.%v] count is %v", db, tbl, fromcount, db, tbl, toRealCount)
				log.Error("shift.count.table.err:%+v", err)
				shift.wg.Done()
				shift.panicMe("shift.count.table.err:", err)
			}
		} else {
			// execute checksum func
			{
				go checksumFunc("from", fromConn, db, tbl, fromchan)
				go checksumFunc("to", toConn, db, tbl, tochan)
			}
			fromchecksum = uint32(<-fromchan)
			tochecksum = uint32(<-tochan)
			if shift.cfg.ToFlavor == config.ToRadonDBFlavor {
				dbTbl := fmt.Sprintf("`%s`.`%s`", db, tbl)
				fromchecksum += shift.cfg.ToTblsChecksumBefore[dbTbl]
			}
			if fromchecksum != tochecksum {
				err := fmt.Errorf("checksum not equivalent: from-table[%v.%v] checksum is %v, to-table[%v.%v] checksum is %v", db, tbl, fromchecksum, db, tbl, tochecksum)
				log.Error("shift.checksum.table.err:%+v", err)
				shift.wg.Done()
				shift.panicMe("shift.checksum.table.err:", err)
			}
			log.Info("shift.checksum.table.from[%v.%v, crc:%v].to[%v.%v, crc:%v].ok", db, tbl, fromchecksum, db, tbl, tochecksum)
		}
	}
	shift.wg.Done()
	return nil
}

/*
 mysql> show table status;
+-----------+--------+---------+------------+--------+----------------+
| Name      | Engine | Version | Row_format | Rows   | Avg_row_length |
+-----------+--------+---------+------------+--------+----------------+
| benchyou0 | InnoDB |      10 | Dynamic    |  95883 |            400 |
| benchyou1 | InnoDB |      10 | Dynamic    |  98833 |            388 |
| benchyou2 | InnoDB |      10 | Dynamic    |  91012 |            421 |
| benchyou3 | InnoDB |      10 | Dynamic    |  95399 |            402 |
| benchyou4 | InnoDB |      10 | Dynamic    |  90460 |            423 |
| benchyou5 | InnoDB |      10 | Dynamic    |  99157 |            386 |
| benchyou6 | InnoDB |      10 | Dynamic    |  92319 |            415 |
| benchyou7 | InnoDB |      10 | Dynamic    | 101641 |            377 |
| benchyou8 | InnoDB |      10 | Dynamic    |  96736 |            363 |
| benchyou9 | InnoDB |      10 | Dynamic    |  94729 |            415 |
+-----------+--------+---------+------------+--------+----------------+
*/
/* count dump progress */
const (
	secondsPerMinute = 60
	secondsPerHour   = secondsPerMinute * 60
	secondsSleep     = 30
)

func (shift *Shift) dumpProgress() error {
	go func(s *Shift) {
		cfg := s.cfg
		log := s.log

		sql := "show table status"
		toConn := s.toPool.Get()
		defer s.toPool.Put(toConn)
		var dumpTime uint64
		log.Info("cfg.FromRows when dump:%d", cfg.FromRows)

		progress := &config.TravelProgress{
			PositionBehinds: "not start yet!",
		}
		// First time to calculate rows and flag
		var baseCalRows uint64
		firstCalFlag := true
		for {
			// If first time we execute, skip to add dumpTime to get baseCalRows first
			// and then we really start to calculate dump progress
			if !firstCalFlag {
				dumpTime += secondsSleep
			}

			// Get rows from toConn
			cfg.ToRows = 0
			for j := 0; j < len(cfg.Databases); j++ {
				// Get tables
				sql = fmt.Sprintf("use `%s`", cfg.Databases[j])
				r, err := toConn.Execute(sql)
				if err != nil {
					log.Error("shift.use.database.sql[%s].error:%+v", sql, err)
					return
				}

				sql = fmt.Sprintf("show table status")
				r, err = toConn.Execute(sql)
				if err != nil {
					log.Error("shift.show.table.status[%s].error:%+v", sql, err)
				}
				if r.RowNumber() == 0 {
					log.Info("shift.check.database.[%+v].no.tables", cfg.Databases[j])
					continue // don`t need return err
				}

				tblRows := "Rows"
				for i := 0; i < r.RowNumber(); i++ {
					rows, _ := r.GetUintByName(i, tblRows)
					cfg.ToRows += rows
				}
			}
			// Store rows when execute loop first time
			if firstCalFlag {
				baseCalRows = cfg.ToRows
				firstCalFlag = false
				time.Sleep(secondsSleep * time.Second)
				continue
			}

			// Calculate remain time
			rowsInc := cfg.ToRows - baseCalRows
			avgRate := rowsInc / dumpTime
			// Unit: second
			remainTime := (cfg.FromRows - rowsInc) / avgRate
			seconds := remainTime % secondsPerMinute
			hours := (remainTime - seconds) / secondsPerHour
			minutes := (remainTime - seconds - hours*secondsPerHour) / secondsPerMinute

			// Calculate progress rate
			// If data is not so large, it may be happened that per > 100
			// as the result of "show table status" is estimate in MySQL
			if cfg.FromRows == 0 {
				cfg.FromRows = 1
			}
			log.Info("cfg.FromRows before cal:%+v", cfg.FromRows)
			per := uint((float64(cfg.ToRows) / float64(cfg.FromRows)) * 100)
			if per > 100 {
				per = 100
			}

			timeStr := fmt.Sprintf("%v%v %v%v", hours, "hours", minutes, "minutes")
			perStr := fmt.Sprintf("%v%v", per, "%")
			progress.DumpRemainTime = timeStr
			progress.DumpProgressRate = perStr
			config.UpdateTravelProgress(progress, cfg.MetaDir)
			log.Info("travel.progress%+v", progress)

			if per > 98 {
				log.Info("wait.dump.shift.dump.done.during.progress:%+v", perStr)
				<-s.canal.WaitDumpDone()
				progress.DumpProgressRate = fmt.Sprintf("%v", "100%")
				progress.DumpRemainTime = fmt.Sprintf("%v %v", "0hours", "0minutes")
				config.UpdateTravelProgress(progress, cfg.MetaDir)
				log.Info("travel.progress%+v", progress)
				break
			}

			time.Sleep(secondsSleep * time.Second)
		}
	}(shift)

	return nil
}

/*
   mysql> show master status;
   +------------------+-----------+--------------+------------------+------------------------------------------------+
   | File             | Position  | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set                              |
   +------------------+-----------+--------------+------------------+------------------------------------------------+
   | mysql-bin.000002 | 112107994 |              |                  | 4dc59763-5431-11e7-90cb-5254281e57de:1-2561361 |
   +------------------+-----------+--------------+------------------+------------------------------------------------+
*/
func (shift *Shift) masterPosition() *mysql.Position {
	position := &mysql.Position{}

	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)
	if fromConn == nil {
		shift.panicMe("shift.get.from.conn.nil.error")
	}

	sql := "show master status"
	r, err := fromConn.Execute(sql)
	if err != nil {
		shift.panicMe("shift.get.master[%s].postion.error:%+v", shift.cfg.From, err)
		return position
	}

	file, err := r.GetString(0, 0)
	if err != nil {
		shift.panicMe("shift.get.master[%s].file.error:%+v", shift.cfg.From, err)
		return position
	}

	pos, err := r.GetUint(0, 1)
	if err != nil {
		shift.panicMe("shift.get.master[%s].pos.error:%+v", shift.cfg.From, err)
		return position
	}
	position.Name = file
	position.Pos = uint32(pos)
	return position
}

// 1. check mysqldump worker done
// 2. check sync binlog pos
func (shift *Shift) behindsCheckStart() error {
	go func(s *Shift) {
		log := s.log
		log.Info("shift.dumping...")
		// If some error happened during dumping, wait dump will be still set dump done.
		<-s.canal.WaitDumpDone()
		// Wait dump worker done.
		log.Info("shift.wait.dumper.background.worker.again...")
		shift.handler.WaitWorkerDone()
		log.Info("shift.wait.dumper.background.worker.done...")
		progress := &config.TravelProgress{
			DumpProgressRate: "100%",
			DumpRemainTime:   "0",
		}

		for range s.behindsTicker.C {
			// If canal get something wrong during dumping or syncing data, we should log error
			if s.CheckCanalStatus() {
				// Get master and sync gtid
				if s.canal.SyncedGTIDSet() != nil {
					syncGtid := s.canal.SyncedGTIDSet().(*mysql.MysqlGTIDSet)
					progress.SynGTID = syncGtid.String()
				}
				if masterGtid, err := s.canal.GetMasterGTIDSet(); err != nil {
					log.Panic("error:%+v", err)
				} else {
					progress.MasterGTID = masterGtid.String()
				}

				masterPos := s.masterPosition()
				syncPos := s.canal.SyncedPosition()
				behinds := int(masterPos.Pos - syncPos.Pos)
				progress.PositionBehinds = fmt.Sprintf("%v", behinds)
				config.UpdateTravelProgress(progress, s.cfg.MetaDir)
				log.Info("travel.progress%+v", progress)
				if behinds <= shift.cfg.Behinds {
					shift.checkAndSetReadlock()
					progress.PositionBehinds = "0"
					config.UpdateTravelProgress(progress, s.cfg.MetaDir)
					log.Info("travel.progress%+v", progress)
					return
				}
			} else {
				log.Error("shift.canal.get.error.during.dump.or.sync")
			}
		}
	}(shift)
	return nil
}

// Start used to start canal and behinds ticker.
func (shift *Shift) Start() error {
	if err := shift.prepareConnection(); err != nil {
		return err
	}
	if shift.cfg.ToFlavor == config.ToRadonDBFlavor {
		if err := shift.checkTableExistForRadonDB(); err != nil {
			return err
		}
		if err := shift.getAutoIncTable(); err != nil {
			return err
		}
		if err := shift.getToChecksumInfos(); err != nil {
			return err
		}
		shift.setIsNotFirstTime()
	} else {
		if err := shift.prepareTable(); err != nil {
			return err
		}
	}
	if err := shift.prepareCanal(); err != nil {
		return err
	}
	if err := shift.dumpProgress(); err != nil {
		return err
	}
	if err := shift.behindsCheckStart(); err != nil {
		return err
	}
	return nil
}

// Close used to destroy all the resource.
func (shift *Shift) Close() {
	log := shift.log
	shift.behindsTicker.Stop()
	shift.canal.Close()

	// if we catch
	if !shift.allDone {
		shift.unLockTables()
		log.Info("datatravel.migrates.data.fail!")
	} else {
		log.Info("datatravel.migrates.all.data.success!")
	}
	shift.fromPool.Close()
	shift.toPool.Close()
}

func (shift *Shift) Done() chan bool {
	return shift.done
}

func (shift *Shift) panicMe(format string, v ...interface{}) {
	shift.Close()
	shift.panicHandler(shift.log, format, v)
}

func (shift *Shift) checkAndSetReadlock() {
	log := shift.log

	// 1. WaitUntilPos
	{
		masterPos := shift.masterPosition()
		log.Info("shift.wait.until.pos[%+v]...", masterPos)
		if err := shift.canal.WaitUntilPos(*masterPos, time.Hour*12); err != nil {
			shift.panicMe("shift.set.radon.wait.until.pos[%#v].error:%+v", masterPos, err)
			return
		}
		log.Info("shift.wait.until.pos.done...")
	}

	// 2. Set a global read lock on from
	{
		if shift.cfg.SetGlobalReadLock {
			log.Info("set.from.a.global.read.lock")
			shift.setGlobalReadLock()
		} else {
			log.Info("continue.travel.without.set.global.read.lock")
		}
	}

	// 3. Wait again
	{
		masterPos := shift.masterPosition()
		log.Info("shift.wait.until.pos[%+v]...", masterPos)
		if err := shift.canal.WaitUntilPos(*masterPos, time.Second*300); err != nil {
			shift.panicMe("shift.wait.until.pos[%+v].error:%+v", masterPos, err)
			return
		}
		log.Info("shift.wait.until.pos.done...")
	}

	// 4. Checksum table.
	if shift.cfg.Checksum {
		switch shift.cfg.ToFlavor {
		case config.ToMySQLFlavor, config.ToMariaDBFlavor, config.ToRadonDBFlavor:
			log.Info("shift.checksum.table...")
			if err := shift.ChecksumTables(); err != nil {
				shift.panicMe("shift.checksum.table.error:%+v", err)
				return
			}
			log.Info("shift.checksum.table.done...")
		default:
			log.Error("shift.checksum.not.support.flavor.:%+v", shift.cfg.ToFlavor)
		}
	}

	// 5. Good, we have all done.
	{
		shift.done <- true
		shift.allDone = true
		log.Info("shift.all.done...")
	}
}

// Func setGlobalReadLock is used to add a global read
// lock "from" when ( master position - syncer position) is less than 2048
func (shift *Shift) setGlobalReadLock() {
	sql1 := "FLUSH TABLES"
	sql2 := "FLUSH TABLES WITH READ LOCK"
	if shift.readLockConn = shift.fromPool.Get(); shift.readLockConn == nil {
		shift.panicMe("datatravel.get.conn.nil")
	}

	if _, err := shift.readLockConn.Execute(sql1); err != nil {
		shift.panicMe("datatravel.execute.master[%s].sql[%s].error:%+v", shift.cfg.From, sql1, err)
		return
	}
	if _, err := shift.readLockConn.Execute(sql2); err != nil {
		shift.panicMe("datatravel.execute.master[%s].sql[%s].error:%+v", shift.cfg.From, sql2, err)
		return
	}
}

// Func unLockTables is used to unlock tables on from.
// It will be called when travel data failed.
func (shift *Shift) unLockTables() {
	// if readLockConn is nil, the read lock will be released automatically
	if shift.readLockConn == nil {
		return
	} else {
		defer shift.fromPool.Put(shift.readLockConn)
		sql := "UNLOCK TABLES"

		if _, err := shift.readLockConn.Execute(sql); err != nil {
			shift.panicMe("datatravel.execute.master[%s].sql[%s].error:%+v", shift.cfg.From, sql, err)
			return
		}
	}
}
