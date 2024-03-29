package controllers

import (
	"fmt"
	"main/handlers"
	"main/models"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SubmitFlag(c *gin.Context) {
	var submitFlag models.SubmitFlag
	var user models.User
	var flag models.Flag
	var challengeFlags []models.Flag
	var userChallenge models.UserChallenge

	if err := c.ShouldBindJSON(&submitFlag); err != nil {
		fmt.Println("Request not binding...")
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	if err := handlers.DB.Preload("Flags").Where("username = ?", submitFlag.Username).Find(&user).Error; err != nil {
		fmt.Println("User not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	if err := handlers.DB.Where("flag = ?", submitFlag.Flag).Find(&flag).Error; err != nil || flag.Score == 0 {
		fmt.Println("Flag not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Flag not found!"})
		return
	}

	if err := handlers.DB.Where("challenge_id = ?", flag.ChallengeID).Find(&challengeFlags).Error; err != nil {
		fmt.Println("Flag not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Challenge flags not found!"})
		return
	}

	handlers.DB.Model(&user).Association("Flags").Append(&models.Flag{Flag: flag.Flag})

	handlers.DB.Save(user)

	userChallenge.UserUsername = user.Username
	userChallenge.ChallengeID = flag.ChallengeID

	handlers.DB.First(&userChallenge)

	userChallenge.State, userChallenge.Status = utils.CompareUserChallengeFlags(user, challengeFlags)

	handlers.DB.Save(&userChallenge)

	c.JSON(http.StatusOK, gin.H{"Flag submitted:": flag})
}

func CreateFlag(c *gin.Context) {
	var flagToCreate models.CreateFlag

	if err := c.ShouldBindJSON(&flagToCreate); err != nil {
		fmt.Println("Not binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": err.Error()})
		return
	}

	flag := models.Flag{Flag: flagToCreate.Flag, Score: flagToCreate.Score}
	handlers.DB.Create(&flag)

	c.JSON(http.StatusOK, gin.H{"Flag created:": flag})
}

func ShowAllFlags(c *gin.Context) {
	var flags []models.Flag
	handlers.DB.Find(&flags)

	c.JSON(http.StatusOK, gin.H{"data": flags})
}

func GetAllUserFlags(c *gin.Context) {
	var user models.User

	fmt.Println(c.Param("username"))
	if err := handlers.DB.Preload("Flags").Where("username = ?", c.Param("username")).First(&user).Error; err != nil {
		fmt.Println("Flags not found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Flags": user.Flags})
}
