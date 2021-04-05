// Package ktsdb acts as a medium/bindings between the program and the Kitsune database
package ktsdb

import (
	"database/sql"
	"fmt"

	//"github.com/TheSlipper/Kitsune/internal/settings"
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

// dbConn is a package instance of sql database connection.
var dbConn *sql.DB

func init() {
	var err error
	//url := settings.BotSettings.DatabaseUsername + ":" + settings.BotSettings.DatabasePassword + "@tcp(" + settings.BotSettings.DatabaseIPAddress + ":3306)/" + settings.BotSettings.DatabaseName
	//dbConn, err = sql.Open("mysql", url)
	//if err != nil {
	//	fmt.Printf("[Error while creating a connection to the database - %s]\r\n", err)
	//	panic(err)
	//}

	// Create connection
	dbConn, err = sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Printf("[Error while creating a connection to the database - %s]\r\n", err)
		panic(err)
	}

	// Create tables if they don't exist
	createServerTable()
	createAliasTable()
	createBlacklistedChannelsTable()
	createServerPrefixesTable()
}

func createServerTable() {
	statement, err := dbConn.Prepare("CREATE TABLE IF NOT EXISTS servers (server_id TEXT PRIMARY KEY, " +
		"server_join_date TEXT, server_prefix TEXT)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func createAliasTable() {
	statement, err := dbConn.Prepare("CREATE TABLE IF NOT EXISTS aliases (alias_id INTEGER PRIMARY KEY, " +
		"server_id TEXT, user_id TEXT, alias_content TEXT, alias_replacement TEXT)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func createBlacklistedChannelsTable() {
	statement, err := dbConn.Prepare("CREATE TABLE IF NOT EXISTS blacklisted_channels (blacklist_id " +
		"INTEGER PRIMARY KEY, channel_id TEXT, server_id TEXT)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}

func createServerPrefixesTable() {
	statement, err := dbConn.Prepare("CREATE TABLE IF NOT EXISTS blacklisted_channels (prefix_id " +
		"INTEGER PRIMARY KEY, server_id TEXT)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec()
	if err != nil {
		panic(err)
	}
}