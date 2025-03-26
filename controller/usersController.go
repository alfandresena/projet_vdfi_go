package controllers

import (
	"net/http"
	"projet_vdfi/initializers"
	"projet_vdfi/models"

	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Clé secrète pour signer le token (même clé utilisée pour l'émission)
var SECRET_KEY = []byte(os.Getenv(("JWT_SECRET")))

// Signup permet d'inscrire un nouvel utilisateur
func Signup(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
	}

	// Vérifier si les données envoyées sont valides
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Hachage du mot de passe
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Création d'un nouvel utilisateur
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword),
	}

	// Enregistrement dans la base de données
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	// Réponse JSON avec l'utilisateur créé (sans mot de passe)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Utilisateur inscrit avec succès",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}

	// Vérifier si les données envoyées sont valides
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Rechercher l'utilisateur par email
	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Vérifier le mot de passe
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // 24 heures
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to generate token"})
		return
	}

	//send it as a cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	// Réponse JSON avec l'utilisateur connecté (sans mot de passe)
	c.JSON(http.StatusOK, gin.H{
		"message": "Utilisateur connecté avec succès",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"token": tokenString,
		},
	})

}

func Logout(c *gin.Context) {
	// Récupérer le token de l’utilisateur
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token found"})
		return
	}

	// Tenter de parser le token (vérifier sa validité)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRET_KEY, nil
	})

	// Si le token est valide, on génère un faux token avec une expiration immédiate
	if token != nil {
		expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"exp": time.Now().Unix(), // Expire immédiatement
		})

		// Signer le token avec la clé secrète
		expiredTokenString, _ := expiredToken.SignedString(SECRET_KEY)

		// Mettre le token invalide dans le cookie (le navigateur ne l’utilisera plus)
		c.SetCookie("Authorization", expiredTokenString, -1, "", "", false, true)
	}

	// Supprimer également le cookie du client pour éviter son réutilisation
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur déconnecté avec succès"})

}

func UpdateUser(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le token après passage par le middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convertir userID en uint
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var body struct {
		Name     string
		Email    string
		Password string
	}

	// Vérifier si les données envoyées sont valides
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Récupérer l'utilisateur à partir de la base de données
	var user models.User
	result := initializers.DB.First(&user, uid)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Mise à jour des champs modifiés
	if body.Name != "" {
		user.Name = body.Name
	}
	if body.Email != "" {
		user.Email = body.Email
	}
	if body.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Sauvegarder les modifications
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func GetUser(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le token après passage par le middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convertir userID en uint
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Récupérer l'utilisateur à partir de la base de données
	var user models.User
	result := initializers.DB.First(&user, uid)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"admin": user.IsAdmin,
		},
	})
}

func DeleteUser(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le token après passage par le middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Convertir userID en uint
	uid, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Récupérer l'utilisateur à partir de la base de données
	var user models.User
	result := initializers.DB.First(&user, uid)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Récupérer et vérifier le mot de passe fourni par l'utilisateur
	var body struct {
		Password string `json:"password"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Vérifier si le mot de passe est correct
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Supprimer l'utilisateur de la base de données
	deleteResult := initializers.DB.Delete(&user)
	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Supprimer le cookie
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	// Répondre avec un message de succès
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func PromoteToAdmin(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le token (déjà ajouté par le middleware)
	userID, _ := c.Get("userID") // Cela retourne l'ID de l'utilisateur authentifié à partir du contexte

	// Définir une structure pour récupérer le corps de la requête
	var body struct {
		Email      string
		KeySpecial string
	}

	// Vérifier si les données envoyées sont valides
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	// Vérifier si la clé spéciale est correcte
	if body.KeySpecial != os.Getenv("KEY_SPECIAL") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid special key"})
		return
	}

	// Rechercher l'utilisateur à promouvoir en utilisant l'ID depuis le token
	var user models.User
	result := initializers.DB.Where("id = ?", userID).First(&user)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Vérifier si l'utilisateur est déjà un administrateur (optionnel)
	if user.IsAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already an admin"})
		return
	}

	// Modifier le rôle de l'utilisateur en "admin"
	user.IsAdmin = true
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user to admin"})
		return
	}

	// Réponse de succès
	c.JSON(http.StatusOK, gin.H{
		"message": "User promoted to admin successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.IsAdmin,
		},
	})
}
