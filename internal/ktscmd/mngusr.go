// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TheSlipper/Kitsune/internal/ktsevt"

	"github.com/bwmarrin/discordgo"
)

// mngusr -kick -target-user @username#1234
// mngusr -ban -target-user @username#1234
// mngusr -assign-role "role name" -target-user @username#1234
// mngusr -remove-role "role name" -target-user @username#1234
// mngusr -voice-mute -target-user @username#1234 -duration "3m50s"
// mngusr -chat-mute -target-user @username#1234 -duration "3m50s"

// MngUsr is an exportable mngUsr struct singleton.
var MngUsr mngUsr

func init() {
	MngUsr.name = "mngusr"
}

// MngUsr is a struct that contains mngusr command information.
type mngUsr struct {
	name       string
	kick       bool
	ban        bool
	vcMute     bool
	txtMute    bool
	addRole    bool
	targetUsr  string
	targetRole string
	duration   string
}

// Run runs the command with the passed command data.
func (mng *mngUsr) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	mng.clearFlags()
	err := mng.parseCmd(args)
	fmt.Println(m.Content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, mng.name+": "+err.Error())
		return nil
	} else if mng.targetUsr == s.State.User.ID {
		s.ChannelMessageSend(m.ChannelID, mng.name+": cannot edit Kitsune with this command!")
	}

	mng.execute(m, s)
	return nil
}

// Name returnp the name of the command.
func (mng *mngUsr) Name() string {
	return mng.name
}

// SetName sets the name of the command.
func (mng *mngUsr) SetName(s string) {
	mng.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (mng *mngUsr) parseCmd(args []string) error {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-kick":
				mng.kick = true
			case "-ban":
				mng.ban = true
			case "-voice-mute":
				mng.vcMute = true
			case "-chat-mute":
				mng.txtMute = true
			case "-assign-role":
				var err error
				mng.addRole = true
				mng.targetRole, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-remove-role":
				var err error
				mng.addRole = false
				mng.targetRole, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-target-user":
				foo, err := extractArgVal(&args, &i)
				if err != nil {
					return err
				} else if !(strings.HasPrefix(foo, "<@") && strings.HasSuffix(foo, ">")) {
					return errors.New("Incorrect user data entered for argument -target-user")
				} else {
					mng.targetUsr = strings.Replace(strings.Replace(foo, "<@", "", 1), ">", "", 1)
				}
			case "-duration":
				var err error
				mng.duration, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			default:
				return errors.New("Unkown argument '" + args[i] + "'")
			}
		} else if strings.HasPrefix(args[i], "<@") && strings.HasSuffix(args[i], ">") {
			args[i] = strings.Replace(args[i], "<@", "", -1)
			args[i] = strings.Replace(args[i], ">", "", -1)
		} else {
			return errors.New("No such command value/flag as: " + args[i])
		}
	}

	return nil
}

// execute runs the command and executes the desired tasks.
func (mng *mngUsr) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	if mng.kick && !(mng.ban || mng.addRole || mng.vcMute || mng.txtMute) {
		mng.kickUsr(m, s)
	} else if mng.ban && !(mng.kick || mng.addRole || mng.vcMute || mng.txtMute) {
		mng.banUsr(m, s)
	} else if mng.addRole && !(mng.ban || mng.kick || mng.vcMute || mng.txtMute) {
		mng.manageRole(m, s, mng.addRole)
	} else if !mng.addRole && !(mng.kick || mng.ban || mng.vcMute || mng.txtMute) {
		mng.manageRole(m, s, mng.addRole)
	} else if mng.vcMute && !(mng.kick || mng.addRole || mng.ban || mng.txtMute) {

	} else if mng.txtMute && !(mng.kick || mng.addRole || mng.vcMute || mng.ban) {

	}
}

// clearFlags clears data after the previous command execution.
func (mng *mngUsr) clearFlags() {
	mng.kick = false
	mng.ban = false
	mng.vcMute = false
	mng.txtMute = false
	mng.addRole = false
	mng.targetUsr = ""
	mng.duration = ""
}

// kick kicks the tagged user from the server in which the command was ran.
func (mng *mngUsr) kickUsr(m *discordgo.MessageCreate, s *discordgo.Session) {
	hasPermissions, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionKickMembers)
	if !hasPermissions || err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": You do not have the permission to kick users!")
		return
	}
	// for i := 0; i < len(mng.targetID); i++ {
	err = s.GuildMemberDelete(m.GuildID, mng.targetUsr)
	if err != nil {
		panic(err)
	}
	// }
}

// ban bans the tagged user from the server in which the command was ran.
func (mng *mngUsr) banUsr(m *discordgo.MessageCreate, s *discordgo.Session) {
	hasPermissions, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionBanMembers)
	if !hasPermissions || err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": You do not have the permission to ban users!")
		return
	}
	err = s.GuildBanCreate(m.GuildID, mng.targetUsr, 0)
	if err != nil {
		panic(err)
	}
}

// manageRole adds or removes the role from a specified user.
func (mng *mngUsr) manageRole(m *discordgo.MessageCreate, s *discordgo.Session, adding bool) {
	hasPermissions, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionManageRoles)
	if (!hasPermissions && !mng.txtMute) || err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": You do not have the permission to manage roles!")
		return
	}

	guild, _ := s.Guild(m.GuildID)
	roles := guild.Roles
	roleID := ""

	// TODO: check if this loop is even necessary
	for _, rl := range roles {
		if rl.Name == mng.targetRole {
			roleID = rl.ID
		}
	}

	if roleID == "" {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": no role with ID: '"+mng.targetRole+"'")
		return
	}

	// for _, usr := range mng.targetID {
	if adding {
		err = s.GuildMemberRoleAdd(m.GuildID, mng.targetUsr, roleID)
	} else {
		err = s.GuildMemberRoleRemove(m.GuildID, mng.targetUsr, roleID)
	}

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": error while assigning/removing user's role *(Does Kitsune have the necessary permissions? If so removing and adding Kitsune back will help)*")
	}
	// }
	fmt.Printf("Finished managing roles\r\n")
}

// mute mutes user in text chats or voice chats of a given server.
func (mng *mngUsr) mute(m *discordgo.MessageCreate, s *discordgo.Session, chatMute bool) {
	hasPermissions, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionVoiceMuteMembers)
	if !hasPermissions || err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": You do not have the permission to mute users!")
		return
	}

	evtTime := time.Now()
	dur, err := time.ParseDuration(mng.duration)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": mute duration \""+mng.duration+"\" formatting is not correct")
		//panic(err.Error())
		return
	}
	evtTime = evtTime.Add(dur)

	var rnp string
	if chatMute {
		//rN.WriteString("chat")
		rnp = "chat"
	} else {
		rnp = "voice"
	}
	rID, err := getRoleID(s, m.GuildID, rnp+"-mute")
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": role \""+rnp+"-mute\" is missing! It shouldn't have been deleted!")
		return
	}
	ktsevt.FreshTxtEvts <- &ktsevt.ChatUnmuteEvt{ktsevt.GetEID(), evtTime, mng.targetUsr, m.GuildID, rID}

	mng.targetRole = rnp + "-mute"
	mng.addRole = true
	mng.manageRole(m, s, true)
}
