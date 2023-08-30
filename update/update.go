package update

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	clientconfig "github.com/ProjectOrangeJuice/vm-manager-client/clientConfig"
)

func UpdateIfNeeded(config *clientconfig.Config) error {
	isLatest, err := IsLatest(config.Version)
	if err != nil {
		return fmt.Errorf("could not check if latest, %s", err)
	}
	if !isLatest {
		log.Printf("Doing update")
		err = Update()
		if err != nil {
			return fmt.Errorf("could not update, %s", err)
		}
	}
	return nil
}

func WhatVersionStartup(config *clientconfig.Config) (string, error) {
	if config.Version == "" {
		// This must be the first "update". Presume we are the latest version
		v, err := getLatestVersion()
		if err != nil {
			return "", fmt.Errorf("could not get latest version, %s", err)
		}
		config.Version = v
		err = clientconfig.UpdateVersion(v)
		if err != nil {
			return "", fmt.Errorf("could not update version, %s", err)
		}
	}

	return config.Version, nil
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func getLatestVersion() (string, error) {
	releaseURL := "https://api.github.com/repos/ProjectOrangeJuice/vm-manager-client/releases/latest"
	resp, err := http.Get(releaseURL)
	if err != nil {
		return "", fmt.Errorf("could not GET, %s", err)
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("could not decode release, %s", err)
	}
	return release.TagName, nil
}

func IsLatest(currentVersion string) (bool, error) {
	latestVersion, err := getLatestVersion()
	if err != nil {
		return false, fmt.Errorf("could not get latest version, %s", err)
	}
	log.Printf("Current version is %s, latest is %s", currentVersion, latestVersion)
	if currentVersion == latestVersion {
		return true, nil
	}
	return false, nil
}

func Update() error {
	releaseURL := "https://api.github.com/repos/ProjectOrangeJuice/vm-manager-client/releases/latest"
	resp, err := http.Get(releaseURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		panic(err)
	}

	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory, %s", err)
	}
	file := path + "/vm-manager-client" + release.TagName
	fmt.Printf("File path -> %s", file)

	for _, asset := range release.Assets {
		if asset.Name == "vm-manager-client" {
			err = downloadAsset(asset.BrowserDownloadURL, file)
			break
		}
	}
	if err != nil {
		return fmt.Errorf("could not download asset, %s", err)
	}

	// Update the config
	err = clientconfig.UpdateVersion(release.TagName)
	if err != nil {
		return fmt.Errorf("could not update version, %s", err)
	}

	//Restart self
	// Prepare to re-execute the same program
	cmd := exec.Command(file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Start the new instance of the program
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("could not restart, %s", err)
	}

	// Exit the current program
	os.Exit(0)

	return nil
}

func downloadAsset(url, fileLoc string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("could not get url, %s", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(fileLoc)
	if err != nil {
		return fmt.Errorf("could not create file, %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("could not write to file, %s", err)
	}

	err = os.Chmod(fileLoc, 0755)
	if err != nil {
		return fmt.Errorf("could not chmod file, %s", err)
	}

	log.Println("Downloaded vm-manager successfully!")
	return nil
}
