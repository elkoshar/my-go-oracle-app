package member

import (
	"encoding/json"

	entity "oracle.com/oracle/my-go-oracle-app/service"
)

type Member struct {
	Name string `db:"NAME"`
	Info string `db:"INFO"`
	entity.BaseEntity
}

type MemberResponse struct {
	Id   int64      `json:"id"`
	Name string     `json:"name"`
	Info MemberInfo `json:"info"`
}

type MemberInfo struct {
	Address string `json:"address"`
	Salary  int    `json:"salary"`
	Age     int    `json:"age"`
}

func (c *Member) ToResponse() MemberResponse {

	var info MemberInfo
	json.Unmarshal([]byte(c.Info), &info)

	return MemberResponse{
		Id:   c.BaseEntity.Id,
		Name: c.Name,
		Info: info,
	}
}
