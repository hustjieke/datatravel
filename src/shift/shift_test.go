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
