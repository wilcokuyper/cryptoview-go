package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/wilcokuyper/cryptoview-go/marketdata"
	"go.uber.org/zap"
)

var logger *zap.Logger
var tmpl *template.Template

func main() {
	godotenv.Load();

	var err error
	if env, ok := os.LookupEnv("APP_ENV"); ok && env == "development" {
		logger, err =  zap.NewDevelopment()
	} else {
		logger, err =  zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("unable to create zap logger")
	}
	defer logger.Sync()

	// Parse templates
	tmpl = template.Must(template.ParseGlob("./templates/*.tmpl"))

	mux := http.NewServeMux()

	// Setup static file server
	fileServer := http.FileServer(http.Dir("./static/public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Setup main handler
	mux.HandleFunc("/", mainHandler)

	client := marketdata.NewCryptocompareClient(
		os.Getenv("CRYPTOCOMPARE_API_KEY"),
		os.Getenv("CRYPTOCOMPARE_BASE_URL"),
		logger,
	)

	server := NewServer(logger, mux, client)

	server.run()
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	tmpl.ExecuteTemplate(w, "index.tmpl", struct{Title string} {Title: "Cryptoview - Manage your crypto assets"})
}