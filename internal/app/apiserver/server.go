package apiserver

import (
	"io"
	"os"
	"warehouse/internal/app/store/sqlstore"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Secret_key = []byte("Some123SecretKeyPremier1")

type Server struct {
	Router *gin.Engine
	Logger *log.Logger
	Store  sqlstore.Store
}

func newServer(store sqlstore.Store) *Server {
	s := &Server{
		Router: gin.New(),
		Logger: log.New(),
		Store:  store,
	}
	f, err := os.OpenFile("logger.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	wrt := io.MultiWriter(os.Stdout, f)

	// s.Router.Use(cors.Default())

	s.Logger.SetOutput(wrt)
	s.Logger.SetFormatter(&log.JSONFormatter{})
	s.configureRouter()
	return s
}

func (s *Server) configureRouter() {
	s.Router.SetTrustedProxies([]string{"localhost"})

	s.Router.POST("users/login", s.Login)             // {"email": string, "password": string}
	s.Router.POST("ware/outcome/file", s.OutcomeFile) //only excel file input
	s.Router.POST("ware/gscode/file", s.GsCodeFile)

	ware := s.Router.Group("/ware") //route for warehouse control
	ware.Use(s.WareCheckRole())
	{
		ware.POST("/components", s.GetAllComponents)                          // {"token": string}
		ware.POST("/components/outcome", s.GetAllComponentsOutCome)           // {"token": string}
		ware.POST("/component", s.GetCompoment)                               // {"id": int, "token": string}
		ware.POST("/component/update", s.UpdateCompoment)                     // {"code":string, "name":string, "checkpoint_id":int, "unit":string, "photo":string, "specs":string, "type_id":int, "weight":float64, "id":int, "token": string}
		ware.POST("/component/add", s.AddComponent)                           // {"code":string, "name":string, "checkpoint_id":int, "unit":string, "photo":string, "specs":string, "type_id":int, "weight":float64, "token": string}
		ware.POST("/component/delete", s.DeleteCompoment)                     // {"id":int, "token": string}
		ware.POST("/checkpoints", s.GetAllCheckpoints)                        // {"token": string}
		ware.POST("/checkpoint/delete", s.DeleteCheckpoint)                   // {"id":int, "token": string}
		ware.POST("/checkpoint/add", s.AddCheckpoint)                         // {"name":string, "photo":string, "token": string}
		ware.POST("/checkpoint/update", s.UpdateCheckpoint)                   // {"name":string, "photo":string, "id":int, "token": string}
		ware.POST("/income", s.Income)                                        // {"component_id":int, "quantity":int, "token": string}
		ware.POST("/income/report", s.IncomeReport)                           // {"date1":string, "date2":string, "token": string}
		ware.POST("/types", s.Types)                                          // {"token": string}
		ware.POST("/models", s.Models)                                        // {"token": string}
		ware.POST("/model", s.Model)                                          // {"id":int, "token": string}
		ware.POST("/outcome/model/check", s.OutcomeModelCheck)                // {"id":int, "token": string}
		ware.POST("/outcome/model/submit", s.OutcomeModelSubmit)              // {"model_id":int, "quantity":float64, "token": string}
		ware.POST("/outcome/component/check", s.OutcomeComponentCheck)        // {"component_id":int, "quantity":float64, "token": string}
		ware.POST("/outcome/component/submit", s.OutcomeComponentSubmit)      // {"component_id":int, "checkpoint_id":int, "quantity":float64, "token": string}
		ware.POST("/outcome/report", s.OutcomeReport)                         // {"date1":string, "date2":string, "token": string}
		ware.POST("/model/update", s.InsertUpdateModel)                       // {"id":int, "code"string, "comment":string, "name":string, "token": string}
		ware.POST("/bom/component", s.BomComponentInfo)                       // {"id":int, "token": string}
		ware.POST("/bom/component/add", s.BomComponentAdd)                    // {"id":int, "token": string}
		ware.POST("/bom/component/delete", s.BomComponentDelete)              // {"model_id":int, "component_id":int, "token": string}
		ware.POST("/production/sector/balance", s.GetSectorBalance)           // {"line":int, "token": string}
		ware.POST("/production/sector/balance/update", s.SectorBalanceUpdate) // {"line":int, "component_id": int, "quantity": float64, "token": string}
		ware.POST("/gscode/get", s.GetKeys)                                   // {"token": string}
		ware.POST("/akt/input", s.AktInput)                                   // {"token": string, "component_id": int, "data": string, "quantity": float64, "checkpoint_id": int} data => comment
	}

	global := s.Router.Group("/api") //Route for global use
	global.Use(s.CheckRole())
	{
		global.POST("/production/last", s.GetLast)                                 // {"line": int, "token": string}
		global.POST("/production/status", s.GetStatus)                             // {"line": int, "token": string}
		global.POST("/production/today", s.GetToday)                               // {"line": int, "token": string}
		global.POST("/production/today/models", s.GetTodayModels)                  // {"line": int, "token": string}
		global.POST("/production/sector/balance", s.GetSectorBalance)              // {"line": int, "token": string}
		global.POST("/production/packing/last", s.GetPackingLast)                  // {"token": string}
		global.POST("/production/packing/today", s.GetPackingToday)                // {"token": string}
		global.POST("/production/packing/today/serial", s.GetPackingTodaySerial)   // {"token": string}
		global.POST("/production/packing/today/models", s.GetPackingTodayModels)   // {"token": string}
		global.POST("/production/lines", s.GetLines)                               // {"token": string}
		global.POST("/production/defects/types", s.GetDefectsTypes)                // {"token": string}
		global.POST("/production/defects/types/delete", s.DeleteDefectsTypes)      // {"id": int,"token": string}
		global.POST("/production/defects/types/add", s.AddDefectsTypes)            // {"id": int,"token": string}
		global.POST("/production/defects/add", s.AddDefects)                       // {"serial": string, "checkpoint_id": int, "defect_id": int, "token": string}
		global.POST("/production/defects/last", s.Last3Defects)                    // {"serial": string, "checkpoint_id": int, "defect_id": int, "token": string}
		global.POST("/production/report/bydate/models/serial", s.GetByDateSerial)  // {"date1": string, "date2": string, "line": int, "token": string}
		global.POST("/production/report/bydate", s.GetCountByDate)                 // {"date1": string, "date2": string, "line": int, "token": string}
		global.POST("/production/report/bydate/models", s.GetByDateModels)         // {"date1": string, "date2": string, "line": int, "token": string}
		global.POST("/production/report/remont", s.GetRemont)                      // {"token": string}
		global.POST("/production/report/remont/today", s.GetRemontToday)           // {"token": string}
		global.POST("/production/report/remont/bydate", s.GetRemontByDate)         // {"date1": string, "date2": string, "token": string}
		global.POST("/production/report/remont/update", s.UpdateRemont)            // {"name": string, "id": int, "token": string} id-> defect id
		global.POST("/production/report/remont/repairedcount", s.GetRepairedCount) // {"name": string, "id": int, "token": string} id-> defect id
		global.POST("/production/serial/info", s.GetInfoBySerial)                  // {"serial": string, "token": string}
		global.POST("/production/galileo/todaymodels", s.GalileoTodayModels)       // {"token": string}
		global.POST("/users/register", s.Create)                                   // {"email":string, "password":string,"token": string}

	}

	production := s.Router.Group("/production") //ONLY FOR PRODUCTION(factory) without check token and role
	production.Use(s.NoCheckRole())
	{
		production.POST("/last", s.GetLast)                                 // {"line": int}
		production.POST("/status", s.GetStatus)                             // {"line": int}
		production.POST("/today", s.GetToday)                               // {"line": int}
		production.POST("/today/models", s.GetTodayModels)                  // {"line": int}
		production.POST("/sector/balance", s.GetSectorBalance)              // {"line": int}
		production.POST("/packing/last", s.GetPackingLast)                  // {}
		production.POST("/packing/today", s.GetPackingToday)                // {}
		production.POST("/packing/today/serial", s.GetPackingTodaySerial)   // {}
		production.POST("/packing/today/models", s.GetPackingTodayModels)   // {}
		production.POST("/packing/serial/input", s.PackingSerialInput)      // {"serial":string, "packing":string}
		production.POST("/lines", s.GetLines)                               // {}
		production.POST("/defects/types", s.GetDefectsTypes)                // {}
		production.POST("/defects/types/delete", s.DeleteDefectsTypes)      // {"id": int}
		production.POST("/defects/types/add", s.AddDefectsTypes)            // {"id": int}
		production.POST("/defects/add", s.AddDefects)                       // {"serial": string, "checkpoint_id": int, "defect_id": int}
		production.POST("/report/bydate/models/serial", s.GetByDateSerial)  // {"date1": string, "date2": string, "line": int}
		production.POST("/report/bydate", s.GetCountByDate)                 // {"date1": string, "date2": string, "line": int}
		production.POST("/report/bydate/models", s.GetByDateModels)         // {"date1": string, "date2": string, "line": int}
		production.POST("/report/remont", s.GetRemont)                      // {}
		production.POST("/report/remont/repairedcount", s.GetRepairedCount) // {}
		production.POST("/report/remont/today", s.GetRemontToday)           // {}
		production.POST("/report/remont/bydate", s.GetRemontByDate)         // {"date1": string, "date2": string}
		production.POST("/report/remont/update", s.UpdateRemont)            // {"name": string, "id": int} id-> defect id
		production.POST("/serial/input", s.SerialInput)                     // {"serial": string, "line": int}
		production.POST("/serial/info", s.GetInfoBySerial)                  // {"serial": string}
		production.POST("/galileo/todaymodels", s.GalileoTodayModels)       // {}
		production.POST("/models", s.Models)                                // {}
		production.POST("/metall/serial", s.MetallSerial)                   // {"id", int}
		// production.POST("/galileo/tcp", s.GalileoTCP)                      // {"id", int}
	}
	s.Router.POST("galileo/input", s.GalileoInput)
}
