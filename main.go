package main

import (
	"crypto/tls"
	"log"
	"time"

	"github.com/ProjectOrangeJuice/vm-manager-client/cert"
	clientconfig "github.com/ProjectOrangeJuice/vm-manager-client/clientConfig"
	"github.com/ProjectOrangeJuice/vm-manager-client/connection"
	"github.com/ProjectOrangeJuice/vm-manager-client/update"
)

func main() {

	// read file to check if this is the first run

	// if this is the first run, run setup
	config, exists, err := clientconfig.ReadConfig()
	if err != nil {
		if exists {
			log.Printf("Error reading config file, %s. As the file exists, we won't create it", err)
			return
		}
		log.Printf("Config file was not there, running setup [%s]", err)
		err = clientconfig.FirstRun()
		if err != nil {
			log.Printf("First run failed, %s", err)
			return
		}
		return
	}

	log.Printf("Config [%+v]", config)
	ver, err := update.WhatVersionStartup(&config)
	if err != nil {
		log.Printf("Error getting version, %s", err)
		return
	}
	log.Printf("Version [%s]", ver)

	// Make it do the update
	err = update.UpdateIfNeeded(&config)
	if err != nil {
		log.Printf("Error updating, %s", err)
		return
	}

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
