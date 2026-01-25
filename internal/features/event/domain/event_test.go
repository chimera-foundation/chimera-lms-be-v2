package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEvent_Validate(t *testing.T) {
	// Setup reusable variables
	orgID := uuid.New()
	userID := uuid.New()
	now := time.Now()
	later := now.Add(time.Hour)
	emptyStr := ""

	tests := []struct {
		name    string
		event   *Event
		wantErr bool
	}{
		// --- 1. Basic Requirements ---
		{
			name:    "Failure: Missing OrganizationID",
			event:   &Event{Title: "No Org", EventType: Vanilla},
			wantErr: true,
		},
		{
			name:    "Failure: Empty Title",
			event:   &Event{OrganizationID: orgID, Title: "   ", EventType: Vanilla},
			wantErr: true,
		},

		// --- 2. EventType Specific Logic ---
		{
			name: "Failure: Session missing times",
			event: &Event{
				OrganizationID: orgID,
				Title:          "Lecture",
				EventType:      Session,
				// StartAt and EndAt are nil
			},
			wantErr: true,
		},
		{
			name: "Failure: Session End before Start",
			event: &Event{
				OrganizationID: orgID,
				Title:          "Time Warp",
				EventType:      Session,
				StartAt:        &later,
				EndAt:          &now,
			},
			wantErr: true,
		},
		{
			name: "Failure: Deadline missing StartAt",
			event: &Event{
				OrganizationID: orgID,
				Title:          "Assignment Due",
				EventType:      Deadline,
				StartAt:        nil,
			},
			wantErr: true,
		},
		{
			name: "Failure: Announcement with empty ImageURL pointer",
			event: &Event{
				OrganizationID: orgID,
				Title:          "News",
				EventType:      Announcement,
				ImageURL:       &emptyStr,
			},
			wantErr: true,
		},

		// --- 3. Scope Logic ---
		{
			name: "Failure: Section scope missing SectionID",
			event: &Event{
				OrganizationID: orgID,
				Title:          "Class Meeting",
				EventType:      Vanilla,
				Scope:          ScopeSection,
				SectionID:      nil,
			},
			wantErr: true,
		},
		{
			name: "Failure: Cohort scope missing CohortID",
			event: &Event{
				OrganizationID: orgID,
				Title:          "Cohort Meetup",
				EventType:      Vanilla,
				Scope:          ScopeCohort,
				CohortID:       nil,
			},
			wantErr: true,
		},
		{
			name: "Failure: Personal scope missing UserID",
			event: &Event{
				OrganizationID: orgID,
				Title:          "Private Note",
				EventType:      Vanilla,
				Scope:          ScopePersonal,
				UserID:         nil,
			},
			wantErr: true,
		},

		// --- 4. Source/Linking Logic ---
		{
			name: "Failure: SourceID present but SourceType missing",
			event: &Event{
				OrganizationID: orgID,
				Title:          "Linked Event",
				EventType:      Vanilla,
				SourceID:       &orgID, // reusing a UUID for test
				SourceType:     nil,
			},
			wantErr: true,
		},

		// --- 5. Happy Paths ---
		{
			name: "Success: Valid Global Session",
			event: &Event{
				OrganizationID: orgID,
				Title:          "Valid Workshop",
				EventType:      Session,
				StartAt:        &now,
				EndAt:          &later,
				Scope:          ScopeGlobal,
			},
			wantErr: false,
		},
		{
			name: "Success: Valid Personal Deadline",
			event: &Event{
				OrganizationID: orgID,
				Title:          "My Due Date",
				EventType:      Deadline,
				StartAt:        &now,
				Scope:          ScopePersonal,
				UserID:         &userID,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Event.Validate() [%s] error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}