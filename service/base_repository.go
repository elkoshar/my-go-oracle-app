package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"oracle.com/oracle/my-go-oracle-app/infra/database"
	oracle "oracle.com/oracle/my-go-oracle-app/infra/database/sql"
	"oracle.com/oracle/my-go-oracle-app/pkg/constants"
)

const logName = "[Base Repository][Operation]"
const maximumCallerDepth = 1
const (
	CONTEXT_TRANSACTION = "trxConn"
)

type BaseRepository struct {
	MasterDB oracle.MasterDB
	SlaveDB  oracle.SlaveDB
}

type BaseRepositoryInterface interface {
	Insert(ctx context.Context, sqlParameter SqlParameter) (int64, error)
	Update(ctx context.Context, sqlParameter SqlParameter) (int64, error)
	Delete(ctx context.Context, sqlParameter SqlParameter) (int64, error)
}

// SelectWithParameter can return multiple row. dest must be pointer to a slice
func (r *BaseRepository) SelectOperations(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	slog.InfoContext(ctx, fmt.Sprintf("query= %v, paramValue=%v,", query, args))
	newContext := database.StartMetrics(ctx, database.Event{
		Name: GetLastFuncCallerName(),
	})

	// Log incoming context deadline information to help debug timeouts
	if dl, ok := ctx.Deadline(); ok {
		slog.InfoContext(ctx, fmt.Sprintf("context deadline=%v, remaining=%v", dl, time.Until(dl)))
	} else {
		slog.InfoContext(ctx, "context has no deadline")
	}

	var err error
	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)
	if txConn == nil {
		err = r.SlaveDB.SelectContext(newContext, dest, query, args...)
	} else {
		err = txConn.(*sqlx.Tx).SelectContext(newContext, dest, query, args...)
	}

	if err != nil {
		return err
	}

	return nil
}

// SelectWithParameter can return multiple row. dest must be pointer to a slice
func (r *BaseRepository) GetOperations(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	slog.InfoContext(ctx, fmt.Sprintf("query= %v, paramValue=%v,", query, args))
	newContext := database.StartMetrics(ctx, database.Event{
		Name: GetLastFuncCallerName(),
	})

	// Log incoming context deadline information to help debug timeouts
	if dl, ok := ctx.Deadline(); ok {
		slog.InfoContext(ctx, fmt.Sprintf("context deadline=%v, remaining=%v", dl, time.Until(dl)))
	} else {
		slog.InfoContext(ctx, "context has no deadline")
	}

	var err error
	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)
	if txConn == nil {
		err = r.SlaveDB.GetContext(newContext, dest, query, args...)
	} else {
		err = txConn.(*sqlx.Tx).GetContext(newContext, dest, query, args...)
	}

	if err != nil {
		return err
	}
	return nil
}

// SelectWithParameter can return multiple row. dest must be pointer to a slice
func (r *BaseRepository) GetOperationsMasterConn(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	slog.InfoContext(ctx, fmt.Sprintf("query= %v, paramValue=%v,", query, args))
	newContext := database.StartMetrics(ctx, database.Event{
		Name: GetLastFuncCallerName(),
	})
	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)
	if txConn != nil {
		err = txConn.(*sqlx.Tx).GetContext(newContext, dest, query, args...)
		if err != nil {
			return err
		}

		return nil
	}
	err = r.MasterDB.GetContext(newContext, dest, query, args...)

	if err != nil {
		return err
	}

	return nil

}

// SelectWithParameter can return multiple row. dest must be pointer to a slice
func (r *BaseRepository) SelectWithParameter(ctx context.Context, dest interface{}, param SqlParameter) error {
	query, args := r.GenerateQuerySelectWithParams("", param)
	slog.InfoContext(ctx, fmt.Sprintf("query= %v, paramValue=%v,", query, args))
	newContext := database.StartMetrics(ctx, database.Event{
		Name: GetLastFuncCallerName(),
	})

	var err error
	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)
	if txConn == nil {
		err = r.SlaveDB.SelectContext(newContext, dest, query, args...)
	} else {
		err = txConn.(*sqlx.Tx).SelectContext(newContext, dest, query, args...)
	}

	if err != nil {

		return err
	}

	return nil

}

