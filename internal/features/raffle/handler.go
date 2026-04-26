package raffle

import (
	"log"
	"net/http"
	"spinLuck/internal/shared/models"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

type RaffleHandler struct {
	service RaffleService
}

func NewRaffleHandler(service RaffleService) *RaffleHandler {
	return &RaffleHandler{service: service}
}

func (h *RaffleHandler) GetAllStatus(c *gin.Context) {
	statuses, err := h.service.GetAllStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching raffle statuses: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": statuses, "message": "Raffle statuses fetched successfully"})
}

func (h *RaffleHandler) GetAll(c *gin.Context) {
	raffles, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching raffles: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": raffles, "message": "Raffles fetched successfully"})
}

func (h *RaffleHandler) GetAllInfoGeneric(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
		return
	}

	raffles, err := h.service.GetAllInfoGeneric(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching raffles info: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": raffles, "message": "Raffles info fetched successfully"})
}

func (h *RaffleHandler) GetAllRecentInfoGeneric(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
		return
	}

	raffles, err := h.service.GetAllRecentInfoGeneric(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching recent raffles info: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": raffles, "message": "Recent raffles info fetched successfully"})
}

func (h *RaffleHandler) GetByID(c *gin.Context) {
	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	log.Printf("Fetching raffle with ID: %d", id)
	raffle, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching raffle: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": raffle, "message": "Raffle fetched successfully"})
}

func (h *RaffleHandler) GetBySlug(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		userID = 0
	}

	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Slug is required"})
		return
	}

	log.Printf("Fetching raffle with slug: %s", slug)
	raffle, err := h.service.GetBySlug(c.Request.Context(), userID, slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching raffle: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": raffle, "message": "Raffle fetched successfully"})
}

func (h *RaffleHandler) GetByOrganizerID(c *gin.Context) {
	organizerID, err := utils.ValidateParamsId(c, "organizer_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Organizer ID: " + err.Error()})
		return
	}

	raffles, err := h.service.GetByOrganizerID(c.Request.Context(), organizerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching raffles by organizer: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": raffles, "message": "Raffles fetched successfully by organizer"})
}

func (h *RaffleHandler) Create(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
		return
	}

	var form models.Raffle
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("Error binding form data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	raffle, err := h.service.Create(c.Request.Context(), userID, form.Title, form.Description, form.Price, form.ImageURL, form.Date, form.QuantityTickets, form.MaxWinners, form.RaffleStatusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating raffle: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": raffle, "message": "Raffle created successfully"})
}

func (h *RaffleHandler) Update(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized: " + err.Error()})
		return
	}

	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	var form models.Raffle
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	raffle, err := h.service.Update(c.Request.Context(), userID, id, form.Title, form.Description, form.Price, form.ImageURL, form.Date, form.QuantityTickets, form.MaxWinners, form.RaffleStatusID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating raffle: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": raffle, "message": "Raffle updated successfully"})
}
