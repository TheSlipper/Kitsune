// Package settings provides a set of tools for configuration of the command's settings.
package settings

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// startupOptions contains all of the parsed command line data.
type startupOptions struct {
	VersionFlag      *bool
	HelpFlag         *bool
	AESKey           *string // Key used for AES ciphertext decryption
	SettingsFilePath *string // Path to the json settings file of the file. "settings.json" on default.
}

// StartupOptions is the global instance of startupOptions that is used to command line data.
var StartupOptions startupOptions

// loadCmdLineArgs parses the command line arguments.
func loadCmdLineArgs() {
	StartupOptions.VersionFlag = flag.Bool("version", false, "Prints out information about the version and the build of the executable.")
	StartupOptions.HelpFlag = flag.Bool("help", false, "Prints out the help page.")
	StartupOptions.AESKey = flag.String("decryption-key", "decryption_key", "Specifies the decryption key used for decrypting the settings file.")
	StartupOptions.SettingsFilePath = flag.String("settings-path", "settings.json", "Specifies a custom path to the configuration file.")
	flag.Parse()

	if *StartupOptions.VersionFlag {
		printVersion()
	}
}

// getArgVal extracts the data for the specified option.
func getArgVal(fcmd string, cmd string, scmd string, i *int) string {
	if strings.HasPrefix(os.Args[*i], fcmd) {
		a := strings.IndexRune(fcmd, '=')
		a = a + 1
		return os.Args[*i][a:]
	} else if os.Args[*i] == cmd || os.Args[*i] == scmd {
		*i = *i + 1
		return os.Args[*i]
	}
	fmt.Println("Incorrect argument '" + os.Args[*i] + "'!")
	os.Exit(1)
	return ""
}

// printVersion prints out the version and build of the command.
func printVersion() {
	fmt.Println("0.0.1-alpha")
	fmt.Println("github.com/TheSlipper build")
	os.Exit(0)
}
