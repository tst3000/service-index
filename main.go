package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/reportportal/commons-go/commons"
	"github.com/reportportal/commons-go/conf"
	"github.com/reportportal/commons-go/server"
	"github.com/reportportal/service-index/traefik"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const httpClientTimeout = 5 * time.Second

func main() {

	cfg := conf.EmptyConfig()

	rpCfg := struct {
		*conf.ServerConfig
		LbURL string `env:"LB_URL" envDefault:"http://localhost:9091"`
	}{
		ServerConfig: cfg,
	}

	err := conf.LoadConfig(&rpCfg)
	if nil != err {
		log.Fatalf("Cannot load config %s", err.Error())
	}

	info := commons.GetBuildInfo()
	info.Name = "Index Service"

	srv := server.New(rpCfg.ServerConfig, info)

	aggregator := traefik.NewAggregator(rpCfg.LbURL, httpClientTimeout)

	srv.WithRouter(func(router *chi.Mux) {
		router.Use(middleware.Logger)
		router.NotFound(func(w http.ResponseWriter, rq *http.Request) {
			http.Redirect(w, rq, "/ui/#notfound", http.StatusFound)
		})

		router.HandleFunc("/composite/info", func(w http.ResponseWriter, r *http.Request) {
			if err := server.WriteJSON(http.StatusOK, aggregator.AggregateInfo(), w); nil != err {
				log.Error(err)
			}
		})
		router.HandleFunc("/composite/health", func(w http.ResponseWriter, r *http.Request) {
			if err := server.WriteJSON(http.StatusOK, aggregator.AggregateHealth(), w); nil != err {
				log.Error(err)
			}
		})
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/ui/", http.StatusFound)
		})
		router.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/ui/", http.StatusFound)
		})

	})
	fmt.Println(info)
	srv.StartServer()
}
