package servicehelper

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"oracle.com/oracle/my-go-oracle-app/pkg/constants"
	"oracle.com/oracle/my-go-oracle-app/pkg/helpers"
	"oracle.com/oracle/my-go-oracle-app/service"
)

const startDateSuffix = " 00:00:00"
const endDateSuffix = " 23:59:59"

// list of param keys that needs to be appended by certain value.
var appendList = map[string]string{}

// GetFilterParamFromRequest will parse request url query and make corresponding filterparam based
func GetFilterParamFromRequest(r *http.Request, mapFilter map[string]service.FilterParam) []service.FilterParam {
	var resp []service.FilterParam
	for key, val := range mapFilter {
		str := r.URL.Query().Get(key)
		if str != "" {
			val.Value = preprocessInput(key, str)

			if ok := helpers.StringExists([]string{constants.IN, constants.NOT_IN}, val.Operand); ok {
				arrString := trimspaceArrString(str)
				val.Value = arrString
			}

			if val.Operand == constants.LIKE || val.Operand == constants.MULTIPLE_LIKE {
				val.Value = "%" + str + "%"
			}

			resp = append(resp, val)
		}
	}

	return resp
}

func trimspaceArrString(str string) []string {
	arrString := strings.Split(str, constants.COMMA)
	for i := range arrString {
		arrString[i] = strings.TrimSpace(arrString[i])
	}
	return arrString
}

// parse will parse input to utc value for time value. if not time then will return as is
func preprocessInput(key, input string) string {
	v := appendList[key]
	if v == "" {
		return input
	}

	// if time already included hh:mm:ss parse, if not add 00:00:00 or 23:59:59 then parse
	t, err := time.ParseInLocation(constants.DATE_TIME_FORMAT, input, constants.JAKARTA_LOCATION)
	if err != nil {
		input += v
		t, _ = time.ParseInLocation(constants.DATE_TIME_FORMAT, input, constants.JAKARTA_LOCATION)
	}
	return t.UTC().Format(constants.DATE_TIME_FORMAT)
}

// GetStatusCode Validate status code
func GetStatusCode(statusCode string) (status string, err error) {
	status, ok := constants.StatusName[strings.ToLower(statusCode)]
	if !ok {
		err = fmt.Errorf("invalid status code")
		return
	}
	return
}

// GelSqlParameterFromRequest will parse request url query and make corresponding sqlparamater
func GelSqlParameterFromRequest(r *http.Request, mapFilter map[string]service.FilterParam, mapOrder map[string]string) service.SqlParameter {
	filterParam := GetFilterParamFromRequest(r, mapFilter)

	limit, _ := strconv.Atoi(r.URL.Query().Get(constants.LIMIT))
	if limit < 1 {
		limit = constants.DEFAULT_LIMIT
	}

	page, _ := strconv.Atoi(r.URL.Query().Get(constants.PAGE))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	orderType := r.URL.Query().Get(constants.ORDER_TYPE_REQ)
	if orderType == "" {
		orderType = constants.ASC
	}

	orderBy := []string{}
	order := r.URL.Query().Get(constants.ORDER_BY_REQ)
	if val, ok := mapOrder[order]; ok {
		orderBy = append(orderBy, fmt.Sprintf("%s %s", val, orderType))
	}

	return service.SqlParameter{
		Params:  filterParam,
		Limit:   limit,
		Offset:  offset,
		OrderBy: orderBy,
	}
}
