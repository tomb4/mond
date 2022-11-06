package health

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

type Health struct {
	srv *http.Server
}

func NewHealth() *Health {
	return &Health{}
}
func (m *Health) Health() error {
	router := mux.NewRouter()
	router.HandleFunc("/health", m.health).Methods("GET")

	m.srv = &http.Server{
		Addr:    ":3333",
		Handler: router,
	}
	err := m.srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (m *Health) health(w http.ResponseWriter, r *http.Request) {
	//TODO： 这里面后面可以加各种逻辑，比如熔断后需要重启
	w.WriteHeader(200)
	w.Write([]byte("Health"))
}

func (m *Health) Shutdown(ctx context.Context) {
	if m.srv != nil {
		m.srv.Shutdown(ctx)

	}
}
