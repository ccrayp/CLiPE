package request

import "encoding/json"

func ToDTO(r Request) RequestDTO {
	var ctx interface{}
	_ = json.Unmarshal(r.Context, &ctx)

	return RequestDTO{
		RequestID: r.RequestID,
		UserID:    r.UserID,
		Context:   ctx,
		Timestamp: r.Timestamp,
	}
}

func FromCreateDTO(dto CreateRequestDTO) Request {
	condBytes, _ := json.Marshal(dto.Context)

	return Request{
		UserID:    dto.UserID,
		Context:   condBytes,
		Timestamp: dto.Timestamp,
	}
}
