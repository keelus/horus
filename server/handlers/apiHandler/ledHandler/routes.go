package ledHandler

import (
	"encoding/json"
	"fmt"
	"horus/config"
	"horus/internal"
	"horus/led"
	"horus/logger"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/forestgiant/sliceutil"
	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	if !internal.IsLogged(c) {
		c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
		return
	}

	mode := c.Param("mode")

	if mode == "StaticGradient" {
		var editingGradient bool
		var hexValues []string
		var previousHexValues []string
		hexValuesStr := strings.ToUpper(c.PostForm("hexValues"))
		previousHexValuesStr := strings.ToUpper(c.PostForm("previousHexValues"))
		fmt.Println(hexValuesStr)
		fmt.Println(previousHexValuesStr)

		editingGradient = false
		if previousHexValuesStr != "[]" {
			editingGradient = true
		}

		err := json.Unmarshal([]byte(hexValuesStr), &hexValues)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"details": "Error with the gradient JSON format."})
			return
		}

		rawGradient := internal.GetGradientStr(hexValues)
		if editingGradient {
			err := json.Unmarshal([]byte(previousHexValuesStr), &previousHexValues)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"details": "Error with the old gradient JSON format."})
				return
			}
			// Check if the gradient we want to edit actually exists
			rawPreviousGradient := internal.GetGradientStr(previousHexValues)
			if !internal.GradientExists(rawPreviousGradient, config.LedPresets.StaticGradient) {
				c.JSON(http.StatusBadRequest, gin.H{"details": "That gradient you want to edit doesn't exist."})
				return
			}

			// Check if the gradient we want to edit to actually exists
			if internal.GradientExists(rawGradient, config.LedPresets.StaticGradient) {
				c.JSON(http.StatusBadRequest, gin.H{"details": "That gradient you want to edit to already exists."})
				return
			}

			for i := 0; i < len(config.LedPresets.StaticGradient); i++ {
				raw := internal.GetGradientStr(config.LedPresets.StaticGradient[i])
				if raw == rawPreviousGradient {
					config.LedPresets.StaticGradient[i] = hexValues
					config.LedActive.Color = hexValues
					led.DrawGradient()
				}
			}

		} else {
			if internal.GradientExists(rawGradient, config.LedPresets.StaticGradient) {
				c.JSON(http.StatusBadRequest, gin.H{"details": "That exact gradient already exists."})
				return
			}

			config.LedPresets.StaticGradient = append(config.LedPresets.StaticGradient, hexValues)
			config.LedActive.Color = hexValues
			led.DrawGradient()
		}

	} else {
		// TODO: Better overall code
		hex := strings.ToUpper(c.PostForm("hexValue"))
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
			led.Draw()

		} else if mode == "BreathingColor" {
			if sliceutil.Contains(config.LedPresets.BreathingColor, hex) {
				c.JSON(http.StatusBadRequest, gin.H{"details": "That color has been already added to this mode."})
				return
			}

			newPreset := config.LedPresets.BreathingColor.Colors
			newPreset = append(newPreset, hex)

			config.LedActive.Color = []string{hex}
			config.LedPresets.BreathingColor.Colors = newPreset
		}
	}

	logger.Log(c, logger.LED, fmt.Sprintf("Led hex/gradient added to mode='%s'", mode))
	internal.SaveFile(&config.LedPresets)
	internal.SaveFile(&config.LedActive) // TODO activate new gradient
	c.Status(http.StatusOK)

}
func Delete(c *gin.Context) {
	if !internal.IsLogged(c) {
		c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
		return
	}

	mode := c.Param("mode")

	if mode == "StaticGradient" {
		rawGradient := c.PostForm("rawGradient")

		existingGradients := len(config.LedPresets.StaticGradient)
		if existingGradients == 1 {
			c.JSON(http.StatusBadRequest, gin.H{"details": "There must be at least 1 preset gradient."})
			return
		}
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
		config.LedActive.Color = config.LedPresets.StaticGradient[0]
		led.DrawGradient()
	} else { // TODO: Better overall code
		hex := c.PostForm("hexValue")

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
			if len(config.LedPresets.BreathingColor.Colors) == 1 {
				c.JSON(http.StatusBadRequest, gin.H{"details": "There must be at least 1 preset color."})
				return
			}
			newPreset := []string{}

			for _, color := range config.LedPresets.BreathingColor.Colors {
				if color != hex {
					newPreset = append(newPreset, color)
				}
			}

			if config.LedActive.Color[0] == hex {
				config.LedActive.Color[0] = newPreset[0]
			}

			config.LedPresets.BreathingColor.Colors = newPreset
		}
	}

	logger.Log(c, logger.LED, fmt.Sprintf("Led hex/gradient deleted from mode='%s'", mode))
	internal.SaveFile(&config.LedActive)
	internal.SaveFile(&config.LedPresets)
	c.Status(http.StatusOK)
}

