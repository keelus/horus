package main

import (
	"horus/config"
	"horus/server"
)

//
//	TODO:
//	# USER CONFIGURATION
//		Switch from .JSON to .YAML
//	# DESIGN & PERSONALISATION
//		Load and save system. Compatibility with scss/css
//

const CUR_VERSION = "0.6.6"

func main() {
	config.Init()
	server.Init()
}
