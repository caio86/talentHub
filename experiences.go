package talenthub

type Experience struct {
	ID          int
	CandidateID int
	Company     string
	Role        string
	Years       int
}

func (e *Experience) Valdiate() error {
	if e.Company == "" {
		return Errorf(EINVALID, "company required")
	}
	if e.Role == "" {
		return Errorf(EINVALID, "role required")
	}

	return nil
}
