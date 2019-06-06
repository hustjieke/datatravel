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
	"testing"
	"xlog"

	"github.com/siddontang/go-mysql/client"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	// Config for normal
	cfg := &config.Config{
		From:         "127.0.0.1:3306",
		FromUser:     "root",
		FromDatabase: "shift_test_from",
		FromTable:    "t1",

		To:         "127.0.0.1:3307",
		ToUser:     "root",
		ToDatabase: "shift_test_to",
		ToTable:    "t1",

		Threads:   16,
		Behinds:   256,
		MySQLDump: "mysqldump",
		Checksum:  true,
	}

	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	poolNormal, err := NewPool(log, 4, cfg.From, cfg.FromUser, cfg.FromPassword, false)
	assert.Nil(t, err)

	var conn *client.Conn
	// Test normal Get(), conn is not nil
	conn = poolNormal.Get()
	assert.NotNil(t, conn)

	// Test pool close first, conn we get is nil
	poolNormal.Close()
	conn = poolNormal.Get()
	assert.Nil(t, conn)

	// Test exception case, when we get conn from pool in one goroutine, and
	// pool is closed in another goroutine, the conns may be closed and will be set nil, e.g.
	// thread1: Get()-->getConns()-->we get a conns pointer and conns is not nil-->thread 2: Closed()-->
	// closed conns and the elements in conns-->get conn from conns(<-conns)-->then we got nil....
	poolError, err1 := fakeNewPool(log, 4, cfg.From, cfg.FromUser, cfg.FromPassword)
	assert.Nil(t, err1)
	// Here conns is not nil althouth it has 4 nil elements in chan.
	conn = poolError.Get()
	assert.Nil(t, conn)
}

func fakeNewPool(log *xlog.Log, cap int, host string, user string, password string) (*Pool, error) {
	conns := make(chan *client.Conn, cap)
	for i := 0; i < cap; i++ {
		conns <- nil
	}
	log.Info("shift.pool[host:%v, cap:%d].done", host, cap)

	return &Pool{
		log:      log,
		conns:    conns,
		host:     host,
		user:     user,
		password: password,
	}, nil
}
