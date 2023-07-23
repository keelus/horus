package models

// USER CONFIGURATION
type Configuration struct {
	UserInfo struct {
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	} `yaml:"UserInfo"`
	ColorModeDark   bool `yaml:"ColorModeDark"`
	SessionSettings struct {
		Lifespan int    `yaml:"Lifespan"`
		Unit     string `yaml:"Unit"`
	} `yaml:"SessionSettings"`
	LedControl  [3]bool `yaml:"LedControl"`
	SystemStats [5]bool `yaml:"SystemStats"`
	Logging     bool    `yaml:"Logging"`
	Security    struct {
		UserInput bool `yaml:"UserInput"`
	} `yaml:"Security"`
	Units struct {
		TemperatureC bool `yaml:"TemperatureC"`
	} `yaml:"Units"`
}

// LED
type LedPresets struct {
	StaticColor    []string   `yaml:"StaticColor"`
	StaticGradient [][]string `yaml:"StaticGradient"`
	FadingRainbow  int        `yaml:"FadingRainbow"`
	BreathingColor struct {
		Cooldown int      `yaml:"Cooldown"`
		Colors   []string `yaml:"Colors"`
	} `yaml:"BreathingColor"`
}
type LedActive struct {
	ActiveMode string   `yaml:"ActiveMode"`
	Color      []string `yaml:"Color"`
	Brightness int      `yaml:"Brightness"`
	LedAmount  int      `yaml:"LedAmount"`
	Cooldown   int      `yaml:"Cooldown"`
}
