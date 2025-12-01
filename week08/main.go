package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"path/filepath"
	"log"
)

func avgHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォーム解析エラー", http.StatusBadRequest)
		log.Println("ParseForm error:", err)
		return
	}

	data := strings.Split(r.FormValue("scores"), ",")
	var sum, count int
	bin := make([]int, 11)

	for _, s := range data {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || n < 0 || n > 100 {
			http.Error(w, "0~100の整数のみ入力してください", http.StatusBadRequest)
			return
		}
		sum += n
		count++
		bin[n/10]++
	}

	if count == 0 {
		http.Error(w, "得点が入力されていません", http.StatusBadRequest)
		return
	}

	avg := float64(sum) / float64(count)

	if _, err := fmt.Fprintf(w, "平均値 = %.2f\n\n10点刻み分布:\n", avg); err != nil {
		log.Println("Write response error:", err)
	}
	for i := 0; i < 11; i++ {
		if _, err := fmt.Fprintf(w, "%2d〜%3d: %d\n", i*10, i*10+9, bin[i]); err != nil {
			log.Println("Write response error:", err)
		}
	}
}

func main() {
	publicPath := filepath.Join("..", "public")
	fs := http.FileServer(http.Dir(publicPath))
	http.Handle("/", fs)
	http.HandleFunc("/api/average", avgHandler)

	fmt.Println("Average server running on :8086")
	if err := http.ListenAndServe(":8086", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
