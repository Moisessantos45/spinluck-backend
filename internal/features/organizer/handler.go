package organizer

import (
	"log"
	"net/http"
	"spinLuck/internal/shared/models"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

type OrganizerHandler struct {
	service OrganizerService
}

func NewOrganizerHandler(service OrganizerService) *OrganizerHandler {
	return &OrganizerHandler{service: service}
}

func (h *OrganizerHandler) GetDashboardInfoDetailGeneric(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
		return
	}

	infoDetail, err := h.service.GetDashboardInfoDetailGeneric(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching dashboard raffle info detail: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": infoDetail, "message": "Dashboard raffle info detail fetched successfully"})
}

func (h *OrganizerHandler) GetByID(c *gin.Context) {

	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	organizer, err := h.service.GetByID(userID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error fetching organizer: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": organizer, "message": "Organizer fetched successfully"})
}

func (h *OrganizerHandler) Create(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid token: " + err.Error()})
		return
	}

	log.Printf("User ID extracted from JWT: %d", userID)

	var input models.Organizer

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(400, gin.H{"message": "Invalid input: " + err.Error()})
		return
	}

	organizer, err := h.service.Create(input.Name, input.Phone, userID)
	if err != nil {
		c.JSON(500, gin.H{"message": "Error creating organizer: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{"data": organizer, "message": "Organizer created successfully"})
}

func (h *OrganizerHandler) Update(c *gin.Context) {

	_, id, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	var input models.Organizer
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(400, gin.H{"message": "Invalid input: " + err.Error()})
		return
	}

	log.Printf("Updating organizer with ID: %d, Name: %s, Phone: %s", id, input.Name, input.Phone)

	organizer, err := h.service.Update(id, input.Name, input.Phone)
	if err != nil {
		log.Printf("Error updating organizer: %v", err)
		c.JSON(500, gin.H{"message": "Error updating organizer: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": organizer, "message": "Organizer updated successfully"})
}
