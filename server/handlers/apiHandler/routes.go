package apiHandler

import (
	"horus/config"
	"horus/internal"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HandleLogin(c *gin.Context) {
	session := sessions.Default(c)

	username := c.PostForm("Username")
	password := c.PostForm("Password")

	// cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"Status": "error",
	// 	})
	// }

	hasError := false

	if config.UserConfiguration.Security.UserInput {
		if username != config.UserConfiguration.UserInfo.Username {
			hasError = true
		}
	}

	err := bcrypt.CompareHashAndPassword([]byte(config.UserConfiguration.UserInfo.Password), []byte(password))
	if err != nil {
		hasError = true
	}

	if hasError {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Status": "error",
		})
	} else {
		session.Clear()
		session.Set("LoggedIn", true)
		session.Save()
		c.JSON(http.StatusOK, gin.H{
			"Status": "OK",
		})
	}
}
func HandleLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
func HandleGetStats(c *gin.Context) {
	// c.Status(http.StatusBadGateway)

	c.JSON(http.StatusOK, gin.H{ // TODO : Change to receive real data. Placeholder for now.
		"Temperature": internal.RandomValue(0, 85),
		"CPU":         internal.RandomValue(0, 100),
		"RAM":         internal.RandomValue(0, 100),
		"Disk":        internal.RandomValue(0, 120000), // MB
		"DiskMax":     120000,                          // MB
		"Uptime":      internal.RandomValue(0, 100000),
	})
}
