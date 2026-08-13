package helper

import (
	"encoding/json"
	"net/http"

	"backend/internal/domain"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*domain.AppError); ok {
		WriteJSON(w, appErr.Code, appErr)
		return
	}

	WriteJSON(w, http.StatusInternalServerError, domain.NewInternalError("internal server error", err.Error()))
}

func DecodeJSON(r *http.Request, v interface{}) error {
	if r.Header.Get("Content-Type") != "" && r.Header.Get("Content-Type") != "application/json" {
		return domain.NewBadRequestError("invalid content type", "Content-Type must be application/json")
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return domain.NewBadRequestError("invalid json body", err.Error())
	}
	return nil
}
