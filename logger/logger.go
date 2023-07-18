package logger

import (
	"fmt"
	"horus/config"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Category string

const (
	POWERON Category = "POWERON"
	UP      Category = "UP"
	VISIT   Category = "VISIT" // Only when visiting /login
	LOGIN   Category = "LOGIN"
	LOGOUT  Category = "LOGOUT"
	ERROR   Category = "ERROR"
	LED     Category = "LED"
	SETTING Category = "SETTING"
)

var CurrentLogger *log.Logger

func Init() {
	if config.UserConfiguration.Logging {

		fileName := time.Now().Format("logs/log_02-01-2006_15-04-05.log")
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			defer f.Close()
			os.Exit(-1)
		}

		CurrentLogger = log.New(f, "", log.LstdFlags)

		Log(nil, POWERON, "Horus has been powered on..")
	}
}

func Log(c *gin.Context, category Category, details string) {
	if !config.UserConfiguration.Logging {
		return
	}

	if CurrentLogger == nil { // When user sets logging on while it was off
		Init()
	}

	ip := "::1"
	if c != nil {
		ip = c.ClientIP()
	}

	logLine := fmt.Sprintf("%s\t%s\t\t%s", ip, category, details)
	CurrentLogger.Println(logLine)

}
