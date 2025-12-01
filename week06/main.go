package main

import (
	"fmt"
	"net/http"
	"strconv"
	"path/filepath"
	"log"
)

func bmiHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォーム解析エラー", http.StatusBadRequest)
		log.Println("ParseForm error:", err)
		return
	}

	weight, err1 := strconv.ParseFloat(r.FormValue("weight"), 64)
	height, err2 := strconv.ParseFloat(r.FormValue("height"), 64)
	if err1 != nil || err2 != nil {
		http.Error(w, "数値を入力してください", http.StatusBadRequest)
		return
	}

	if height == 0 {
		http.Error(w, "身長が0です", http.StatusBadRequest)
		return
	}

	h := height / 100.0
	bmi := weight / (h * h)

	if _, err := fmt.Fprintf(w, "BMI = %.2f\n", bmi); err != nil {
		log.Println("Write response error:", err)
	}
}

func main() {
	publicPath := filepath.Join("..", "public")
	fs := http.FileServer(http.Dir(publicPath))
	http.Handle("/", fs)
	http.HandleFunc("/api/bmi", bmiHandler)

	fmt.Println("BMI server running on :8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}

