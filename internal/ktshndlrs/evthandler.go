// Package ktshndlrs provides functions that process commands and events.
package ktshndlrs

import (
	"time"

	"github.com/TheSlipper/Kitsune/internal/ktsevt"
)

// BotEventHandler is a handler that processes all of the scheduled events.
func BotEventHandler() {
	for {
		ktsevt.ProcessTxtChanEvents()
		ktsevt.ProcessNewVcEvts()
		// process text events
		for _, txtEvtSlice := range ktsevt.ServerTxtEvents {
			for i := 0; i < len(txtEvtSlice); i++ {
				if txtEvtSlice[i] != nil && txtEvtSlice[i].Time().Before(time.Now()) {
					txtEvtSlice[i].ExecFun(ktsevt.Session)
					if txtEvtSlice[i].ExecsCmd() {
						m := txtEvtSlice[i].CmdData()
						go MsgHandler(ktsevt.Session, &m)
					}
					txtEvtSlice[i] = nil
				}
			}
		}
		// process voice chat events
		for sID, vcQueue := range ktsevt.VcEvents {
			if !vcQueue.IsEventPlaying() && vcQueue.GetQueueEvt(0) != nil {
				go ktsevt.VcEvents[sID].StartEvt()
			}
		}
		// timeout for 1 second to put less stress on the CPU
		ktsevt.EvtBreak = true
		time.Sleep(time.Second)
		ktsevt.EvtBreak = false
	}
}
