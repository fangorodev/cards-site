package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"      // HTTP router
	_ "github.com/mattn/go-sqlite3" // SQLite driver

	_ "github.com/GoAdminGroup/go-admin/adapter/gin"  // Gin adapter (registers itself)
	engine "github.com/GoAdminGroup/go-admin/engine"         // GoAdmin engine
	config "github.com/GoAdminGroup/go-admin/modules/config" // GoAdmin config
	db "github.com/GoAdminGroup/go-admin/modules/db" // Database reference
	adminPkg "github.com/GoAdminGroup/go-admin/plugins/admin"  // Admin plugin
	table "github.com/GoAdminGroup/go-admin/plugins/admin/modules/table" // Table setup
	context "github.com/GoAdminGroup/go-admin/context"
	_ "github.com/GoAdminGroup/themes/sword" // Sword theme
)

// runMigrations ensures the session table exists.
func runMigrations(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Adjustments: CSRF tokens require it to be key-value
	// const sqlStmt = `
	// CREATE TABLE IF NOT EXISTS goadmin_session (
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	session_key TEXT NOT NULL UNIQUE,
	// 	session_value TEXT NOT NULL,
	// 	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	// 	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	// );`

	const sqlStmt = `
	CREATE TABLE IF NOT EXISTS goadmin_session (
		id		INTEGER		PRIMARY KEY AUTOINCREMENT,
		key		TEXT		NOT NULL UNIQUE,
		"values"	TEXT		NOT NULL,
		created_at	DATETIME	DEFAULT CURRENT_TIMESTAMP,
		updated_at	DATETIME	DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(sqlStmt); err != nil {
		return err
	}

	return nil
}

func main() {
	const dbFile = "./database.db"
	// Ensure uploads directory exists for file‐uploads in GoAdmin.
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		log.Fatalf("could not create uploads dir: %v", err)
	}

	// Run migrations before starting the server.
	if err := runMigrations(dbFile); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	// Initialize Gin engine.
	r := gin.Default()

	// Initialize GoAdmin.
	eng := engine.Default()
	cfg := config.Config{
		Theme:     "sword",
		Databases: map[string]config.Database{
			"default": {
				Driver: config.DriverSqlite,
				File:   dbFile, // file‐based SQLite
			},
		},
		UrlPrefix: "admin", // access via /admin
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "/uploads",
		},
		Language: "en",
		IndexUrl: "/",
	}

	// Admin object for plugin handling
	adminPlugin := adminPkg.NewAdmin()

	// Register a default generator so it triggers full setup

	adminPlugin.AddGenerator("user", func(ctx *context.Context) table.Table {
		// This creates a default CRUD table for you... NOT ANYMORE!
		// return table.NewDefaultTable(table.DefaultConfigWithDriver("sqlite"))
		// This was removed at some point in time (1.2.x),
		// Now you get to do it yourself, HAVE FUN

		// 1) Define a table.Config for SQLite
		tblCfg := table.Config{
			// which driver to use (must match your config.Databases["default"].Driver)
			Driver:     config.DriverSqlite,
			// the connection name; for a single‐DB app this is always DefaultConnectionName
			Connection: table.DefaultConnectionName,
			// primary key settings:
			PrimaryKey: table.PrimaryKey{
				Type: db.Int,
				Name: table.DefaultPrimaryKeyName,
			},
			// CRUD capabilities
			CanAdd:     true,
			Editable:   true,
			Deletable:  true,
			Exportable: true,
		}

		// 2) Create the default CRUD table with that config
		userTable := table.NewDefaultTable(ctx, tblCfg)

		// 3) (Optional) customize which columns to show
		info := userTable.GetInfo()
		info.
		AddField("ID", "id", db.Int).FieldSortable().
		AddField("Session Key", "session_key", db.Varchar).
		AddField("Value", "session_value", db.Varchar)

		return userTable
	}) 

	// Wire up GoAdmin with Gin.
	if err := eng.AddConfig(&cfg).
	AddPlugins(adminPlugin).
	Use(r); err != nil {
		log.Fatalf("failed to initialize GoAdmin: %v", err)
	}

	// A simple root handler.
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello, GoAdmin on Gin!")
	})

	log.Println("Listening on :8080 ...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
