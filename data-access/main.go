package main

import (
	"database/sql"
	"fmt"
	"log"
	// "os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID		int64
	Title	string
	Artist	string
	Price	float32
}

func main(){
	//Capture connection properties
	cfg := mysql.Config{
		// User: os.Getenv("DBUSER"),
		User: "root",
		// Passwd: os.Getenv("DBPASS"),
		Passwd: "password",
		Net: "tcp",
		Addr: "127.0.0.1:3306",
		DBName: "recordings",
	}
	//Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// get all albums by artist
	albums, err := albumsByArtist("John Coltrane")
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	//get album by id
	alb, err := albumById(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	//add new album
	albId, err := addAlbum(Album{
		Title: "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price: 49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID od added album: %v\n", albId)
}

// albumsByArtist queries for albums that hace the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil{
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next(){
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil{
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

func albumById(id int64) (Album, error){
	//An album to hold data for the returned row.
	var alb Album
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumById %d: %v", id, err)
	}
	return alb, nil
}

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err  := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}