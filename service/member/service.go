package member

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	service "oracle.com/oracle/my-go-oracle-app/service"
)

type memberService struct {
	mr MemberRepository
}

type MemberService interface {
	FindById(ctx context.Context, id int64) (MemberResponse, error)
	FindAll(ctx context.Context, param service.SqlParameter) ([]MemberResponse, service.Pagination, error)
	CreateMember(ctx context.Context, data *MemberRequest) (MemberResponse, error)
	UpdateMember(ctx context.Context, id int64, data *MemberRequest) (MemberResponse, error)
	DeleteMember(ctx context.Context, id int64) (bool, error)
}

func NewMemberService(mr MemberRepository) MemberService {
	return &memberService{mr: mr}
}

func (m *memberService) FindById(ctx context.Context, id int64) (MemberResponse, error) {
	var (
		response MemberResponse
		member   Member
	)

	member, err := m.mr.FindById(ctx, id)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("Failed to get member data: %v", err), slog.Int64("id", id))
		return response, err
	}

	response = member.ToResponse()
	return response, nil
}

func (m *memberService) FindAll(ctx context.Context, param service.SqlParameter) (memberResponse []MemberResponse, page service.Pagination, err error) {
	memberEntities, err := m.mr.GetAllMembers(ctx, param)

	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("Failed to get member data: %v", err))
		return
	}
	result := make([]MemberResponse, len(memberEntities))
	for idx, data := range memberEntities {
		result[idx] = data.ToResponse()
	}

	count, err := m.mr.CountAll(ctx, param)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("Failed to find count total member. err =%v\n", err))
		return
	}
	len := len(memberEntities)
	page = service.MakePagination(count, param, len)

	return result, page, nil
}

func (m *memberService) CreateMember(ctx context.Context, data *MemberRequest) (MemberResponse, error) {

	var (
		response MemberResponse
	)

	baseEntity := service.BaseEntity{
		CreatedDate: time.Now(),
		IsDeleted:   "0",
	}

	member := data.ToEntity(baseEntity)

	id, err := m.mr.CreateMember(ctx, &member)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed create member = %v, err = %v", data, err))
		return response, fmt.Errorf("err:%s", err.Error())
	}

	response = member.ToResponse()
	response.Id = id

	return response, nil

}
func (m *memberService) UpdateMember(ctx context.Context, id int64, data *MemberRequest) (MemberResponse, error) {

	var (
		response MemberResponse
	)

	baseEntity := service.BaseEntity{
		UpdatedDate: sql.NullTime{Time: time.Now(), Valid: true},
		IsDeleted:   "0",
	}

	member := data.ToEntity(baseEntity)

	// Set the member's ID since we're updating
	member.Id = id

	_, err := m.mr.UpdateMember(ctx, id, &member)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed update member = %v, err = %v", data, err))
		return response, fmt.Errorf("err:%s", err.Error())
	}

	response = member.ToResponse()
	response.Id = id

	return response, nil

}

func (m *memberService) DeleteMember(ctx context.Context, id int64) (bool, error) {

	_, err := m.mr.DeleteMember(ctx, id)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("failed delete member id = %v, err = %v", id, err))
		return false, fmt.Errorf("err:%s", err.Error())
	}

	return true, nil
}
