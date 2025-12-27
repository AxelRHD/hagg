package storesqlite

import "github.com/nullism/bqb"

func qCreateUser(uid string) *bqb.Query {
	return bqb.New(`
		INSERT INTO users (uid)
		VALUES (?)
		RETURNING
			id,
			uid,
			description,
			created_at,
			updated_at
	`, uid)
}

func qUserByUID(uid string) *bqb.Query {
	sel := bqb.New("SELECT id, uid, description, created_at, updated_at FROM users")
	sel.Concat("\nWHERE uid = ?", uid)

	return sel
}

func qAllUsers() *bqb.Query {
	sel := bqb.New("SELECT id, uid, description, created_at, updated_at FROM users")

	return sel
}
