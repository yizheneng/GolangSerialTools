// mainwindow
package main

import (
	"fmt"
	_ "reflect"

	"github.com/therecipe/qt/serialport"
	"github.com/therecipe/qt/widgets"
)

type MainWindow struct {
	*widgets.QWidget
}

func NewMainwindow() (mainWindow *MainWindow) {
	mainWindow = &MainWindow{}
	mainWindow.QWidget = widgets.NewQWidget(nil, 0)
	mainWindow.SetMinimumSize2(500, 600)
	mainWindow.SetWindowTitle("串口调试工具")

	mainLayout := widgets.NewQVBoxLayout()
	widgetsLayout := widgets.NewQHBoxLayout()
	settingLayout := widgets.NewQVBoxLayout()
	dataDisplay := widgets.NewQVBoxLayout()
	menuBar := widgets.NewQMenuBar(nil)

	mainLayout.AddWidget(menuBar, 0, 0)
	mainLayout.AddLayout(widgetsLayout, 0)
	widgetsLayout.AddLayout(settingLayout, 0)
	widgetsLayout.AddLayout(dataDisplay, 0)

	portSettingGroup := widgets.NewQGroupBox2("串口配置", nil)
	portSettingLayout := widgets.NewQGridLayout2()
	portSettingGroup.SetLayout(portSettingLayout)
	receiveSettingGroup := widgets.NewQGroupBox2("接收设置", nil)
	receiveSeetingLayout := widgets.NewQGridLayout2()
	receiveSettingGroup.SetLayout(receiveSeetingLayout)
	sendSettingGroup := widgets.NewQGroupBox2("发送设置", nil)
	sendSettingLayout := widgets.NewQGridLayout2()
	sendSettingGroup.SetLayout(sendSettingLayout)

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

	asciiReceiveButton := widgets.NewQRadioButton2("ASCII", nil)
	hexReceiveButton := widgets.NewQRadioButton2("Hex", nil)
	autoNewLineReciveCheckBox := widgets.NewQCheckBox2("自动换行", nil)
	displayTimeCheckBox := widgets.NewQCheckBox2("显示时间", nil)
	receiveSeetingLayout.AddWidget(asciiReceiveButton, 0, 0, 0)
	receiveSeetingLayout.AddWidget(hexReceiveButton, 0, 1, 0)
	receiveSeetingLayout.AddWidget(autoNewLineReciveCheckBox, 1, 0, 0)
	receiveSeetingLayout.AddWidget(displayTimeCheckBox, 2, 0, 0)

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

	settingLayout.AddWidget(portSettingGroup, 0, 0)
	settingLayout.AddWidget(receiveSettingGroup, 0, 0)
	settingLayout.AddWidget(sendSettingGroup, 0, 0)
	settingLayout.AddStretch(20)

	mainWindow.SetLayout(mainLayout)
	serialsInfo := serialport.QSerialPortInfo{}
	for _, serialInfo := range serialsInfo.AvailablePorts() {
		fmt.Println(serialInfo.PortName())
	}

	return
}
