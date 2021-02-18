package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"adv-backend-trainee-assignment/src/db"
	log "github.com/sirupsen/logrus"
)

type APIServer struct {
	DBManager db.DatabaseConnection
}

type (
	Route struct {
		Name        string
		Method      string
		Pattern     string
		HandlerFunc http.HandlerFunc
	}
	Routes []Route
)

func GenerateRoutes(apiServer APIServer) Routes {
	return Routes{
		Route{
			"create ad",
			"POST",
			"/ad",
			apiServer.NewAd,
		},
		Route{
			"get ad",
			"GET",
			"/ads/{adID}",
			apiServer.SelectAd,
		},
		Route{
			Name:        "get ads",
			Method:      "GET",
			Pattern:     "/ads",
			HandlerFunc: apiServer.GetAllAds,
		},
	}
}

func (server APIServer) parseRequest(r *http.Request, parseStruct interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return fmt.Errorf("can't read body")
	}
	err = json.Unmarshal(body, &parseStruct)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("can't parse body to struct")
	}
	return nil
}
