package poller

import (
	"io/ioutil"
	"net/http"
)

type configHttpHandler struct {
	config *Config
}

// Create a handler function that is usable by http.Handle.
// This handler will be able response to GET, POST and PUT requests.
// * GET the list of checks as a JSON array
// * POST will create a new check and add it to the CheckList.
// * After any of POST or PUT operation, the configuration is persisted to it's store.
func NewConfigHttpHandler(config *Config) http.Handler {
	return &configHttpHandler{config}
}

func (h *configHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO: Log the error on the server
			http.Error(w, err.Error(), 500)
			return
		}
		defer r.Body.Close()

		check, err := NewCheckFromJSON(data)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := h.config.Add(check); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(201)
		return
	}
}
