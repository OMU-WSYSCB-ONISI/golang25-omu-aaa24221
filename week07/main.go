package main

import (
	"fmt"
	"net/http"
	"strconv"
	"path/filepath"
	"log"
)

func calcHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "フォーム解析エラー", http.StatusBadRequest)
		log.Println("ParseForm error:", err)
		return
	}

	x, err1 := strconv.Atoi(r.FormValue("x"))
	y, err2 := strconv.Atoi(r.FormValue("y"))
	op := r.FormValue("op")
	if err1 != nil || err2 != nil {
		http.Error(w, "整数を入力してください", http.StatusBadRequest)
		return
	}

	var result int
	switch op {
	case "+":
		result = x + y
	case "-":
		result = x - y
	case "*":
		result = x * y
	case "/":
		if y == 0 {
			http.Error(w, "0で割れません", http.StatusBadRequest)
			return
		}
		result = x / y
	default:
		http.Error(w, "演算子が不正です", http.StatusBadRequest)
		return
	}

	if _, err := fmt.Fprintf(w, "結果 = %d", result); err != nil {
		log.Println("Write response error:", err)
	}
}

func main() {
	publicPath := filepath.Join("..", "public")
	fs := http.FileServer(http.Dir(publicPath))
	http.Handle("/", fs)
	http.HandleFunc("/api/calc", calcHandler)

	fmt.Println("Calc server running on :8085")
	if err := http.ListenAndServe(":8085", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
