package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	godotenv.Load()
	connectToMysqlDatabase()

	router := gin.Default()
	router.GET("/albums", getAllAlbums)
	router.POST("/albums", addAlbum)
	router.GET("/albums/:id", getAlbumById)

	router.Run("localhost:8080")
}

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func connectToMysqlDatabase() {
	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_DB_USER"),
		Passwd:               os.Getenv("MYSQL_DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_DB_URI"),
		DBName:               os.Getenv("MYSQL_DB_NAME"),
		AllowNativePasswords: true,
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
}

func getAllAlbums(context *gin.Context) {
	var albums []Album
	albums, err := getAlbumsFromDatabase()

	if err != nil {
		log.Fatal(err)
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	context.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(context *gin.Context) {
	var newAlbum Album

	if err := context.BindJSON(&newAlbum); err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
	}

	id, err := addAlbumToDatabase(newAlbum)

	if err != nil {
		log.Fatal(err)
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	album, err := getAlbumFromDatabase(id)

	if err != nil {
		log.Fatal(err)
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	context.IndentedJSON(http.StatusCreated, album)
}

func getAlbumById(context *gin.Context) {
	var album Album
	var albumId int64
	var errConv error

	albumIdParam := context.Param("id")
	albumId, errConv = strconv.ParseInt(albumIdParam, 10, 64)

	if errConv != nil {
		log.Fatal(errConv)
		context.AbortWithError(http.StatusBadRequest, errConv)
	}

	album, err := getAlbumFromDatabase(albumId)

	if err != nil {
		log.Fatal(err)
		context.AbortWithError(http.StatusInternalServerError, err)
	}

	context.IndentedJSON(http.StatusOK, album)
}

func getAlbumsFromDatabase() ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")

	if err != nil {
		return albums, err
	}

	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return albums, err
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return albums, err
	}

	return albums, nil
}

func addAlbumToDatabase(album Album) (int64, error) {
	query := "INSERT INTO album (title, artist, price) VALUES (?, ?, ?)"

	insertResult, err := db.ExecContext(context.Background(), query, album.Title, album.Artist, album.Price)

	id, err := insertResult.LastInsertId()

	if err != nil {
		return id, err
	}

	return id, nil
}

func getAlbumFromDatabase(albumId int64) (Album, error) {
	var album Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", albumId)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, err
		}
		return album, err
	}

	return album, nil
}
