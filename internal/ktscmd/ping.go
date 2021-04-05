// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"github.com/bwmarrin/discordgo"
)

// Ping is an exportable ping struct singleton.
var Ping ping

func init() {
	Ping.name = "ping"
}

// Ping is a struct that contains ping command information.
type ping struct {
	name string
}

// Run runs the command with the passed command data.
func (p *ping) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	p.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (p *ping) Name() string {
	return p.name
}

// SetName sets the name of the command.
func (p *ping) SetName(s string) {
	p.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (p *ping) parseCmd(args []string) error {
	return nil
}

// execute runs the command and executes the desired tasks.
func (p *ping) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "pong!")
}

// clearFlags clears data after the previous command execution.
func (p *ping) clearFlags() {
}
