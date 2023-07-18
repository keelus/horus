package internal

import (
	"fmt"
	"horus/models"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

func FileName(fileType interface{}) string {
	var filename string

	switch fileType.(type) {
	case *models.Configuration:
		filename = "userConfig.yaml"
		break
	case *models.LedActive:
		filename = "ledActive.yaml"
		break
	case *models.LedPresets:
		filename = "ledPresets.yaml"
		break
	case models.Configuration:
		filename = "userConfig.yaml"
		break
	case models.LedActive:
		filename = "ledActive.yaml"
		break
	case models.LedPresets:
		filename = "ledPresets.yaml"
		break
	}

	return filename
}

func LoadFile(loadLocation interface{}, fromDefaults bool) error {
	cError := color.New(color.FgRed, color.Bold)

	var location string
	filename := FileName(loadLocation)

	if fromDefaults {
		location = fmt.Sprintf("config/defaults/%s", filename)
	} else {
		location = fmt.Sprintf("config/%s", filename)
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s", location))
	if err != nil {
		if fromDefaults {
			cError.Printf("Error loading a default configuration file. '%s'\n", location)
			os.Exit(-1)
		}
		return err
	}

	_ = yaml.Unmarshal(data, loadLocation)
	if err != nil {
		if fromDefaults {
			cError.Printf("Error parsing a default configuration file. '%s'\n", location)
			os.Exit(-1)
		}
		cError.Printf("Error parsing a configuration file. '%s'\n", location)
		return err
	}

	return nil
}

func SaveFile(saveData interface{}) error {
	cError := color.New(color.FgRed, color.Bold)
	filename := FileName(saveData)

	data, err := yaml.Marshal(saveData)
	if err != nil {
		cError.Printf("Error parsing a configuration file. '%s'\n", filename)
		os.Exit(-1)
	}

	err = ioutil.WriteFile(fmt.Sprintf("config/%s", filename), data, 0644)
	if err != nil {
		cError.Printf("Error saving a configuration file. '%s'\n", filename)
		os.Exit(-1)
	}

	return nil
}
