package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ilyakaznacheev/sdk-server/internal/model"
)

var testML = modelDict{
	"proxy": module{
		version: 5,
		module: &model.Module{
			File: model.ModuleFile{
				Path: "some/path/proxy",
				Sum:  "12345",
			},
			Meta: model.ModuleMeta{
				Name:    "ProxyTest",
				ID:      "proxy",
				Version: 5,
				Extension: map[string]string{
					"ext1": "test1",
					"ext2": "test2",
				},
			},
		},
	},
	"megatest": module{
		version: 2,
		module: &model.Module{
			File: model.ModuleFile{
				Path: "some/path/megatest",
				Sum:  "7890738978979",
			},
			Meta: model.ModuleMeta{
				Name:    "MegaTest",
				ID:      "megatest",
				Version: 2,
			},
		},
	},
}

func testHandle(path string, w http.ResponseWriter, req *http.Request, h func(http.ResponseWriter, *http.Request)) {
	r := mux.NewRouter()
	r.HandleFunc(path, h).Methods(req.Method)
	r.ServeHTTP(w, req)
}

func testCheckResponse(loc string, t *testing.T, w *httptest.ResponseRecorder, respStatus int, respBody string) {
	if w.Code != respStatus {
		t.Errorf("[%s]:\twrong StatusCode: got %d, expected %d",
			loc, w.Code, respStatus)
	}

	body, _ := ioutil.ReadAll(w.Result().Body)
	bodyStr := string(body)
	if bodyStr != respBody {
		t.Errorf("[%s]:\twrong Response: got \n%s\n, expected \n%s",
			loc, bodyStr, respBody)
	}
}

func TestHandlerHandleSync(t *testing.T) {
	cases := []struct {
		Name     string
		Request  string
		Response string
		Status   int
	}{
		{
			Name:     "single",
			Request:  `{ "deviceId": "testID", "installedModules": [ { "id": "proxy", "version": 3 } ] }`,
			Response: `{"id":["proxy"]}`,
			Status:   http.StatusOK,
		},
		{
			Name:     "several",
			Request:  `{ "deviceId": "testID", "installedModules": [ { "id": "megatest", "version": 1 }, { "id": "proxy", "version": 3 } ] }`,
			Response: `{"id":["megatest","proxy"]}`,
			Status:   http.StatusOK,
		},
		{
			Name:     "one new, one old",
			Request:  `{ "deviceId": "testID", "installedModules": [ { "id": "megatest", "version": 5 }, { "id": "proxy", "version": 1 } ] }`,
			Response: `{"id":["proxy"]}`,
			Status:   http.StatusOK,
		},
		{
			Name:     "both new",
			Request:  `{ "deviceId": "testID", "installedModules": [ { "id": "megatest", "version": 15 }, { "id": "proxy", "version": 10 } ] }`,
			Response: ``,
			Status:   http.StatusOK,
		},
	}

	for _, c := range cases {
		h := Handler{testML}

		url := "/sync"
		req := httptest.NewRequest("GET", url, bytes.NewBuffer([]byte(c.Request)))
		w := httptest.NewRecorder()

		testHandle("/sync", w, req, h.HandleSync)
		testCheckResponse("Sync Versions:"+c.Name, t, w, c.Status, c.Response)
	}
}

func TestHandlerHandleGetModule(t *testing.T) {
	cases := []struct {
		Name     string
		ID       string
		Response string
		Status   int
	}{
		{
			Name:     "first",
			ID:       "proxy",
			Response: `{"package":"ProxyTest","id":"proxy","version":5,"extension":{"ext1":"test1","ext2":"test2"},"checksum":"12345","url":"http://example.com/module/bin/proxy"}`,
			Status:   http.StatusOK,
		},
		{
			Name:     "second",
			ID:       "megatest",
			Response: `{"package":"MegaTest","id":"megatest","version":2,"checksum":"7890738978979","url":"http://example.com/module/bin/megatest"}`,
			Status:   http.StatusOK,
		},
		{
			Name:     "none",
			ID:       "outsider",
			Response: http.StatusText(http.StatusNotFound) + "\n",
			Status:   http.StatusNotFound,
		},
	}

	for _, c := range cases {
		h := Handler{testML}

		url := "/module/" + c.ID
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		testHandle("/module/{id}", w, req, h.HandleGetModule)
		testCheckResponse("Module Info:"+c.Name, t, w, c.Status, c.Response)
	}
}
