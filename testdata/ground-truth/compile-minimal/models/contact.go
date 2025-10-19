package models

type Contact struct {
	Email string
	ID    uint
	Phone string
}

func (m Contact) GetIDPrimary() ContactIDPrimary {
	return ContactIDPrimary{
		ID: m.ID,
	}
}
