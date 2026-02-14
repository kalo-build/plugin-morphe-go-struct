package models

type ContactInfo struct {
	Email            string
	ID               uint
	PersonID         *uint
	Person           *Person
	RelatedContactID *uint
	RelatedContact   *Contact
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
