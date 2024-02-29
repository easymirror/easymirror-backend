package tests

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func uploadWithPipe(url string, file *os.File) (resp *http.Response, err error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() error {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("myFile", "foo.txt")
		if err != nil {
			return err
		}
		// file, err := os.Open(name)
		// if err != nil {
		// 	return err
		// }
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return err
		}
		return nil
	}()

	req, _ := http.NewRequest("POST", url, r)
	req.Header = http.Header{
		"Content-Type":  {m.FormDataContentType()},
		"Authorization": {"Basic OjlkZmM4ZmJhLTY2YjctNGMwNy04MzU1LWI1ZGNmMTA3ZGJiYw=="},
	}
	c := &http.Client{}
	return c.Do(req)
}

func uploadWithBuf(url string, file *os.File) (resp *http.Response, err error) {
	// // Open the file from AWS S3
	// fileResp, err := http.Get(s3FileURL)
	// if err != nil {
	// 	panic(err)
	// }
	// defer fileResp.Body.Close()

	// // Create a buffer to store the file
	var buf bytes.Buffer
	// io.Copy(&buf, fileResp.Body)
	io.Copy(&buf, file)

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the form data
	part, err := writer.CreateFormFile("file", "filename.jpg")
	if err != nil {
		panic(err)
	}
	io.Copy(part, &buf)

	// Close the multipart writer
	writer.Close()

	req, _ := http.NewRequest("POST", url, body)
	req.Header = http.Header{
		"Content-Type":  {writer.FormDataContentType()},
		"Authorization": {"Basic OjlkZmM4ZmJhLTY2YjctNGMwNy04MzU1LWI1ZGNmMTA3ZGJiYw=="},
	}
	c := &http.Client{}
	return c.Do(req)
}

func createFile(sizeMB int) *os.File {
	// Define the file size in bytes
	fileSize := int64(sizeMB * 1024 * 1024) // 50MB

	// Create a new file
	file, err := os.Create("50mb_file.txt")
	if err != nil {
		log.Fatal(err)
	}
	// defer file.Close()

	// Set the file size
	if err := file.Truncate(fileSize); err != nil {
		log.Fatal(err)
	}
	return file
}

func BenchmarkUploadWithPipe(b *testing.B) {
	// Create file
	b.StopTimer()
	file := createFile(50)
	b.StartTimer()

	// Upload
	_, err := uploadWithPipe("https://pixeldrain.com/api/file/", file)
	if err != nil {
		b.Fatalf("Error uploading with pipe: %v", err)
	}

	// Close file
	b.StopTimer()
	file.Close()
	b.StartTimer()
}

func BenchmarkUploadWithBuf(b *testing.B) {
	// Create file
	b.StopTimer()
	file := createFile(50)
	b.StartTimer()

	// Upload
	_, err := uploadWithBuf("https://pixeldrain.com/api/file/", file)
	if err != nil {
		b.Fatalf("Error uploading with pipe: %v", err)
	}

	// Close file
	b.StopTimer()
	file.Close()
	b.StartTimer()
}
