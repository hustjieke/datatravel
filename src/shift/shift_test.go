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
	"sync"
	"testing"
	"time"
	"xlog"

	"github.com/stretchr/testify/assert"
)

func assertChecksumEqual(t *testing.T, shift *Shift) {
	// if checksum ok, we`ll get true(done) finally.
	assert.True(t, <-shift.Done())
}

func assertChecksumNotEqual(t *testing.T, shift *Shift) {
	// check checksum.
	<-shift.Done()

	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)

	sql := fmt.Sprintf("checksum table `%s`.`%s`", shift.cfg.FromDatabase, shift.cfg.FromTable)
	r, err := fromConn.Execute(sql)
	assert.Nil(t, err)
	from, err := r.GetString(0, 1)
	assert.Nil(t, err)

	sql = fmt.Sprintf("checksum table `%s`.`%s`", shift.cfg.ToDatabase, shift.cfg.ToTable)
	r, err = fromConn.Execute(sql)
	assert.Nil(t, err)
	to, err := r.GetString(0, 1)
	assert.Nil(t, err)
	assert.NotEqual(t, from, to)
}

func TestShiftInsert(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		begin = begin + step
		log.Debug("test.shift.insert.done")
	}

	// MultiInserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d'),(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i, i+1, i+1, i+1)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.multi.insert.done")
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
		log.Debug("checksum done, do cleanup")
	}
}

/*
func TestShiftInsertJson(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into %s.%s(a,b,c,e) values(%d,%d,'%d', '{\"screen\": \"50 inch\"}')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		begin = begin + step
		log.Debug("test.shift.insert.done")
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}
*/

