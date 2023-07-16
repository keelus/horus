package server

import (
	"encoding/json"
	"fmt"
	"horus/internal"
	"horus/led"
	"net/http"
	"strconv"
	"strings"

	"horus/config"

	"github.com/forestgiant/sliceutil"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Init() {
	r := gin.Default()

	cookieAge := 0

	if config.UserConfiguration.SessionSettings.Lifespan == -1 || config.UserConfiguration.SessionSettings.Unit == "" {
		cookieAge = 0
	} else {
		if config.UserConfiguration.SessionSettings.Unit == "min" {
			cookieAge = 60 * config.UserConfiguration.SessionSettings.Lifespan
		} else if config.UserConfiguration.SessionSettings.Unit == "hour" {
			cookieAge = 60 * 60 * config.UserConfiguration.SessionSettings.Lifespan
		} else if config.UserConfiguration.SessionSettings.Unit == "day" {
			cookieAge = 60 * 60 * 24 * config.UserConfiguration.SessionSettings.Lifespan
		}
	}

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path: "/",
		// Secure:   true,
		Secure: false,
		// HttpOnly: true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieAge,
	})
	r.Use(sessions.Sessions("sesion", store))

	r.HTMLRender = Renderer()
	r.Use(gin.Recovery())

	r.Static("/static", "web/static")

	// ##### FRONT PAGES #####
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/panel/main")
	})
	r.GET("/panel", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/panel/main")
	})
	r.GET("/login", func(c *gin.Context) {
		if internal.IsLogged(c) {
			c.Redirect(http.StatusTemporaryRedirect, "/panel/main")
			return
		}
		c.HTML(http.StatusOK, "login", gin.H{
			"UserConfiguration": config.UserConfiguration,
		})
	})
	r.GET("/panel/:category", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		latestVersion := "X.X.X" // TODO
		usingLatest := true

		if latestVersion > config.UserConfiguration.Version {
			usingLatest = false
		}

		// colors := []string{"#48d051", "#c05bef", "#0fcbdc", "#8c799c", "#12c670", "#7a53a7", "#7d29a2", "#5f552b", "#5191f2", "#03ba62", "#f69d97", "#4cc856", "#a1e180", "#c30113", "#85b864", "#b5b437", "#d51bf5", "#e13ad1", "#466acd", "#ecc95f", "#703fdc"}

		category := c.Param("category")
		c.HTML(http.StatusOK, "panel", gin.H{
			"Active":            category,
			"UserConfiguration": config.UserConfiguration,
			"LatestVersion":     latestVersion,
			"UsingLatest":       usingLatest,
			"CurrentLED":        "#00FF00",
			"LedPresets":        config.LedPresets,
			"LedActive":         config.LedActive,
		})
	})

	r.POST("/back/login", func(c *gin.Context) {
		session := sessions.Default(c)

		username := c.PostForm("Username")
		password := c.PostForm("Password")

		// cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"Status": "error",
		// 	})
		// }

		hasError := false

		if config.UserConfiguration.Security.UserInput {
			if username != config.UserConfiguration.UserInfo.Username {
				hasError = true
			}
		}

		err := bcrypt.CompareHashAndPassword([]byte(config.UserConfiguration.UserInfo.Password), []byte(password))
		if err != nil {
			hasError = true
		}

		if hasError {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Status": "error",
			})
		} else {
			session.Clear()
			session.Set("LoggedIn", true)
			session.Save()
			c.JSON(http.StatusOK, gin.H{
				"Status": "OK",
			})
		}
	})

	r.GET("/back/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	r.POST("/back/saveConfiguration/:category", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in.", "field": ""})
			return
		}

		category := c.Param("category")

		returnError := []map[string]string{}

		switch category {
		case "UserInfo":
			username := c.PostForm("Username")
			password := c.PostForm("Password")

			if len(username) < 3 || len(username) > 20 {
				returnError = append(returnError, map[string]string{"details": "The username must be at least 3 characters long, with a maximum of 20.", "field": "userInfo0"})
			}

			if strings.TrimSpace(password) != "" {
				if len(password) < 4 {
					returnError = append(returnError, map[string]string{"details": "The password must be at least 4 characters long.", "field": "userInfo1"})
				}
			}

			if len(returnError) == 0 {
				if strings.TrimSpace(password) != "" {
					cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
					if err != nil {
						returnError = append(returnError, map[string]string{"details": "The password must be at least 4 characters long.", "field": "userInfo1"})
					} else { // If we have no error, we will store the username, as well as the password
						config.UserConfiguration.UserInfo.Username = username
						config.UserConfiguration.UserInfo.Password = string(cryptedPassword)
					}
				} else { // If there is no password to store, username will always be store
					config.UserConfiguration.UserInfo.Username = username
				}
			}
			break
		case "SessionSettings":
			lifespanString := c.PostForm("Lifespan")
			unit := c.PostForm("Unit")

			allowedUnits := []string{"min", "hour", "day", ""}

			fmt.Printf("lifespan:'%s'\n", lifespanString)
			fmt.Printf("unit:'%s'\n", unit)
			if !sliceutil.Contains(allowedUnits, unit) { // Replace to array instead of slice TODO
				returnError = append(returnError, map[string]string{"details": "Unexpected time unit.", "field": "sessionDuration3"})
			}

			lifespan, err := strconv.Atoi(lifespanString)
			if err != nil {
				returnError = append(returnError, map[string]string{"details": "Unexpected lifespan format.", "field": "sessionDuration2"})
			}

			if len(returnError) == 0 {
				if unit == "" {
					config.UserConfiguration.SessionSettings.Unit = ""
					config.UserConfiguration.SessionSettings.Lifespan = -1
				} else {
					config.UserConfiguration.SessionSettings.Unit = unit
					config.UserConfiguration.SessionSettings.Lifespan = lifespan
				}
			}
			break
		case "LedControl":
			ledControlString := c.PostForm("LedControl")

			var ledControl [3]bool
			err := json.Unmarshal([]byte(ledControlString), &ledControl)
			if err != nil {
				returnError = append(returnError, map[string]string{"details": "Unexpected format.", "field": ""})
			}

			if len(returnError) == 0 {
				config.UserConfiguration.LedControl = ledControl
			}
			break
		case "SystemStats":
			systemStatsString := c.PostForm("SystemStats")

			var systemStats [5]bool
			err := json.Unmarshal([]byte(systemStatsString), &systemStats)
			if err != nil {
				returnError = append(returnError, map[string]string{"details": "Unexpected format.", "field": ""})
			}

			if len(returnError) == 0 {
				config.UserConfiguration.SystemStats = systemStats
			}
			break
		case "Logging":
			loggingString := c.PostForm("Logging")

			logging, err := strconv.ParseBool(loggingString)
			if err != nil {
				returnError = append(returnError, map[string]string{"details": "Unexpected format.", "field": ""})
			}

			if len(returnError) == 0 {
				config.UserConfiguration.Logging = logging
			}
			break
		case "Security":
			userInputString := c.PostForm("UserInput")

			userInput, err := strconv.ParseBool(userInputString)
			if err != nil {
				returnError = append(returnError, map[string]string{"details": "Unexpected format.", "field": ""})
			}

			if len(returnError) == 0 {
				config.UserConfiguration.Security.UserInput = userInput
			}
			break
		case "Units":
			timeMode := c.PostForm("TimeMode")
			temperature := c.PostForm("Temperature")

			allowedTimeMode := []string{"12", "24"}
			allowedTemperature := []string{"C", "F"}

			if !sliceutil.Contains(allowedTimeMode, timeMode) { // Replace to array instead of slice TODO
				returnError = append(returnError, map[string]string{"details": "Unexpected time unit.", "field": "units0"})
			}
			if !sliceutil.Contains(allowedTemperature, temperature) { // Replace to array instead of slice TODO
				returnError = append(returnError, map[string]string{"details": "Unexpected temperature unit.", "field": "units1"})
			}

			if len(returnError) == 0 {
				config.UserConfiguration.Units.TimeMode12 = timeMode == "12"
				config.UserConfiguration.Units.TemperatureC = temperature == "C"
			}
			break
		case "Design":
			// TODO
			break
		}

		if len(returnError) == 0 {
			err := internal.SaveFile(&config.UserConfiguration)
			if err != nil {
				c.JSON(http.StatusInternalServerError, map[string]string{"details": "Error while saving the user configuration. 500", "field": ""})
			} else {
				c.Status(http.StatusOK)
			}
			return
		}

		c.JSON(http.StatusBadRequest, returnError)
	})

	r.POST("/back/ledControl/cooldown/:mode/:amount", func(c *gin.Context) {
		amountStr := c.Param("amount")
		mode := c.Param("mode")
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"details": "Unexpected type of amount."})
			return
		}

		if mode == "BreathingColor" {
			if config.LedActive.ActiveMode != "BreathingColor" {
				go led.BreathingColor()
			}

			//config.LedPresets.BreathingColor TODO save cooldown
		} else if mode == "FadingRainbow" {
			if config.LedActive.ActiveMode != "FadingRainbow" {
				go led.Rainbow()
			}
			config.LedPresets.FadingRainbow = amount
		}

		config.LedActive.Cooldown = amount
		internal.SaveFile(&config.LedActive)
		internal.SaveFile(&config.LedPresets)
		c.Status(http.StatusOK)
	})

	r.POST("/back/ledControl/activate/:mode", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
			return
		}

		mode := c.Param("mode")

		if mode == "StaticColor" {
			// By default first color is activated. Always will be one at least.
			led.SetColor([]string{config.LedPresets.StaticColor[0]}) // By default first color will be initialized
			config.LedActive.ActiveMode = "StaticColor"
			config.LedActive.Cooldown = 0
		} else if mode == "StaticGradient" {
			config.LedActive.ActiveMode = "StaticGradient"
			config.LedActive.Color = config.LedPresets.StaticGradient[0] // All colors are activated on Fading Colors
			config.LedActive.Cooldown = 0                                // TODO
			led.DrawGradient()
		} else if mode == "FadingRainbow" {
			previousMode := config.LedActive.ActiveMode
			config.LedActive.ActiveMode = "FadingRainbow"
			config.LedActive.Color = []string{"000000"}
			config.LedActive.Cooldown = config.LedPresets.FadingRainbow

			if previousMode != "FadingRainbow" {
				if previousMode == "BreathingColor" {
					led.StopBreathing = true
				}
				go led.Rainbow()
			}
		} else if mode == "BreathingColor" {
			previousMode := config.LedActive.ActiveMode
			config.LedActive.ActiveMode = "BreathingColor"
			config.LedActive.Color = []string{config.LedPresets.BreathingColor[0]} // By default first color will be initialized
			config.LedActive.Cooldown = 0

			if previousMode != "BreathingColor" {
				if previousMode == "FadingRainbow" {
					led.StopRainbow = true
				}
				go led.BreathingColor()
			}
		}

		internal.SaveFile(&config.LedActive)
		c.JSON(http.StatusOK, gin.H{"Color": config.LedActive.Color, "Brightness": config.LedActive.Brightness, "Cooldown": config.LedActive.Cooldown})
	})

	r.POST("/back/ledControl/activateGradient", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
			return
		}

		rawGradient := c.PostForm("rawGradient")

		if !internal.GradientExists(rawGradient, config.LedPresets.StaticGradient) {
			c.JSON(http.StatusBadRequest, gin.H{"details": "That gradient doesn't exist."})
			return
		}

		for _, gradient := range config.LedPresets.StaticGradient {
			if internal.GetGradientStr(gradient) == rawGradient {
				config.LedActive.Color = gradient
			}
		}

		led.DrawGradient()

		internal.SaveFile(&config.LedPresets)
		internal.SaveFile(&config.LedActive)
		c.Status(http.StatusOK)

	})

	r.POST("/back/ledControl/activate/:mode/:hex", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
			return
		}

		mode := c.Param("mode")
		hex := c.Param("hex")

		if mode == "StaticGradient" {
			c.JSON(http.StatusBadRequest, gin.H{"details": "Unexpected mode static gradient."})
			return
		} else {
			led.SetColor([]string{hex})
			config.LedActive.Color = []string{hex}
		}

		c.Status(http.StatusOK)
	})

	r.POST("/back/ledControl/brightness/:valuePercent", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
			return
		}

		valuePercent, err := strconv.Atoi(c.Param("valuePercent"))
		if err != nil || !(valuePercent >= 0 && valuePercent <= 100) {
			c.JSON(http.StatusBadGateway, gin.H{"details": "Unexpected brightness value type."})
			return
		}

		led.SetBrightness(int(valuePercent * 255 / 100))

		internal.SaveFile(&config.LedActive)
		c.JSON(http.StatusOK, gin.H{"Color": config.LedActive.Color, "Brightness": config.LedActive.Brightness, "Cooldown": config.LedActive.Cooldown})
	})

	r.POST("/back/ledControl/deleteGradient", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
			return
		}

		rawGradient := c.PostForm("rawGradient")
		fmt.Println(rawGradient)

		if !internal.GradientExists(rawGradient, config.LedPresets.StaticGradient) {
			c.JSON(http.StatusBadRequest, gin.H{"details": "That gradient doesn't exist."})
			return
		}

		newGradientSlice := [][]string{}
		for _, gradient := range config.LedPresets.StaticGradient {
			if internal.GetGradientStr(gradient) != rawGradient {
				newGradientSlice = append(newGradientSlice, gradient)
			}
		}

		config.LedPresets.StaticGradient = newGradientSlice

		internal.SaveFile(&config.LedPresets)
		internal.SaveFile(&config.LedActive) // TODO activate first & check minimum 1
		c.Status(http.StatusOK)

	})
	r.POST("/back/ledControl/delete/:mode/:hex", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
			return
		}

		mode := c.Param("mode")
		hex := c.Param("hex")

		if mode == "StaticGradient" {
			c.JSON(http.StatusBadRequest, gin.H{"details": "Unexpected mode static gradient."})
			return
		}

		if mode == "StaticColor" {
			if len(config.LedPresets.StaticColor) == 1 {
				c.JSON(http.StatusBadRequest, gin.H{"details": "There must be at least 1 preset color."})
				return
			}
			newPreset := []string{}

			for _, color := range config.LedPresets.StaticColor {
				if color != hex {
					newPreset = append(newPreset, color)
				}
			}

			if config.LedActive.Color[0] == hex {
				config.LedActive.Color[0] = newPreset[0]
			}

			config.LedPresets.StaticColor = newPreset
		} else if mode == "BreathingColor" {
			if len(config.LedPresets.BreathingColor) == 1 {
				c.JSON(http.StatusBadRequest, gin.H{"details": "There must be at least 1 preset color."})
				return
			}
			newPreset := []string{}

			for _, color := range config.LedPresets.BreathingColor {
				if color != hex {
					newPreset = append(newPreset, color)
				}
			}

			if config.LedActive.Color[0] == hex {
				config.LedActive.Color[0] = newPreset[0]
			}

			config.LedPresets.BreathingColor = newPreset
		} else if mode == "StaticGradient" {
			// if len(config.LedPresets.StaticGradient) == 2 {
			// 	c.JSON(http.StatusBadRequest, gin.H{"details": "There must be at least 2 preset colors."})
			// 	return
			// }

			// newPreset := []string{}
			// for _, color := range config.LedPresets.StaticGradient {
			// 	if color != hex {
			// 		newPreset = append(newPreset, color)
			// 	}
			// }

			// config.LedActive.Color = newPreset
			// config.LedPresets.StaticGradient = newPreset
		}

		internal.SaveFile(&config.LedActive)
		internal.SaveFile(&config.LedPresets)
		c.Status(http.StatusOK)
	})

	r.POST("/back/ledControl/addGradient/", func(c *gin.Context) {
		var hexValues []string
		hexValuesStr := c.PostForm("hexValues")

		err := json.Unmarshal([]byte(hexValuesStr), &hexValues)
		if err != nil {
			fmt.Println("Error! ", err) // TODO
			return
		}

		rawGradient := internal.GetGradientStr(hexValues)
		if internal.GradientExists(rawGradient, config.LedPresets.StaticGradient) {
			c.JSON(http.StatusBadRequest, gin.H{"details": "That exact gradient already exists."})
			return
		}

		config.LedPresets.StaticGradient = append(config.LedPresets.StaticGradient, hexValues)
		internal.SaveFile(&config.LedPresets)
		internal.SaveFile(&config.LedActive) // TODO activate new gradient
		c.Status(http.StatusOK)
	})

	r.POST("/back/ledControl/add/:mode/:hex", func(c *gin.Context) {
		if !internal.IsLogged(c) {
			c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
			return
		}

		mode := c.Param("mode")
		hex := strings.ToUpper(c.Param("hex"))

		hexChars := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}

		bytes := []byte(hex)

		for _, b := range bytes {
			if !sliceutil.Contains(hexChars, b) {
				c.JSON(http.StatusForbidden, gin.H{"details": fmt.Sprintf("Hex byte '%c' not allowed.", b)})
				return
			}
		}

		if mode == "StaticGradient" {
			c.JSON(http.StatusBadRequest, gin.H{"details": "Unexpected mode static gradient."})
			return
		}

		if mode == "StaticColor" {
			if sliceutil.Contains(config.LedPresets.StaticColor, hex) {
				c.JSON(http.StatusBadRequest, gin.H{"details": "That color has been already added to this mode."})
				return
			}

			newPreset := config.LedPresets.StaticColor
			newPreset = append(newPreset, hex)

			config.LedActive.Color = []string{hex}
			config.LedPresets.StaticColor = newPreset

		} else if mode == "BreathingColor" {
			if sliceutil.Contains(config.LedPresets.BreathingColor, hex) {
				c.JSON(http.StatusBadRequest, gin.H{"details": "That color has been already added to this mode."})
				return
			}

			newPreset := config.LedPresets.BreathingColor
			newPreset = append(newPreset, hex)

			config.LedActive.Color = []string{hex}
			config.LedPresets.BreathingColor = newPreset
		} else if mode == "FadingRainbow" {
			c.JSON(http.StatusBadRequest, gin.H{"details": "Unexpected color for this mode."})
			return
		}

		internal.SaveFile(&config.LedActive)
		internal.SaveFile(&config.LedPresets)
		c.Status(http.StatusOK)
	})

	r.GET("/back/getStats", func(c *gin.Context) {
		// c.Status(http.StatusBadGateway)

		c.JSON(http.StatusOK, gin.H{ // TODO : Change to receive real data. Placeholder for now.
			"Temperature": internal.RandomValue(0, 85),
			"CPU":         internal.RandomValue(0, 100),
			"RAM":         internal.RandomValue(0, 100),
			"Disk":        internal.RandomValue(0, 120000), // MB
			"DiskMax":     120000,                          // MB
			"Uptime":      internal.RandomValue(0, 100000),
		})
	})

	r.Run(":80")
}
