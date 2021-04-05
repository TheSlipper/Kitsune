package ktscmd

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ytplay -url "https://www.youtube.com/watch?v=14cgH6MiYf8" -target-channel "chan name"
// ytplay -tags "super mario bros 2 theme" -target-channel "chan name"
// ytplay -clear-queue

// YTPlay is an exportable ytPlay struct singleton.
var YTPlay ytPlay

const tempFilesPathPrefix = "temp/audio/"

func init() {
	YTPlay.name = "ytplay"
	YTPlay.saveFolderPath = "/temp/audio"
}

// ytPlay is a struct that contains ytPlay command information.
type ytPlay struct {
	name           string
	url            string
	tags           string
	targetChan     string
	saveFolderPath string
	clear          bool
	s              *discordgo.Session
}

// Run runs the command with the passed command data.
func (yt *ytPlay) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	if yt.s == nil {
		yt.s = s
	}
	yt.clearFlags()
	err := yt.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, yt.name+": "+err.Error())
	}

	yt.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (yt *ytPlay) Name() string {
	return yt.name
}

// SetName sets the name of the command.
func (yt *ytPlay) SetName(s string) {
	yt.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (yt *ytPlay) parseCmd(args []string) (err error) {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-url":
				yt.url, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-tags":
				yt.tags, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-target-channel":
				yt.targetChan, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-clear-queue":
				yt.clear = true
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
func (yt *ytPlay) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	// songID := strings.Replace(yt.url, "https://youtu.be/", "", 1)
	// songID = strings.Replace(songID, "https://www.youtube.com/", "", 1)
	// go downloadSongFromUrl(yt.url, m.GuildID, songID)

	// var channelID string
	// g, _ := s.Guild(m.GuildID)
	// for _, ch := range g.Channels {
	// 	if ch.Name == yt.targetChan && ch.Type == discordgo.ChannelTypeGuildVoice {
	// 		channelID = ch.ID
	// 	}
	// }

	// ktsevt.FreshVcEvents <- &ytEntryPlayEvent{ktsevt.GetEID(), false, channelID,
	// 	songID, m.GuildID, nil, false}
	s.ChannelMessageSend(m.ChannelID, "command has been deactivated due to constant crashes. It will be rewritten in "+
		"the future.")
}

// clearFlags clears data after the previous command execution.
func (yt *ytPlay) clearFlags() {
	yt.url = ""
	yt.tags = ""
	yt.targetChan = ""
	yt.clear = false
}
