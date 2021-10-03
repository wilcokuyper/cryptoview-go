package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/wilcokuyper/cryptoview-go/marketdata"
)

var tmpl *template.Template

func main() {
	godotenv.Load();

	// Parse templates
	tmpl = template.Must(template.ParseGlob("./templates/*.tmpl"))

	mux := http.NewServeMux()

	// Setup static file server
	fileServer := http.FileServer(http.Dir("./static/public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Setup main handler
	mux.HandleFunc("/", mainHandler)

	// Setup API routes
	setupAPIRoutes(mux)

	// Lookup port and start server
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	log.Println("Starting webserver. Listening on :" + port)

	err := http.ListenAndServe(":" + port, mux)
	if err != nil {
		log.Fatal(err)
	}
}	

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	tmpl.ExecuteTemplate(w, "index.tmpl", struct{Title string} {Title: "Cryptoview - Manage your crypto assets"})
}

func setupAPIRoutes(mux *http.ServeMux) {
	client := marketdata.NewCryptocompareClient(os.Getenv("CRYPTOCOMPARE_API_KEY"), os.Getenv("CRYPTOCOMPARE_BASE_URL"))
	server := &Server{client: client}

	mux.HandleFunc("/api/price", server.GetPriceHandler)
	mux.HandleFunc("/api/symbols", server.GetSymbolsHandler)
}