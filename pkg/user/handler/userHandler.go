package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	dto "github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/DTO"
	pb "github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/user/userpb"
	"github.com/shivaraj-shanthaiah/user-management-apigateway/utility"
)

func CreateUserHandler(c *gin.Context, client pb.UserServiceClient) {
	timeout := time.Second * 100
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var user dto.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "error while binding json",
			"Error":   err.Error()})
		return
	}

	// Validate struct
	validate := validator.New()
	validate.RegisterValidation("email", utility.EmailValidation)
	validate.RegisterValidation("phone", utility.PhoneNumberValidation)
	validate.RegisterValidation("user_name", utility.AlphaSpace)

	// Validate struct
	err := validate.Struct(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Status": http.StatusBadRequest,
			"Errors": utility.ExtractValidationErrors(err),
		})
		return
	}

	response, err := client.CreateUser(ctx, &pb.Create{
		User_Name: user.UserName,
		Email:     user.Email,
		Phone:     user.Phone,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "error in client response",
			"Data":    response,
			"Error":   err.Error()})
		return
	}

	if response.Status == pb.Response_ERROR {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": response.Message,
			"Data":    response.Payload,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"Status":  http.StatusAccepted,
		"Message": "User created successfully",
		"Data":    response,
	})
}

func GetUserByIDHandler(c *gin.Context, client pb.UserServiceClient) {
	timeout := time.Second * 100
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Invalid user ID",
			"Error":   err.Error(),
		})
		return
	}

	response, err := client.GetUserByID(ctx, &pb.ID{
		ID: uint32(userID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Error fetching user details",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Status":  http.StatusOK,
		"Message": "Responded successfully",
		"Data":    response,
	})
}

func UpdateUserByIDHandler(c *gin.Context, client pb.UserServiceClient) {
	timeout := time.Second * 100
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Invalid user ID",
			"Error":   err.Error(),
		})
		return
	}

	var user dto.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "error while binding json",
			"Error":   err.Error()})
		return
	}

	response, err := client.UpdateUser(ctx, &pb.Profile{
		User_ID:   uint32(userID),
		User_Name: user.UserName,
		Phone:     user.Phone,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Error in client response",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"Status":  http.StatusAccepted,
		"Message": "User details updated successfully",
		"Data":    response,
	})
}

func DeleteUserByIDHandler(c *gin.Context, client pb.UserServiceClient) {
	timeout := time.Second * 100
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	userIDParam := c.Param("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  http.StatusBadRequest,
			"Message": "Invalid user ID",
			"Error":   err.Error(),
		})
		return
	}

	response, err := client.DeleteUserBYID(ctx, &pb.ID{
		ID: uint32(userID),
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"Status":  http.StatusInternalServerError,
			"Message": "Error in client response",
			"Error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"Status":  http.StatusAccepted,
		"Message": "User deleted successfully",
		"Data":    response,
	})
}

// func Method1Handler(ctx *gin.Context, client pb.UserServiceClient, waitTime int) {
// 	time.Sleep(time.Duration(waitTime) * time.Second)

// 	ctx.JSON(http.StatusOK, gin.H{"status": "Method 1 executed after waiting", "waitTime": waitTime})
// }

// func Method2Handler(c *gin.Context, client pb.UserServiceClient, waitTime int) {
// 	timeout := time.Second * 10
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()

// 	usersResp, err := client.GetAllUsers(ctx, &pb.NoParams{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
// 		return
// 	}

// 	var userNames []string
// 	for _, user := range usersResp.Users {
// 		userNames = append(userNames, user.User_Name)
// 	}

// 	time.Sleep(time.Duration(waitTime) * time.Second)

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":    "successfully fetched user using 2nd method",
// 		"users":     userNames,
// 		"Wait time": strconv.Itoa(waitTime) + " seconds",
// 	})
// }
