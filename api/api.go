package api

import (
	"net/http"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/signatures"
	log "github.com/Sirupsen/logrus"

	"encoding/json"

	"fmt"

	"net/url"

	"github.com/colebrumley/operator/tasks"
	"github.com/gorilla/mux"
)

// OperatorAPI is the wrapper object for the API server.
type OperatorAPI struct {
	ListenAddress, Cert, Key, Password string
	UseTLS, UseBasicAuth               bool
	Version                            int
}

// Start runs the API web server
func (o *OperatorAPI) Start(server *machinery.Server) {
	r := mux.NewRouter().StrictSlash(true)

	basePath := fmt.Sprintf("/api/v%v", o.Version)

	// Default to a 404
	r.HandleFunc("/*", DefaultHandler)

	// Determine routes to add from what's in TaskList
	for path := range tasks.TaskList {
		var (
			handler, helphandler func(http.ResponseWriter, *http.Request)
		)
		p := fmt.Sprintf("%s/%s", basePath, path)
		if o.UseBasicAuth && len(o.Password) > 0 {
			handler = BasicAuth(o.Password, makeHandler(path, server))
			helphandler = BasicAuth(o.Password, makeHelpHandler(path))
		} else {
			handler = makeHandler(path, server)
			helphandler = makeHelpHandler(path)
		}

		log.Debug("Registering API route " + p)
		r.HandleFunc(p, handler)

		log.Debug("Registering API route " + p + "/help")
		r.HandleFunc(p+"/help", helphandler)
	}

	// Check for TLS settings and start a TLS listener if necessary
	// otherwise use HTTP
	if o.UseTLS {
		log.Fatal(http.ListenAndServeTLS(o.ListenAddress, o.Cert, o.Key, r))
	} else {
		log.Fatal(http.ListenAndServe(o.ListenAddress, r))
	}
}

// DefaultHandler always returns a 404. This is just a catchall for unhandled requests.
func DefaultHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte("There is no handler defined for this request."))
}

func makeHelpHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Debug("Received API request " + req.RequestURI)
		rw.Write([]byte(tasks.TaskList[path].Description))
	}
}

func makeHandler(name string, server *machinery.Server) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Debug("Received API request " + req.RequestURI)

		// Separate out URL query params from POST body params
		formValues, urlValues, err := extractURLValues(req)
		if err != nil {
			log.Error("Could not parse request: ", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		wait := urlValues.Get("wait")

		taskargs := []signatures.TaskArg{}
		// We are expecting a single JSON doc as the body
		for name := range formValues {
			if len(name) > 0 {
				taskargs = append(taskargs, signatures.TaskArg{
					Type:  "string",
					Value: name,
				})
			}
		}

		sendResult, err := server.SendTask(&signatures.TaskSignature{
			Name: name,
			Args: taskargs,
		})

		if err != nil {
			errmsg := "Task failed: " + err.Error()
			log.Error(errmsg)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		state := sendResult.GetState()
		if len(wait) > 0 && wait == "true" && !state.IsCompleted() {
			for {
				state = sendResult.GetState()
				if state.IsCompleted() {
					break
				}
			}
		}
		data, err := json.Marshal(state)
		if err != nil {
			errmsg := "Task failed: " + err.Error()
			log.Error(errmsg)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Write(data)
	}
}

// Workaround for separating URL and POST body params, which normally get merged into req.Form
// https://github.com/golang/go/issues/3630
func extractURLValues(req *http.Request) (form url.Values, query url.Values, err error) {
	if query, err = url.ParseQuery(req.URL.RawQuery); err != nil {
		return
	}
	// Blank RawQuery before parsing the rest of the form
	req.URL.RawQuery = ""

	err = req.ParseForm()
	form = req.Form
	return
}
