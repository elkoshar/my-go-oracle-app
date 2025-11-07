package api

import (
	"context"

	"oracle.com/oracle/my-go-oracle-app/service/member"
)

type MemberService interface {
	FindById(ctx context.Context, id int64) (member.MemberResponse, error)
}
