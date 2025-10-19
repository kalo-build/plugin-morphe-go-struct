package models

import (
	"github.com/kalo-build/dummy/enums"
)

type Person struct {
	FirstName         string
	ID                uint `morphe:"mandatory"`
	LastName          string
	Nationality       enums.Nationality
	CompanyID         *uint
	Company           *Company
	ContactInfoID     *uint
	ContactInfo       *ContactInfo
	NoteIDs           []uint
	Notes             []Comment
	PersonalContactID *uint
	PersonalContact   *Contact
	WorkContactID     *uint
	WorkContact       *Contact
}

func (m Person) GetIDName() PersonIDName {
	return PersonIDName{
		FirstName: m.FirstName,
		LastName:  m.LastName,
	}
}

func (m Person) GetIDPrimary() PersonIDPrimary {
	return PersonIDPrimary{
		ID: m.ID,
	}
}
