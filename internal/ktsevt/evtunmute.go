package ktsevt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ChatUnmuteEvt struct {
	EvtID      int
	UnmuteTime time.Time
	UsrID      string
	GldID      string
	MuteRoleID string
}

func (t *ChatUnmuteEvt) Time() time.Time {
	return t.UnmuteTime
}

func (t *ChatUnmuteEvt) ExecFun(s *discordgo.Session) {
	err := s.GuildMemberRoleRemove(t.GldID, t.UsrID, t.MuteRoleID)
	if err != nil {
		fmt.Printf("[ChatUnmuteEvt - Unexpected error occurred while removing mute role from a user: %s]", err)
		panic(err.Error())
	}
}

func (t *ChatUnmuteEvt) ExecsCmd() bool {
	return false
}

func (t ChatUnmuteEvt) CmdData() discordgo.MessageCreate {
	return discordgo.MessageCreate{}
}

func (t ChatUnmuteEvt) Server() string {
	return t.GldID
}

func (t ChatUnmuteEvt) String() string {
	return "EID: " + strconv.Itoa(t.EvtID) + "\r\nTxtEvent Type: Chat Unmute TxtEvent\r\nPlanned time of execution: " + t.UnmuteTime.String() + "\r\nTarget user ID: " + t.UsrID
}

func (t ChatUnmuteEvt) ID() int {
	return t.EvtID
}

func (t *ChatUnmuteEvt) AssignID(id int) {
	t.EvtID = id
}
