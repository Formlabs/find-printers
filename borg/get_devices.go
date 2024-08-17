package borg

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Device struct {
	IPAddress       string `json:"ip_address"`
	Serial          string `json:"serial"`
	MachineTypeID   string `json:"machine_type_id"`
	FirmwareVersion string `json:"firmware_version"`
	MachineType     string `json:"machine_type"`
	IsDebugFirmware bool   `json:"is_debug_firmware"`
}

func GetDevices(c *http.Client, headers map[string]string) ([]Device, error) {
	req, _ := http.NewRequest("GET", "https://case.factory.priv.prod.gcp.formlabs.cloud/api/v1/devices/", nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("Borg is unreachable, make sure you are in Formlabs WIFI or connected to VPN!\nError: %s\n", err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("Failed to close response body", slog.Any("error", err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error during parsing response", slog.Any("error", err))
		return nil, err
	}
	var devices []Device
	err = json.Unmarshal(body, &devices)
	if err != nil {
		slog.Error("Error during JSON unmarshalling response", slog.Any("error", err))
		return nil, err
	}
	return devices, nil
}
