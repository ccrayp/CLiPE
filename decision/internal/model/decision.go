package model

type Decision struct {
	Result bool
	Policy struct {
		Id   uint
		Name string
	}
	RequestId  uint
	DecisionId uint
}
