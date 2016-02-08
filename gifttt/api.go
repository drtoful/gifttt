package gifttt

import (
	"encoding/json"
	"net/http"

	"github.com/drtoful/gifttt/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"github.com/drtoful/gifttt/Godeps/_workspace/src/github.com/gorilla/mux"
)

type APIServer struct {
	ip, port string
	handler  *negroni.Negroni
}

var (
	internalVars = []string{
		"time:second",
		"time:minute",
		"time:hour",
		"date:day",
		"date:month",
		"date:year",
	}
)

func postVar(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	varname := vars["var"]

	for _, c := range internalVars {
		if c == varname {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	var value Value
	if err := json.NewDecoder(r.Body).Decode(&value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	vm := GetManager()
	if err := vm.Set(varname, value.Value); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getVar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varname := vars["var"]

	vm := GetManager()
	v, err := vm.Get(varname)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	val := &Value{Value: v}
	b, err := json.Marshal(val)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func NewAPIServer(ip, port string) *APIServer {
	router := mux.NewRouter()

	api := router.PathPrefix("/v").Subrouter()
	api = api.StrictSlash(true)
	api.Path("/{var}").Methods("POST").HandlerFunc(postVar)
	api.Path("/{var}").Methods("GET").HandlerFunc(getVar)

	n := negroni.New(negroni.NewRecovery())
	n.UseHandler(router)

	return &APIServer{
		ip:      ip,
		port:    port,
		handler: n,
	}
}

func (a *APIServer) Run() {
	a.handler.Run(a.ip + ":" + a.port)
}
