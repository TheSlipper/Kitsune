// Package ktscmd contains structs and functions used for execution of kitsune commands.
package ktscmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	//_ "github.com/go-sql-driver/mysql"
)

// sendArtEmbed sends a discord embed with art data to the specified channel.
// func sendArtEmbed(s *discordgo.Session, chID string, searchRes artsearch.SearchResult) {
// 	emb := &discordgo.MessageEmbed{
// 		Author: &discordgo.MessageEmbedAuthor{},
// 		Color:  0x0097e6,
// 		Title:  searchRes.Query(),
// 		//Description: searchRes.Artist(),
// 		Fields: []*discordgo.MessageEmbedField{
// 			&discordgo.MessageEmbedField{
// 				Name:   "Artist: ",
// 				Value:  searchRes.Artist(),
// 				Inline: true,
// 			},
// 			&discordgo.MessageEmbedField{
// 				Name:   "Uploaded by: ",
// 				Value:  searchRes.Uploader(),
// 				Inline: true,
// 			},
// 			&discordgo.MessageEmbedField{
// 				Name:   "Tags: ",
// 				Value:  searchRes.Tags(),
// 				Inline: true,
// 			},
// 		},
// 		Image: &discordgo.MessageEmbedImage{
// 			URL: searchRes.URL(),
// 		},
// 	}
// 	_, _ = s.ChannelMessageSendEmbed(chID, emb)
// }

// sendRawEmbed sends a simple discord embed to the specified channel.
func sendRawEmbed(s *discordgo.Session, chID string, hdr string, bdy string) {
	emb := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x0097e6,
		Title:       hdr,
		Description: bdy,
	}
	_, err := s.ChannelMessageSendEmbed(chID, emb)
	if err != nil {
		fmt.Printf("[sendRawEmbed(discordgo.Session, string, string, string) - error while sending the embed (hdr length: %d; bdy length: %d): %s\r\n", len(hdr), len(bdy), err)
		panic(err.Error())
	}
}

// sendRawEmbedMultiple sends multiple "raw" discord embeds to the specified channel.
func sendRawEmbedMultiple(s *discordgo.Session, chID string, n string, hdr []string, bdy []string) {
	flds := make([]*discordgo.MessageEmbedField, len(hdr))
	allEmpty := true
	fldsLen := 0
	for i, j := 0, 0; j < len(flds); i, j = i+1, j+1 {
		if hdr[j] == "" || bdy[j] == "" {
			i--
			continue
		}
		allEmpty = false
		flds[i] = &discordgo.MessageEmbedField{
			Name:   hdr[j],
			Value:  bdy[j],
			Inline: true,
		}
		fldsLen = i + 1
	}

	var emb *discordgo.MessageEmbed
	if allEmpty {
		emb = &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0x0097e6,
			Title:       n,
			Description: "No items",
		}
	} else {
		emb = &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color:  0x0097e6,
			Title:  n,
			Fields: flds[0:fldsLen],
		}
	}
	_, err := s.ChannelMessageSendEmbed(chID, emb)
	if err != nil {
		fmt.Printf("[sendRawEmbedMultiple(discordgo.Session, string, string, string) - error while sending the embed (hdr length: %d; bdy length: %d): %s\r\n", len(hdr), len(bdy), err)
		panic(err.Error())
	}
}

// extractInsideCharsAndIterate extracts all of the string data inside the two of the specified string patterns.
func extractInsideCharsAndIterate(char string, i *int, arr []string) (string, error) {
	var sb strings.Builder
	sb.WriteString(arr[*i])
	if strings.HasSuffix(arr[*i], char) {

		return strings.Replace(sb.String(), char, "", -1), nil
	}
	for {
		*i++
		if *i == len(arr) {
			return "", errors.New("no closing character `" + char + "` found")
		}

		sb.WriteString(" " + strings.Replace(arr[*i], "/"+char, "{placeHolder}", -1))
		if strings.HasSuffix(sb.String(), char) {
			return strings.TrimSuffix(strings.TrimPrefix(strings.Replace(sb.String(), "{placeHolder}", char, -1), char), char), nil
		}
	}
}

// TODO: Example of this ^.
// https://www.reddit.com/r/golang/comments/2yqe1d/getting_code_examples_to_show_up_in_godoc/

// memberHasPermission checks if a member has a  specified permission.
func memberHasPermission(s *discordgo.Session, serverID string, userID string, permission int) (bool, error) {
	member, err := s.State.Member(serverID, userID)
	if err != nil {
		if member, err = s.GuildMember(serverID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(serverID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions& int64(permission) != 0 {
			return true, nil
		}
	}

	return false, nil
}

// getRoleID gets the specified role's ID by its name.
func getRoleID(s *discordgo.Session, gID string, r string) (string, error) {
	guild, _ := s.Guild(gID)
	roles := guild.Roles
	for _, rl := range roles {
		if rl.Name == r {
			return rl.ID, nil
		}
	}

	return "", InvalidRoleIDError{RoleName: r}
}

// extractArgVal extracts a value from a command argument.
func extractArgVal(args *[]string, i *int) (val string, err error) {
	vn := (*args)[*i]
	if strings.Contains(vn, "=") {
		val = vn[strings.Index(vn, "="):]
		return val, nil
	} else {
		if len(*args) == *i+1 {
			return "", errors.New("Sudden end of command while parsing: No argument value found")
		} else if !strings.Contains((*args)[*i+1], "\"") {
			// return "", errors.New("Argument value should be inside quotation marks")
			*i = *i + 1
			val := (*args)[*i]
			return val, nil
		}
		*i = *i + 1
		arg, err := extractInsideCharsAndIterate("\"", i, *args)
		return arg, err
	}
}

// InvalidRoleIDError error that indicates an error in role's ID evaluation.
type InvalidRoleIDError struct {
	RoleName string
}

// Error text that is to be displayed when this error occurs.
func (i InvalidRoleIDError) Error() string {
	return "[InvalidRoleIDError: no role with name '" + i.RoleName + "]\r\n"
}
