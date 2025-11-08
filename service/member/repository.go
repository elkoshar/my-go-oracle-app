package member

import (
	"context"
	"fmt"
	"log/slog"

	"oracle.com/oracle/my-go-oracle-app/pkg/constants"
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
	tableName                 = "MEMBER"
)

type memberRepository struct {
	service.BaseRepository
}

type MemberRepository interface {
	FindById(ctx context.Context, ID int64) (Member, error)
	GetAllMembers(ctx context.Context, param service.SqlParameter) ([]Member, error)
	CountAll(ctx context.Context, params service.SqlParameter) (int64, error)
	//	CreateMember(ctx context.Context, data *Member) (int64, error)

	service.BaseRepositoryInterface
}

func NewMemberRepository(baseRepository service.BaseRepository) MemberRepository {
	return &memberRepository{
		baseRepository,
	}

}

func (mr *memberRepository) FindById(ctx context.Context, ID int64) (member Member, err error) {
	err = mr.GetOperations(ctx, &member, findByIdQuery, ID)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to fetch data: %v", err), slog.String("query", findByIdQuery), slog.Int64("ID", ID))
		return
	}
	return
}

func (mr *memberRepository) GetAllMembers(ctx context.Context, param service.SqlParameter) (members []Member, err error) {
	query, args := mr.GenerateQuerySelectWithParams(getAllMemberQuery, param)
	fmt.Printf("\n query %s \n", query)
	fmt.Printf("\n args %v \n", args)

	err = mr.SelectOperations(ctx, &members, query, args...)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to fetch data: %v", err), slog.String("query", query), slog.Any("args", args))
		return
	}
	fmt.Printf("\n members %v \n", members)

	return
}

func (mr *memberRepository) CountAll(ctx context.Context, params service.SqlParameter) (count int64, err error) {
	params.TableName = fmt.Sprintf("%s m", tableName)
	params.Columns = []string{constants.COUNT_COL}
	params.Limit = 0
	params.Offset = 0
	var resp int64
	err = mr.GetWithParameter(ctx, &resp, params)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf(FAILED_FETCH_DATA_ERR_MSG, err))
		return 0, err
	}

	return resp, nil

}
