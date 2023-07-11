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

	err = internal.LoadFile(&UserConfiguration)
	if err != nil {
		fmt.Println("User configuration could not be loaded...")
	}
	err = internal.LoadFile(&LedPresets)
	if err != nil {
		fmt.Println("Led presets could not be loaded...")
	}
	err = internal.LoadFile(&LedActive)
	if err != nil {
		fmt.Println("Led active could not be loaded...")
	}
}
