package net

import (
	crypto "DevSpace/Crypto"
	db "DevSpace/DB"
	system "DevSpace/System"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	Config *system.Config
	DB     *db.DBManager
}

func (r *Rest) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	bindErr := c.ShouldBindJSON(&req)

	if bindErr != nil {
		c.JSON(400, gin.H{
			"error": "json error",
		})

		return
	}

	hash, err := crypto.EncodePassword(req.Password, r.Config)
	if err != nil {
		c.JSON(500, gin.H{"error": "error occured while hashing"})
		return
	}

	r.DB.InsertNewEntity(&db.User{Email: req.Email, Nickname: req.Nickname, PasswordHash: hash, AvatarUrl: "ЗАГЛУШКА", Status: "ЗАГЛУШКА", Bio: "ЗАГЛУШКА"})

	c.JSON(200, nil)

}

func (r *Rest) AuthUser(c *gin.Context) {
	var req AuthUserRequest

	bindErr := c.ShouldBindJSON(&req)
	if bindErr != nil {
		c.JSON(400, gin.H{
			"error": "json error",
		})

		return
	}

	if req.Email == nil && req.Login == nil {
		c.JSON(400, gin.H{
			"error": "no email and login",
		})

		return
	}

	var userHash string

	if req.Login != nil {
		r.DB.Connection.Table("User").Select("password_hash").Where("login = ?", req.Login).First(&userHash)
	} else {
		r.DB.Connection.Table("User").Select("password_hash").Where("email = ?", req.Email).First(&userHash)
	}

	if len(userHash) == 0 {
		c.JSON(404, nil)
		return
	}

	correct, err := crypto.VerifyPassword(req.Password, userHash)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "error while fetching password",
		})

		return
	} else if correct {
		c.JSON(200, nil)
		return
	} else {
		c.JSON(401, nil)
	}
}
