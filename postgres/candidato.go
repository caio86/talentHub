package postgres

import (
	"context"
	"net/mail"

	talenthub "github.com/caio86/talentHub"
	"github.com/caio86/talentHub/postgres/repository"
)

var _ talenthub.CandidatoService = (*CandidatoService)(nil)

type CandidatoService struct {
	db   *DB
	repo *repository.Queries
}

func NewCandidatoService(db *DB) *CandidatoService {
	return &CandidatoService{db: db, repo: repository.New(db.conn)}
}

func (s *CandidatoService) FindCandidatoByID(ctx context.Context, id int) (*talenthub.Candidato, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	result, err := repoTx.GetCandidateByID(ctx, int32(id))
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "candidato not found")
	}

	candidate := &talenthub.Candidato{
		ID:       int(result.ID),
		Name:     result.Name,
		Email:    result.Email,
		Password: result.Password,
	}

	// Checking for nils
	if result.Phone != nil {
		candidate.Phone = *result.Phone
	}
	if result.Address != nil {
		candidate.Address = *result.Address
	}
	if result.Linkedin != nil {
		candidate.Linkedin = *result.Linkedin
	}
	if result.ResumeUrl != nil {
		candidate.ResumeLink = *result.ResumeUrl
	}

	s.attachCandidateAssociations(ctx, tx, candidate)

	return candidate, nil
}

func (s *CandidatoService) FindCandidatoByEmail(ctx context.Context, email string) (*talenthub.Candidato, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, talenthub.Errorf(talenthub.EINVALID, "email invalid")
	}

	result, err := repoTx.GetCandidateByEmail(ctx, email)
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "candidato not found")
	}

	candidate := &talenthub.Candidato{
		ID:       int(result.ID),
		Name:     result.Name,
		Email:    result.Email,
		Password: result.Password,
	}

	// Checking for nils
	if result.Phone != nil {
		candidate.Phone = *result.Phone
	}
	if result.Address != nil {
		candidate.Address = *result.Address
	}
	if result.Linkedin != nil {
		candidate.Linkedin = *result.Linkedin
	}
	if result.ResumeUrl != nil {
		candidate.ResumeLink = *result.ResumeUrl
	}

	s.attachCandidateAssociations(ctx, tx, candidate)

	return candidate, nil
}

func (s *CandidatoService) FindCandidatos(ctx context.Context, filter talenthub.CandidatoFilter) ([]*talenthub.Candidato, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	arg := repository.ListCandidatesParams{
		Limit:  filter.Limit,
		Offset: filter.Offset * filter.Limit,
	}

	if arg.Limit <= 0 {
		return s.findAllCandidatos(ctx, tx)
	} else {
		return s.findCandidatos(ctx, tx, arg)
	}
}

func (s *CandidatoService) findAllCandidatos(ctx context.Context, tx *Tx) ([]*talenthub.Candidato, int, error) {
	repoTx := s.repo.WithTx(tx)

	total, err := repoTx.CountCandidates(ctx)
	if err != nil {
		return nil, 0, err
	}

	candidatos, err := repoTx.ListAllCandidates(ctx)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*talenthub.Candidato, len(candidatos))
	for i, v := range candidatos {
		res[i] = &talenthub.Candidato{
			ID:       int(v.ID),
			Name:     v.Name,
			Email:    v.Email,
			Password: v.Password,
		}

		// Checking for nils
		if v.Phone != nil {
			res[i].Phone = *v.Phone
		}
		if v.Address != nil {
			res[i].Address = *v.Address
		}
		if v.Linkedin != nil {
			res[i].Linkedin = *v.Linkedin
		}
		if v.ResumeUrl != nil {
			res[i].ResumeLink = *v.ResumeUrl
		}

		s.attachCandidateAssociations(ctx, tx, res[i])
	}

	return res, int(total), nil
}

