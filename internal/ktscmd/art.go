// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/TheSlipper/Kitsune/internal/settings"

	"github.com/TheSlipper/Kitsune/pkg/ktsart"
	"github.com/bwmarrin/discordgo"
)

// art -gallery gelbooru -random -tags "tag1 tag2" -amount 5 -uncompressed

// Art is an exportable art struct singleton.
var Art art

func init() {
	// on init initiate connections to all of the galleries
	Art.name = "art"
	var clients [6]http.Client

	// danbooru
	dbAuth := ktsart.DBAuth{User: settings.BotSettings.DanbooruLogin, Hash: settings.BotSettings.DanbooruToken}
	Art.dbClient = ktsart.NewDB(&clients[0], settings.BotSettings.Addresses.DanbooruURL, &dbAuth)

	// gelbooru
	gbAuth := ktsart.GBAuth{User: settings.BotSettings.GelbooruUsrID, Hash: settings.BotSettings.GelbooruToken}
	Art.gbClient = ktsart.NewGB(&clients[1], settings.BotSettings.Addresses.GelbooruURL, &gbAuth)

	// konachan
	Art.kClient = ktsart.NewKB(&clients[2], settings.BotSettings.Addresses.KonachanURL)

	// rule34
	Art.r34Client = ktsart.NewR34(&clients[3], settings.BotSettings.Addresses.Rule34URL)

	// safebooru
	Art.sbClient = ktsart.NewSB(&clients[4], settings.BotSettings.Addresses.SafebooruURL)

	// yandere
	Art.yaClient = ktsart.NewYA(&clients[5], settings.BotSettings.Addresses.YandereURL)
}

// Art is a struct that contains art command information.
type art struct {
	name        string
	galleryName string
	tags        []string
	artCount    int
	random      bool
	compressed  bool
	dbClient    *ktsart.DBAPI
	gbClient    *ktsart.GBAPI
	kClient     *ktsart.KAPI
	r34Client   *ktsart.R34API
	sbClient    *ktsart.SBAPI
	yaClient    *ktsart.YAAPI
}

// Run runs the command with the passed command data.
func (a *art) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	a.clearFlags()
	err := a.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, a.name+": "+err.Error())
		return nil
	}

	a.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (a *art) Name() string {
	return a.name
}

// SetName sets the name of the command.
func (a *art) SetName(s string) {
	a.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (a *art) parseCmd(args []string) (err error) {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-gallery":
				a.galleryName, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-random":
				a.random = true
			case "-tags":
				tags, err := extractArgVal(&args, &i)
				if err == nil {
					a.tags = strings.Fields(tags)
					if len(a.tags) > 2 {
						return errors.New("Only 2 tags per search are allowed")
					}
				} else {
					return err
				}
			case "-amount":
				amnt, err := extractArgVal(&args, &i)
				if err == nil {
					a.artCount, err = strconv.Atoi(amnt)
					if err != nil {
						return errors.New("incorrect number format `" + amnt + "`")
					} else if a.artCount > 5 {
						return errors.New("Only 5 results per search are allowed")
					}
				} else {
					return err
				}
			case "-uncompressed":
				a.compressed = false
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
func (a *art) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	msg, err := a.booruImgURLSearch(a.galleryName, a.random, a.compressed) // TODO: Change to embed
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, a.name+": "+err.Error())
	} else if msg == "" {
		s.ChannelMessageSend(m.ChannelID, a.name+": No results found.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, msg)
}

// clearFlags clears data after the previous command execution.
func (a *art) clearFlags() {
	a.galleryName = ""
	a.artCount = 1
	a.tags = a.tags[:0]
	a.compressed = true
	a.random = false
}

// booruImgURLSearch searches for a set of image urls on galleries called "boorus" and returns a formatted ready-to-send message.
func (a *art) booruImgURLSearch(g string, rand bool, compr bool) (msg string, err error) {
	var res *[]ktsart.BooruPost
	var booru ktsart.BooruAPI

	switch g {
	case "danbooru":
		booru = *a.dbClient
	case "gelbooru":
		booru = *a.gbClient
	case "konachan":
		booru = *a.kClient
	case "rule34":
		booru = *a.r34Client
	case "safebooru":
		booru = *a.sbClient
	case "yandere":
		booru = *a.yaClient
	default:
		return "", errors.New(a.name + ": Incorrect gallery name")
	}

	if rand {
		res, err = booru.GetByTagsRandGeneric(a.tags, a.artCount)
	} else {
		res, err = booru.GetByTagsGeneric(a.tags, a.artCount)
	}
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, post := range *res {
		if compr {
			sb.WriteString(*post.ComprIMGURL() + "\n")
		} else {
			sb.WriteString(*post.IMGURL() + "\n")
		}
	}

	return sb.String(), nil
}
