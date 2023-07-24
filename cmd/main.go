package main

import (
	"bufio"
	"fmt"
	"horus/config"
	"horus/internal"
	"horus/led"
	"horus/logger"
	"horus/server"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
)

var BlueColor = color.New(color.FgBlue, color.Bold)
var RedColor = color.New(color.FgRed, color.Bold)
var GreenColor = color.New(color.FgGreen, color.Bold)
var YellowColor = color.New(color.FgYellow, color.Bold)

func init() {
	YellowColor.Println(`
	⠀⠀⠀⠀⣀⣤⣶⠾⠿⠿⠿⠿⢶⣦⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
	⠀⠀⣤⠾⠛⠉⠀⠀⠀⠀⠀⠀⠀⠀⠉⠙⠛⠻⠷⣶⣤⣤⣤⣀⣀⣀⣀⣀⠀⠀  ██╗  ██╗ ██████╗ ██████╗ ██╗   ██╗███████╗
	⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠉⠉⠉⠉⠉⠀⠀	██║  ██║██╔═══██╗██╔══██╗██║   ██║██╔════╝
	⠀⠀⠀⠀⣠⡾⢛⣽⣿⣿⣏⠙⠛⠻⠷⣦⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⡀⠀⠀	███████║██║   ██║██████╔╝██║   ██║███████╗
	⠀⠀⢠⣾⣋⡀⢸⣿⣿⣿⣿⠀⠀⢀⣀⣤⣽⡿⠿⠛⠿⠿⠷⠾⠿⠿⠛⠋⠀⠀	██╔══██║██║   ██║██╔══██╗██║   ██║╚════██║
	⠀⠀⠻⠛⠛⠻⣶⣽⣿⣿⣿⡶⠿⠛⠋⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀	██║  ██║╚██████╔╝██║  ██║╚██████╔╝███████║
	⠀⠀⠀⠀⠀⠀⣠⣿⡏⠻⣷⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣶⠶⢶⣤⠀⠀⠀	╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝
	⠀⠀⠀⠀⠀⠀⢹⣯⠁⠀⠈⠛⢷⣤⡀⠀⠀⠀⠀⠀⠀⠀⠸⠧⠀⠀⢹⡇⠀⠀	
	⠀⠀⠀⠀⠀⠀⠈⣿⠀⠀⠀⠀⠀⠉⠻⠷⣦⣤⣤⣀⣀⣀⣀⣠⣤⡶⠟⠀⠀⠀	Initializing version ` + internal.VERSION_CURRENT + `...
	⠀⠀⠀⠀⠀⠀⠀⠛⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠉⠉⠉⠉⠉⠁⠀⠀⠀⠀⠀`)

	if fileExists("config/userConfig.yaml") {
		GreenColor.Println("✓ User configuration file found")
		config.Init()

		if fileExists("config/ledActive.yaml") {
			GreenColor.Println("✓ Led Active file found")
		} else {
			RedColor.Println("⨉ Led Active file not found, generating a new one...")
			internal.SaveFile(&config.LedActive)
		}

		if fileExists("config/ledPresets.yaml") {
			GreenColor.Println("✓ Led Presets file found")
		} else {
			RedColor.Println("⨉ Led Presets file not found, generating a new one...")
			internal.SaveFile(&config.LedPresets)
		}

	} else {
		RedColor.Println("⨉ - User configuration file not found")
		config.Init()
		setupConfig()
	}
}

func main() {
	logger.Init()

	// Check for sudo
	euid := os.Geteuid()
	if euid == 0 {
		BlueColor.Println("✓ Sudo privileges found, starting WS281X Led module...")
		led.Init()
	} else {
		RedColor.Println("⨉ Sudo privileges were not found. They are required to access the GPIO of your Raspberry Pi, in order to control your WS281X Led strip... WS281X Led module not initialized.")
		logger.Log(nil, logger.ERROR, "⨉ Sudo privileges were not found. They are required to access the GPIO of your Raspberry Pi, in order to control your WS281X Led strip... WS281X Led module not initialized.")
	}

	userOS := runtime.GOOS
	if userOS != "linux" {
		RedColor.Printf("⨉ '%s' OS detected. Keep in mind this software is Linux & Raspberry Pi focused. Most features won't work/take effect on other systems/operative systems.\n", strings.ToUpper(userOS))
		logger.Log(nil, logger.ERROR, fmt.Sprintf("⨉ '%s' OS detected. Keep in mind this software is Linux & Raspberry Pi focused. Most features won't work/take effect on other systems/operative systems.", strings.ToUpper(userOS)))
	} else {
		GreenColor.Println("✓ Linux detected (keep in mind some features only work on Raspberry Pi systems).")
		logger.Log(nil, logger.POWERON, "✓ Linux detected (keep in mind some features only work on Raspberry Pi systems).")
	}

	if euid == 0 {
		BlueColor.Println("✓ Sudo privileges found, initializing GIN web server...")
		logger.Log(nil, logger.UP, "Server UP and running on port :80.") // TODO: Check
	} else {
		RedColor.Println("⨉ Sudo privileges were not found. They are required to access to the root-privileged port 80. Gin web server couldn't be initialized.") // TODO: Be able to change the port and use a non-privileged one.
		logger.Log(nil, logger.ERROR, "⨉ Sudo privileges were not found. They are required to access the GPIO of your Raspberry Pi, in order to control your WS281X Led strip... WS281X Led module not initialized.")
		os.Exit(-1)
	}

	internal.CheckLatestVersion()
	time.Sleep(1 * time.Second)
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
	YellowColor.Println("Horus first time setup:")
	var username string
	var password string

	for {
		username = readLine("username")
		if username != "" {
			break
		}
	}

	for {
		password = readLine("password")
		if username != "" {
			break
		}
	}

	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Unexpected error happened while encrypting your password.")
		os.Exit(-1)
	}

	config.UserConfiguration.UserInfo.Username = username
	config.UserConfiguration.UserInfo.Password = string(cryptedPassword)

	config.SaveAll()
	GreenColor.Println("✓ Configuration saved.")
}

func readLine(field string) string {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	switch field {
	case "username":
		YellowColor.Print("\t-Username:")
	case "password":
		YellowColor.Print("\t-Password:")
		fmt.Print("\033[8m")
	}

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		signal.Stop(sigChan)

		switch field {
		case "username":
			if len(scanner.Text()) < 3 || len(scanner.Text()) > 20 {
				fmt.Println("The username must be at least 3 characters long, with a maximum of 20.")
				return ""
			} else {
				return scanner.Text()
			}
		case "password":
			if len(scanner.Text()) < 4 {
				fmt.Println("The password must be at least 4 characters long.")
				return ""
			} else {
				fmt.Print("\033[28m")
				return scanner.Text()
			}
		}

	}

	signal.Stop(sigChan)

	if field == "password" {
		fmt.Print("\033[28m") // Prevent empty output, as password is hidden
	}

	RedColor.Println("\nHorus set up process cancelled. Exiting...")
	os.Exit(0)
	return ""
}
