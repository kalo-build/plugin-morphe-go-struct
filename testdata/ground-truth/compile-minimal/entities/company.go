package entities

type Company struct {
	ID        uint `morphe:"immutable"`
	Name      string
	TaxID     string
	PersonIDs []uint
	People    []Person
}

func (e Company) GetIDPrimary() CompanyIDPrimary {
	return CompanyIDPrimary{
		ID: e.ID,
	}
}
