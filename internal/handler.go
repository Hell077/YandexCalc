package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result *float64 `json:"result,omitempty"`
	Error  *string  `json:"error,omitempty"`
	Code   *int     `json:"code,omitempty"`
}

func isValidExpression(expression string) bool {
	return regexp.MustCompile(`^[0-9+\-*/().\s]+$`).MatchString(expression)
}

func evaluateExpression(expression string) (float64, error) {
	result, err := strconv.ParseFloat(expression, 64)
	if err != nil {
		return 0, errors.New("error in expression evaluation")
	}
	return result, nil
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid JSON format", http.StatusInternalServerError)
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
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
	response := Response{
		Result: &result,
		Code:   new(int),
	}
	*response.Code = http.StatusOK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, errorMsg string, statusCode int) {
	response := Response{
		Error: &errorMsg,
		Code:  new(int),
	}
	*response.Code = statusCode
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
