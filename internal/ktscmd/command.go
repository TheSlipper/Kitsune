// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"github.com/bwmarrin/discordgo"
)

// KitsuneCommand defines the structure and behaviour of kitsune's commands
type KitsuneCommand interface {
	Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error
	Name() string
	SetName(s string)
	parseCmd(args []string) error
	execute(m *discordgo.MessageCreate, s *discordgo.Session)
	clearFlags()
}
