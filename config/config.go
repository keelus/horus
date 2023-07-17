package config

import (
	"horus/internal"
	"horus/models"
)

var UserConfiguration models.Configuration
var LedPresets models.LedPresets
var LedActive models.LedActive

func Init() {
	var err error

	err = internal.LoadFile(&UserConfiguration, false)
	if err != nil { // File not found. Use template
		err = internal.LoadFile(&UserConfiguration, true)
	}
	err = internal.LoadFile(&LedPresets, false)
	if err != nil {
		err = internal.LoadFile(&LedPresets, true)
	}
	err = internal.LoadFile(&LedActive, false)
	if err != nil {
		err = internal.LoadFile(&LedActive, true)
	}
}

func SaveAll() {
	internal.SaveFile(&UserConfiguration)
	internal.SaveFile(&LedPresets)
	internal.SaveFile(&LedActive)
}
