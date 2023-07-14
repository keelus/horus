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

var StopRainbow = false
var StopPulsating = false

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

	if config.LedActive.ActiveMode == "StaticColor" {
		Draw()
	} else {
		if config.LedActive.ActiveMode == "FadingRainbow" {
			go Rainbow()
		}
	}
}

func SetColor(color []string) {
	config.LedActive.Color = color
	internal.SaveFile(&config.LedActive)

	Draw()
}

func SetBrightness(brightness int) { // TODO: Will return true once transition is finished to avoid glitching while sending multiple SetBrightness from client

	if config.LedActive.ActiveMode == "StaticColor" { // Brightness fading will only occur on Static Color to prevent unexpected flashing
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
	} else {
		config.LedActive.Brightness = brightness
		LedStrip.SetBrightness(0, brightness)
		internal.SaveFile(&config.LedActive)
	}
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

// Code translated to Golang from Python [rpi-ws281x-python example.py]

func wheel(pos int) uint32 {
	if pos < 85 {
		return uint32(pos*3)<<16 | uint32(255-pos*3)<<8
	} else if pos < 170 {
		pos -= 85
		return uint32(255-pos*3)<<16 | uint32(pos*3)
	} else {
		pos -= 170
		return uint32(pos*3)<<8 | uint32(255-pos*3)
	}
}

func Rainbow() {
	fmt.Println("Rainbow loop start")
	for j := 0; j < 256; j++ {
		for i := 0; i < 120; i++ {
			LedStrip.Leds(0)[i] = wheel((i + j) & 255)
		}

		if config.LedActive.ActiveMode != "FadingRainbow" || StopRainbow { // TODO: Better way
			StopRainbow = false
			return
		}
		LedStrip.Render()
		time.Sleep(time.Duration(config.LedActive.Cooldown) * time.Millisecond)

	}
	if config.LedActive.ActiveMode == "FadingRainbow" { // TODO: Better way
		Rainbow()
	}

}

func PulsatingColor() {
	currentBrightness := config.LedActive.Brightness
	Draw()

	// Down
	for i := currentBrightness; i > 0; i-- {
		LedStrip.SetBrightness(0, i)
		fmt.Printf("Brightness: %d\n", i)
		ForceDraw(config.LedActive.Color, i)
		time.Sleep(time.Duration(config.LedActive.Cooldown) * time.Millisecond)
		if config.LedActive.ActiveMode != "PulsatingColor" || StopPulsating { // TODO: Better way
			StopPulsating = false
			return
		}
	}

	// Up
	for i := 0; i < currentBrightness; i++ {
		LedStrip.SetBrightness(0, i)
		fmt.Printf("Brightness: %d\n", i)
		ForceDraw(config.LedActive.Color, i)
		time.Sleep(time.Duration(config.LedActive.Cooldown) * time.Millisecond)
		if config.LedActive.ActiveMode != "PulsatingColor" || StopPulsating { // TODO: Better way
			StopPulsating = false
			return
		}
	}

	if config.LedActive.ActiveMode != "PulsatingColor" || StopPulsating { // TODO: Better way
		StopPulsating = false
		return
	}

	PulsatingColor()
}
