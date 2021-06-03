package clinic

// Clinic represents the structure of both the dental and vet clinics
type Clinic struct {
	Name         string       `json:"name"`
	State        string       `json:"state"`
	Availability Availability `json:"availability"`
}

// Availability contains the period during which a clinic is available
type Availability struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type DentalClinic struct {
	Name         string       `json:"name"`
	State        string       `json:"stateName"`
	Availability Availability `json:"availability"`
}

type VetClinic struct {
	Name         string       `json:"clinicName"`
	State        string       `json:"stateCode"`
	Availability Availability `json:"opening"`
}

type SearchParams struct {
	Name  string `json:"name"`
	State string `json:"state"`
	From  string `json:"from"`
	To    string `json:"to"`
}
