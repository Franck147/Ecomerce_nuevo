package httpapi

import (
	"encoding/json"
	"net/http"

	"ecommerce-patrones/internal/catalog"
)

func (a *API) listProducts(w http.ResponseWriter, r *http.Request) {
	items, err := a.Catalog.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (a *API) getProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	item, err := a.Catalog.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) createSimpleProduct(w http.ResponseWriter, r *http.Request) {
	var in catalog.CreateSimpleInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	res, err := a.Catalog.CreateSimple(r.Context(), in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (a *API) createBundleProduct(w http.ResponseWriter, r *http.Request) {
	var in catalog.CreateBundleInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	res, err := a.Catalog.CreateBundle(r.Context(), in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, res)
}

func (a *API) updateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var in catalog.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "cuerpo JSON inválido")
		return
	}
	res, err := a.Catalog.Update(r.Context(), id, in)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (a *API) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := a.Catalog.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "producto eliminado exitosamente"})
}
