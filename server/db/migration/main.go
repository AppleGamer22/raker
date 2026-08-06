package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/AppleGamer22/raker/server/db"
	old "github.com/AppleGamer22/raker/server/db/mongo"
	"github.com/AppleGamer22/raker/server/handlers"
	"github.com/charmbracelet/log"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
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
		_, err := pgdb.UserAdd(ctx, db.UserAddParams{
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
			Post:     history.Post,
			Username: "",
			PostType: db.PostType(history.Post),
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

func coordinates(ctx context.Context, username, storageRoot string, pgdb *db.Queries) {
	storage := rakerServer.NewStorageHandler(rakerServer.Configuration.Storage, rakerServer.Configuration.Directories)

	user, err := pgdb.UserGetByUsername(ctx, username)
	if err != nil {
		log.Fatal(err)
	}

	count, err := pgdb.HistoryCount(ctx, db.HistoryCountParams{
		PostTypes:      []db.PostType{db.PostTypeVsco},
		Categories:     user.Categories,
		UserCategories: user.Categories,
		Username:       user.Username,
	})
	if err != nil {
		log.Fatal(err)
	}

	pageSize := 50.0
	pages := int(math.Ceil(float64(count) / pageSize))

	for page := 0; page < pages; page++ {
		histories, err := pgdb.HistoryGetPage(ctx, db.HistoryGetPageParams{
			PostTypes:      []db.PostType{db.PostTypeVsco},
			Categories:     user.Categories,
			UserCategories: user.Categories,
			Username:       user.Username,
			Page:           int32(page),
			PageSize:       int32(pageSize),
		})
		if err != nil {
			log.Fatal(err, "page", page, "pageSize", pageSize)
		}
		for _, history := range histories {
			latitude, longitude := storage.LocationEXIF(user, history.PostType, history.PostOwner, history.Files[0])
			if latitude == 0 && longitude == 0 {
				continue
			} else if history.Coordinates.Valid {
				log.Debug(
					"skipping",
					"page", page,
					"pageSize", pageSize,
					"history", fmt.Sprintf("%s/%s/%s", history.PostType, history.PostOwner, history.Post),
					"coordinates", history.Coordinates.Point.String(),
				)
			}

			err := pgdb.HistoryUpdateCoordinates(ctx, db.HistoryUpdateCoordinatesParams{
				Latitude:  latitude,
				Longitude: longitude,
				PostType:  history.PostType,
				PostOwner: history.PostOwner,
				Username:  history.Username,
			})

			if err != nil {
				log.Fatal(err, "page", page, "pageSize", pageSize, "history", history)
			}
		}
	}

}

var rakerServer handlers.RakerServer

func main() {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&rakerServer.Configuration); err != nil {
		log.Fatal(err)
	}

	log.Debug(rakerServer.Configuration.Database)

	ctx := context.Background()
	connection, err := sql.Open("postgres", rakerServer.URI)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	// if err := connection.Ping(); err != nil {
	// 	log.Fatal(err)
	// }

	pgdb := db.New(connection)

	// users(ctx, pgdb)
	// histories(ctx, pgdb)
	coordinates(ctx, rakerServer.Configuration.Database, rakerServer.Configuration.Storage, pgdb)
}
