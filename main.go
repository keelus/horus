package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/forestgiant/sliceutil"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

//
//	TODO:
//	# LED MANAGEMENT
//		Led management[save on a JSON and load in memory]
//		Struct for types of LED [static, rainbow, gradient]
//		For -static: just hex
//		For -rainbow: if all-> just save mode. if specific-> save color spectrum
//		For -gradient: if all-> just save mode. if specific-> save color spectrum
//
//		{
//				"mode":"X",
//			"spectrum": ["","",...] // if all or static: "spectrum":[]
//		}
//
//		Option to save presets

//	# DESIGN & PERSONALISATION
//		Load and save system. Compatibility with scss/css
//

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

var userConfiguration Configuration

func renderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	funcs := template.FuncMap{
		"fPreventCaching": func() string {
			rand.Seed(time.Now().UnixNano())
			characters := "0123456789abcdefghijklmnopqrstuvwxyz"

			randomBytes := make([]byte, 5)
			for i := range randomBytes {
				randomBytes[i] = characters[rand.Intn(len(characters))]
			}

			return fmt.Sprintf("?preventCaching=%s", string(randomBytes))
		},
	}

	r.AddFromFilesFuncs("login", funcs, "templates/login.html")
	r.AddFromFilesFuncs("panel", funcs, "templates/panel.html")

	return r
}

func init() {
	var err error
	userConfiguration, err = loadUserConfiguration()
	if err != nil {
		fmt.Println("User configuration could not be loaded...")
	}
}
func main() {
	r := gin.Default()

	cookieAge := 0

	if userConfiguration.SessionSettings.Lifespan == -1 || userConfiguration.SessionSettings.Unit == "" {
		cookieAge = 0
	} else {
		if userConfiguration.SessionSettings.Unit == "min" {
			cookieAge = 60 * userConfiguration.SessionSettings.Lifespan
		} else if userConfiguration.SessionSettings.Unit == "hour" {
			cookieAge = 60 * 60 * userConfiguration.SessionSettings.Lifespan
		} else if userConfiguration.SessionSettings.Unit == "day" {
			cookieAge = 60 * 60 * 24 * userConfiguration.SessionSettings.Lifespan
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

	r.HTMLRender = renderer()
	r.Use(gin.Recovery())

	r.Static("/static", "./static")

	// ##### FRONT PAGES #####
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/panel/main")
	})
	r.GET("/panel", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/panel/main")
	})
	r.GET("/login", func(c *gin.Context) {
		if isLogged(c) {
			c.Redirect(http.StatusTemporaryRedirect, "/panel/main")
			return
		}
		c.HTML(http.StatusOK, "login", gin.H{
			"UserConfiguration": userConfiguration,
		})
	})
	r.GET("/panel/:category", func(c *gin.Context) {
		if !isLogged(c) {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		latestVersion := "0.5.0" // TODO
		usingLatest := true

		if latestVersion > userConfiguration.Version {
			usingLatest = false
		}

		category := c.Param("category")
		c.HTML(http.StatusOK, "panel", gin.H{
			"Active":            category,
			"UserConfiguration": userConfiguration,
			"LatestVersion":     latestVersion,
			"UsingLatest":       usingLatest,
			"CurrentLED":        "#00FF00",
		})
	})

	r.POST("/back/login", func(c *gin.Context) {
		session := sessions.Default(c)

		username := c.PostForm("Username")
		password := c.PostForm("Password")

		hasError := false

		if userConfiguration.Security.UserInput {
			if username != userConfiguration.UserInfo.Username {
				hasError = true
			}
		}

		if password != userConfiguration.UserInfo.Password {
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
				userConfiguration.UserInfo.Username = username
				if strings.TrimSpace(password) != "" {
					userConfiguration.UserInfo.Password = password
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
					userConfiguration.SessionSettings.Unit = ""
					userConfiguration.SessionSettings.Lifespan = -1
				} else {
					userConfiguration.SessionSettings.Unit = unit
					userConfiguration.SessionSettings.Lifespan = lifespan
				}
			}
			break
		case "LedControl":
			ledControlString := c.PostForm("LedControl")

			var ledControl [2]bool
			err := json.Unmarshal([]byte(ledControlString), &ledControl)
			if err != nil {
				returnError = append(returnError, map[string]string{"details": "Unexpected format.", "field": ""})
			}

			if len(returnError) == 0 {
				userConfiguration.LedControl = ledControl
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
				userConfiguration.SystemStats = systemStats
			}
			break
		case "Logging":
			loggingString := c.PostForm("Logging")

			logging, err := strconv.ParseBool(loggingString)
			if err != nil {
				returnError = append(returnError, map[string]string{"details": "Unexpected format.", "field": ""})
			}

			if len(returnError) == 0 {
				userConfiguration.Logging = logging
			}
			break
		case "Security":
			userInputString := c.PostForm("UserInput")

			userInput, err := strconv.ParseBool(userInputString)
			if err != nil {
				returnError = append(returnError, map[string]string{"details": "Unexpected format.", "field": ""})
			}

			if len(returnError) == 0 {
				userConfiguration.Security.UserInput = userInput
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
				userConfiguration.Units.TimeMode12 = timeMode == "12"
				userConfiguration.Units.TemperatureC = temperature == "C"
			}
			break
		case "Design":
			// TODO
			break
		}

		if len(returnError) == 0 {
			err := saveUserConfiguration()
			if err != nil {
				c.JSON(http.StatusInternalServerError, map[string]string{"details": "Error while saving the user configuration. 500", "field": ""})
			} else {
				c.Status(http.StatusOK)
			}
			return
		}

		c.JSON(http.StatusBadRequest, returnError)
	})

	r.Run(":80")
}

func isLogged(c *gin.Context) bool {
	session := sessions.Default(c)

	status := session.Get("LoggedIn")
	if status == nil {
		return false
	}
	if status.(bool) {
		return true
	}
	return false
}

func loadUserConfiguration() (Configuration, error) {
	var config Configuration

	data, err := ioutil.ReadFile("userConfig.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return config, err
	}

	return config, nil
}

func saveUserConfiguration() error {
	data, err := json.MarshalIndent(userConfiguration, "", "	")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	err = ioutil.WriteFile("userConfig.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return err
	}

	return nil
}
