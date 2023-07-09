package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Configuration struct {
	Version         string
	UserInfo        UserInfo
	SessionSettings SessionSettings
	LedControl      [2]bool
	SystemStats     [5]bool
	Logging         bool
	Security        Security
	Units           Units
	Design          Design
}

type UserInfo struct {
	Username string
	Password string
}

type SessionSettings struct {
	Lifespan int
	Unit     string
}

type Security struct {
	UserInput bool
}

type Units struct {
	TimeMode12   bool
	TemperatureC bool
}

type Design struct {
	Accent []string
	Fonts  []FontDetails
}

type FontDetails struct {
	Title  string
	Source string
}

var userConfiguration Configuration

func renderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	r.AddFromFiles("login", "templates/login.html")
	r.AddFromFiles("panel", "templates/panel.html")

	return r
}

func init() {
	// Read the JSON file
	data, err := ioutil.ReadFile("userConfig.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Unmarshal JSON data into the Configuration struct
	err = json.Unmarshal(data, &userConfiguration)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}
}
func main() {
	r := gin.Default()

	cookieAge := 0

	if userConfiguration.SessionSettings.Lifespan == -1 || userConfiguration.SessionSettings.Unit == "" {
		cookieAge = 0
	} else {
		if userConfiguration.SessionSettings.Unit == "min" {
			cookieAge = 60 * userConfiguration.SessionSettings.Lifespan
		} else if userConfiguration.SessionSettings.Unit == "hour" {
			cookieAge = 60 * 60 * userConfiguration.SessionSettings.Lifespan
		} else if userConfiguration.SessionSettings.Unit == "day" {
			cookieAge = 60 * 60 * 24 * userConfiguration.SessionSettings.Lifespan
		}
	}

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieAge,
	})
	r.Use(sessions.Sessions("sesion", store))

	r.HTMLRender = renderer()
	r.Use(gin.Recovery())

	r.Static("/static", "./static")

	// ##### FRONT PAGES #####
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/panel/main")
	})
	r.GET("/panel", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/panel/main")
	})
	r.GET("/login", func(c *gin.Context) {
		if isLogged(c) {
			c.Redirect(http.StatusTemporaryRedirect, "/panel/main")
			return
		}
		c.HTML(http.StatusOK, "login", gin.H{
			"UserConfiguration": userConfiguration,
		})
	})
	r.GET("/panel/:category", func(c *gin.Context) {
		if !isLogged(c) {
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		latestVersion := "1.3.1"
		usingLatest := true

		if latestVersion > userConfiguration.Version {
			usingLatest = false
		}

		category := c.Param("category")
		c.HTML(http.StatusOK, "panel", gin.H{
			"Active":            category,
			"UserConfiguration": userConfiguration,
			"LatestVersion":     latestVersion,
			"UsingLatest":       usingLatest,
		})
	})

	r.POST("/back/login", func(c *gin.Context) {
		session := sessions.Default(c)

		username := c.PostForm("Username")
		password := c.PostForm("Password")

		hasError := false

		if userConfiguration.Security.UserInput {
			if username != userConfiguration.UserInfo.Username {
				hasError = true
			}
		}

		if password != userConfiguration.UserInfo.Password {
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
	})

	r.GET("/back/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusTemporaryRedirect, "/")
	})

	r.POST("/back/saveConfiguration/:category", func(c *gin.Context) {
		category := c.Param("category")

		switch category {
		case "userInfo":
			break
		case "sessionSettings":
			break
		case "ledControl":
			break
		case "systemStats":
			break
		case "logging":
			break
		case "security":
			break
		case "region":
			break
		case "design":
			break
		case "updates":
			break
		case "deletion":
			break
		}
	})

	r.Run(":80")
}

func isLogged(c *gin.Context) bool {
	session := sessions.Default(c)

	status := session.Get("LoggedIn")
	if status == nil {
		return false
	}
	if status.(bool) {
		return true
	}
	return false
}
