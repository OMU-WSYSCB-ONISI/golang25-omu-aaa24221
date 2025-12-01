package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
    // 現在時刻
    now := time.Now().Format("15:04:05")

    // User-Agent を取得
    ua := r.UserAgent()

    // 出力
    fmt.Fprintf(w, "今の時刻は %s で，利用しているブラウザは「%s」ですね。\n", now, ua)
}

func main() {
    // ハンドラ登録
    http.HandleFunc("/info", infoHandler)

    fmt.Println("Server running at http://localhost:8080/info")
    // サーバー起動
    log.Fatal(http.ListenAndServe(":8082", nil))
}

