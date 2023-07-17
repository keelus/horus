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
//	# DESIGN & PERSONALISATION
//		Load and save system. Compatibility with scss/css

const CUR_VERSION = "0.7.8"

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
	server.Init()
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
