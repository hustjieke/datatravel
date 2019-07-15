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
	AutoIncTable map[string]bool     // key:db.table, value: bool

	// FromRows and ToRows are used to count dump progress
	FromRows uint64
	ToRows   uint64
}
