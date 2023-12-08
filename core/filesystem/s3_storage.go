// Package filesystem
// Reference https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html
package filesystem

import (
	"app/core/errors"
	"app/core/log"
	"app/core/utils"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"os"
	"path/filepath"
	"time"
)

//---------------------------------------
//				Structure				|
//---------------------------------------

// NewS3Storage Create S3 Storage with basics info.
func NewS3Storage() *S3Storage {
	// Load the Shared AWS Configuration (~/.aws/config). Note: Also load combine .env file.
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	return &S3Storage{
		S3Client: s3.NewFromConfig(cfg),
	}
}

/*
Make sure S3 below setting:

*** Bucket policy ***
[code]
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicRead",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:GetObject",
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:PutObjectAcl",
                "s3:GetObjectAcl",
				"s3:GetObjectAttributes"
            ],
            "Resource": "arn:aws:s3:::902-local/*"
        }
    ]
}
[/code]

*** Cross-origin resource sharing (CORS) ***
[code]
[
    {
        "AllowedHeaders": [
            "*"
        ],
        "AllowedMethods": [
            "PUT",
            "POST",
            "DELETE"
        ],
        "AllowedOrigins": [
            "*"
        ],
        "ExposeHeaders": []
    }
]
[/code]
*/

type S3Storage struct {
	S3Client *s3.Client
}

// ---------------------------------------
//			Implement IStorage			|
// ---------------------------------------

func (s *S3Storage) Put(path, contents string, options ...interface{}) bool {
	localStorage := NewLocalStorage()

	// Put content to temporary dir at local.
	tempPath := fmt.Sprintf("%s/%s", os.Getenv("TEMP_DIR"), utils.FileName(path))
	localStorage.Put(tempPath, contents)

	// Open file source
	file, err := os.Open(filepath.Clean(tempPath))
	if err != nil {
		log.Errorf("Unable create file %q. Here's why: %v\n", tempPath, err)

		return false
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Errorf("Unable to close file %q. Here's why: %v\n", tempPath, err)
		}
	}(file)

	return s.PutFile(path, file, options)
}

func (s *S3Storage) PutFile(path string, fileSource *os.File, options ...interface{}) bool {
	_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(path),
		Body:   fileSource,
	})
	if err != nil {
		log.Errorf("Unable to write file %q. Here's why: %v\n", path, err)

		return false
	}

	return true
}

func (s *S3Storage) PutFilepath(path, filePath string, options ...interface{}) bool {
	fileSource, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		log.Errorf("Unable to read file %q. Here's why: %v\n", filePath, err)

		return false
	}

	_, err = s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(path),
		Body:   fileSource,
	})
	if err != nil {
		log.Errorf("Unable to write file %q. Here's why: %v\n", path, err)

		return false
	}

	return true
}

func (s *S3Storage) Delete(path string) bool {
	_, err := s.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Errorf("Unable to delete file %q. Here's why: %v\n", path, err)

		return false
	}

	return true
}

func (s *S3Storage) Copy(from, to string) bool {
	_, err := s.S3Client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(os.Getenv("AWS_BUCKET")),
		CopySource: aws.String(fmt.Sprintf("%s/%s", os.Getenv("AWS_BUCKET"), from)),
		Key:        aws.String(to),
	})
	if err != nil {
		log.Errorf("Unable to copy file %s to %s. Here's why: %v\n", from, to, err)

		return false
	}

	return true
}

func (s *S3Storage) Move(from, to string) bool {
	if s.Copy(from, to) {
		return s.Delete(from)
	}

	return false
}

func (s *S3Storage) Exists(path string) bool {
	return s.Size(path) != 0
}
func (s *S3Storage) Get(path string) ([]byte, error) {
	result, err := s.getObject(path)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Unable to close body. Here's why: %v\n", err)
		}
	}(result.Body)

	body, err := io.ReadAll(result.Body)
	if err != nil {
		log.Errorf("Unable read object body from %v. Here's why: %v\n", path, err)
	}

	return body, nil
}

func (s *S3Storage) Size(path string) int64 {
	result, err := s.S3Client.GetObjectAttributes(context.TODO(), &s3.GetObjectAttributesInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(path),
		ObjectAttributes: []s3Types.ObjectAttributes{
			s3Types.ObjectAttributesObjectSize,
		},
	})
	if err != nil {
		return 0
	}

	return *result.ObjectSize
}

func (s *S3Storage) LastModified(path string) time.Time {
	result, err := s.S3Client.GetObjectAttributes(context.TODO(), &s3.GetObjectAttributesInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(path),
		ObjectAttributes: []s3Types.ObjectAttributes{
			s3Types.ObjectAttributesObjectParts,
		},
	})

	if err != nil {
		log.Errorf("Unable to get info of %s. Here's why: %v\n", path, err)

		return time.Time{}
	}

	return *result.LastModified
}

func (s *S3Storage) MakeDir(dir string) bool {
	return s.Put(fmt.Sprintf("%s/.info", dir), "Info")
}

func (s *S3Storage) DeleteDir(dir string) bool {
	// Get all objects in dir
	// Note: Can not delete a dir have children object.
	result, err := s.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Prefix: aws.String(dir),
	})
	var contents []s3Types.Object
	if err != nil {
		log.Errorf("Unable to list objects from dir %v. Here's why: %v\n", dir, err)
	} else {
		contents = result.Contents
	}

	var objectIds []s3Types.ObjectIdentifier

	// Collect children objects
	for _, object := range contents {
		objectIds = append(objectIds, s3Types.ObjectIdentifier{Key: object.Key})
	}

	// Append current object
	objectIds = append(objectIds, s3Types.ObjectIdentifier{Key: aws.String(dir)})

	// Delete objects
	_, err = s.S3Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Delete: &s3Types.Delete{Objects: objectIds},
	})

	if err != nil {
		log.Errorf("Unable to delete object from bucket %v. Here's why: %v\n", dir, err)

		return false
	}

	return true
}

func (s *S3Storage) Append(path, data string) bool {
	log.Errorf("Unable to append data %s into %s. Here's why: %v\n", path, data, errors.NotYetImplemented.Error())

	return false
}

func (s *S3Storage) getObject(path string) (*s3.GetObjectOutput, error) {
	result, err := s.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Errorf("Unable to get object %s. Here's why: %v\n", path, err)

		return nil, err
	}

	return result, nil
}
