package member

import (
	"encoding/json"

	entity "oracle.com/oracle/my-go-oracle-app/service"
)

// todo:
// buat contoh untuk type data lain
// integer, booleah (cahar(1)), date/timestamp

type Member struct {
	Name   string          `db:"NAME"`
	Info   string          `db:"INFO"`
	Detail json.RawMessage `db:"DETAIL"`
	Policy string          `db:"POLICY"`
	entity.BaseEntity
}

type MemberDetail struct {
	Category string `json:"category"`
	Level    int64  `json:"level"`
}

type MemberRequest struct {
	Name string     `json:"name"`
	Info MemberInfo `json:"info"`
}

type Policy struct {
	EmergencyContact string `json:"emergencyContact"`
	Status           string `json:"status"`
}

type MemberResponse struct {
	Id     int64        `json:"id"`
	Name   string       `json:"name"`
	Info   MemberInfo   `json:"info"`
	Detail MemberDetail `json:"detail"`
	Policy Policy       `json:"policy"`
}

type MemberInfo struct {
	Address string `json:"address"`
	Salary  int    `json:"salary"`
	Age     int    `json:"age"`
}

func (m *Member) ToResponse() MemberResponse {

	var info MemberInfo
	var detail MemberDetail
	var policy Policy
	json.Unmarshal([]byte(m.Info), &info)
	json.Unmarshal(m.Detail, &detail)
	json.Unmarshal([]byte(m.Policy), &policy)

	return MemberResponse{
		Id:     m.BaseEntity.Id,
		Name:   m.Name,
		Info:   info,
		Detail: detail,
		Policy: policy,
	}
}
func (m *MemberRequest) ToEntity() Member {
	infoBytes, _ := json.Marshal(m.Info)

	return Member{
		Name: m.Name,
		Info: string(infoBytes),
	}
}
