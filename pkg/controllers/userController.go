package controllers

import (
	"context"
	"encoding/json"
	"example/STRUCTURE/pkg/database/mongodb/models"
	database "example/STRUCTURE/pkg/database/mongodb/repository"
	helper "example/STRUCTURE/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenConnection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	result, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(result)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is not correct")
		check = false
	}

	return check, msg
}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Fatalln(err)
			return
		}

		validationErr := validate.Struct(&user)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			log.Fatalln("Validation", validationErr)

			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while checking for the email"})
			log.Fatalln("Error occurred while checking for the email")
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email is already exist"})
			return
		}

		var user_type = "USER"
		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		user.User_type = &user_type

		token, refreshtoken, _ := helper.GenerateAllTokens(*user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshtoken

		_, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.User_type, *&foundUser.User_id)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully", "access_token": &token, "refresh_token": &refreshToken})

	}
}
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		projection := bson.M{"password": 0, "refresh_token": 0, "token": 0}

		cursor, err := userCollection.Find(context.Background(), bson.M{}, options.Find().SetProjection(projection))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while finding users"})
			return
		}
		defer cursor.Close(context.Background())

		var allUsers []bson.M
		if err = cursor.All(context.Background(), &allUsers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding all users"})
			return
		}

		c.JSON(http.StatusOK, allUsers)
	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func GetUserFromDatabase(user_id string) (models.User, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User

	err := userCollection.FindOne(ctx, bson.M{"user_id": user_id}).Decode(&user)
	defer cancel()
	if err != nil {
		return models.User{}, fmt.Errorf("user not found: %v", err)
	}

	return user, nil
}
func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody struct {
			RefreshToken string `json:"refresh_token"`
		}

		err := json.NewDecoder(ctx.Request.Body).Decode(&reqBody)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := jwt.Parse(reqBody.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return helper.SECRET_KEY, nil
		})

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userId := strconv.FormatFloat(claims["id"].(float64), 'f', -1, 64)

		foundUser, err := GetUserFromDatabase(userId)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		index := -1

		for i, token := range *foundUser.Refresh_token {
			if string(token) == reqBody.RefreshToken {
				index = i
				break
			}
		}

		if index == -1 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  userId,
			"exp": time.Now().Add(time.Minute).Unix(),
		})

		accessTokenString, err := accessToken.SignedString(helper.SECRET_KEY)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.Header("Content-Type", "application/json")
		ctx.JSON(http.StatusOK, gin.H{"access_token": accessTokenString})
	}
}
