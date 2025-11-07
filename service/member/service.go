package member

import (
	"context"
	"fmt"
	"log/slog"
)

type memberService struct {
	mr MemberRepository
}

type MemberService interface {
	FindById(ctx context.Context, id int64) (MemberResponse, error)
}

func NewMemberService(mr MemberRepository) MemberService {
	return &memberService{mr: mr}
}

func (c *memberService) FindById(ctx context.Context, id int64) (MemberResponse, error) {
	var (
		response MemberResponse
		member   Member
	)

	member, err := c.mr.FindById(ctx, id)
	if err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("Failed to get member data: %v", err), slog.Int64("id", id))
		return response, err
	}

	response = member.ToResponse()
	return response, nil
}
