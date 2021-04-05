// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"strconv"
	"strings"

	"github.com/TheSlipper/Kitsune/internal/ktsevt"

	"github.com/bwmarrin/discordgo"
)

// kill -eid 1

// Kill is an exportable kill struct singleton.
var Kill kill

func init() {
	Kill.name = "kill"
}

// kill is a struct that contains kill command information.
type kill struct {
	EID    int
	action string
	name   string
}

// Run runs the command with the passed command data.
func (k *kill) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	k.clearFlags()
	err := k.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, k.name+": "+err.Error())
		return nil
	} else if k.EID < 0 {
		s.ChannelMessageSend(m.ChannelID, k.name+": The specified EID '"+strconv.Itoa(k.EID)+"' is incorrect")
		return nil
	}

	k.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (k *kill) Name() string {
	return k.name
}

// SetName sets the name of the command.
func (k *kill) SetName(s string) {
	k.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (k *kill) parseCmd(args []string) (err error) {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-eid":
				i++
				k.EID, err = strconv.Atoi(args[i])
				if err != nil {
					return errors.New("The specified EID '" + args[i] + "' is incorrect")
				}
			default:
				return errors.New("Unkown argument '" + args[i] + "'")
			}
		} else {
			return errors.New("No such command value/flag as: " + args[i])
		}
	}
	return nil
}

// execute runs the command and executes the desired tasks.
func (k *kill) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	if !k.hasEvent(m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, k.name+": no event with EID '"+strconv.Itoa(k.EID)+"' found")
		return
	}
	ktsevt.KillTxtEvent(m.GuildID, k.EID)
	ktsevt.KillVcEvent(m.GuildID, k.EID)
}

// clearFlags clears data after the previous command execution.
func (k *kill) clearFlags() {
	k.EID = -1
	k.action = ""
}

// hasEvent checks if the given server has an event with the specified event ID.
func (k *kill) hasEvent(sID string) bool {
	txtEvts, vcEvts := ktsevt.GetScheduledTxtEvents(sID), ktsevt.GetScheduledVcEvents(sID)
	for _, evt := range txtEvts {
		if evt.ID() == k.EID && evt.Server() == sID {
			return true
		}
	}
	for _, evt := range vcEvts {
		if evt.EID() == k.EID && evt.Server() == sID {
			return true
		}
	}
	return false
}
