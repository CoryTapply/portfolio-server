package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "host.docker.internal"
	pgPort   = 5432
	user     = "postgres"
	password = "mypassword"
	dbname   = "portfolio"
)

func getDB() (db *sql.DB) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, pgPort, user, password, dbname,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	return db
}

func createVideo(id string, title string, tags string, game string, hasVoice bool, duration string) {
	db := getDB()
	defer db.Close()

	sqlStatement := `
		INSERT INTO VideoSchema.tblVideo (VideoID, Title, Tags, Game, HasVoice, Duration)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING VideoID
	`
	returnedId := ""
	err := db.QueryRow(sqlStatement, id, title, tags, game, hasVoice, duration).Scan(&returnedId)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", returnedId)
}

func getVideosFromDB() (videos []Video) {
	db := getDB()
	defer db.Close()

	sqlStatement := `
			SELECT VideoID, Title, Tags, Game, HasVoice, ViewCount, Duration FROM VideoSchema.tblVideo
	`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		video := Video{}
		err := rows.Scan(&video.ID, &video.Title, &video.Tags, &video.Game, &video.HasVoice, &video.ViewCount, &video.Duration)
		if err != nil {
			log.Println(err)
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return
	}

	return videos
}
