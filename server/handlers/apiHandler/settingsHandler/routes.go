package settingsHandler

import (
	"encoding/json"
	"fmt"
	"horus/config"
	"horus/internal"
	"net/http"
	"strconv"
	"strings"

	"github.com/forestgiant/sliceutil"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SaveConfiguration(c *gin.Context) {
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
}
