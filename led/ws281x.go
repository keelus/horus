package led

import (
	"fmt"
	"strconv"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	DefaultBrightness = 90
	DefaultLedCount   = 120
)

var LedStrip *ws2811.WS2811

func Init() {
	options := ws2811.DefaultOptions
	options.Channels[0].Brightness = DefaultBrightness
	options.Channels[0].LedCount = DefaultLedCount

	LedStrip, err := ws2811.MakeWS2811(&options)
	if err != nil {
		fmt.Printf("failed initializing LED strip: %v\n", err)
	}

	err = LedStrip.Init()
	if err != nil {
		fmt.Printf("failed to initialize LED strip: %v\n", err)
	}

	defer LedStrip.Fini()
}

func SetColor(color string) {
	if LedStrip == nil {
		fmt.Printf("LED strip is not initialized\n")
	}

	colorHex, err := strconv.ParseUint(color, 16, 32)
	if err != nil {
		fmt.Printf("error parsing color hex code: %v\n", err)
	}

	for i := 0; i < DefaultLedCount; i++ {
		LedStrip.Leds(0)[i] = uint32(colorHex)
	}

	LedStrip.Render()
}
