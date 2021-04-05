// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/TheSlipper/Kitsune/internal/ktsdb"

	"github.com/TheSlipper/Kitsune/internal/ktsevt"

	"github.com/bwmarrin/discordgo"
)

// List is an exportable list struct singleton.
var List list

func init() {
	List.name = "list"
}

// List is a struct that contains list command information.
type list struct {
	lsTarget string
	name     string
}

// Run runs the command with the passed command data.
func (l *list) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	l.clearFlags()
	err := l.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, l.name+": "+err.Error())
	}
	l.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (l *list) Name() string {
	return l.name
}

// SetName sets the name of the command.
func (l *list) SetName(s string) {
	l.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (l *list) parseCmd(args []string) error {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			l.lsTarget = strings.TrimPrefix(args[i], "-")
		} else {
			return errors.New("No such command value/flag as: " + args[i])
		}
	}
	return nil
}

// execute runs the command and executes the desired tasks.
func (l *list) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	if l.lsTarget == "txt-events" {
		l.sendTxtEvtData(s, m.GuildID, m.ChannelID)
	} else if l.lsTarget == "vc-events" {
		l.sendVcEvtData(s, m.GuildID, m.ChannelID)
	} else if l.lsTarget == "aliases" {
		l.sendAliasData(s, m.GuildID, m.ChannelID)
	} else if l.lsTarget == "blacklists" {
		l.sendBlacklistData(s, m.GuildID, m.ChannelID)
	} else {
		s.ChannelMessageSend(m.ChannelID, l.name+": Incorrect listing target '"+l.lsTarget+"'!")
	}
}

// clearFlags clears data after the previous command execution.
func (l *list) clearFlags() {
	l.lsTarget = ""
}

// sendTxtEvtData sends the list of the specified server's text channel events.
func (l *list) sendTxtEvtData(s *discordgo.Session, sID, chID string) {
	evts := ktsevt.GetScheduledTxtEvents(sID)
	headers, bodies := make([]string, len(evts)), make([]string, len(evts))
	for i, evt := range evts {
		if evt == nil {
			continue
		}
		headers[i] = "TxtEvent " + strconv.Itoa(i+1)
		bodies[i] = evt.String()
	}
	sendRawEmbedMultiple(s, chID, "Text Chat Data Listing", headers, bodies)
}

// sendVcEvtData sends the list of the specified server's voice channel events.
func (l *list) sendVcEvtData(s *discordgo.Session, sID, chID string) {
	evts := ktsevt.GetScheduledVcEvents(sID)
	headers, bodies := make([]string, len(evts)), make([]string, len(evts))
	for i, evt := range evts {
		if evt == nil {
			continue
		}
		headers[i] = "Voice Channel Event (EID:" + strconv.Itoa(evt.EID()) + ")"
		bodies[i] = evt.String()
	}
	sendRawEmbedMultiple(s, chID, "Voice Chat Event Data Listing", headers, bodies)
}

// sendAliasData sends the list of the specified server's aliases.
func (l *list) sendAliasData(s *discordgo.Session, sID string, chID string) {
	aliases, err := ktsdb.AliasGetByServerID(sID)
	if err != nil {
		_, _ = s.ChannelMessageSend(chID, l.name+": unexpected database error occurred - could not get aliases!")
		return
	}
	headers, bodies := make([]string, len(aliases)), make([]string, len(aliases))
	for i, a := range aliases {
		headers[i] = "Alias '" + a.AliasReplacement + "'"
		bodies[i] = "\r\nAlias ID: " + strconv.FormatInt(a.AliasID, 10) + "\r\nCreator ID: " + strconv.FormatInt(a.UserID, 10) + "\r\nAlias Content: " + a.AliasContent + "\r\n"
	}
	sendRawEmbedMultiple(s, chID, "Alias Data Listing", headers, bodies)
}

// sendBlacklistData sends the list of the specified server's blacklisted channels.
func (l *list) sendBlacklistData(s *discordgo.Session, sID string, chID string) {
	blacklists := ktsdb.ChanGetBlacklistedByServer(sID)
	headers, bodies := make([]string, len(blacklists)), make([]string, len(blacklists))
	for i, str := range blacklists {
		channel, err := s.Channel(str)
		if err != nil {
			fmt.Printf("["+l.name+" - Error while creating a channel struct!: %s]\r\n", err)
			panic(err)
		}
		headers[i] = "Channel '" + channel.Name + "'"
		bodies[i] = "ID: " + channel.ID
	}
	sendRawEmbedMultiple(s, chID, "Blacklisted Channels Listing", headers, bodies)
}
