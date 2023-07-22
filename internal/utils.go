package internal

import (
	"math/rand"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const CUR_VERSION = "0.9.5 beta"

func IsLogged(c *gin.Context) bool {
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

func RandomValue(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// TODO: Optimisation
func GetGradientStr(gradientColors []string) string {
	gradientStr := ""
	last := len(gradientColors) - 1
	for i, hexStr := range gradientColors {
		gradientStr += "#" + hexStr
		if i != last {
			gradientStr += ","
		}
	}

	return gradientStr
}

func GradientExists(gradientToCheck string, gradientsSlice [][]string) bool {
	for _, gradient := range gradientsSlice {
		if GetGradientStr(gradient) == gradientToCheck {
			return true
		}
	}
	return false
}
