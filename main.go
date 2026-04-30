package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "modernc.org/sqlite"
)

// EnvEntry 对应数据库记录和接口格式
type EnvEntry struct {
	Ident string `json:"ident"`
	Owner string `json:"owner"`
	Date  string `json:"date"`
}

// DeleteRequest 对应删除接口的参数格式
type DeleteRequest struct {
	Idents []string `json:"idents"`
}

var db *sql.DB

func main() {
	// 1. 支持指定端口
	port := flag.String("port", "9301", "服务运行端口")
	flag.Parse()

	// 2. 初始化 SQLite 数据库
	var err error
	db, err = sql.Open("sqlite", "./env_records.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建表结构
	createTableSQL := `CREATE TABLE IF NOT EXISTS env_usage (
		ident TEXT PRIMARY KEY,
		owner TEXT,
		date TEXT
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}

	// 3. 路由注册
	http.HandleFunc("/env", envHandler)

	fmt.Printf("服务已启动，监听端口: %s...\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

// 总入口路由分发
func envHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getEnv(w, r)
	case http.MethodPost:
		postEnv(w, r)
	case http.MethodDelete:
		deleteEnv(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// --- 接口实现 ---

// 1. 查询 (GET)
func getEnv(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT ident, owner, date FROM env_usage")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := []EnvEntry{}
	for rows.Next() {
		var e EnvEntry
		if err := rows.Scan(&e.Ident, &e.Owner, &e.Date); err != nil {
			continue
		}
		results = append(results, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// 2. 添加/更新 (POST)
func postEnv(w http.ResponseWriter, r *http.Request) {
	var entries []EnvEntry
	if err := json.NewDecoder(r.Body).Decode(&entries); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tx, _ := db.Begin()
	for _, e := range entries {
		// 使用 REPLACE INTO 实现：如果 ident 存在则更新，不存在则插入
		_, err := tx.Exec("REPLACE INTO env_usage (ident, owner, date) VALUES (?, ?, ?)", e.Ident, e.Owner, e.Date)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	tx.Commit()
	w.WriteHeader(http.StatusOK)
}

// 3. 删除 (DELETE)
func deleteEnv(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Idents) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 构建参数化查询：DELETE FROM env_usage WHERE ident IN (?, ?, ...)
	placeholders := make([]string, len(req.Idents))
	args := make([]interface{}, len(req.Idents))
	for i, id := range req.Idents {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM env_usage WHERE ident IN (%s)", strings.Join(placeholders, ","))
	_, err := db.Exec(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
