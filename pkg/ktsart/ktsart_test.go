// Package ktsart provides bindings to the booru APIs.
package ktsart

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/TheSlipper/Kitsune/internal/settings"
)

// decryptionKey contains the key to the encrypted json file with all of the tokens.
const decryptionKey = ""

// getGbConn creates and returns the connection to gelbooru.
func getGbConn(t *testing.T) *GBAPI {
	*settings.StartupOptions.AESKey = decryptionKey
	decryptSettings(t)

	var client http.Client
	gbAuth := GBAuth{User: settings.BotSettings.GelbooruUsrID, Hash: settings.BotSettings.GelbooruToken}
	gb := NewGB(&client, settings.BotSettings.Addresses.GelbooruURL, &gbAuth)
	return gb
}

// decryptSettings loads an encrypted file, decrypts it and puts it in the BotSettings.
func decryptSettings(t *testing.T) {
	data, err := ioutil.ReadFile(*settings.StartupOptions.SettingsFilePath + ".encrypted")
	if err != nil {
		fmt.Printf("[Error while reading the encrypted settings file: %s]\r\n", err)
		panic(err)
	}
	key := settings.CreateHash(*settings.StartupOptions.AESKey)
	dataDecrypted := settings.Decrypt(data, key)
	err = json.Unmarshal(dataDecrypted, &settings.BotSettings)
	if err != nil {
		t.Errorf("error while decrypting: %s", err.Error())
		panic(err)
	}
}

// TestNewGb tests the creation of new gelbooru API connection
func TestNewGb(t *testing.T) {
	gb := getGbConn(t)

	p, err := gb.GetByTagsRaw([]string{"asuka_langley_soryu", "plugsuit"}, 5)
	if err != nil {
		t.Error("Error while fetching the results")
	}

	t.Log("URL: " + p.List[0].FileURL)
}

// TestGBGetByTagsRaw tests if tag search can handle all of the possible cases of search works correctly.
func TestGBGetByTagsRaw(t *testing.T) {
	gb := getGbConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"asuka_langley_soryu", "plugsuit", 1},
		{"asuka_langley_soryu", "plugsuit", 5},
		{"asuka_langley_soryu", "plugsuit", 10},
		{"asuka_langley_soryu", "plugsuit", 20},
		{"asuka_langley_soryu", "plugsuit", 100},
		{"asukalangley_soryu", "plugsuit", 1},
		{"asuka_langley_soryu", "plugsui", 1},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := gb.GetByTagsRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// TestGBGetByTagsRandRaw tests if tag search with randomly sorted results works correctly.
func TestGBGetByTagsRandRaw(t *testing.T) {
	gb := getGbConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"asuka_langley_soryu", "plugsuit", 1},
		{"asuka_langley_soryu", "plugsuit", 5},
		{"asuka_langley_soryu", "plugsuit", 10},
		{"asuka_langley_soryu", "plugsuit", 20},
		{"asuka_langley_soryu", "plugsuit", 100},
		{"asukalangley_soryu", "plugsuit", 1},
		{"asuka_langley_soryu", "plugsui", 1},
		{"asuka_langley_soryu", "plugsuit", 101},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := gb.GetByTagsRandRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			if err.Error() == "amount of requested posts too big" {
				continue
			} else {
				t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
			}
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// getDBConn creates and returns the connection to gelbooru.
func getDBConn(t *testing.T) *DBAPI {
	*settings.StartupOptions.AESKey = decryptionKey
	decryptSettings(t)

	var client http.Client
	dbAuth := DBAuth{User: settings.BotSettings.DanbooruLogin, Hash: settings.BotSettings.DanbooruToken}
	db := NewDB(&client, settings.BotSettings.Addresses.DanbooruURL, &dbAuth)
	return db
}

