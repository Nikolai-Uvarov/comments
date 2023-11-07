//объектная модель сервиса комментариев
package obj

// Комментарий к новости или другому комментарию (ненулевой PostID или CommentID)
type Comment struct {
	ID        int       //ID комментария
	PostID    int       //ID поста (новости) к которой осавлен комментарий
	CommentID int       //ID комментария, к которому привязан данный коммент, если этот комментарий является ответом на другой комментарий
	Text      string    //Текст комментария
	Answers   []Comment //Ответы на данный комментарий
}

type DB interface {
	SaveComment(Comment) error
	GetComments(int) ([]Comment, error)
}

func (c *Comment) BuildCommentTree(list []Comment) {
	for _, x := range list {
		if x.CommentID == c.ID {
			x.BuildCommentTree(list)
			c.Answers = append(c.Answers, x)
		}
	}
}
