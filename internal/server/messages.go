package server

import "github.com/ilyakaznacheev/sdk-server/internal/model"

// RequestMessageSync is a request structure
type RequestMessageSync struct {
	DeviceID         string `json:"deviceId"`
	InstalledModules []struct {
		ID      string `json:"id"`
		Version int    `json:"version"`
	} `json:"installedModules"`
}

// ResponseMessageSync is a sync response message structure
type ResponseMessageSync struct {
	ID []string `json:"id"`
}

// ResponseMessageModule is a module download lloiiio
type ResponseMessageModule struct {
	model.ModuleMeta
	Sum          string `json:"checksum"`
	DownloafPath string `json:"url"`
}
