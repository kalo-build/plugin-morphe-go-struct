package models

type ContactInfo struct {
	Email    string
	ID       uint `mandatory`
	PersonID uint
	Person   *Person
}

func (m ContactInfo) GetIDEmail() ContactInfoIDEmail {
	return ContactInfoIDEmail{
		Email: m.Email,
	}
}

func (m ContactInfo) GetIDPrimary() ContactInfoIDPrimary {
	return ContactInfoIDPrimary{
		ID: m.ID,
	}
}
