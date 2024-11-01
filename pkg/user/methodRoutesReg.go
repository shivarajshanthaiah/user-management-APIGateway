package user

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	dto "github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/DTO"
	pb "github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/user/userpb"
)

func (u *User) HandleMethods(ctx *gin.Context) {
	var req dto.MethodRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}
	// Create a channel to handle timeout
	done := make(chan struct{})

	go func() {
		switch req.Method {
		case 1:
			u.Method1(ctx, req.WaitTime)
		case 2:
			u.Method2(ctx, req.WaitTime)
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid method number"})
		}
		done <- struct{}{}
	}()

	// Set a global timeout of 30 seconds
	select {
	case <-done:
		return
	case <-time.After(30 * time.Second):
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Request timeout",
		})
		return
	}
}

var method1Mutex sync.Mutex

func (u *User) Method1(ctx *gin.Context, waitTime int) {
	// Try to acquire the lock
	method1Mutex.Lock()
	defer method1Mutex.Unlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(waitTime+5)*time.Second)
	defer cancel()

	// Create channels for results and errors
	usersChan := make(chan []*pb.Profile)
	errChan := make(chan error)

	// Fetch users asynchronously first
	go func() {
		resp, err := u.Client.GetAllUsers(timeoutCtx, &pb.NoParams{})
		if err != nil {
			errChan <- err
			return
		}
		usersChan <- resp.Users
	}()

	// Wait for user fetch operation
	var users []*pb.Profile
	select {
	case err := <-errChan:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch users",
			"details": err.Error(),
		})
		return
	case users = <-usersChan:
		// Continue processing
	case <-timeoutCtx.Done():
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Operation timed out while fetching users",
		})
		return
	}

	// Process user names
	var userNames []string
	for _, user := range users {
		userNames = append(userNames, user.User_Name)
	}

	// After getting the data, handle the wait time
	waitChan := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(waitTime) * time.Second)
		waitChan <- struct{}{}
	}()

	// Wait for either timeout or completion
	select {
	case <-timeoutCtx.Done():
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"error":   "Operation timed out",
			"details": "Request exceeded maximum allowed time",
		})
		return
	case <-waitChan:
		ctx.JSON(http.StatusOK, gin.H{
			"status":     "Method 1 executed successfully",
			"users":      userNames,
			"waitTime":   waitTime,
			"message":    "Request processed sequentially",
			"totalUsers": len(userNames),
		})
	}
}

func (u *User) Method2(ctx *gin.Context, waitTime int) {
	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(waitTime+5)*time.Second)
	defer cancel()

	// Create channels for results and errors
	usersChan := make(chan []*pb.Profile)
	errChan := make(chan error)

	// Fetch users asynchronously
	go func() {
		resp, err := u.Client.GetAllUsers(timeoutCtx, &pb.NoParams{})
		if err != nil {
			errChan <- err
			return
		}
		usersChan <- resp.Users
	}()

	// Wait for user fetch operation
	var users []*pb.Profile
	select {
	case err := <-errChan:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch users",
			"details": err.Error(),
		})
		return
	case users = <-usersChan:
		// Continue processing
	case <-timeoutCtx.Done():
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Operation timed out while fetching users",
		})
		return
	}

	// Process user names
	var userNames []string
	for _, user := range users {
		userNames = append(userNames, user.User_Name)
	}

	// Handle wait time asynchronously
	waitChan := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(waitTime) * time.Second)
		waitChan <- struct{}{}
	}()

	// Wait for either timeout or completion
	select {
	case <-timeoutCtx.Done():
		ctx.JSON(http.StatusGatewayTimeout, gin.H{
			"error": "Operation timed out during wait period",
		})
		return
	case <-waitChan:
		ctx.JSON(http.StatusOK, gin.H{
			"status":     "Method 2 executed successfully",
			"users":      userNames,
			"waitTime":   waitTime,
			"totalUsers": len(userNames),
		})
	}
}
