package service

import (
	"Log-Ingestor/contracts"
	"Log-Ingestor/repo"
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

type InjestLogService struct {
	repo.LogInjestorRepo
}

func NewInjestLog(repo repo.LogInjestorRepo) LogInjestorService {
	return &InjestLogService{
		repo,
	}
}

type LogInjestorService interface {
	InjestLogs() http.Response
}

func (i InjestLogService) InjestLogs() http.Response {

	wg := new(sync.WaitGroup)
	filePath := "logs.txt"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return http.Response{
			Status:     "InternalServerError",
			StatusCode: 500,
		}
	}
	defer file.Close()

	logsChannel := make(chan *contracts.LogEntry)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go readLogsWorker(i, file, 10, logsChannel, wg)
	}
	go func() {
		wg.Wait()
		close(logsChannel)
	}()

	logs := make([]*contracts.LogEntry, 0)
	for _ = range logsChannel {
		logs = append(logs, <-logsChannel)
	}
	return i.LogInjestorRepo.InjestLogs(logs)
}

func readLogsWorker(id int, file *os.File, numGoroutines int, logsChannel chan *contracts.LogEntry, wg *sync.WaitGroup) {

	defer wg.Done()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := fileInfo.Size()
	chunkSize := fileSize / int64(numGoroutines)
	startOffset := int64(id) * chunkSize
	endOffset := (int64(id) + 1) * chunkSize

	_, err = file.Seek(startOffset, 0)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var logEntry *contracts.LogEntry
		err := json.Unmarshal([]byte(line), &logEntry)
		if err != nil {
			log.Fatal(err)
		}
		logsChannel <- logEntry

		offset, err := file.Seek(0, 1)
		if err != nil {
			return
		}
		if offset >= endOffset {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
