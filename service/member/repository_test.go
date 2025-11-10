package member_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	oracle "oracle.com/oracle/my-go-oracle-app/infra/database/sql"
	entity "oracle.com/oracle/my-go-oracle-app/service"
	"oracle.com/oracle/my-go-oracle-app/service/member"
)

// MockSlaveDB implements oracle.SlaveDB interface for testing
type MockSlaveDB struct {
	mock.Mock
}

func (m *MockSlaveDB) Rebind(query string) string {
	args := m.Called(query)
	return args.String(0)
}

func (m *MockSlaveDB) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSlaveDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSlaveDB) PreparexContext(ctx context.Context, query string) (oracle.SlaveStatement, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(oracle.SlaveStatement), args.Error(1)
}

func (m *MockSlaveDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	callArgs := m.Called(ctx, dest, query, args)
	return callArgs.Error(0)
}

func (m *MockSlaveDB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(*sqlx.Row)
}

func (m *MockSlaveDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	callArgs := m.Called(ctx, dest, query, args)
	return callArgs.Error(0)
}

func (m *MockSlaveDB) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(*sqlx.Rows), callArgs.Error(1)
}

// MockMasterDB implements oracle.MasterDB interface for testing
type MockMasterDB struct {
	mock.Mock
}

func (m *MockMasterDB) Rebind(query string) string {
	args := m.Called(query)
	return args.String(0)
}

func (m *MockMasterDB) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMasterDB) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMasterDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *MockMasterDB) Beginx() (*sqlx.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}

func (m *MockMasterDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(sql.Result), callArgs.Error(1)
}

func (m *MockMasterDB) PreparexContext(ctx context.Context, query string) (oracle.MasterStatement, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(oracle.MasterStatement), args.Error(1)
}

func (m *MockMasterDB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(*sqlx.Row)
}

func (m *MockMasterDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(*sql.Rows), callArgs.Error(1)
}

func (m *MockMasterDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	callArgs := m.Called(ctx, dest, query, args)
	return callArgs.Error(0)
}

func (m *MockMasterDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	callArgs := m.Called(ctx, dest, query, args)
	return callArgs.Error(0)
}

func (m *MockMasterDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	callArgs := m.Called(ctx, query, args)
	if row, ok := callArgs.Get(0).(*sql.Row); ok {
		return row
	}
	return &sql.Row{}
}

const findByIdQuery = "SELECT ID,NAME,INFO FROM MEMBER m WHERE id = :1 "

func setupTestRepo() (member.MemberRepository, *MockMasterDB, *MockSlaveDB) {
	mockMaster := new(MockMasterDB)
	mockSlave := new(MockSlaveDB)
	baseRepo := entity.BaseRepository{
		MasterDB: mockMaster,
		SlaveDB:  mockSlave,
	}
	repo := member.NewMemberRepository(baseRepo)
	return repo, mockMaster, mockSlave
}

func TestFindById_Success(t *testing.T) {
	// Setup
	repo, _, mockSlave := setupTestRepo()
	ctx := context.Background()
	expectedID := int64(1)
	expectedMember := member.Member{
		Name: "Test User",
		Info: "Test Info",
		BaseEntity: entity.BaseEntity{
			Id: expectedID,
		},
	}

	// Mock behavior
	mockSlave.On("GetContext",
		mock.Anything, // Accept any context
		mock.AnythingOfType("*member.Member"),
		findByIdQuery,
		mock.MatchedBy(func(args []interface{}) bool {
			// Verify the args array contains our expected ID
			if len(args) != 1 {
				return false
			}
			return args[0] == expectedID
		}),
	).Run(func(args mock.Arguments) {
		// Copy the expected member into the destination
		arg := args.Get(1).(*member.Member)
		*arg = expectedMember
	}).Return(nil)

	// Execute
	result, err := repo.FindById(ctx, expectedID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedMember.Id, result.Id)
	assert.Equal(t, expectedMember.Name, result.Name)
	assert.Equal(t, expectedMember.Info, result.Info)
	mockSlave.AssertExpectations(t)
}

