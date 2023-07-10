package models

// USER CONFIGURATION
type Configuration struct {
	Version         string
	UserInfo        UserInfo
	SessionSettings SessionSettings
	LedControl      [2]bool
	SystemStats     [5]bool
	Logging         bool
	Security        Security
	Units           Units
	Design          Design
}

type UserInfo struct {
	Username string
	Password string
}

type SessionSettings struct {
	Lifespan int
	Unit     string
}

type Security struct {
	UserInput bool
}

type Units struct {
	TimeMode12   bool
	TemperatureC bool
}

type Design struct {
	Accent []string
	Fonts  []FontDetails
}

type FontDetails struct {
	Title  string
	Source string
}

// LED
type LedPresets struct {
	StaticColor    []string
	CyclingColors  []string
	PulsatingColor []string
}
type LedActive struct {
	ActiveMode string
	Color      []string
	Brightness int
	Cooldown   int
}
