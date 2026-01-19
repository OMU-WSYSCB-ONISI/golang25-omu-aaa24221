package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

const memoDir = "data/memos"

// ===== データ構造 =====
type Memo struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"` // Markdown
	Created time.Time `json:"created"`
}

// ===== main =====
func main() {
	os.MkdirAll(memoDir, 0755)

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/api/memos", basicAuth(memosHandler))
	http.HandleFunc("/api/memos/", basicAuth(memoHandler))

	fmt.Println("Server started at http://localhost:8090")
	http.ListenAndServe(":8090", nil)
}

// ===== 認証 =====
func basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok ||
			user != os.Getenv("MEMO_USER") ||
			pass != os.Getenv("MEMO_PASS") {
			w.Header().Set("WWW-Authenticate", `Basic realm="memo"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// ===== 画面 =====
func indexPage(w http.ResponseWriter, r *http.Request) {
	memos := loadAllMemos()

	fmt.Fprint(w, `<html><body>
	<h1>メモ一覧</h1>
	<form method="post" action="/api/memos">
	タイトル:<br>
	<input name="title"><br>
	内容(Markdown):<br>
	<textarea name="content" rows="6" cols="50"></textarea><br>
	<button type="submit">保存</button>
	</form><hr>`)

	for _, m := range memos {
		html := blackfriday.Run([]byte(m.Content))
		fmt.Fprintf(w,
			"<h2>%s</h2><div>%s</div>"+
				"<form method='post' action='/api/memos/%d?_method=DELETE'>"+
				"<button>削除</button></form><hr>",
			m.Title, html, m.ID)
	}

	fmt.Fprint(w, "</body></html>")
}

// ===== /api/memos =====
func memosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(loadAllMemos())
	case "POST":
		r.ParseForm()
		m := Memo{
			ID:      int(time.Now().Unix()),
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
			Created: time.Now(),
		}
		saveMemo(m)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// ===== /api/memos/{id} =====
func memoHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/memos/")
	id, _ := strconv.Atoi(idStr)

	if r.Method == "POST" && r.URL.Query().Get("_method") == "DELETE" {
		deleteMemo(id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// ===== ファイル操作 =====
func saveMemo(m Memo) {
	b, _ := json.MarshalIndent(m, "", "  ")
	os.WriteFile(filepath.Join(memoDir, fmt.Sprintf("%d.json", m.ID)), b, 0644)
}

func deleteMemo(id int) {
	os.Remove(filepath.Join(memoDir, fmt.Sprintf("%d.json", id)))
}

func loadAllMemos() []Memo {
	files, _ := os.ReadDir(memoDir)
	var memos []Memo
	for _, f := range files {
		b, _ := os.ReadFile(filepath.Join(memoDir, f.Name()))
		var m Memo
		json.Unmarshal(b, &m)
		memos = append(memos, m)
	}
	return memos
}
