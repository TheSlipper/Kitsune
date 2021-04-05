// Package ktshndlrs provides functions that process commands and events.
package ktshndlrs

import (
	"fmt"
	"strings"

	"github.com/TheSlipper/Kitsune/internal/ktscmd"

	"github.com/TheSlipper/Kitsune/internal/ktsdb"
	"github.com/bwmarrin/discordgo"
)

var (
	cmds [15]ktscmd.KitsuneCommand
)

func init() {
	// TODO: Unit tests for commands
	cmds[0] = &ktscmd.Alias
	cmds[1] = &ktscmd.Art
	cmds[2] = &ktscmd.At
	cmds[3] = &ktscmd.ChngPrefix
	cmds[4] = &ktscmd.Clear
	cmds[5] = &ktscmd.Help
	cmds[6] = &ktscmd.Kill
	cmds[7] = &ktscmd.List
	cmds[8] = &ktscmd.MngChan
	cmds[9] = &ktscmd.MngRole
	cmds[10] = &ktscmd.MngUsr
	cmds[11] = &ktscmd.Ping
	cmds[12] = &ktscmd.YTSearch
	cmds[13] = &ktscmd.Honk
	cmds[14] = &ktscmd.YTPlay
}

// processAliases processes aliases of a given server.
func processAliases(cmd string, sID string) string {
	content := cmd
	aliases, err := ktsdb.AliasGetByServerID(sID)
	if err != nil {
		return content
	}
	for _, alias := range aliases {
		if strings.Contains(cmd, alias.AliasReplacement) {
			content = strings.Replace(content, alias.AliasReplacement, alias.AliasContent, -1)
		}
	}

	return content
}

// MsgHandler is a handler that processes the recieved messages.
func MsgHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := ktsdb.PrefixGetByServerID(m.GuildID)
	if m.Author.ID == s.State.User.ID || ktsdb.ChanBlacklisted(m.ChannelID) || !isASCII(m.Content) || !strings.HasPrefix(m.Content, prefix) {
		return
	}
	var cmd string
	if !strings.HasPrefix(m.Content, prefix+"alias") { // Prefix + alias command
		cmd = processAliases(m.Content, m.GuildID)
	} else {
		cmd = m.Content
	}
	tokens := strings.Split(cmd, " ")
	tokens[0] = strings.Replace(tokens[0], prefix, "", -1)
	for _, cmd := range cmds {
		if cmd.Name() == tokens[0] {
			err := cmd.Run(tokens, m, s)
			if err != nil {
				fmt.Printf("[Unexpected error occurred: %s]\r\n", err)
			}
			return
		}
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, "Incorrect command. '"+tokens[0]+"' is not a command!")
}
