package tests

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
)

func uploadWithPipe(uri string, file *os.File) (resp *http.Response, err error) {
	// Copy into buffer.
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("file", file.Name())
		if err != nil {
			return
		}
		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()

	// Upload
	req, _ := http.NewRequest("POST", uri, r)
	req.Header = http.Header{
		"Content-Type":  {m.FormDataContentType()},
		"Authorization": {"Basic OjlkZmM4ZmJhLTY2YjctNGMwNy04MzU1LWI1ZGNmMTA3ZGJiYw=="},
	}
	client := &http.Client{}
	// proxyURL, _ := url.Parse("http://10.0.0.58:8888")
	// proxy := http.ProxyURL(proxyURL)
	// transport := &http.Transport{Proxy: proxy}
	// client := &http.Client{Transport: transport}
	return client.Do(req)
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
	part, err := writer.CreateFormFile("file", file.Name())
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

func createFile(sizeMB int, name string) *os.File {
	// Define the file size in bytes
	fileSize := int64(sizeMB * 1024 * 1024) // 50MB

	// Create a new file
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file.WriteString("Hello")

	// Set the file size
	if err := file.Truncate(fileSize); err != nil {
		log.Fatal(err)
	}
	return file
}

func BenchmarkUploadWithPipe(b *testing.B) {
	// Create file
	b.StopTimer()
	file := createFile(10, "uploaded_with_pipe.txt")
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
	file := createFile(10, "uploaded_with_buf.txt")
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

func BenchmarkRealWorldPipe(b *testing.B) {
	// Get Body
	URL := "https://easymirror.s3.us-east-1.amazonaws.com/43bb847b-d89c-4db4-803a-522a0148b8d0/Matthew_Lugo_Resume.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA4LVGZIJF45Z3KRFR%2F20240303%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240303T025927Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=31720faf643cd52ef0385574ec1563615bc1f7213c18777bd14e6c259582f704"
	// URL := "https://easymirror.s3.us-east-1.amazonaws.com/43bb847b-d89c-4db4-803a-522a0148b8d0/MRL_0053.MOV?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA4LVGZIJF45Z3KRFR%2F20240303%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240303T022137Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=881482ad26ab8295686e526fc1b7ab2f495e1b84812d1a63bb610c157f1ed7d9"

	req, _ := http.NewRequest("GET", URL, nil)
	resp1, err := http.DefaultClient.Do(req)
	if err != nil {
		b.Fatalf("Error getting body: %v", err)
	}

	// Copy into buffer.
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		defer resp1.Body.Close()
		part, err := m.CreateFormFile("file", "Matthew_Lugo_Resume.pdf")
		if err != nil {
			return
		}
		if _, err = io.Copy(part, resp1.Body); err != nil {
			return
		}
	}()
	// http.Post(url, m.FormDataContentType(), r)

	// Upload
	// b.StopTimer()
	// server := getServer()
	// uploadURL := fmt.Sprintf("https://%v.gofile.io/uploadFile", server)
	uploadURL := "https://pixeldrain.com/api/file"
	// b.StartTimer()

	req, _ = http.NewRequest("POST", uploadURL, r)
	req.Header = http.Header{
		"Content-Type":  {m.FormDataContentType()},
		"Authorization": {"Basic OjlkZmM4ZmJhLTY2YjctNGMwNy04MzU1LWI1ZGNmMTA3ZGJiYw=="},
	}

	// proxyURL, _ := url.Parse("http://10.0.0.58:8888")
	// proxy := http.ProxyURL(proxyURL)
	// transport := &http.Transport{Proxy: proxy}
	// client := &http.Client{Transport: transport}

	client := &http.Client{}
	resp2, err := client.Do(req)
	if err != nil {
		b.Fatalf("Error uploading: %v", err)
	}
	body, _ := io.ReadAll(resp2.Body)
	fmt.Println(string(body))
}