func TestFindById_NotFound(t *testing.T) {
	// Setup
	repo, _, mockSlave := setupTestRepo()
	ctx := context.Background()
	expectedID := int64(999)

	// Mock behavior - simulate no rows found
	mockSlave.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*member.Member"),
		findByIdQuery,
		mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 1 && args[0] == expectedID
		}),
	).Return(sql.ErrNoRows)

	// Execute
	result, err := repo.FindById(ctx, expectedID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Empty(t, result)
	mockSlave.AssertExpectations(t)
}

func TestFindById_DBError(t *testing.T) {
	// Setup
	repo, _, mockSlave := setupTestRepo()
	ctx := context.Background()
	expectedID := int64(1)
	expectedError := errors.New("database error")

	// Mock behavior - simulate database error
	mockSlave.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*member.Member"),
		findByIdQuery,
		mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 1 && args[0] == expectedID
		}),
	).Return(expectedError)

	// Execute
	result, err := repo.FindById(ctx, expectedID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, result)
	mockSlave.AssertExpectations(t)
}

func TestFindById_ContextTimeout(t *testing.T) {
	// Setup
	repo, _, mockSlave := setupTestRepo()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	expectedID := int64(1)

	// Mock behavior - simulate slow DB response
	mockSlave.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*member.Member"),
		findByIdQuery,
		mock.MatchedBy(func(args []interface{}) bool {
			return len(args) == 1 && args[0] == expectedID
		}),
	).Run(func(args mock.Arguments) {
		time.Sleep(5 * time.Millisecond) // Simulate slow DB
	}).Return(context.DeadlineExceeded)

	// Execute
	result, err := repo.FindById(ctx, expectedID)

	// Assert
	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
	assert.Empty(t, result)
	mockSlave.AssertExpectations(t)
}

const (
	getAllMembersQuery = "SELECT ID,NAME,INFO FROM MEMBER m"
	countAllQuery      = "SELECT COUNT(*) as count FROM MEMBER m"
)

func TestGetAllMembers_Success(t *testing.T) {
	// Setup
	repo, _, mockSlave := setupTestRepo()
	ctx := context.Background()
	param := entity.SqlParameter{
		TableName: "MEMBER m",
		Columns:   []string{"ID", "NAME", "INFO"},
	}
	expectedMembers := []member.Member{
		{
			Name: "User 1",
			Info: "Info 1",
			BaseEntity: entity.BaseEntity{
				Id: 1,
			},
		},
		{
			Name: "User 2",
			Info: "Info 2",
			BaseEntity: entity.BaseEntity{
				Id: 2,
			},
		},
	}

	// Mock behavior
	mockSlave.On("SelectContext",
		mock.Anything,
		mock.AnythingOfType("*[]member.Member"),
		getAllMembersQuery,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		// Copy the expected members into the destination slice
		dest := args.Get(1).(*[]member.Member)
		*dest = expectedMembers
	}).Return(nil)

	// Execute
	members, err := repo.GetAllMembers(ctx, param)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, members, 2)
	assert.Equal(t, expectedMembers[0].Id, members[0].Id)
	assert.Equal(t, expectedMembers[0].Name, members[0].Name)
	assert.Equal(t, expectedMembers[1].Id, members[1].Id)
	assert.Equal(t, expectedMembers[1].Name, members[1].Name)
	mockSlave.AssertExpectations(t)
}

func TestGetAllMembers_Error(t *testing.T) {
	// Setup
	repo, _, mockSlave := setupTestRepo()
	ctx := context.Background()
	param := entity.SqlParameter{
		TableName: "MEMBER m",
		Columns:   []string{"ID", "NAME", "INFO"},
	}
	expectedError := errors.New("database error")

	// Mock behavior
	mockSlave.On("SelectContext",
		mock.Anything,
		mock.AnythingOfType("*[]member.Member"),
		getAllMembersQuery,
		mock.Anything,
	).Return(expectedError)

	// Execute
	members, err := repo.GetAllMembers(ctx, param)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, members)
	mockSlave.AssertExpectations(t)
}

