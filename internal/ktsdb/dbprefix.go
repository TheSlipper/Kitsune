// Package ktsdb acts as a medium/bindings between the program and the Kitsune database
package ktsdb

import "fmt"

// PrefixChange changes the prefix used for Kitsune commands in a given server.
func PrefixChange(p string, gID string) bool {
	q := "UPDATE `servers` SET `server_prefix` = '" + p + "' WHERE `servers`.`server_id` = '" + gID + "'"
	_, err := dbConn.Exec(q)
	if err != nil {
		fmt.Printf("[dbprefix - Error while trying to delete an alias: %s]\r\n", err)
		return false
	}
	return true
}

// PrefixGetByServerID gets the command prefix of a specific guild.
func PrefixGetByServerID(gID string) string {
	var prefix string
	rows, err := dbConn.Query("SELECT server_prefix FROM servers WHERE server_id='" + gID + "'")
	defer func() {
		rows.Close()
		if err != nil {
			fmt.Println("test")
		}
	}()
	if err != nil {
		fmt.Printf("[dbprefix - Unregistered guild of ID '%s': %s]\r\n", gID, err)
		return "$"
	} else if !rows.Next() {
		fmt.Printf("[dbprefix - No prefix for guild of ID '%s': %s]\r\n", gID, err)
		return "$"
	}
	rows.Scan(&prefix)
	return prefix
}
