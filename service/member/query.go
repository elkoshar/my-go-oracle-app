package member

const (
	getAllMemberQuery = `SELECT ID,NAME,INFO FROM MEMBER m`
	findByIdQuery     = getAllMemberQuery + ` WHERE id = :1 `
)
