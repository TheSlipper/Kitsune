// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// mngrole -create-role "role name"
// mngrole -delete-role  "role name"
// mngrole -set-color "#1A5723" -target-role "role name"
// mngrole -set-color "rgb(69, 96, 60)" -target-role "role name"
// mngrole -add-permission "permission" -target-role "role name"
// mngrole -remove-permission "permission" -target-role "role name"

// MngRole is an exportable mngRole struct singleton.
var MngRole mngRole

func init() {
	MngRole.name = "mngrole"
}

// MngRole is a struct that contains mngrole command information.
type mngRole struct {
	name        string
	targetRoles []string
	color       string
	rolePerm    string
	createRole  bool
	dlRole      bool
	setColor    bool
	addPerm     bool
	rmPerm      bool
}

// Run runs the command with the passed command data.
func (mng *mngRole) Run(args []string, m *discordgo.MessageCreate, s *discordgo.Session) error {
	mng.clearFlags()
	err := mng.parseCmd(args)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, mng.name+": "+err.Error())
		return nil
	}

	mng.execute(m, s)
	return nil
}

// Name returns the name of the command.
func (mng *mngRole) Name() string {
	return mng.name
}

// SetName sets the name of the command.
func (mng *mngRole) SetName(s string) {
	mng.name = s
}

// parseCmd parses the passed command data and sets up the command flags and data.
func (mng *mngRole) parseCmd(args []string) (err error) {
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			switch args[i] {
			case "-create-role":
				mng.createRole = true
				foo, err := extractArgVal(&args, &i)
				if err != nil {
					return err
				}
				mng.targetRoles = append(mng.targetRoles, foo)
			case "-delete-role":
				mng.dlRole = true
				foo, err := extractArgVal(&args, &i)
				if err != nil {
					return err
				}
				mng.targetRoles = append(mng.targetRoles, foo)
			case "-set-color":
				mng.setColor = true
				mng.color, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-add-permission":
				mng.addPerm = true
				mng.rolePerm, err = extractArgVal(&args, &i)
				mng.color, err = extractArgVal(&args, &i)
				if err != nil {
					return err
				}
			case "-remove-permission":
				mng.rmPerm = true
				mng.rolePerm, err = extractArgVal(&args, &i)
				mng.color, err = extractArgVal(&args, &i)
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
	return nil
}

// execute runs the command and executes the desired tasks.
func (mng *mngRole) execute(m *discordgo.MessageCreate, s *discordgo.Session) {
	b, err := memberHasPermission(s, m.GuildID, m.Author.ID, discordgo.PermissionManageRoles)
	if err != nil || !b {
		_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": you do not have the permission to edit roles!")
		return
	}

	if mng.createRole && !(mng.dlRole || mng.setColor || mng.addPerm || mng.rmPerm) {
		// TODO:
		// BUG(slipper): This throws an HTTP bad request error: err code 400
		mng.addRole(m, s)
	} else if mng.dlRole && !(mng.createRole || mng.setColor || mng.addPerm || mng.rmPerm) {
		mng.deleteRole(m, s)
	} else if mng.setColor && !(mng.createRole || mng.dlRole || mng.addPerm || mng.rmPerm) {
		mng.changeRoleColor(m, s)
	} else if mng.addPerm && !(mng.createRole || mng.dlRole || mng.setColor || mng.rmPerm) {
		// TODO:
	} else if mng.rmPerm && !(mng.createRole || mng.dlRole || mng.setColor || mng.addPerm) {
		// TODO:
	} else {
		s.ChannelMessageSend(m.ChannelID, mng.name+": Actions specified incorrectly - you must select **only one** action!")
	}

}

// clearFlags clears data after the previous command execution.
func (mng *mngRole) clearFlags() {
	mng.targetRoles = mng.targetRoles[:0]
	mng.color = ""
	mng.rolePerm = ""
	mng.createRole = false
	mng.dlRole = false
	mng.setColor = false
	mng.addPerm = false
	mng.rmPerm = false
}

// addRole creates role with the specified data.
func (mng *mngRole) addRole(m *discordgo.MessageCreate, s *discordgo.Session) {
	for _, rn := range mng.targetRoles {
		r, err := s.GuildRoleCreate(m.GuildID)
		if err != nil {
			fmt.Printf("["+mng.name+" - unexpected error occurred while creating a role: %s]\r\n", err)
			return
		}
		_, _ = s.GuildRoleEdit(m.GuildID, r.ID, rn, 0, false, discordgo.PermissionReadMessages, false)
	}
}

// deleteRole deletes role with the specified name.
func (mng *mngRole) deleteRole(m *discordgo.MessageCreate, s *discordgo.Session) {
	g, err := s.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("["+mng.name+" - unexpected error occurred while deleting a role: %s]\r\n", err)
		return
	}
	for _, r := range g.Roles {
		for _, rn := range mng.targetRoles {
			if rn == r.Name {
				_ = s.GuildRoleDelete(m.GuildID, r.ID)
			}
		}
	}
}

// setColor sets the color of a role to the specified RGB or HEX value.
func (mng *mngRole) changeRoleColor(m *discordgo.MessageCreate, s *discordgo.Session) {
	var val int
	var err error
	if mng.color[0] == '#' {
		decimal, err := strconv.ParseInt(mng.color[1:], 0, 64)
		if err != nil {
			fmt.Printf("["+mng.name+" - unexpected error occurred while converting hex color value to decimal: %s\r\n]", err)

			_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": incorrect hex color value \""+mng.color+"\"")
			return
		}
		val = int(decimal)
	} else if strings.HasPrefix(mng.color, "rgb(") && strings.HasSuffix(mng.color, ")") {
		clrstr := mng.color
		clrstr = strings.TrimPrefix(clrstr, "rgb(")
		clrstr = strings.TrimSuffix(clrstr, ")")
		clrstr = strings.Replace(clrstr, " ", "", -1)
		rgb := strings.Split(clrstr, ",")

		var sb strings.Builder
		for _, v := range rgb {
			dec, err := strconv.ParseInt(v, 0, 64)
			if err != nil {
				fmt.Printf("["+mng.name+" - unexpected error occurred while converting to int from string in rgb case: %s\r\n]", err)
			}
			s := strconv.FormatInt(dec, 16)
			sb.WriteString(s)
		}
		fmt.Println()
		decimal, err := strconv.ParseInt(sb.String(), 0, 64)
		if err != nil {
			fmt.Printf("["+mng.name+" - unexpected error occurred while converting hex color value to decimal in rgb case: %s\r\n]", err)
			return
		}
		val = int(decimal)
	} else {
		s.ChannelMessageSend(m.ChannelID, mng.name+": incorrect color value \""+mng.color+"\". You can read `$help "+mng.name+"` to see what colors are valid.")
	}

	g, err := s.Guild(m.GuildID)
	if err != nil {
		fmt.Printf("["+mng.name+" - unexpected error occurred while setting a color of a role: %s]\r\n", err)
		return
	}

	for _, r := range g.Roles {
		for _, rn := range mng.targetRoles {
			if rn == r.Name {
				if err != nil {
					fmt.Printf("["+mng.name+" - unexpected error occurred while converting the action data to decimal color representation: %s\r\n]", err)
					_, _ = s.ChannelMessageSend(m.ChannelID, mng.name+": incorrect color value \""+mng.color+"\"")
					return
				}
				_, _ = s.GuildRoleEdit(m.GuildID, r.ID, r.Name, val, r.Hoist, r.Permissions, r.Mentionable)
			}
		}
	}
}

// addPermission TODO
func (mng *mngRole) addPermission() {
	// TODO: Adding permissions to a given role
}
