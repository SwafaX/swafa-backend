package initializers

import (
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinio(config *Config) {
	endpoint := config.MinioEndpoint
	accessKeyID := config.MinioAccessKey
	secretAccessKey := config.MinioSecretKey
	useSSL := false

	// Initialize minio client object.
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln("Could not connect to MinIO")
	}

	fmt.Println("✔ Successfully connected to Minio.")
}
