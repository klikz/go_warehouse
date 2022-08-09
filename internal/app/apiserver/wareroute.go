package apiserver

import (
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetAllComponents(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetAllComponents()
	if err != nil {
		s.Logger.Error("GetAllComponents: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetCompoment(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt("id")

	data, err := s.Store.Repo().GetComponent(id)
	if err != nil {
		s.Logger.Error("GetComponent: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) UpdateCompoment(c *gin.Context) {
	component := models.Component{}
	resp := models.Responce{}

	component.Available = c.GetFloat64("available")
	component.ID = c.GetInt("id")
	component.Code = c.GetString("code")
	component.Name = c.GetString("name")
	component.Checkpoint = c.GetString("checkpoint")
	component.Checkpoint_id = c.GetInt("checkpoint_id")
	component.Unit = c.GetString("unit")
	component.Specs = c.GetString("specs")
	component.Photo = c.GetString("photo")
	component.Time = c.GetString("time")
	component.Type = c.GetString("type")
	component.Type_id = c.GetInt("type_id")
	component.Weight = c.GetFloat64("weight")

	if component.ID == 0 {
		s.Logger.Error("GetComponent: ", "blank id")
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}

	err := s.Store.Repo().UpdateComponent(&component)
	if err != nil {
		s.Logger.Error("GetComponent: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) AddComponent(c *gin.Context) {
	component := models.Component{}
	resp := models.Responce{}

	component.Code = c.GetString("code")
	component.Name = c.GetString("name")
	component.Checkpoint_id = c.GetInt("checkpoint_id")
	component.Unit = c.GetString("unit")
	component.Specs = c.GetString("specs")
	component.Photo = c.GetString("photo")
	component.Type_id = c.GetInt("type_id")
	component.Weight = c.GetFloat64("weight")
	s.Logger.Info(component)
	err := s.Store.Repo().AddComponent(&component)
	if err != nil {
		s.Logger.Error("AddComponent: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) DeleteCompoment(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt("id")

	err := s.Store.Repo().DeleteComponent(id)
	if err != nil {
		s.Logger.Error("GetComponent: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) GetAllCheckpoints(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetAllCheckpoints()
	if err != nil {
		s.Logger.Error("GetAllCheckpoints: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) DeleteCheckpoint(c *gin.Context) {

	resp := models.Responce{}

	id := c.GetInt("id")

	err := s.Store.Repo().DeleteCheckpoint(id)
	if err != nil {
		s.Logger.Error("DeleteCheckpoints: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) AddCheckpoint(c *gin.Context) {

	resp := models.Responce{}

	name := c.GetString("name")
	photo := c.GetString("photo")

	err := s.Store.Repo().AddCheckpoint(name, photo)
	if err != nil {
		s.Logger.Error("AddCheckpoint: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) UpdateCheckpoint(c *gin.Context) {

	resp := models.Responce{}

	name := c.GetString("name")
	photo := c.GetString("photo")
	id := c.GetInt(("id"))

	err := s.Store.Repo().UpdateCheckpoint(name, photo, id)
	if err != nil {
		s.Logger.Error("UpdateCheckpoint: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) Income(c *gin.Context) {

	resp := models.Responce{}
	quantity := c.GetFloat64("quantity")
	id := c.GetInt(("id"))

	err := s.Store.Repo().Income(id, quantity)
	if err != nil {
		s.Logger.Error("Income: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) Types(c *gin.Context) {

	resp := models.Responce{}

	data, err := s.Store.Repo().Types()
	if err != nil {
		s.Logger.Error("Types: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, data)
}
