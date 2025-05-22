package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"

	"github.com/GoAdminGroup/go-admin/engine"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/plugins/admin"
	_ "github.com/GoAdminGroup/themes/sword"

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"
)

func main() {
	// Connect to Postgres
	db, err := sql.Open("postgres", "postgres://postgres:example@db:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal("DB-Open: ", err)
	}
	defer db.Close() // Late return

	// Run migrations... in docker, not here

	// GoAdmin Engine
	eng := engine.Default()
	cfg := config.Config{
		Theme:	"sword",
		Databases: map[string]config.Database{
			"default": {
				Host:       "db",
				Port:       "5432",
				User:       "postgres",
				Pwd:        "example",
				Name:       "postgres",
				MaxIdleConns: 50,
				MaxOpenConns: 150,
				Driver:     config.DriverPostgresql,
			},
		},
		UrlPrefix: "admin",        // your dashboard at /admin
		Store:     config.Store{    // file upload settings
			Path:   "./uploads",
			Prefix: "/uploads",
		},
		Language: "en",
		IndexUrl: "/",
	}

	// Add Gin to the Tonic
	r := gin.Default()

	if err:= eng.AddConfig(&cfg).AddPlugins(
		admin.NewAdmin(),
	).Use(r); err != nil {
		panic(err)
	}



	// Serve... (routing)
	r.GET("/", func (ctx * gin.Context) {
		ctx.String(200, "Hello World!")
	})

	// Start Server
	log.Println("Listen on :8080 ...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
