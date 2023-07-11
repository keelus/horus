package internal

import (
	"encoding/json"
	"fmt"
	"horus/models"
	"io/ioutil"
)

func FileName(fileType interface{}) string {
	var filename string

	switch fileType.(type) {
	case *models.Configuration:
		filename = "userConfig.json"
		break
	case *models.LedActive:
		filename = "ledActive.json"
		break
	case *models.LedPresets:
		filename = "ledPresets.json"
		break
	}

	return filename
}

func LoadFile(loadLocation interface{}) error {
	filename := FileName(loadLocation)

	data, err := ioutil.ReadFile(fmt.Sprintf("config/%s", filename))
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return err
	}

	err = json.Unmarshal(data, loadLocation)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return err
	}

	return nil
}
func SaveFile(saveData interface{}) error {
	filename := FileName(saveData)

	data, err := json.MarshalIndent(saveData, "", "	")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("config/%s", filename), data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return err
	}

	return nil
}
