package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var fortunes = []string{
	"大吉",
	"中吉",
	"吉",
	"小吉",
	"末吉",
	"凶",
}

func webFortuneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fortune := fortunes[rand.Intn(len(fortunes))]

	fmt.Fprintf(w, "今の運勢は %s です。", fortune)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/webfortune", webFortuneHandler)

	addr := ":8081"
	log.Printf("Server running at http://localhost%s/webfortune", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
