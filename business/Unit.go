package business

import (
	"fmt"

	"github.com/szabogabriel/foodview/database"
)

func ToBasicUnit(unit database.UnitEntity, value float32) (*database.UnitEntity, float32, error) {
	if unit.ParentId == 0 && unit.QuotientToParent == 0.0 {
		return &unit, value, nil
	}

	var retUnit *database.UnitEntity = &unit
	var err error
	var ret float32 = value

	for unit.ParentId != 0 && unit.QuotientToParent != 0.0 {
		ret *= float32(retUnit.QuotientToParent)
		retUnit, err = database.GetUnitByID(retUnit.ParentId)
		if err != nil {
			fmt.Println("Error getting parent unit:", err)
			return nil, 0.0, err
		}
	}

	return retUnit, ret, nil
}
