package worker

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	AwsAccessKeyID     = ""
	AwsSecretAccessKey = ""
	Bucket             = "key-master"
	Region             = "ap-southeast-1"
	Folder             = "keys"
)

type AwsMaster struct {
	Session   *session.Session
	S3        *s3.S3
	EmptyFile *os.File
	Count     int64
	WorkerID  int
}

func NewAwsMaster(workerID int) *AwsMaster {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewStaticCredentials(AwsAccessKeyID, AwsSecretAccessKey, ""),
	})
	if err != nil {
		fmt.Println("Error creating AWS session:", err)
		return nil
	}
	file, err := os.Create("./temp")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}
	return &AwsMaster{
		Session:   sess,
		S3:        s3.New(sess),
		EmptyFile: file,
		Count:     0,
		WorkerID:  workerID,
	}
}

func (a AwsMaster) Upload(privateKey, address string) {
	p := 2
	set := 3
	address = strings.ToLower(address)
	addressSplits := make([]string, 0)
	for ; p+set <= len(address)-2; p = p + set {
		addressSplits = append(addressSplits, address[p:p+set])
	}
	addressSplits = append(addressSplits, address[p:])
	addressPath := strings.Join(addressSplits, "/")
	filePath := Folder + "/" + addressPath + "/"
	fileKey := filePath + privateKey
	a.UploadFile2S3(fileKey)
}

// CreateLocalFile costs a lot , unless you don't care your local disk
func (a AwsMaster) CreateLocalFile(dir, file string) {
	err := os.MkdirAll(dir, fs.FileMode(0755))
	if err != nil {
		return
	}
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return
	}
}

func (a AwsMaster) UploadFile2S3(fileKey string) {
	_, err := a.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(fileKey),
		Body:   a.EmptyFile,
	})
	if err != nil {
		fmt.Println("Error uploading file to S3:", err)
		return
	}
	a.Count++
}

func (a AwsMaster) Stop() {
	fmt.Println("S3 uploaded :", a.Count)
	_ = a.EmptyFile.Close()
}
