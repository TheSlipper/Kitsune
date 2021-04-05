package ktsevt

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

var FreshVcEvents chan VoiceEvent
var VcEvents map[string]*voiceChanQueue
var evtCounter int

type VoiceEvent interface {
	ExecFun(*discordgo.Session)
	Server() string
	Channel() string
	EID() int
	Kill()
	String() string
}

type voiceChanQueue struct {
	playing bool
	events []VoiceEvent
}

func (v *voiceChanQueue) IsEventPlaying() bool {
	return v.playing
}

func (v *voiceChanQueue) GetQueueEvt(i int) VoiceEvent {
	return v.events[i]
}

func (v *voiceChanQueue) EvtOut() {
	v.events = append(v.events[1:], nil)
}

func (v *voiceChanQueue) StartEvt() {
	v.playing = true
	v.events[0].ExecFun(Session)
	v.EvtOut()
	v.playing = false
}

func NewVoiceChanQueue() *voiceChanQueue {
	return &voiceChanQueue{false, make([]VoiceEvent, 10)}
}

func init() {
	FreshVcEvents = make(chan VoiceEvent, 10)
	VcEvents = make(map[string]*voiceChanQueue)
	evtCounter = 0
}

func ProcessNewVcEvts() {
	select {
	case evt, ok := <-FreshVcEvents:
		if ok {
			exists := false
			for sID, q := range VcEvents {
				if sID == evt.Server() {
					exists = true
					for i := 0; i < len(q.events); i++ {
						if q.events[i] == nil {
							q.events[i] = evt
							break
						}
					}
				}
			}
			if !exists {
				VcEvents[evt.Server()] = NewVoiceChanQueue()
				VcEvents[evt.Server()].events[0] = evt
			}
			ProcessNewVcEvts()
		}
	default:
		return
	}
}

func GetScheduledVcEvents(serverID string) []VoiceEvent {
	for sID, evtQueue := range VcEvents {
		if sID == serverID {
			return evtQueue.events
		}
	}
	return nil
}

func KillVcEvent(sID string, eid int) {
	evts := GetScheduledVcEvents(sID)
	for i := 0; i < len(evts); i++ {
		if evts[i].EID() == eid && evts[i].Server() == sID {
			if !EvtBreak {
				time.Sleep(time.Second/2)
				continue
			}
			evts[i].Kill()
			break
		}
	}
}

