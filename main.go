package main

import (
	"net/http"
	"reflect"
	"time"

	"github.com/kavkaz/gol"
)

type Ctrl struct {
	Config *Config
	H      func(*Ctrl, http.ResponseWriter, *http.Request) (int, error)
}

func FilesCreate(ctrl *Ctrl, wr http.ResponseWriter, r *http.Request) (int, error) {
	gol.Debugf("params: %#v", r.Method)
	return 200, nil
}

func FilesUpdate(ctrl *Ctrl, wr http.ResponseWriter, r *http.Request) (int, error) {
	gol.Debugf("incoming %s", r.Method)
	return 200, nil
}

func FileServe(ctrl *Ctrl, rw http.ResponseWriter, r *http.Request) (int, error) {
	http.ServeFile(rw, r, "./"+r.URL.Path[1:])
	return 0, nil
}

// Return http status from ResponseWriter (*http.response)
func GetStatus(rw http.ResponseWriter) int64 {
	value := reflect.ValueOf(rw).Elem()
	return value.FieldByName("status").Int()
}

func (ctrl *Ctrl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t1 := time.Now()
	status, err := ctrl.H(ctrl, w, r)
	if err != nil {
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
	t2 := time.Now()

	if status == 0 {
		status = int(GetStatus(w))
	}

	gol.Infof("%s [%d] %s %s [%v]", r.RemoteAddr, status, r.Method, r.URL.Path, t2.Sub(t1))
}

func main() {
	config := ParseArgs()
	gol.SetLevel(gol.DEBUG)

	http.Handle("/files", &Ctrl{config, FilesCreate})
	http.Handle("/files/", &Ctrl{config, FilesUpdate})
	http.Handle("/", &Ctrl{config, FileServe})
	//http.Handle("/", http.FileServer(http.Dir("./")))

	gol.Infof("Start server on %s", config.Host)
	gol.Fatal(http.ListenAndServe(config.Host, nil))

}
