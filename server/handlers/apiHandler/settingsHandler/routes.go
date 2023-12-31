package settingsHandler

import (
	"encoding/json"
	"fmt"
	"horus/config"
	"horus/internal"
	"horus/logger"
	"horus/server/handlers/mainHandler"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/forestgiant/sliceutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SaveConfiguration(c *gin.Context) {
	returnError := []map[string]string{}

	category := c.Param("category")

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
	case "ColorTheme":
		colorTheme := c.PostForm("ColorTheme")
		switch colorTheme {
		case "light":
			config.UserConfiguration.ColorModeDark = false
		case "dark":
			config.UserConfiguration.ColorModeDark = true
		default:
			returnError = append(returnError, map[string]string{"details": "Unexpected color mode.", "field": "null"})
		}
	case "SessionSettings":
		lifespanString := c.PostForm("Lifespan")
		unit := c.PostForm("Unit")

		allowedUnits := []string{"min", "hour", "day", ""}

		if !sliceutil.Contains(allowedUnits, unit) {
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
		temperature := c.PostForm("Temperature")

		if temperature != "C" && temperature != "F" {
			returnError = append(returnError, map[string]string{"details": "Unexpected temperature unit.", "field": "units1"})
		}

		if len(returnError) == 0 {
			config.UserConfiguration.Units.TemperatureC = temperature == "C"
		}
		break
	case "Deletion":
		password := c.PostForm("Password")

		err := bcrypt.CompareHashAndPassword([]byte(config.UserConfiguration.UserInfo.Password), []byte(password))
		if err != nil {
			returnError = append(returnError, map[string]string{"details": "Incorrect password.", "field": "deletion0"})
		}

		if len(returnError) == 0 {
			_ = os.Remove("config/ledActive.yaml")
			_ = os.Remove("config/ledPresets.yaml")
			_ = os.Remove("config/userConfig.yaml")
			_ = os.Remove("avatar.jpg")

			logger.Log(c, logger.SETTING, "Horus data deleted successfully.")

			mainHandler.AllowDataRemovedPage = true

			go Shutdown(2)

			c.JSON(http.StatusOK, gin.H{
				"Status": "OK",
			})
			return
		}
		break
	}

	if len(returnError) == 0 {
		err := internal.SaveFile(&config.UserConfiguration)
		if err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{"details": "Error while saving the user configuration. 500", "field": ""})
		} else {
			logger.Log(c, logger.SETTING, fmt.Sprintf("Settings changed [category='%s'].", category)) // TODO: Specify
			c.Status(http.StatusOK)
		}
		return
	}

	c.JSON(http.StatusBadRequest, returnError)
}

func Shutdown(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
	os.Exit(0)
}
