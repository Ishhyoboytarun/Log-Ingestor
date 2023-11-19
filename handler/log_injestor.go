package handler

import (
	"Log-Ingestor/service"
	"net/http"
)

type InjestLog struct {
	service.LogInjestorService
}

func NewInjestLogs(service service.LogInjestorService) *InjestLog {
	return &InjestLog{
		service,
	}
}

func (i *InjestLog) InjestLogs(w http.ResponseWriter, req *http.Request) {

	response := i.LogInjestorService.InjestLogs()
	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Status))
}
