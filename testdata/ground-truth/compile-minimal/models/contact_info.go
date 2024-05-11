package models

type ContactInfo struct {
	Email string
	ID    uint `mandatory`
}

func (m ContactInfoIDEmail) GetIDEmail() ContactInfoIDEmail {
	return ContactInfoIDEmail{
		Email: m.Email,
	}
}

func (m ContactInfoIDPrimary) GetIDPrimary() ContactInfoIDPrimary {
	return ContactInfoIDPrimary{
		ID: m.ID,
	}
}
