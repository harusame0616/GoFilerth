package tui

import (
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/harusame0616/GoFilerth/gofilerth/infrastructure/local"
	"github.com/harusame0616/GoFilerth/gofilerth/usecase"
	"github.com/rivo/tview"
)

type pathChangeObserver func(newPath string)

type filesView struct {
	table       *tview.Table
	fileQuery   *usecase.FileQueryUsecase
	currentPath string
	pathChangeObserver
}

// ファイル一覧ビューを作成する
// path : 初期表示パス
func NewFilesView(path string) *filesView {
	fileQuery, err := usecase.NewFileQuery(local.NewFileQuery())
	if err != nil {
		log.Fatal(err)
	}

	table := tview.NewTable()
	table.SetSelectable(true, false)

	// 初期パス表示
	fv := &filesView{table: table, fileQuery: fileQuery}
	fv.open(path)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'o':
				row, _ := table.GetSelection()
				selectedCell := table.GetCell(row, 0)
				path := fv.currentPath
				if !strings.HasSuffix(path, "/") {
					path += "/"
				}

				fv.open(path + selectedCell.Text)

				return nil
			case 'O':
				fv.upDirectory()
				return nil
			}
		}

		return event
	})

	return fv
}

// パスを開く
func (fv *filesView) open(path string) error {
	files, err := fv.fileQuery.ListFiles(path)
	if err != nil {
		return err
	}

	fv.currentPath = path
	fv.updateTable(files)
	if fv.pathChangeObserver != nil {
		fv.pathChangeObserver(path)
	}

	return nil
}

// カレントパスから一階層上に上がる
func (fv *filesView) upDirectory() error {
	paths := strings.Split(fv.currentPath, "/")

	upperPath := strings.Join(paths[0:len(paths)-1], "/")
	if upperPath == "" {
		upperPath = "/"
	}

	return fv.open(upperPath)
}

// テーブルを更新する
func (fv *filesView) updateTable(files []usecase.FileDto) {
	fv.table.Clear()
	for row, file := range files {
		fileNameCell := tview.NewTableCell(file.Name)
		fileNameCell.SetExpansion(1)
		fv.table.SetCell(row, 0, fileNameCell)
	}
	fv.table.ScrollToBeginning()
}

// currentPathが変更されたときに実行されるオブザーバーを登録する
// 複数回呼ばれると最後に登録されたオブザーバーのみ有効
func (fv *filesView) observePathChange(observer pathChangeObserver) {
	fv.pathChangeObserver = observer
}
