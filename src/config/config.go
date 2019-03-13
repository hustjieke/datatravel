/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package config

// Use flavor for different destinations,
const (
	ToMySQLFlavor   = "mysql"
	ToMariaDBFlavor = "mariadb"
	ToRadonDBFlavor = "radondb"
)

type Config struct {
	ToFlavor          string
	SetGlobalReadLock bool

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

	Cleanup   bool
	MySQLDump string
	Threads   int
	Behinds   int
	RadonURL  string
	Checksum  bool

	Databases    []string
	DBTablesMaps map[string][]string // key:db, value: tables
}
