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

	s.Logger.SetOutput(wrt)
	s.Logger.SetFormatter(&log.JSONFormatter{})
	s.configureRouter()
	return s
}

func (s *Server) configureRouter() {
	s.Router.POST("users/create", s.Create)
	s.Router.POST("users/login", s.Login)

	ware := s.Router.Group("/ware") //route for warehouse control
	ware.Use(s.CheckRole())
	{
		ware.POST("/components", s.GetAllComponents) // {"token": string}
		ware.POST("/component", s.GetCompoment)      // {"id": int, "token": string}
	}

	global := s.Router.Group("/api") //Route for global use
	global.Use(s.CheckRole())
	{
		global.POST("/production/last", s.GetLast)                                // {"line": int, "token": string}
		global.POST("/production/status", s.GetStatus)                            // {"line": int, "token": string}
		global.POST("/production/today", s.GetToday)                              // {"line": int, "token": string}
		global.POST("/production/today/models", s.GetTodayModels)                 // {"line": int, "token": string}
		global.POST("/production/sector/balance", s.GetSectorBalance)             // {"line": int, "token": string}
		global.POST("/production/packing/last", s.GetPackingLast)                 // {"token": string}
		global.POST("/production/packing/today", s.GetPackingToday)               // {"token": string}
		global.POST("/production/lines", s.GetLines)                              // {"token": string}
		global.POST("/production/defects/types", s.GetDefectsTypes)               // {"token": string}
		global.POST("/production/defects/types/delete", s.DeleteDefectsTypes)     // {"id": int,"token": string}
		global.POST("/production/defects/types/add", s.AddDefectsTypes)           // {"id": int,"token": string}
		global.POST("/production/defects/add", s.AddDefects)                      // {"serial": string, "checkpoint_id": int, "defect_id": int, "token": string}
		global.POST("/production/report/bydate/models/serial", s.GetByDateSerial) // {"date1": string, "date2": string, "line": int, "token": string}

	}

	production := s.Router.Group("/production") //ONLY FOR PRODUCTION(factory) without check token and role
	production.Use(s.NoCheckRole())
	{
		production.POST("/last", s.GetLast)                                // {"line": int}
		production.POST("/status", s.GetStatus)                            // {"line": int}
		production.POST("/today", s.GetToday)                              // {"line": int}
		production.POST("/today/models", s.GetTodayModels)                 // {"line": int}
		production.POST("/sector/balance", s.GetSectorBalance)             // {"line": int}
		production.POST("/packing/last", s.GetPackingLast)                 // {}
		production.POST("/packing/today", s.GetPackingToday)               // {}
		production.POST("/packing/today/serial", s.GetPackingTodaySerial)  // {}
		production.POST("/packing/today/models", s.GetPackingTodayModels)  // {}
		production.POST("/lines", s.GetLines)                              // {}
		production.POST("/defects/types", s.GetDefectsTypes)               // {}
		production.POST("/defects/types/delete", s.DeleteDefectsTypes)     // {"id": int}
		production.POST("/defects/types/add", s.AddDefectsTypes)           // {"id": int}
		production.POST("/defects/add", s.AddDefects)                      // {"serial": string, "checkpoint_id": int, "defect_id": int}
		production.POST("/report/bydate/models/serial", s.GetByDateSerial) // {"date1": string, "date2": string, "line": int, "token": string}
	}

}