// GetWithParameter will return only 1 row. dest must be pointer to a struct
func (r *BaseRepository) GetWithParameter(ctx context.Context, dest interface{}, param SqlParameter) error {
	query, args := r.GenerateQuerySelectWithParams("", param)
	slog.InfoContext(ctx, fmt.Sprintf("query= %v, paramValue=%v,", query, args))
	newContext := database.StartMetrics(ctx, database.Event{
		Name: GetLastFuncCallerName(),
	})
	var err error
	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)
	if txConn == nil {
		err = r.SlaveDB.GetContext(newContext, dest, query, args...)
	} else {
		err = txConn.(*sqlx.Tx).GetContext(newContext, dest, query, args...)
	}

	if err != nil {

		return err
	}

	return nil
}

// WriteOrUpdateOperation will execute query.
func (r *BaseRepository) WriteOrUpdateOperation(ctx context.Context, query string, returnedID *int64, args ...interface{}) (int64, error) {
	slog.InfoContext(ctx, fmt.Sprintf("query= %v, paramValue=%v, returnedID=%v", query, args, returnedID))

	var (
		result sql.Result
		err    error
	)

	// Check if the operation is a RETURNING INTO query (must have a pointer to the returned ID)
	isReturningInto := strings.HasPrefix(strings.ToUpper(query), "INSERT") && returnedID != nil

	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)

	if isReturningInto {
		// --- Special Handling for Oracle RETURNING INTO ---
		var row *sql.Row
		if txConn == nil {
			// Use QueryRowContext for the master database connection
			row = r.MasterDB.QueryRowContext(ctx, query, args...)
		} else {
			// Use QueryRowContext for the transaction connection
			row = txConn.(*sqlx.Tx).QueryRowContext(ctx, query, args...)
		}

		// Scan the result into the provided ID pointer.
		// Note: godror sets the output parameter before Scan, but Scan is needed to check for errors.
		err = row.Scan(returnedID)

		// sql.ErrNoRows might be returned if the RETURNING INTO clause fails,
		// but for a successful INSERT, we still treat it as a successful row change.
		if err != nil && err != sql.ErrNoRows {
			return 0, err
		}

		// Return 1 row affected since the INSERT was successful and the ID was captured.
		return 1, nil

	} else {
		// --- Standard ExecContext for UPDATE, DELETE, or simple INSERT ---

		if txConn == nil {
			result, err = r.MasterDB.ExecContext(ctx, query, args...)
		} else {
			result, err = txConn.(*sqlx.Tx).ExecContext(ctx, query, args...)
		}

		if err != nil {
			return 0, err
		}

		// Return RowsAffected for standard UPDATE, DELETE, and simple INSERTs.
		return result.RowsAffected()
	}
}

func (r *BaseRepository) WriteOrUpdateOperation2(ctx context.Context, query string, args ...interface{}) (int64, error) {
	slog.InfoContext(ctx, fmt.Sprintf("query= %v, paramValue=%v,", query, args))

	var (
		result sql.Result
		err    error
	)

	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)
	if txConn == nil {
		result, err = r.MasterDB.ExecContext(ctx, query, args...)
	} else {
		result, err = txConn.(*sqlx.Tx).ExecContext(ctx, query, args...)
	}

	if err != nil {

		return 0, err
	}

	if strings.HasPrefix(strings.ToUpper(query), "INSERT") {
		return result.LastInsertId()
	}
	return result.RowsAffected()
}

func (r *BaseRepository) Insert(ctx context.Context, sqlParameter SqlParameter) (int64, error) {
	sql, args := r.GenerateQueryInsert(sqlParameter)

	res, err := r.WriteOrUpdateOperation(ctx, sql, nil, args...)

	return res, err
}

func (r *BaseRepository) Update(ctx context.Context, sqlParameter SqlParameter) (int64, error) {
	sql, args := r.GenerateQueryUpdate(sqlParameter)

	res, err := r.WriteOrUpdateOperation(ctx, sql, nil, args...)

	return res, err
}

func (r *BaseRepository) Delete(ctx context.Context, sqlParameter SqlParameter) (int64, error) {
	sql := fmt.Sprintf("DELETE FROM %s", sqlParameter.TableName)
	conditional, args := r.GenerateConditional(sqlParameter)
	sql += conditional
	res, err := r.WriteOrUpdateOperation(ctx, sql, nil, args...)

	return res, err
}

