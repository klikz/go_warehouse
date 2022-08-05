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
	temp, _ := c.Get("id")

	id := temp.(int)
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

func (s *Server) GetLast(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("line")
	line := temp.(int)

	data, err := s.Store.Repo().GetLast(line)
	if err != nil {
		s.Logger.Error("GetLast: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetStatus(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("line")
	line := temp.(int)

	data, err := s.Store.Repo().GetStatus(line)
	if err != nil {
		s.Logger.Error("GetLast: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetToday(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("line")
	line := temp.(int)

	today, err := s.Store.Repo().GetToday(line)
	if err != nil {
		s.Logger.Error("GetLast: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, today)
}

func (s *Server) GetTodayModels(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("line")
	line := temp.(int)

	data, err := s.Store.Repo().GetTodayModels(line)
	if err != nil {
		s.Logger.Error("GetTodayModels: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetSectorBalance(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("line")
	line := temp.(int)

	data, err := s.Store.Repo().GetSectorBalance(line)
	if err != nil {
		s.Logger.Error("GetSectorBalance: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetPackingLast(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetPackingLast()
	if err != nil {
		s.Logger.Error("GetSectorBalance: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetPackingToday(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetPackingToday()
	if err != nil {
		s.Logger.Error("GetPackingToday: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetPackingTodaySerial(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetPackingTodaySerial()
	if err != nil {
		s.Logger.Error("GetPackingTodaySerial: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetPackingTodayModels(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetPackingTodayModels()
	if err != nil {
		s.Logger.Error("GetPackingTodayModels: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetLines(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetLines()
	if err != nil {
		s.Logger.Error("GetLines: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}
func (s *Server) GetDefectsTypes(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetDefectsTypes()
	if err != nil {
		s.Logger.Error("GetDefectsTypes: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) DeleteDefectsTypes(c *gin.Context) {
	resp := models.Responce{}

	temp, _ := c.Get("id")
	id := temp.(int)

	err := s.Store.Repo().DeleteDefectsTypes(id)
	if err != nil {
		s.Logger.Error("GetDefectsTypes: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) AddDefectsTypes(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("line")
	temp2, _ := c.Get("name")
	line := temp.(int)
	name := temp2.(string)

	err := s.Store.Repo().AddDefectsTypes(line, name)
	if err != nil {
		s.Logger.Error("AddDefectsTypes: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) AddDefects(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("serial")
	temp2, _ := c.Get("checkpoint")
	temp4, _ := c.Get("name")
	temp5, _ := c.Get("defect")

	checkpoint := temp2.(int)
	name := temp4.(string)
	serial := temp.(string)
	defect := temp5.(int)
	// u.Serial, name, u.Line, id.ID, u.Defect
	s.Logger.Info("serial:", serial, " name: ", name, " check: ", checkpoint, " def ", defect)
	err := s.Store.Repo().AddDefects(serial, name, checkpoint, defect)
	if err != nil {
		s.Logger.Error("AddDefects: ", err)
		resp.Result = "error"
		resp.Err = "Serial xato"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) GetByDateSerial(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("date1")
	temp2, _ := c.Get("date2")
	date1 := temp.(string)
	date2 := temp2.(string)

	data, err := s.Store.Repo().GetByDateSerial(date1, date2)
	if err != nil {
		s.Logger.Error("GetByDateSerial: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}
