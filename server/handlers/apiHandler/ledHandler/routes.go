package ledHandler

import (
	"encoding/json"
	"fmt"
	"horus/config"
	"horus/internal"
	"horus/led"
	"net/http"
	"strconv"
	"strings"

	"github.com/forestgiant/sliceutil"
	"github.com/gin-gonic/gin"
)

func SetCooldown(c *gin.Context) {
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
}
func Activate(c *gin.Context) {
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
}
func ActivateGradient(c *gin.Context) {
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

}
func ActivateHex(c *gin.Context) {
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
}
func SetBrightness(c *gin.Context) {
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
}
func DeleteGradient(c *gin.Context) {
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
}
func Delete(c *gin.Context) {
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
}
func AddGradient(c *gin.Context) {
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
}
func AddHex(c *gin.Context) {
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
}
