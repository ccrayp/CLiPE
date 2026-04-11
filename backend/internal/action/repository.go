package action

import "clipe/pkg/database"

type ActionRep struct {
	db_ *database.DB
}

func NewActionRep(db *database.DB) *ActionRep {
	return &ActionRep{
		db_: db,
	}
}

func (a *ActionRep) GetPagination(limit int, offset int) ([]ActionDTO, error) {
	var actions []Action

	if err := a.db_.Conn().Limit(limit).Offset(offset).Find(&actions).Error; err != nil {
		return nil, err
	}

	var result []ActionDTO
	for _, a := range actions {
		result = append(result, ToDTO(a))
	}

	return result, nil
}