func Activate(c *gin.Context) {
	if !internal.IsLogged(c) {
		c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
		return
	}
	mode := c.Param("mode")

	if mode == "StaticGradient" {
		rawGradient := c.PostForm("rawGradient")

		if rawGradient == "" {
			config.LedActive.Color = config.LedPresets.StaticGradient[0] // All colors are activated on Fading Colors
		} else {
			if !internal.GradientExists(rawGradient, config.LedPresets.StaticGradient) {
				c.JSON(http.StatusBadRequest, gin.H{"details": "That gradient doesn't exist."})
				return
			}
			for _, gradient := range config.LedPresets.StaticGradient {
				if internal.GetGradientStr(gradient) == rawGradient {
					config.LedActive.Color = gradient
				}
			}
		}

		config.LedActive.ActiveMode = "StaticGradient"
		config.LedActive.Cooldown = 0 // TODO

		led.DrawGradient()
	} else {
		hex := c.PostForm("hexValue")

		if mode == "StaticColor" {
			// By default first color is activated. Always will be one at least.
			if hex == "" {
				led.SetColor([]string{config.LedPresets.StaticColor[0]}) // By default first color will be initialized
			} else {
				led.SetColor([]string{hex})
			}

			config.LedActive.ActiveMode = "StaticColor"
			config.LedActive.Cooldown = 0
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
			config.LedActive.Cooldown = config.LedPresets.BreathingColor.Cooldown

			if hex == "" {
				led.SetColor([]string{config.LedPresets.BreathingColor.Colors[0]}) // By default first color will be initialized
			} else {
				led.SetColor([]string{hex})
			}

			if previousMode != "BreathingColor" {
				if previousMode == "FadingRainbow" {
					led.StopRainbow = true
				}
				go led.BreathingColor()
			}
		}

	}
	internal.SaveFile(&config.LedActive)
	logger.Log(c, logger.LED, fmt.Sprintf("Led hex/gradient activated. Mode='%s'", mode))
	c.Status(http.StatusOK)
}

func SetBrightness(c *gin.Context) {
	if !internal.IsLogged(c) {
		c.JSON(http.StatusForbidden, gin.H{"details": "User not logged in."})
		return
	}

	if led.ApplyingBrightness {
		convertedBrightness := int(math.Ceil(float64(config.LedActive.Brightness) * 100 / 255)) // We return brightness 100% to set the slider value again.
		c.JSON(http.StatusBadGateway, gin.H{"details": "Brightness is being applied, please wait.", "brightness": convertedBrightness})
		return
	}

	valuePercent, err := strconv.Atoi(c.Param("valuePercent"))
	if err != nil || !(valuePercent >= 0 && valuePercent <= 100) {
		c.JSON(http.StatusBadGateway, gin.H{"details": "Unexpected brightness value type."})
		return
	}

	newVal := int(valuePercent * 255 / 100)
	led.SetBrightness(newVal)

	internal.SaveFile(&config.LedActive)
	logger.Log(c, logger.LED, fmt.Sprintf("Led brightness set to '%d'", newVal))
	c.JSON(http.StatusOK, gin.H{"Color": config.LedActive.Color, "Brightness": config.LedActive.Brightness, "Cooldown": config.LedActive.Cooldown})
	return
}

func SetCooldown(c *gin.Context) {
	amountStr := c.Param("amount")
	mode := c.Param("mode")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": "Unexpected type of amount."})
		return
	}

	if mode == "BreathingColor" {
		config.LedPresets.BreathingColor.Cooldown = amount
		if config.LedActive.ActiveMode != "BreathingColor" {
			go led.BreathingColor()
		}

	} else if mode == "FadingRainbow" {
		config.LedPresets.FadingRainbow = amount
		if config.LedActive.ActiveMode != "FadingRainbow" {
			go led.Rainbow()
		}
	}

	config.LedActive.Cooldown = amount
	logger.Log(c, logger.LED, fmt.Sprintf("Led cooldown set to '%d' in mode='%s'", amount, mode))
	internal.SaveFile(&config.LedActive)
	internal.SaveFile(&config.LedPresets)
	c.Status(http.StatusOK)
}

func SetLedAmount(c *gin.Context) {
	ledAmountStr := c.Param("amount")

	ledAmount, err := strconv.Atoi(ledAmountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"details": "Unexpected led amount type."})
		return
	}

	if ledAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"details": "Led amount should be a positive amount."})
		return
	}

	config.LedActive.LedAmount = ledAmount
	internal.SaveFile(&config.LedActive)
	led.Init()
	c.Status(http.StatusOK)
}
