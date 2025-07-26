package handlers

import (
	"encoding/json"
	"net/http"
	"showroom/sensors"
)

func HeightsHandler(w http.ResponseWriter, r *http.Request) {
	left, right, err := sensors.GetBothHeights()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{
		"left_cm":  left,
		"right_cm": right,
		"diff":     left - right,
	})
}
