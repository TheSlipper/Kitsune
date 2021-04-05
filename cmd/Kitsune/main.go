package main

// TODO: More graceful way of displaying errors (use println and not printf)

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheSlipper/Kitsune/internal/ktsevt"
	"github.com/TheSlipper/Kitsune/internal/ktshndlrs"
	"github.com/TheSlipper/Kitsune/internal/settings"

	"github.com/bwmarrin/discordgo"
)

// TODO: mutexes on run commands
func main() {
	fmt.Printf("[Attempting to create connection with discord using the '%s' token]\r\n", settings.BotSettings.BotToken)
	dg, err := discordgo.New("Bot " + settings.BotSettings.BotToken)
	if err != nil {
		fmt.Printf("[Error creating a discord session: %s]\r\n", err)
		panic(err)
	}
	ktsevt.Session = dg

	fmt.Printf("[Attaching discord API handlers]\r\n")
	dg.AddHandler(ktshndlrs.MsgHandler)
	dg.AddHandler(ktshndlrs.ServerHandler)

	fmt.Printf("[Creating goroutine for the internal bot event handling]\r\n")
	go ktshndlrs.BotEventHandler()

	err = dg.Open()
	if err != nil {
		fmt.Printf("[Error opening connection: %s]\r\n", err)
		panic(err)
	}
	defer dg.Close()

	fmt.Printf("[Successfully initiated kitsune]\r\n")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
