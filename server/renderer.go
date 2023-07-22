package server

import (
	"fmt"
	"horus/config"
	"horus/internal"
	"html/template"
	"math"
	"math/rand"
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
			finalString := ""
			for i, color := range gradientColors {
				finalString += "#" + color
				if i != len(gradientColors)-1 {
					finalString += ","
				}
			}
			return finalString
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
	}

	r.AddFromFilesFuncs("login", funcs, "web/templates/login.html")
	r.AddFromFilesFuncs("panel", funcs, "web/templates/panel.html", "web/templates/panels/LedControl.html", "web/templates/panels/SystemStats.html", "web/templates/panels/Settings.html")
	r.AddFromFilesFuncs("sitewebmanifest", funcs, "web/static/images/icos/site.manifest")

	return r
}
