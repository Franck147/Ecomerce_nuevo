package httpapi

import (
	"encoding/json"
	"net/http"

	"ecommerce-patrones/internal/catalog"
	"ecommerce-patrones/internal/pricing"
)

type API struct {
	Catalog *catalog.Service
	Pricing *pricing.Service
}

func New(catalogSvc *catalog.Service, pricingSvc *pricing.Service) *API {
	return &API{
		Catalog: catalogSvc,
		Pricing: pricingSvc,
	}
}

func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()

	// Healthcheck
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "app": "ecommerce-patrones"})
	})

	// CRUD 1: COMPOSITE (Productos Simples y Combos)
	mux.HandleFunc("GET /api/v1/products", a.listProducts)
	mux.HandleFunc("GET /api/v1/products/{id}", a.getProduct)
	mux.HandleFunc("POST /api/v1/products/simple", a.createSimpleProduct)
	mux.HandleFunc("POST /api/v1/products/bundle", a.createBundleProduct)
	mux.HandleFunc("PUT /api/v1/products/{id}", a.updateProduct)
	mux.HandleFunc("DELETE /api/v1/products/{id}", a.deleteProduct)

	// CRUD 2: DECORATOR (Decoradores de Precios)
	mux.HandleFunc("GET /api/v1/decorators", a.listDecorators)
	mux.HandleFunc("GET /api/v1/decorators/{id}", a.getDecorator)
	mux.HandleFunc("POST /api/v1/decorators", a.createDecorator)
	mux.HandleFunc("PUT /api/v1/decorators/{id}", a.updateDecorator)
	mux.HandleFunc("DELETE /api/v1/decorators/{id}", a.deleteDecorator)

	// Simulador / Cotizador con patron Decorator
	mux.HandleFunc("POST /api/v1/pricing/quote", a.calculateQuote)

	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
