package panelHandler

import (
	"horus/config"
	"horus/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/panel/main")
}

func HandleCategory(c *gin.Context) {
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
}
