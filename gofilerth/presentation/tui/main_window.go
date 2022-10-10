package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type mainWindow struct {
	app       *tview.Application
	leftPain  *filesView
	rightPain *filesView
}

func NewMainWindow() *mainWindow {
	window := &mainWindow{app: tview.NewApplication(), leftPain: NewFilesView("/"), rightPain: NewFilesView("/")}

	// ファイル一覧横２画面
	fileListViews := tview.NewFlex()
	fileListViews.AddItem(window.leftPain.table, 0, 1, false)
	fileListViews.AddItem(window.rightPain.table, 0, 1, false)
	window.leftPain.table.SetBorder(true)
	window.rightPain.table.SetBorder(true)

	// パス横２画面
	pathViews := tview.NewFlex()
	path1 := tview.NewTextArea()
	path2 := tview.NewTextArea()
	path1.SetText("/", false)
	path2.SetText("/", false)
	path1.SetBorder(true)
	path2.SetBorder(true)
	window.leftPain.observePathChange(func(newPath string) {
		path1.SetText(newPath, false)
	})
	window.rightPain.observePathChange(func(newPath string) {
		path2.SetText(newPath, false)
	})

	pathViews.AddItem(path1, 0, 1, false)
	pathViews.AddItem(path2, 0, 1, false)

	root := tview.NewFlex()
	root.SetDirection(tview.FlexRow)
	root.AddItem(pathViews, 3, 1, false)
	root.AddItem(fileListViews, 0, 1, false)

	window.app.SetRoot(root, true).EnableMouse(true)
	window.app.SetFocus(window.leftPain.table)

	window.app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyCtrlJ:
				// CtrlFを押した時の処理を記述
				return event // CtrlFをInputFieldのdefaultのキー設定へ伝える

			case tcell.KeyRune:
				switch event.Rune() {
				case 'h':
					window.app.SetFocus(window.leftPain.table)
				case 'l':
					window.app.SetFocus(window.rightPain.table)
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
