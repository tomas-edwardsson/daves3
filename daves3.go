package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func authenticate(r *http.Request, username string, password string) bool {
	requestUsername, requestPassword, ok := r.BasicAuth()
	if username == requestUsername &&
		password == requestPassword && ok == true {
		return true
	}
	return false
}

func logger(logChan chan string, msg string) {
	logChan <- fmt.Sprintf("%s: %s", time.Now().String(), msg)
}

func handleFunc(
	sess *session.Session,
	logChan chan string,
	s3bucket string,
	username string,
	password string) func(w http.ResponseWriter, r *http.Request) {
	svc := s3.New(sess)
	uploader := s3manager.NewUploader(sess)
	return func(w http.ResponseWriter, r *http.Request) {
		logger(logChan, fmt.Sprintf("%s %s", r.Method, r.URL))

		if authenticate(r, username, password) == false {
			w.Header().Add("WWW-Authenticate", `Basic realm="daves3"`)
			http.Error(w, "Bad auth", 401)
			return
		}
		if r.Method == "PUT" {
			result, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(s3bucket),
				Key:    aws.String(r.URL.EscapedPath()),
				Body:   r.Body,
			})
			if err != nil {
				logger(logChan, err.Error())
				http.Error(w, err.Error(), 500)
				return
			}
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "Created %s", result.UploadID)
			return
		}
		if r.Method == "HEAD" {
			_, err := svc.HeadObject(&s3.HeadObjectInput{
				Bucket: aws.String(s3bucket),
				Key:    aws.String(r.URL.EscapedPath()),
			})
			if err != nil {
				http.Error(w, err.Error(), 404)
				return
			}
			return
		}
		if r.Method == "GET" {
			req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(s3bucket),
				Key:    aws.String(r.URL.EscapedPath()),
			})
			str, err := req.Presign(10 * time.Second)
			if err != nil {
				logger(logChan, err.Error())
				http.Error(w, err.Error(), 500)
				return
			}
			http.Redirect(w, r, str, 302)
			return
		}
		http.NotFound(w, r)
		return
	}
}

func logStdout(logChan chan string) {
	for {
		fmt.Println(<-logChan)
	}
}

func main() {
	sess := session.Must(session.NewSession())

	logChan := make(chan string, 512)
	username := os.Getenv("DAVES3_USERNAME")
	password := os.Getenv("DAVES3_PASSWORD")
	s3bucket := os.Getenv("DAVES3_BUCKET")

	if username == "" {
		fmt.Println("DAVES3_USERNAME empty")
		os.Exit(3)
	}
	if password == "" {
		fmt.Println("DAVES3_PASSWORD empty")
		os.Exit(3)
	}
	if s3bucket == "" {
		fmt.Println("DAVES3_BUCKET empty")
		os.Exit(3)
	}
	go logStdout(logChan)
	http.HandleFunc("/", handleFunc(sess, logChan, s3bucket, username, password))
	http.ListenAndServe(":8080", nil)
}
