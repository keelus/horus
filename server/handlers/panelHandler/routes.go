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

	category := c.Param("category")
	c.HTML(http.StatusOK, "panel", gin.H{
		"Active":            category,
		"UserConfiguration": config.UserConfiguration,
		"LatestVersion":     latestVersion,
		"UsingLatest":       usingLatest,
		"LedPresets":        config.LedPresets,
		"LedActive":         config.LedActive,
	})
}
