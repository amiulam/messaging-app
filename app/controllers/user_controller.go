package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	jwttoken "github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	user := new(models.User)

	err := ctx.BodyParser(user)
	if err != nil {
		fmt.Println("Failed to parse request ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	err = user.Validate()
	if err != nil {
		fmt.Println("Failed to validate user request ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Fail to encrypt password")
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	user.Password = string(hashedPassword)

	err = repository.InsertNewUser(ctx.Context(), user)
	if err != nil {
		fmt.Println("Failed to insert user data", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	resp := user
	resp.Password = ""

	return response.SendSuccessResponse(ctx, resp)
}

func Login(ctx *fiber.Ctx) error {
	loginReq := new(models.LoginRequest)
	resp := models.LoginResponse{}

	err := ctx.BodyParser(loginReq)
	if err != nil {
		fmt.Println("Failed to parse request ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	err = loginReq.Validate()
	if err != nil {
		fmt.Println("Failed to validate user request ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	user, err := repository.GetUserByUsername(ctx.Context(), loginReq.Username)
	if err != nil {
		fmt.Println("Failed to get user by username ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "invalid credentials", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		fmt.Println("Failed to check password ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "invalid credentials", nil)
	}

	token, err := jwttoken.GenerateToken(ctx.Context(), user.Username, user.Fullname, "token")
	if err != nil {
		fmt.Println("Failed to generate token ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	refreshToken, err := jwttoken.GenerateToken(ctx.Context(), user.Username, user.Fullname, "refresh_token")
	if err != nil {
		fmt.Println("Failed to generate refresh token ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	userSession := &models.UserSession{
		UserID:              user.ID,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        time.Now().Add(jwttoken.MapTokenType["token"]),
		RefreshTokenExpired: time.Now().Add(jwttoken.MapTokenType["refresh_token"]),
	}

	err = repository.InsertNewUserSession(ctx.Context(), userSession)
	if err != nil {
		fmt.Println("Failed to insert user session ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	resp.Username = user.Username
	resp.Fullname = user.Fullname
	resp.Token = token
	resp.RefreshToken = refreshToken

	return response.SendSuccessResponse(ctx, resp)
}

func Logout(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")
	err := repository.DeleteUserSessionByToken(ctx.Context(), token)
	if err != nil {
		fmt.Println("Failed to delete user session ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return response.SendSuccessResponse(ctx, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {
	refreshToken := ctx.Get("Authorization")
	username := ctx.Locals("username").(string)
	fullname := ctx.Locals("full_name").(string)

	token, err := jwttoken.GenerateToken(ctx.Context(), username, fullname, "token")
	if err != nil {
		fmt.Println("Failed to generate token from refresh token", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	err = repository.UpdateUserSessionToken(ctx.Context(), token, refreshToken, time.Now().Add(jwttoken.MapTokenType["token"]))
	if err != nil {
		fmt.Println("Failed to update user session from refresh token", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}
