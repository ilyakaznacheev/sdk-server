package server

import (
	"github.com/ilyakaznacheev/sdk-server/internal/model"
)

type module struct {
	version int
	module  *model.Module
}

type modelDict map[string]module
