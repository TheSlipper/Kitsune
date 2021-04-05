// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"strconv"
	"strings"

	"github.com/TheSlipper/Kitsune/internal/ktsdb"

	"github.com/bwmarrin/discordgo"
)

// alias -add-alias -alias-name "yeet" -alias-content "mngusr kick blablabla"
// alias -remove-alias -alias-name "yeet"

// Alias is an exportable alias struct singleton.
var Alias alias

func init() {
	Alias.name = "alias"
}

// Alias is a struct that contains alias command information.
type alias struct {
	name      string
	add       bool
	remove    bool
	aliasData ktsdb.AliasData
}

// Run runs the command with the passed command data.
func (a *alias) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	hasPerm, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAll)
	if err != nil || !hasPerm {
		_, _ = s.ChannelMessageSend(m.ChannelID, a.name+": you need to be an administrator of the server in order to run this command!")
		return nil
	}
	a.clearFlags()
	err = a.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, a.name+": "+err.Error())
		return nil
	}

	a.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (a *alias) Name() string {
	return a.name
}

// SetName sets the name of the command.
func (a *alias) SetName(s string) {
	a.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (a *alias) parseCmd(args []string) (err error) {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-add-alias":
				a.add = true
			case "-remove-alias":
				a.remove = true
			case "-alias-name":
				a.aliasData.AliasReplacement, err = extractArgVal(&args, &i)
				if err != nil {
					return errors.New("incorrect alias replacement entered")
				}
			case "-alias-content":
				a.aliasData.AliasContent, err = extractArgVal(&args, &i)
				if err != nil {
					return errors.New("incorrect alias content entered")
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
func (a *alias) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	if a.add && a.remove {
		s.ChannelMessageSend(m.ChannelID, a.name+": You cannot add and remove an alias in one command!")
	} else if a.add {
		a.createAlias(m, s)
	} else if a.remove {
		ktsdb.AliasDelete(a.aliasData.AliasReplacement, m.GuildID)
	}
}

// clearFlags clears data after the previous command execution.
func (a *alias) clearFlags() {
	a.aliasData.Reset()
	a.add = false
	a.remove = false
}

// createAlias puts the alias in the kitsune database with the use of ktsdb package.
func (a *alias) createAlias(m *discordgo.MessageCreate, s *discordgo.Session) {
	placeHolder, _ := strconv.ParseInt(m.GuildID, 10, 64)
	a.aliasData.ServerID = placeHolder
	placeHolder, _ = strconv.ParseInt(m.Author.ID, 10, 64)
	a.aliasData.UserID = placeHolder

	if ktsdb.AliasExists(a.aliasData.AliasReplacement, m.GuildID) {
		_, _ = s.ChannelMessageSend(m.ChannelID, a.name+": alias '"+a.aliasData.AliasReplacement+"' is already taken. Please delete it if you want to assign something else to it.")
		return
	} else if a.aliasData.AliasContent == "" || a.aliasData.AliasReplacement == "" {
		s.ChannelMessageSend(m.ChannelID, a.name+": alias content or alias replacement not specified!")
		return
	}
	if err := ktsdb.AliasAdd(a.aliasData); err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, a.name+": unexpected database error occurred - could not add '"+a.aliasData.AliasReplacement+"' alias.")
	}
}