func TestShiftDelete(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Delete.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			sql := fmt.Sprintf("delete from `%s`.`%s` where a=%d", shift.cfg.FromDatabase, shift.cfg.FromTable, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.delete.done")
	}

	// MultiDeletes.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		sql := fmt.Sprintf("delete from `%s`.`%s`", shift.cfg.FromDatabase, shift.cfg.FromTable)
		_, err := fromConn.Execute(sql)
		assert.Nil(t, err)
		log.Debug("test.shift.multi.delete.done")
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftDeleteWithoutPK(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, false, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Delete.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			sql := fmt.Sprintf("delete from `%s`.`%s` where a=%d", shift.cfg.FromDatabase, shift.cfg.FromTable, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.delete.done")
	}

	// MultiDeletes.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		sql := fmt.Sprintf("delete from `%s`.`%s`", shift.cfg.FromDatabase, shift.cfg.FromTable)
		_, err := fromConn.Execute(sql)
		assert.Nil(t, err)
		log.Debug("test.shift.multi.delete.done")
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftUpdate(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Update.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			sql := fmt.Sprintf("update `%s`.`%s` set b=1 where a=%d", shift.cfg.FromDatabase, shift.cfg.FromTable, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.update.done")
	}

	//MultiUpdates
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		sql := fmt.Sprintf("update `%s`.`%s` set b=1 where a in (1,3,5)", shift.cfg.FromDatabase, shift.cfg.FromTable)
		_, err := fromConn.Execute(sql)
		assert.Nil(t, err)
		log.Debug("test.shift.multi.update.done")
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftUpdateWithoutPK(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, false, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Update.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			sql := fmt.Sprintf("update `%s`.`%s` set b=1 where a=%d", shift.cfg.FromDatabase, shift.cfg.FromTable, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.update.done")
	}

	//MultiUpdates
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		sql := fmt.Sprintf("update `%s`.`%s` set b=1 where a in (1,3,5)", shift.cfg.FromDatabase, shift.cfg.FromTable)
		_, err := fromConn.Execute(sql)
		assert.Nil(t, err)
		log.Debug("test.shift.multi.update.done")
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftReplace(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// replace.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin - 2; i < begin+step; i++ {
			sql := fmt.Sprintf("replace into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftIntegerUnsigned(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, false, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)

	// Inserts.
	sql := fmt.Sprintf("insert into `%s`.`%s`(a,f,g,h,i,j,k,l,m,n) values(-2147483648,0,-9223372036854775808,0,-128,0,-32768,0,-8388608,0)", shift.cfg.FromDatabase, shift.cfg.FromTable)
	_, err := fromConn.Execute(sql)
	assert.Nil(t, err)

	sql = fmt.Sprintf("insert into `%s`.`%s`(a,f,g,h,i,j,k,l,m,n) values(2147483647,4294967295,9223372036854775807,18446744073709551615,127,255,32767,65535,8388607,16777214)", shift.cfg.FromDatabase, shift.cfg.FromTable)
	_, err = fromConn.Execute(sql)
	assert.Nil(t, err)

	sql = fmt.Sprintf("insert into `%s`.`%s`(a,f,g,h,i,j,k,l,m,n) values(2147483646,4294967294,9223372036854775806,18446744073709551614,121,252,32761,65532,8388605,16777213)", shift.cfg.FromDatabase, shift.cfg.FromTable)
	_, err = fromConn.Execute(sql)
	assert.Nil(t, err)

	// update
	sql = fmt.Sprintf("update `%s`.`%s` set f=4294967291 where h=18446744073709551615", shift.cfg.FromDatabase, shift.cfg.FromTable)
	_, err = fromConn.Execute(sql)
	assert.Nil(t, err)

	// delete
	sql = fmt.Sprintf("delete from `%s`.`%s` where h=18446744073709551614", shift.cfg.FromDatabase, shift.cfg.FromTable)
	_, err = fromConn.Execute(sql)
	assert.Nil(t, err)

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}

}

func TestShiftXACommit(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShiftXa(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)

	step := 7
	begin := 0
	// Inserts.
	{
		for i := begin; i < begin+step; i++ {
			// xa start.
			sql := fmt.Sprintf("xa start 'xc%d'", i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)

			sql = fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa end.
			sql = fmt.Sprintf("xa end 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa prepare.
			sql = fmt.Sprintf("xa prepare 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa commit.
			sql = fmt.Sprintf("xa commit 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Update.
	{
		for i := begin; i < begin+step; i += 2 {
			// xa start.
			sql := fmt.Sprintf("xa start 'xc%d'", i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)

			sql = fmt.Sprintf("update `%s`.`%s` set b=1 where a=%d", shift.cfg.FromDatabase, shift.cfg.FromTable, i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa end.
			sql = fmt.Sprintf("xa end 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa prepare.
			sql = fmt.Sprintf("xa prepare 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa commit.
			sql = fmt.Sprintf("xa commit 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Delete.
	{
		sql := fmt.Sprintf("use %s", shift.cfg.FromDatabase)
		_, err := fromConn.Execute(sql)
		assert.Nil(t, err)

		for i := begin; i < begin+step; i += 2 {
			// xa start.
			sql := fmt.Sprintf("xa start 'xc%d'", i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)

			sql = fmt.Sprintf("delete from `%s` where a=%d", shift.cfg.FromTable, i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa end.
			sql = fmt.Sprintf("xa end 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa prepare.
			sql = fmt.Sprintf("xa prepare 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa commit.
			sql = fmt.Sprintf("xa commit 'xc%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftXARollback(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShiftXa(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			// xa start.
			sql := fmt.Sprintf("xa start 'xr%d'", i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)

			sql = fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa end.
			sql = fmt.Sprintf("xa end 'xr%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa prepare.
			sql = fmt.Sprintf("xa prepare 'xr%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			if i%3 == 0 {
				sql = fmt.Sprintf("xa commit 'xr%d'", i)
			} else {
				sql = fmt.Sprintf("xa rollback 'xr%d'", i)
			}
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Update.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			// xa start.
			sql := fmt.Sprintf("xa start 'xr%d'", i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)

			sql = fmt.Sprintf("update `%s`.`%s` set b=1 where a=%d", shift.cfg.FromDatabase, shift.cfg.FromTable, i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa end.
			sql = fmt.Sprintf("xa end 'xr%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa prepare.
			sql = fmt.Sprintf("xa prepare 'xr%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa.
			if i%2 == 0 {
				sql = fmt.Sprintf("xa commit 'xr%d'", i)
			} else {
				sql = fmt.Sprintf("xa rollback 'xr%d'", i)
			}
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Delete.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			// xa start.
			sql := fmt.Sprintf("xa start 'xr%d'", i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)

			sql = fmt.Sprintf("delete from `%s`.`%s` where a=%d", shift.cfg.FromDatabase, shift.cfg.FromTable, i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa end.
			sql = fmt.Sprintf("xa end 'xr%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa prepare.
			sql = fmt.Sprintf("xa prepare 'xr%d'", i)
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)

			// xa.
			if i%2 == 0 {
				sql = fmt.Sprintf("xa commit 'xr%d'", i)
			} else {
				sql = fmt.Sprintf("xa rollback 'xr%d'", i)
			}
			_, err = fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftInsertWithDump(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShiftWithData(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		begin = begin + step
		log.Debug("test.shift.insert.done")
	}

	// MultiInserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d'),(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i, i+1, i+1, i+1)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.multi.insert.done")
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftChecksumTable(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShift(log, false, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	// Delete.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			sql := fmt.Sprintf("delete from `%s`.`%s` where a=%d", shift.cfg.FromDatabase, shift.cfg.FromTable, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
	}

	{
		assertChecksumEqual(t, shift)
	}

	err := shift.ChecksumTables()
	assert.Nil(t, err)
}

func TestShiftMySQLTable(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShiftMysqlTable(log, false, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	// Delete.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		sql := fmt.Sprintf("delete from `%s`.`%s` where user='%s'", shift.cfg.FromDatabase, shift.cfg.FromTable, "root")
		_, err := fromConn.Execute(sql)
		assert.Nil(t, err)
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

func TestShiftMySQLTableWithDatas(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShiftMysqlTableWithData(log, false, "mysql")
	defer cleanup()

	// Checksum check.
	{
		assertChecksumNotEqual(t, shift)
	}
}

func TestShiftWithCleanup(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShiftWithCleanup(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockPanicMe

	step := 7
	begin := 0
	// Inserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		begin = begin + step
		log.Debug("test.shift.insert.done")
	}

	// MultiInserts.
	{
		fromConn := shift.fromPool.Get()
		defer shift.fromPool.Put(fromConn)

		for i := begin; i < begin+step; i += 2 {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d'),(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i, i+1, i+1, i+1)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.multi.insert.done")
	}

	// Checksum check.
	{
		assertChecksumEqual(t, shift)
	}
}

/*
func TestShiftDDLEvent(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	shift, cleanup := MockShiftDDL(log, true, "mysql")
	defer cleanup()
	shift.panicHandler = mockRecoverPanicMe
	fromConn := shift.fromPool.Get()
	defer shift.fromPool.Put(fromConn)

	begin := 0
	step := 7
	// Inserts.
	{
		for i := begin; i < begin+step; i++ {
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", shift.cfg.FromDatabase, shift.cfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.insert.done")
	}

	// Create database event.
	{
		sql := fmt.Sprintf("create database %s_shift_test", shift.cfg.FromDatabase)
		_, err := fromConn.Execute(sql)
		assert.NotNil(t, err)
	}

	// Drop database event.
	{
		sql := fmt.Sprintf("drop database %s_shift_test", shift.cfg.FromDatabase)
		_, err := fromConn.Execute(sql)
		assert.NotNil(t, err)
	}

	// Truncate db.table event.
	{
		sql := fmt.Sprintf("truncate table %s.%s", shift.cfg.FromDatabase, shift.cfg.FromTable)
		_, err := fromConn.Execute(sql)
		assert.NotNil(t, err)
	}

	// Alter db.table event.
	{
		sql := fmt.Sprintf("alter table %s.%s add xxx int", shift.cfg.FromDatabase, shift.cfg.FromTable)
		_, err := fromConn.Execute(sql)
		assert.NotNil(t, err)
	}

	// Alter table event.
	{
		sql := fmt.Sprintf("use %s", shift.cfg.FromDatabase)
		_, err := fromConn.Execute(sql)
		assert.NotNil(t, err)

		sql = fmt.Sprintf("alter table %s engine=MyISAM", shift.cfg.FromTable)
		_, err = fromConn.Execute(sql)
		assert.NotNil(t, err)
	}

	// Create table xx event.
	{
		sql := fmt.Sprintf("create table if not exists xx(a int)")
		_, err := fromConn.Execute(sql)
		assert.NotNil(t, err)
	}

	// Drop table xx event.
	{
		sql := fmt.Sprintf("drop table xx")
		_, err := fromConn.Execute(sql)
		assert.NotNil(t, err)
	}

	// Checksum check.
	{
		assertChecksumNotEqual(t, shift)
	}
}
*/

func TestShiftStart(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	mockCfg.ToFlavor = "mysql"
	mockCfg.DBTablesMaps = make(map[string][]string) // init map
	mockCfg.Databases = make([]string, 0, 0)         // init dbs
	mockCfg.FromRows = 0
	mockCfg.ToRows = 0
	shift := NewShift(log, mockCfg)
	defer shift.Close()

	err := shift.Start()
	assert.Nil(t, err)
}

// Fix bug for issue #4
func TestShiftCanalClose(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	mockCfg.ToFlavor = "mysql"
	mockCfg.DBTablesMaps = make(map[string][]string) // init map
	mockCfg.Databases = make([]string, 0, 0)         // init dbs
	mockCfg.FromRows = 0
	mockCfg.ToRows = 0

	shift := NewShift(log, mockCfg)
	// Replace time ticker every 5s to 1s, 5s is to long for us to test
	shift.behindsTicker = time.NewTicker(time.Duration(1000) * time.Millisecond)

	fromPool, errfrom := NewPool(log, 4, mockCfg.From, mockCfg.FromUser, mockCfg.FromPassword, false)
	assert.Nil(t, errfrom)
	toPool, errto := NewPool(log, 4, mockCfg.To, mockCfg.ToUser, mockCfg.ToPassword, false)
	assert.Nil(t, errto)

	var c = make(chan bool, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	// Keep continuing insert so that canal will be always running and we
	// have enough time to use shift.Close() to simulate signal like kill
	go func() {
		begin := 0
		step := 5000
		fromConn := fromPool.Get()
		defer fromPool.Put(fromConn)
		toConn := toPool.Get()
		defer toPool.Put(toConn)

		// Cleanup From table first.
		sql := fmt.Sprintf("drop table if exists `%s`.`%s`", mockCfg.FromDatabase, mockCfg.FromTable)
		if _, err := fromConn.Execute(sql); err != nil {
			log.Panicf("test.drop.from.table.error:%+v", err)
		}

		// Cleanup To table.
		sql = fmt.Sprintf("drop table if exists `%s`.`%s`", mockCfg.FromDatabase, mockCfg.FromTable)
		if _, err := toConn.Execute(sql); err != nil {
			log.Panicf("test.drop.to.table.error:%+v", err)
		}

		// Create database on from.
		sql = fmt.Sprintf("create database if not exists `%s`", mockCfg.FromDatabase)
		if _, err := fromConn.Execute(sql); err != nil {
			log.Panicf("test.prepare.database.error:%+v", err)
		}

		// Create table on from.
		sql = fmt.Sprintf("create table `%s`.`%s`(a int primary key, b int, c varchar(200), d DOUBLE NULL DEFAULT NULL, e json DEFAULT NULL, f INT UNSIGNED DEFAULT NULL, g BIGINT DEFAULT NULL, h BIGINT UNSIGNED DEFAULT NULL, i TINYINT NULL, j TINYINT UNSIGNED DEFAULT NULL, k SMALLINT DEFAULT NULL, l SMALLINT UNSIGNED DEFAULT NULL, m MEDIUMINT DEFAULT NULL, n INT UNSIGNED DEFAULT NULL)", mockCfg.FromDatabase, mockCfg.FromTable)
		if _, err := fromConn.Execute(sql); err != nil {
			log.Panicf("test.prepare.from.table.error:%+v", err)
		}

		for i := begin; i < begin+step; i++ {
			select {
			case <-c:
				log.Info("test.gets.signal.done.and.insert.nums:%v", i-1)
				i = 5000
			default:
			}
			sql := fmt.Sprintf("insert into `%s`.`%s`(a,b,c) values(%d,%d,'%d')", mockCfg.FromDatabase, mockCfg.FromTable, i, i, i)
			_, err := fromConn.Execute(sql)
			assert.Nil(t, err)
		}
		log.Debug("test.shift.insert.done")
		wg.Done()
	}()

	// Sleep 1s so that tb1 has some rows before we start canal and dump
	time.Sleep(time.Second * 1)

	// Start
	var errstart error
	errstart = shift.Start()
	assert.Nil(t, errstart)

	// Sleep 2s to make sure we have enough time that shift.Start() having executed
	// canal.Run(), then we can use shift.Close() to simulate signal like kill
	time.Sleep(time.Second * 2)
	log.Debug("sleep 2s")
	shift.Close()
	c <- true
	wg.Wait()
	assert.False(t, shift.allDone)
}

// Fix bug for issue #30
func TestShiftParseBOM(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	mockCfg.ToFlavor = "mysql"
	mockCfg.DBTablesMaps = make(map[string][]string) // init map
	mockCfg.Databases = make([]string, 0, 0)         // init dbs
	mockCfg.FromRows = 0
	mockCfg.ToRows = 0

	fromPool, errfrom := NewPool(log, 4, mockCfg.From, mockCfg.FromUser, mockCfg.FromPassword, false)
	assert.Nil(t, errfrom)
	defer fromPool.Close()
	toPool, errto := NewPool(log, 4, mockCfg.To, mockCfg.ToUser, mockCfg.ToPassword, false)
	assert.Nil(t, errto)
	defer toPool.Close()

	fromConn := fromPool.Get()
	toConn := toPool.Get()
	defer func() {
		// Clean
		sql := `DROP DATABASE IF EXISTS DTtest`
		if _, err := fromConn.Execute(sql); err != nil {
			log.Panicf("test.drop.databse.error:%+v", err)
		}
		fromPool.Put(fromConn)
		toPool.Put(toConn)
	}()

	// Drop test DTtest DB
	sql := `DROP DATABASE IF EXISTS DTtest`
	if _, err := fromConn.Execute(sql); err != nil {
		log.Panicf("test.drop.from.databse.error:%+v", err)
	}
	if _, err := toConn.Execute(sql); err != nil {
		log.Panicf("test.drop.from.databse.error:%+v", err)
	}
	sql = `DROP DATABASE IF EXISTS shift_test_from`
	if _, err := fromConn.Execute(sql); err != nil {
		log.Panicf("test.drop.from.databse.error:%+v", err)
	}
	if _, err := toConn.Execute(sql); err != nil {
		log.Panicf("test.drop.to.databse.error:%+v", err)
	}

	// Create test DTtest DB
	sql = `CREATE DATABASE DTtest`
	if _, err := fromConn.Execute(sql); err != nil {
		log.Panicf("test.create.databse.error:%+v", err)
	}
	sql = `USE DTtest`
	if _, err := fromConn.Execute(sql); err != nil {
		log.Panicf("test.use.databse.error:%+v", err)
	}
	// Create test table
	sql = `CREATE TABLE normal(
            id BIGINT(64) UNSIGNED  NOT NULL AUTO_INCREMENT,
            str VARCHAR(256),
            f FLOAT,
            d DOUBLE,
            de DECIMAL(10,2),
            i INT,
            bi BIGINT,
            e enum ("e1", "e2"),
            b BIT(8),
            y YEAR,
            da DATE,
            ts TIMESTAMP,
            dt DATETIME,
            tm TIME,
            t TEXT,
            bb BLOB,
            se SET('a', 'b', 'c'),
            PRIMARY KEY (id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8`
	if _, err := fromConn.Execute(sql); err == nil {
		fromConn.Execute(`INSERT INTO normal(str, f, i, e, b, y, da, ts, dt, tm, de, t, bb, se)
        VALUES ("3", -3.14, 10, "e1", 0b0011, 1985,
        "2012-05-07", "2012-05-07 14:01:01", "2012-05-07 14:01:01",
        "14:01:01", -45363.64, "abc", "12345", "a,b")`)
	}

	sql = `CREATE TABLE test_json (
            id BIGINT(64) UNSIGNED  NOT NULL AUTO_INCREMENT,
            c1 JSON,
            c2 DECIMAL(10, 0),
            PRIMARY KEY (id)
            ) ENGINE=InnoDB`

	if _, err := fromConn.Execute(sql); err == nil {
		fromConn.Execute(`INSERT INTO test_json (c2) VALUES (1)`)
		fromConn.Execute(`INSERT INTO test_json (c1, c2) VALUES ('{"key1": "value1", "key2": "value2"}', 1)`)
	}

	sql = `CREATE TABLE test_json_v2 (
            id INT,
            c JSON,
            PRIMARY KEY (id)
            ) ENGINE=InnoDB`

	if _, err := fromConn.Execute(sql); err == nil {
		tbls := []string{
			// Refer: https://github.com/shyiko/mysql-binlog-connector-java/blob/c8e81c879710dc19941d952f9031b0a98f8b7c02/src/test/java/com/github/shyiko/mysql/binlog/event/deserialization/json/JsonBinaryValueIntegrationTest.java#L84
			// License: https://github.com/shyiko/mysql-binlog-connector-java#license
			`INSERT INTO test_json_v2 VALUES (0, NULL)`,
			`INSERT INTO test_json_v2 VALUES (1, '{\"a\": 2}')`,
			`INSERT INTO test_json_v2 VALUES (2, '[1,2]')`,
			`INSERT INTO test_json_v2 VALUES (3, '{\"a\":\"b\", \"c\":\"d\",\"ab\":\"abc\", \"bc\": [\"x\", \"y\"]}')`,
			`INSERT INTO test_json_v2 VALUES (4, '[\"here\", [\"I\", \"am\"], \"!!!\"]')`,
			`INSERT INTO test_json_v2 VALUES (5, '\"scalar string\"')`,
			`INSERT INTO test_json_v2 VALUES (6, 'true')`,
			`INSERT INTO test_json_v2 VALUES (7, 'false')`,
			`INSERT INTO test_json_v2 VALUES (8, 'null')`,
			`INSERT INTO test_json_v2 VALUES (9, '-1')`,
			`INSERT INTO test_json_v2 VALUES (10, CAST(CAST(1 AS UNSIGNED) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (11, '32767')`,
			`INSERT INTO test_json_v2 VALUES (12, '32768')`,
			`INSERT INTO test_json_v2 VALUES (13, '-32768')`,
			`INSERT INTO test_json_v2 VALUES (14, '-32769')`,
			`INSERT INTO test_json_v2 VALUES (15, '2147483647')`,
			`INSERT INTO test_json_v2 VALUES (16, '2147483648')`,
			`INSERT INTO test_json_v2 VALUES (17, '-2147483648')`,
			`INSERT INTO test_json_v2 VALUES (18, '-2147483649')`,
			`INSERT INTO test_json_v2 VALUES (19, '18446744073709551615')`,
			`INSERT INTO test_json_v2 VALUES (20, '18446744073709551616')`,
			`INSERT INTO test_json_v2 VALUES (21, '3.14')`,
			`INSERT INTO test_json_v2 VALUES (22, '{}')`,
			`INSERT INTO test_json_v2 VALUES (23, '[]')`,
			`INSERT INTO test_json_v2 VALUES (24, CAST(CAST('2015-01-15 23:24:25' AS DATETIME) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (25, CAST(CAST('23:24:25' AS TIME) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (125, CAST(CAST('23:24:25.12' AS TIME(3)) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (225, CAST(CAST('23:24:25.0237' AS TIME(3)) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (26, CAST(CAST('2015-01-15' AS DATE) AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (27, CAST(TIMESTAMP'2015-01-15 23:24:25' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (127, CAST(TIMESTAMP'2015-01-15 23:24:25.12' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (227, CAST(TIMESTAMP'2015-01-15 23:24:25.0237' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (327, CAST(UNIX_TIMESTAMP('2015-01-15 23:24:25') AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (28, CAST(ST_GeomFromText('POINT(1 1)') AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (29, CAST('[]' AS CHAR CHARACTER SET 'ascii'))`,
			// TODO: 30 and 31 are BIT type from JSON_TYPE, may support later.
			`INSERT INTO test_json_v2 VALUES (30, CAST(x'cafe' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (31, CAST(x'cafebabe' AS JSON))`,
			`INSERT INTO test_json_v2 VALUES (100, CONCAT('{\"', REPEAT('a', 64 * 1024 - 1), '\":123}'))`,
		}

		for _, query := range tbls {
			fromConn.Execute(query)
		}

		// If MySQL supports JSON, it must supports GEOMETRY.
		sql = `CREATE TABLE test_geo (g GEOMETRY)`
		_, err = fromConn.Execute(sql)
		assert.Nil(t, err)

		tbls = []string{
			`INSERT INTO test_geo VALUES (POINT(1, 1))`,
			`INSERT INTO test_geo VALUES (LINESTRING(POINT(0,0), POINT(1,1), POINT(2,2)))`,
			// TODO: add more geometry tests
		}

		for _, query := range tbls {
			fromConn.Execute(query)
		}
	}

	// Must allow zero time.
	fromConn.Execute(`SET sql_mode=''`)
	str := `CREATE TABLE test_parse_time (
            a1 DATETIME,
            a2 DATETIME(3),
            a3 DATETIME(6),
            b1 TIMESTAMP,
            b2 TIMESTAMP(3) ,
            b3 TIMESTAMP(6))`
	fromConn.Execute(str)

	fromConn.Execute(`INSERT INTO test_parse_time VALUES
        ("2014-09-08 17:51:04.123456", "2014-09-08 17:51:04.123456", "2014-09-08 17:51:04.123456",
        "2014-09-08 17:51:04.123456","2014-09-08 17:51:04.123456","2014-09-08 17:51:04.123456"),
        ("0000-00-00 00:00:00.000000", "0000-00-00 00:00:00.000000", "0000-00-00 00:00:00.000000",
        "0000-00-00 00:00:00.000000", "0000-00-00 00:00:00.000000", "0000-00-00 00:00:00.000000"),
        ("2014-09-08 17:51:04.000456", "2014-09-08 17:51:04.000456", "2014-09-08 17:51:04.000456",
        "2014-09-08 17:51:04.000456","2014-09-08 17:51:04.000456","2014-09-08 17:51:04.000456")`)

	// Start
	shift := NewShift(log, mockCfg)
	err := shift.Start()
	time.Sleep(time.Second * 1)
	assert.Nil(t, err)
}
