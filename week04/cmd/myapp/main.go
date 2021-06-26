package main

import (
	setting "Week04/internal/pkg"
	"Week04/internal/service"
	"fmt"
	"net/http"
)

func main() {
	router := service.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}