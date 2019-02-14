package server

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/sdk-server/internal/model"
)

// RunServer start server on default port
func RunServer(moduleListPath string) error {
	md, err := loadModules(moduleListPath)
	if err != nil {
		return err
	}

	h := NewHandler(*md)

	r := mux.NewRouter()
	r.HandleFunc("/sync", h.HandleSync).Methods("GET")
	r.HandleFunc("/module/{id}", h.HandleGetModule).Methods("GET")
	r.HandleFunc("/module/bin/{id}", h.HandleGetModuleBinary).Methods("GET")

	fmt.Println("Starting server")
	return http.ListenAndServe(":8000", r)
}

// loadModules loads module information from json file
func loadModules(path string) (*modelDict, error) {
	var modList []moduleFS
	md := make(modelDict)

	// read json file
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonParser := json.NewDecoder(jsonFile)
	if err = jsonParser.Decode(&modList); err != nil {
		return nil, err
	}

	// process data from config
	for _, mod := range modList {
		// calculate md5
		file, err := os.Open(mod.Path)
		if err != nil {
			return nil, err
		}

		hash := md5.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			return nil, err
		}
		sum := hex.EncodeToString(hash.Sum(nil)[:16])

		file.Close()

		// fill data
		md[mod.ModuleMeta.ID] = module{
			version: mod.ModuleMeta.Version,
			module: &model.Module{
				File: model.ModuleFile{
					Path: mod.Path,
					Sum:  sum,
				},
				Meta: mod.ModuleMeta,
			},
		}
	}

	return &md, nil
}
