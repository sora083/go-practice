package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Sheet struct {
	ID    int    `json:"id,omitempty"`
	Rank  string `json:"rank,omitempty"`
	Num   string `json:"num,omitempty"`
	Price string `json:"price,omitempty"`
}

var db *sql.DB

func main() {

	dbconf := "root:@tcp(localhost:3306)/torb"
	db, err := sql.Open("mysql", dbconf)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	e.GET("/sheets", func(c echo.Context) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Commit()

		rows, err := tx.Query("SELECT * FROM sheets")

		if err != nil {
			log.Fatal(err)
		}

		var sheets []Sheet

		for rows.Next() {
			var s Sheet
			if err := rows.Scan(&s.ID, &s.Rank, &s.Num, &s.Price); err != nil {
				log.Fatal(err)
			}
			sheets = append(sheets, s)
		}

		return c.JSON(http.StatusOK, sheets)
	})

	// サーバー起動
	e.Start(":8000")
}
