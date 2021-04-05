// Package settings provides a set of tools for configuration of the command's settings.
package settings

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	mathrand "math/rand"
	"os"
	"time"
)

// urls defines the urls used by the program. Those are mostly addresses to REST API servers.
type urls struct {
	DanbooruURL  string `json:"danbooruURL"`
	SafebooruURL string `json:"safebooruURL"`
	GelbooruURL  string `json:"gelbooruURL"`
	KonachanURL  string `json:"konachanURL"`
	Rule34URL    string `json:"rule34URL"`
	YandereURL   string `json:"yandereURL"`
}

// settings defines basic settings of the command.
type settings struct {
	BotToken          string `json:"botToken"`
	YTToken           string `json:"ytToken"`
	DanbooruToken     string `json:"danbooruToken"`
	DanbooruLogin     string `json:"danbooruLogin"`
	GelbooruToken     string `json:"gelbooruToken"`
	GelbooruUsrID     string `json:"gelbooruUsrID"`
	DatabaseIPAddress string `json:"databaseIPAddress"`
	DatabaseUsername  string `json:"databaseUsername"`
	DatabasePassword  string `json:"databasePassword"`
	DatabaseName      string `json:"databaseName"`
	Addresses         urls   `json:"urls"`
	SettingsProcessed bool   `json:"settingsProcessed"`
}

// BotSettings is the global instance of settings that is used to access all of the command configuration.
var BotSettings settings

func init() {
	loadCmdLineArgs()
	loadSettings()

	// If unit testing decrypt in the testing functions
	if flag.Lookup("test.v") != nil {
		return
	}

	if !BotSettings.SettingsProcessed {
		encryptSettings()
	} else {
		decryptSettings()
	}
}

// loadSettings Loads the settings in raw form from the specified settings file
func loadSettings() {
	// TODO: Make local variables with strings or add dereference
	sfp, sfpe := *StartupOptions.SettingsFilePath, *StartupOptions.SettingsFilePath+".encrypted"
	if _, err := os.Stat(sfp); os.IsNotExist(err) {
		if _, err := os.Stat(sfpe); os.IsNotExist(err) {
			createSettingsFile()
		}
	}

	if _, err := os.Stat(sfpe); os.IsNotExist(err) {
		dat, err := ioutil.ReadFile(sfp)
		if err != nil {
			fmt.Printf("[Error reading the bot settings file: %s]\r\n", err)
			panic(err)
		}

		err = json.Unmarshal([]byte(string(dat)), &BotSettings)
		if err != nil {
			fmt.Printf("[Error while unmarshaling the bot settings file: %s]\r\n", err)
			panic(err)
		}
	} else {
		BotSettings.SettingsProcessed = true
	} // If the settings are encrypted and this is not the first launch, the settings will be loaded in decryption stage
}

// decryptSettings loads an encrypted file, decrypts it and puts it in the BotSettings.
func decryptSettings() {
	data, err := ioutil.ReadFile(*StartupOptions.SettingsFilePath + ".encrypted")
	if err != nil {
		fmt.Printf("[Error while reading the encrypted settings file: %s]\r\n", err)
		panic(err)
	}
	key := CreateHash(*StartupOptions.AESKey)
	dataDecrypted := Decrypt(data, key)
	fmt.Println(string(dataDecrypted))
	err = json.Unmarshal(dataDecrypted, &BotSettings)
	if err != nil {
		fmt.Printf("[Error while unmarshaling the encrypted settings file: %s]\r\n", err)
		fmt.Printf("Have you used the [-decryption-key key] argument?")
		panic(err)
	}
}

// encryptSettings encrypts all of the necessary settings data and saves it to a file.
func encryptSettings() {
	// Get some random passphrase and turn it into a MD5 hashed key
	passphrase := string(encryptGetKey())
	key := CreateHash(passphrase)

	// Marshal into json and encrypt the marshaled data
	data, err := json.Marshal(BotSettings)
	if err != nil {
		fmt.Printf("[Error while marshaling the BotSettings: %s]\r\n", err)
		panic(err)
	}
	encryptedData := Encrypt(data, key)

	// Save the changes and delete the json file
	err = ioutil.WriteFile(*StartupOptions.SettingsFilePath+".encrypted", encryptedData, 0644)
	if err != nil {
		fmt.Printf("[Error while saving encrypted settings: %s]\r\n", err)
		panic(err)
	}
	_ = os.Remove(*StartupOptions.SettingsFilePath)

	BotSettings.SettingsProcessed = true
	StartupOptions.AESKey = &passphrase // TODO: Does this disappear because of the end of the scope

	// Notify user about the changes
	fmt.Println("ENCRYPTED THE '" + *StartupOptions.SettingsFilePath + "' FILE TO '" + *StartupOptions.SettingsFilePath + ".encrypted' WITH KEY: '" + *StartupOptions.AESKey + "' !!!")
	fmt.Println("Please save this key and use it each time you run the bot to decrypt your bot data from '" + *StartupOptions.SettingsFilePath + ".encrypted'")
}

// encryptGetKey generates a 32-bit ASCII encryption key.
func encryptGetKey() []byte {
	mathrand.Seed(time.Now().UnixNano())
	key := make([]byte, 12) // TODO: Check if works after changing to 64 from 32
	// letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letters := "ABCDEF0123456789"
	for i := 0; i < len(key); i++ {
		key[i] = letters[mathrand.Intn(len(letters))]
	}
	return key
}

// createSettingsFile creates the default settings file.
func createSettingsFile() {
	sfc := []byte(`{
	"botToken": "token",
	"ytToken": "token",
	"danbooruToken": "token",
	"danbooruLogin": "login",
	"gelbooruToken": "token",
	"gelbooruUsrID": "user id",
	"databaseIPAddress": "localhost",
	"databaseUsername": "username",
	"databasePassword": "password",
	"databaseName": "kitsune",
	"urls": {
		"danbooruURL": "https://danbooru.donmai.us/",
		"safebooruURL": "https://safebooru.org/",
		"gelbooruURL": "https://gelbooru.com/",
		"konachanURL": "https://konachan.com/",
		"rule34URL": "https://rule34.xxx/",
		"yandereURL": "https://yande.re/"
	},
	"settingsProcessed": false
}`)
	err := ioutil.WriteFile(*StartupOptions.SettingsFilePath, sfc, 0644)
	if err != nil {
		fmt.Printf("[Error creating a new settings file: %s]\r\n", err)
		panic(err)
	}
	fmt.Println("A new settings file has been created under the path: " + *StartupOptions.SettingsFilePath + ". Please fill out the settings and re-run the bot.")
	os.Exit(0)
}

// CreateHash creates a hash of the encryption key
func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt encrypts a sequence of bytes with a given passphrase. Uses Galois/Counter Mode.
func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

// Decrypt decrypts a sequence of bytes with a given passphrase. Uses Galois/Counter Mode.
func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(CreateHash(passphrase))
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, cipherText, nil)
	return plaintext
}
