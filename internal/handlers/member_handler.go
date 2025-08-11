package handlers

import (
	"fmt"
	"haslaw-be-services/internal/service"
	"haslaw-be-services/internal/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// MemberHandler handles member requests
type MemberHandler struct {
	memberService service.MemberService
}

// NewMemberHandler creates a new member handler
func NewMemberHandler(memberService service.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
	}
}

// GetAll gets all members
func (h *MemberHandler) GetAll(c *gin.Context) {
	members, err := h.memberService.GetAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch members", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Members retrieved successfully", members)
}

// GetByID gets member by ID
func (h *MemberHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid member ID", err.Error())
		return
	}

	member, err := h.memberService.GetByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Member not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Member retrieved successfully", member)
}

// Create creates new member
func (h *MemberHandler) Create(c *gin.Context) {
	var req service.CreateMemberRequest
	
	// Check content type - support both JSON and form-data
	contentType := c.GetHeader("Content-Type")
	
	if strings.Contains(contentType, "application/json") {
		// Handle JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.BadRequestResponse(c, "Invalid request body", err.Error())
			return
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		// Handle form-data request
		req.FullName = c.PostForm("full_name")
		req.TitlePosition = c.PostForm("title_position")
		req.Email = c.PostForm("email")
		req.PhoneNumber = c.PostForm("phone_number")
		req.LinkedIn = c.PostForm("linkedin")
		req.Biography = c.PostForm("biography")
		req.Education = c.PostForm("education")
		
		// Handle array fields
		if practiceFocus := c.PostForm("practice_focus"); practiceFocus != "" {
			req.PracticeFocus = strings.Split(practiceFocus, ",")
		}
		if language := c.PostForm("language"); language != "" {
			req.Language = strings.Split(language, ",")
		}
		
		// Validate required fields
		if req.FullName == "" || req.TitlePosition == "" || req.Email == "" {
			utils.BadRequestResponse(c, "Required fields missing", "full_name, title_position, and email are required")
			return
		}
		
		// Handle file uploads
		if displayImageFile, err := c.FormFile("display_image"); err == nil {
			displayImagePath := fmt.Sprintf("uploads/members/%d_display_%s", time.Now().Unix(), displayImageFile.Filename)
			if err := c.SaveUploadedFile(displayImageFile, displayImagePath); err != nil {
				utils.InternalServerErrorResponse(c, "Failed to save display image", err.Error())
				return
			}
			req.DisplayImage = displayImagePath
		}
		
		if detailImageFile, err := c.FormFile("detail_image"); err == nil {
			detailImagePath := fmt.Sprintf("uploads/members/%d_detail_%s", time.Now().Unix(), detailImageFile.Filename)
			if err := c.SaveUploadedFile(detailImageFile, detailImagePath); err != nil {
				utils.InternalServerErrorResponse(c, "Failed to save detail image", err.Error())
				return
			}
			req.DetailImage = detailImagePath
		}
		
		if businessCardFile, err := c.FormFile("business_card"); err == nil {
			businessCardPath := fmt.Sprintf("uploads/members/%d_business_%s", time.Now().Unix(), businessCardFile.Filename)
			if err := c.SaveUploadedFile(businessCardFile, businessCardPath); err != nil {
				utils.InternalServerErrorResponse(c, "Failed to save business card", err.Error())
				return
			}
			req.BusinessCard = businessCardPath
		}
	} else {
		utils.BadRequestResponse(c, "Unsupported content type", "Use application/json or multipart/form-data")
		return
	}

	member, err := h.memberService.Create(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create member", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Member created successfully", member)
}

// Update updates member
func (h *MemberHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid member ID", err.Error())
		return
	}

	var req service.UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err.Error())
		return
	}

	member, err := h.memberService.Update(uint(id), &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to update member", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Member updated successfully", member)
}

// Delete deletes member
func (h *MemberHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid member ID", err.Error())
		return
	}

	if err := h.memberService.Delete(uint(id)); err != nil {
		utils.BadRequestResponse(c, "Failed to delete member", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Member deleted successfully", nil)
}