func (r *BaseRepository) GenerateQueryInsert(sqlParameter SqlParameter) (string, []interface{}) {
	var s strings.Builder
	s.WriteString("INSERT INTO ")
	s.WriteString(sqlParameter.TableName)
	s.WriteString(" (")
	for i := 0; i < len(sqlParameter.Values); i++ {
		s.WriteString(sqlParameter.Values[i].Field)
		if i < len(sqlParameter.Values)-1 {
			s.WriteString(constants.COMMA)
		}
	}
	var args []interface{}
	s.WriteString(") VALUES (")
	for i := 0; i < len(sqlParameter.Values); i++ {
		s.WriteString("?")
		args = append(args, sqlParameter.Values[i].Value)
		if i < len(sqlParameter.Values)-1 {
			s.WriteString(constants.COMMA)
		}
	}
	s.WriteString(")")
	return s.String(), args
}

func (r *BaseRepository) GenerateQueryUpdate(sqlParameter SqlParameter) (string, []interface{}) {
	var s strings.Builder
	s.WriteString("UPDATE ")
	s.WriteString(sqlParameter.TableName)
	s.WriteString(" SET ")
	var args []interface{}
	for i := 0; i < len(sqlParameter.Values); i++ {
		s.WriteString(sqlParameter.Values[i].Field)
		args = append(args, sqlParameter.Values[i].Value)
		s.WriteString("=?")
		if i < len(sqlParameter.Values)-1 {
			s.WriteString(constants.COMMA)
		}
	}

	conditional, conditionalArgs := r.GenerateConditional(sqlParameter)
	s.WriteString(conditional)
	return s.String(), append(args, conditionalArgs...)
}

func (r *BaseRepository) GenerateQuerySelectFrom(sqlParameter SqlParameter) string {
	sql := "SELECT "
	for i := 0; i < len(sqlParameter.Columns); i++ {
		sql += sqlParameter.Columns[i]
		if i < len(sqlParameter.Columns)-1 {
			sql += constants.COMMA
		}
	}
	sql += constants.FROM
	sql += sqlParameter.TableName

	return sql
}

func (r *BaseRepository) GenerateConditional(sqlParameter SqlParameter) (string, []interface{}) {
	var sql strings.Builder
	var args []interface{}
	if len(sqlParameter.Params) != 0 {
		sql.WriteString(constants.WHERE)
		for i := 0; i < len(sqlParameter.Params); i++ {
			if sqlParameter.Params[i].Operand == "" {
				sqlParameter.Params[i].Operand = "="
			}
			switch sqlParameter.Params[i].Operand {
			case constants.IN:
				tempQuery := fmt.Sprintf("%s %s (:%d)", sqlParameter.Params[i].Field, constants.IN, i)
				tempQuery, tempArgs, _ := sqlx.In(tempQuery, sqlParameter.Params[i].Value)
				sql.WriteString(tempQuery)
				args = append(args, tempArgs...)
			case constants.REVERSE_IN:
				sql.WriteString(fmt.Sprintf("? %s (%s)", constants.IN, sqlParameter.Params[i].Field))
				args = append(args, sqlParameter.Params[i].Value)
			case constants.MULTIPLE_LIKE:
				or, orarg := r.generateOrClause(strings.Split(sqlParameter.Params[i].Field, ","), constants.LIKE, sqlParameter.Params[i].Value)
				sql.WriteString(or)
				args = append(args, orarg...)
			case constants.MULTIPLE_EQUAL:
				or, orarg := r.generateOrClause(strings.Split(sqlParameter.Params[i].Field, ","), constants.EQUAL, sqlParameter.Params[i].Value)
				sql.WriteString(or)
				args = append(args, orarg...)
			case constants.IN_LIKE_STRING:
				or, orarg := r.generateInLikeClause(sqlParameter.Params[i].Field, sqlParameter.Params[i].Value)
				sql.WriteString(or)
				args = append(args, orarg...)
			case constants.IS_NULL:
				sql.WriteString(fmt.Sprintf("%s %s", sqlParameter.Params[i].Field, sqlParameter.Params[i].Operand))
			case constants.NOT_IN:
				tempQuery := fmt.Sprintf("%s %s (:%d)", sqlParameter.Params[i].Field, constants.NOT_IN, i)
				tempQuery, tempArgs, _ := sqlx.In(tempQuery, sqlParameter.Params[i].Value)
				sql.WriteString(tempQuery)
				args = append(args, tempArgs...)
			default:
				sql.WriteString(fmt.Sprintf("%s %s :%d", sqlParameter.Params[i].Field, sqlParameter.Params[i].Operand, i))
				args = append(args, sqlParameter.Params[i].Value)
			}
			if i < len(sqlParameter.Params)-1 {
				sql.WriteString(constants.AND)
			}
		}
	}

	return sql.String(), args
}

