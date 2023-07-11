package internal

import (
	"fmt"
	"horus/models"
	"io/ioutil"

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
	}

	return filename
}

func LoadFile(loadLocation interface{}) error {
	filename := FileName(loadLocation)

	data, err := ioutil.ReadFile(fmt.Sprintf("config/%s", filename))
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return err
	}

	err = yaml.Unmarshal(data, loadLocation)
	if err != nil {
		fmt.Println("Error unmarshaling YAML:", err)
		return err
	}

	return nil
}
func SaveFile(saveData interface{}) error {
	filename := FileName(saveData)

	data, err := yaml.Marshal(saveData)
	if err != nil {
		fmt.Println("Error marshaling YAML:", err)
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("config/%s", filename), data, 0644)
	if err != nil {
		fmt.Println("Error writing YAML file:", err)
		return err
	}

	return nil
}
