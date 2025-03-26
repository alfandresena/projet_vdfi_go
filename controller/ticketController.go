package controllers

import (
	"net/http"
	"projet_vdfi/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Obtenir un ticket pour un événement
func CreateTicket(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Get("userID")

	var event models.Event
	if err := db.Where("id = ?", c.Param("event_id")).First(&event).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Événement non trouvé"})
		return
	}

	// Vérifier la date de l'événement
	if time.Now().After(event.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'événement est terminé"})
		return
	}

	// Vérifier si l'utilisateur a déjà un ticket
	var existingTicket models.Ticket
	if err := db.Where("user_id = ? AND event_id = ?", userID, event.ID).First(&existingTicket).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Vous avez déjà un ticket pour cet événement"})
		return
	}

	// Créer le ticket
	ticket := models.Ticket{
		UserID:  userID.(uint),
		EventID: event.ID,
	}

	if err := db.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer le ticket"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Ticket obtenue avec succès", "ticket": ticket})
}

// Récupérer les tickets de l'utilisateur connecté
func GetUserTickets(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Get("userID")

	var tickets []models.Ticket
	if err := db.Preload("Event").Where("user_id = ?", userID).Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// Supprimer un ticket
func DeleteTicket(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID, _ := c.Get("userID")

	var ticket models.Ticket
	if err := db.Where("id = ? AND user_id = ?", c.Param("id"), userID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket non trouvé"})
		return
	}

	db.Delete(&ticket)
	c.JSON(http.StatusOK, gin.H{"message": "Ticket supprimé avec succès"})
}
