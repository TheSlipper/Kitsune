// Package ktshndlrs provides functions that process commands and events.
package ktshndlrs

import (
	"github.com/bwmarrin/discordgo"
)

// addMuteRoles adds chat and voice mute roles.
func addMuteRoles(s *discordgo.Session, g *discordgo.Guild) {
	//// TODO: Make this work
	//scol := "FF0000"
	//col, err := strconv.ParseInt(scol, 0, 64)
	//if err != nil {
	//	fmt.Printf("[Error while converting color red from hex to dec in gldhandler: %s]\r\n", err)
	//	panic(err)
	//}
	//r, err := s.GuildRoleCreate(g.ID)
	//r.Name = "chat-mute"
	//// r.Color = "#ff0000"
	//r.Color = int(col)
	//s.GuildRoleEdit(g.ID, r.ID, r.Name, r.Color, false, discordgo.PermissionReadMessages^discordgo.PermissionSendMessages, false)
	//
	//r, err = s.GuildRoleCreate(g.ID)
	//r.Name = "voice-mute"
	//r.Color = int(col)
	//s.GuildRoleEdit(g.ID, r.ID, r.Name, r.Color, false, discordgo.PermissionReadMessages^discordgo.PermissionVoiceSpeak, false)
}

// ServerHandler adds information about a server after the bot joins it.
func ServerHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	//// if guildIsRegistered(g) { // TODO: If not registered make a role "chat-mute" and make it so chat-mute can't send text messages on any text channel and do the same for "voice-mute" role with voice channels
	//b, err := ktsdb.ServerExists(g.ID)
	//if b {
	//	return
	//} else if err != nil {
	//	panic(err)
	//}
	//
	//// TODO: Fix this
	//// addMuteRoles(s, g.Guild)
	//
	//b, err = ktsdb.ServerOwnerRegistered(g.OwnerID)
	//if err != nil {
	//	fmt.Printf("[gldhandler - Error while trying to check if server was registered]\r\n")
	//} else if !b {
	//	// if !ownerIsRegistered(g) {
	//	// err = addGuildOwner(g)
	//	b = ktsdb.ServerAddOwner(g.OwnerID)
	//	if !b {
	//		fmt.Printf("[gldhandler - Error while registering the server owner!]\r\n")
	//		return
	//	}
	//}
	//
	//// err = addGuild(g)
	//b = ktsdb.ServerAdd(g)
	//if b == false {
	//	fmt.Printf("[Error while adding guild data to database: %s]\r\n", err)
	//	return
	//}
}
