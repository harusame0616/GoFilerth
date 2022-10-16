package tui

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/harusame0616/GoFilerth/gofilerth/infrastructure/inmemory"
	"github.com/harusame0616/GoFilerth/gofilerth/infrastructure/local"
	"github.com/harusame0616/GoFilerth/gofilerth/infrastructure/shell"
	"github.com/harusame0616/GoFilerth/gofilerth/presentation/setting"
	"github.com/harusame0616/GoFilerth/gofilerth/usecase"
	"github.com/rivo/tview"
)

type pathChangeObserver func(newPath string)

type filesView struct {
	app          *tview.Application
	table        *tview.Table
	fileQuery    *usecase.FileQueryUsecase
	filerUsecase *usecase.FilerCommand
	filerId      string
	pathChangeObserver
	files []usecase.FileDto
}

// ファイル一覧ビューを作成する
// path : 初期表示パス
func NewFilesView(path string, app *tview.Application) *filesView {
	fv := &filesView{app: app}

	if fileQuery, err := usecase.NewFileQuery(local.NewFileQuery()); err == nil {
		fv.fileQuery = fileQuery
	} else {
		log.Fatal(err)
	}

	fv.filerUsecase = usecase.NewFilerCommand(local.NewFileRepository(), inmemory.NewFilerRepository(), shell.NewCommandGateway())

	fv.table = tview.NewTable()
	fv.table.SetSelectable(true, false)

	// 初期パス表示
	if filerId, err := fv.filerUsecase.CreateNewFiler(path); err == nil {
		fv.filerId = filerId
	} else {
		log.Fatal(err)
	}
	fv.openPath(path)

	fv.table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'o':
				row, _ := fv.table.GetSelection()
				fv.openByIndex(row)
				return nil
			case 'O':
				fv.upDirectory()
				return nil
			case 'S':
				fv.openShell()
				return nil
			}
		}

		return event
	})

	return fv
}

// パスを開く
func (fv *filesView) openPath(path string) error {
	if err := fv.filerUsecase.ChangeDirectory(fv.filerId, path); err != nil {
		log.Fatal(err)
	}

	fv.updatePath(path)

	return nil
}

func (fv *filesView) openByIndex(index int) {
	fv.openPath(fv.files[index].FullPath)
}

// カレントパスから一階層上に上がる
func (fv *filesView) upDirectory() error {
	if path, err := fv.filerUsecase.UpDirectory(fv.filerId); err == nil {
		fv.updatePath(path)
		return nil
	} else {
		return err
	}
}

// テーブルを更新する
func (fv *filesView) updatePath(path string) {
	if files, err := fv.fileQuery.ListFiles(path); err == nil {
		fv.files = files
	} else {
		log.Fatal(err)
	}

	if fv.pathChangeObserver != nil {
		fv.pathChangeObserver(path)
	}

	fv.table.Clear()
	for row, file := range fv.files {
		fileNameCell := tview.NewTableCell(file.Name)
		fileNameCell.SetExpansion(1)
		fv.table.SetCell(row, 0, fileNameCell)
		if color, err := setting.GetColor(file); err == nil {
			fileNameCell.SetTextColor(tcell.NewHexColor(color))
		} else {
			log.Fatal(err)
		}
	}

	fv.table.ScrollToBeginning()
	fv.table.Select(0, 0)
}

func (fv *filesView) openShell() {
	fv.app.Suspend(func() {
		fv.filerUsecase.OpenShell(fv.filerId)
	})
}

// currentPathが変更されたときに実行されるオブザーバーを登録する
// 複数回呼ばれると最後に登録されたオブザーバーのみ有効
func (fv *filesView) observePathChange(observer pathChangeObserver) {
	fv.pathChangeObserver = observer
}
