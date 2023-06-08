package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type YoutubeStats struct {
	Subscribers    int
	ChannelName    string
	MinutesWatched int
	Views          int
}

func getChannelStats() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		yt := YoutubeStats{
			Subscribers:    5,
			ChannelName:    "yo",
			MinutesWatched: 10,
			Views:          9,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(yt); err != nil {
			panic(err)
		}
	}
}
