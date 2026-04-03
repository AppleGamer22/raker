package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/AppleGamer22/raker/server/db"
	old "github.com/AppleGamer22/raker/server/db/mongo"
	_ "github.com/lib/pq"
)

// mongoexport -d raker -c users -o users.json --jsonArray
func users(ctx context.Context, pgdb *db.Queries) {
	usersFile, err := os.Open("users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer usersFile.Close()

	decoder := json.NewDecoder(usersFile)
	token, err := decoder.Token()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	for decoder.More() {
		var user old.User
		if err := decoder.Decode(&user); err != nil {
			fmt.Println(err)
			break
		}
		err := pgdb.UserAdd(ctx, db.UserAddParams{
			Username:           user.Username,
			PasswordHash:       user.Hash,
			InstagramSessionID: user.Instagram.SessionID,
			InstagramUserID:    user.Instagram.UserID,
			Categories:         user.Categories,
		})
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("Inserted user %+v\n", user)
	}
}

// mongoexport -d raker -c histories -o histories.json --jsonArray
func histories(ctx context.Context, pgdb *db.Queries) {
	historiesFile, err := os.Open("histories.json")
	if err != nil {
		log.Fatal(err)
	}
	defer historiesFile.Close()

	decoder := json.NewDecoder(historiesFile)
	decoder.Token()

	for decoder.More() {
		var history old.HistoryArchive
		if err := decoder.Decode(&history); err != nil {
			fmt.Println(err)
			var sb strings.Builder
			io.Copy(&sb, decoder.Buffered())
			fmt.Println(sb.String())
			break
		}
		// fmt.Printf("Read history %+v\n", history)
		_, err := pgdb.HistoryGet(ctx, db.HistoryGetParams{
			Type: db.PostType(history.Type),
			Post: history.Post,
		})
		// skip if already inserted
		if err == nil {
			// fmt.Println(err)
			continue
		}
		row := db.HistoryAddFromArchiveParams{
			PostType:   db.PostType(history.Type),
			PostOwner:  history.Owner,
			Post:       history.Post,
			Files:      history.URLs,
			Categories: history.Categories,
			PostDate:   history.Date.Value,
		}
		_, err = pgdb.HistoryAddFromArchive(ctx, row)
		if err != nil {
			fmt.Println(err, row)
			continue
		}
		// fmt.Printf("Inserted history %s\n", h.Post)
	}
}

func main() {
	ctx := context.Background()
	connection, err := sql.Open("postgres", "postgres://applegamer22:postgres@rpi4b12/raker?sslmode=prefer")

	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	if err := connection.Ping(); err != nil {
		log.Fatal(err)
	}

	pgdb := db.New(connection)

	// users(ctx, pgdb)
	histories(ctx, pgdb)
}
