// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// help commandName

// Help is an exportable help struct singleton.
var Help help

func init() {
	Help.name = "help"
}

// TODO: Check if u can use Sprintf with just strings

// newLine contains the newline character of the currently running OS (defined by Fprintln go compiler interpretation).
var newLine strings.Builder

// StringFormatterWriter struct that implements io.Writer used to extract the newline character.
type StringFormatterWriter struct {
}

// Write method that writes to the newline character/s to newLine strings.Builder.
func (s StringFormatterWriter) Write(p []byte) (n int, err error) {
	n, err = newLine.Write(p)
	return n, err
}

func init() {
	var formatter StringFormatterWriter
	fmt.Fprintln(formatter, "")
}

// Help is a struct that contains help command information.
type help struct {
	name      string
	targetCmd string
}

// Run runs the command with the passed command data.
func (h *help) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	h.clearFlags()
	h.parseCmd(args)
	h.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (h *help) Name() string {
	return h.name
}

// SetName sets the name of the command.
func (h *help) SetName(s string) {
	h.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (h *help) parseCmd(args []string) error {
	if len(args) < 2 {
		h.targetCmd = "help"
		return nil
	}
	h.targetCmd = args[1]
	return nil
}

// execute runs the command and executes the desired tasks.
func (h *help) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	dat, err := ioutil.ReadFile("mans/" + h.targetCmd)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, h.name+": no such command like '"+h.targetCmd+"'")
		return
	}

	var bdy strings.Builder
	var hdr string
	nl := newLine.String()
	for i, str := range strings.Split(string(dat), nl) {
		if i == 0 {
			hdr = str
		} else {
			bdy.WriteString(str + nl)
		}
	}

	sendRawEmbed(s, m.ChannelID, hdr, bdy.String())
}

// clearFlags clears data after the previous command execution.
func (h *help) clearFlags() {
	h.targetCmd = ""
}