func (s *CandidatoService) findCandidatos(ctx context.Context, tx *Tx, arg repository.ListCandidatesParams) ([]*talenthub.Candidato, int, error) {
	repoTx := s.repo.WithTx(tx)

	total, err := repoTx.CountCandidates(ctx)
	if err != nil {
		return nil, 0, err
	}

	candidatos, err := repoTx.ListCandidates(ctx, arg)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*talenthub.Candidato, len(candidatos))
	for i, v := range candidatos {
		res[i] = &talenthub.Candidato{
			ID:       int(v.ID),
			Name:     v.Name,
			Email:    v.Email,
			Password: v.Password,
		}

		// Checking for nils
		if v.Phone != nil {
			res[i].Phone = *v.Phone
		}
		if v.Address != nil {
			res[i].Address = *v.Address
		}
		if v.Linkedin != nil {
			res[i].Linkedin = *v.Linkedin
		}
		if v.ResumeUrl != nil {
			res[i].ResumeLink = *v.ResumeUrl
		}

		s.attachCandidateAssociations(ctx, tx, res[i])
	}

	return res, int(total), nil
}

func (s *CandidatoService) CreateCandidato(ctx context.Context, candidato *talenthub.Candidato) (*talenthub.Candidato, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	if _, err := mail.ParseAddress(candidato.Email); err != nil {
		return nil, talenthub.Errorf(talenthub.EINVALID, "email invalid")
	}

	arg := repository.CreateCandidateParams{
		Name:      candidato.Name,
		Email:     candidato.Email,
		Password:  candidato.Password,
		Phone:     &candidato.Phone,
		Address:   &candidato.Name,
		Linkedin:  &candidato.Name,
		ResumeUrl: &candidato.Name,
	}
	newCandidate, err := repoTx.CreateCandidate(ctx, arg)
	if err != nil {
		switch err.Error() {
		case "ERROR: duplicate key value violates unique constraint \"candidates_email_key\" (SQLSTATE 23505)":
			return nil, talenthub.Errorf(talenthub.ECONFLICT, "email already exists")
		default:
			return nil, err
		}
	}

	res := &talenthub.Candidato{
		ID:       int(newCandidate.ID),
		Name:     newCandidate.Name,
		Email:    newCandidate.Email,
		Password: newCandidate.Password,
	}

	// Checking for nils
	if newCandidate.Phone != nil {
		res.Phone = *newCandidate.Phone
	}
	if newCandidate.Address != nil {
		res.Address = *newCandidate.Address
	}
	if newCandidate.Linkedin != nil {
		res.Linkedin = *newCandidate.Linkedin
	}
	if newCandidate.ResumeUrl != nil {
		res.ResumeLink = *newCandidate.ResumeUrl
	}

	// Education
	for _, v := range candidato.Education {
		argEdu := repository.AddCandidateEducationParams{
			CandidateID: int32(newCandidate.ID),
			Institution: v.Institution,
			Course:      v.Course,
			Level:       v.Level,
		}

		newEdu, err := repoTx.AddCandidateEducation(ctx, argEdu)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "erro interno: %s", err)
		}

		res.Education = append(res.Education, &talenthub.Education{
			ID:          int(newEdu.ID),
			CandidateID: int(newCandidate.ID),
			Institution: newEdu.Institution,
			Course:      newEdu.Course,
			Level:       newEdu.Level,
		})
	}

	// Experiences
	for _, v := range candidato.Experiences {
		argExp := repository.AddCandidateExperienceParams{
			CandidateID: int32(newCandidate.ID),
			Company:     v.Company,
			Role:        v.Role,
			Years:       int32(v.Years),
		}

		newExp, err := repoTx.AddCandidateExperience(ctx, argExp)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "erro interno: %s", err)
		}

		res.Experiences = append(res.Experiences, &talenthub.Experience{
			ID:          int(newExp.ID),
			CandidateID: int(newCandidate.ID),
			Company:     newExp.Company,
			Role:        newExp.Role,
			Years:       int(newExp.Years),
		})
	}

	// Skills
	for _, v := range candidato.Skills {
		argSkill := repository.AddCandidateSkillParams{
			CandidateID: int32(newCandidate.ID),
			Skill:       v,
		}

		err := repoTx.AddCandidateSkill(ctx, argSkill)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "erro interno: %s", err)
		}

	}
	res.Skills = candidato.Skills

	// Interests
	for _, v := range candidato.Interests {
		argInterest := repository.AddCandidateInterestParams{
			CandidateID: int32(newCandidate.ID),
			Interest:    v,
		}

		err := repoTx.AddCandidateInterest(ctx, argInterest)
		if err != nil {
			return nil, talenthub.Errorf(talenthub.EINTERNAL, "erro interno: %s", err)
		}
	}
	res.Interests = candidato.Interests

	tx.Commit(ctx)

	return res, nil
}

