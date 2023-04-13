package apiserver

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"time"
	"warehouse/internal/app/models"

	"github.com/bingoohuang/xlsx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) TodayStatistics(c *gin.Context) {
	type AllInfo struct {
		Counters         interface{} `json:"counters"`
		DefectCounters   interface{} `json:"defect_counters"`
		MetallModels     interface{} `json:"metall_models"`
		SborkaModels     interface{} `json:"sborka_models"`
		PpuModels        interface{} `json:"ppu_models"`
		AgregatModels    interface{} `json:"agregat_models"`
		FreonModels      interface{} `json:"freon_models"`
		LaboratoryModels interface{} `json:"laboratory_models"`
		PackingModels    interface{} `json:"packing_models"`
		ModelsInMonth    interface{} `json:"models_month"`
		PlanMonth        int         `json:"plan_month"`
		CountMonth       int         `json:"count_month"`
		PercentMonth     float64     `json:"percent_month"`
		PlanDaily        int         `json:"plan_daily"`
		CountDaily       int         `json:"count_daily"`
		PercentDaily     float64     `json:"percent_daily"`
		CurrentMonth211  int         `json:"model_211"`
		CurrentMonth261  int         `json:"model_261"`
	}

	resp := models.Responce{}
	allInfo := AllInfo{}

	counters, err := s.Store.Repo().GetCounters()
	if err != nil {
		s.Logger.Error("GetCounters: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.DefectCounters, err = s.Store.Repo().GetDefectCounters()
	if err != nil {
		s.Logger.Error("allInfo.DefectCounters: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.MetallModels, err = s.Store.Repo().GetTodayModels(9)
	if err != nil {
		s.Logger.Error("allInfo.MetallModels: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.SborkaModels, err = s.Store.Repo().GetTodayModels(2)
	if err != nil {
		s.Logger.Error("allInfo.SborkaModels: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.PpuModels, err = s.Store.Repo().GetTodayModels(10)
	if err != nil {
		s.Logger.Error("allInfo.PpuModels: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.AgregatModels, err = s.Store.Repo().GetTodayModels(19)
	if err != nil {
		s.Logger.Error("allInfo.AgregatModels: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.FreonModels, err = s.Store.Repo().GalileoTodayModels()
	if err != nil {
		s.Logger.Error("allInfo.FreonModels: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.LaboratoryModels, err = s.Store.Repo().GetTodayModels(11)
	if err != nil {
		s.Logger.Error("allInfo.LaboratoryModels: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.PackingModels, err = s.Store.Repo().GetPackingTodayModels()
	if err != nil {
		s.Logger.Error("allInfo.PackingModels: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.CountMonth, err = s.Store.Repo().GetBajarildiCountMonth()
	if err != nil {
		s.Logger.Error("allInfo.GetPackingCountOfMonth: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.PlanDaily, err = s.Store.Repo().GetPlanToday()
	if err != nil {
		s.Logger.Error("allInfo.GetPlanToday: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	allInfo.CountDaily, err = s.Store.Repo().GetPlanCountToday()
	if err != nil {
		s.Logger.Error("allInfo.GetPlanCountToday: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	allInfo.PlanMonth, err = s.Store.Repo().GetPlanCountMonth()
	if err != nil {
		s.Logger.Error("allInfo.GetPlanCountMonth: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	allInfo.Counters = counters

	if allInfo.CountMonth > 0 && allInfo.PlanMonth > 0 {
		allInfo.PercentMonth = float64((allInfo.CountMonth * 100) / allInfo.PlanMonth)
	}

	if allInfo.CountDaily > 0 && allInfo.PlanDaily > 0 {
		allInfo.PercentDaily = float64((allInfo.CountDaily * 100) / allInfo.PlanDaily)
	}

	allInfo.CurrentMonth211 = s.Store.Repo().Get211ModelCurrentMonth()
	allInfo.CurrentMonth261 = s.Store.Repo().Get261ModelCurrentMonth()

	// s.Logger.Info("percent: ", float64((5000*100)/allInfo.Plan))

	resp.Result = "ok"
	resp.Data = allInfo

	c.JSON(200, resp)
}

func (s *Server) PlanUpdate(c *gin.Context) {
	resp := models.Responce{}

	id := c.GetInt("id")
	quantity := c.GetFloat64("quantity")

	if quantity < 0 {
		resp.Result = "error"
		resp.Err = "Noto'g'ri reja kiritildi"
		c.JSON(200, resp)
		return
	}

	err := s.Store.Repo().PlanUpdate(id, quantity)
	if err != nil {
		s.Logger.Error("PlanUpdate: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Result = "ok"
	// resp.Data = data

	c.JSON(200, resp)
}

func (s *Server) GetByMonthPlan(c *gin.Context) {
	resp := models.Responce{}

	month := c.GetString("date1")

	data, err := s.Store.Repo().GetByMonthPlan(month)
	if err != nil {
		s.Logger.Error("GetByMonthPlan: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Result = "ok"
	resp.Data = data

	c.JSON(200, resp)
}

func (s *Server) GetPlanToday(c *gin.Context) {
	resp := models.Responce{}

	// c_time := time.Now()
	// c_time.Month().
	data := ""
	// data = "2023" + "-" + c_time.Month().String() + "-" + strconv.Itoa(c_time.Day()) + " " + "8-00"
	data = "2023-04-12 8-00"
	time1, err := time.Parse("2023-04-12 8-00", data)
	if err != nil {
		fmt.Println("err: ", err)
	}
	time2 := time.Now()

	// fmt.Println("c_time: ", c_time)
	fmt.Println("time1: ", time1)
	fmt.Println("time2: ", time2)

	oraliq := time2.Sub(time1)

	fmt.Println("oraliq time: ", oraliq.Hours())

	// data, err := s.Store.Repo().GetPlanToday()
	// if err != nil {
	// 	s.Logger.Error("GetPlanToday: ", err)
	// 	resp.Result = "error"
	// 	resp.Err = err.Error()
	// 	c.JSON(200, resp)
	// 	return
	// }

	resp.Result = "ok"
	resp.Data = data

	c.JSON(200, resp)
}

func (s *Server) GetCurrentPlan(c *gin.Context) {
	resp := models.Responce{}
	data, err := s.Store.Repo().GetCurrentPlan()
	if err != nil {
		s.Logger.Error("GetCurrentPlan: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Result = "ok"
	resp.Data = data

	c.JSON(200, resp)
}

func (s *Server) CellAddComponent(c *gin.Context) {
	resp := models.Responce{}
	lot_id := c.GetInt("lot_id")
	component_id := c.GetInt("component_id")
	quantity := c.GetFloat64("quantity")
	cell_id := c.GetInt("cell_id")
	err := s.Store.Repo().CellAddComponent(quantity, component_id, lot_id, cell_id)
	if err != nil {
		s.Logger.Error("CellAddComponent: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	err = s.Store.Repo().UpdateMinusComponentIncome(component_id, quantity)
	if err != nil {
		s.Logger.Error("UpdateMinusComponentIncome: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Result = "ok"

	c.JSON(200, resp)
}
func (s *Server) CellGetAll(c *gin.Context) {
	resp := models.Responce{}

	data, err := s.Store.Repo().CellGetAll()
	if err != nil {
		s.Logger.Error("CellGetAll: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Result = "ok"
	resp.Data = data

	c.JSON(200, resp)
}
func (s *Server) CellGetByComponent(c *gin.Context) {
	resp := models.Responce{}

	component_id := c.GetInt("component_id")

	data, err := s.Store.Repo().CellGetByComponent(component_id)
	if err != nil {
		s.Logger.Error("CellGetAll: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Result = "ok"
	resp.Data = data

	c.JSON(200, resp)
}

func (s *Server) CellGetByComponentAll(c *gin.Context) {
	resp := models.Responce{}

	component_id := c.GetInt("component_id")

	data, err := s.Store.Repo().CellGetByComponentAll(component_id)
	if err != nil {
		s.Logger.Error("CellGetAll: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Result = "ok"
	resp.Data = data

	c.JSON(200, resp)
}

func (s *Server) CellGetEmpty(c *gin.Context) {
	resp := models.Responce{}
	// lot_id := c.GetInt("lot_id")
	component_id := c.GetInt("component_id")
	notFilter := []string{"890044818",
		"890112180",
		"890112185",
		"890112495",
		"890119914",
		"890171747",
		"890171942",
		"890192215",
		"890226858",
		"890231378",
		"890235373",
		"890242209",
		"890253279",
		"890254383",
		"890290041",
		"890290446",
		"890290773",
		"890292726",
		"890296695",
		"890296871",
		"890296951",
		"890304132",
		"890306273",
		"890308659",
		"890310001",
		"890310252",
		"890311664",
		"890314043",
		"890314047",
		"890315140",
		"890315147",
		"890317361",
		"890317363",
		"890317712",
		"890327833",
		"890333818",
		"890340264",
		"PR00100MB256",
		"PR00100SB260",
		"PR00100SO258",
		"PR211003K",
		"PR211261F",
		"PR211261I",
		"PR211261P",
		"PR211261S",
		"PR211261TH",
		"PR261003K",
		"PRP211261"}

	comp_name, err := s.Store.Repo().GetComponentName(component_id)
	if err != nil {
		s.Logger.Error("GetComponentName: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	skip_check := false

	fmt.Println("comp_name: ", comp_name)

	for i := 0; i < len(notFilter); i++ {
		if notFilter[i] == comp_name {
			skip_check = true
		}

	}

	fmt.Println("skip_check: ", skip_check)
	if skip_check == true {
		data, err := s.Store.Repo().CellGetNoFilter(comp_name)
		if err != nil {
			s.Logger.Error("CellGetEmpty: ", err)
			resp.Result = "error"
			resp.Err = err.Error()
			c.JSON(200, resp)
			return
		}
		resp.Result = "ok"
		resp.Data = data
	} else {
		data, err := s.Store.Repo().CellGetEmpty(component_id)
		if err != nil {
			s.Logger.Error("CellGetEmpty: ", err)
			resp.Result = "error"
			resp.Err = err.Error()
			c.JSON(200, resp)
			return
		}
		resp.Result = "ok"
		resp.Data = data
	}

	c.JSON(200, resp)
}

func (s *Server) AktReport(c *gin.Context) {
	resp := models.Responce{}
	date1 := c.GetString(("date1"))
	date2 := c.GetString(("date2"))

	data, err := s.Store.Repo().AktReport(date1, date2)
	if err != nil {
		s.Logger.Error("AktReport: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	c.JSON(200, data)
}

func (s *Server) AktInput(c *gin.Context) {

	account := models.Akt{}
	resp := models.Responce{}

	account.Component_id = c.GetInt("component_id")
	account.UserName = c.GetString("username")
	account.Comment = c.GetString("comment")
	account.Quantity = c.GetFloat64("quantity")
	account.Checkpoint_id = c.GetInt("checkpoint_id")
	account.Photo = c.GetString("photo")
	type_id := c.GetInt("type_id")
	// s.Logger.Info("photo: ", account.Photo)

	dec, err := base64.StdEncoding.DecodeString(account.Photo)
	if err != nil {
		s.Logger.Error("error1: ", err)
	}

	fileName := uuid.New().String() + ".jpg"
	f, err := os.Create(`g:\premier\server_V2\ware_spa\dist\assets\photo\` + fileName)
	// f, err := os.Create(`..\global\media\` + fileName)
	if err != nil {
		s.Logger.Error("error2: ", err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		s.Logger.Error("error3: ", err)
	}
	if err := f.Sync(); err != nil {
		s.Logger.Error("error4: ", err)
	}
	s.Logger.Info("user: ", account.UserName, " update sector: ", account.Checkpoint_id, account.Component_id, account.Quantity)

	if account.Checkpoint_id == 28 {

	} else {
		if err := s.Store.Repo().SectorBalanceUpdateByQuantity(account.Checkpoint_id, account.Component_id, account.Quantity); err != nil {
			s.Logger.Error("AktInput SectorBalanceUpdateByQuantity: ", err)
			resp.Result = "error"
			resp.Err = err.Error()
			c.JSON(200, resp)
			c.Abort()
			return
		}
	}

	err = s.Store.Repo().AktInput(account, fileName, type_id)
	if err != nil {
		s.Logger.Error("AktInput: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		c.Abort()
		return
	}

	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) AktInputWare(c *gin.Context) {

	account := models.Akt{}
	resp := models.Responce{}

	account.Component_id = c.GetInt("component_id")
	account.UserName = c.GetString("username")
	account.Comment = c.GetString("comment")
	account.Quantity = c.GetFloat64("quantity")
	account.Checkpoint_id = c.GetInt("checkpoint_id")
	account.Photo = c.GetString("photo")
	type_id := c.GetInt("type_id")
	cell_id := c.GetInt("cell_id")
	lot_id := c.GetInt("lot_id")
	// s.Logger.Info("Component_id: ", account.Component_id)
	// s.Logger.Info("UserName: ", account.UserName)
	// s.Logger.Info("Quantity: ", account.Quantity)
	// s.Logger.Info("type_id: ", type_id)
	s.Logger.Info("lot_id: ", lot_id)

	if err := s.Store.Repo().CellRemoveComponent(cell_id, account.Quantity); err != nil {
		s.Logger.Error("CellRemoveComponent: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	if err := s.Store.Repo().CellCheckEmpty(cell_id); err != nil {
		s.Logger.Error("CellCheckEmpty: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	if err := s.Store.Repo().RemoveComponentFromWare(account.Component_id, account.Quantity); err != nil {
		s.Logger.Error("RemoveComponentFromWare: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	dec, err := base64.StdEncoding.DecodeString(account.Photo)
	if err != nil {
		s.Logger.Error("error1: ", err)
	}

	fileName := uuid.New().String() + ".jpg"
	f, err := os.Create(`g:\premier\server_V2\ware_spa\dist\assets\photo\` + fileName)
	// f, err := os.Create(`..\global\media\` + fileName)
	if err != nil {
		s.Logger.Error("error2: ", err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		s.Logger.Error("error3: ", err)
	}
	if err := f.Sync(); err != nil {
		s.Logger.Error("error4: ", err)
	}
	s.Logger.Info("user: ", account.UserName, " update sector: ", account.Checkpoint_id, account.Component_id, account.Quantity)

	id, err := s.Store.Repo().AktInputWare(account, fileName, type_id, lot_id)
	if err != nil {
		s.Logger.Error("AktInput: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		c.Abort()
		return
	}

	t := time.Now()
	tString := t.Format("2006-01-02 15:04:05")
	tString += fmt.Sprintf(" Списано ID: %v", id)

	if err := s.Store.Repo().OutcomeInsert(account.Component_id, account.Checkpoint_id, account.Quantity, tString); err != nil {
		s.Logger.Error("OutcomeInsert: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		c.Abort()
		return
	}

	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) GetAllComponents(c *gin.Context) {

	data, err := s.Store.Repo().GetAllComponents()
	if err != nil {
		resp := models.Responce{}
		s.Logger.Error("GetAllComponents: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}
func (s *Server) GetAllComponentsOutCome(c *gin.Context) {

	data, err := s.Store.Repo().GetAllComponentsOutcome()
	if err != nil {
		resp := models.Responce{}
		s.Logger.Error("GetAllComponents: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GetAllComponentsOutComeByQuantity(c *gin.Context) {

	data, err := s.Store.Repo().GetAllComponentsOutComeByQuantity()
	if err != nil {
		resp := models.Responce{}
		s.Logger.Error("GetAllComponents: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}
func (s *Server) GetGPCompontents(c *gin.Context) {

	data, err := s.Store.Repo().GetGPCompontents()
	if err != nil {
		resp := models.Responce{}
		s.Logger.Error("GetGPCompontents: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) GPCompontentsAdd(c *gin.Context) {

	resp := models.Responce{}
	line := c.GetInt("checkpoint_id")
	component := c.GetInt("component_id")
	model := c.GetInt("model_id")

	err := s.Store.Repo().GPCompontentsAdd(line, component, model)
	if err != nil {
		s.Logger.Error("GPCompontentsAdd: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"

	c.JSON(200, resp)
}

func (s *Server) GPCompontentsRemove(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt("id")

	err := s.Store.Repo().GPCompontentsRemove(id)
	if err != nil {
		s.Logger.Error("GPCompontentsRemove: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"

	c.JSON(200, resp)
}

func (s *Server) GPCompontentsAdded(c *gin.Context) {

	resp := models.Responce{}
	components, err := s.Store.Repo().GPCompontentsAdded()
	if err != nil {
		s.Logger.Error("GPCompontentsAdded: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	resp.Data = components

	c.JSON(200, resp)
}

func (s *Server) GetCompoment(c *gin.Context) {

	id := c.GetInt("id")

	data, err := s.Store.Repo().GetComponent(id)
	if err != nil {
		resp := models.Responce{}
		s.Logger.Error("GetComponent: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
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
	component.InnerCode = c.GetString("inner_code")

	if component.ID == 0 {
		s.Logger.Error("GetComponent: ", "blank id")
		resp.Result = "error"
		resp.Err = "component id == 0"
		c.JSON(200, resp)
		return
	}

	err := s.Store.Repo().UpdateComponent(&component)
	if err != nil {
		s.Logger.Error("GetComponent: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) AddComponent(c *gin.Context) {
	component := models.Component{}
	resp := models.Responce{}

	component.ID = c.GetInt("id")
	component.Code = c.GetString("code")
	component.Name = c.GetString("name")
	component.Checkpoint_id = c.GetInt("checkpoint_id")
	component.Unit = c.GetString("unit")
	component.Specs = c.GetString("specs")
	component.Photo = c.GetString("photo")
	component.Type_id = c.GetInt("type_id")
	component.Weight = c.GetFloat64("weight")
	component.InnerCode = c.GetString("inner_code")
	s.Logger.Info(component)
	id, err := s.Store.Repo().AddComponent(&component)
	if err != nil {
		s.Logger.Error("AddComponent: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	s.Logger.Info("id: ", id)
	if id > 0 {
		err := s.Store.Repo().AddComponentsIncome(id)
		if err != nil {
			s.Logger.Error("AddComponentsIncome: ", err)
			resp.Result = "error"
			resp.Err = err.Error()
			c.JSON(200, resp)
		}
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
		resp.Err = err.Error()
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
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) DeleteCheckpoint(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt("checkpoint_id")

	err := s.Store.Repo().DeleteCheckpoint(id)
	if err != nil {
		s.Logger.Error("DeleteCheckpoints: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
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
		resp.Err = err.Error()
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
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) Income(c *gin.Context) {

	resp := models.Responce{}
	quantity := c.GetFloat64("quantity")
	id := c.GetInt("id")

	err := s.Store.Repo().Income(id, quantity)
	if err != nil {
		s.Logger.Error("Income: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	err = s.Store.Repo().UpdateComponentIncome(id, quantity)
	if err != nil {
		s.Logger.Error("UpdateComponentIncome: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) IncomeReport(c *gin.Context) {
	resp := models.Responce{}
	date1 := c.GetString("date1")
	date2 := c.GetString("date2")

	data, err := s.Store.Repo().IncomeReport(date1, date2)
	if err != nil {
		s.Logger.Error("IncomeReport: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	c.JSON(200, data)
}

func (s *Server) Types(c *gin.Context) {

	resp := models.Responce{}

	data, err := s.Store.Repo().Types()
	if err != nil {
		s.Logger.Error("Types: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) Models(c *gin.Context) {

	resp := models.Responce{}

	data, err := s.Store.Repo().Models()
	if err != nil {
		s.Logger.Error("Models: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) Model(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt(("id"))

	data, err := s.Store.Repo().Model(id)
	if err != nil {
		s.Logger.Error("Model: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) InsertUpdateModel(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt(("id"))
	code := c.GetString(("code"))
	comment := c.GetString(("comment"))
	name := c.GetString(("name"))

	err := s.Store.Repo().InsertUpdateModel(name, code, comment, id)
	if err != nil {
		s.Logger.Error("InsertUpdateModel: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) OutcomeModelCheck(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt(("id"))
	quantity := c.GetFloat64(("quantity"))

	data, err := s.Store.Repo().OutcomeModelCheck(id, quantity)
	if err != nil {
		s.Logger.Error("OutcomeModelCheck: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) OutcomeComponentCheck(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt(("id"))
	quantity := c.GetFloat64(("quantity"))

	data, err := s.Store.Repo().OutcomeComponentCheck(id, quantity)
	if err != nil {
		s.Logger.Error("OutcomeComponentCheck: ", err)
		resp.Result = "error"
		resp.Data = data
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Data = data
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) OutcomeComponentSubmit(c *gin.Context) {

	resp := models.Responce{}
	component_id := c.GetInt("component_id")
	checkpoint_id := c.GetInt("checkpoint_id")
	quantity := c.GetFloat64("quantity")
	cell_id := c.GetInt("cell_id")

	err := s.Store.Repo().OutcomeComponentSubmit(component_id, checkpoint_id, quantity)
	if err != nil {
		s.Logger.Error("OutcomeComponentSubmit: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	if err = s.Store.Repo().CellRemoveComponent(cell_id, quantity); err != nil {
		s.Logger.Error("CellRemoveComponent: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	if err = s.Store.Repo().CellCheckEmpty(cell_id); err != nil {
		s.Logger.Error("CellCheckEmpty: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) OutcomeModelSubmit(c *gin.Context) {

	resp := models.Responce{}
	model_id := c.GetInt(("model_id"))
	quantity := c.GetFloat64(("quantity"))

	err := s.Store.Repo().OutcomeModelSubmit(model_id, quantity)
	if err != nil {
		s.Logger.Error("OutcomeModelSubmit: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) OutcomeReport(c *gin.Context) {
	resp := models.Responce{}
	date1 := c.GetString(("date1"))
	date2 := c.GetString(("date2"))

	data, err := s.Store.Repo().OutcomeReport(date1, date2)
	if err != nil {
		s.Logger.Error("OutcomeReport: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	c.JSON(200, data)
}

func (s *Server) OutcomeFile(c *gin.Context) {
	s.Logger.Info("outcome file")
	resp := models.Responce{}

	type Form struct {
		File *multipart.FileHeader `form:"excel" binding:"required"`
	}

	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		s.Logger.Error("OutcomeFile: ", err)
		resp.Result = "error"
		resp.Err = err
		c.JSON(200, resp)
		return
	}
	// Get raw file bytes - no reader method
	// openedFile, _ := form.File.Open()
	// file, _ := ioutil.ReadAll(openedFile)
	c.SaveUploadedFile(form.File, "temp.xlsx")
	// myString := string(file[:])

	var file []models.FileInput
	x, _ := xlsx.New(xlsx.WithInputFile("temp.xlsx"))
	defer x.Close()

	if err := x.Read(&file); err != nil {
		s.Logger.Error(err)
	}
	// fmt.Println(file)
	res, err := s.Store.Repo().FileInput(file)
	if err != nil {
		s.Logger.Error(err)
		resp.Result = "error"
		resp.Data = res
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}

	// defer openedFile.Close()
	resp.Result = "ok"
	resp.Data = res
	c.JSON(200, resp)
}

func (s *Server) BomComponentInfo(c *gin.Context) {

	resp := models.Responce{}
	id := c.GetInt(("id"))

	data, err := s.Store.Repo().BomComponentInfo(id)
	if err != nil {
		s.Logger.Error("BomComponentInfo: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}

func (s *Server) BomComponentAdd(c *gin.Context) {

	resp := models.Responce{}
	component_id := c.GetInt(("id"))
	model_id := c.GetInt(("model_id"))
	quantity := c.GetFloat64(("quantity"))
	comment := c.GetString(("id"))

	err := s.Store.Repo().BomComponentAdd(component_id, model_id, quantity, comment)
	if err != nil {
		s.Logger.Error("BomComponentAdd: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func (s *Server) BomComponentDelete(c *gin.Context) {

	resp := models.Responce{}
	component_id := c.GetInt(("component_id"))
	model_id := c.GetInt(("model_id"))

	err := s.Store.Repo().BomComponentDelete(component_id, model_id)
	if err != nil {
		s.Logger.Error("BomComponentDelete: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	resp.Result = "ok"
	c.JSON(200, resp)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (s *Server) GsCodeFile(c *gin.Context) {

	s.Logger.Info("GS code file")
	resp := models.Responce{}

	type FileToken struct {
		Model int    `json:"model"`
		Token string `json:"token"`
	}

	type Form struct {
		File  *multipart.FileHeader `form:"gscode" binding:"required"`
		File2 *multipart.FileHeader `form:"data" binding:"required"`
	}

	var form Form
	err := c.ShouldBind(&form)
	if err != nil {
		s.Logger.Error("OutcomeFile: ", err)
		resp.Result = "error"
		resp.Err = err
		c.JSON(200, resp)
		return
	}
	c.SaveUploadedFile(form.File, "temp.csv")
	c.SaveUploadedFile(form.File2, "temp.json")

	plan, _ := ioutil.ReadFile("temp.json")
	data := &FileToken{}
	err = json.Unmarshal(plan, &data)
	if err != nil {
		s.Logger.Error(err)
	}

	parsedToken, err := ParseToken(data.Token)
	if err != nil {
		s.Logger.Error("WareCheckRole Wrong Token: ", data.Token, " error: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(401, resp)
		c.Abort()
		return
	}

	res, err := s.Store.Repo().CheckRole(c.Request.URL.String(), parsedToken.Email)
	if err != nil {
		s.Logger.Error("WareCheckRole: ", data.Token, " error: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(401, resp)
		c.Abort()
		return
	}

	if !res {
		s.Logger.Error("WareCheckRole: ", data.Token, " error: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(401, resp)
		c.Abort()
		return
	}

	stringArray, err := readLines("temp.csv")
	if err != nil {
		s.Logger.Error("string err: ", err)
	}

	s.Logger.Info("string array len: ", len(stringArray))

	badCode, err := s.Store.Repo().InsertGsCode(stringArray, data.Model)

	if err != nil {
		for i := range badCode {
			s.Logger.Error("error key: ", i)
		}
		resp.Result = "error"
		resp.Data = badCode
		c.JSON(200, resp)
		return
	}
	s.Logger.Info("added keys: ", len(stringArray), " model: ", data.Model)

	// file, err := os.Open("temp.csv")
	// if err != nil {
	// 	s.Logger.Error(err)
	// }

	// reader := csv.NewReader(file)
	// reader.Comma = '@'
	// reader.LazyQuotes = true
	// for {
	// 	record, err := reader.Read()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		s.Logger.Error("csv decoding error: ", err)
	// 	}
	// res := strings.ReplaceAll(record[0], "", "")
	// if err := s.Store.Repo().InsertGsCode(res, data.Model); err != nil {
	// 	resp.Result = "error"
	// 	resp.Err = err
	// 	c.JSON(200, resp)
	// 	return
	// }
	// 	// s.Logger.Info(record[0])

	// }
	// defer file.Close()
	resp.Result = "ok"

	c.JSON(200, resp)
}

func (s *Server) GetKeys(c *gin.Context) {

	resp := models.Responce{}
	data, err := s.Store.Repo().GetKeys()
	if err != nil {
		s.Logger.Error("GetKeys: ", err)
		resp.Result = "error"
		resp.Err = err.Error()
		c.JSON(200, resp)
		return
	}
	c.JSON(200, data)
}
