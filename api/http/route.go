package http

import (
	"encoding/json"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"oracle.com/oracle/my-go-oracle-app/api"
	"oracle.com/oracle/my-go-oracle-app/api/http/member"
	config "oracle.com/oracle/my-go-oracle-app/configs"
	"oracle.com/oracle/my-go-oracle-app/pkg/helpers"
	"oracle.com/oracle/my-go-oracle-app/pkg/logger"
	"oracle.com/oracle/my-go-oracle-app/pkg/panics"
)

func root(w http.ResponseWriter, r *http.Request) {
	app := map[string]interface{}{
		"name":        "my-go-oracle-app",
		"description": "my-go-oracle-app",
	}

	data, _ := json.Marshal(app)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func handler(checker api.HealthChecker, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(panics.HTTPRecoveryMiddleware)
	r.Use(middleware.Timeout(cfg.HttpInboundTimeout))

	//skip middleware group
	r.Get("/application/health", checker.HealthChi)
	r.Get("/", root)
	r.Handle("/metrics", promhttp.Handler())
	if helpers.GetEnvString() != helpers.EnvProduction {
		r.Get("/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./docs/swagger.json")
		}))

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/swagger.json"),
		))

	}

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequestLogger(&logger.CustomLogFormatter{Logger: logger.NewSlogWrapper(cfg)}))

		cors := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		})
		r.Use(cors.Handler)

		// Test Panics to Slack function
		r.Handle("/panics", panics.CaptureHandler(func(w http.ResponseWriter, r *http.Request) {
			panic("Panics from /test/panics endpoint")
		}))

		r.With(api.InterceptorRequest()).Route("/my-go-oracle-app", func(r chi.Router) {
			r.Use(api.NewMetricMiddleware())
			// members group
			r.Route("/members", func(r chi.Router) {
				r.Get("/", member.GetAllMembers)
				r.Get("/{id}", member.GetMemberById)
			})

		})
	})

	return r
}
