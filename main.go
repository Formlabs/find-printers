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

const BEARER = "81ebcb62699d5cd7a4b8fd0f0ffb012b962ab5ec"
const TIMEOUT_SECONDS = 4

func main() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Error initializing clipboard!")
		panic(err)
	}

	client := &http.Client{Timeout: TIMEOUT_SECONDS * time.Second}
	headers := map[string]string{
		"Authorization": "Bearer " + BEARER,
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
