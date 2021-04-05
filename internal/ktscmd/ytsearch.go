// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"github.com/TheSlipper/Kitsune/internal/settings"
	"github.com/bwmarrin/discordgo"
)

// YTSearch is an exportable ytSearch struct singleton.
var YTSearch ytSearch

func init() {
	YTSearch.name = "ytsearch"
	YTSearch.client = &http.Client{
		Transport: &transport.APIKey{Key: settings.BotSettings.YTToken},
	}

	var err error
	YTSearch.service, err = youtube.New(YTSearch.client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}
}

// TODO: Standardize the style of commands (for example '=' crashes the program)

// YTSearch is a struct that contains ytsearch command information.
type ytSearch struct {
	name        string
	tags        string
	contentType string
	results     int
	client      *http.Client // TODO: Init this somewhere
	service     *youtube.Service
}

// Run runs the command with the passed command data.
func (yt *ytSearch) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	yt.clearFlags()
	err := yt.parseCmd(args)
	if err != nil {
		return err
	}
	yt.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (yt *ytSearch) Name() string {
	return yt.name
}

// SetName sets the name of the command.
func (yt *ytSearch) SetName(s string) {
	yt.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (yt *ytSearch) parseCmd(args []string) error {
	for i := 1; i < len(args); i++ {
		if strings.Contains(args[i], "content-type") || strings.HasPrefix(args[i], "-c") {
			i = i + 1
			yt.contentType = args[i]
		} else if strings.Contains(args[i], "tags") || strings.HasPrefix(args[i], "-t") {
			i = i + 1
			// TODO:
			yt.tags, _ = extractInsideCharsAndIterate("\"", &i, args)
		} else if strings.Contains(args[i], "results") || strings.HasPrefix(args[i], "-r") {
			i = i + 1
			var err error
			yt.results, err = strconv.Atoi(args[i])
			if err != nil {
				return err
			}
		}
	}

	if yt.results > 5 {
		yt.results = 5
	}
	return nil
}

// execute runs the command and executes the desired tasks.
func (yt *ytSearch) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	//call := yt.service.Search.List("id,snippet").Q(yt.tags).MaxResults(int64(yt.results))
	call := yt.service.Search.List([]string{"id", "snippet"}).Q(yt.tags).MaxResults(int64(yt.results))
	response, err := call.Do()
	// handleError(err, "")
	if err != nil {
		panic(err.Error())
	}

	var sb strings.Builder

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			if yt.contentType == "video" {
				sb.WriteString("https://www.youtube.com/watch?v=" + item.Id.VideoId + "\r\n")
			}
			break
		case "youtube#channel":
			if yt.contentType == "channel" { // TODO: Test if this URL is correct
				sb.WriteString("https://www.youtube.com/channel/" + item.Id.ChannelId + "\r\n")
			}
			break
		case "youtube#playlist":
			if yt.contentType == "playlist" {
				sb.WriteString("https://www.youtube.com/playlist?list=" + item.Id.PlaylistId + "\r\n")
			}
			break
		}
	}

	if sb.String() != "" {
		s.ChannelMessageSend(m.ChannelID, sb.String())
	} else {
		s.ChannelMessageSend(m.ChannelID, "ytsearch: No valid search result for type '"+yt.contentType+"' in pool of "+strconv.Itoa(yt.results)+" search results")
	}
}

// clearFlags clears data after the previous command execution.
func (yt *ytSearch) clearFlags() {
	yt.tags = ""
	yt.contentType = ""
	yt.results = 1
}
