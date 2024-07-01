package main

import (
	"database/sql"
	"fmt"
	"livecode-catatan-keuangan/config"
	"livecode-catatan-keuangan/controller"
	"livecode-catatan-keuangan/middleware"
	"livecode-catatan-keuangan/repository"
	"livecode-catatan-keuangan/service"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// server ini menghubungkan semua komponen
type Server struct {
	uS      service.UserService
	jS      service.JwtService
	aM      middleware.AuthMiddleware
	engine  *gin.Engine //untuk start engine gin
	portApp string
}

// method untuk memanggil route yang di controller
func (s *Server) initiateRoute() {
	//bisa menambah grouping lagi disini
	routerGroup := s.engine.Group("/api/v1")
	controller.NewUserController(s.uS, routerGroup).Route()
}

// func untuk running
func (s *Server) Start() {
	s.initiateRoute()
	s.engine.Run(s.portApp)
}

// constructur, agar dipanggil main.go
func NewServer() *Server {
	//memanggil hasil config .env
	co, _ := config.NewConfig()

	//melakukan koneksi database
	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", co.Host, co.Port, co.User, co.Password, co.Name)

	db, err := sql.Open(co.Driver, urlConnection)
	if err != nil {
		log.Fatal(err)
	}

	portApp := co.AppPort
	userRepo := repository.NewUserRepository(db)

	jwtService := service.NewJwtService(co.SecurityConfig)
	userService := service.NewUserService(userRepo, jwtService)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	//menginject repo ke service
	return &Server{
		uS:      userService,
		jS:      jwtService,
		aM:      authMiddleware,
		portApp: portApp,
		engine:  gin.Default(),
	}
}
