package apiserver

import (
	"fmt"
	"net/http"
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *Server) Create(c *gin.Context) {

	user := models.User{}
	resp := models.Responce{}

	if err := c.ShouldBind(&user); err != nil {
		s.Logger.Error("Create: Error Parsing body: ", err)
		resp.Result = "error"
		resp.Err = err
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	parsedToken, err := ParseToken(user.Token)
	if err != nil {
		s.Logger.Error("Create: Wrong Token: ", user.Email, " token: ", user.Token, " error: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if parsedToken.Role != "admin" {
		s.Logger.Error("Create: Wrong Role: ", user.Email, " error: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if !isEmailValid(user.Email) {
		s.Logger.Error("Create: Wrong email to login: ", user.Email, " error: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if !isPasswordValid(user.Password) {
		s.Logger.Error("Create: Try wrong password: ", user.Email, " error: ", err)
		resp.Result = "error"
		resp.Err = "password length < 6"
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	enc, err := encryptString(user.Password)
	if err != nil {
		resp.Result = "error"
		resp.Err = err
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	user.EncryptedPassword = enc
	err = s.Store.Repo().Create(&user)
	if err != nil {
		s.Logger.Error("Create: Error in create user: ", err)
		resp.Result = "error"
		resp.Err = err
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Result = fmt.Sprint(user.ID)
	c.JSON(200, resp)

}

func (s *Server) Login(c *gin.Context) {
	user := models.User{}
	resp := models.Responce{}

	if err := c.ShouldBind(&user); err != nil {
		logrus.Error("Login: Error Parsing body: ", err)
	}
	if err := s.Store.Repo().FindByEmail(&user); err != nil {
		resp.Result = "error"
		resp.Err = "wrong email or password"
		s.Logger.Error("Login: Try incorrect email to login: ", user.Email, " error: ", err)
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if !ComparePassword(user.Password, user.EncryptedPassword) {
		resp.Result = "error"
		resp.Err = "wrong email or password"
		s.Logger.Error("Login: Try incorrect email to password: ", user.Email)
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := GetToken(&user); err != nil {
		s.Logger.Error("Login: GetToken: ", user.Email, " error: ", err)
		resp.Result = "error"
		resp.Err = err
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	s.Logger.Info("User Logged: ", user.Email, " client ip: ", c.ClientIP(), " remote ip: ", c.RemoteIP())
	c.JSON(200, user)

}