// TestNewGb tests the creation of new gelbooru API connection
func TestNewDB(t *testing.T) {
	db := getDBConn(t)

	p, err := db.GetByTagsRaw([]string{"souryuu_asuka_langley", "plugsuit"}, 5)
	if err != nil {
		t.Error("Error while fetching the results")
	}

	t.Log("URL: " + p.List[0].FileURL)
}

// TestDBGetByTagsRaw tests if tag search can handle all of the possible cases of search works correctly.
func TestDBGetByTagsRaw(t *testing.T) {
	db := getDBConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"souryuu_asuka_langley", "plugsuit", 1},
		{"souryuu_asuka_langley", "plugsuit", 5},
		{"souryuu_asuka_langley", "plugsuit", 10},
		{"souryuu_asuka_langley", "plugsuit", 20},
		{"souryuu_asuka_langley", "plugsuit", 100},
		{"souryuuasuka_langley", "plugsuit", 1},
		{"souryuu_asuka_langley", "plugsui", 1},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := db.GetByTagsRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// TestGBGetByTagsRandRaw tests if tag search with randomly sorted results works correctly.
func TestDBGetByTagsRandRaw(t *testing.T) {
	db := getDBConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"souryuu_asuka_langley", "plugsuit", 1},
		{"souryuu_asuka_langley", "plugsuit", 5},
		{"souryuu_asuka_langley", "plugsuit", 10},
		{"souryuu_asuka_langley", "plugsuit", 20},
		{"souryuu_asuka_langley", "plugsuit", 100},
		{"souryuuasuka_langley", "plugsuit", 1},
		{"souryuu_asuka_langley", "plugsui", 1},
		{"souryuu_asuka_langley", "plugsuit", 101},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := db.GetByTagsRandRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			if err.Error() == "amount of requested posts too big" {
				continue
			} else {
				t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
			}
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// getDBConn creates and returns the connection to gelbooru.
func getKBConn(t *testing.T) *KAPI {
	*settings.StartupOptions.AESKey = decryptionKey
	decryptSettings(t)

	var client http.Client
	db := NewKB(&client, settings.BotSettings.Addresses.KonachanURL)
	return db
}

// TestNewGb tests the creation of new gelbooru API connection
func TestNewKB(t *testing.T) {
	db := getKBConn(t)

	p, err := db.GetByTagsRaw([]string{"soryu_asuka_langley", "bodysuit"}, 5)
	if err != nil {
		t.Error("Error while fetching the results")
	}

	t.Log("URL: " + p.List[0].FileURL)
}

// TestDBGetByTagsRaw tests if tag search can handle all of the possible cases of search works correctly.
func TestKBGetByTagsRaw(t *testing.T) {
	db := getDBConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"soryu_asuka_langley", "bodysuit", 1},
		{"soryu_asuka_langley", "bodysuit", 5},
		{"soryu_asuka_langley", "bodysuit", 10},
		{"soryu_asuka_langley", "bodysuit", 20},
		{"soryu_asuka_langley", "bodysuit", 100},
		{"soryuuasuka_langley", "bodysuit", 1},
		{"soryu_asuka_langley", "bodysui", 1},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := db.GetByTagsRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// TestGBGetByTagsRandRaw tests if tag search with randomly sorted results works correctly.
func TestKBGetByTagsRandRaw(t *testing.T) {
	db := getDBConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"soryu_asuka_langley", "bodysuit", 1},
		{"soryu_asuka_langley", "bodysuit", 5},
		{"soryu_asuka_langley", "bodysuit", 10},
		{"soryu_asuka_langley", "bodysuit", 20},
		{"soryu_asuka_langley", "bodysuit", 100},
		{"soryuuasuka_langley", "bodysuit", 1},
		{"soryu_asuka_langley", "bodysui", 1},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := db.GetByTagsRandRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			if err.Error() == "amount of requested posts too big" {
				continue
			} else {
				t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
			}
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// TODO: Tests for safebooru
// getDBConn creates and returns the connection to gelbooru.
func getSBConn(t *testing.T) *SBAPI {
	*settings.StartupOptions.AESKey = decryptionKey
	decryptSettings(t)

	var client http.Client
	sb := NewSB(&client, settings.BotSettings.Addresses.SafebooruURL)
	return sb
}

