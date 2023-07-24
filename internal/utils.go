package internal

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const VERSION_CURRENT = "1.0.0"

var VERSION_LAST string
var VERSION_CHECK int
var VERSION_UPDATED bool

type GithubData struct {
	VersionTag string `json:"tag_name"`
}

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

func CheckLatestVersion() {
	VERSION_CHECK = int(time.Now().Unix())
	response, err := http.Get("https://api.github.com/repos/keelus/horus/releases/latest")
	if err != nil {
		VERSION_LAST = "error"
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		VERSION_LAST = "error"
		return
	}

	data := &GithubData{}

	err = json.Unmarshal(body, data)
	if err != nil {
		VERSION_LAST = "error"
		return
	}

	if data.VersionTag == "" {
		VERSION_LAST = "error"
		return
	}

	VERSION_LAST = strings.Replace(data.VersionTag, "v", "", 1)

	if VERSION_CURRENT == VERSION_LAST || VERSION_CURRENT > VERSION_LAST { // You can't have a greater version, but just in case
		VERSION_UPDATED = true
	} else if VERSION_CURRENT < VERSION_LAST {
		VERSION_UPDATED = false
	}
}
