package constants

import (
	"time"
)

func init() {
	JAKARTA_LOCATION, _ = time.LoadLocation("Asia/Jakarta")
}

var JAKARTA_LOCATION *time.Location

const (
	ACCEPT_LANG   = contextKey("Accept-Language")
	AUTHORIZATION = contextKey("Authorization")

	COMMA                  = ","
	FROM                   = " FROM "
	WHERE                  = " WHERE "
	AND                    = " AND "
	OR                     = " OR "
	ORDER_BY               = " ORDER BY "
	GROUP_GY               = " GROUP BY "
	HAVING                 = " HAVING "
	ASC                    = "asc"
	DESC                   = "desc"
	LIMIT                  = "limit"
	LAST                   = "last"
	OFFSET                 = "offset"
	SEARCH_BY              = "searchBy"
	STATUS                 = "status"
	PRODUCT_GROUP_ID_FIELD = "product_group_id"
	PRODUCT_GROUP_ID_KEY   = "productGroupId"
	UPDATE_FROM            = "lastUpdateStart"
	UPDATE_TO              = "lastUpdateEnd"
	ISSUED_FROM            = "issuedStart"
	ISSUED_TO              = "issuedEnd"
	CODE                   = "code"
	DEFAULT_LIMIT          = 10
	DEFAULT_OFFSET         = 0
	SUCCESS                = "SUCCESS"
	PENDING                = "PENDING"
	FAILED                 = "FAILED"
	ORDER_CONFLICT         = "ORDER_CONFLICT"
	ORDER_CHANGES          = "ORDER_CHANGES"
	LANG_ID                = "ID"
	LANG_EN                = "EN"
	EMPTY_STR              = ""
	PAGE                   = "page"
	ORDER_BY_REQ           = "orderBy"
	ORDER_TYPE_REQ         = "orderType"
	COUNT_COL              = "COUNT(*) as count"
	UPDATED_DATE           = "updated_date"
	ISSUED_DATE            = "issued_date"
	CREATED_DATE           = "created_date"
	TRANSACTION_FROM       = "transactionStart"
	TRANSACTION_TO         = "transactionEnd"
	MSG_ORDER_ID_KEY       = "orderId"
	MSG_PRODUCT_NAME_KEY   = "productName"
	SERVICE_NAME           = "my-go-oracle-app"
	REQUEST_DATE_START     = "requestDateStart"
	REQUEST_DATE_END       = "requestDateEnd"
	REQUEST_DATE           = "request_date"
	ARRIVAL_DATE           = "ta.travel_start"
	REFUND_DATE            = "ac.created_date"
	EXPIRED_DATE           = "td.expired_at"
	CLAIM_DATE             = "ic.claim_date"
	SUBJECT_KEY            = "subject"
	DELETED_COLUMN         = "is_deleted"
	UNIQUE_IDENTIFIER      = "0712"
	UNIQUE_IDENTIFIER_INT  = 712
	INDONESIAN_RUPIAH      = "IDR"
	CLAIM_NOTIF_DATE       = "fcn.created_date"

	TD_ISSUED_DATE = "td.issued_date"
)

const (
	DATA_NOT_FOUND      = "Data not found"
	CLIENT_NOT_FOUND    = "CLIENT_NOT_FOUND"
	DATA_NOT_FOUND_CODE = "DATA_NOT_FOUND"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

type StatusType struct {
	Active, Inactive string
}

var Status = StatusType{
	"ACTIVE", "INACTIVE",
}

var StatusName = map[string]string{
	"active":   Status.Active,
	"inactive": Status.Inactive,
}

const (
	EQUAL              = "="
	NOT_EQUAL          = "<>"
	LESS_THAN          = "<"
	LESS_THAN_EQUAL    = "<="
	GREATER_THAN       = ">"
	GREATER_THAN_EQUAL = ">="
	LIKE               = "LIKE"
	IS_NULL            = "IS NULL"
	IN                 = "IN"    // column IN (vals)
	REVERSE_IN         = "RIN"   // val IN (columns)
	MULTIPLE_LIKE      = "MLIKE" // ( col1 LIKE "%sSTR%s" or  col2 LIKE "%sSTR%s")
	NOT                = "!"
	NOT_IN             = "NOT IN"
	CONTAINS           = "∋"
	NOT_CONTAINS       = "∌"
	IN_LIKE_STRING     = "IN_LIKE" // ( col LIKE "%sSTR1%s" or  col LIKE "%sSTR2%s")
	MULTIPLE_EQUAL     = "M_EQUAL"
)

const (
	ACTION_INSERT = "insert"
	ACTION_UPDATE = "update"
	ACTION_DELETE = "delete"
	ACTION_FIND   = "find"
	ACTION_REMOVE = "remove"
)

const (
	CONTEXT_TRANSACTION = "trxConn"
)
const (
	REPORT_DATE_FORMAT     = "02/Jan/2006"
	REPORT_TIME_FORMAT     = "15:04:05"
	DATE_TIME_FORMAT       = "2006-01-02 15:04:05"
	DATE_TIME_FORMAT_SHORT = "2006-01-02 15:04"
	DATE_FORMAT            = "2006-01-02"
	TIME_FORMAT            = "15:04"

	VARIABLE_DATE_TIME_FORMAT = "2006-01-02T15:04:05Z0700"

	EMAIL_FORMAT_EN = "Mon, Jan 2 2006"
	EMAIL_FORMAT_ID = "2 Jan 2006"

	LONG_DATE_TIME_FORMAT = "Mon, 2 Jan 2006, 15:04"
)