func TestCountAll_Success(t *testing.T) {
	// Setup
	repo, _, mockSlave := setupTestRepo()
	ctx := context.Background()
	param := entity.SqlParameter{
		TableName: "MEMBER m",
	}
	expectedCount := int64(5)

	// Mock behavior
	mockSlave.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*int64"),
		countAllQuery,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*int64)
		*dest = expectedCount
	}).Return(nil)

	// Execute
	count, err := repo.CountAll(ctx, param)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockSlave.AssertExpectations(t)
}

func TestCountAll_Error(t *testing.T) {
	// Setup
	repo, _, mockSlave := setupTestRepo()
	ctx := context.Background()
	param := entity.SqlParameter{
		TableName: "MEMBER m",
	}
	expectedError := errors.New("database error")

	// Mock behavior
	mockSlave.On("GetContext",
		mock.Anything,
		mock.AnythingOfType("*int64"),
		countAllQuery,
		mock.Anything,
	).Return(expectedError)

	// Execute
	count, err := repo.CountAll(ctx, param)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Zero(t, count)
	mockSlave.AssertExpectations(t)
}

type mockResult struct {
	lastId       int64
	rowsAffected int64
}

func (m mockResult) LastInsertId() (int64, error) {
	return m.lastId, nil
}

func (m mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func TestCreateMember_Success(t *testing.T) {
	// Setup
	repo, mockMaster, _ := setupTestRepo()
	ctx := context.Background()
	newMember := &member.Member{
		Name: "New User",
		Info: "New Info",
	}
	expectedLastID := int64(1)

	// Mock behavior for ExecContext (Oracle RETURNING INTO clause)
	mockMaster.On("ExecContext",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Run(func(args mock.Arguments) {
		queryArgs := args.Get(2).([]interface{})
		if len(queryArgs) >= 3 {
			if out, ok := queryArgs[2].(sql.Out); ok {
				*(out.Dest.(*int64)) = expectedLastID
			}
		}
	}).Return(mockResult{rowsAffected: 1}, nil)

	// Execute
	lastID, err := repo.CreateMember(ctx, newMember)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedLastID, lastID)
	assert.Equal(t, expectedLastID, newMember.Id)
	mockMaster.AssertExpectations(t)
}

func TestUpdateMember_Success(t *testing.T) {
	// Setup
	repo, mockMaster, _ := setupTestRepo()
	ctx := context.Background()
	updateID := int64(1)
	updateMember := &member.Member{
		Name: "Updated User",
		Info: "Updated Info",
	}
	expectedRowsAffected := int64(1)
	mockResult := &mockResult{rowsAffected: expectedRowsAffected}

	// Mock behavior
	mockMaster.On("ExecContext",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mockResult, nil)

	// Execute
	rowsAffected, err := repo.UpdateMember(ctx, updateID, updateMember)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRowsAffected, rowsAffected)
	mockMaster.AssertExpectations(t)
}

func TestDeleteMember_Success(t *testing.T) {
	// Setup
	repo, mockMaster, _ := setupTestRepo()
	ctx := context.Background()
	deleteID := int64(1)
	expectedRowsAffected := int64(1)
	mockResult := &mockResult{rowsAffected: expectedRowsAffected}

	// Mock behavior
	mockMaster.On("ExecContext",
		mock.Anything,
		mock.AnythingOfType("string"),
		mock.Anything,
	).Return(mockResult, nil)

	// Execute
	rowsAffected, err := repo.DeleteMember(ctx, deleteID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRowsAffected, rowsAffected)
	mockMaster.AssertExpectations(t)
}
