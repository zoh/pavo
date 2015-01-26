package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kavkaz/gol"
)

type Ctrl struct {
	Config *Config
}

func (ctrl *Ctrl) FilesCreate(wr http.ResponseWriter, r *http.Request, params httprouter.Params) (int, error) {
	gol.Debugf("params: %#v", params)
	return 200, nil
}

func wrap(action Action) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		action.ServeHTTP(rw, r, params)
	}
}

type Action func(wr http.ResponseWriter, r *http.Request, params httprouter.Params) (int, error)

func (action Action) ServeHTTP(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if status, err := action(w, r, p); err != nil {
		// We could also log our errors centrally:
		// i.e. log.Printf("HTTP %d: %v", err)
		switch status {
		// We can have cases as granular as we like, if we wanted to
		// return custom errors for specific status codes.
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		default:
			// Catch any other errors we haven't explicitly handled
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func main() {
	config := ParseArgs()
	gol.SetLevel(gol.DEBUG)

	ctrl := &Ctrl{config}

	router := httprouter.New()
	router.GET("/hello/:name", wrap(ctrl.FilesCreate))

	gol.Infof("Start server on %s", config.Host)
	gol.Fatal(http.ListenAndServe(config.Host, router))

}
