package clientconfig

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Name             string
	KeyLocation      string
	ServerAddress    string
	Version          string
	AutoUpdate       bool
	AllowInsecureSSL bool
}

func FirstRun() error {
	config := Config{
		Name:             "NOT_SET",
		KeyLocation:      "/etc/vm-manager-client/keys/",
		ServerAddress:    "localhost:8080",
		AutoUpdate:       true,
		AllowInsecureSSL: false,
	}

	// write json to file
	file, err := os.Create("config.json")
	if err != nil {
		return fmt.Errorf("could not create config file, %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("could not encode config file, %s", err)
	}
	return nil
}

func ReadConfig() (Config, bool, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return Config{}, false, fmt.Errorf("could not open config file, %s", err)
	}
	defer file.Close()

	config := Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, true, fmt.Errorf("could not decode config file, %s", err)
	}
	return config, true, nil
}

func UpdateVersion(version string) error {
	file, err := os.Open("config.json")
	if err != nil {
		return fmt.Errorf("could not open config file, %s", err)
	}
	defer file.Close()

	config := Config{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return fmt.Errorf("could not decode config file, %s", err)
	}
	config.Version = version

	// write json to file
	file, err = os.Create("config.json")
	if err != nil {
		return fmt.Errorf("could not create config file, %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(config)
	if err != nil {
		return fmt.Errorf("could not encode config file, %s", err)
	}
	return nil
}
