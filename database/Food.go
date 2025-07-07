package database

import "fmt"

var sinceVersion_food = 1

type FoodEntity struct {
	ID              int     `db:"id"`
	Name            string  `db:"name"`
	Description     string  `db:"description"`
	UnitId          int     `db:"unit_id"`
	Quantity        float32 `db:"quantity"`
	Calories        float32 `db:"calories"`
	Protein         float32 `db:"protein"`
	Fat             float32 `db:"fat"`
	NonSaturatedFat float32 `db:"non_saturated_fat"`
	Carbohydrates   float32 `db:"carbohydrates"`
	Sugar           float32 `db:"sugar"`
	Source          string  `db:"source"`
}

func (FoodEntity) InitTable(targetVersion int) {
	if targetVersion < sinceVersion_food {
		fmt.Println("Target version is lower than the minimum required version for food table.")
		return
	}
	if targetVersion == sinceVersion_food {
		_, err := DB.Exec(`
			CREATE TABLE IF NOT EXISTS food (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				description TEXT,
				unit_id INTEGER,
				quantity REAL DEFAULT 1.0,
				calories REAL DEFAULT 0.0,
				protein REAL DEFAULT 0.0,
				fat REAL DEFAULT 0.0,
				non_saturated_fat REAL DEFAULT 0.0,
				carbohydrates REAL DEFAULT 0.0,
				sugar REAL DEFAULT 0.0,
				source TEXT,
				FOREIGN KEY (unit_id) REFERENCES unit(id)
			);
		`)
		if err != nil {
			fmt.Println("Error creating food table: " + err.Error())
			return
		}
		fmt.Println("Food table initialized successfully.")
	}
}

func InsertFood(food FoodEntity) (*FoodEntity, error) {
	if food.ID != 0 {
		return nil, fmt.Errorf("food ID must be zero for insertion, got %d", food.ID)
	}

	result, err := DB.Exec(`
		INSERT INTO food (name, description, unit_id, quantity, calories, protein, fat, non_saturated_fat, carbohydrates, sugar, source) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		food.Name, food.Description, food.UnitId, food.Quantity,
		food.Calories, food.Protein, food.Fat,
		food.NonSaturatedFat, food.Carbohydrates,
		food.Sugar, food.Source)
	if err != nil {
		return nil, fmt.Errorf("error inserting food: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	var ret FoodEntity = food
	ret.ID = int(id)

	return &ret, nil
}
