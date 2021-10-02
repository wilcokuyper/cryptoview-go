package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/wilcokuyper/cryptoview-go/services"
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
	api := services.NewCryptocompare(os.Getenv("CRYPTOCOMPARE_API_KEY"), os.Getenv("CRYPTOCOMPARE_BASE_URL"))

	mux.HandleFunc("/api/price", func (w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if !r.URL.Query().Has("symbol") {
			w.WriteHeader(http.StatusNotAcceptable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Missing symbol parameter",
			})
			return
		}

		symbol := r.URL.Query().Get("symbol")
		price, err := api.GetPrice(symbol, "EUR")

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Price not found",
			})

			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"price": map[string]interface{}{
				"symbol": symbol,
				"price": price,
			},
		})
	})

	mux.HandleFunc("/api/symbols", func (w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		
		symbols, err := api.GetSymbols()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Unable to retrieve symbols",
			})

			return
		}

		json.NewEncoder(w).Encode(symbols)
	})
}