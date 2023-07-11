package server

import (
	"fmt"
	"horus/config"
	"html/template"
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
	}

	r.AddFromFilesFuncs("login", funcs, "web/templates/login.html")
	r.AddFromFilesFuncs("panel", funcs, "web/templates/panel.html", "web/templates/panels/LedControl.html", "web/templates/panels/SystemStats.html", "web/templates/panels/Settings.html")

	return r
}
