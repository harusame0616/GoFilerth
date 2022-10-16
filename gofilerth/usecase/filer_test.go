package usecase

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/harusame0616/GoFilerth/gofilerth/domain/file"
	"github.com/harusame0616/GoFilerth/gofilerth/infrastructure/inmemory"
)

type fileRepositoryMock struct {
	files map[string]file.Dto
}

func (fileRepository *fileRepositoryMock) GetOneByPath(destinationPath string) (*file.File, error) {
	if _file, ok := fileRepository.files[destinationPath]; ok {
		return file.FromDto(_file), nil
	} else {
		return nil, errors.New("file is not found")
	}
}

func TestCreateNewFiler(t *testing.T) {

	fileRepository := &fileRepositoryMock{files: map[string]file.Dto{
		"/":                  {Path: "/", IsDirectory: true},
		"/dir1":              {Path: "/dir1", IsDirectory: true},
		"/dir_nest/dir_nest": {Path: "/dir_nest/dir_nest", IsDirectory: true},
		"/dir_nest/file":     {Path: "/dir_nest/file", IsDirectory: false},
	}}
	filerCommand := NewFilerCommand(fileRepository, inmemory.NewFilerRepository())

	t.Run("正常系", func(t *testing.T) {
		for _, path := range []string{"/", "/dir1", "/dir_nest/dir_nest"} {
			t.Run("ディレクトリが存在する", func(t *testing.T) {
				if id, err := filerCommand.CreateNewFiler(path); err != nil {
					t.Errorf("err = %v; want err == nil", err)
				} else {
					currentPath := filerCommand.CurrentPath(id)

					expectPath := path
					if strings.HasSuffix(path, "/") && len(path) > 1 {
						expectPath = path[:len(path)-1]
					}

					if expectPath != currentPath {
						t.Errorf("CurrentPath = %s; want %s", currentPath, expectPath)
					}
				}
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		t.Run("パスが存在しない", func(t *testing.T) {
			if id, err := filerCommand.CreateNewFiler("/not_exists"); id != "" || err == nil {
				t.Errorf("id == %s, err == %v; want id == \"\", err != nil", id, err)
			}
		})

		t.Run("パスがディレクトリではない", func(t *testing.T) {
			if id, err := filerCommand.CreateNewFiler("/dir_nest/file"); id != "" || err == nil {
				t.Errorf("id == %s, err == %v; want id == \"\", err != nil", id, err)
			}
		})

		t.Run("パスがカラ文字", func(t *testing.T) {
			if id, err := filerCommand.CreateNewFiler(""); id != "" || err == nil {
				t.Errorf("id == %s, err == %v; want id == \"\", err != nil", id, err)
			}
		})
	})
}

func TestChangeDirectory(t *testing.T) {

	fileRepository := &fileRepositoryMock{files: map[string]file.Dto{
		"/":                  {Path: "/", IsDirectory: true},
		"/dir1":              {Path: "/dir1", IsDirectory: true},
		"/dir_nest/dir_nest": {Path: "/dir_nest/dir_nest", IsDirectory: true},
		"/dir_nest/file":     {Path: "/dir_nest/file", IsDirectory: false},
	}}
	filerCommand := NewFilerCommand(fileRepository, inmemory.NewFilerRepository())
	id, err := filerCommand.CreateNewFiler("/")
	if err != nil {
		t.Errorf("prepare error")
	}

	t.Run("正常系", func(t *testing.T) {
		for _, path := range []string{"/", "/dir1", "/dir_nest/dir_nest"} {
			t.Run("ディレクトリが存在する", func(t *testing.T) {
				if err := filerCommand.ChangeDirectory(id, path); err != nil {
					t.Errorf("err = %v; want err == nil", err)
				}

				currentPath := filerCommand.CurrentPath(id)
				if path != currentPath {
					t.Errorf("CurrentPath = (%s); want %s", currentPath, path)
				}
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		t.Run("パスが存在しない", func(t *testing.T) {
			if err := filerCommand.ChangeDirectory(id, "/not_exists"); err == nil {
				t.Errorf("err == %v; want id == \"\", err != nil", err)
			}
		})

		t.Run("パスがディレクトリではない", func(t *testing.T) {
			if err := filerCommand.ChangeDirectory(id, "/dir_nest/file"); err == nil {
				t.Errorf("err == %v; want id == \"\", err != nil", err)
			}
		})

		t.Run("パスがカラ文字", func(t *testing.T) {
			if err := filerCommand.ChangeDirectory(id, ""); err == nil {
				t.Errorf("err == %v; want id == \"\", err != nil", err)
			}
		})
	})
}

func TestUpDirectory(t *testing.T) {
	fileRepository := &fileRepositoryMock{files: map[string]file.Dto{
		"/":          {Path: "/", IsDirectory: true},
		"/dir1":      {Path: "/dir1", IsDirectory: true},
		"/dir1/dir1": {Path: "/dir2", IsDirectory: true},
		"/dir2/dir2": {Path: "/dir2/dir2", IsDirectory: true},
	}}
	filerCommand := NewFilerCommand(fileRepository, inmemory.NewFilerRepository())

	t.Run("正常系", func(t *testing.T) {
		for _, testParam := range [][]string{{"/dir1", "/"}, {"/dir1/dir1", "/dir1"}} {
			testPath, expectPath := testParam[0], testParam[1]

			t.Run(fmt.Sprintf("上位ディレクトリが存在する(test:%s,expect:%s)", testPath, expectPath), func(t *testing.T) {
				id, err := filerCommand.CreateNewFiler(testPath)
				if err != nil {
					t.Errorf("err == %v; want err == nil", err)
				}

				if path, err := filerCommand.UpDirectory(id); path == expectPath && err == nil {
					// pass
				} else {
					t.Errorf("path = %s, err = %v; want path == %s, err == nil", filerCommand.CurrentPath((id)), err, expectPath)
				}
			})
		}
	})

	t.Run("異常系", func(t *testing.T) {
		t.Run("カレントパスがルート", func(t *testing.T) {
			id, err := filerCommand.CreateNewFiler("/")
			if err != nil {
				t.Errorf("prepare error, %v", err)
			}

			if path, err := filerCommand.UpDirectory(id); err == nil || path != "" {
				t.Errorf("path == %s, err == %v; want path == \"\", err != nil", path, err)
			}
		})

		t.Run("上位パスが存在しない", func(t *testing.T) {
			id, err := filerCommand.CreateNewFiler("/dir2/dir2")
			if err != nil {
				t.Errorf("err == %v; err == nil", err)
			}

			if path, err := filerCommand.UpDirectory(id); err == nil || path != "" {
				t.Errorf("path == %s, err == %v; want path == \"\", err != nil", path, err)
			}
		})

	})
}
