// Package ktsdb acts as a medium/bindings between the program and the Kitsune database
package ktsdb

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ServerData golang representation of the Servers' table entry.
type ServerData struct {
	ServerID       int64     `json:"server_id"`
	ServerJoinDate time.Time `json:"server_join_date"`
	ServerName     string    `json:"server_name"`
	//ServerOwnerID  int64     `json:"server_owner_id"`
	ServerPrefix   string    `json:"server_prefix"`
}

//// ServerAddOwner adds information on the owner of the server to the database.
//func ServerAddOwner(id string) bool {
//	q := "INSERT INTO users (`user_id`, `user_registration_date`, `user_privilege_group`, `user_kitsune_servers`) VALUES (" + id + ", current_timestamp(), 'discord_user', '0')"
//	_, err := dbConn.Exec(q)
//	if err != nil {
//		fmt.Printf("[dbserver - could not add the server owner to the database: %s]\r\n", err)
//		return false
//	}
//	return true
//}

// ServerAdd adds information on the server to the database.
func ServerAdd(g *discordgo.GuildCreate) bool {
	statement, err := dbConn.Prepare("INSERT INTO servers (server_id, server_join_date, server_prefix) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Printf("[dbserver - Error while preparing the statement: %s]\r\n", err)
		return false
	}

	_, err = statement.Exec(g.Guild.ID, time.Now().String(), "$")
	if err != nil {
		fmt.Printf("[dbserver - Error while adding guild data to database: %s]\r\n", err)
		return false
	}

	//var rows *sql.Rows
	//var servCount int64
	//rows, err = dbConn.Query("SELECT user_kitsune_servers FROM `users` WHERE `users`.`user_id`=" + g.OwnerID)
	//if err != nil {
	//	fmt.Printf("[dbserver - Error while fetching admin user data from the database: %s]\r\n", err)
	//	return false
	//}
	//defer rows.Close()
	//rows.Next()
	//rows.Scan(&servCount)
	//servCount++
	//
	//q = "UPDATE `users` SET `user_kitsune_servers` = '" + string(servCount) + "' WHERE `users`.`user_id` = 123"
	//_, err = dbConn.Exec(q)
	//if err != nil {
	//	fmt.Printf("[dbserver - Error while incrementing user_kitsune_servers in kitsune database: %s]\r\n", err)
	//	return false
	//}

	return true
}

// ServerExists checks if server with the specified ID exists.
func ServerExists(sID string) (b bool, err error) {
	rows, err := dbConn.Query("SELECT server_id FROM servers WHERE server_id = '" + sID + "'")
	if err != nil {
		fmt.Printf("[Error while checking if the guild is registered in kitsune database: %s]\r\n", err)
		// panic(err.Error())
		return false, err
	}
	return rows.Next(), nil
}

// ServerOwnerRegistered checks if the server owner is already registered.
//func ServerOwnerRegistered(oID string) (b bool, err error) {
//	rows, err := dbConn.Query("SELECT user_id FROM `users` WHERE `users`.`user_id` = " + oID)
//	if err != nil {
//		fmt.Printf("[Error while checking if the owner of the guild is registered in kitsune database: %s]\r\n", err)
//		// panic(err.Error())
//		return b, err
//	}
//	return rows.Next(), nil
//}