func (r *BaseRepository) GenerateQuerySelectWithParams(query string, sqlParameter SqlParameter) (string, []interface{}) {
	var sql strings.Builder
	var argsJoin []interface{}
	if query == "" {
		query = r.GenerateQuerySelectFrom(sqlParameter)
	}
	sql.WriteString(query)
	totalParams := len(sqlParameter.Params)
	totalOrder := len(sqlParameter.OrderBy)
	if len(sqlParameter.Joins) > 0 {
		for _, join := range sqlParameter.Joins {
			sql.WriteString(fmt.Sprintf(" %s JOIN %s", join.JoinType, join.Table))
			if join.Alias != "" {
				sql.WriteString(fmt.Sprintf(" %s", join.Alias))
			}
			sql.WriteString(fmt.Sprintf(" ON %s", join.On))

			for i, cond := range join.Conditions {
				// to be: AND io.tablename IN (?, ?)
				sqlCond, condArgs := buildCondition(cond, i)
				sql.WriteString(" AND " + sqlCond)

				argsJoin = append(argsJoin, condArgs...)
			}
		}
	}

	conditional, args := r.GenerateConditional(sqlParameter)
	sql.WriteString(conditional)

	if len(sqlParameter.GroupBy) != 0 {
		sql.WriteString(constants.GROUP_GY)
		sql.WriteString(strings.Join(sqlParameter.GroupBy, ", "))
	}

	if len(sqlParameter.Having) != 0 {
		sql.WriteString(constants.HAVING)
		sql.WriteString("(")
		sql.WriteString(strings.Join(sqlParameter.Having, " AND "))
		sql.WriteString(")")
	}

	if len(sqlParameter.OrderBy) != 0 {
		sql.WriteString(constants.ORDER_BY)
		for i := 0; i < len(sqlParameter.OrderBy); i++ {
			sql.WriteString(sqlParameter.OrderBy[i])
			if i < len(sqlParameter.OrderBy)-1 {
				sql.WriteString(constants.COMMA)
			}
		}
	}

	if sqlParameter.Limit != 0 {
		sql.WriteString(fmt.Sprintf(" OFFSET :%d ROWS", totalParams+totalOrder+1))
		args = append(args, sqlParameter.Offset)

		sql.WriteString(fmt.Sprintf(" FETCH NEXT :%d ROWS ONLY", totalParams+totalOrder+2))
		args = append(args, sqlParameter.Limit)
	}

	if len(argsJoin) > 0 {
		args = append(argsJoin, args...)
	}

	return sql.String(), args
}

func buildCondition(param FilterParam, i int) (string, []interface{}) {
	switch v := param.Value.(type) {
	case []string:
		placeholders := make([]string, len(v))
		args := make([]interface{}, len(v))
		for i, val := range v {
			placeholders[i] = fmt.Sprintf(":%d", i)
			args[i] = val
		}
		return fmt.Sprintf("%s %s (%s)", param.Field, param.Operand, strings.Join(placeholders, ", ")), args
	case string:
		return fmt.Sprintf("%s %s :%d", param.Field, param.Operand, i), []interface{}{v}
	default:
		return fmt.Sprintf("%s %s :%d", param.Field, param.Operand, i), []interface{}{param.Value}
	}
}

func MakeFilterParam(field, operand string, val interface{}) FilterParam {
	return FilterParam{
		Field:   field,
		Operand: operand,
		Value:   val,
	}
}

type TxDb struct {
	Tx oracle.TransactionDB
}

func New(tx oracle.TransactionDB) *TxDb {
	return &TxDb{
		Tx: tx,
	}
}

