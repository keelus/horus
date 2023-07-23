package mainHandler

import (
	"horus/config"
	"horus/internal"
	"horus/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var AllowDataRemovedPage = false

func HandleIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/panel/main")
}

func HandleLogin(c *gin.Context) {
	if internal.IsLogged(c) {
		c.Redirect(http.StatusTemporaryRedirect, "/panel/main")
		return
	}
	c.HTML(http.StatusOK, "login", gin.H{
		"CurrentVersion":    internal.VERSION_CURRENT,
		"UserConfiguration": config.UserConfiguration,
	})
	logger.Log(c, logger.VISIT, "Visit to /login.")
}

func SiteManifestHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "sitewebmanifest", gin.H{})
}

func DataRemovedHandler(c *gin.Context) {
	if !AllowDataRemovedPage {
		c.Status(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "dataRemoved", gin.H{})
}
