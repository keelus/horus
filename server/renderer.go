package server

import (
	"fmt"
	"horus/config"
	"horus/internal"
	"html/template"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-contrib/multitemplate"
)

func Renderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	funcs := template.FuncMap{
		"fPreventCaching": func() string {
			rand.Seed(time.Now().UnixNano())
			characters := "0123456789abcdefghijklmnopqrstuvwxyz"

			randomBytes := make([]byte, 5)
			for i := range randomBytes {
				randomBytes[i] = characters[rand.Intn(len(characters))]
			}

			return fmt.Sprintf("?preventCaching=%s", string(randomBytes))
		},
		"showSystemStats": func() bool { // System Stats tab will show only if at least 1 element is shown.
			show := false
			for _, stat := range config.UserConfiguration.SystemStats {
				if stat {
					show = true
				}
			}

			return show
		},
		"gradientString": func(gradientColors []string) string {
			return "#" + strings.Join(gradientColors, ",#")
		},
		"isActiveColor": func(color string) bool {
			activeColor := config.LedActive.Color[0]
			return color == activeColor
		},
		"isActiveGradient": func(gradientColors []string) bool {
			activeGradient := internal.GetGradientStr(config.LedActive.Color)
			return internal.GetGradientStr(gradientColors) == activeGradient
		},
		"isActiveMode": func(mode string) bool {
			activeMode := config.LedActive.ActiveMode
			return mode == activeMode
		},
		"convertBrightness": func(brightness int) int {
			return int(math.Ceil(float64(brightness) * 100 / 255))
		},
		"renderDate": func(unixDateInt int) string {
			var visualAmount int
			var pluralS string

			if unixDateInt == 0 {
				return "never"
			}

			now := time.Now()

			unixDate := time.Unix(int64(unixDateInt), 0)

			diffSeconds := int(math.Floor(now.Sub(unixDate).Seconds()))

			if diffSeconds < 60 {
				visualAmount = diffSeconds
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d second%s ago", visualAmount, pluralS)
			} else if diffSeconds < 3600 {
				visualAmount = int(math.Floor(float64(diffSeconds) / float64(60)))
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d minute%s ago", visualAmount, pluralS)
			} else if diffSeconds < 86400 {
				visualAmount = int(math.Floor(float64(diffSeconds) / float64(3600)))
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d hour%s ago", visualAmount, pluralS)
			} else if diffSeconds < 604800 {
				visualAmount = int(math.Floor(float64(diffSeconds) / float64(86400)))
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d day%s ago", visualAmount, pluralS)
			} else if diffSeconds < 2419200 {
				visualAmount = int(math.Floor(float64(diffSeconds) / float64(604800)))
				if visualAmount != 1 {
					pluralS = "s"
				}
				return fmt.Sprintf("%d week%s ago, at %d:%d", visualAmount, pluralS, unixDate.Hour(), unixDate.Minute())
			} else {
				return fmt.Sprintf("%d of %s of %d, at %d:%d", unixDate.Day(), unixDate.Month(), unixDate.Year(), unixDate.Hour(), unixDate.Minute())
			}
		},
	}

	r.AddFromFilesFuncs("login", funcs, "web/templates/login.html")
	r.AddFromFilesFuncs("panel", funcs, "web/templates/panel.html", "web/templates/panels/LedControl.html", "web/templates/panels/SystemStats.html", "web/templates/panels/Settings.html")
	r.AddFromFilesFuncs("dataRemoved", funcs, "web/templates/dataRemoved.html")

	return r
}
