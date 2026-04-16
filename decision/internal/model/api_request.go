package model

import "time"

type ApiRequest struct {
	User struct {
		Name   string   `json:"username"`
		UID    uint     `json:"uid"`
		GID    uint     `json:"gid"`
		Groups []string `json:"groups"`
	} `json:"user"`

	Host struct {
		IP       string `json:"ip"`
		HostName string `json:"hostname"`
	} `json:"host"`

	Service string `json:"service"`
	Action  string `json:"action"`

	Time struct {
		Timestamp time.Time `json:"timestamp"`
		Weekday   string    `json:"weekday"`
	} `json:"time"`
}

func GetValue(req *ApiRequest, type_ string) any {
	switch type_ {

	case "username":
		return req.User.Name

	case "uid":
		return req.User.UID

	case "gid":
		return req.User.GID

	case "groups":
		return req.User.Groups

	case "ip":
		return req.Host.IP

	case "hostname":
		return req.Host.HostName

	case "service":
		return req.Service

	case "action":
		return req.Action

	case "timestamp":
		return req.Time.Timestamp

	case "weekday":
		return req.Time.Timestamp

	}

	return nil
}

//
// {
//   "user": {
//     "name": "roman",
//     "uid": 1001,
//     "gid": 1001,
//     "groups": ["sudo"]
//   },
//   "host": {
//     "ip": "192.168.1.10",
//     "hostname": "prod-1"
//   },
//   "service": "sshd",
//   "action": "auth",
//   "time": {
//     "timestamp": "2026-04-15T12:00:00Z",
//     "weekday": "tue"
//   }
// }
//
