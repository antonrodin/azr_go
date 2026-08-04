package photos

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Storage guarda y sirve los ficheros de fotos en S3. Las fotos viven bajo
// {status}/{id}/ ("private/" sin acceso público, "public/" con lectura
// pública vía bucket policy). Publicar = mover objetos de prefijo.
type Storage struct {
	client *s3.Client
	psign  *s3.PresignClient
	bucket string
	region string
}

// NewStorage crea el cliente S3 con las credenciales del entorno
// (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION del .env).
func NewStorage(ctx context.Context) (*Storage, error) {
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		return nil, fmt.Errorf("falta S3_BUCKET en el entorno")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("configuración AWS: %w", err)
	}

	client := s3.NewFromConfig(cfg)
	return &Storage{
		client: client,
		psign:  s3.NewPresignClient(client),
		bucket: bucket,
		region: cfg.Region,
	}, nil
}

// Key construye la clave S3 de un fichero de foto.
func Key(status string, id int, name string) string {
	return fmt.Sprintf("%s/%d/%s", status, id, name)
}

// PublicURL devuelve la URL pública estable de un objeto bajo public/.
func (s *Storage) PublicURL(key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, key)
}

// Presign firma una URL GET temporal para objetos privados.
func (s *Storage) Presign(ctx context.Context, key string, dur time.Duration) (string, error) {
	req, err := s.psign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(dur))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// UploadFile sube un fichero local. Cache-Control largo: las claves llevan
// el id y los derivados no cambian nunca.
func (s *Storage) UploadFile(ctx context.Context, key, path, contentType string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(s.bucket),
		Key:          aws.String(key),
		Body:         f,
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	})
	return err
}

// listKeys devuelve las claves bajo un prefijo.
func (s *Storage) listKeys(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	p := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(prefix),
	})
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			keys = append(keys, *obj.Key)
		}
	}
	return keys, nil
}

// MovePhoto mueve todos los ficheros de una foto entre prefijos
// (private→public al publicar, public→private al despublicar).
// Copia primero y borra después: si algo falla a medias, no se pierde nada.
func (s *Storage) MovePhoto(ctx context.Context, id int, fromStatus, toStatus string) error {
	fromPrefix := fmt.Sprintf("%s/%d/", fromStatus, id)
	keys, err := s.listKeys(ctx, fromPrefix)
	if err != nil {
		return err
	}

	for _, key := range keys {
		newKey := fmt.Sprintf("%s/%d/%s", toStatus, id, key[len(fromPrefix):])
		_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(s.bucket),
			CopySource: aws.String(url.PathEscape(s.bucket + "/" + key)),
			Key:        aws.String(newKey),
		})
		if err != nil {
			return fmt.Errorf("copiar %s: %w", key, err)
		}
	}

	return s.deleteKeys(ctx, keys)
}

// DeletePhoto borra todos los ficheros de una foto de su prefijo actual.
func (s *Storage) DeletePhoto(ctx context.Context, id int, status string) error {
	keys, err := s.listKeys(ctx, fmt.Sprintf("%s/%d/", status, id))
	if err != nil {
		return err
	}
	return s.deleteKeys(ctx, keys)
}

func (s *Storage) deleteKeys(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}
	objs := make([]types.ObjectIdentifier, len(keys))
	for i, k := range keys {
		objs[i] = types.ObjectIdentifier{Key: aws.String(k)}
	}
	_, err := s.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.bucket),
		Delete: &types.Delete{Objects: objs},
	})
	return err
}
