package handlers

import (
	"haslaw-be-services/internal/service"
	"haslaw-be-services/internal/utils"
	"net/http"
	"strconv"

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
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body", err.Error())
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
