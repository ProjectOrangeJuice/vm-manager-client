package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ProjectOrangeJuice/vm-manager-client/cert"
	clientconfig "github.com/ProjectOrangeJuice/vm-manager-client/clientConfig"
	"github.com/ProjectOrangeJuice/vm-manager-client/connection"
	"github.com/ProjectOrangeJuice/vm-manager-client/update"
)

func main() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		os.Exit(0)
	}()

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
	ver, err := update.WhatVersionStartup(&config)
	if err != nil {
		log.Printf("Error getting version, %s", err)
		return
	}
	log.Printf("Version [%s]", ver)
	log.Printf("Config [%+v]", config)
	if config.AutoUpdate {
		// err := update.FinishUpdate()
		// if err != nil {
		// 	log.Printf("Error finishing update, %s", err)
		// 	return
		// }
		// Make it do the update
		err = update.UpdateIfNeeded(&config)
		if err != nil {
			log.Printf("Error updating, %s", err)
			return
		}
	}

	if config.Name == "NOT_SET" {
		log.Println("The client name is not set. Please set it in the config file and restart the client")
		return
	}

	TLSConfig, err := cert.SetupTLSConfig(&config)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := tls.Dial("tcp", config.ServerAddress, TLSConfig)
		if err != nil {
			log.Printf("Error dialing, trying again in 5 seconds: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		serverConnection := connection.NewConnection(conn, &config)
		serverConnection.ProcessLines() // Loops forever until disconnected

		log.Printf("Disconnected, trying again in 5 seconds")
		time.Sleep(5 * time.Second)
	}

}
