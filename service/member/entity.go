package member

import (
	entity "oracle.com/oracle/my-go-oracle-app/service"
)

type Member struct {
	Name string `db:"NAME"`
	Info string `db:"INFO"`
	entity.BaseEntity
}

type MemberResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Info string `json:"info"`
}

func (c *Member) ToResponse() MemberResponse {
	return MemberResponse{
		Id:   c.BaseEntity.Id,
		Name: c.Name,
		Info: c.Info,
	}
}
