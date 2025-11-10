package member

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"oracle.com/oracle/my-go-oracle-app/api"
	"oracle.com/oracle/my-go-oracle-app/pkg/constants"
	"oracle.com/oracle/my-go-oracle-app/pkg/helpers"
	"oracle.com/oracle/my-go-oracle-app/pkg/response"
	"oracle.com/oracle/my-go-oracle-app/pkg/servicehelper"
	"oracle.com/oracle/my-go-oracle-app/service"
	entity "oracle.com/oracle/my-go-oracle-app/service/member"
)

var (
	memberService api.MemberService
)

const (
	ErrParseUrlParamMsg = "Parse Url Param Failed. %v"
	ErrCreateDataMsg    = "Create Data Failed. %+v"
	ErrParseValidateMsg = "Failed to Parse and Validate. err=%v"
)

func Init(service api.MemberService) {
	memberService = service
}

var variableFilterMapping = map[string]service.FilterParam{
	"name":        {Field: "UPPER(M.NAME)", Operand: constants.LIKE},
	"address":     {Field: "JSON_VALUE(INFO, '$.address')", Operand: constants.LIKE},
	"ageStart":    {Field: "JSON_VALUE(INFO, '$.age')", Operand: constants.GREATER_THAN_EQUAL},
	"ageEnd":      {Field: "JSON_VALUE(INFO, '$.age')", Operand: constants.LESS_THAN_EQUAL},
	"salaryStart": {Field: "JSON_VALUE(INFO, '$.salary')", Operand: constants.GREATER_THAN_EQUAL},
	"salaryEnd":   {Field: "JSON_VALUE(INFO, '$.salary')", Operand: constants.LESS_THAN_EQUAL},
}

var variableOrderMapping = map[string]string{
	"name": "M.NAME",
}

// GetMemberById : HTTP Handler for Get Member by Id
// @Summary Get Member by Id
// @Description GetMemberById handles request for Get Member by Id
// @Tags Member
// @Accept json
// @Produce json
// @Param Accept-Language header string true "accept language" default(id)
// @Param id path string true "id of Member"
// @Success 200 {object} response.Response{data=entity.MemberResponse} "Success Response"
// @Failure 400 "Bad Request"
// @Failure 500 "InternalServerError"
// @Router /members/{id} [GET]
// GetMemberById
func GetMemberById(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.Render(w, r)

	var (
		err    error
		result entity.MemberResponse
	)

	id, err := helpers.GetUrlPathInt64(r, "id")
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf(ErrParseUrlParamMsg, err))
		resp.SetError(err, http.StatusBadRequest)
		return
	}

	result, err = memberService.FindById(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.WarnContext(r.Context(), fmt.Sprintf("Not Found. err=%v", err), slog.Int64("id", id))
			resp.SetError(fmt.Errorf("DATA_NOT_EXIST"), http.StatusNotFound)
			return
		}

		slog.WarnContext(r.Context(), fmt.Sprintf("Get by Id Failed. err=%v", err), slog.Int64("id", id))
		resp.SetError(err, http.StatusInternalServerError)
		return
	}
	resp.Data = result
}

// GetAllMembers : HTTP Handler for Get All Member
// @Summary Get All  by Id
// @Description GetAllMembers handles request for Get Member by Id
// @Tags Member
// @Accept json
// @Produce json
// @Param Accept-Language header string true "accept language" default(id)
// @Param limit query string false "limit data"
// @Param page query integer false "page data"
// @Param name query string false "name filter"
// @Param address query string false "address filter"
// @Param ageStart query int false "ageStart filter"
// @Param ageEnd query int false "ageEnd filter"
// @Param salaryStart query string false "salaryStart filter"
// @Param salaryEnd query string false "salaryEnd filter"
// @Success 200 {object} response.Response{data=[]entity.MemberResponse} "Success Response"
// @Failure 400 "Bad Request"
// @Failure 500 "InternalServerError"
// @Router /members/ [GET]
// GetAllMembers
func GetAllMembers(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.Render(w, r)

	params := servicehelper.GelSqlParameterFromRequest(r, variableFilterMapping, variableOrderMapping)

	result, page, err := memberService.FindAll(r.Context(), params)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf("Failed. %+v", err))
		resp.SetError(err, http.StatusInternalServerError)
		return
	}
	resp.Data = result
	resp.Pagination = page
}

