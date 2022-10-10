package gcs

import (
	"context"
	"regexp"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/harusame0616/GoFilerth/gofilerth/usecase"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FileQuery struct {
	client *storage.Client
	ctx    context.Context
}

func NewFileQueryByCredential(credentialFilePath string) (*FileQuery, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))

	if err != nil {
		return nil, err
	}

	return &FileQuery{client, ctx}, nil
}

func (fileQuery *FileQuery) ListFiles(path string) ([]usecase.FileDto, error) {
	pathFormat := regexp.MustCompile(`^gcs://(.+?)(/(.*))*$`)
	matches := pathFormat.FindAllStringSubmatch(path, -1)

	bucketName, path := matches[0][1], matches[0][3]
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// バケットルートを取得したいときに / を指定してしまうと
	// オブジェクト一覧が取得できない
	if path == "/" {
		path = ""
	}

	var files []usecase.FileDto
	for it := fileQuery.client.Bucket(bucketName).Objects(fileQuery.ctx, &storage.Query{Prefix: path, Delimiter: "/"}); ; {
		attrs, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		file := usecase.FileDto{
			Size:       attrs.Size,
			Path:       path,
			ModifiedAt: attrs.Updated,
		}

		if len(attrs.Prefix) > 0 {
			paths := strings.Split(attrs.Prefix, "/")
			file.Name = paths[len(paths)-2] // 最後が空白要素になるので -2
			file.IsDirectory = true
		} else if len(attrs.Name) > 0 {
			paths := strings.Split(attrs.Name, "/")
			file.Name = paths[len(paths)-1]
			file.IsDirectory = false
		}

		files = append(files, file)
	}

	return files, nil
}