func (r *BaseRepository) ExecuteWithTx(ctx context.Context, fn func(*TxDb) error) error {
	tx, err := r.MasterDB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	t := New(tx)
	err = fn(t)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			slog.Error(fmt.Sprintf("tx err: %v, rb err: %v", err, rbErr))
			return rbErr
		}
		slog.Error(fmt.Sprintf("Rollback tx err: %v", err))
		return err
	}
	return tx.Commit()
}

func (r *BaseRepository) MasterClose() {
	r.MasterDB.Close()
}

func (r *BaseRepository) SlaveClose() {
	r.SlaveDB.Close()
}

// generateOrClause will generatate or clause
func (r *BaseRepository) generateOrClause(listCol []string, operator string, arg interface{}) (string, []interface{}) {
	count := len(listCol)
	var s strings.Builder
	var args []interface{}
	s.WriteString("(")
	for i := range listCol {
		s.WriteString(listCol[i])
		s.WriteString(" ")
		s.WriteString(operator)
		s.WriteString(" ?")
		args = append(args, arg)
		if i < count-1 {
			s.WriteString(constants.OR)
		}
	}
	s.WriteString(")")
	return s.String(), args
}

func (r *BaseRepository) generateInLikeClause(listCol string, arg interface{}) (string, []interface{}) {
	var s strings.Builder
	var args []interface{}
	str, _ := arg.(string)
	argString := strings.Split(str, ",")
	count := len(argString)

	s.WriteString("(")
	for i := range argString {
		s.WriteString(listCol)
		s.WriteString(" LIKE ?")
		args = append(args, "%"+argString[i]+"%")
		if i < count-1 {
			s.WriteString(constants.OR)
		}
	}
	s.WriteString(")")

	return s.String(), args
}

// SetDBConnInContext set db connection in context
func SetTxConnInContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, constants.CONTEXT_TRANSACTION, tx)
}

func RemoveTxConnInContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, constants.CONTEXT_TRANSACTION, nil)
}

func GetTxConnInContext(ctx context.Context) (*sqlx.Tx, bool) {
	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)
	if txConn != nil {
		return txConn.(*sqlx.Tx), true
	}

	return nil, false
}

func GetLastFuncCallerName() (name string) {
	funcCall := make([]uintptr, maximumCallerDepth)
	_ = runtime.Callers(3, funcCall)
	if len(funcCall) <= 0 {
		return
	}

	eventCaller := runtime.FuncForPC(funcCall[0])
	indexedName := strings.Split(eventCaller.Name(), "/")
	if len(indexedName) <= 0 {
		return
	}
	name = indexedName[len(indexedName)-1]
	return
}

type TimeFilter struct {
	IssuedStart      time.Time
	IssuedEnd        time.Time
	RequestStart     time.Time
	RequestEnd       time.Time
	ArrivalTimeStart time.Time
	ArrivalTimeEnd   time.Time
	RefundDateStart  time.Time
	RefundDateEnd    time.Time
	ExpiredStart     time.Time
	ExpiredEnd       time.Time
	ClaimStart       time.Time
	ClaimEnd         time.Time
	ClaimNotifStart  time.Time
	ClaimNotifEnd    time.Time
}

