package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // регистрирует /debug/pprof/*

	"Prak_13/internal/work"
)

func main() {
	// Используем DefaultServeMux: пакет net/http/pprof регистрирует
	// свои хендлеры в нём при импорте _ "net/http/pprof".
	// Эндпоинт, вызывающий “тяжёлую” работу.
	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		defer work.TimeIt("Fib(38)")()
		res := work.Fib(38)
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte((fmtInt(res))))
	})

	log.Println("Server on :8080; pprof on /debug/pprof/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fmtInt(v int) string { return fmt.Sprintf("%d\n", v) }