// TestNewGb tests the creation of new gelbooru API connection
func TestNewSB(t *testing.T) {
	db := getSBConn(t)

	p, err := db.GetByTagsRaw([]string{"souryuu_asuka_langley", "plugsuit"}, 5)
	if err != nil {
		t.Error("Error while fetching the results")
	}

	t.Log("URL: " + p.List[0].FileURL)
}

// TestDBGetByTagsRaw tests if tag search can handle all of the possible cases of search works correctly.
func TestSBGetByTagsRaw(t *testing.T) {
	sb := getSBConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"souryuu_asuka_langley", "plugsuit", 1},
		{"souryuu_asuka_langley", "plugsuit", 5},
		{"souryuu_asuka_langley", "plugsuit", 10},
		{"souryuu_asuka_langley", "plugsuit", 20},
		{"souryuu_asuka_langley", "plugsuit", 100},
		{"souryuuasuka_langley", "plugsuit", 1},
		{"souryuu_asuka_langley", "plugsui", 1},
	}

	for _, entry := range tagTables {
		var sbuilder strings.Builder
		p, err := sb.GetByTagsRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
		}

		sbuilder.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sbuilder.WriteString(", ")
			}
			sbuilder.WriteString(pst.FileURL)
		}
		sbuilder.WriteString("\n")
		sbuilder.WriteString("\n")
		t.Logf(sbuilder.String())
	}
}

// TestGBGetByTagsRandRaw tests if tag search with randomly sorted results works correctly.
func TestSBGetByTagsRandRaw(t *testing.T) {
	sb := getSBConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"souryuu_asuka_langley", "plugsuit", 1},
		{"souryuu_asuka_langley", "plugsuit", 5},
		{"souryuu_asuka_langley", "plugsuit", 10},
		{"souryuu_asuka_langley", "plugsuit", 20},
		{"souryuu_asuka_langley", "plugsuit", 100},
		{"souryuuasuka_langley", "plugsuit", 1},
		{"souryuu_asuka_langley", "plugsui", 1},
		{"souryuu_asuka_langley", "plugsuit", 101},
	}

	for _, entry := range tagTables {
		var sbuilder strings.Builder
		p, err := sb.GetByTagsRandRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			if err.Error() == "amount of requested posts too big" {
				continue
			} else {
				t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
			}
		}

		sbuilder.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sbuilder.WriteString(", ")
			}
			sbuilder.WriteString(pst.FileURL)
		}
		sbuilder.WriteString("\n")
		sbuilder.WriteString("\n")
		t.Logf(sbuilder.String())
	}
}

// getR34Conn creates and returns the connection to rule34.
func getR34Conn(t *testing.T) *R34API {
	*settings.StartupOptions.AESKey = decryptionKey
	decryptSettings(t)

	var client http.Client
	gb := NewR34(&client, settings.BotSettings.Addresses.Rule34URL)
	return gb
}

// TestNewR34 tests the creation of new rule34 API connection.
func TestNewR34(t *testing.T) {
	r34 := getR34Conn(t)

	p, err := r34.GetByTagsRaw([]string{"asuka_langley_sohryu", "plugsuit"}, 5)
	if err != nil {
		t.Error("Error while fetching the results")
	}

	t.Log("URL: " + p.List[0].FileURL)
}

