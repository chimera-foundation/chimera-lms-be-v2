package http

import (
	"net/http"

	"github.com/chimera-foundation/chimera-lms-be-v2/internal/features/attachment/service"
	"github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/auth"
	response "github.com/chimera-foundation/chimera-lms-be-v2/internal/shared/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AttachmentHandler struct {
	attachmentService service.AttachmentService
	log               *logrus.Logger
}

func NewAttachmentHandler(attachmentService service.AttachmentService, log *logrus.Logger) *AttachmentHandler {
	return &AttachmentHandler{
		attachmentService: attachmentService,
		log:               log,
	}
}

func (h *AttachmentHandler) ProtectedRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/assessment/{assessmentID}", h.UploadAssessmentAttachment)
	r.Post("/submission/{submissionID}", h.UploadSubmissionAttachment)
	r.Delete("/{id}", h.DeleteAttachment)
	r.Get("/assessment/{assessmentID}", h.ListAssessmentAttachments)
	r.Get("/submission/{submissionID}", h.ListSubmissionAttachments)

	return r
}

func (h *AttachmentHandler) UploadAssessmentAttachment(w http.ResponseWriter, r *http.Request) {
	assessmentIDStr := chi.URLParam(r, "assessmentID")
	assessmentID, err := uuid.Parse(assessmentIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid assessment ID")
		return
	}

	h.handleUpload(w, r, &assessmentID, nil)
}

func (h *AttachmentHandler) UploadSubmissionAttachment(w http.ResponseWriter, r *http.Request) {
	submissionIDStr := chi.URLParam(r, "submissionID")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid submission ID")
		return
	}

	h.handleUpload(w, r, nil, &submissionID)
}

func (h *AttachmentHandler) handleUpload(w http.ResponseWriter, r *http.Request, assessmentID *uuid.UUID, submissionID *uuid.UUID) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	orgID, ok := auth.GetOrgID(r.Context())
	if !ok {
		response.Unauthorized(w, "Organization not found in context")
		return
	}

	if err := r.ParseMultipartForm(26 << 20); err != nil {
		h.log.WithError(err).Warn("failed to parse multipart form")
		response.BadRequest(w, "File too large or invalid form data")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.log.WithError(err).Warn("file field missing from form")
		response.BadRequest(w, "File is required (field name: 'file')")
		return
	}
	defer file.Close()

	req := service.UploadRequest{
		OrgID:        orgID,
		UploadedBy:   userID,
		AssessmentID: assessmentID,
		SubmissionID: submissionID,
		FileName:     header.Filename,
		FileSize:     header.Size,
		MIMEType:     header.Header.Get("Content-Type"),
		File:         file,
	}

	result, err := h.attachmentService.Upload(r.Context(), req)
	if err != nil {
		h.log.WithError(err).Error("failed to upload attachment")
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, result)
}

func (h *AttachmentHandler) DeleteAttachment(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "User not found in context")
		return
	}

	attachmentIDStr := chi.URLParam(r, "id")
	attachmentID, err := uuid.Parse(attachmentIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid attachment ID")
		return
	}

	if err := h.attachmentService.Delete(r.Context(), userID, attachmentID); err != nil {
		h.log.WithError(err).WithField("attachment_id", attachmentID).Error("failed to delete attachment")
		response.BadRequest(w, err.Error())
		return
	}

	response.NoContent(w)
}

func (h *AttachmentHandler) ListAssessmentAttachments(w http.ResponseWriter, r *http.Request) {
	assessmentIDStr := chi.URLParam(r, "assessmentID")
	assessmentID, err := uuid.Parse(assessmentIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid assessment ID")
		return
	}

	attachments, err := h.attachmentService.ListByAssessment(r.Context(), assessmentID)
	if err != nil {
		h.log.WithError(err).Error("failed to list assessment attachments")
		response.InternalServerError(w, err.Error())
		return
	}

	response.OK(w, attachments)
}

func (h *AttachmentHandler) ListSubmissionAttachments(w http.ResponseWriter, r *http.Request) {
	submissionIDStr := chi.URLParam(r, "submissionID")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid submission ID")
		return
	}

	attachments, err := h.attachmentService.ListBySubmission(r.Context(), submissionID)
	if err != nil {
		h.log.WithError(err).Error("failed to list submission attachments")
		response.InternalServerError(w, err.Error())
		return
	}

	response.OK(w, attachments)
}
