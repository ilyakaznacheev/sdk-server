package server

import "github.com/ilyakaznacheev/sdk-server/internal/model"

type moduleFS struct {
	model.ModuleMeta
	Path string `json:"path"`
}
