package main

import (
	"fmt"
	"horus/config"
	"horus/led"
	"horus/server"
	"os"

	"golang.org/x/crypto/bcrypt"
)

//	TODO:
//	# USER CONFIGURATION
//		Switch from .JSON to .YAML
//	# DESIGN & PERSONALISATION
//		Load and save system. Compatibility with scss/css

const CUR_VERSION = "0.7.0"

func init() {
	fmt.Println("Initializing horus...")
	if !fileExists("config/userConfig.yaml") {
		fmt.Println("User configuration has not been found. Initializing horus setup...")
		config.Init()
		setupConfig()
		config.SaveAll()
		// os.Exit(-1)
	} else {
		config.Init()
	}

}

func main() {
	led.Init()
	// config.Init()
	server.Init()

	// Read the YAML file
	// yamlFile, err := ioutil.ReadFile("config/defaults/userConfig.yaml")
	// if err != nil {
	// 	log.Fatalf("Failed to read YAML file: %v", err)
	// }

	// // Unmarshal the YAML into a Config struct
	// var config Config
	// err = yaml.Unmarshal(yamlFile, &config)
	// if err != nil {
	// 	log.Fatalf("Failed to unmarshal YAML: %v", err)
	// }

	// // Print the contents of the config
	// fmt.Printf("Version: %s\n", config.Version)
	// fmt.Printf("Username: %s\n", config.UserInfo.Username)
	// fmt.Printf("Password: %s\n", config.UserInfo.Password)
	// fmt.Printf("Lifespan: %d %s\n", config.SessionSettings.Lifespan, config.SessionSettings.Unit)
	// fmt.Printf("LedControl: %v\n", config.LedControl)
	// fmt.Printf("SystemStats: %v\n", config.SystemStats)
	// fmt.Printf("Logging: %v\n", config.Logging)
	// fmt.Printf("UserInput: %v\n", config.Security.UserInput)
	// fmt.Printf("TimeMode12: %v\n", config.Units.TimeMode12)
	// fmt.Printf("TemperatureC: %v\n", config.Units.TemperatureC)
	// fmt.Printf("Accent: %v\n", config.Design.Accent)
	// for _, font := range config.Design.Fonts {
	// 	fmt.Printf("Font - Title: %s, Source: %s\n", font.Title, font.Source)
	// }
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func setupConfig() {
	var username string
	var password string

	fmt.Println("HORUS setup:")

	fmt.Print("\t-Username:")
	fmt.Scan(&username)
	fmt.Print("\t-Password:\033[8m")

	fmt.Scan(&password)
	fmt.Print("\033[28m")

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Unexpected error happened while encrypting your password.")
		os.Exit(-1)
	}

	config.UserConfiguration.UserInfo.Username = username
	config.UserConfiguration.UserInfo.Password = string(cryptedPassword)

	fmt.Println("Configuration done. Initializing horus...")
}
