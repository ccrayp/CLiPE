package aggregator

import (
	"clipe/internal/action"
	"clipe/internal/host"
	"clipe/internal/policy"
	"clipe/internal/rule"
	"clipe/internal/service"
	"clipe/internal/user"
	"clipe/pkg/database"
)

type AggreagtorRepository struct {
	db *database.DB
}

func NewAggregatorRepository(db *database.DB) *AggreagtorRepository {
	return &AggreagtorRepository{
		db: db,
	}
}

func (r *AggreagtorRepository) FindUserIdByName(userName string) (*uint, error) {
	filter := user.UserDTO{
		UserName: userName,
	}

	record := user.User{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record.UserID, nil
}

func (r *AggreagtorRepository) FinHostIdByIp(hostIp string) (*uint, error) {
	filter := host.HostDTO{
		IP: hostIp,
	}

	record := host.Host{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record.HostID, nil
}

func (r *AggreagtorRepository) FindServiceIdByName(serviceName string) (*uint, error) {
	filter := service.ServiceDTO{
		ServiceName: serviceName,
	}

	record := service.Service{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record.ServiceID, nil
}

func (r *AggreagtorRepository) FindActionIdByName(actionName string) (*uint, error) {
	filter := action.ActionDTO{
		ActionName: actionName,
	}

	record := action.Action{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record.ActionID, nil
}

func (r *AggreagtorRepository) FindPolicy(userId, hostId, serviceId, actionId uint) (*policy.Policy, error) {
	filter := policy.PolicyDTO{
		UserID:    &userId,
		HostID:    &hostId,
		ServiceID: &serviceId,
		ActionID:  &actionId,
	}

	record := policy.Policy{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *AggreagtorRepository) FindRuleById(ruleId uint) (*rule.Rule, error) {
	var record rule.Rule

	err := r.db.Conn().
		Where("rule_id = ?", ruleId).
		First(&record).Error

	if err != nil {
		return nil, err
	}

	return &record, nil
}
