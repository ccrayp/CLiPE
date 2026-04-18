package rule

import (
	"clipe/pkg/database"
	"encoding/json"
)

type RuleRepository struct {
	db_ *database.DB
}

func NewRuleRep(db *database.DB) *RuleRepository {
	return &RuleRepository{db_: db}
}

func (r *RuleRepository) Select(filter *RuleDTO, limit int, offset int) ([]RuleDTO, error) {

	var rules []Rule

	query := r.db_.Conn().
		Limit(limit).
		Offset(offset)

	if filter.RuleID != 0 {
		query = query.Where("rule_id = ?", filter.RuleID)
	}
	if filter.RuleName != "" {
		query = query.Where("rule_name = ?", filter.RuleName)
	}
	query = query.Where("effect = ?", filter.Effect)

	if err := query.Find(&rules).Error; err != nil {
		return nil, err
	}

	var result []RuleDTO
	for _, r := range rules {
		result = append(result, ToDTO(r))
	}

	return result, nil
}

func (r *RuleRepository) Create(dto *CreateRuleDTO) (*uint, error) {

	model := FromCreateDTO(*dto)

	if err := r.db_.Conn().Create(&model).Error; err != nil {
		return nil, err
	}

	return &model.RuleID, nil
}

func (r *RuleRepository) Update(id uint, dto *CreateRuleDTO) error {

	var model Rule

	if err := r.db_.Conn().First(&model, id).Error; err != nil {
		return err
	}

	condBytes, err := json.Marshal(dto.Condition)
	if err != nil {
		return err
	}

	model.RuleName = dto.RuleName
	model.Condition = condBytes
	model.Effect = dto.Effect

	return r.db_.Conn().Save(&model).Error
}

func (r *RuleRepository) Delete(id uint) error {
	return r.db_.Conn().Unscoped().Delete(&Rule{}, id).Error
}
