// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"strings"

	"github.com/TheSlipper/Kitsune/internal/ktsdb"
	"github.com/bwmarrin/discordgo"
)

// chngprefix "p"

// ChngPrefix is an exportable chngPrefix struct singleton.
var ChngPrefix chngPrefix

func init() {
	ChngPrefix.name = "chngprefix"
}

// ChngPrefix is a struct that contains chngprefix command information.
type chngPrefix struct {
	name      string
	newPrefix string
}

// Run runs the command with the passed command data.
func (ch *chngPrefix) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	ch.clearFlags()
	ch.parseCmd(args)
	if ch.newPrefix == "" {
		_, _ = s.ChannelMessageSend(m.ChannelID, ch.name+": no new prefix specified!")
	} else if len(ch.newPrefix) > 1 {
		_, _ = s.ChannelMessageSend(m.ChannelID, ch.name+": prefix should be a single character! *(e.g.: the dollar character - \"$\"")
	}
	ch.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (ch *chngPrefix) Name() string {
	return ch.name
}

// SetName sets the name of the command.
func (ch *chngPrefix) SetName(s string) {
	ch.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (ch *chngPrefix) parseCmd(args []string) error {
	if len(args) < 2 {
		return nil
	}
	ch.newPrefix = args[1]
	ch.newPrefix = strings.Replace(ch.newPrefix, "\"", "", -1)
	return nil
}

// execute runs the command and executes the desired tasks.
func (ch *chngPrefix) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	if !ktsdb.PrefixChange(ch.newPrefix, m.GuildID) {
		// fmt.Printf("["+ch.name+" - error while changing the server's prefix: %s\r\n]", err)
		// messageErrLog(err, m)
		_, _ = s.ChannelMessageSend(m.ChannelID, ch.name+": unexpected error occurred while trying to change the server's prefix")
	}
}

// clearFlags clears data after the previous command execution.
func (ch *chngPrefix) clearFlags() {
	ch.newPrefix = ""
}
