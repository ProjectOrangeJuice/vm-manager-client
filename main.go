package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ProjectOrangeJuice/vm-manager-client/cert"
	"github.com/ProjectOrangeJuice/vm-manager-client/connection"
)

func main() {

	// read file to check if this is the first run

	// if this is the first run, run setup
	config, exists, err := readConfig()
	if err != nil {
		if exists {
			log.Printf("Error reading config file, %s. As the file exists, we won't create it", err)
			return
		}
		log.Printf("Config file was not there, running setup [%s]", err)
		err = firstRun()
		if err != nil {
			log.Printf("First run failed, %s", err)
			return
		}
		return
	}

	log.Printf("Config [%+v]", config)
	TLSConfig, err := cert.SetupTLSConfig(config.KeyLocation, config.Name)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := tls.Dial("tcp", config.ServerAddress, TLSConfig)
		if err != nil {
			log.Printf("Error dialing, trying again in 5 seconds: %s", err)
			time.Sleep(5 * time.Second)
			return
		}

		serverConnection := connection.NewConnection(conn)
		serverConnection.ProcessLines() // Loops forever until disconnected

		log.Printf("Disconnected, trying again in 5 seconds")
		conn.Close()
		time.Sleep(5 * time.Second)

	}

}

type Config struct {
	Name          string
	KeyLocation   string
	ServerAddress string
}

func firstRun() error {
	config := Config{
		Name:          "Test client",
		KeyLocation:   "./keys/",
		ServerAddress: "localhost:8080",
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

func readConfig() (Config, bool, error) {
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
