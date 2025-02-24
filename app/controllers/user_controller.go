package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/models"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	jwttoken "github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"go.elastic.co/apm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Register", "controller")
	defer span.End()

	user := new(models.User)

	err := ctx.BodyParser(user)
	if err != nil {
		log.Println("Failed to parse request ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	err = user.Validate()
	if err != nil {
		log.Println("Failed to validate user request ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Fail to encrypt password")
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	user.Password = string(hashedPassword)

	err = repository.InsertNewUser(spanCtx, user)
	if err != nil {
		log.Println("Failed to insert user data", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	resp := user
	resp.Password = ""

	return response.SendSuccessResponse(ctx, resp)
}

func Login(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Login", "controller")
	defer span.End()

	loginReq := new(models.LoginRequest)
	resp := models.LoginResponse{}

	err := ctx.BodyParser(loginReq)
	if err != nil {
		log.Println("Failed to parse request ", err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	err = loginReq.Validate()
	if err != nil {
		log.Println("Failed to validate user request ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}

	user, err := repository.GetUserByUsername(spanCtx, loginReq.Username)
	if err != nil {
		log.Println("Failed to get user by username ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "invalid credentials", nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		log.Println("Failed to check password ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusNotFound, "invalid credentials", nil)
	}

	token, err := jwttoken.GenerateToken(spanCtx, user.Username, user.Fullname, "token")
	if err != nil {
		log.Println("Failed to generate token ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	refreshToken, err := jwttoken.GenerateToken(spanCtx, user.Username, user.Fullname, "refresh_token")
	if err != nil {
		log.Println("Failed to generate refresh token ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	userSession := &models.UserSession{
		UserID:              user.ID,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        time.Now().Add(jwttoken.MapTokenType["token"]),
		RefreshTokenExpired: time.Now().Add(jwttoken.MapTokenType["refresh_token"]),
	}

	err = repository.InsertNewUserSession(spanCtx, userSession)
	if err != nil {
		log.Println("Failed to insert user session ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	resp.Username = user.Username
	resp.Fullname = user.Fullname
	resp.Token = token
	resp.RefreshToken = refreshToken

	return response.SendSuccessResponse(ctx, resp)
}

func Logout(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "Logout", "controller")
	defer span.End()

	token := ctx.Get("Authorization")
	err := repository.DeleteUserSessionByToken(spanCtx, token)
	if err != nil {
		log.Println("Failed to delete user session ", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return response.SendSuccessResponse(ctx, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.Context(), "RefreshToken", "controller")
	defer span.End()

	refreshToken := ctx.Get("Authorization")
	username := ctx.Locals("username").(string)
	fullname := ctx.Locals("full_name").(string)

	token, err := jwttoken.GenerateToken(spanCtx, username, fullname, "token")
	if err != nil {
		log.Println("Failed to generate token from refresh token", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	err = repository.UpdateUserSessionToken(spanCtx, token, refreshToken, time.Now().Add(jwttoken.MapTokenType["token"]))
	if err != nil {
		log.Println("Failed to update user session from refresh token", err.Error())
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	return response.SendSuccessResponse(ctx, fiber.Map{
		"token": token,
	})
}
