// advancedsendwidget
package main

import (
	"github.com/therecipe/qt/widgets"
)

type AdvancedSendWidget struct {
	*widgets.QWidget
}

func NewAdvancedSendWidget() *AdvancedSendWidget {
	widget := &AdvancedSendWidget{}
	widget.QWidget = widgets.NewQWidget(nil, 0)
	//	mainLayout := widgets.NewQVBoxLayout()

	return widget
}
