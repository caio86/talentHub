package talenthub

type Education struct {
	ID          int
	CandidateID int
	Institution string
	Course      string
	Level       string
}

func (e *Education) Validate() error {
	if e.Institution == "" {
		return Errorf(EINVALID, "institution required")
	}
	if e.Course == "" {
		return Errorf(EINVALID, "course required")
	}
	if e.Level == "" {
		return Errorf(EINVALID, "level required")
	}

	return nil
}
