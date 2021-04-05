// Package ktsdb acts as a medium/bindings between the program and the Kitsune database
package ktsdb

import (
	"fmt"
)

// AliasData is a golang representation of the mysql aliases' table entry.
type AliasData struct {
	AliasID          int64  `json:"alias_id"`          // Unique identification number of the alias
	ServerID         int64  `json:"alias_server_id"`   // Identification number of the server that the alias is declared in
	UserID           int64  `json:"alias_user_id"`     // Identification number of the user that has created the alias
	AliasContent     string `json:"alias_content"`     // Content of the alias
	AliasReplacement string `json:"alias_replacement"` // Replacement of the alias
}

// Reset resets the values of the AliasData struct.
func (a *AliasData) Reset() {
	a.AliasID = 0
	a.ServerID = 0
	a.UserID = 0
	a.AliasContent = ""
	a.AliasReplacement = ""
}

// AliasAdd adds an alias to the database.
func AliasAdd(a AliasData) error {
	//_, err := dbConn.Exec("INSERT INTO `aliases` (`alias_id`, `alias_server_id`, `alias_user_id`, `alias_content`, `alias_replacement`) VALUES (NULL, '" + strconv.FormatInt(a.ServerID, 10) + "', '" + strconv.FormatInt(a.UserID, 10) + "', '" + a.AliasContent + "', '" + a.AliasReplacement + "')")
	statement, err := dbConn.Prepare("INSERT INTO aliases (alias_id, server_id, user_id, alias_content, alias_replacement)" +
		" VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	_, err = statement.Exec(a.AliasID, a.ServerID, a.UserID, a.AliasContent, a.AliasReplacement)
	if err != nil {
		fmt.Printf("[dbalias - Error while trying to register an alias in the database: %s]\r\n", err)
		return err
	}
	return nil
}

// AliasDelete deletes the alias with the specified replacement and server id.
func AliasDelete(ar string, sID string) error {
	q := "DELETE FROM aliases WHERE alias_replacement='" + ar + "' AND server_id='" + sID + "'"
	_, err := dbConn.Exec(q)
	if err != nil {
		fmt.Printf("[dbalias - Error while trying to delete an alias: %s]\r\n", err)
		return err
	}
	return nil
}

// AliasExists checks if alias with the given replacement exists in a given guild.
func AliasExists(ar string, sID string) bool {
	rows, err := dbConn.Query("SELECT alias_id FROM aliases WHERE alias_replacement = '" + ar +
		"' AND server_id='" + sID + "'")
	defer rows.Close()
	if err != nil {
		fmt.Printf("[dbalias - Error while trying to check if alias is a duplicate: %s]\r\n", err)
		return false
	}
	return rows.Next()
}

// AliasGetByServerID gets an array of aliases that are assigned to a specific guild.
func AliasGetByServerID(sID string) (arr []AliasData, e error) {
	var aliases []AliasData
	ar, err := dbConn.Query("SELECT * FROM aliases WHERE server_id='" + sID + "'")
	defer ar.Close()
	if err != nil {
		fmt.Printf("[dbalias - Error while trying to get aliases of a guild: %s]\r\n", err)
		return make([]AliasData, 0), err
	}
	for ar.Next() {
		var alias AliasData
		err = ar.Scan(&alias.AliasID, &alias.ServerID, &alias.UserID, &alias.AliasContent, &alias.AliasReplacement)
		if err != nil {
			fmt.Printf("[dbalias - Error while processing an alias: %s]\r\n", err)
			return make([]AliasData, 0), err
		}
		aliases = append(aliases, alias)
	}
	return aliases, nil
}
