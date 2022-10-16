package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type mainWindow struct {
	app       *tview.Application
	leftPane  *filesView
	rightPane *filesView
}

func NewMainWindow() *mainWindow {
	app := tview.NewApplication()
	window := &mainWindow{app: app, leftPane: NewFilesView("/", app), rightPane: NewFilesView("/", app)}

	// ファイル一覧横２画面
	fileListViews := tview.NewFlex()
	fileListViews.AddItem(window.leftPane.table, 0, 1, false)
	fileListViews.AddItem(window.rightPane.table, 0, 1, false)
	window.leftPane.table.SetBorder(true)
	window.rightPane.table.SetBorder(true)

	// パス横２画面
	pathViews := tview.NewFlex()
	path1 := tview.NewTextArea()
	path2 := tview.NewTextArea()
	path1.SetText("/", false)
	path2.SetText("/", false)
	path1.SetBorder(true)
	path2.SetBorder(true)
	window.leftPane.observePathChange(func(newPath string) {
		path1.SetText(newPath, false)
	})
	window.rightPane.observePathChange(func(newPath string) {
		path2.SetText(newPath, false)
	})

	pathViews.AddItem(path1, 0, 1, false)
	pathViews.AddItem(path2, 0, 1, false)

	root := tview.NewFlex()
	root.SetDirection(tview.FlexRow)
	root.AddItem(pathViews, 3, 1, false)
	root.AddItem(fileListViews, 0, 1, false)

	window.app.SetRoot(root, true).EnableMouse(true)
	window.app.SetFocus(window.leftPane.table)

	window.app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyRune:
				switch event.Rune() {
				case 'h':
					window.app.SetFocus(window.leftPane.table)
				case 'l':
					window.app.SetFocus(window.rightPane.table)
				}

			}
			return event // 上記以外のキー入力をdefaultのキーアクションへ伝える
		})

	return window
}

func (window *mainWindow) Run() {
	if err := window.app.Run(); err != nil {
		panic(err)
	}
}
