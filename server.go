package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func main() {

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	type Todo struct {
		Id        int       `db:"pk" column:"id" json:"id"`
		Title     string    `json:"title"`
		Done      bool      `json:"done"`
		CreatedAt time.Time `json:"created_at"`
	}

	goji.Get("/api/:id", func(c web.C, w http.ResponseWriter, r *http.Request) {
		var todos []Todo
		id := c.URLParams["id"]
		row := db.QueryRow("SELECT * FROM todo WHERE id =", id)
		err = row.Scan(&todos)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todos); err != nil {
			panic(err.Error())
		}
	})

	//UPDATE一旦削除

	goji.Post("/api", func(c web.C, w http.ResponseWriter, r *http.Request) {
		todo := &Todo{
			Title: r.PostFormValue("title"),
		}
		_, err := db.Exec("INSERT todo VALUES", todo)

		if err != nil {
			panic(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err.Error())
		}
	})
	goji.Get("/api", func(c web.C, w http.ResponseWriter, r *http.Request) {
	})
	goji.Get("/*", http.FileServer(http.Dir("./assets")))
	goji.Serve()

}
