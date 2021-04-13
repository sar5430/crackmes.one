package controller

import (
	"app/model"
	"encoding/xml"
	"log"
	"net/http"
	"time"
)

type item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Author      string `xml:"author"`
	Category    string `xml:"category"`
	Guid        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

type rss struct {
	Version     string `xml:"version,attr"`
	Title       string `xml:"channel>title"`
	Description string `xml:"channel>description"`
	Link        string `xml:"channel>link"`
	Items       []item `xml:"channel>item"`
}

var diffs = []string{"Very Easy", "Easy", "Medium", "Hard", "Very Hard", "Insane"}

func RssCrackmesGET(w http.ResponseWriter, r *http.Request) {
	crackmes, err := model.LastCrackMes(1)
	if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	var items []item
	for _, v := range(crackmes) {

        var difficulty float64
        difficulties, err := model.RatingDifficultyByCrackme(v.HexId)
        if err != nil {
            log.Println(err)
            Error500(w, r)
            return
        }

        for _, d := range difficulties {
            difficulty += float64(d.Rating)
        }
        difficulty /= float64(len(difficulties))

		items = append(items, item{
			Title: v.Name+" ["+v.Platform+" - "+v.Lang+" - "+diffs[int(difficulty) - 1]+"]",
			Description: v.Info,
			Author: v.Author,
			PubDate: v.CreatedAt.Format(time.RFC1123Z),
			Category: v.Platform,
			Link: "https://crackmes.one/crackme/"+v.HexId,
			Guid: "https://crackmes.one/crackme/"+v.HexId,
		})
	}
	crss := rss{
		Version: "2.0",
		Title: "Latest crackmes - crackmes.one",
		Link: "https://crackmes.one/lasts",
		Description: "The latest 50 crackmes from crackmes.one",
		Items: items,
	}

	b, err := xml.Marshal(crss)
	if err != nil {
		log.Println(err)
		Error500(w, r)
		return
	}

	w.Header().Set("content-type", "application/rss+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
