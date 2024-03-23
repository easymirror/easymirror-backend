package tests

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

/*

Test to see which is more efficient (memory and speed wise)
for uploading to other hosts from AWS S3

Tips:
use io.Copy

Option 1: Download from presigned URL
Option 2: Download from download manager
*/

var s3Client *s3.Client

func init() {
	// Load the env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("no env file loaded.")
	}
	loadS3Client()
}

func loadS3Client() {
	var err error
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	s3Client = s3.NewFromConfig(cfg)
}

func getPresignURL() string {
	presignClient := s3.NewPresignClient(s3Client)
	presignedUrl, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    aws.String("50mb_file.txt"),
		},
		s3.WithPresignExpires(time.Hour*15))
	if err != nil {
		log.Fatal(err)
	}
	return presignedUrl.URL
}

func downloadFromS3Manager() (*manager.WriteAtBuffer, error) {
	var partMiBs int64 = 64
	downloader := manager.NewDownloader(s3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
	})
	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String("50mb_file.txt"),
	})
	if err != nil {
		fmt.Println("Failed to download:", err)
		return nil, err
	}
	return buffer, nil
}

// Downloads a given file from a presigned URL
func downloadFromPresigned(url string) error {

	r, err := http.Get(url)
	if err != nil {
		log.Println("Cannot get from URL", err)
	}
	defer r.Body.Close()

	file, _ := os.Create("copy.data")
	defer file.Close()

	writer := bufio.NewWriter(file)
	io.Copy(writer, r.Body)
	return writer.Flush()
}

// func uploadToPixelDrain() {
// 	// req, _ := http.NewRequest("POST", "https://pixeldrain.com/api/file/", nil)
// 	r, w := io.Pipe()
// 	m := multipart.NewWriter(w)
// 	go func() {
// 		defer w.Close()
// 		defer m.Close()
// 		part, err := m.CreateFormFile("myFile", "foo.txt")
// 		if err != nil {
// 			return
// 		}
// 		file, err := os.Open(name)
// 		if err != nil {
// 			return
// 		}
// 		defer file.Close()
// 		if _, err = io.Copy(part, file); err != nil {
// 			return
// 		}
// 	}()
// 	http.Post(url, m.FormDataContentType(), r)
// }

func TestURL(t *testing.T) {
	url := getPresignURL()
	fmt.Println(url)
}

func BenchmarkDownloadPresigned(b *testing.B) {
	b.StopTimer()
	url := getPresignURL()
	b.StartTimer()

	if err := downloadFromPresigned(url); err != nil {
		b.Fatalf("Error downloading from presigned: %v", err)
	}
}

func BenchmarkDownloadS3Manager(b *testing.B) {
	_, err := downloadFromS3Manager()
	if err != nil {
		b.Fatalf("Error downloading from s3 manager: %v", err)
	}
}
