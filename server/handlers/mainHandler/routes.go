package mainHandler

import (
	"fmt"
	"horus/config"
	"horus/internal"
	"horus/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var AllowDataRemovedPage = false

func HandleIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/panel/LedControl")
}

func HandleLogin(c *gin.Context) {
	if internal.IsLogged(c) {
		c.Redirect(http.StatusTemporaryRedirect, "/panel/LedControl")
		return
	}
	c.HTML(http.StatusOK, "login", gin.H{
		"CurrentVersion":    internal.VERSION_CURRENT,
		"UserConfiguration": config.UserConfiguration,
	})
	logger.Log(c, logger.VISIT, "Visit to /login.")
}

func SiteManifestHandler(c *gin.Context) {
	themeColor := "181818"
	if !config.UserConfiguration.ColorModeDark {
		themeColor = "F2F2F2"
	}
	manifest := fmt.Sprintf(`{
		"name": "Horus",
		"short_name": "Horus",
		"icons": [
		  {
			"src": "/static/images/icos/android-chrome-192x192.png",
			"sizes": "192x192",
			"type": "image/png"
		  },
		  {
			"src": "/static/images/icos/android-chrome-512x512.png",
			"sizes": "512x512",
			"type": "image/png"
		  }
		],
		"theme_color": "#%s",
		"background_color": "#%s",
		"display": "standalone"
	  }`, themeColor, themeColor)
	c.Header("Content-Type", "application/json")
	c.Writer.WriteString(manifest)
}

func DataRemovedHandler(c *gin.Context) {
	if !AllowDataRemovedPage {
		c.Status(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "dataRemoved", gin.H{})
}
