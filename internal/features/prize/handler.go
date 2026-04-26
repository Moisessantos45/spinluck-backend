package prize

import (
	"log"
	"net/http"
	"spinLuck/internal/shared/models"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

type PrizeHandler struct {
	service PrizeService
}

func NewPrizeHandler(service PrizeService) *PrizeHandler {
	return &PrizeHandler{service: service}
}

func (h *PrizeHandler) GetAllByIdRaffleByOrganizerID(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from token: " + err.Error()})
		return
	}

	raffleID, err := utils.ValidateParamsId(c, "raffleID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Raffle ID: " + err.Error()})
		return
	}

	prizes, err := h.service.GetAllByIdRaffleByOrganizerID(c.Request.Context(), userID, raffleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching prizes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": prizes, "message": "Prizes fetched successfully"})
}

func (h *PrizeHandler) GetAllByIdRafflePublic(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Slug is required"})
		return
	}

	prizes, err := h.service.GetAllByIdRafflePublic(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching prizes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": prizes, "message": "Prizes fetched successfully"})
}

func (h *PrizeHandler) GetByID(c *gin.Context) {
	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	prize, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching prize: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": prize, "message": "Prize fetched successfully"})
}

func (h *PrizeHandler) Create(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from token: " + err.Error()})
		return
	}

	var prizer models.Prize
	if err := c.ShouldBind(&prizer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input: " + err.Error()})
		return
	}

	log.Printf("Creating prize with title: %s, description: %s, imageURL: %s, raffleID: %d", prizer.Title, prizer.Description, prizer.ImageURL, prizer.RaffleID)
	prize, err := h.service.Create(c.Request.Context(), userID, prizer.Title, prizer.Description, prizer.ImageURL, prizer.RaffleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating prize: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": prize, "message": "Prize created successfully"})
}

func (h *PrizeHandler) Update(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from token: " + err.Error()})
		return
	}

	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	var prizer models.Prize
	if err := c.ShouldBind(&prizer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input: " + err.Error()})
		return
	}

	prize, err := h.service.Update(c.Request.Context(), userID, id, prizer.Title, prizer.Description, prizer.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating prize: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": prize, "message": "Prize updated successfully"})
}

func (h *PrizeHandler) Delete(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from token: " + err.Error()})
		return
	}

	raffleID, err := utils.ValidateParamsId(c, "raffleID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	prizeID, err := utils.ValidateParamsId(c, "prizeID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	err = h.service.Delete(c.Request.Context(), userID, raffleID, prizeID)
	if err != nil {
		log.Printf("Error deleting prize with ID %d from raffle %d: %v", prizeID, raffleID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting prize: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prize deleted successfully"})
}
