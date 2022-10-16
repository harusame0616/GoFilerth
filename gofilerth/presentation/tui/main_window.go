package tui

import (
	"os"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	FOCUS_LEFT_PANE  = 1
	FOCUS_RIGHT_PANE = 2
)

type mainWindow struct {
	app       *tview.Application
	leftPane  *filesView
	rightPane *filesView
	focusArea int
}

func NewMainWindow() *mainWindow {
	window := &mainWindow{app: tview.NewApplication(), leftPane: NewFilesView("/"), rightPane: NewFilesView("/")}

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
	window.focusArea = FOCUS_LEFT_PANE

	window.app.SetInputCapture(
		func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyRune:
				switch event.Rune() {
				case 'h':
					window.app.SetFocus(window.leftPane.table)
					window.focusArea = FOCUS_LEFT_PANE
				case 'l':
					window.app.SetFocus(window.rightPane.table)
					window.focusArea = FOCUS_RIGHT_PANE
				case 'S':
					window.openShell()
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

func (window *mainWindow) openShell() {
	window.app.Suspend(func() {
		shell := exec.Command("bash")
		shell.Stdin = os.Stdin
		shell.Stdout = os.Stdout
		shell.Stderr = os.Stderr

		if window.focusArea == FOCUS_LEFT_PANE {
			shell.Dir = window.leftPane.CurrentPath()
		} else {
			shell.Dir = window.rightPane.CurrentPath()
		}
		shell.Run()
	})
}
