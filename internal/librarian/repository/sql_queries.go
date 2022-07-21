package repository

const (
	createLibrarianQuery = `INSERT INTO librarians (first_name, last_name, email, password, avatar) 
		VALUES ($1, $2, $3, $4, COALESCE(NULLIF($5, ''), null)) 
		RETURNING librarian_id, first_name, last_name, email, password, avatar, created_at, updated_at`

	findByEmailQuery = `SELECT librarian_id, email, first_name, last_name, avatar, password, created_at, updated_at FROM librarians WHERE email = $1`

	findByIdQuery = `SELECT librarian_id, email, first_name, last_name, avatar, password, created_at, updated_at FROM librarians WHERE librarian_id = $1`

	findAllQuery = `SELECT librarian_id, email, first_name, last_name, avatar, password, created_at, updated_at FROM librarians LIMIT $1 OFFSET $2`

	updateByIdQuery = `UPDATE librarians SET first_name = $2, last_name = $3, email = $4, password = $5, avatar = $6 WHERE librarian_id = $1
		RETURNING librarian_id, first_name, last_name, email, password, avatar, created_at, updated_at`

	deleteByIdQuery = `DELETE FROM librarians WHERE librarian_id = $1`
)
