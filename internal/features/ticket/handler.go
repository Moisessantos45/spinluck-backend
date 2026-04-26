package ticket

import (
	"log"
	"net/http"
	"spinLuck/internal/shared/models"
	"spinLuck/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	service TicketService
}

func NewTicketHandler(service TicketService) *TicketHandler {
	return &TicketHandler{service: service}
}

func (h *TicketHandler) GetAllStatus(c *gin.Context) {
	statuses, err := h.service.GetAllStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching ticket statuses: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": statuses, "message": "Ticket statuses fetched successfully"})
}

func (h *TicketHandler) GetAll(c *gin.Context) {
	tickets, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching tickets: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tickets, "message": "Tickets fetched successfully"})
}

func (h *TicketHandler) GetByID(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from JWT: " + err.Error()})
		return
	}

	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	ticket, err := h.service.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching ticket: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket, "message": "Ticket fetched successfully"})
}

func (h *TicketHandler) GetAllByRaffleIDToOrgaizer(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from JWT: " + err.Error()})
		return
	}

	raffleID, err := utils.ValidateParamsId(c, "raffleID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid raffle ID: " + err.Error()})
		return
	}

	tickets, err := h.service.GetAllByRaffleIDToOrgaizer(c.Request.Context(), raffleID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching tickets: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets, "message": "Tickets fetched successfully"})
}

func (h *TicketHandler) GenerateVoucher(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from JWT: " + err.Error()})
		return
	}

	ticketID, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket ID: " + err.Error()})
		return
	}

	voucherBytes, err := h.service.GenerateVoucher(c.Request.Context(), ticketID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating voucher: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", voucherBytes)
}

func (h *TicketHandler) GetRandomSoldTicket(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from JWT: " + err.Error()})
		return
	}

	raffleID, err := utils.ValidateParamsId(c, "raffleID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid raffle ID: " + err.Error()})
		return
	}

	ticket, err := h.service.GetRandomSoldTicket(c.Request.Context(), userID, raffleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching ticket: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket, "message": "Random sold ticket fetched successfully"})
}

func (h *TicketHandler) GetRandomAvailableTicket(c *gin.Context) {
	raffleID, err := utils.ValidateParamsId(c, "raffleID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid raffle ID: " + err.Error()})
		return
	}

	log.Printf("Fetching random available ticket for raffle ID: %d", raffleID)

	ticket, err := h.service.GetRandomAvailableTicket(c.Request.Context(), raffleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching ticket: " + err.Error()})
		return
	}

	log.Printf("Random available ticket fetched: %+v", ticket)

	c.JSON(http.StatusOK, gin.H{"data": ticket, "message": "Random available ticket fetched successfully"})
}

func (h *TicketHandler) Create(c *gin.Context) {
	var input models.Ticket
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input: " + err.Error()})
		return
	}

	ticket, err := h.service.Create(c.Request.Context(), input.Number, input.ParticipantName, input.ParticipantPhone, input.RaffleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating ticket: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": ticket, "message": "Ticket created successfully"})
}

func (h *TicketHandler) Update(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from JWT: " + err.Error()})
		return
	}

	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	var input models.Ticket
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input: " + err.Error()})
		return
	}

	log.Printf("Received update request for ticket ID %d with data: %+v", id, input)

	ticket, err := h.service.Update(c.Request.Context(), id, input.ParticipantName, input.ParticipantPhone, input.TicketStatusID, userID)
	if err != nil {
		log.Printf("Error updating ticket: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating ticket: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket, "message": "Ticket updated successfully"})
}

func (h *TicketHandler) UpdateWinner(c *gin.Context) {
	_, userID, err := utils.ExtractedParamsJwt(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error extracting user ID from JWT: " + err.Error()})
		return
	}

	raffleID, err := utils.ValidateParamsId(c, "raffleID")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid raffle ID: " + err.Error()})
		return
	}

	id, err := utils.ValidateParamsId(c, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ID: " + err.Error()})
		return
	}

	if err := h.service.UpdateWinner(c.Request.Context(), id, raffleID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating ticket winner: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket winner updated successfully"})
}
