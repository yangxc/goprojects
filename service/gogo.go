package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRouters(mx, formatter)
	n.UseHandler(mx)
	return n
}

func createMatchHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newMatchRequest newMatchRequest
		err := json.Unmarshal(payload, &newMatchRequest)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse create match request")
		}
		if !newMatchRequest.isValid() {
			formatter.Text(w, http.StatusBadRequest, "Invalid new match request")
			return
		}
		newMatch := gogo.NewMatch(newMatchRequest.GridSize, newMatchRequest.PlayerBlack, newMatchRequest.PlayerWhite)
		repo.addMatch(newMatch)
		w.Header().Add("Location", "/matches/"+newMatch.ID)
		formatter.JSON(w, http.StatusCreated, &newMatchResponse{ID: newMatch.ID, GridSize: newMatch.GridSize, PlayerBlack: newMatchRequest.PlayerBlack, PlayerWhite: newMatchRequest.PlayerWhite})
	}
}

func initRouters(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/test", testHandler(formatter)).Methods("GET")
}

func testHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"This is a test"})
	}
}
