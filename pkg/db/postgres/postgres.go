// Пакет для работы приложения с sql DB postgres
package postgres

import (
	"comments/pkg/db/obj"
	"context"
	"log"
	"os"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	DB  *pgxpool.Pool
	ctx context.Context
}

// создает новое подключение к БД
func New() *DB {
	db := new(DB)
	db.ctx = context.Background()
	// Подключение к БД
	dbpass := os.Getenv("dbpass")
	var err error
	db.DB, err = pgxpool.Connect(db.ctx, "postgres://postgres:"+dbpass+"@192.168.1.35:5432/comments")

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return db
}

// Сохраняет комментарий, представленный объектом obj.Comment, в БД
func (db *DB) SaveComment(c obj.Comment) error {
	_, err := db.DB.Exec(db.ctx,
		`INSERT INTO comments (postid, commentid, content) 
		VALUES (($1), ($2), ($3))`,
		c.PostID, c.CommentID, c.Text)

	if err != nil {
		return err
	}

	return nil
}

// Возвращает массив комментариев по ID новости, содержащий вложенные комментарии
func (db *DB) GetComments(id int) ([]obj.Comment, error) {

	allComments, err:= db.GetAllComments()

	if err != nil {
		return nil, err
	}

	var comments []obj.Comment

	for _, c := range allComments {
		if c.PostID == id {
			c.BuildCommentTree(allComments)
			comments = append(comments, c)
		}
	}

	return comments, nil
}

// Возвращает массив всех комментариев из базы
func (db *DB) GetAllComments() ([]obj.Comment, error) {

	rows, err := db.DB.Query(db.ctx, `SELECT * FROM comments ORDER BY ID;`) 

	if err != nil {
		return nil, err
	}

	var allComments []obj.Comment

	for rows.Next() {
		var c obj.Comment

		err = rows.Scan(
			&c.ID,
			&c.PostID,
			&c.CommentID,
			&c.Text,
			&c.Censored)

		if err != nil {
			return nil, err
		}
		allComments = append(allComments, c)
	}

	return allComments, rows.Err()
}

//Записывает в базу признак цензуры по id комментария.
func (db *DB) SetCensored(id int) error {

	_, err := db.DB.Exec(db.ctx,
		`UPDATE comments SET censored = true WHERE id = ($1);`, id)

	if err != nil {
		return err
	}

	return nil
}