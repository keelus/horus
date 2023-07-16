package models

// USER CONFIGURATION
type Configuration struct {
	Version         string          `yaml:"Version"`
	UserInfo        UserInfo        `yaml:"UserInfo"`
	SessionSettings SessionSettings `yaml:"SessionSettings"`
	LedControl      [3]bool         `yaml:"LedControl"`
	SystemStats     [5]bool         `yaml:"SystemStats"`
	Logging         bool            `yaml:"Logging"`
	Security        Security        `yaml:"Security"`
	Units           Units           `yaml:"Units"`
	Design          Design          `yaml:"Design"`
}

type UserInfo struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
}

type SessionSettings struct {
	Lifespan int    `yaml:"Lifespan"`
	Unit     string `yaml:"Unit"`
}

type Security struct {
	UserInput bool `yaml:"UserInput"`
}

type Units struct {
	TimeMode12   bool `yaml:"TimeMode12"`
	TemperatureC bool `yaml:"TemperatureC"`
}

type Design struct {
	Accent []string      `yaml:"Accent"`
	Fonts  []FontDetails `yaml:"Fonts"`
}

type FontDetails struct {
	Title  string `yaml:"Title"`
	Source string `yaml:"Source"`
}

// LED
type LedPresets struct {
	StaticColor    []string `yaml:"StaticColor"`
	StaticGradient []string `yaml:"StaticGradient"`
	FadingRainbow  int      `yaml:"FadingRainbow"`
	BreathingColor []string `yaml:"BreathingColor"`
}
type LedActive struct {
	ActiveMode string   `yaml:"ActiveMode"`
	Color      []string `yaml:"Color"`
	Brightness int      `yaml:"Brightness"`
	Cooldown   int      `yaml:"Cooldown"`
}
