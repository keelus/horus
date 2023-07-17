package server

import (
	"net/http"

	"horus/config"
	"horus/server/handlers/apiHandler"
	"horus/server/handlers/apiHandler/ledHandler"
	"horus/server/handlers/apiHandler/settingsHandler"
	"horus/server/handlers/mainHandler"
	"horus/server/handlers/panelHandler"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	cookieAge := 0

	if config.UserConfiguration.SessionSettings.Lifespan == -1 || config.UserConfiguration.SessionSettings.Unit == "" {
		cookieAge = 0
	} else {
		if config.UserConfiguration.SessionSettings.Unit == "min" {
			cookieAge = 60 * config.UserConfiguration.SessionSettings.Lifespan
		} else if config.UserConfiguration.SessionSettings.Unit == "hour" {
			cookieAge = 60 * 60 * config.UserConfiguration.SessionSettings.Lifespan
		} else if config.UserConfiguration.SessionSettings.Unit == "day" {
			cookieAge = 60 * 60 * 24 * config.UserConfiguration.SessionSettings.Lifespan
		}
	}

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		Path: "/",
		// Secure:   true,
		Secure: false,
		// HttpOnly: true,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   cookieAge,
	})
	r.Use(sessions.Sessions("sesion", store))

	r.HTMLRender = Renderer()
	r.Use(gin.Recovery())

	r.Static("/static", "web/static")

	mainGroup := r.Group("/")
	{
		mainGroup.GET("/", mainHandler.HandleIndex)
		mainGroup.GET("/login", mainHandler.HandleLogin)
	}

	panelGroup := r.Group("/panel")
	{
		panelGroup.GET("/", panelHandler.HandleIndex)
		panelGroup.GET("/:category", panelHandler.HandleCategory)
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.POST("/login", apiHandler.HandleLogin)
		apiGroup.GET("/logout", apiHandler.HandleLogout)
		apiGroup.GET("/getStats", apiHandler.HandleGetStats)

		ledControlGroup := apiGroup.Group("/ledControl")
		{
			ledControlGroup.POST("/cooldown/:mode/:amount", ledHandler.SetCooldown)
			ledControlGroup.POST("/activate/:mode", ledHandler.Activate)
			ledControlGroup.POST("/activateGradient", ledHandler.ActivateGradient)
			ledControlGroup.POST("/activate/:mode/:hex", ledHandler.ActivateHex)
			ledControlGroup.POST("/brightness/:valuePercent", ledHandler.SetBrightness)
			ledControlGroup.POST("/deleteGradient", ledHandler.DeleteGradient)
			ledControlGroup.POST("/delete/:mode/:hex", ledHandler.Delete)
			ledControlGroup.POST("/addGradient/", ledHandler.AddGradient)
			ledControlGroup.POST("/add/:mode/:hex", ledHandler.AddHex)
		}

		settingsGroup := apiGroup.Group("/settings")
		{
			settingsGroup.POST("/saveConfiguration/:category", settingsHandler.SaveConfiguration)
		}
	}

	// ##### FRONT PAGES #####

	r.Run(":80")
	return r
}
