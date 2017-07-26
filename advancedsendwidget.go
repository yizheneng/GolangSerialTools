// advancedsendwidget
package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type AdvancedSendWidget struct {
	*widgets.QWidget

	tableWidget    *widgets.QTableWidget
	sendDataOnce   func(data string, sendMode int)
	dataTimer      *core.QTimer
	sendDataConfig []AdvancedSendStruct /// 发送数据配置
	sendCounters   []int                /// 发送数据计数器
	sendTimes      []int                /// 发送数据的定时器间隔
}

type AdvancedSendStruct struct {
	Data      string
	InputMode int
	Interval  int
	Enable    bool
}

func NewAdvancedSendWidget() *AdvancedSendWidget {
	widget := &AdvancedSendWidget{}
	widget.QWidget = widgets.NewQWidget(nil, 0)
	widget.dataTimer = core.NewQTimer(nil)

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
	startSendButton := widgets.NewQPushButton2("开始发送", nil)

	buttonLayout := widgets.NewQHBoxLayout()
	buttonLayout.AddStretch(10)
	buttonLayout.AddWidget(addItemButton, 0, 0)
	buttonLayout.AddWidget(removeItemButton, 0, 0)
	buttonLayout.AddWidget(startSendButton, 0, 0)

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
	/// 单击立即发送按钮
	widget.tableWidget.ConnectCellClicked(func(row, column int) {
		if column != 4 {
			return
		}

		fmt.Println("clicked row:", row)
		widget.sendDataOnce(widget.tableWidget.Item(row, 0).Text(), widgets.NewQComboBoxFromPointer(widget.tableWidget.CellWidget(row, 1).Pointer()).CurrentIndex())
	})
	/// 开始发送按钮按下
	startSendButton.ConnectClicked(func(clicked bool) {
		if startSendButton.Text() == "开始发送" {
			startSendButton.SetText("停止发送")
			widget.tableWidget.SetDisabled(true)
			addItemButton.SetDisabled(true)
			removeItemButton.SetDisabled(true)

			widget.sendCounters = []int{}
			widget.sendTimes = []int{}
			widget.sendDataConfig = widget.GetSettings()
			var temp []AdvancedSendStruct
			var tempNums []int
			for _, sendStruct := range widget.sendDataConfig {
				if sendStruct.Enable {
					temp = append(temp, sendStruct)
					tempNums = append(tempNums, sendStruct.Interval)
					widget.sendCounters = append(widget.sendCounters, 0)
				}
			}

			widget.sendDataConfig = temp

			gcd := GetGCD(tempNums)
			for _, num := range tempNums {
				widget.sendTimes = append(widget.sendTimes, num/gcd)
			}

			widget.dataTimer.Start(gcd)
		} else {
			startSendButton.SetText("开始发送")
			widget.tableWidget.SetDisabled(false)
			addItemButton.SetDisabled(false)
			removeItemButton.SetDisabled(false)
			widget.dataTimer.Stop()
		}
	})
	/// 定时器溢出
	widget.dataTimer.ConnectTimeout(func() {
		for i, counter := range widget.sendCounters {
			if (counter + 1) >= widget.sendTimes[i] {
				widget.sendCounters[i] = 0
				widget.sendDataOnce(widget.sendDataConfig[i].Data, widget.sendDataConfig[i].InputMode)
			} else {
				widget.sendCounters[i] = counter + 1
			}
		}
	})

	return widget
}

/// 添加一行配置
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

	widget.tableWidget.SetItem(rowCount, 0, widgets.NewQTableWidgetItem2(content, 0))
	widget.tableWidget.SetCellWidget(rowCount, 1, inputModeCellWidget)
	widget.tableWidget.SetCellWidget(rowCount, 2, intervalCellWidget)
	widget.tableWidget.SetCellWidget(rowCount, 3, enableCellWidget)
	item := widgets.NewQTableWidgetItem2("立即发送", 0)
	item.SetFlags(core.Qt__ItemIsEnabled)
	widget.tableWidget.SetItem(rowCount, 4, item)
}

/// 获取配置数据
func (widget *AdvancedSendWidget) GetSettings() []AdvancedSendStruct {
	var result []AdvancedSendStruct

	for i := 0; i < widget.tableWidget.RowCount(); i++ {
		var temp AdvancedSendStruct

		temp.Data = widget.tableWidget.Item(i, 0).Text()
		temp.InputMode = widgets.NewQComboBoxFromPointer(widget.tableWidget.CellWidget(i, 1).Pointer()).CurrentIndex()
		temp.Interval = widgets.NewQSpinBoxFromPointer(widget.tableWidget.CellWidget(i, 2).Pointer()).Value()
		temp.Enable = widgets.NewQCheckBoxFromPointer(widget.tableWidget.CellWidget(i, 3).Pointer()).IsChecked()
		fmt.Println("temp:", temp)
		result = append(result, temp)
	}

	return result
}

/// 设置高级发送配置
func (widget *AdvancedSendWidget) SetSettings(settings []AdvancedSendStruct) {
	for _, setting := range settings {
		fmt.Println("setting:", setting)
		widget.addNewRow(setting.Data, setting.InputMode, setting.Interval, setting.Enable)
	}
}

func (widget *AdvancedSendWidget) ConnectSendDataOnce(f func(data string, sendMode int)) {
	widget.sendDataOnce = f
}

/// 获取最大公约数
func GetGCD(nums []int) int {
	var i, j, temp int
	for j = 0; j < len(nums)-1; j++ {
		for i = 0; i < len(nums)-1-j; i++ {
			if nums[i] > nums[i+1] {
				temp = nums[i]
				nums[i] = nums[i+1]
				nums[i+1] = temp
			}
		}
	}

	if nums[0] == 1 {
		return 1
	}

	result := nums[0] + 1

	for {
		result--
		if result == 1 {
			break
		}

		i := 0
		for ; i < len(nums); i++ {
			if nums[i]%result != 0 {
				break
			}
		}

		if i == len(nums) {
			return result
		}
	}

	return 1
}
