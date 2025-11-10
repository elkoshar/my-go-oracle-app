package member_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"oracle.com/oracle/my-go-oracle-app/service"
	"oracle.com/oracle/my-go-oracle-app/service/member"
)

// MockMemberRepository implements member.MemberRepository for testing
type MockMemberRepository struct {
	mock.Mock
}

func (m *MockMemberRepository) FindById(ctx context.Context, ID int64) (member.Member, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(member.Member), args.Error(1)
}

func (m *MockMemberRepository) GetAllMembers(ctx context.Context, param service.SqlParameter) ([]member.Member, error) {
	args := m.Called(ctx, param)
	return args.Get(0).([]member.Member), args.Error(1)
}

func (m *MockMemberRepository) CountAll(ctx context.Context, params service.SqlParameter) (int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMemberRepository) CreateMember(ctx context.Context, data *member.Member) (int64, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMemberRepository) UpdateMember(ctx context.Context, id int64, data *member.Member) (int64, error) {
	args := m.Called(ctx, id, data)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockMemberRepository) DeleteMember(ctx context.Context, id int64) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func setupTestService() (member.MemberService, *MockMemberRepository) {
	mockRepo := new(MockMemberRepository)
	service := member.NewMemberService(mockRepo)
	return service, mockRepo
}

func TestService_FindById_Success(t *testing.T) {
	// Setup
	svc, mockRepo := setupTestService()
	ctx := context.Background()
	expectedID := int64(1)
	memberEntity := member.Member{
		Name: "Test User",
		Info: `{"address":"Test Address","salary":5000,"age":30}`,
		BaseEntity: service.BaseEntity{
			Id: expectedID,
		},
	}

	// Mock behavior
	mockRepo.On("FindById", ctx, expectedID).Return(memberEntity, nil)

	// Execute
	result, err := svc.FindById(ctx, expectedID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedID, result.Id)
	assert.Equal(t, memberEntity.Name, result.Name)
	assert.Equal(t, "Test Address", result.Info.Address)
	assert.Equal(t, 5000, result.Info.Salary)
	assert.Equal(t, 30, result.Info.Age)
	mockRepo.AssertExpectations(t)
}

func TestService_FindById_NotFound(t *testing.T) {
	// Setup
	svc, mockRepo := setupTestService()
	ctx := context.Background()
	expectedID := int64(999)

	// Mock behavior
	mockRepo.On("FindById", ctx, expectedID).Return(member.Member{}, errors.New("not found"))

	// Execute
	result, err := svc.FindById(ctx, expectedID)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestService_FindAll_Success(t *testing.T) {
	// Setup
	svc, mockRepo := setupTestService()
	ctx := context.Background()
	param := service.SqlParameter{
		TableName: "MEMBER m",
		Columns:   []string{"ID", "NAME", "INFO"},
		Limit:     10,
		Offset:    0,
	}

	members := []member.Member{
		{
			Name: "User 1",
			Info: `{"address":"Address 1","salary":5000,"age":30}`,
			BaseEntity: service.BaseEntity{
				Id: 1,
			},
		},
		{
			Name: "User 2",
			Info: `{"address":"Address 2","salary":6000,"age":35}`,
			BaseEntity: service.BaseEntity{
				Id: 2,
			},
		},
	}

	expectedCount := int64(2)

	// Mock behavior
	mockRepo.On("GetAllMembers", ctx, param).Return(members, nil)
	mockRepo.On("CountAll", ctx, param).Return(expectedCount, nil)

	// Execute
	results, pagination, err := svc.FindAll(ctx, param)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, expectedCount, pagination.TotalData)
	assert.Equal(t, 2, len(results))
	mockRepo.AssertExpectations(t)
}

func TestService_CreateMember_Success(t *testing.T) {
	// Setup
	svc, mockRepo := setupTestService()
	ctx := context.Background()
	request := &member.MemberRequest{
		Name: "New User",
		Info: member.MemberInfo{
			Address: "New Address",
			Salary:  5000,
			Age:     25,
		},
	}

	expectedID := int64(1)

	// Mock behavior - will check that the correct member entity is created
	mockRepo.On("CreateMember", ctx, mock.AnythingOfType("*member.Member")).
		Run(func(args mock.Arguments) {
			memberArg := args.Get(1).(*member.Member)
			assert.Equal(t, request.Name, memberArg.Name)
			// Info will be JSON string
			assert.Contains(t, memberArg.Info, `"address":"New Address"`)
			assert.Contains(t, memberArg.Info, `"salary":5000`)
			assert.Contains(t, memberArg.Info, `"age":25`)
		}).
		Return(expectedID, nil)

	// Execute
	result, err := svc.CreateMember(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedID, result.Id)
	assert.Equal(t, request.Name, result.Name)
	assert.Equal(t, request.Info.Address, result.Info.Address)
	assert.Equal(t, request.Info.Salary, result.Info.Salary)
	assert.Equal(t, request.Info.Age, result.Info.Age)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateMember_Success(t *testing.T) {
	// Setup
	svc, mockRepo := setupTestService()
	ctx := context.Background()
	updateID := int64(1)
	request := &member.MemberRequest{
		Name: "Updated User",
		Info: member.MemberInfo{
			Address: "Updated Address",
			Salary:  6000,
			Age:     26,
		},
	}

	// Mock for UpdateMember
	// Mock for UpdateMember
	mockRepo.On("UpdateMember", ctx, updateID, mock.AnythingOfType("*member.Member")).
		Run(func(args mock.Arguments) {
			memberArg := args.Get(2).(*member.Member)
			assert.Equal(t, request.Name, memberArg.Name)
			assert.Contains(t, memberArg.Info, `"address":"Updated Address"`)
			assert.Contains(t, memberArg.Info, `"salary":6000`)
			assert.Contains(t, memberArg.Info, `"age":26`)
		}).
		Return(int64(1), nil)

	// Execute
	result, err := svc.UpdateMember(ctx, updateID, request)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, updateID, result.Id)
	assert.Equal(t, request.Name, result.Name)
	assert.Equal(t, request.Info.Address, result.Info.Address)
	assert.Equal(t, request.Info.Salary, result.Info.Salary)
	assert.Equal(t, request.Info.Age, result.Info.Age)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteMember_Success(t *testing.T) {
	// Setup
	svc, mockRepo := setupTestService()
	ctx := context.Background()
	deleteID := int64(1)

	// Mock behavior
	mockRepo.On("DeleteMember", ctx, deleteID).Return(int64(1), nil)

	// Execute
	success, err := svc.DeleteMember(ctx, deleteID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, success)
	mockRepo.AssertExpectations(t)
}

func TestService_DeleteMember_NotFound(t *testing.T) {
	// Setup
	svc, mockRepo := setupTestService()
	ctx := context.Background()
	deleteID := int64(999)

	// Mock behavior - no rows affected
	mockRepo.On("DeleteMember", ctx, deleteID).Return(int64(0), nil)

	// Execute
	success, _ := svc.DeleteMember(ctx, deleteID)

	// Assert
	assert.True(t, success) // The operation succeeded but no rows were affected
	mockRepo.AssertExpectations(t)
}
