package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"

	"github.com/droptheplot/rand_chat/config"
	"github.com/droptheplot/rand_chat/room"
	"github.com/droptheplot/rand_chat/telegram"
)

var DB *gorm.DB
var templates = template.Must(template.ParseGlob("templates/*.html"))

func init() {
	var err error
	DB, err = gorm.Open("postgres", config.Store.Database)

	if err != nil {
		panic(err)
	}

	DB.LogMode(true)

	room.DB = DB
}

func main() {
	driver, _ := postgres.WithInstance(DB.DB(), &postgres.Config{})
	migrations, _ := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	migrations.Up()

	telegram.SetWebhook()

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/api/update", UpdateHandler).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", map[string]string{"Title": "qwe"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var update telegram.Update
	json.NewDecoder(r.Body).Decode(&update)

	switch update.Message.Text {
	case "/start":
		room.Join(update.Message.Chat.ID)
	case "/stop":
		room.Stop(update.Message.Chat.ID)
	default:
		_, targetID := room.Find(update.Message.Chat.ID)

		telegram.SendMessage(targetID, update.Message.Text)
	}

	w.WriteHeader(http.StatusOK)
}