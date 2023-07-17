package led

import (
	"fmt"
	"horus/config"
	"horus/internal"
	"math"
	"os"
	"strconv"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const (
	DefaultBrightness = 90
	NumberOfLeds      = 120
)

type Color struct {
	Red   int
	Green int
	Blue  int
}

var StopRainbow = false
var StopBreathing = false

var LedStrip *ws2811.WS2811

func Init() {
	var err error

	options := ws2811.DefaultOptions
	options.Channels[0].Brightness = config.LedActive.Brightness
	options.Channels[0].LedCount = NumberOfLeds // TODO: Led amount given by client

	LedStrip, err = ws2811.MakeWS2811(&options)
	if err != nil {
		fmt.Printf("failed initializing LED strip: %v\n", err)
	}

	err = LedStrip.Init()
	if err != nil {
		fmt.Printf("failed to initialize LED strip: %v\n", err)
	}

	switch config.LedActive.ActiveMode {
	case "StaticColor":
		Draw()
		break
	case "StaticGradient":
		DrawGradient()
		break
	case "FadingRainbow":
		go Rainbow()
		break
	case "BreathingColor":
		go BreathingColor()
		break
	default:
		fmt.Println("Unexpected LED color mode.")
		os.Exit(-1)
		break
	}
}

func SetColor(color []string) {
	config.LedActive.Color = color
	internal.SaveFile(&config.LedActive)

	Draw()
}

func SetBrightness(brightness int) { // TODO: Will return true once transition is finished to avoid glitching while sending multiple SetBrightness from client
	if config.LedActive.ActiveMode != "StaticColor" {
		config.LedActive.Brightness = brightness
		LedStrip.SetBrightness(0, brightness)
		internal.SaveFile(&config.LedActive)
	} else { // Brightness fading will only occur on Static Color to prevent unexpected flashing
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
}

func Draw() {
	if LedStrip == nil {
		fmt.Printf("LED strip is not initialized\n")
	}

	colorHex, err := strconv.ParseUint(config.LedActive.Color[0], 16, 32) // TODO
	if err != nil {
		fmt.Printf("error parsing color hex code: %v\n", err)
	}

	for i := 0; i < NumberOfLeds; i++ {
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
	for i := 0; i < NumberOfLeds; i++ {
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
	LedStrip.SetBrightness(0, config.LedActive.Brightness)

	for {
		for j := 0; j < 256; j++ {
			for i := 0; i < NumberOfLeds; i++ {
				LedStrip.Leds(0)[i] = wheel((i + j) & 255)
			}

			if config.LedActive.ActiveMode != "FadingRainbow" || StopRainbow { // TODO: Better way
				StopRainbow = false
				return
			}
			LedStrip.Render()
			time.Sleep(time.Duration(config.LedActive.Cooldown) * time.Millisecond)

		}
	}
}

func gaussVal(x int, speed int) int {
	mi := 0.5
	sigma := 0.14
	val := 255 * math.Exp(-((math.Pow((float64(x)/float64(speed))-mi, 2)) / (2 * math.Pow(sigma, 2))))
	return int(val)
}

func BreathingColor() {
	Draw()

	// Normal value := 50 => 500 of smoothness | 10 => 100 (faster)
	//

	speed := config.LedActive.Cooldown * 10
	for {
		for x := 0; x < speed; x++ {
			val := gaussVal(x, speed)

			// fmt.Printf("Bright: %d\n", val)
			ForceDraw(config.LedActive.Color, val)
			LedStrip.SetBrightness(0, val)
			time.Sleep(5 * time.Millisecond)

			if config.LedActive.ActiveMode != "BreathingColor" || StopBreathing {
				StopBreathing = false
				return
			}
		}
	}
}

func DrawGradient() {
	LedStrip.SetBrightness(0, config.LedActive.Brightness)

	if LedStrip == nil {
		fmt.Printf("LED strip is not initialized\n")
	}

	colorGradientIndexes := []Color{}
	for _, colorStr := range config.LedActive.Color {
		red, _ := strconv.ParseInt(colorStr[0:2], 16, 64)
		green, _ := strconv.ParseInt(colorStr[2:4], 16, 64)
		blue, _ := strconv.ParseInt(colorStr[4:6], 16, 64)

		redInt := int(red)
		greenInt := int(green)
		blueInt := int(blue)

		newColor := Color{Red: redInt, Green: greenInt, Blue: blueInt}

		colorGradientIndexes = append(colorGradientIndexes, newColor)
	}
	generatedGradient := generateGradient(colorGradientIndexes)
	for ledIdx, color := range generatedGradient {
		LedStrip.Leds(0)[ledIdx] = uint32(color.Red)<<16 | uint32(color.Green)<<8 | uint32(color.Blue)
	}

	LedStrip.Render()
}

func generateGradient(colors []Color) [NumberOfLeds]Color {
	leds := [NumberOfLeds]Color{}

	increment := float64(NumberOfLeds-1) / float64(len(colors)-1)

	for i := 0; i < len(colors); i++ {
		idx := int(math.Floor(float64(i) * increment))
		leds[idx] = colors[i]
	}

	pairs := len(colors) - 1

	for i := 0; i < pairs; i++ {
		currentColor := colors[i]
		nextColor := colors[i+1]

		currentIndex := int(math.Floor(float64(i) * increment))
		nextIndex := int(math.Floor(float64(i+1) * increment))

		dif := nextIndex - currentIndex

		redIncrement := float64(nextColor.Red-currentColor.Red) / float64(dif)
		greenIncrement := float64(nextColor.Green-currentColor.Green) / float64(dif)
		blueIncrement := float64(nextColor.Blue-currentColor.Blue) / float64(dif)

		for j := currentIndex + 1; j < nextIndex; j++ {
			newRed := int(math.Ceil(float64(currentColor.Red) + float64(j-currentIndex)*redIncrement))
			newGreen := int(math.Ceil(float64(currentColor.Green) + float64(j-currentIndex)*greenIncrement))
			newBlue := int(math.Ceil(float64(currentColor.Blue) + float64(j-currentIndex)*blueIncrement))

			leds[j] = Color{Red: newRed, Green: newGreen, Blue: newBlue}
		}
	}

	return leds
}
