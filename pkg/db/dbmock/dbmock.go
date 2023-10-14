// Пакет in-mem database для разработки и тестов
package dbmock

import (
	"comments/pkg/db/obj"
	"sort"
)

type DB struct {
	comments []obj.Comment  
	nextid int
}

// создает новое подключение к БД
func New() *DB {
	db := new(DB)
	db.nextid = 1
	return db
}

// Сохраняет комментарий, представленный объектом obj.Comment, в БД
func (db *DB) SaveComment(c obj.Comment) error {
	c.ID = db.nextid
	db.comments = append(db.comments, c)
	db.nextid++
	return nil
}

// Возвращает массив комментариев по ID новости
func (db *DB) GetComments(id int) ([]obj.Comment, error) {
	
	var comments []obj.Comment

	for _,c := range db.comments{
		if c.PostID==id{
			c.BuildCommentTree(db.comments)
			comments = append(comments, c)
		}
	}
	
	sort.Slice(comments, func(i, j int) bool { return comments[i].ID > comments[j].ID })

	return comments, nil
}

