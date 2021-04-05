// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/TheSlipper/Kitsune/internal/ktsevt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// honk -target-channel "chanel name"

// Honk is an exportable honk struct singleton.
var Honk honk

func init() {
	var err error
	Honk.name = "honk"
	Honk.honkFile, err = os.Open("airhorn.dca")
	if err != nil {
		// TODO: Correct error handling!
		fmt.Println("Error opening dca file :", err)
		panic(err)
	}
	//Honk.buffer = make([][]byte, 0)
	//Honk.loadSound()
}

type honkEvent struct {
	eid int
	repeats bool
	channel string
	server string
	buffer [][]byte
	interrupt bool
}

func (h *honkEvent) loadHonk() error {
	file, err := os.Open("airhorn.dca")
	if err != nil {
		fmt.Println("[Error opening from dca file]")
		fmt.Println(err)
		return err
	}

	var opuslen int16
	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				fmt.Println("[Error reading from dca file]")
				fmt.Println(err)
				return err
			}
			return nil
		} else if err != nil {
			fmt.Println("[Error reading from dca file]")
			fmt.Println(err)
			return err
		}

		// Read encoded pcm from dca file
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)
		if err != nil {
			fmt.Println("[Error reading from dca file]")
			fmt.Println(err)
			return err
		}

		// Append encoded pcm data to the buffer
		h.buffer = append(h.buffer, InBuf)
	}
}

func (h *honkEvent) ExecFun(s *discordgo.Session) {
	h.buffer = make([][]byte, 0)
	err := h.loadHonk()
	if err != nil {
		return
	}

	vc, err := s.ChannelVoiceJoin(h.server, h.channel, false, true)
	if err != nil {
		vc.Close()
		return
	}

	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)
	for _, buff := range h.buffer {
		if h.interrupt {
			break
		}
		vc.OpusSend <- buff
	}

	vc.Speaking(false)
	vc.Disconnect()
}

func (h *honkEvent) Server() string {
	return h.server
}

func (h *honkEvent) Channel() string {
	return h.channel
}

func (h *honkEvent) EID() int {
	return h.eid
}

func (h *honkEvent) Kill() {
	h.interrupt = true
}

func (h *honkEvent) String() string {
	return "Event Type: Honking voice test\r\nRepeats: " + strconv.FormatBool(h.repeats) + "\r\nCommand: " + Honk.name
}

// honk is a struct that contains honk command information.
type honk struct {
	name     string
	target   string
	honkFile *os.File
	buffer   [][]byte
}

// Run runs the command with the passed command data.
func (h *honk) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	h.clearFlags()
	err := h.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, h.name+": "+err.Error())
		return nil
	}

	h.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (h *honk) Name() string {
	return h.name
}

// SetName sets the name of the command.
func (h *honk) SetName(s string) {
	h.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (h *honk) parseCmd(args []string) (err error) {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-target-channel":
				h.target, err = extractArgVal(&args, &i)
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

	if h.target == "" {
		return errors.New("No target channel specified")
	} else {
		return nil
	}
}

// execute runs the command and executes the desired tasks.
func (h *honk) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	g, err := s.Guild(m.GuildID)
	if err != nil {
		return
	}

	for _, ch := range g.Channels {
		if ch.Name == h.target && ch.Type == discordgo.ChannelTypeGuildVoice {
			h.target = ch.ID
		}
	}
	ktsevt.FreshVcEvents <- &honkEvent{ktsevt.GetEID(), false, h.target, m.GuildID, nil, false}
}

// clearFlags clears data after the previous command execution.
func (h *honk) clearFlags() {
	h.target = ""
}

// loadSound attempts to load an encoded sound file from disk.
func (h *honk) loadSound() error {
	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err := binary.Read(h.honkFile, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := h.honkFile.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(h.honkFile, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		h.buffer = append(h.buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func (h *honk) playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range h.buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