func GetTimeFilter(params SqlParameter) TimeFilter {
	var resp TimeFilter
	for i := range params.Params {
		if params.Params[i].Field == constants.ISSUED_DATE && params.Params[i].Operand == constants.GREATER_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.IssuedStart = t

		}
		if params.Params[i].Field == constants.ISSUED_DATE && params.Params[i].Operand == constants.LESS_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.IssuedEnd = t
		}

		if params.Params[i].Field == constants.REQUEST_DATE && params.Params[i].Operand == constants.GREATER_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.RequestStart = t
		}
		if params.Params[i].Field == constants.REQUEST_DATE && params.Params[i].Operand == constants.LESS_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.RequestEnd = t
		}

		if params.Params[i].Field == constants.ARRIVAL_DATE && params.Params[i].Operand == constants.GREATER_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.ArrivalTimeStart = t
		}

		if params.Params[i].Field == constants.ARRIVAL_DATE && params.Params[i].Operand == constants.LESS_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.ArrivalTimeEnd = t
		}

		if params.Params[i].Field == constants.REFUND_DATE && params.Params[i].Operand == constants.GREATER_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.RefundDateStart = t
		}

		if params.Params[i].Field == constants.REFUND_DATE && params.Params[i].Operand == constants.LESS_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.RefundDateEnd = t
		}

		if params.Params[i].Field == constants.EXPIRED_DATE && params.Params[i].Operand == constants.GREATER_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.ExpiredStart = t
		}

		if params.Params[i].Field == constants.EXPIRED_DATE && params.Params[i].Operand == constants.LESS_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.ExpiredEnd = t
		}

		if params.Params[i].Field == constants.CLAIM_DATE && params.Params[i].Operand == constants.GREATER_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.ClaimStart = t
		}

		if params.Params[i].Field == constants.CLAIM_DATE && params.Params[i].Operand == constants.LESS_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.ClaimEnd = t
		}

		if params.Params[i].Field == constants.CLAIM_NOTIF_DATE && params.Params[i].Operand == constants.GREATER_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.ClaimNotifStart = t
		}

		if params.Params[i].Field == constants.CLAIM_NOTIF_DATE && params.Params[i].Operand == constants.LESS_THAN_EQUAL {
			var t time.Time
			t, _ = time.Parse(constants.DATE_TIME_FORMAT, fmt.Sprintf("%v", params.Params[i].Value))
			resp.ClaimNotifEnd = t
		}

	}

	return resp
}

func IsTimeFilterValid(filter TimeFilter, max int) error {
	emptyTime := time.Time{}
	validIssued := filter.IssuedEnd != emptyTime && filter.IssuedStart != emptyTime
	onRangeIssued := validIssued && filter.IssuedStart.Add(24*time.Hour*time.Duration(max)).Unix() >= filter.IssuedEnd.Unix()

	validRequest := filter.RequestEnd != emptyTime && filter.RequestStart != emptyTime
	onRangeRequest := validRequest && filter.RequestStart.Add(24*time.Hour*time.Duration(max)).Unix() >= filter.RequestEnd.Unix()

	validArrival := filter.ArrivalTimeEnd != emptyTime && filter.ArrivalTimeStart != emptyTime
	onRangeArrival := validArrival && filter.ArrivalTimeStart.Add(24*time.Hour*time.Duration(max)).Unix() >= filter.ArrivalTimeEnd.Unix()

	validRefund := filter.RefundDateEnd != emptyTime && filter.RefundDateStart != emptyTime
	onRangeRefund := validRefund && filter.ArrivalTimeStart.Add(24*time.Hour*time.Duration(max)).Unix() >= filter.ArrivalTimeEnd.Unix()

	validExpire := filter.ExpiredEnd != emptyTime && filter.ExpiredStart != emptyTime
	onRangeExpire := validExpire && filter.ExpiredStart.Add(24*time.Hour*time.Duration(max)).Unix() >= filter.ExpiredEnd.Unix()

	validClaim := filter.ClaimEnd != emptyTime && filter.ClaimStart != emptyTime
	onRangeClaim := validClaim && filter.ClaimStart.Add(24*time.Hour*time.Duration(max)).Unix() >= filter.ClaimEnd.Unix()

	validClaimNotif := filter.ClaimNotifEnd != emptyTime && filter.ClaimNotifStart != emptyTime
	onRangeClaimNotif := validClaimNotif && filter.ClaimNotifStart.Add(24*time.Hour*time.Duration(max)).Unix() >= filter.ClaimNotifEnd.Unix()

	if !validRequest && !validIssued && !validArrival && !validRefund && !validExpire && !validClaim && !validClaimNotif {
		return fmt.Errorf("date filter is mandatory")
	}

	if !onRangeArrival && !onRangeIssued && !onRangeRefund && !onRangeRequest && !onRangeExpire && !onRangeClaim &&
		!onRangeClaimNotif && max > 0 {
		return fmt.Errorf("max range is %d days", max)
	}

	return nil
}

func (b *BaseRepository) PreparexContext(ctx context.Context, query string) (oracle.MasterStatement, error) {
	txConn := ctx.Value(constants.CONTEXT_TRANSACTION)
	if txConn == nil {
		return b.MasterDB.PreparexContext(ctx, query)
	}
	return txConn.(*sqlx.Tx).PreparexContext(ctx, query)
}