// TestR34GetByTagsRaw tests if tag search can handle all of the possible cases of search works correctly.
func TestR34GetByTagsRaw(t *testing.T) {
	r34 := getR34Conn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"asuka_langley_sohryu", "plugsuit", 1},
		{"asuka_langley_sohryu", "plugsuit", 5},
		{"asuka_langley_sohryu", "plugsuit", 10},
		{"asuka_langley_sohryu", "plugsuit", 20},
		{"asuka_langley_sohryu", "plugsuit", 100},
		{"asukalangley_sohryu", "plugsuit", 1},
		{"asuka_langley_sohryu", "plugsui", 1},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := r34.GetByTagsRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// TestR34GetByTagsRandRaw tests if tag search with randomly sorted results works correctly.
func TestR34GetByTagsRandRaw(t *testing.T) {
	r34 := getR34Conn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"asuka_langley_sohryu", "plugsuit", 1},
		{"asuka_langley_sohryu", "plugsuit", 5},
		{"asuka_langley_sohryu", "plugsuit", 10},
		{"asuka_langley_sohryu", "plugsuit", 20},
		{"asuka_langley_sohryu", "plugsuit", 100},
		{"asukalangley_sohryu", "plugsuit", 1},
		{"asuka_langley_sohryu", "plugsui", 1},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := r34.GetByTagsRandRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			if err.Error() == "amount of requested posts too big" {
				continue
			} else {
				t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
			}
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// getDBConn creates and returns the connection to gelbooru.
func getYAConn(t *testing.T) *YAAPI {
	*settings.StartupOptions.AESKey = decryptionKey
	decryptSettings(t)

	var client http.Client
	ya := NewYA(&client, settings.BotSettings.Addresses.YandereURL)
	return ya
}

// TestNewYA tests the creation of new yande.re API connection
func TestNewYA(t *testing.T) {
	ya := getYAConn(t)

	p, err := ya.GetByTagsRaw([]string{"soryu_asuka_langley", "bodysuit"}, 5)
	if err != nil {
		// t.Error("Error while fetching the results")
		t.Error(err)
	}

	t.Log("URL: " + p.List[0].FileURL)
}

// TestYAGetByTagsRaw tests if tag search can handle all of the possible cases of search works correctly.
func TestYAGetByTagsRaw(t *testing.T) {
	ya := getYAConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"soryu_asuka_langley", "bodysuit", 1},
		{"soryu_asuka_langley", "bodysuit", 5},
		{"soryu_asuka_langley", "bodysuit", 10},
		{"soryu_asuka_langley", "bodysuit", 20},
		{"soryu_asuka_langley", "bodysuit", 100},
		{"soryuuasuka_langley", "bodysuit", 1},
		{"soryu_asuka_langley", "bodysui", 1},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := ya.GetByTagsRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}

// TestYAGetByTagsRandRaw tests if tag search with randomly sorted results works correctly.
func TestYAGetByTagsRandRaw(t *testing.T) {
	db := getDBConn(t)

	tagTables := []struct {
		firstTag     string
		secondTag    string
		resultAmount int
	}{
		{"soryu_asuka_langley", "bodysuit", 1},
		{"soryu_asuka_langley", "bodysuit", 5},
		{"soryu_asuka_langley", "bodysuit", 10},
		{"soryu_asuka_langley", "bodysuit", 20},
		{"soryu_asuka_langley", "bodysuit", 100},
		{"soryuuasuka_langley", "bodysuit", 1},
		{"soryu_asuka_langley", "bodysui", 1},
	}

	for _, entry := range tagTables {
		var sb strings.Builder
		p, err := db.GetByTagsRandRaw([]string{entry.firstTag, entry.secondTag}, entry.resultAmount)
		if err != nil {
			if err.Error() == "amount of requested posts too big" {
				continue
			} else {
				t.Error("Error while fetching the results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}")
			}
		}

		sb.WriteString("Results for: {" + entry.firstTag + ", " + entry.secondTag + ", " + strconv.Itoa(entry.resultAmount) + "}:\n")
		for i, pst := range p.List {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(pst.FileURL)
		}
		sb.WriteString("\n")
		sb.WriteString("\n")
		t.Logf(sb.String())
	}
}
