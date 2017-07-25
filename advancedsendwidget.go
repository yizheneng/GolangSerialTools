// advancedsendwidget
package main

import (
	"github.com/therecipe/qt/widgets"
)

type AdvancedSendWidget struct {
	*widgets.QWidget

	tableWidget *widgets.QTableWidget
}

func NewAdvancedSendWidget() *AdvancedSendWidget {
	widget := &AdvancedSendWidget{}
	widget.QWidget = widgets.NewQWidget(nil, 0)

	widget.tableWidget = widgets.NewQTableWidget(nil)
	widget.tableWidget.SetColumnCount(5)
	widget.tableWidget.SetHorizontalHeaderLabels([]string{"数据输入", "数据格式", "数据发送间隔", "数据使能", "----"})
	//	widget.tableWidget.HorizontalHeader().SetStretchLastSection(true)
	widget.tableWidget.HorizontalHeader().SetSectionResizeMode2(0, widgets.QHeaderView__Stretch)
	widget.tableWidget.HorizontalHeader().SetSectionResizeMode2(1, widgets.QHeaderView__ResizeToContents)
	widget.tableWidget.HorizontalHeader().SetSectionResizeMode2(2, widgets.QHeaderView__ResizeToContents)
	widget.tableWidget.HorizontalHeader().SetSectionResizeMode2(3, widgets.QHeaderView__ResizeToContents)
	widget.tableWidget.HorizontalHeader().SetSectionResizeMode2(4, widgets.QHeaderView__ResizeToContents)

	addItemButton := widgets.NewQPushButton2("添加", nil)
	removeItemButton := widgets.NewQPushButton2("删除", nil)

	buttonLayout := widgets.NewQHBoxLayout()
	buttonLayout.AddStretch(10)
	buttonLayout.AddWidget(addItemButton, 0, 0)
	buttonLayout.AddWidget(removeItemButton, 0, 0)

	mainLayout := widgets.NewQVBoxLayout()
	mainLayout.AddWidget(widget.tableWidget, 0, 0)
	mainLayout.AddLayout(buttonLayout, 0)

	widget.SetLayout(mainLayout)

	/// 空间功能绑定
	/// 添加按钮
	addItemButton.ConnectClicked(func(clicked bool) {
		widget.addNewRow("", 1, 1000, true)
	})
	/// 删除按钮
	removeItemButton.ConnectClicked(func(clicked bool) {
		widget.tableWidget.RemoveRow(widget.tableWidget.CurrentRow())
	})
	return widget
}

func (widget *AdvancedSendWidget) addNewRow(content string, inputMode int, interval int, enable bool) {
	inputModeCellWidget := widgets.NewQComboBox(nil)
	inputModeCellWidget.AddItems([]string{"十六进制", "ASCII"})
	inputModeCellWidget.SetCurrentIndex(inputMode)

	intervalCellWidget := widgets.NewQSpinBox(nil)
	intervalCellWidget.SetMinimum(1)
	intervalCellWidget.SetMaximum(999999)
	intervalCellWidget.SetValue(interval)

	enableCellWidget := widgets.NewQCheckBox(nil)
	enableCellWidget.SetChecked(enable)

	rowCount := widget.tableWidget.RowCount()
	widget.tableWidget.InsertRow(rowCount)

	widget.tableWidget.SetCellWidget(rowCount, 1, inputModeCellWidget)
	widget.tableWidget.SetCellWidget(rowCount, 2, intervalCellWidget)
	widget.tableWidget.SetCellWidget(rowCount, 3, enableCellWidget)
	widget.tableWidget.SetItem(rowCount, 4, widgets.NewQTableWidgetItem2("立即发送", 0))
}
