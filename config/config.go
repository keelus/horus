package config

import (
	"fmt"
	"horus/internal"
	"horus/models"
)

var UserConfiguration models.Configuration
var LedPresets models.LedPresets
var LedActive models.LedActive

func Init() {
	var err error
	UserConfiguration, err = internal.LoadUserConfiguration()
	if err != nil {
		fmt.Println("User configuration could not be loaded...")
	}

	LedPresets, err = internal.LoadLedPresets()
	if err != nil {
		fmt.Println("Led presets could not be loaded...")
	}

	LedActive, err = internal.LoadLedActive()
	if err != nil {
		fmt.Println("Active led presets could not be loaded...")
	}
}
