package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func doSomething(s int) ([]byte,  error) {
	logData := map[string]interface{}{
		"level":   "info",
		"time":    time.Now().Format(time.RFC3339),
		"message": "Processing input",
		"context": map[string]interface{}{
			"type": "this is a int",
		},
	}
	
	x, err := json.Marshal(logData)
	if err != nil {
		log.Error().Err(err).Msg("Error occurred in doSomething")
		return createErrorResponse(err), err
	}
	if s == 0 {
		log.Info().Str("type", "this is an int").Msg("Processing input 0")
		return x, errors.New("sample error for log 0")
	}
	if s == 1 {
		log.Info().Str("type", "this is an int").Msg("Processing input 1")
		return x, errors.New("sample error for log 1")
	}
	return x, nil
}

type ErrorResponse struct {
	Error error `json:"error"`
}

func createErrorResponse(err error) []byte {
	errorResponse := ErrorResponse{
		Error: err,
	}
	jsonResponse, _ := json.Marshal(errorResponse)
	return jsonResponse
}


type LogEntry struct {
    Level     string `json:"level"`
    Timestamp string `json:"time"`
    Message   string `json:"message"`
    Context   map[string]interface{} `json:"context"`
}

type LogEntryResponse struct {
    Level     string 
	Timestamp time.Time 
    Message   string 
    Context   []byte
}

func Write(p []byte)(LogEntryResponse , []byte){
	var entry LogEntry
	var err error
	if err = json.Unmarshal(p, &entry); err != nil {
		return LogEntryResponse{} , createErrorResponse(err)
	}

	level := entry.Level
	timestamp, err := time.Parse(time.RFC3339, entry.Timestamp)
	if err != nil {
		return LogEntryResponse{}, createErrorResponse(err)
	}
	message := entry.Message
	context, err := json.Marshal(entry.Context)
	if err != nil {
		return LogEntryResponse{}, createErrorResponse(err)
	}
	data := LogEntryResponse{
		Level: level,
		Timestamp: timestamp,
		Message: message,
		Context: context,
	}

	return data, nil
}

func main(){
	if d, err := doSomething(1); err != nil {
		r , err := Write(d)
		if err != nil {
			fmt.Println("Error:",string(err))
		}
		fmt.Printf("log details %+v\n", r)
		
	}

	if  d, err := doSomething(0); err != nil {
		r , err := Write(d)
		if err != nil {
			fmt.Println("Error:",string(err))	
		}
		fmt.Printf("log details:%+v\n",r)
		
	}
}
