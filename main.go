package main

import (
	"find-printers/borg"
	"find-printers/ui"
	"fmt"
	"golang.design/x/clipboard"
	"net/http"
	"os"
	"sort"
	"time"
)

const BEARER = "81ebcb62699d5cd7a4b8fd0f0ffb012b962ab5ec"

func main() {
	err := clipboard.Init()
	if err != nil {
		fmt.Println("Error initializing clipboard!")
		panic(err)
	}

	client := &http.Client{Timeout: 2 * time.Second}
	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.9,hu;q=0.8,es;q=0.7",
		"Authorization":   "Bearer " + BEARER,
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
