package api

import (
	"net/http"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/signatures"
	log "github.com/Sirupsen/logrus"

	"encoding/json"

	"fmt"

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
		p := fmt.Sprintf("%s/%s/", basePath, path)
		handler := func(http.ResponseWriter, *http.Request) {}
		helphandler := func(http.ResponseWriter, *http.Request) {}
		if o.UseBasicAuth && len(o.Password) > 0 {
			handler = BasicAuth(o.Password, makeHandler(path, server))
			helphandler = BasicAuth(o.Password, makeHelpHandler(path))
		} else {
			handler = makeHandler(path, server)
			helphandler = makeHelpHandler(path)
		}

		log.Debug("Registering API route " + p)
		r.HandleFunc(p, handler)

		log.Debug("Registering API route " + p + "help")
		r.HandleFunc(p+"help", helphandler)
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

		if err := req.ParseForm(); err != nil {
			log.Error("Could not parse form: ", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		taskargs := []signatures.TaskArg{}
		for name, arg := range req.Form {
			if len(name) > 0 {
				taskargs = append(taskargs, signatures.TaskArg{
					Type:  "string",
					Value: name,
				})
			}

			if len(arg) > 0 {
				for _, a := range arg {
					taskargs = append(taskargs, signatures.TaskArg{
						Type:  "string",
						Value: a,
					})
				}
			}
		}

		result := ""
		sendResult, err := server.SendTask(&signatures.TaskSignature{
			Name: name,
			Args: taskargs,
		})

		if err != nil {
			errmsg := "Task failed: " + err.Error()
			log.Error(errmsg)
			result = errmsg
			rw.WriteHeader(http.StatusBadRequest)
		}

		if sendResult != nil {
			state := sendResult.GetState()
			data, err := json.Marshal(state)
			if err != nil {
				errmsg := "Task failed: " + err.Error()
				log.Error(errmsg)
				result = errmsg
				rw.WriteHeader(http.StatusInternalServerError)
			}
			result = string(data)
		}
		rw.Write([]byte(result))
	}
}
