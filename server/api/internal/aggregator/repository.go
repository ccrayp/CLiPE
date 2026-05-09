package aggregator

import (
	"clipe/internal/policy"
	policycontent "clipe/internal/policy_content"
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

func (r *AggreagtorRepository) FindPolicy(userId uint) (*policy.Policy, error) {
	filter := policy.PolicyDTO{
		UserID: userId,
		Status: &[]bool{true}[0],
	}

	record := policy.Policy{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *AggreagtorRepository) FindService(sericeName string) (*service.Service, error) {
	filter := service.ServiceDTO{
		ServiceName: sericeName,
	}

	record := service.Service{}

	err := r.db.Conn().Where(filter).Find(&record).Error
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *AggreagtorRepository) FindRule(serviceId, policyId uint) (*rule.Rule, error) {
	filterPolicyContent := policycontent.PolicyContentDTO{
		ServiceID: serviceId,
		PolicyID:  policyId,
	}

	recordPolicyContent := policycontent.PolicyContent{}

	err := r.db.Conn().
		Where(filterPolicyContent).
		First(&recordPolicyContent).Error
	if err != nil {
		return nil, err
	}

	recordRule := rule.Rule{}

	err = r.db.Conn().
		First(&recordRule, recordPolicyContent.RuleID).Error
	if err != nil {
		return nil, err
	}

	return &recordRule, nil
}

func (r *AggreagtorRepository) CreateService(serviceName string) error {
	service := service.Service{
		ServiceName: serviceName,
	}

	return r.db.Conn().Create(&service).Error
}
