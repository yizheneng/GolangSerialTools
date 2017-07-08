// mainwindow
package main

import (
	_ "fmt"
	_ "reflect"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/serialport"
	"github.com/therecipe/qt/widgets"
)

type MainWindow struct {
	*widgets.QWidget
	portOpenFlag bool
}

func NewMainwindow() (mainWindow *MainWindow) {
	mainWindow = &MainWindow{portOpenFlag: false}
	mainWindow.QWidget = widgets.NewQWidget(nil, 0)
	mainWindow.SetMinimumSize2(300, 200)
	mainWindow.SetWindowTitle("串口调试工具")

	mainLayout := widgets.NewQVBoxLayout()
	widgetsLayout := widgets.NewQHBoxLayout()
	settingLayout := widgets.NewQVBoxLayout()
	dataDisplayLayout := widgets.NewQVBoxLayout()
	toolBar := widgets.NewQToolBar("工具栏", nil)
	historySendListWidget := widgets.NewQListWidget(nil) ///< 发送历史LIST

	mainLayout.AddWidget(toolBar, 0, 0)
	mainLayout.AddLayout(widgetsLayout, 0)
	widgetsLayout.AddLayout(settingLayout, 0)
	widgetsLayout.AddLayout(dataDisplayLayout, 0)
	widgetsLayout.AddWidget(historySendListWidget, 0, 0)

	/// 主要布局
	portSettingGroup := widgets.NewQGroupBox2("串口配置", nil)
	portSettingLayout := widgets.NewQGridLayout2()
	portSettingGroup.SetLayout(portSettingLayout)
	receiveSettingGroup := widgets.NewQGroupBox2("接收设置", nil)
	receiveSeetingLayout := widgets.NewQGridLayout2()
	receiveSettingGroup.SetLayout(receiveSeetingLayout)
	sendSettingGroup := widgets.NewQGroupBox2("发送设置", nil)
	sendSettingLayout := widgets.NewQGridLayout2()
	sendSettingGroup.SetLayout(sendSettingLayout)
	/// 串口设置
	portNameLabel := widgets.NewQLabel2("串口号:", nil, 0)
	buadLabel := widgets.NewQLabel2("波特率:", nil, 0)
	dataBitLabel := widgets.NewQLabel2("数据位:", nil, 0)
	checkBitLabel := widgets.NewQLabel2("校验位:", nil, 0)
	stopBitLabel := widgets.NewQLabel2("停止位:", nil, 0)
	portNameCombox := widgets.NewQComboBox(nil)
	buadCombox := widgets.NewQComboBox(nil)
	dataBitCombox := widgets.NewQComboBox(nil)
	checkBitCombox := widgets.NewQComboBox(nil)
	stopBitCombox := widgets.NewQComboBox(nil)
	portSettingLayout.AddWidget(portNameLabel, 0, 0, 0)
	portSettingLayout.AddWidget(buadLabel, 1, 0, 0)
	portSettingLayout.AddWidget(dataBitLabel, 2, 0, 0)
	portSettingLayout.AddWidget(checkBitLabel, 3, 0, 0)
	portSettingLayout.AddWidget(stopBitLabel, 4, 0, 0)
	portSettingLayout.AddWidget(portNameCombox, 0, 1, 0)
	portSettingLayout.AddWidget(buadCombox, 1, 1, 0)
	portSettingLayout.AddWidget(dataBitCombox, 2, 1, 0)
	portSettingLayout.AddWidget(checkBitCombox, 3, 1, 0)
	portSettingLayout.AddWidget(stopBitCombox, 4, 1, 0)
	/// 接收设置
	asciiReceiveButton := widgets.NewQRadioButton2("ASCII", nil)
	hexReceiveButton := widgets.NewQRadioButton2("Hex", nil)
	autoNewLineReciveCheckBox := widgets.NewQCheckBox2("自动换行", nil)
	displayTimeCheckBox := widgets.NewQCheckBox2("显示时间", nil)
	receiveSeetingLayout.AddWidget(asciiReceiveButton, 0, 0, 0)
	receiveSeetingLayout.AddWidget(hexReceiveButton, 0, 1, 0)
	receiveSeetingLayout.AddWidget(autoNewLineReciveCheckBox, 1, 0, 0)
	receiveSeetingLayout.AddWidget(displayTimeCheckBox, 2, 0, 0)
	/// 发送设置
	asciiSendButton := widgets.NewQRadioButton2("ASCII", nil)
	hexSendButton := widgets.NewQRadioButton2("Hex", nil)
	reSendCheckButton := widgets.NewQCheckBox2("重复发送:", nil)
	reSendSpinBox := widgets.NewQSpinBox(nil)
	reSendLabel := widgets.NewQLabel2("ms", nil, 0)
	sendSettingLayout.AddWidget(asciiSendButton, 0, 0, 0)
	sendSettingLayout.AddWidget(hexSendButton, 0, 1, 0)
	sendSettingLayout.AddWidget3(reSendCheckButton, 1, 0, 1, 2, 0)
	sendSettingLayout.AddWidget(reSendSpinBox, 2, 0, 0)
	sendSettingLayout.AddWidget(reSendLabel, 2, 1, 0)

	/// 发送数据显示
	sendDataDisplayLayout := widgets.NewQHBoxLayout()
	sendDataDisplay := widgets.NewQPlainTextEdit(nil)
	sendButtonLayout := widgets.NewQVBoxLayout()
	sendButton := widgets.NewQPushButton2("打开串口", nil)
	advancedButton := widgets.NewQPushButton2("高级发送>>", nil)
	sendButtonLayout.AddWidget(sendButton, 0, 0)
	sendButtonLayout.AddWidget(advancedButton, 0, 0)
	sendDataDisplayLayout.AddWidget(sendDataDisplay, 0, 0)
	sendDataDisplayLayout.AddLayout(sendButtonLayout, 0)

	/// 数据显示
	receiveDataDisplay := widgets.NewQPlainTextEdit(nil)
	dataDisplayLayout.AddWidget(receiveDataDisplay, 0, 0)
	dataDisplayLayout.AddLayout(sendDataDisplayLayout, 0)

	settingLayout.AddWidget(portSettingGroup, 0, 0)
	settingLayout.AddWidget(receiveSettingGroup, 0, 0)
	settingLayout.AddWidget(sendSettingGroup, 0, 0)
	settingLayout.AddStretch(20)

	widgetsLayout.SetStretch(0, 0)
	widgetsLayout.SetStretch(1, 4)
	widgetsLayout.SetStretch(2, 1)
	dataDisplayLayout.SetStretch(0, 5)
	dataDisplayLayout.SetStretch(1, 1)
	mainWindow.SetLayout(mainLayout)

	/// 高级发送显示
	advancedWidget := widgets.NewQWidget(nil, 0)
	advancedWidget.SetMinimumHeight(100)
	advancedWidget.Hide()
	dataDisplayLayout.AddWidget(advancedWidget, 0, 0)

	/// 工具栏
	toolBar.SetObjectName("toolbar")
	toolBar.SetMinimumHeight(30)
	openPortToolButton := widgets.NewQToolButton(nil)
	openPortToolButton.SetObjectName("openPortToolButton")
	openPortToolButton.SetToolTip("打开串口")
	openPortToolButton.SetCheckable(true)
	toolBar.AddWidget(openPortToolButton)

	/// 控件数据初始化
	serialsInfo := serialport.QSerialPortInfo{}
	for _, serialInfo := range serialsInfo.AvailablePorts() {
		portNameCombox.AddItem(serialInfo.PortName(), core.NewQVariant())
	}

	buadCombox.AddItems([]string{"115200", "57600", "38400", "19200", "9600"})
	dataBitCombox.AddItems([]string{"8", "7"})
	checkBitCombox.AddItems([]string{"None", "Even", "Odd", "Mark", "Space"})
	stopBitCombox.AddItems([]string{"1", "1.5", "2"})
	asciiReceiveButton.SetChecked(true)
	asciiSendButton.SetChecked(true)

	/// 控件功能绑定
	sendButton.ConnectClicked(func(checked bool) { ///< 发送按钮
		if mainWindow.portOpenFlag {
			portNameCombox.SetDisabled(true)
			buadCombox.SetDisabled(true)
			checkBitCombox.SetDisabled(true)
			dataBitCombox.SetDisabled(true)
			stopBitCombox.SetDisabled(true)
			sendButton.SetText("发  送")
		} else {
			portNameCombox.SetDisabled(false)
			buadCombox.SetDisabled(false)
			checkBitCombox.SetDisabled(false)
			dataBitCombox.SetDisabled(false)
			stopBitCombox.SetDisabled(false)
			sendButton.SetText("打开串口")
		}

		mainWindow.portOpenFlag = !mainWindow.portOpenFlag
	})

	openPortToolButton.ConnectClicked(func(checked bool) {
		if checked {
			openPortToolButton.SetToolTip("关闭串口")
		} else {
			openPortToolButton.SetToolTip("打开串口")
		}
	})

	advancedButton.ConnectClicked(func(checked bool) { ///< 高级发送按钮
		if advancedWidget.IsHidden() {
			advancedWidget.Show()
		} else {
			advancedWidget.Hide()
		}
	})

	return
}
