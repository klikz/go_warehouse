package apiserver

import (
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetAllComponents(c *gin.Context) {
	resp := models.Responce{}
	comp, err := s.Store.Repo().GetAllComponents()
	if err != nil {
		s.Logger.Error("GetAllComponents: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(401, resp)
		return
	}
	c.JSON(200, comp)
}

func (s *Server) GetCompoment(c *gin.Context) {

	resp := models.Responce{}
	temp, _ := c.Get("id")

	id := temp.(int)
	component, err := s.Store.Repo().GetComponent(id)
	if err != nil {
		s.Logger.Error("GetComponent: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(401, resp)
		return
	}
	c.JSON(200, component)
}

func (s *Server) GetLast(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("line")
	line := temp.(int)

	last, err := s.Store.Repo().GetLast(line)
	if err != nil {
		s.Logger.Error("GetLast: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(401, resp)
		return
	}
	c.JSON(200, last)
}

func (s *Server) GetStatus(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("line")
	line := temp.(int)

	status, err := s.Store.Repo().GetStatus(line)
	if err != nil {
		s.Logger.Error("GetLast: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(401, resp)
		return
	}
	c.JSON(200, status)
}
