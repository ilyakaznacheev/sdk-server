package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Handler is a request handler
type Handler struct {
	m modelDict
}

// NewHandler creates new handler
func NewHandler(m modelDict) *Handler {
	return &Handler{m: m}
}

// HandleSync checks versions and sends update information
func (h *Handler) HandleSync(w http.ResponseWriter, r *http.Request) {
	var syncMsg RequestMessageSync
	var respList ResponseMessageSync

	// read request body
	dc := json.NewDecoder(r.Body)
	err := dc.Decode(&syncMsg)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, m := range syncMsg.InstalledModules {
		if module, ok := h.m[m.ID]; ok && m.Version < module.version {
			respList.ID = append(respList.ID, m.ID)
		}
	}

	if len(respList.ID) != 0 {
		// fill response json
		resp, err := json.Marshal(&respList)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	}

	// setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

// HandleGetModule sends module information and download lonk
func (h *Handler) HandleGetModule(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	downloadURL := fmt.Sprintf("http://%s/module/bin/%s", r.Host, id)

	if module, ok := h.m[id]; !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		respMsg := ResponseMessageModule{
			DownloafPath: downloadURL,
			Sum:          module.module.File.Sum,
			ModuleMeta:   module.module.Meta,
		}
		// fill response json
		resp, err := json.Marshal(&respMsg)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// HandleGetModuleBinary serves file for download
func (h *Handler) HandleGetModuleBinary(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if module, ok := h.m[id]; !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		file, err := os.Open(module.module.File.Path)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		http.ServeContent(w, r, module.module.Meta.Name, time.Now(), file)
	}

}
