package led

import (
	"fmt"
	"strconv"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	brightness = 90
	ledCounts  = 120
)

var leds *ws2811.WS2811

func Init() {
	var err error

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	leds, err = ws2811.MakeWS2811(&opt)
	if err != nil {
		fmt.Println("Failed initializing ws2811 strip. Err: ", err)
		return
	}

	err = leds.Init()
	if err != nil {
		fmt.Println("Failed on leds.Init(). Err: ", err)
		return
	}
	defer leds.Fini()
}

func SetColor(color string) {
	colorHex, err := strconv.ParseUint(color, 16, 32)
	if err != nil {
		fmt.Println("Error at parsing the color hex code. Setting to #FF0000")
		colorHex = 0xFFFFFF
	}

	for i := 0; i < ledCounts; i++ {
		leds.Leds(0)[i] = uint32(colorHex)
	}
	leds.Render()
}
