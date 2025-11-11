package member

import (
	"encoding/json"

	entity "oracle.com/oracle/my-go-oracle-app/service"
)

type Member struct {
	Name   string          `db:"NAME"`
	Info   string          `db:"INFO"`
	Detail json.RawMessage `db:"DETAIL"`
	entity.BaseEntity
}

type MemberDetail struct {
	Category string `json:"category"`
	Level    string `json:"level"`
}

type MemberRequest struct {
	Name string     `json:"name"`
	Info MemberInfo `json:"info"`
}

type MemberResponse struct {
	Id     int64        `json:"id"`
	Name   string       `json:"name"`
	Info   MemberInfo   `json:"info"`
	Detail MemberDetail `json:"detail"`
}

type MemberInfo struct {
	Address string `json:"address"`
	Salary  int    `json:"salary"`
	Age     int    `json:"age"`
}

func (m *Member) ToResponse() MemberResponse {

	var info MemberInfo
	var detail MemberDetail
	json.Unmarshal([]byte(m.Info), &info)
	json.Unmarshal(m.Detail, &detail)

	return MemberResponse{
		Id:     m.BaseEntity.Id,
		Name:   m.Name,
		Info:   info,
		Detail: detail,
	}
}
func (m *MemberRequest) ToEntity() Member {
	infoBytes, _ := json.Marshal(m.Info)

	return Member{
		Name: m.Name,
		Info: string(infoBytes),
	}
}
