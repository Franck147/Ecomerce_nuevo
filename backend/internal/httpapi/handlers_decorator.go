package httpapi

import (
	"encoding/json"
	"net/http"

	"ecommerce-patrones/internal/pricing"
)

func (a *API) listDecorators(w http.ResponseWriter, r *http.Request) {
	items, err := a.Pricing.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) getDecorator(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := a.Pricing.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) createDecorator(w http.ResponseWriter, r *http.Request) {
	var in pricing.CreateDecoratorInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	res, err := a.Pricing.Create(r.Context(), in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (a *API) updateDecorator(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in pricing.UpdateDecoratorInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	res, err := a.Pricing.Update(r.Context(), id, in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (a *API) deleteDecorator(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := a.Pricing.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "decorador eliminado exitosamente"})
}

func (a *API) calculateQuote(w http.ResponseWriter, r *http.Request) {
	var req pricing.QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "solicitud JSON inválida")
		return
	}
	breakdown, err := a.Pricing.CalculateQuote(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, breakdown)
}
