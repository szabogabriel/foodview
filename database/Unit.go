package database

import (
	"database/sql"
	"fmt"
)

var sinceVersion_unit = 1

type Unit struct {
	ID               int     `db:"id"`
	Name             string  `db:"name"`
	ParentId         int     `db:"parent_id"`
	QuotientToParent float64 `db:"quotient_to_parent"`
}

func (Unit) InitTable(db *sql.DB, targetVersion int) {
	if targetVersion < sinceVersion_unit {
		fmt.Println("Target version is lower than the minimum required version.")
		return
	}
	if targetVersion == sinceVersion_unit {
		_, err := db.Exec(`
			CREATE TABLE IF NOT EXISTS unit (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				parent_id INTEGER,
				quotient_to_parent REAL,
				FOREIGN KEY (parent_id) REFERENCES unit(id)
			);
		`)
		if err != nil {
			fmt.Println("Error creating unit table:", err)
			return
		}

		_, err = db.Exec(`
			INSERT INTO unit (name, parent_id, quotient_to_parent) VALUES
			('gram', NULL, 1.0),
			('kilogram', 1, 1000.0),
			('milligram', 1, 0.001),
			('liter', NULL, 1.0),
			('milliliter', 4, 0.001),
			('deciliter', 4, 0.1),
			('centiliter', 4, 0.01),
			('teaspoon', 4, 0.005),
			('tablespoon', 4, 0.015),
			('cup', 4, 0.24),
			('ounce', 1, 28.3495),
			('pound', 1, 453.592),
			('fluid ounce', 4, 0.0295735),
			('gallon', 4, 3.78541),
			('pinch', 5, 0.36),
			('dash', 5, 0.6),
			('pint', 5, 473.176),
			('piece', NULL, 0.0),
			('dozen', 18, 12.0),
			('half dozen', 18, 6.0)`)
		if err != nil {
			fmt.Println("Error inserting into unit table:", err)
			return
		}
		fmt.Println("Unit table initialized successfully.")
	}
}

func InsertUnit(db *sql.DB, unit Unit) (*Unit, error) {
	if unit.ID != 0 {
		return nil, fmt.Errorf("unit ID must be zero for insertion, got %d", unit.ID)
	}

	if unit.Name == "" {
		return nil, fmt.Errorf("unit name cannot be empty")
	}

	foundUnit, err := GetUnitByName(db, unit.Name)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking for existing unit: %w", err)
	}
	if foundUnit != nil {
		return nil, fmt.Errorf("unit with name '%s' already exists", unit.Name)
	}

	result, err := db.Exec(`
		INSERT INTO unit (name, parent_id, quotient_to_parent) 
		VALUES (?, ?, ?)`, unit.Name, unit.ParentId, unit.QuotientToParent)

	if err != nil {
		return nil, fmt.Errorf("error inserting unit: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %w", err)
	}
	ret := unit
	ret.ID = int(id)
	return &ret, nil
}

func GetUnitByID(db *sql.DB, id int) (*Unit, error) {
	unit := &Unit{}

	if id <= 0 {
		return nil, fmt.Errorf("invalid unit ID: %d", id)
	}
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	err := db.QueryRow(`
		SELECT id, name, parent_id, quotient_to_parent 
		FROM unit WHERE id = ?`, id).Scan(&unit.ID, &unit.Name, &unit.ParentId, &unit.QuotientToParent)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

func GetUnitByName(db *sql.DB, name string) (*Unit, error) {
	if name == "" {
		return nil, fmt.Errorf("unit name cannot be empty")
	}

	unit := &Unit{}
	err := db.QueryRow(`
		SELECT id, name, parent_id, quotient_to_parent 
		FROM unit WHERE name = ?`, name).Scan(&unit.ID, &unit.Name, &unit.ParentId, &unit.QuotientToParent)
	if err != nil {
		return nil, err
	}
	return unit, nil
}
