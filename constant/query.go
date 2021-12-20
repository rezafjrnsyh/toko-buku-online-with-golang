package constant

const (
	SCHEMA = `
	CREATE TABLE "mst_books" ("id" bigserial,"title" text,"description" text,"price" bigint,"stock" bigint,"purchase_amount" bigint,PRIMARY KEY ("id"));
	CREATE TABLE "mst_members" ("id" bigserial,"first_name" text,"last_name" text,"email" text,"password" text,"status" bigint,PRIMARY KEY ("id"));
	CREATE TABLE "members_books" ("member_id" bigint,"book_id" bigint,"quantity" bigint,PRIMARY KEY ("member_id","book_id"),CONSTRAINT "fk_members_books_member" FOREIGN KEY ("member_id") REFERENCES "mst_members"("id"),CONSTRAINT "fk_members_books_book" FOREIGN KEY ("book_id") REFERENCES "mst_books"("id"))`

	INSERT_BOOK     = `INSERT INTO mst_books(title, description, price, stock) VALUES($1,$2,$3,$4) RETURNING id;`
	FIND_BOOKS      = `SELECT id, title, description, price, stock FROM mst_books`
	FIND_BOOK_BY_ID = `SELECT id, title, description, price, stock FROM mst_books WHERE id=$1`
	UPDATE_BOOK     = `UPDATE mst_books SET title = :title, description = :description, price = :price, stock = :price WHERE id = :id`
	DELETE_BOOK     = `DELETE FROM mst_books WHERE id = $1`

	INSERT_MEMBER        = `INSERT INTO mst_members(first_name, last_name, email, password,status) VALUES($1,$2,$3,$4,$5) RETURNING id;`
	FIND_MEMBER          = `SELECT id, first_name, last_name, email, password, status FROM mst_members WHERE id=$1`
	UPDATE_STATUS_MEMBER = `UPDATE mst_members SET first_name=:first_name, last_name=:last_name, email=:email, password=:password, status = :status WHERE email = :email`
)
