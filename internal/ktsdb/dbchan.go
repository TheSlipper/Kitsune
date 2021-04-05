// Package ktsdb acts as a medium/bindings between the program and the Kitsune database
package ktsdb

import "fmt"

// ChanBlacklist blacklists a channel (mutes it - Kitsune won't listen to commands from this channel).
func ChanBlacklist(chID string, sID string) bool {
	statement, err := dbConn.Prepare("INSERT INTO blacklisted_channels (blacklist_id, channel_id, server_id)" +
		" VALUES (NULL, ?, ?)")
	if err != nil {
		fmt.Printf("[dbchan - Error while preparing the statement to blacklist a channel: %s]\r\n", err)
		return false
	}
	_, err = statement.Exec(chID, sID)
	if err != nil {
		fmt.Printf("[dbchan - Error while trying to add the channel to the blacklist: %s]\r\n", err)
		return false
	}

	return true
}

// ChanWhitelist whitelists a channel (unmutes it - Kitsune will now listen to commands from this channel).
func ChanWhitelist(chID string) bool {
	_, err := dbConn.Exec("DELETE FROM blacklisted_channels WHERE channel_id='" + chID + "'")
	if err != nil {
		fmt.Printf("[dbchan - Error while trying to whitelist a channel: %s]\r\n", err)
		return false
	}
	return true
}

// ChanBlacklisted checks if a channel with the specified ID is blacklisted.
func ChanBlacklisted(chID string) bool {
	rows, err := dbConn.Query("SELECT blacklist_id FROM blacklisted_channels WHERE channel_id='" + chID + "'")
	if err != nil {
		fmt.Printf("[dbchan - Error while checking if channel is blacklisted!: %s]\r\n", err)
		// panic(err)
		return false
	}

	return rows.Next()
}

// ChanGetBlacklistedByServer gets the list of blacklisted channels of a server.
func ChanGetBlacklistedByServer(sID string) []string {
	var chans []string
	rows, err := dbConn.Query("SELECT channel_id FROM blacklisted_channels WHERE server_id='" + sID + "'")
	if err != nil {
		fmt.Printf("[dbchan - Error while getting all the blacklisted channels!]\r\n")
		return nil
	}
	for rows.Next() {
		var str string
		rows.Scan(&str)
		chans = append(chans, str)
	}
	return chans
}
