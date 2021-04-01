package logger

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

// ServiceRequestHttpLog is log for http request
func ServiceRequestHttpLog(url *url.URL,header http.Header,body string) {
	log.Println("Request:")
	log.Println(url)
	log.Println("Header:", header)
	log.Println("Body", body)
}

// ServiceResponseLog is log for service response
func ServiceResponseLog(url string, body string) {
	dateFmt := "2006-01-02T15:04:05-07:00"
	date := time.Now()
	colorGreen := "\033[36m"
	colorReset := "\033[0m"

	fmt.Println(string(colorGreen),"INFO"+string(colorReset),fmt.Sprintf("[%s]",date.Format(dateFmt))+"Response Body: ", body)

}

// ServiceRequestLog is log for service request
func ServiceRequestLog(url string, req *http.Request,body string) {
	dateFmt := "2006-01-02T15:04:05-07:00"
	date := time.Now()
	cyan := "\033[36m"
	colorReset := "\033[0m"

	fmt.Println(string(cyan), "INFO"+string(colorReset), fmt.Sprintf("[%s]", date.Format(dateFmt))+"Request:")
	fmt.Println(string(cyan), "INFO"+string(colorReset), fmt.Sprintf("[%s]", date.Format(dateFmt))+url)
	fmt.Println(string(cyan), "INFO"+string(colorReset), fmt.Sprintf("[%s]", date.Format(dateFmt))+fmt.Sprintf("%s", req.Header))
	fmt.Println(string(cyan), "INFO"+string(colorReset), fmt.Sprintf("[%s]", date.Format(dateFmt))+body)

}