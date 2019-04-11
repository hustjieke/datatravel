/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package main

import (
	"build"
	"config"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"shift"
	"syscall"
	"xlog"
)

var (
	toFlavor          = flag.String("to-flavor", "", "Destination db flavor, like mysql/mariadb/radondb")
	setGlobalReadLock = flag.Bool("set-global-read-lock", true, "Add a read lock when src MySQL data is going done")
	metaDir           = flag.String("meta-dir", "./datatravel-meta", "meta dir to store database meta data")

	from         = flag.String("from", "", "Source MySQL backend")
	fromUser     = flag.String("from-user", "", "MySQL user, must have replication privilege")
	fromPassword = flag.String("from-password", "", "MySQL user password")
	fromDatabase = flag.String("from-database", "", "Source database")
	fromTable    = flag.String("from-table", "", "Source table")

	to         = flag.String("to", "", "Destination MySQL backend")
	toUser     = flag.String("to-user", "", "MySQL user, must have replication privilege")
	toPassword = flag.String("to-password", "", "MySQL user password")
	toDatabase = flag.String("to-database", "", "Destination database")
	toTable    = flag.String("to-table", "", "Destination table")

	cleanup   = flag.Bool("cleanup", true, "Cleanup the from table after shifted(defaults true)")
	checksum  = flag.Bool("checksum", true, "Checksum the from table and to table after shifted(defaults true)")
	mysqlDump = flag.String("mysqldump", "mysqldump", "mysqldump path")
	threads   = flag.Int("threads", 16, "shift threads num(defaults 16)")
	behinds   = flag.Int("behinds", 2048, "seconds behinds num(default 2048)")
	radonURL  = flag.String("radonurl", "http://127.0.0.1:8080", "Radon RESTful api(defaults http://127.0.0.1:8080)")

	debug = flag.Bool("debug", false, "Set log to debug mode(defaults false)")
)

func check(log *xlog.Log) {
	if *toFlavor == "" || *from == "" || *fromUser == "" || *to == "" || *toUser == "" {
		log.Panic("usage: datatravel --to-flavor=[radondb/mariadb/mysql] --from=[host:port] --from-password=[password] --to=[host:port] --to-user=[user] --to-password=[password] --cleanup=[false|true]")
	}
}

func main() {
	log := xlog.NewStdLog(xlog.Level(xlog.INFO))
	runtime.GOMAXPROCS(runtime.NumCPU())

	build := build.GetInfo()
	fmt.Printf("datatravel:[%+v]\n", build)

	// flags.
	flag.Parse()

	// log.
	if *debug {
		log = xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	}
	check(log)
	fmt.Println(`
           IMPORTANT: Please check that the shift run completes successfully.
           At the end of a successful shift run prints "datatravel.migrates.all.data.success!".`)

	cfg := &config.Config{
		ToFlavor:          *toFlavor,
		SetGlobalReadLock: *setGlobalReadLock,
		MetaDir:           *metaDir,
		From:              *from,
		FromUser:          *fromUser,
		FromPassword:      *fromPassword,
		FromDatabase:      *fromDatabase,
		FromTable:         *fromTable,
		To:                *to,
		ToUser:            *toUser,
		ToPassword:        *toPassword,
		ToDatabase:        *toDatabase,
		ToTable:           *toTable,
		Cleanup:           *cleanup,
		MySQLDump:         *mysqlDump,
		Threads:           *threads,
		Behinds:           *behinds,
		RadonURL:          *radonURL,
		Checksum:          *checksum,
	}
	cfg.DBTablesMaps = make(map[string][]string)
	log.Info("datatravel.cfg:%+v", cfg)
	shift := shift.NewShift(log, cfg)
	if err := shift.Start(); err != nil {
		log.Panicf("shift.start.error:%+v", err)
	}
	defer shift.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	select {
	case <-shift.Done():
		fmt.Println("datatravel.exit.normal!")
	case <-sc:
		fmt.Println("datatravel.catch.signal...")
	}
}
