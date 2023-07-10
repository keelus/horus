package internal

import (
	"encoding/json"
	"fmt"
	"horus/models"
	"io/ioutil"
)

func LoadUserConfiguration() (models.Configuration, error) {
	var config models.Configuration

	data, err := ioutil.ReadFile("internal/data/userConfig.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return config, err
	}

	return config, nil
}

func SaveUserConfiguration(userConfiguration models.Configuration) error {
	data, err := json.MarshalIndent(userConfiguration, "", "	")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	err = ioutil.WriteFile("internal/data/userConfig.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return err
	}

	return nil
}

func LoadLedPresets() (models.LedPresets, error) {
	var presets models.LedPresets

	data, err := ioutil.ReadFile("internal/data/ledPresets.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return presets, err
	}

	err = json.Unmarshal(data, &presets)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return presets, err
	}

	return presets, nil
}

func SaveLedPresets(ledPresets models.LedPresets) error {
	data, err := json.MarshalIndent(ledPresets, "", "	")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	err = ioutil.WriteFile("internal/data/ledPresets.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return err
	}

	return nil
}

func LoadLedActive() (models.LedActive, error) {
	var active models.LedActive

	data, err := ioutil.ReadFile("internal/data/ledActive.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return active, err
	}

	err = json.Unmarshal(data, &active)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return active, err
	}

	return active, nil
}

func SaveLedActive(ledActive models.LedActive) error {
	data, err := json.MarshalIndent(ledActive, "", "	")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	err = ioutil.WriteFile("internal/data/ledActive.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return err
	}

	return nil
}
