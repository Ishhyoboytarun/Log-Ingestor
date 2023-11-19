package handler

import (
	"Log-Ingestor/repo"
	"Log-Ingestor/service"
	"database/sql"
)

type Handlers struct {
	Injest *InjestLog
}

var Handler *Handlers

func Init() {

	Handler = new(Handlers)
	Db := new(sql.DB)
	Repo := repo.NewInjestLog(Db)
	Service := service.NewInjestLog(Repo)
	InjestLogHandler := NewInjestLogs(Service)
	Handler.Injest = InjestLogHandler
}
