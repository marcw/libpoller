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
// * PUT will CLEAR the CheckList (AND the scheduler) and recreates a new CheckList from the Input
// * After any of POST or PUT operation, the configuration is persisted to it's store.
func NewConfigHttpHandler(config *Config) http.Handler {
	return &configHttpHandler{config}
}

func (h *configHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := h.config.checks.JSON()
		if err != nil {
			// TODO: Log the error on the server and do not output error to the client
			http.Error(w, err.Error(), 500)
		}

		w.Write(data)

		return
	}

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
		h.config.Add(check)
		h.config.Persist()

		w.WriteHeader(201)
		return
	}

	// PUT is an idempotent method. We could use a more clever way of
	// achieving idempotence than by calling Clear and Add but that will be for later
	if r.Method == "PUT" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// TODO: Log the error on the server
			http.Error(w, err.Error(), 500)
			return
		}
		defer r.Body.Close()

		cl, err := NewCheckListFromJSON(data)
		if err != nil {
			// TODO: Log the error on the server
			println(err.Error())
			http.Error(w, err.Error(), 400)
			return
		}

		// We don't clear anything until we're sure the data is valid
		h.config.SetCheckList(cl)
		h.config.Persist()

		w.WriteHeader(201)
		return
	}
}
