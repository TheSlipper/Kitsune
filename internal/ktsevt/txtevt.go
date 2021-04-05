package ktsevt

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

var FreshTxtEvts chan TxtEvent
var ServerTxtEvents map[string][]TxtEvent
var Session *discordgo.Session
var txtEvtCounter int
var EvtBreak bool

func init() {
	ServerTxtEvents = make(map[string][]TxtEvent)
	FreshTxtEvts = make(chan TxtEvent, 10)
	txtEvtCounter = 0
	EvtBreak = false
}

type TxtEvent interface {
	Time() time.Time
	ExecFun(*discordgo.Session)
	ExecsCmd() bool
	CmdData() discordgo.MessageCreate
	Server() string
	String() string
	ID() int
	AssignID(int)
}

func ProcessTxtChanEvents() {
	select {
	case evt, ok := <-FreshTxtEvts:
		if ok {
			added := false
			for i := 0; i < len(ServerTxtEvents[evt.Server()]); i++ {
				if ServerTxtEvents[evt.Server()][i] == nil {
					ServerTxtEvents[evt.Server()][i] = evt
					added = true
					break
				}
			}
			if !added {
				ServerTxtEvents[evt.Server()] = append(ServerTxtEvents[evt.Server()], evt)
			}
			ProcessTxtChanEvents()
		}
	default:
		return
	}
}

func GetScheduledTxtEvents(sID string) []TxtEvent {
	for index, evtList := range ServerTxtEvents {
		if index == sID {
			return evtList
		}
	}
	return nil
}

func KillTxtEvent(sID string, EID int) {
	evts := GetScheduledTxtEvents(sID)
	for i := 0; i < len(evts); i++ {
		if evts[i].ID() == EID && evts[i].Server() == sID {
			for {
				// TODO: Add an event for killing the events
				if !EvtBreak {
					// Divided by two so that ticking of killing and event handling aren't in sync
					time.Sleep(time.Second / 2)
					continue
				}
				evts[i] = nil
				break
			}
		}
	}
}

func GetEID() int {
	val := txtEvtCounter + 1
	txtEvtCounter++
	return val
}
