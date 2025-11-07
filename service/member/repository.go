package member

import (
	"context"
	"fmt"
	"log/slog"

	service "oracle.com/oracle/my-go-oracle-app/service"
)

var (
	id   = "id"
	name = "name"
	info = "info"

	allColumns = []string{id, name, info}
)

const (
	FAILED_FETCH_DATA_ERR_MSG = "failed fetch data. err=%v"
)

type memberRepository struct {
	service.BaseRepository
}

type MemberRepository interface {
	FindById(ctx context.Context, ID int64) (Member, error)
	//	CreateMember(ctx context.Context, data *Member) (int64, error)
	//	GetAllMembers(ctx context.Context, param service.SqlParameter) ([]Member, error)

	service.BaseRepositoryInterface
}

func NewMemberRepository(baseRepository service.BaseRepository) MemberRepository {
	return &memberRepository{
		baseRepository,
	}

}

func (cr *memberRepository) FindById(ctx context.Context, ID int64) (member Member, err error) {
	err = cr.GetOperations(ctx, &member, findByIdQuery, ID)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to fetch data: %v", err), slog.String("query", findByIdQuery), slog.Int64("ID", ID))
		return
	}
	return
}