// CreateMember : HTTP Handler for Create Member
// @Summary Create Member
// @Description CreateMember handles request for creating a new member
// @Tags Member
// @Accept json
// @Produce json
// @Param Accept-Language header string true "accept language" default(id)
// @Param member body entity.MemberRequest true "Member Request Body"
// @Success 200 {object} response.Response{data=entity.MemberResponse} "Success Response"
// @Failure 400 "Bad Request"
// @Failure 500 "InternalServerError"
// @Router /members/ [POST]
// CreateMember
func CreateMember(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.Render(w, r)

	var (
		err error
		req entity.MemberRequest
	)

	err = helpers.ParseBodyAndValidate(r, &req)
	if err != nil {
		resp.SetError(err, http.StatusBadRequest)
		slog.WarnContext(r.Context(), fmt.Sprintf(ErrParseValidateMsg, err))
		return
	}

	result, err := memberService.CreateMember(r.Context(), &req)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		slog.WarnContext(r.Context(), fmt.Sprintf(ErrCreateDataMsg, err),
			slog.Any("request", req))
		return
	}

	resp.Data = result

}

// UpdateMember : HTTP Handler for Update Member
// @Summary Update Member
// @Description UpdateMember handles request for updating a member
// @Tags Member
// @Accept json
// @Produce json
// @Param Accept-Language header string true "accept language" default(id)
// @Param id path string true "id of Member"
// @Param member body entity.MemberRequest true "Member Request Body"
// @Success 200 {object} response.Response{data=entity.MemberResponse} "Success Response"
// @Failure 400 "Bad Request"
// @Failure 500 "InternalServerError"
// @Router /members/{id} [PUT]
// UpdateMember
func UpdateMember(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.Render(w, r)

	var (
		err error
		req entity.MemberRequest
	)

	id, err := helpers.GetUrlPathInt64(r, "id")
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf(ErrParseUrlParamMsg, err))
		resp.SetError(err, http.StatusBadRequest)
		return
	}

	err = helpers.ParseBodyAndValidate(r, &req)
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf(ErrParseValidateMsg, err))
		resp.SetError(err, http.StatusBadRequest)
		return
	}

	result, err := memberService.UpdateMember(r.Context(), id, &req)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		slog.WarnContext(r.Context(), fmt.Sprintf("failed to update member data: %v", err),
			slog.Any("request", req))
		return
	}

	resp.Data = result

}

// DeleteMember : HTTP Handler for Delete Member
// @Summary Delete Member
// @Description DeleteMember handles request for deleting a member
// @Tags Member
// @Accept json
// @Produce json
// @Param Accept-Language header string true "accept language" default(id)
// @Param id path string true "id of Member"
// @Success 200 {object} response.Response{data=bool} "Success Response"
// @Failure 400 "Bad Request"
// @Failure 500 "InternalServerError"
// @Router /members/{id} [DELETE]
// DeleteMember
func DeleteMember(w http.ResponseWriter, r *http.Request) {
	resp := response.Response{}
	defer resp.Render(w, r)

	var (
		err error
	)

	id, err := helpers.GetUrlPathInt64(r, "id")
	if err != nil {
		slog.WarnContext(r.Context(), fmt.Sprintf(ErrParseUrlParamMsg, err))
		resp.SetError(err, http.StatusBadRequest)
		return
	}

	result, err := memberService.DeleteMember(r.Context(), id)
	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		slog.WarnContext(r.Context(), fmt.Sprintf("failed to delete member data: %v", err),
			slog.Int64("id", id))
		return
	}

	resp.Data = result

}
