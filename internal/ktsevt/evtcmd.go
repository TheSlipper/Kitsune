package ktsevt

import (
	"strconv"
	"time"

	"github.com/TheSlipper/Kitsune/internal/ktsdb"
	"github.com/bwmarrin/discordgo"
)

type CmdEvent struct {
	EvtID         int
	ExecutionTime time.Time
	RepeatDelay   time.Duration
	InitMsgDat    discordgo.MessageCreate
	Command       string
	Repeat        bool
}

func (c CmdEvent) Time() time.Time {
	return c.ExecutionTime
}

func (c CmdEvent) ExecFun(s *discordgo.Session) {
	if c.Repeat {
		FreshTxtEvts <- &CmdEvent{c.EvtID, time.Now().Add(c.RepeatDelay), c.RepeatDelay, c.InitMsgDat, c.Command, true}
	}
}

func (c CmdEvent) ExecsCmd() bool {
	return true
}

func (c CmdEvent) CmdData() discordgo.MessageCreate {
	c.InitMsgDat.Message.Content = ktsdb.PrefixGetByServerID(c.InitMsgDat.Message.GuildID) + c.Command
	return c.InitMsgDat
}

func (c CmdEvent) Server() string {
	return c.InitMsgDat.GuildID
}

func (c CmdEvent) String() string {
	return "EID: " + strconv.Itoa(c.EvtID) + "\r\nTxtEvent Type: Command TxtEvent\r\nPlanned time of execution: " + c.ExecutionTime.String() + "\r\nRepeats: " + strconv.FormatBool(c.Repeat) + "\r\nRepeat Delay: " + c.RepeatDelay.String() + "\r\nCommand: " + c.Command
}

func (c CmdEvent) ID() int {
	return c.EvtID
}

func (c *CmdEvent) AssignID(id int) {
	c.EvtID = id
}
