// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// clear

// Clear is an exportable clear struct singleton.
var Clear clear

func init() {
	Clear.name = "clear"
}

// Clear is a struct that contains clear command information.
type clear struct {
	name string
}

// Run runs the command with the passed command data.
func (c *clear) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	c.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (c *clear) Name() string {
	return c.name
}

// SetName sets the name of the command.
func (c *clear) SetName(s string) {
	c.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (c *clear) parseCmd(args []string) error {
	return nil
}

// execute runs the command and executes the desired tasks.
func (c *clear) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	var sb strings.Builder
	for i := 0; i < 50; i++ {
		sb.WriteString("__\r\n")
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, sb.String())
}

// clearFlags clears data after the previous command execution.
func (c *clear) clearFlags() {
}
