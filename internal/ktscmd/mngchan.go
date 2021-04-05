// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TheSlipper/Kitsune/internal/ktsdb"
	"github.com/bwmarrin/discordgo"
)

// mngchan -create-channel "channel name" -type "text"
// mngchan -remove-channel #tagged-channel
// mngchan -whitelist-channel #tagged-channel
// mngchan -blacklist-channel #tagged-channel

// MngChan is an exportable mngChan struct singleton.
var MngChan mngChan

func init() {
	MngChan.name = "mngchan"
}

// MngChan is a struct that contains mngchan command information.
type mngChan struct {
	name          string
	targetID      []string
	newChanType   string
	newChanName   string
	createChan    bool
	rmChan        bool
	whitelistChan bool
	blacklistChan bool
}

// Run runs the command with the passed command data.
func (mng *mngChan) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	mng.clearFlags()
	hasPerm, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionAll)
	if err != nil || !hasPerm {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": you need to be an administrator of the server in order to run this command!")
		return nil
	}

	err = mng.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, mng.name+": "+err.Error())
		return nil
	}

	if mng.newChanType == "voice" {
		mng.proccessVoiceChanNames(m.GuildID, s)
	}
	mng.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (mng *mngChan) Name() string {
	return mng.name
}

// SetName sets the name of the command.
func (mng *mngChan) SetName(s string) {
	mng.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (mng *mngChan) parseCmd(args []string) (err error) {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-create-channel":
				mng.createChan = true
				mng.newChanName, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-remove-channel":
				mng.rmChan = true
			case "-whitelist-channel":
				mng.whitelistChan = true
			case "-blacklist-channel":
				mng.blacklistChan = true
			case "-type":
				mng.newChanType, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			default:
				return errors.New("Unkown argument '" + args[i] + "'")
			}
		} else if strings.HasPrefix(args[i], "<#") && strings.HasSuffix(args[i], ">") {
			args[i] = strings.Replace(args[i], "<#", "", -1)
			args[i] = strings.Replace(args[i], ">", "", -1)
			mng.targetID = append(mng.targetID, args[i])
		} else {
			return errors.New("No such command value/flag as: " + args[i])
		}
	}

	return nil
}

// execute runs the command and executes the desired tasks.
func (mng *mngChan) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	if mng.createChan && !(mng.rmChan || mng.blacklistChan || mng.whitelistChan) {
		mng.createChannel(m, s)
	} else if mng.rmChan && !(mng.createChan || mng.blacklistChan || mng.whitelistChan) {
		mng.removeChannel(m, s)
	} else if mng.blacklistChan && !(mng.rmChan || mng.createChan || mng.whitelistChan) {
		mng.addToBlacklist(m, s)
	} else if mng.whitelistChan && !(mng.rmChan || mng.blacklistChan || mng.createChan) {
		mng.removeFromBlacklist(m, s)
	} else {
		s.ChannelMessageSend(m.ChannelID, mng.name+": Actions specified incorrectly - you must select **only one** action!")
	}
}

// clearFlags clears data after the previous command execution.
func (mng *mngChan) clearFlags() {
	mng.targetID = mng.targetID[:0]
	mng.createChan = false
	mng.rmChan = false
	mng.whitelistChan = false
	mng.blacklistChan = false
	mng.newChanType = ""
	mng.newChanName = ""
}

// proccessVoiceChanNames processes if the specified target is a voice channel.
func (mng *mngChan) proccessVoiceChanNames(gID string, s *discordgo.Session) {
	for i := 0; i < len(mng.targetID); i++ {
		if strings.HasPrefix(mng.targetID[i], "#") {
			mng.targetID[i] = strings.Replace(mng.targetID[i], "#", "", -1)
			g, err := s.Guild(gID)
			if err != nil {
				fmt.Printf("["+mng.name+" - error while getting a server when processing voice channel names: '%s'\r\n", err)
				panic(err.Error())
			}
			for _, ch := range g.Channels {
				if ch.Name == mng.targetID[i] {
					mng.targetID[i] = ch.ID
				}
			}
		}
	}
}

//  createChannel creates a discord channel on a given server with the specified data.
func (mng *mngChan) createChannel(m *discordgo.MessageCreate, s *discordgo.Session) {
	var err error
	if mng.newChanType == "text" {
		_, err = s.GuildChannelCreate(m.GuildID, mng.newChanName, discordgo.ChannelTypeGuildText)
	} else if mng.newChanType == "voice" {
		_, err = s.GuildChannelCreate(m.GuildID, mng.newChanName, discordgo.ChannelTypeGuildVoice)
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": no such channel type like '"+mng.newChanType+"'")
		return
	}
	if err != nil {
		fmt.Printf("["+mng.name+" - error while creating guild channel \""+mng.newChanName+"\": %s", err)
	}
}

// removeChannel removes the specified channel on a given server.
func (mng *mngChan) removeChannel(m *discordgo.MessageCreate, s *discordgo.Session) {
	for _, ch := range mng.targetID {
		_, err := s.ChannelDelete(ch)
		if err != nil {
			fmt.Printf("["+mng.name+" - Error while deleting a channel: '%s']\r\n", err)
			panic(err.Error())
		}
	}
}

// addToBlacklist adds targeted channel to the list of blacklisted channels.
func (mng *mngChan) addToBlacklist(m *discordgo.MessageCreate, s *discordgo.Session) {
	for _, str := range mng.targetID {
		ktsdb.ChanBlacklist(str, m.GuildID)
	}
}

// removeFromBlacklist removes the targeted channel from the blacklisted channels list.
func (mng *mngChan) removeFromBlacklist(m *discordgo.MessageCreate, s *discordgo.Session) {
	for _, str := range mng.targetID {
		ktsdb.ChanWhitelist(str)
	}
}