func (s *CandidatoService) UpdateCandidato(ctx context.Context, id int, upd talenthub.CandidatoUpdate) (*talenthub.Candidato, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	repoTx := s.repo.WithTx(tx)

	candidato, err := repoTx.GetCandidateByID(ctx, int32(id))
	if err != nil {
		return nil, talenthub.Errorf(talenthub.ENOTFOUND, "candidato not found")
	}

	arg := repository.UpdateCandidateParams{
		ID:        int32(id),
		Name:      candidato.Name,
		Phone:     candidato.Phone,
		Address:   candidato.Address,
		Linkedin:  candidato.Linkedin,
		ResumeUrl: candidato.ResumeUrl,
	}

	// Apply updates
	if upd.Name != nil {
		arg.Name = *upd.Name
	}
	if upd.Phone != nil {
		arg.Phone = upd.Phone
	}
	if upd.Address != nil {
		arg.Address = upd.Address
	}
	if upd.Linkedin != nil {
		arg.Linkedin = upd.Linkedin
	}
	if upd.ResumeLink != nil {
		arg.ResumeUrl = upd.ResumeLink
	}

	updated, err := repoTx.UpdateCandidate(ctx, arg)
	if err != nil {
		return nil, talenthub.Errorf(talenthub.EINTERNAL, "internal error: %v", err)
	}

	res := &talenthub.Candidato{
		ID:         int(updated.ID),
		Name:       updated.Name,
		Phone:      *updated.Phone,
		Address:    *updated.Address,
		Linkedin:   *updated.Linkedin,
		ResumeLink: *updated.ResumeUrl,
	}

	tx.Commit(ctx)

	return res, nil
}

func (s *CandidatoService) attachCandidateAssociations(ctx context.Context, tx *Tx, candidate *talenthub.Candidato) {
	repoTx := s.repo.WithTx(tx)

	// educations := make([]*talenthub.Education)
	if resEducation, err := repoTx.GetCandidateEducations(ctx, int32(candidate.ID)); err != nil {
		candidate.Education = make([]*talenthub.Education, 0)
	} else {
		candidate.Education = make([]*talenthub.Education, len(resEducation))
		for k, v := range resEducation {
			candidate.Education[k] = &talenthub.Education{
				CandidateID: int(v.CandidateID),
				ID:          int(v.ID),
				Institution: v.Institution,
				Course:      v.Course,
				Level:       v.Level,
			}
		}
	}

	// experiences := make([]*talenthub.Experience)
	if resExperiences, err := repoTx.GetCandidateExperiences(ctx, int32(candidate.ID)); err != nil {
		candidate.Experiences = make([]*talenthub.Experience, 0)
	} else {
		candidate.Experiences = make([]*talenthub.Experience, len(resExperiences))
		for k, v := range resExperiences {
			candidate.Experiences[k] = &talenthub.Experience{
				CandidateID: int(v.CandidateID),
				ID:          int(v.ID),
				Company:     v.Company,
				Role:        v.Role,
				Years:       int(v.Years),
			}
		}
	}

	// Candiate Skills
	if resSkills, err := repoTx.GetCandidateSkills(ctx, int32(candidate.ID)); err != nil {
		candidate.Skills = make([]string, 0)
	} else {
		candidate.Skills = make([]string, len(resSkills))
		for k, v := range resSkills {
			candidate.Skills[k] = v.Skill
		}
	}

	// Candidate Interests
	if resInterests, err := repoTx.GetCandidateInterests(ctx, int32(candidate.ID)); err != nil {
		candidate.Interests = make([]string, 0)
	} else {
		candidate.Interests = make([]string, len(resInterests))
		for k, v := range resInterests {
			candidate.Interests[k] = v.Interest
		}
	}
}
