package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"adv-backend-trainee-assignment/src/models"
	log "github.com/sirupsen/logrus"
)

func convertDBAdToBasicAd(dbAd *models.DbAd) *models.BasicAd {
	return &models.BasicAd{
		AdID:          dbAd.AdID,
		Title:         dbAd.Title,
		Price:         dbAd.Price,
		MainPhotoLink: dbAd.PhotoLinks[0],
	}
}

func (server APIServer) GetAllAds(w http.ResponseWriter, r *http.Request) {
	var err error
	q := r.URL.Query()
	sortBy := q.Get("sortby")
	if sortBy == "createdat" {
		sortBy = "created_at"
	}
	if !(sortBy == "created_at" || sortBy == "price") {
		sortBy = "created_at"
	}
	sortDirection := q.Get("sortdirection")
	if !(sortDirection == "asc" || sortDirection == "desc") {
		sortDirection = "desc"
	}
	pageStr := q.Get("page")
	var page int
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
	}
	perPageStr := q.Get("perpage")
	var perPage int
	if perPageStr == "" {
		perPage = 10
	} else {
		perPage, err = strconv.Atoi(perPageStr)
		if err != nil {
			perPage = 10
		}
	}
	if perPage > 100 {
		perPage = 100
	} else if perPage < 1 {
		perPage = 1
	}
	adData, err := server.DBManager.GetAllAds(sortBy, sortDirection, page, perPage)
	if err != nil {
		log.Errorf("couldn't get ads from db. err: [%s]", err)
		http.Error(w, "error getting ads from db", http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		var resp []*models.BasicAd
		if adData != nil {
			for _, tmpData := range adData {
				resp = append(resp, convertDBAdToBasicAd(tmpData))
			}
			_ = json.NewEncoder(w).Encode(resp)
		} else {
			_, _ = w.Write([]byte("{}"))
		}
	}
}
