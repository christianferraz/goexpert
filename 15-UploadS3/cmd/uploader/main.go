package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       sync.WaitGroup
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		panic(err)
	}
	s3Client = s3.New(sess)
	s3Bucket = os.Getenv("S3_BUCKET")
}

func main() {
	dir, err := os.Open("../../tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()
	uploadControl := make(chan struct{}, 100)
	// se houver erro, indica qual arquivo falhou
	errorFileUpload := make(chan string)
	go func() {
		for {
			select {
			case filename := <-errorFileUpload:
				fmt.Printf("Error uploading file: %s\nRetrying\n", filename)
				uploadControl <- struct{}{}
				wg.Add(1)
				go uploadFile(filename, uploadControl, errorFileUpload)
			}
		}
	}()
	for {
		files, err := dir.ReadDir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading directory: %s\n", err)
			continue
		}
		wg.Add(1)
		// coloca uma struct vazia, vai liberar até 100 channels e esperar ele esvaziar
		uploadControl <- struct{}{}
		go uploadFile(files[0].Name(), uploadControl, errorFileUpload)
	}
	wg.Wait()
}

func uploadFile(filename string, uploadControl <-chan struct{}, errorFileUpload chan<- string) {
	defer wg.Done()
	completeFileName := fmt.Sprintf("../../tmp/%s", filename)
	fmt.Printf("Uploading file: %s\n", completeFileName)
	f, err := os.Open(completeFileName)
	if err != nil {
		errorFileUpload <- filename
		// se tiver erro, esvazia o channel
		<-uploadControl
		fmt.Printf("Error opening file: %s\n", err)
		return
	}
	defer f.Close()
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		errorFileUpload <- filename
		// se tiver erro, esvazia o channel
		<-uploadControl
		fmt.Printf("Error uploading file: %s\n", err)
		return
	}
	fmt.Printf("Successfully uploaded file: %s\n", completeFileName)
	// esvazia o channel ao finalizar
	<-uploadControl
}
