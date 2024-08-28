package main

import (
	"find-printers/borg"
	"find-printers/ui"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"golang.design/x/clipboard"
)

const TOKEN_STORE = ".config/find-printers"
const TOKEN_FILE_NAME = "token"
const TIMEOUT_SECONDS = 4

func setToken(token string) error {
	UserHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	// open TOKEN_STORE and set the token variable
	os.MkdirAll(
		fmt.Sprintf("%s/%s", UserHome, TOKEN_STORE),
		0755,
	)
	f, err := os.OpenFile(
		fmt.Sprintf("%s/%s/%s", UserHome, TOKEN_STORE, TOKEN_FILE_NAME),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	defer f.Close()
	if err != nil {
		return err
	}
	f.Write([]byte(token))
	return nil
}

func main() {
	args := os.Args

	if len(args) > 1 {
		if args[1] == "set-token" {
			if len(args) < 3 {
				fmt.Println("Usage: find-printers set-token <token>")
				os.Exit(-1)
			}
			err := setToken(args[2])
			if err != nil {
				fmt.Println("Error setting token!")
				panic(err)
			}
			fmt.Println("Token set!")
			os.Exit(0)
		} else {
			fmt.Println("Usage:")
			fmt.Println("find-printers set-token <token> - set Borg token")
			fmt.Println("find-printers - run the application")
			os.Exit(-1)
		}
	}

	// Read token from file
	UserHome, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory!")
		panic(err)
	}
	token, err := os.ReadFile(fmt.Sprintf("%s/%s/%s", UserHome, TOKEN_STORE, TOKEN_FILE_NAME))
	if err != nil {
		fmt.Println("Error reading token file!")
		panic(err)
	}
	tokenStr := string(token)

	err = clipboard.Init()
	if err != nil {
		fmt.Println("Error initializing clipboard!")
		panic(err)
	}

	client := &http.Client{Timeout: TIMEOUT_SECONDS * time.Second}
	headers := map[string]string{
		"Authorization": "Bearer " + tokenStr,
	}
	devices, err := borg.GetDevices(client, headers)
	if err != nil {
		os.Exit(-1)
	}

	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Serial < devices[j].Serial
	})

	uiModel := ui.NewModel()
	uiModel.AddDevices(devices)
	ui.Run(uiModel)
}
