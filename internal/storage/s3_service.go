package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	_ "github.com/joho/godotenv/autoload"
)

type S3Service struct {
	client   *s3.Client
	bucket   string
	basePath string
}

func NewS3Service() (*S3Service, error) {
	if os.Getenv("WATCH_PATH") == "" {
		return nil, errors.New("S3Service: 未設定環境變數WATCH_PATH")
	}

	if os.Getenv("AWS_REGION") == "" {
		return nil, errors.New("S3Service: 未設定環境變數AWS_REGION")
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		return nil, errors.New("S3Service: 未設定環境變數AWS_ACCESS_KEY_ID")
	}

	if os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		return nil, errors.New("S3Service: 未設定環境變數AWS_SECRET_ACCESS_KEY")
	}

	if os.Getenv("S3_BUCKET_NAME") == "" {
		return nil, errors.New("S3Service: 未設定環境變數S3_BUCKET_NAME")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client:   client,
		bucket:   os.Getenv("S3_BUCKET_NAME"),
		basePath: os.Getenv("WATCH_PATH"),
	}, nil
}

func (m *S3Service) GetFileList() ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: &m.bucket,
	}

	resp, err := m.client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	fileList := []string{}
	for _, item := range resp.Contents {
		fileList = append(fileList, *item.Key)
	}

	log.Printf("S3Service: 成功從S3擷取檔案清單")
	return fileList, nil
}

func (m *S3Service) CheckFileIsExist(fileName string) (bool, error) {
	input := &s3.HeadObjectInput{
		Bucket: &m.bucket,
		Key:    &fileName,
	}

	_, err := m.client.HeadObject(context.TODO(), input)
	if err == nil {
		log.Printf("S3Service: 檔案 %s 已存在於S3上", fileName)
		return true, nil
	}

	if strings.Index(err.Error(), "404") != -1 {
		return false, nil
	} else {
		return false, err
	}
}

func (m *S3Service) PutFile(fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		return errors.New(fmt.Sprintf("S3Service: 開啟本地檔案 %s 時發生錯誤", fileName))
	}

	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket:       &m.bucket,
		Key:          &fileName,
		Body:         file,
		StorageClass: types.StorageClassStandard,
	}

	_, err = m.client.PutObject(context.TODO(), input)
	if err != nil {
		return err
	}

	log.Printf("S3Service: 檔案 %s 成功儲存至S3", fileName)
	return nil
}
