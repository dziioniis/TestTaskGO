package utils

import (
	"encoding/json"
	"net/http"
	"project/models"
)



func Respond(w http.ResponseWriter, data *models.TokenDetails) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
