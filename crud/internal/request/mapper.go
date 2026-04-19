package request

import "encoding/json"

func ToDTO(r Request) RequestDTO {
	var ctx interface{}
	_ = json.Unmarshal(r.Context, &ctx)

	return RequestDTO{
		RequestID: r.RequestID,
		UserID:    r.UserID,
		Context:   ctx,
	}
}

func FromCreateDTO(dto CreateRequestDTO) Request {
	condBytes, _ := json.Marshal(dto.Context)

	clean := func(v *uint) *uint {
		if v == nil || *v == 0 {
			return nil
		}
		return v
	}

	return Request{
		UserID:  clean(dto.UserID),
		Context: condBytes,
	}
}
