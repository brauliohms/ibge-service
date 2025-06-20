package http

import (
	"net/http"
	"time"

	"github.com/brauliohms/ibge-service/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRouter(handler *IBGEHandler) http.Handler {
	r := chi.NewRouter()

	// 1. Carregar configurações
	cfg := config.Load()

	// Configuração CORS baseada em variáveis de ambiente
	allowedOrigins := cfg.AllowedOrigins

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		ExposedHeaders: []string{
			"Link",
			"X-Total-Count",
		},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutos
	}))

	// Middlewares para produção
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Timeout para requisições
	r.Use(middleware.Timeout(30 * time.Second))

	// Compressão para reduzir bandwidth
	r.Use(middleware.Compress(5))

	// Headers de segurança
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			next.ServeHTTP(w, r)
		})
	})

	// Rate limiting (opcional - requer biblioteca externa)
	r.Use(httprate.LimitByIP(cfg.RateLimit, 1*time.Minute)) // 100 requests por minuto por IP

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/estados", handler.GetAllEstados)
		r.Get("/estados/{uf}", handler.GetEstadoByUF)
		r.Get("/estados/{uf}/cidades", handler.GetCidadesByEstadoUF)
		r.Get("/cidades/{codigo_ibge}", handler.GetCidadeByCodigo)
		r.Get("/cidades/{codigo_tom}/tom", handler.GetCidadeByCodigoTOM)
		r.Get("/docs/*", httpSwagger.WrapHandler)
	})

	return r
}
