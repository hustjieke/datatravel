/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package config

// Use flavor for different target cluster
const (
	ToMySQLFlavor   = "mysql"
	ToMariaDBFlavor = "mariadb"
	ToRadonDBFlavor = "radondb"
)

type Config struct {
	ToFlavor           string
	SetGlobalReadLock  bool
	MetaDir            string
	FkCheck            bool
	MaxAllowedPacketMB int

	From         string
	FromUser     string
	FromPassword string
	FromDatabase string
	FromTable    string

	To         string
	ToUser     string
	ToPassword string
	ToDatabase string
	ToTable    string

	MySQLDump string
	Threads   int
	Behinds   int
	Checksum  bool

	// Will override Databases
	TableDB   string
	Tables    []string
	Databases []string

	DBTablesMaps map[string][]string // key:db, value: tables

	// FromRows and ToRows are used to count dump progress
	FromRows uint64
	ToRows   uint64

	// Next args are used in case mysql-->radondb
	AutoIncTable map[string]bool // key:db.table, value: bool
	// count(*) of to db.tbls, here we only count tables include
	// auto_increment columns key: db.tbl, value: count(*) rows
	ToTblsRowsBefore map[string]uint64
	// checksum of db.tbls of to,key: db.tbl, value: checksum table db.tbl
	ToTblsChecksumBefore map[string]uint32
	IsNotFisrtTime       bool
}
