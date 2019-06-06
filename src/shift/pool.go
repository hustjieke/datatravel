/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package shift

import (
	"sync"
	"xlog"

	"github.com/siddontang/go-mysql/client"
)

// Blocked connection pool.
type Pool struct {
	log   *xlog.Log
	conns chan *client.Conn
	mu    sync.Mutex

	host     string
	user     string
	password string
}

func NewPool(log *xlog.Log, cap int, host string, user string, password string, fkCheck bool) (*Pool, error) {
	conns := make(chan *client.Conn, cap)
	for i := 0; i < cap; i++ {
		to, err := client.Connect(host, user, password, "")
		if err != nil {
			log.Error("shift.new.pool.connection.error:%+v", err)
			return nil, err
		}
		if !fkCheck {
			to.SetForeignKeyCheckDisable()
		}
		conns <- to
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

func (p *Pool) Get() *client.Conn {
	log := p.log
	var err error
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conns == nil {
		log.Error("datatravel.pool.get.conns.is.nil")
		return nil
	}

	conn := <-p.conns
	if conn == nil {
		log.Error("datatravel.pool.get.conn.is.nil")
		return nil
	}

	if err = conn.Ping(); err != nil {
		log.Warning("shift.get.connection.was.bad, prepare.a.new.connection")
		conn, err = client.Connect(p.host, p.user, p.password, "")
		if err != nil {
			log.Error("shift.get.connection.error:%+v", err)
			return nil
		}
	}
	return conn
}

func (p *Pool) Put(conn *client.Conn) {
	log := p.log
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conns == nil || conn == nil {
		log.Error("datatravel.pool.put.conns.or.conn.is.nil")
		return
	}
	p.conns <- conn
}

func (p *Pool) Close() {
	log := p.log
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conns == nil {
		log.Error("datatravel.pool.close.conns.is.nil")
		return
	}
	close(p.conns)
	for conn := range p.conns {
		conn.Close()
	}
	p.conns = nil
}
