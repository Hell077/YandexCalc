package internal

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result *float64 `json:"result,omitempty"`
	Error  *string  `json:"error,omitempty"`
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
		// Если ошибка - ErrInvalidExpression, возвращаем код 422
		if errors.Is(err, ErrInvalidExpression) {
			sendErrorResponse(w, "Expression is not valid", http.StatusUnprocessableEntity)
		} else {
			sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	sendSuccessResponse(w, result)
}

func sendSuccessResponse(w http.ResponseWriter, result float64) {
	response := Response{Result: &result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := Response{Error: &errorMsg}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
