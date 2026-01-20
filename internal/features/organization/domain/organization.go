package domain

import (
	u "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/user/domain"
	ap "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/academic_period/domain"
	p "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/program/domain"
	c "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/course/domain"
	co "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/cohort/domain"
	e "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/education_level/domain"
	s "github.com/chimera-foundation/chimera-lms-be-v2/internal/features/subject/domain"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type OrgType int

const (
	University OrgType = iota
	HighSchool
	MiddleSchool
	GradeSchool
)

type Organization struct {
	shared.Base
	
	Name string `validate:"required,min=3,max=100"`
	Slug string `validate:"required,alphanum"`
	Type OrgType `validate:"required"`
	Address string `validate:"required,min=10,max=200"`
	Users []u.User 
	AcademicPeriods []ap.AcademicPeriod
	Programs []p.Program
	Courses []c.Course
	Cohorts []co.Cohort
	EducationLevels []e.EducationLevel
	Subjects []s.Subject

	IsActive bool
}