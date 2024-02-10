package controllers

import (
	"context"
	"example/STRUCTURE/pkg/database/mongodb/models"
	database "example/STRUCTURE/pkg/database/mongodb/repository"
	helper "example/STRUCTURE/pkg/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var organizationCollection *mongo.Collection = database.OpenConnection(database.Client, "organizations")

func CreateOrganization() gin.HandlerFunc {
	// Parse the request body
	return func(c *gin.Context) {

		var org models.Organization
		if err := c.ShouldBindJSON(&org); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if a valid token is provided in the Authorization header
		if tokenString := c.GetHeader("Authorization"); tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			return
		} else {
			// Verify the token and extract the user ID
			userIdString, err := helper.VerifyToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}

			// Convert the user ID string to an ObjectID value
			userId, err := StringToObjectID(userIdString)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Set the user ID on the organization object
			org.ID = userId
		}

		// Insert the organization document into the database
		result, err := organizationCollection.InsertOne(context.Background(), org)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the ID of the newly created organization
		c.JSON(http.StatusCreated, gin.H{"organization_id": result.InsertedID})

	}

}

func StringToObjectID(id string) (primitive.ObjectID, error) {

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {

		return primitive.ObjectID{}, err

	}

	return objectID, nil

}

func GetOrganization() gin.HandlerFunc {
	return func(c *gin.Context) {
		orgId := c.Param("org_id")

		if orgId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "org_id parameter is required"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var org models.Organization

		err := organizationCollection.FindOne(ctx, bson.M{"org_id": orgId}).Decode(&org)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, org)
	}
}

func GetAllOrganizations() gin.HandlerFunc {
	return func(c *gin.Context) {
		cursor, err := organizationCollection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while finding Organizations"})
			return
		}
		defer cursor.Close(context.Background())

		var allOrganizations []bson.M
		if err = cursor.All(context.Background(), &allOrganizations); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while decoding all Organization"})
			return
		}
		c.JSON(http.StatusOK, allOrganizations)
	}
}

func UpdateOrganization(c *gin.Context) {
	// Get the organization ID from the URL path parameter
	organizationID := c.Param("org_id")

	// Parse the organization ID as a MongoDB ObjectID
	id, err := primitive.ObjectIDFromHex(organizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	// Check if a valid token is provided in the Authorization header
	if tokenString := c.GetHeader("Authorization"); tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		return
	} else {
		// Verify the token and extract the user ID
		userIdString, err := helper.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Convert the user ID string to an ObjectID value
		userId, err := StringToObjectID(userIdString)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(userId)
	}

	// Parse the request body as an Organization object
	var org models.Organization
	if err := c.ShouldBindJSON(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the ID field of the organization object
	org.ID = id

	// Update the organization document in the database
	result, err := organizationCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: org.Name}}},
			{Key: "$set", Value: bson.D{{Key: "description", Value: org.Description}}},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the organization was updated
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	// Return the updated organization object
	c.JSON(http.StatusOK, org)
}

func InviteUserToOrganization(c *gin.Context) {
	orgId := c.Param("org_id")

	if orgId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "org_id parameter is required"})
		return
	}

	var inviteUserRequest models.InviteUserRequest
	if err := c.ShouldBindJSON(&inviteUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var organization models.Organization

	err := organizationCollection.FindOne(context.Background(), bson.M{"_id": orgId}).Decode(&organization)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		return
	}

	// Check if the user is already a member of the organization
	for _, member := range organization.Members {
		if member.UserID == inviteUserRequest.UserID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is already a member of the organization"})
			return
		}
	}

	// Update the organization document in the database
	result, err := organizationCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": orgId},
		bson.D{
			{Key: "$push", Value: bson.D{{Key: "members", Value: bson.D{{Key: "user_id", Value: inviteUserRequest.UserID}, {Key: "role", Value: inviteUserRequest.Role}}}}},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the user was invited to the organization
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		return
	}

	// Return the updated organization object
	c.JSON(http.StatusOK, organization)
}
