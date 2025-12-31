package storesqlite

import "github.com/nullism/bqb"

func qCreateUser(uid, displayName string) *bqb.Query {
	return bqb.New(`
		INSERT INTO users (uid, display_name)
		VALUES (?, ?)
		RETURNING
			id,
			uid,
			display_name,
			last_name,
			first_name,
			created_at,
			updated_at
	`, uid, displayName)
}

func qUserSelector() *bqb.Query {
	return bqb.New("SELECT id, uid, display_name, last_name, first_name, created_at, updated_at FROM users")
}

func qUserByUID(uid string) *bqb.Query {
	sel := qUserSelector()
	sel.Concat("\nWHERE uid = ?", uid)

	return sel
}

func qUserByDisplayName(displayName string) *bqb.Query {
	sel := qUserSelector()
	sel.Concat("\nWHERE display_name = ?", displayName)

	return sel
}

func qAllUsers() *bqb.Query {
	return qUserSelector()
}
