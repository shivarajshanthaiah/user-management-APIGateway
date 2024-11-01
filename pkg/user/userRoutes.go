package user

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/config"
	"github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/user/handler"
	pb "github.com/shivaraj-shanthaiah/user-management-apigateway/pkg/user/userpb"
)

type User struct {
	cfg    *config.Config
	Client pb.UserServiceClient
}

func NewUserRoute(c *gin.Engine, cfg config.Config) {
	client, err := ClientDial(cfg)
	if err != nil {
		log.Fatalf("Error while connecting with grpc client :%v", err.Error())
	}

	userHandler := &User{
		cfg:    &cfg,
		Client: client,
	}
	apiVersion := c.Group("/api/v1")
	user := apiVersion.Group("/user")
	{
		user.POST("/create", userHandler.CreateUser)
		user.GET("/get/:id", userHandler.GetUserByID)
		user.PATCH("update/:id", userHandler.UpdateUserByID)
		user.DELETE("delete/:id", userHandler.DeleteUserBYID)
		user.POST("/methods", userHandler.HandleMethods)

	}
}

func (u *User) CreateUser(ctx *gin.Context) {
	handler.CreateUserHandler(ctx, u.Client)
}

func (u *User) GetUserByID(ctx *gin.Context) {
	handler.GetUserByIDHandler(ctx, u.Client)
}

func (u *User) UpdateUserByID(ctx *gin.Context) {
	handler.UpdateUserByIDHandler(ctx, u.Client)
}

func (u *User) DeleteUserBYID(ctx *gin.Context) {
	handler.DeleteUserByIDHandler(ctx, u.Client)
}

// func (u *User) HandleMethods(ctx *gin.Context) {
// 	var req dto.MethodRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
// 		return
// 	}

// 	switch req.Method {
// 	case 1:
// 		u.Method1(ctx, req.WaitTime)
// 	case 2:
// 		u.Method2(ctx, req.WaitTime)
// 	default:
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid method number"})
// 	}
// }

// func (u *User) Method1(ctx *gin.Context, waitTime int) {
// 	handler.Method1Handler(ctx, u.Client, waitTime)
// }

// func (u *User) Method2(ctx *gin.Context, waitTime int) {
// 	handler.Method2Handler(ctx, u.Client, waitTime)
// }
