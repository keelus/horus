package led

import (
	"fmt"
	"horus/config"
	"horus/internal"
	"strconv"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	DefaultBrightness = 90
	DefaultLedCount   = 120
)

var LedStrip *ws2811.WS2811

func Init() {
	var err error

	options := ws2811.DefaultOptions
	options.Channels[0].Brightness = config.LedActive.Brightness
	options.Channels[0].LedCount = DefaultLedCount

	LedStrip, err = ws2811.MakeWS2811(&options)
	if err != nil {
		fmt.Printf("failed initializing LED strip: %v\n", err)
	}

	err = LedStrip.Init()
	if err != nil {
		fmt.Printf("failed to initialize LED strip: %v\n", err)
	}
}

func SetColor(color []string) {
	config.LedActive.Color = color
	internal.SaveFile(&config.LedActive)

	Draw()
}

func SetBrightness(brightness int) {
	currentBrightness := config.LedActive.Brightness
	duration := 1 * time.Second

	steps := 1
	stepSize := 1
	if brightness > currentBrightness {
		steps = brightness - currentBrightness
		stepSize = 1
	} else {
		steps = currentBrightness - brightness
		stepSize = -1
	}

	interval := duration / time.Duration(steps)

	for i := 0; i < steps; i++ {
		newBrightness := currentBrightness + stepSize*i
		ForceDraw(config.LedActive.Color, newBrightness)
		time.Sleep(interval)
	}

	config.LedActive.Brightness = brightness
	LedStrip.SetBrightness(0, brightness)
	internal.SaveFile(&config.LedActive)
	Draw()
}

func Draw() {
	if LedStrip == nil {
		fmt.Printf("LED strip is not initialized\n")
	}

	colorHex, err := strconv.ParseUint(config.LedActive.Color[0], 16, 32) // TODO
	if err != nil {
		fmt.Printf("error parsing color hex code: %v\n", err)
	}

	for i := 0; i < DefaultLedCount; i++ {
		LedStrip.Leds(0)[i] = uint32(colorHex)
	}

	LedStrip.Render()
}

func ForceDraw(color []string, brightness int) {
	if LedStrip == nil {
		fmt.Printf("LED strip is not initialized\n")
	}

	colorHex, err := strconv.ParseUint(config.LedActive.Color[0], 16, 32) // TODO
	if err != nil {
		fmt.Printf("error parsing color hex code: %v\n", err)
	}

	LedStrip.SetBrightness(0, brightness)
	for i := 0; i < DefaultLedCount; i++ {
		LedStrip.Leds(0)[i] = uint32(colorHex)
	}

	LedStrip.Render()
}
