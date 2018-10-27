package models

// Deletes old session on login for space efficiency
func DeleteSession(userId int) error {
	sqlSession := `
		DELETE FROM sessions
		WHERE userid = $1;`

	_, err := db.Query(sqlSession, userId)
	if err != nil {
		return err
	}
	return nil
}

// CreateSession inserts user id, random hex value for
// Fetching user for auth
func CreateSession(hex string, userId int) error {
	sqlSession := `
		INSERT INTO sessions ( hex, userid )
		VALUES( $1, $2);`

	_, err := db.Query(sqlSession, hex, userId)
	if err != nil {
		return err
	}
	return nil
}
