package rest

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	ase "github.com/x-color/calendar/app/rest/auth"
	cse "github.com/x-color/calendar/app/rest/calendar"
	"github.com/x-color/calendar/app/rest/middlewares"
	as "github.com/x-color/calendar/auth/service"
	cs "github.com/x-color/calendar/calendar/service"
	"github.com/x-color/calendar/logging"
)

func StartServer(authService as.Service, calService cs.Service, l logging.Logger, port string) {
	r := newRouter(authService, calService, l)
	http.ListenAndServe(":"+port, r)
}

func newRouter(authService as.Service, calService cs.Service, l logging.Logger) *mux.Router {
	r := mux.NewRouter()
	r.NotFoundHandler = http.NotFoundHandler()
	r.Use(middlewares.ReqIDMiddleware)
	r.Use(middlewares.LoggingMiddleware(l))

	apiRouter := r.PathPrefix("/api").Subrouter()

	ar := apiRouter.PathPrefix("/auth").Subrouter()
	ase.NewRouter(ar, authService)

	ur := apiRouter.PathPrefix("/register").Subrouter()
	cse.NewUserRouter(ur, calService, authService)

	cr := apiRouter.PathPrefix("/calendars").Subrouter()
	cse.NewCalendarRouter(cr, calService, authService)

	pr := apiRouter.PathPrefix("/plans").Subrouter()
	cse.NewPlanRouter(pr, calService, authService)

	spa := spaHandler{staticPath: "web/calendar/dist", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	return r
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	path = filepath.Join(h.staticPath, path)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
