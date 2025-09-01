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
	dsnFlag := flag.String("dsn", "", "Overriden DSN")
	flag.Parse()

	// Load config
	conf, err := loadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("Failure at config load: %v", err)
	}

	env := os.Getenv("DB_DSN")

	var dsn string
	switch {
	case *dsnFlag != "":
		dsn = *dsnFlag
	case env != "":
		dsn = env
	case conf.DSN != "":
		dsn = conf.DSN
	default:
		log.Fatal("No DSN provided (set --dsn, DB_DSN, or dsn in config.yaml)")
	}

	// DB Connection
	db, err := sqlx.Connect("postgres", dsn) 
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
