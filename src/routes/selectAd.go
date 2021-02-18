package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"adv-backend-trainee-assignment/src/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func convertDBAdToExtendedAd(dbAd *models.DbAd) *models.ExtendedAd {
	return &models.ExtendedAd{
		Title:         dbAd.Title,
		Price:         dbAd.Price,
		MainPhotoLink: dbAd.PhotoLinks[0],
		Description:   dbAd.Description,
		PhotoLinks:    dbAd.PhotoLinks,
	}
}

func (server APIServer) SelectAd(w http.ResponseWriter, r *http.Request) {
	if adID, ok := mux.Vars(r)["adID"]; ok {
		adData, err := server.DBManager.SelectAd(adID)
		if err != nil {
			log.Errorf("couldn't get ad with id %s from db. err: [%s]", adID, err)
			http.Error(w, "error getting ad from db", http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			if adData != nil {
				fullAd := convertDBAdToExtendedAd(adData)
				fields := r.URL.Query().Get("fields")
				fullAd.Description = ""
				fullAd.PhotoLinks = []string{}
				if fields != "" {
					for _, word := range strings.Split(fields, ",") {
						switch word {
						case "description":
							fullAd.Description = adData.Description
						case "photolinks":
							fullAd.PhotoLinks = adData.PhotoLinks
						}
						if fullAd.Description != "" && len(fullAd.PhotoLinks) != 0 {
							break
						}
					}
				}
				_ = json.NewEncoder(w).Encode(fullAd)
			} else {
				_, _ = w.Write([]byte("{}"))
			}
		}
	} else {
		http.Error(w, "couldn't extract ad id from urlFormat", http.StatusBadRequest)
	}
}
