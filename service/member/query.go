package member

const (
	getAllMemberQuery = `SELECT ID,NAME,INFO FROM MEMBER m`
	findByIdQuery     = getAllMemberQuery + ` WHERE id = :1 `
	createMemberQuery = `INSERT INTO MEMBER (NAME, INFO) VALUES (:1, :2) RETURNING ID INTO :3`
	updateMemberQuery = `UPDATE MEMBER SET NAME = :1, INFO = :2 WHERE ID = :3`
	DeleteMemberQuery = `DELETE FROM MEMBER WHERE ID = :1`
)
