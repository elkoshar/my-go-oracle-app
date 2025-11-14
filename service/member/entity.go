package member

import (
	"database/sql"
	"encoding/json"
	"time"

	entity "oracle.com/oracle/my-go-oracle-app/service"
)

// todo:
// buat contoh untuk type data lain
// integer, booleah (cahar(1)), date/timestamp

type Member struct {
	Name   string           `db:"NAME"`
	Info   string           `db:"INFO"`
	Detail sql.Null[[]byte] `db:"DETAIL"`
	Policy sql.NullString   `db:"POLICY"`
	entity.BaseEntity
}

type MemberDetail struct {
	MemberId        string `json:"memberId"`
	OnboardingStage string `json:"onboardingStage"`
	RiskRating      string `json:"riskRating"`
}

type MemberRequest struct {
	Name   string       `json:"name"`
	Info   MemberInfo   `json:"info"`
	Detail MemberDetail `json:"detail"`
	Policy Policy       `json:"policy"`
}

type Policy struct {
	EffectiveDate  string   `json:"effectiveDate"`
	Status         string   `json:"status"`
	DataCategories []string `json:"dataCategories"`
}

type MemberResponse struct {
	Id          int64        `json:"id"`
	Name        string       `json:"name"`
	Info        MemberInfo   `json:"info"`
	Detail      MemberDetail `json:"detail"`
	Policy      Policy       `json:"policy"`
	CreatedDate time.Time    `json:"createdDate,omitempty"`
	IsDeleted   bool         `json:"isDeleted"`
}

type MemberInfo struct {
	Address Address `json:"address"`
	Salary  int     `json:"salary"`
	Age     int     `json:"age"`
}

type Address struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
}

func (m *Member) ToResponse() MemberResponse {

	var info MemberInfo
	var detail MemberDetail
	var policy Policy
	var isDeleted bool
	json.Unmarshal([]byte(m.Info), &info)
	json.Unmarshal(m.Detail.V, &detail)
	json.Unmarshal([]byte(m.Policy.String), &policy)
	isDeleted = false
	if m.IsDeleted == "1" {
		isDeleted = true
	}
	return MemberResponse{
		Id:          m.BaseEntity.Id,
		Name:        m.Name,
		Info:        info,
		Detail:      detail,
		Policy:      policy,
		CreatedDate: m.CreatedDate,
		IsDeleted:   isDeleted,
	}
}
func (m *MemberRequest) ToEntity(base entity.BaseEntity) Member {
	infoBytes, _ := json.Marshal(m.Info)
	policyBytes, _ := json.Marshal(m.Policy)
	detailBytes, _ := json.Marshal(m.Detail)

	return Member{
		Name:       m.Name,
		Info:       string(infoBytes),
		Detail:     sql.Null[[]byte]{V: detailBytes, Valid: true},
		Policy:     sql.NullString{String: string(policyBytes), Valid: true},
		BaseEntity: base,
	}
}
