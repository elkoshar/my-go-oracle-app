package member

const (
	getAllMemberQuery = `SELECT ID,NAME,INFO,DETAIL,POLICY FROM MEMBER m`
	findByIdQuery     = getAllMemberQuery + ` WHERE id = :1 `
	createMemberQuery = `INSERT INTO MEMBER (NAME, INFO) VALUES (:1, :2) RETURNING ID INTO :3`
	updateMemberQuery = `UPDATE MEMBER SET NAME = :1, INFO = :2, DETAIL = :3, POLICY = :4, UPDATED_DATE = :5, IS_DELETED = :6 WHERE ID = :7`
	DeleteMemberQuery = `DELETE FROM MEMBER WHERE ID = :1`
)
