package member

import (
	"context"
	"fmt"
	"log/slog"

	service "oracle.com/oracle/my-go-oracle-app/service"
)

type memberService struct {
	mr MemberRepository
}

type MemberService interface {
	FindById(ctx context.Context, id int64) (MemberResponse, error)
	FindAll(ctx context.Context, param service.SqlParameter) ([]MemberResponse, service.Pagination, error)
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
