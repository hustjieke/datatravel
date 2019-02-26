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
	"net/http"
	"strings"
	"time"
)

func (shift *Shift) setRadonReadOnly(v bool) error {
	log := shift.log
	cfg := shift.cfg
	path := cfg.RadonURL + "/v1/radon/readonly"

	type request struct {
		Readonly bool `json:"readonly"`
	}
	req := &request{
		Readonly: v,
	}
	log.Info("shift.set.radon[%s].readonlly.req[%+v]", path, req)

	resp, cleanup, err := HTTPPut(path, req)
	defer cleanup()
	if err != nil {
		return err
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("shift.set.radon.readonly[%s].response.error:%+s", path, HTTPReadBody(resp))
	}
	return nil
}

func (shift *Shift) setRadonRule() error {
	log := shift.log
	cfg := shift.cfg
	path := cfg.RadonURL + "/v1/shard/shift"

	if _, isSystem := sysDatabases[strings.ToLower(shift.cfg.FromDatabase)]; isSystem {
		log.Info("shift.set.radon.rune.skip.system.table:[%s.%s]", shift.cfg.FromDatabase, shift.cfg.FromTable)
		return nil
	}

	type request struct {
		Database    string `json:"database"`
		Table       string `json:"table"`
		FromAddress string `json:"from-address"`
		ToAddress   string `json:"to-address"`
	}
	req := &request{
		Database:    cfg.FromDatabase,
		Table:       cfg.FromTable,
		FromAddress: cfg.From,
		ToAddress:   cfg.To,
	}
	log.Info("shift.set.radon[%s].rule.req[%+v]", path, req)

	resp, cleanup, err := HTTPPost(path, req)
	defer cleanup()
	if err != nil {
		return err
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("shift.set.radon.shard.rule[%s].response.error:%+s", path, HTTPReadBody(resp))
	}
	return nil
}

var (
	radon_limits_min = 2400
	radon_limits_max = 80000
)

func (shift *Shift) setRadonThrottle(factor float32) error {
	log := shift.log
	cfg := shift.cfg
	path := cfg.RadonURL + "/v1/radon/throttle"

	type request struct {
		Limits int `json:"limits"`
	}

	// limits =0 means unlimits.
	limits := int(float32(radon_limits_max) * factor)
	if limits != 0 && limits < radon_limits_min {
		limits = radon_limits_min
	}
	req := &request{
		Limits: limits,
	}
	log.Info("shift.set.radon[%s].throttle.to.req[%+v].by.factor[%v]", path, req, factor)

	resp, cleanup, err := HTTPPut(path, req)
	defer cleanup()
	if err != nil {
		return err
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("shift.set.radon.throttle[%s].response.error:%+s", path, HTTPReadBody(resp))
	}
	return nil
}

func (shift *Shift) setRadon() {
	log := shift.log

	// 1. WaitUntilPos
	{
		masterPos := shift.masterPosition()
		log.Info("shift.wait.until.pos[%#v]...", masterPos)
		if err := shift.canal.WaitUntilPos(*masterPos, time.Hour*12); err != nil {
			shift.panicMe("shift.set.radon.wait.until.pos[%#v].error:%+v", masterPos, err)
			return
		}
		log.Info("shift.wait.until.pos.done...")
	}

	// 2. Set radon to readonly.
	{
		log.Info("shift.set.radon.readonly...")
		if err := shift.setRadonReadOnly(true); err != nil {
			shift.panicMe("shift.set.radon.readonly.error:%+v", err)
			return
		}
		log.Info("shift.set.radon.readonly.done...")
	}

	// 3. Wait again.
	{
		masterPos := shift.masterPosition()
		log.Info("shift.wait.until.pos[%#v]...", masterPos)
		if err := shift.canal.WaitUntilPos(*masterPos, time.Second*300); err != nil {
			shift.panicMe("shift.wait.until.pos[%#v].error:%+v", masterPos, err)
			return
		}
		log.Info("shift.wait.until.pos.done...")
	}

	// 4. Checksum table.
	if shift.cfg.Checksum {
		log.Info("shift.checksum.table...")
		if err := shift.ChecksumTable(); err != nil {
			shift.panicMe("shift.checksum.table.error:%+v", err)
			return
		}
		log.Info("shift.checksum.table.done...")
	}

	// 5. Set radon rule.
	{
		log.Info("shift.set.radon.rule...")
		if err := shift.setRadonRule(); err != nil {
			shift.panicMe("shift.set.radon.rule.error:%+v", err)
			return
		}
		log.Info("shift.set.radon.rule.done...")
	}

	// 6. Set radon to read/write.
	{
		log.Info("shift.set.radon.to.write...")
		if err := shift.setRadonReadOnly(false); err != nil {
			shift.panicMe("shift.set.radon.write.error:%+v", err)
			return
		}
		log.Info("shift.set.radon.to.write.done...")
	}

	// 7. Set radon throttle to unlimits.
	{
		log.Info("shift.set.radon.throttle.to.unlimits...")
		if err := shift.setRadonThrottle(0); err != nil {
			shift.panicMe("shift.set.radon.throttle.to.unlimits.error:%+v", err)
			return
		}
		log.Info("shift.set.radon.throttle.to.unlimits.done...")
	}

	// 8. Good, we have all done.
	{
		shift.done <- true
		shift.allDone = true
		log.Info("shift.all.done...")
	}
}
