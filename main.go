package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	store "github.com/singhkshitij/GOShortener/store"
	factory "github.com/singhkshitij/GOShortener/utils"
)

var factoryUtils *factory.Factory
var dB *store.DB

type urlStruct struct {
	URL string `json:"URL"`
}

func showWelcomeMessage(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/html")
	response.Write([]byte("<center><h2>Welcome to Investment Tracker</h2></center>"))
}

func shortenLongURL(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(request.Body)
	var urlData urlStruct
	err := decoder.Decode(&urlData)
	if err != nil {
		panic(err)
	}

	longURL := urlData.URL

	key, err := factoryUtils.Gen(longURL)

	if err != nil {
		panic(err)
	}

	if err = dB.Set(key, longURL); err != nil {
		response.Header().Set("Content-Type", "text/html")
		response.Write([]byte("Internal error while saving the URL"))
	} else {
		var shortenedURL string
		response.Header().Set("Content-Type", "application/json")
		if request.TLS != nil {
			shortenedURL = "https://" + request.Host + "/u/" + key
		} else {
			shortenedURL = "http://" + request.Host + "/u/" + key
		}
		urlData.URL = shortenedURL
		log.Println("Generated Short Url : " + shortenedURL)
		json.NewEncoder(response).Encode(&urlData)
	}
}

func redirectToLongURL(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	shortURL := params["shortkey"]
	longURL := dB.Get(shortURL)
	log.Println("Fetched Long URL : " + longURL)
	http.Redirect(response, request, longURL, http.StatusTemporaryRedirect)
}

func backupDB(response http.ResponseWriter, request *http.Request) {
	response = dB.Backup("backup.db", response)
}

func main() {
	dB = store.NewDB("shortener.db")
	factoryUtils = factory.NewFactory(factory.DefaultGenerator, dB)
	port := ":8080"
	router := mux.NewRouter()

	router.HandleFunc("/", showWelcomeMessage).Methods("GET")
	router.HandleFunc("/api/v1/url/shorten", shortenLongURL).Methods("POST")
	router.HandleFunc("/u/{shortkey}", redirectToLongURL).Methods("GET")
	router.HandleFunc("/api/v1/url/backup", backupDB).Methods("GET")

	log.Println("Starting server at port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}
