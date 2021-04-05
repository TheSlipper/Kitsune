// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"strings"
	"time"

	"github.com/TheSlipper/Kitsune/internal/ktsevt"
	"github.com/bwmarrin/discordgo"
)

// at -time "+3h20m10s" -command "ping" -repeat "1m"

// At is an exportable at struct singleton.
var At at

func init() {
	At.name = "at"
}

// At is a struct that contains at command information.
type at struct {
	actionDelay    time.Duration
	repeatDelay    time.Duration
	delayedCommand string
	name           string
}

// Run runs the command with the passed command data.
func (a *at) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	err := a.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, a.name+": "+err.Error())
		return nil
	}

	a.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (a *at) Name() string {
	return a.name
}

// SetName sets the name of the command.
func (a *at) SetName(s string) {
	a.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (a *at) parseCmd(args []string) (err error) {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-time":
				var foo string
				foo, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
				a.actionDelay, err = time.ParseDuration(strings.TrimPrefix(foo,
					"+"))
				if err != nil {
					return err
				}
			case "-command":
				a.delayedCommand, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-repeat":
				var foo string
				foo, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
				a.repeatDelay, err = time.ParseDuration(strings.TrimPrefix(foo,
					"+"))
				if err != nil {
					return err
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
func (a *at) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	ktsevt.FreshTxtEvts <- &ktsevt.CmdEvent{ktsevt.GetEID(), time.Now().Add(a.actionDelay), a.repeatDelay, *m, a.delayedCommand, a.repeatDelay != 0}
}

// clearFlags clears data after the previous command execution.
func (a *at) clearFlags() {
	a.actionDelay = 0
	a.repeatDelay = 0
	a.delayedCommand = ""
}
