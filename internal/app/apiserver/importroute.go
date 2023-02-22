package apiserver

import (
	"encoding/base64"
	"io/ioutil"
	"warehouse/internal/app/models"

	"github.com/bingoohuang/xlsx"
	"github.com/gin-gonic/gin"
)

func (s *Server) InsertLot(c *gin.Context) {
	resp := models.Responce{}
	name := c.GetString("name")
	comment := c.GetString("comment")

	if name == "" {
		resp.Result = "error"
		resp.Err = "lot nomi kiritilmadi"
		c.JSON(200, resp)
		return
	}

	err := s.Store.Repo().InsertLot(name, comment)
	if err != nil {
		s.Logger.Error("InsertLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) DeleteLot(c *gin.Context) {
	resp := models.Responce{}
	lot_id := c.GetInt("lot_id")

	err := s.Store.Repo().DeleteLot(lot_id)
	if err != nil {
		s.Logger.Error("DeleteLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) UpdateLot(c *gin.Context) {
	resp := models.Responce{}
	name := c.GetString("name")
	comment := c.GetString("comment")
	lot_id := c.GetInt("lot_id")

	err := s.Store.Repo().UpdateLot(name, comment, lot_id)
	if err != nil {
		s.Logger.Error("UpdateLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) BlockLot(c *gin.Context) {
	resp := models.Responce{}
	lot_id := c.GetInt("lot_id")

	err := s.Store.Repo().BlockLot(lot_id)
	if err != nil {
		s.Logger.Error("BlockLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) UnBlockLot(c *gin.Context) {
	resp := models.Responce{}
	lot_id := c.GetInt("lot_id")

	err := s.Store.Repo().UnBlockLot(lot_id)
	if err != nil {
		s.Logger.Error("UnBlockLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) ActivateLot(c *gin.Context) {
	resp := models.Responce{}
	lot_id := c.GetInt("lot_id")

	err := s.Store.Repo().ActivateLot(lot_id)
	if err != nil {
		s.Logger.Error("ActivateLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) DeActivateLot(c *gin.Context) {
	resp := models.Responce{}
	lot_id := c.GetInt("lot_id")

	err := s.Store.Repo().DeActivateLot(lot_id)
	if err != nil {
		s.Logger.Error("DeActivateLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) GetAllLot(c *gin.Context) {
	resp := models.Responce{}

	data, err := s.Store.Repo().GetAllLot()
	if err != nil {
		s.Logger.Error("GetAllLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(200, resp)
}
func (s *Server) GetAllLotActive(c *gin.Context) {
	resp := models.Responce{}

	data, err := s.Store.Repo().GetAllLotActive()
	if err != nil {
		s.Logger.Error("GetAllLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(200, resp)
}

//routes for Batches

func (s *Server) InsertBatch(c *gin.Context) {
	resp := models.Responce{}
	name := c.GetString("name")
	comment := c.GetString("comment")
	lot_id := c.GetInt("lot_id")

	err := s.Store.Repo().InsertBatch(lot_id, name, comment)
	if err != nil {
		s.Logger.Error("InsertBatch: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) DeleteBatch(c *gin.Context) {
	resp := models.Responce{}
	batch_id := c.GetInt("batch_id")

	err := s.Store.Repo().DeleteBatch(batch_id)
	if err != nil {
		s.Logger.Error("DeleteBatch: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) UpdateBatch(c *gin.Context) {
	resp := models.Responce{}
	name := c.GetString("name")
	comment := c.GetString("comment")
	batch_id := c.GetInt("batch_id")
	print(name, " ", comment, " ", batch_id)
	err := s.Store.Repo().UpdateBatch(name, comment, batch_id)
	if err != nil {
		s.Logger.Error("UpdateLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) GetBatchByLot(c *gin.Context) {
	resp := models.Responce{}
	lot_id := c.GetInt("lot_id")

	data, err := s.Store.Repo().GetBatchByLot(lot_id)
	if err != nil {
		s.Logger.Error("GetBatchByLot: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(200, resp)
}

//routes for container
func (s *Server) InsertContainer(c *gin.Context) {
	resp := models.Responce{}
	name := c.GetString("name")
	comment := c.GetString("comment")
	lot_id := c.GetInt("lot_id")
	batch_id := c.GetInt("batch_id")

	err := s.Store.Repo().InsertContainer(name, comment, lot_id, batch_id)
	if err != nil {
		s.Logger.Error("InsertBatch: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) DeleteContainer(c *gin.Context) {
	resp := models.Responce{}
	container_id := c.GetInt("container_id")

	err := s.Store.Repo().DeleteContainer(container_id)
	if err != nil {
		s.Logger.Error("DeleteContainer: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) UpdateContainer(c *gin.Context) {
	resp := models.Responce{}
	name := c.GetString("name")
	comment := c.GetString("comment")
	container_id := c.GetInt("container_id")

	err := s.Store.Repo().UpdateContainer(name, comment, container_id)
	if err != nil {
		s.Logger.Error("UpdateContainer: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) GetContainerByBatch(c *gin.Context) {
	resp := models.Responce{}
	batch_id := c.GetInt("batch_id")

	data, err := s.Store.Repo().GetContainerByBatch(batch_id)
	if err != nil {
		s.Logger.Error("GetContainerByBatch: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(200, resp)
}

func (s *Server) GetContainerComponents(c *gin.Context) {
	resp := models.Responce{}
	container_id := c.GetInt("container_id")

	data, err := s.Store.Repo().GetContainerComponents(container_id)
	if err != nil {
		s.Logger.Error("GetContainerByBatch: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	resp.Data = data
	c.JSON(200, resp)
}

func (s *Server) ContainerComponentsDelete(c *gin.Context) {
	resp := models.Responce{}
	component_id := c.GetInt("component_id")

	err := s.Store.Repo().ContainerComponentsDelete(component_id)
	if err != nil {
		s.Logger.Error("ContainerComponentsDelete: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}
func (s *Server) ContainerComponentsUpdate(c *gin.Context) {
	resp := models.Responce{}
	r_quantity := c.GetFloat64("r_quantity")
	component_id := c.GetInt("component_id")
	comment := c.GetString("comment")

	err := s.Store.Repo().ContainerComponentsUpdate(r_quantity, comment, component_id)
	if err != nil {
		s.Logger.Error("ContainerComponentsDelete: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) ImportIncomeRegister(c *gin.Context) {

	resp := models.Responce{}
	lot_id := c.GetInt("lot_id")
	batch_id := c.GetInt("batch_id")
	container_id := c.GetInt("container_id")
	component_id := c.GetInt("component_id")
	quantity := c.GetFloat64("quantity")
	comment := c.GetString("comment")
	// unit := c.GetString("unit")

	err := s.Store.Repo().ImportIncomeRegister(lot_id, batch_id, container_id, component_id, quantity, comment, "")
	if err != nil {
		s.Logger.Error("ImportIncomeRegister: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}
func (s *Server) ImportIncomeAdd(c *gin.Context) {

	resp := models.Responce{}
	quantity := c.GetFloat64("quantity")
	income_id := c.GetInt("income_id")

	err := s.Store.Repo().ImportIncomeAdd(quantity, income_id)
	if err != nil {
		s.Logger.Error("ImportIncomeAdd: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) FileFromContainer(c *gin.Context) {
	file64 := c.GetString("file64")
	lot_id := c.GetInt("lot_id")
	batch_id := c.GetInt("batch_id")
	container_id := c.GetInt("container_id")

	resp := models.Responce{}

	data, err := base64Decode(file64)
	if err != nil {
		s.Logger.Error(err)
		return
	}
	err = ioutil.WriteFile("import.xlsx", []byte(data), 0644)
	if err != nil {
		s.Logger.Error(err)
	}

	var file []models.FileInput2
	x, _ := xlsx.New(xlsx.WithInputFile("import.xlsx"))
	defer x.Close()

	if err := x.Read(&file); err != nil {
		s.Logger.Error(err)
	}

	s.Logger.Info(file)

	err = s.Store.Repo().ImportIncomeRegisterFromFile(lot_id, batch_id, container_id, file)
	if err != nil {
		s.Logger.Error("ImportIncomeRegister: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		c.Abort()
		return
	}

	resp.Result = "ok"
	c.JSON(200, resp)

}
