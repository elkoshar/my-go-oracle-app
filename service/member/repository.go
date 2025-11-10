package member

import (
	"context"
	"database/sql"
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
	CreateMember(ctx context.Context, data *Member) (int64, error)
	UpdateMember(ctx context.Context, id int64, data *Member) (int64, error)
	DeleteMember(ctx context.Context, id int64) (int64, error)
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

	err = mr.SelectOperations(ctx, &members, query, args...)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to fetch data: %v", err), slog.String("query", query), slog.Any("args", args))
		return
	}

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

func (m memberRepository) CreateMember(ctx context.Context, data *Member) (lastInsertId int64, err error) {
	var returnedID int64
	args := []interface{}{data.Name, data.Info, sql.Out{Dest: &returnedID}}
	_, err = m.WriteOrUpdateOperation(ctx, createMemberQuery, &returnedID, args...)

	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to execute query, member = %v, errInsert = %v", data, err))
		return 0, err
	}

	data.Id = returnedID
	return returnedID, nil

}

func (m memberRepository) UpdateMember(ctx context.Context, id int64, data *Member) (rowsAffected int64, err error) {
	args := []interface{}{data.Name, data.Info, id}
	result, errExec := m.WriteOrUpdateOperation(ctx, updateMemberQuery, nil, args...)
	if errExec != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to execute query, member = %v, errExec = %v", data, errExec))
		return 0, errExec
	}

	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to get rows affected, member = %v, err = %v", data, err))
		return 0, err
	}

	return result, nil
}

func (m memberRepository) DeleteMember(ctx context.Context, id int64) (rowsAffected int64, err error) {
	args := []interface{}{id}
	result, errExec := m.WriteOrUpdateOperation(ctx, DeleteMemberQuery, nil, args...)
	if errExec != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to execute query, id = %v, errExec = %v", id, errExec))
		return 0, errExec
	}

	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed to get rows affected, id = %v, err = %v", id, err))
		return 0, err
	}

	return result, nil
}
