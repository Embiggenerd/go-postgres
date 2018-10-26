package models

func DeleteSession(userId int) error {
	sqlSession := `
		DELETE FROM sessions
		WHERE userid = $1`

	_, err := db.Query(sqlSession, userId)
	if err != nil {
		return err
	}
	return nil
}

func CreateSession(hex string, userId int) error {
	sqlSession := `
		INSERT INTO sessions ( hex, userid )
		VALUES( $1, $2)`

	_, err := db.Query(sqlSession, hex, userId)
	if err != nil {
		return err
	}
	return nil
}
