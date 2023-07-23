package panelHandler

import (
	"horus/config"
	"horus/internal"
	"net/http"

	"github.com/forestgiant/sliceutil"
	"github.com/gin-gonic/gin"
)

func HandleIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/panel/ledControl")
}

func HandleCategory(c *gin.Context) {
	if !internal.IsLogged(c) {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	allowedCategories := []string{"LedControl", "SystemStats", "Settings"}
	category := c.Param("category")

	if !sliceutil.Contains(allowedCategories, category) {
		c.Status(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "panel", gin.H{
		"CurrentVersion":     internal.VERSION_CURRENT,
		"Active":             category,
		"UserConfiguration":  config.UserConfiguration,
		"LatestVersion":      internal.VERSION_LAST,
		"UsingLatestVersion": internal.VERSION_UPDATED,
		"LedPresets":         config.LedPresets,
		"LedActive":          config.LedActive,
		"LastChecked":        internal.VERSION_CHECK,
	})
}
