package canal

import (
	"regexp"
	"time"

	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/siddontang/go-mysql/schema"
)

var (
	expAlterTable = regexp.MustCompile("(?i)^ALTER\\s{1,}TABLE\\s{1,}.*?`{0,1}(.*?)`{0,1}\\.{0,1}`{0,1}([^`\\.]+?)`{0,1}\\s.*")
)

func (c *Canal) startSyncBinlog() error {
	pos := c.master.Position()

	log.Infof("start sync binlog at %v", pos)

	// prepare TableSchema
	_, err := c.GetTable(c.dumper.TableDB, c.dumper.Tables[0])
	if err != nil {
		return errors.Trace(err)
	}

	s, err := c.syncer.StartSync(pos)
	if err != nil {
		return errors.Errorf("start sync replication at %v error %v", pos, err)
	}

	for {
		ev, err := s.GetEvent(c.ctx)

		if err != nil {
			return errors.Trace(err)
		}

		curPos := pos.Pos
		//next binlog pos
		pos.Pos = ev.Header.LogPos

		// We only save position with RotateEvent and XIDEvent.
		// For RowsEvent, we can't save the position until meeting XIDEvent
		// which tells the whole transaction is over.
		// TODO: If we meet any DDL query, we must save too.
		switch e := ev.Event.(type) {
		case *replication.RotateEvent:
			pos.Name = string(e.NextLogName)
			pos.Pos = uint32(e.Position)
			log.Infof("rotate binlog to %s", pos)

			if err = c.eventHandler.OnRotate(e); err != nil {
				return errors.Trace(err)
			}
		case *replication.RowsEvent:
			// we only focus row based event
			err = c.handleRowsEvent(ev)
			if err != nil && errors.Cause(err) != schema.ErrTableNotExist {
				// We can ignore table not exist error
				log.Errorf("handle rows event at (%s, %d) error %v", pos.Name, curPos, err)
				return errors.Trace(err)
			}
			c.master.Update(pos)
			continue
		case *replication.XIDEvent:
			// try to save the position later
			if err := c.eventHandler.OnXID(pos); err != nil {
				return errors.Trace(err)
			}
		case *replication.XAPrepareEvent:
			if err = c.handleXAPrepareEvent(pos, ev); err != nil {
				return errors.Trace(err)
			}
		case *replication.QueryEvent:
			if err = c.handleQueryEvent(pos, ev); err != nil {
				return errors.Trace(err)
			}
			c.master.Update(pos)
			continue
		default:
			if pos.Pos > 0 {
				c.master.Update(pos)
			}
			continue
		}

		c.master.Update(pos)
	}

	return nil
}

func (c *Canal) handleRowsEvent(e *replication.BinlogEvent) error {
	ev := e.Event.(*replication.RowsEvent)

	// Caveat: table may be altered at runtime.
	schema := string(ev.Table.Schema)
	table := string(ev.Table.Table)

	t, err := c.GetTable(schema, table)
	if err != nil {
		return errors.Trace(err)
	}
	var action string
	switch e.Header.EventType {
	case replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
		action = InsertAction
	case replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
		action = DeleteAction
	case replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
		action = UpdateAction
	default:
		return errors.Errorf("%s not supported now", e.Header.EventType)
	}
	events := newRowsEvent(t, action, ev.Rows, BINLOGDATA)
	return c.eventHandler.OnRow(events)
}

func (c *Canal) handleQueryEvent(pos mysql.Position, e *replication.BinlogEvent) error {
	ev := e.Event.(*replication.QueryEvent)

	switch ev.Type {
	case replication.QueryEvent_ALTER:
		if mb := expAlterTable.FindSubmatch(ev.Query); mb != nil {
			if len(mb[1]) == 0 {
				mb[1] = ev.Schema
			}
			c.ClearTableCache(mb[1], mb[2])
			log.Infof("table structure changed, clear table cache: %s.%s\n", mb[1], mb[2])
			log.Infof("get new structure for table: %s.%s\n", mb[1], mb[2])
			_, err := c.GetTable(string(mb[1]), string(mb[2]))
			if err != nil {
				return errors.Trace(err)
			}
		}
	}
	return c.eventHandler.OnDDL(pos, ev)
}

func (c *Canal) handleXAPrepareEvent(pos mysql.Position, e *replication.BinlogEvent) error {
	ev := e.Event.(*replication.XAPrepareEvent)
	queryEv := &replication.QueryEvent{
		Type:  replication.QueryEvent_XA,
		Query: ev.Query,
	}
	return c.eventHandler.OnDDL(pos, queryEv)
}

func (c *Canal) WaitUntilPos(pos mysql.Position, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			return errors.Errorf("wait position %v too long > %s", pos, timeout)
		default:
			curPos := c.master.Position()
			if curPos.Compare(pos) >= 0 {
				return nil
			} else {
				log.Debugf("master pos is %v, wait catching %v", curPos, pos)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	return nil
}

func (c *Canal) CatchMasterPos(timeout time.Duration) error {
	rr, err := c.Execute("SHOW MASTER STATUS")
	if err != nil {
		return errors.Trace(err)
	}

	name, _ := rr.GetString(0, 0)
	pos, _ := rr.GetInt(0, 1)

	return c.WaitUntilPos(mysql.Position{name, uint32(pos)}, timeout)
}
