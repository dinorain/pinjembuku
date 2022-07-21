package repository

const (
	createOrderQuery = `INSERT INTO orders (user_id, librarian_id, item, status, pickup_schedule) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING order_id, user_id, librarian_id, item, status, pickup_schedule, created_at, updated_at`

	findByIdQuery = `SELECT order_id, user_id, librarian_id, item, status, pickup_schedule, created_at, updated_at FROM orders WHERE order_id = $1`

	findAllQuery = `SELECT order_id, user_id, librarian_id, item, status, pickup_schedule, created_at, updated_at FROM orders LIMIT $1 OFFSET $2`

	findByUserIdQuery = `SELECT order_id, user_id, librarian_id, item, status, pickup_schedule, created_at, updated_at FROM orders WHERE user_id = $1 LIMIT $2 OFFSET $3`

	findAllByLibrarianIdQuery = `SELECT order_id, user_id, librarian_id, item, status, pickup_schedule, created_at, updated_at FROM orders WHERE librarian_id = $1 LIMIT $2 OFFSET $3`

	findAllByUserIdLibrarianIDQuery = `SELECT order_id, user_id, librarian_id, item, status, pickup_schedule, created_at, updated_at FROM orders WHERE user_id = $1 AND librarian_id = $2 LIMIT $3 OFFSET $4`

	updateByIdQuery = `UPDATE orders SET user_id = $2, librarian_id = $3, item = $4, status = $5, pickup_schedule = $6 WHERE order_id = $1
		RETURNING order_id, user_id, librarian_id, item, status, pickup_schedule, created_at, updated_at`

	deleteByIdQuery = `DELETE FROM orders WHERE order_id = $1`
)
