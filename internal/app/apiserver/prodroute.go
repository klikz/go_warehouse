package apiserver

import (
	"warehouse/internal/app/models"

	"github.com/gin-gonic/gin"
)

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

func (s *Server) GetCountByDate(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("date1")
	temp2, _ := c.Get("date2")
	temp3, _ := c.Get("line")
	date1 := temp.(string)
	date2 := temp2.(string)
	line := temp3.(int)

	data, err := s.Store.Repo().GetCountByDate(date1, date2, line)
	if err != nil {
		s.Logger.Error("GetByDate: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetByDateModels(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("date1")
	temp2, _ := c.Get("date2")
	temp3, _ := c.Get("line")
	date1 := temp.(string)
	date2 := temp2.(string)
	line := temp3.(int)

	data, err := s.Store.Repo().GetByDateModels(date1, date2, line)
	if err != nil {
		s.Logger.Error("GetByDateModels: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetRemont(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetRemont()
	if err != nil {
		s.Logger.Error("GetRemont: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetRemontToday(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetRemontToday()
	if err != nil {
		s.Logger.Error("GetRemontToday: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetRemontByDate(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("date1")
	temp2, _ := c.Get("date2")
	date1 := temp.(string)
	date2 := temp2.(string)

	data, err := s.Store.Repo().GetRemontByDate(date1, date2)
	if err != nil {
		s.Logger.Error("GetByDate: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) UpdateRemont(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("name")
	temp2, _ := c.Get("id")
	name := temp.(string)
	id := temp2.(int)

	err := s.Store.Repo().UpdateRemont(name, id)
	if err != nil {
		s.Logger.Error("GetByDate: ", err)
		resp.Result = "error"
		resp.Err = "Wrong Credentials"
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) SerialInput(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("serial")
	temp2, _ := c.Get("line")
	serial := temp.(string)
	line := temp2.(int)

	err := s.Store.Repo().SerialInput(line, serial)
	if err != nil {
		s.Logger.Error("SerialInput: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) PackingSerialInput(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("serial")
	temp2, _ := c.Get("packing")
	serial := temp.(string)
	packing := temp2.(string)

	if packing == serial || packing == "" {
		s.Logger.Error("SerialInput: ", "serial: ", serial, "packing: ", packing)
		resp.Result = "error"
		resp.Err = "Serial xato"
		c.JSON(200, resp)
		return
	}

	err := s.Store.Repo().PackingSerialInput(serial, packing)
	if err != nil {
		s.Logger.Error("SerialInput: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) GetInfoBySerial(c *gin.Context) {
	resp := models.Responce{}
	temp, _ := c.Get("serial")
	serial := temp.(string)

	data, err := s.Store.Repo().GetInfoBySerial(serial)
	if err != nil {
		s.Logger.Error("SerialInput: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GalileoInput(c *gin.Context) {

	req := models.Galileo{}
	resp := models.Responce{}

	if err := c.ShouldBind(&req); err != nil {
		s.Logger.Error("Error Pasing body in NoCheckRole(): ", err)
		resp.Result = "error"
		resp.Err = err
		c.JSON(401, resp)
		c.Abort()
		return
	}

	err := s.Store.Repo().GalileoInput(&req)
	if err != nil {
		s.Logger.Error("SerialInput: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) GalileoTodayModels(c *gin.Context) {
	resp := models.Responce{}

	data, err := s.Store.Repo().GalileoTodayModels()
	if err != nil {
		s.Logger.Error("SerialInput: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}