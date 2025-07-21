package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"

	"cards-site/internal/routes"
)

// App config
type Config struct {
	Port	int	`yaml:"port"`
	DSN	string	`yaml:"dsn"`
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var conf Config
	decoder := yaml.NewDecoder(file)
	
	if err := decoder.Decode(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func main() {
	// Parse flag/s
	cfgPath := flag.String("config", "cmd/web/config.yaml", "Path to config file")
	flag.Parse()

	// Load config
	conf, err := loadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("Failure at config load: %v", err)
	}

	// DB Connection
	db, err := sqlx.Connect("postgres", conf.DSN) 
	if err != nil {
		log.Fatalf("Failure at DB connection: %v", err)
	}

	defer db.Close()
	// Note: Ensure migrations have been run through containers...

	// Routes + Server
	router := routes.NewRouter(db)

	addr := fmt.Sprintf(":%d", conf.Port)
	log.Printf("Server On %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failure at Server link: %v", err)
	}
}
