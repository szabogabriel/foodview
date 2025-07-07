package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite" // Import the SQLite driver
)

var DB *sql.DB

const SCHEMA_META_TABLE_NAME = "schema_meta"
const CURRENT_VERSION = 1

type dbtable interface {
	InitTable(db *sql.DB, targetVersion int)
}

var DB_TABLES = []dbtable{
	Unit{},
}

func Connect() {
	createConnection()

	initDatabase()

	updateDatabase()
}

func tableExists(name string) (bool, error) {
	var foundName string

	err := DB.QueryRow(`
        SELECT name FROM sqlite_master 
        WHERE type='table' AND name=?`, name).Scan(&foundName)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func createSchemaVersionTable() error {
	var err error
	_, err = DB.Exec("CREATE TABLE " + SCHEMA_META_TABLE_NAME + " (version INTEGER PRIMARY KEY, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		fmt.Println("Error creating schema meta table:", err)
		return err
	}
	return nil
}

func getSchemaVersion() (int, error) {
	var version int
	err := DB.QueryRow("SELECT version FROM " + SCHEMA_META_TABLE_NAME + " ORDER BY version DESC LIMIT 1").Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil // Not initialized
	}
	return version, err
}

func setSchemaVersion(version int) error {
	_, err := DB.Exec("INSERT INTO "+SCHEMA_META_TABLE_NAME+" (version) VALUES (?)", version)
	if err != nil {
		return fmt.Errorf("failed to set schema version: %w", err)
	}
	return nil
}

func createConnection() error {
	var sqliteFile = os.Getenv("DB_SQLITE_FILE")
	fmt.Println("Connecting to database file", sqliteFile)
	var err error
	DB, err = sql.Open("sqlite", sqliteFile)
	if err != nil {
		fmt.Println("Error opening the database", err)
		return err
	}
	return nil
}

func initDatabase() {
	exists, err := tableExists(SCHEMA_META_TABLE_NAME)
	if err != nil {
		fmt.Println("Error checking schema meta table:", err)
		return
	}
	if !exists {
		createSchemaVersionTable()
	}
}

func updateDatabase() {
	version, err := getSchemaVersion()
	if err != nil {
		fmt.Println("Error getting schema version:", err)
		return
	}

	if version < CURRENT_VERSION {
		fmt.Println("Updating database schema to version", CURRENT_VERSION)
		for _, table := range DB_TABLES {
			for i := version + 1; i <= CURRENT_VERSION; i++ {
				table.InitTable(DB, i)
			}
		}

	}
}
