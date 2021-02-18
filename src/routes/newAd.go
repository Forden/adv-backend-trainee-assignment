package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"adv-backend-trainee-assignment/src/models"
)

func (server APIServer) NewAd(w http.ResponseWriter, r *http.Request) {
	var adData models.CreatingAd
	err := server.parseRequest(r, &adData)
	if err == nil {
		if 1 <= len(adData.Title) && len(adData.Title) <= 200 && 1 <= len(adData.Description) && len(adData.Description) <= 1000 && 1 <= len(adData.PhotoLinks) && len(adData.PhotoLinks) <= 3 && 1 <= adData.Price {
			adId, err := server.DBManager.NewAd(adData)
			if err != nil {
				http.Error(w, "error creating ad in db", http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				_ = json.NewEncoder(w).Encode(models.CreatedAd{AdID: adId})
			}
		} else {
			http.Error(w, "exceeding data limitations", http.StatusBadRequest)
		}
	} else {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
	}
}
