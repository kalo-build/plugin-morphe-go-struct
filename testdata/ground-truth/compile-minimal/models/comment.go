package models

type Comment struct {
	ID              uint
	Text            string
	CommentableType string
	CommentableID   string
}

func (m Comment) GetIDPrimary() CommentIDPrimary {
	return CommentIDPrimary{
		ID: m.ID,
	}
}
