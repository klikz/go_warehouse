package apiserver

import (
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) CheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := models.Request{}
		resp := models.Responce{}

		if err := c.ShouldBind(&req); err != nil {
			s.Logger.Error("Error Pasing body in CheckRole(): ", err)
			resp.Result = "error"
			resp.Err = err
			c.JSON(401, resp)
			c.Abort()
			return
		}

		parsedToken, err := ParseToken(req.Token)
		if err != nil {
			s.Logger.Error("Wrong Token: ", req.Token, " error: ", err)
			resp.Result = "error"
			resp.Err = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}

		res, err := s.Store.Repo().CheckRole(c.Request.URL.String(), parsedToken.Email)
		if err != nil {
			s.Logger.Error("CheckRole: ", req.Token, " error: ", err)
			resp.Result = "error"
			resp.Err = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}

		if !res {
			s.Logger.Error("CheckRole: ", req.Token, " error: ", err)
			resp.Result = "error"
			resp.Err = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}
		s.Logger.Info("Action URL: ", c.Request.URL.String(), " user: ", parsedToken.Email)
		c.Set("id", req.ID)
		c.Set("date1", req.Date1)
		c.Set("date2", req.Date2)
		c.Set("email", req.Email)
		c.Set("line", req.Line)
		c.Set("name", parsedToken.Email)
		c.Set("serial", req.Serial)
		c.Set("defect", req.Defect)
		c.Set("checkpoint", req.Checkpoint)
		c.Set("packing", req.Packing)
		c.Set("password", req.Password)

	}
}

func (s *Server) NoCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := models.Request{}
		resp := models.Responce{}

		if err := c.ShouldBind(&req); err != nil {
			s.Logger.Error("Error Pasing body in NoCheckRole(): ", err)
			resp.Result = "error"
			resp.Err = err
			c.JSON(401, resp)
			c.Abort()
			return
		}

		c.Set("id", req.ID)
		c.Set("date1", req.Date1)
		c.Set("date2", req.Date2)
		c.Set("email", req.Email)
		c.Set("line", req.Line)
		c.Set("name", req.Name)
		c.Set("serial", req.Serial)
		c.Set("defect", req.Defect)
		c.Set("checkpoint", req.Checkpoint)
		c.Set("packing", req.Packing)

	}
}

func (s *Server) WareCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := models.Component{}
		resp := models.Responce{}

		if err := c.ShouldBind(&req); err != nil {
			s.Logger.Error("Error Pasing body in CheckRole(): ", err)
			resp.Result = "error"
			resp.Err = err
			c.JSON(401, resp)
			c.Abort()
			return
		}

		parsedToken, err := ParseToken(req.Token)
		if err != nil {
			s.Logger.Error("WareCheckRole Wrong Token: ", req.Token, " error: ", err)
			resp.Result = "error"
			resp.Err = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}

		res, err := s.Store.Repo().CheckRole(c.Request.URL.String(), parsedToken.Email)
		if err != nil {
			s.Logger.Error("WareCheckRole: ", req.Token, " error: ", err)
			resp.Result = "error"
			resp.Err = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}

		if !res {
			s.Logger.Error("WareCheckRole: ", req.Token, " error: ", err)
			resp.Result = "error"
			resp.Err = "Wrong Credentials"
			c.JSON(401, resp)
			c.Abort()
			return
		}
		s.Logger.Info("Action URL: ", c.Request.URL.String(), " user: ", parsedToken.Email)

		c.Set("available", req.Available)
		c.Set("id", req.ID)
		c.Set("code", req.Code)
		c.Set("name", req.Name)
		c.Set("checkpoint", req.Checkpoint)
		c.Set("checkpoint_id", req.Checkpoint_id)
		c.Set("unit", req.Unit)
		c.Set("specs", req.Specs)
		c.Set("photo", req.Photo)
		c.Set("time", req.Time)
		c.Set("type", req.Type)
		c.Set("type_id", req.Type_id)
		c.Set("weight", req.Weight)
		c.Set("quantity", req.Quantity)
		c.Set("comment", req.Comment)
		c.Set("model_id", req.Model_ID)
		c.Set("component_id", req.Component_id)
		c.Set("date1", req.Date1)
		c.Set("date2", req.Date2)

	}
}
