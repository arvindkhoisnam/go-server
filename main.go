package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type Response struct{
	Message string `json:"message"`
}

func healthCheckpoint(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	 temp :=  Response{
		Message: "Abenao thibaikooooonnn 🐽🐽",
	 }
	 json.NewEncoder(w).Encode(temp)
}

func main(){
	ctx,stop := signal.NotifyContext(context.Background(),syscall.SIGINT,syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.Handle("GET /",LoggerMiddleware(http.HandlerFunc(healthCheckpoint)))
	gracefulServer := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	go func ()  {
		log.Println("Listening on port 8080")
		if err := gracefulServer.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			log.Fatal(err)
		}	
	}()

	<- ctx.Done()
	log.Println("Server shutting down.")

	timerCtx,close := context.WithTimeout(context.Background(), 5 *time.Second)
	defer close()

	if err := gracefulServer.Shutdown(timerCtx); err != nil{
		fmt.Printf("Error shutting down %s",err.Error())
	}

	log.Println("Server successfully shutdown.")
}

func LoggerMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.Header.Get("X-Forwarded-For")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}
		log.Println("Request from:", ip)
		next.ServeHTTP(w,r)
	})
}