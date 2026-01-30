package domain

import (
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared"
)

type OrgType string

const (
	University OrgType = "university"
	HighSchool OrgType = "high_school"
	MiddleSchool OrgType = "middle_school"
	GradeSchool OrgType = "grade_school"
)

type Organization struct {
	shared.Base
	
	Name string `validate:"required,min=3,max=100"`
	Slug string `validate:"required,alphanum"`
	Type OrgType `validate:"required"`
	Address string `validate:"required,min=10,max=200"`
	AcademicPeriods []AcademicPeriod

	IsActive bool
	IsSystemOrg *bool
}

func NewOrganization(
	name string,
	slug string,
	org_type OrgType,
	address string,
	academicPeriods []AcademicPeriod,
) *Organization {
	return &Organization{
		Name: name,
		Slug: slug,
		Type: org_type,
		Address: address,
		AcademicPeriods: academicPeriods,
	}
}