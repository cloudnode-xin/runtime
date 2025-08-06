package runtime

import (
	"context"
	"flag"
	"net/http"
	"os"
)

type HealthChecker interface {
	IsHealthy() bool
}

type router struct {
	checker HealthChecker
}

func (r *router) ServeHTTP(rep http.ResponseWriter, req *http.Request) {
	if r.checker.IsHealthy() {
		rep.Write([]byte("OK"))
	} else {
		rep.WriteHeader(http.StatusBadRequest)
		rep.Write([]byte("ERR"))
	}

}

type healthChecker struct {
	check  bool
	server *http.Server
	root   *Service
}

func (h *healthChecker) Name() string {
	return "#healthcheck"
}

func (h *healthChecker) IsHealthy() bool {
	if h.root == nil {
		return true
	}

	return h.root.IsHealthy()
}

func (h *healthChecker) Load(f Finder) error {
	check := flag.Bool("healthcheck", false, "do health check")
	flag.Parse()

	h.check = *check
	return nil
}

func (h *healthChecker) Start(f Finder, ctx context.Context) error {
	log := f.MustGet("#logger").(*Logger).New("healthcheck")

	if !h.check {
		h.server = &http.Server{
			Addr: ":9180",
			Handler: &router{
				checker: h,
			},
		}

		go func() {
			err := h.server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				log.Error(err)
			}
		}()
	} else {
		rep, err := http.Get("http://localhost:9180/healthcheck")
		if err != nil {
			os.Exit(1)
		} else if rep.StatusCode != 200 {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	hn, err := os.Hostname()
	if err != nil {
		return err
	}

	log.Infof("hostname: %s", hn)
	return nil
}

func (h *healthChecker) Stop(f Finder) error {
	if !h.check && h.server != nil {
		return h.server.Shutdown(context.Background())
	}

	return nil
}

func GlobalHealthChecker() Servicer {
	return &healthChecker{}
}
