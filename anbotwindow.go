// anbotwindow.go
package main

import (
	"github.com/therecipe/qt/widgets"
)

type AboutWindow struct {
	*widgets.QDialog
}

func NewAboutWindow(parent widgets.QWidget_ITF) (aboutWindow *AboutWindow) {
	aboutWindow = &AboutWindow{}
	aboutWindow.QDialog = widgets.NewQDialog(parent, 0)
	aboutWindow.SetFixedSize2(400, 100)
	aboutWindow.SetWindowTitle("关于该软件")

	mainLayout := widgets.NewQVBoxLayout()
	aboutWindow.SetLayout(mainLayout)

	plainText := widgets.NewQPlainTextEdit(nil)
	mainLayout.AddWidget(plainText, 0, 0)

	plainText.SetPlainText(`欢迎使用本软件！本软件是一款多功能串口调试工具。
软件作者:yizheneng <kennalee@163.com> 
软件开发环境:go + github.com/therecipe/qt 
软件授权:本软件完全免费，欢迎使用和commit 
软件项目地址:https://github.com/yizheneng/GolangSerialTools`)

	return aboutWindow
}
