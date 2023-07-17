package mainHandler

import (
	"horus/config"
	"horus/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/panel/main")
}

func HandleLogin(c *gin.Context) {
	if internal.IsLogged(c) {
		c.Redirect(http.StatusTemporaryRedirect, "/panel/main")
		return
	}
	c.HTML(http.StatusOK, "login", gin.H{
		"UserConfiguration": config.UserConfiguration,
	})
}
