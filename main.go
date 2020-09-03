package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thoas/go-funk"
	_ "github.com/walf443/go-sql-tracer"

	_ "net/http/pprof"

	"github.com/felixge/fgprof"

	"github.com/labstack/echo-contrib/prometheus"

	"github.com/gomodule/redigo/redis"
)

type Sheet struct {
	ID    int    `json:"id,omitempty"`
	Rank  string `json:"rank,omitempty"`
	Num   string `json:"num,omitempty"`
	Price string `json:"price,omitempty"`
}

var db *sql.DB
var conn *redis.Conn

var sheets []*Sheet // declare as global variable

func getSheets(tx *sql.Tx) ([]*Sheet, error) {

	rows, err := tx.Query("SELECT * FROM sheets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sheets []*Sheet

	for rows.Next() {
		var s Sheet
		if err := rows.Scan(&s.ID, &s.Rank, &s.Num, &s.Price); err != nil {
			log.Print(err)
			return nil, err
		}
		sheets = append(sheets, &s)
	}

	return sheets, nil
}

func main() {

	// DB connection
	dbconf := "root:@tcp(localhost:3306)/torb"
	db, err := sql.Open("mysql:trace", dbconf)
	if err != nil {
		log.Fatal(err)
	}

	// Redis connection
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// fgprof
	http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	e := echo.New()
	// Required to define Middleware if want to apply for all requests(logger etc.)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Routing
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	e.GET("/sheets", func(c echo.Context) error {

		// fetch from DB if sheets has no data
		if sheets == nil {
			tx, err := db.Begin()
			if err != nil {
				return err
			}
			defer tx.Commit()

			log.Println("Call getSheets")
			sheets, err = getSheets(tx)
			if err != nil {
				log.Fatal(err)
			}
		}

		return c.JSON(http.StatusOK, sheets)
	})

	// e.g http://localhost:8000/sheetsByRank?rank=S
	e.GET("/sheetsByRank", func(c echo.Context) error {

		rank := c.QueryParam("rank")

		if rank == "" {
			return c.JSON(http.StatusBadRequest, "rankが未指定です")
		}

		// fetch from DB if sheets has no data
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Commit()

		var list []*Sheet
		list, err = getSheets(tx)

		filtered := funk.Filter(list, func(s *Sheet) bool {
			return s.Rank == rank
		})

		return c.JSON(http.StatusOK, filtered)
	})

	e.POST("/redis", func(c echo.Context) error {
		// 値の書き込み
		r, err := conn.Do("SET", "temperature", "25")
		if err != nil {
			panic(err)
		}
		log.Println(r)
		return c.JSON(http.StatusOK, "PUT SUCCESS")
	})

	e.GET("/redis", func(c echo.Context) error {
		// 値の読み出し
		s, err := redis.String(conn.Do("GET", "temperature"))
		if err != nil {
			panic(err)
		}
		log.Println(s)
		return c.JSON(http.StatusOK, "GET SUCCESS")
	})

	e.Logger.Fatal(e.Start(":8000"))
}
