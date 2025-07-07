package database

import (
	"database/sql"
	"fmt"
)

var sinceVersion_unit = 1

type UnitEntity struct {
	ID               int     `db:"id"`
	Name             string  `db:"name"`
	ParentId         int     `db:"parent_id"`
	QuotientToParent float64 `db:"quotient_to_parent"`
	ShortName        string  `db:"short_name"`
}

func (UnitEntity) InitTable(targetVersion int) {
	if targetVersion < sinceVersion_unit {
		fmt.Println("Target version is lower than the minimum required version.")
		return
	}
	if targetVersion == sinceVersion_unit {
		_, err := DB.Exec(`
			CREATE TABLE IF NOT EXISTS unit (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				parent_id INTEGER,
				quotient_to_parent REAL,
				FOREIGN KEY (parent_id) REFERENCES unit(id),
				short_name TEXT
			);
		`)
		if err != nil {
			fmt.Println("Error creating unit table:", err)
			return
		}

		_, err = DB.Exec(`
			INSERT INTO unit (name, parent_id, quotient_to_parent) VALUES
			('gram', NULL, 1.0, 'g'),
			('kilogram', 1, 1000.0, 'kg'),
			('milligram', 1, 0.001, 'mg),
			('liter', NULL, 1.0, 'l'),
			('milliliter', 4, 0.001, 'ml'),
			('deciliter', 4, 0.1, 'dl'),
			('centiliter', 4, 0.01, 'cl'),
			('teaspoon', 4, 0.005, 'tsp'),
			('tablespoon', 4, 0.015, 'tbsp'),
			('cup', 4, 0.24, 'cup'),
			('ounce', 1, 28.3495, 'oz'),
			('pound', 1, 453.592, 'lb'),
			('fluid ounce', 4, 0.0295735, 'fl oz'),
			('gallon', 4, 3.78541, 'gal'),
			('pinch', 5, 0.36, 'pinch'),
			('dash', 5, 0.6, 'dash'),
			('pint', 5, 473.176, 'pt'),
			('piece', NULL, 0.0, 'pc'),
			('dozen', 18, 12.0, 'dz'),
			('half dozen', 18, 6.0, 'hdz')`)
		if err != nil {
			fmt.Println("Error inserting into unit table:", err)
			return
		}
		fmt.Println("Unit table initialized successfully.")
	}
}

func InsertUnit(unit UnitEntity) (*UnitEntity, error) {
	if unit.ID != 0 {
		return nil, fmt.Errorf("unit ID must be zero for insertion, got %d", unit.ID)
	}

	if unit.Name == "" {
		return nil, fmt.Errorf("unit name cannot be empty")
	}

	foundUnit, err := GetUnitByName(unit.Name)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking for existing unit: %w", err)
	}
	if foundUnit != nil {
		return nil, fmt.Errorf("unit with name '%s' already exists", unit.Name)
	}

	result, err := DB.Exec(`
		INSERT INTO unit (name, parent_id, quotient_to_parent, short_name) 
		VALUES (?, ?, ?, ?)`, unit.Name, unit.ParentId, unit.QuotientToParent, unit.ShortName)

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

func GetUnitByID(id int) (*UnitEntity, error) {
	unit := &UnitEntity{}

	if id <= 0 {
		return nil, fmt.Errorf("invalid unit ID: %d", id)
	}
	if DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	err := DB.QueryRow(`
		SELECT id, name, parent_id, quotient_to_parent, short_name 
		FROM unit WHERE id = ?`, id).Scan(&unit.ID, &unit.Name, &unit.ParentId, &unit.QuotientToParent, &unit.ShortName)
	if err != nil {
		return nil, err
	}
	return unit, nil
}

func GetUnitByName(name string) (*UnitEntity, error) {
	if name == "" {
		return nil, fmt.Errorf("unit name cannot be empty")
	}

	unit := &UnitEntity{}
	err := DB.QueryRow(`
		SELECT id, name, parent_id, quotient_to_parent, short_name 
		FROM unit WHERE name = ?`, name).Scan(&unit.ID, &unit.Name, &unit.ParentId, &unit.QuotientToParent, &unit.ShortName)
	if err != nil {
		return nil, err
	}
	return unit, nil
}
