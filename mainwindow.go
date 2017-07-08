// mainwindow
package main

import (
	"fmt"
	_ "reflect"
	"strconv"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/serialport"
	"github.com/therecipe/qt/widgets"
)

type MainWindow struct {
	*widgets.QWidget
	serialPort   *serialport.QSerialPort
	portOpenFlag bool

	/// 串口参数控件
	portNameCombox     *widgets.QComboBox
	buadCombox         *widgets.QComboBox
	dataBitCombox      *widgets.QComboBox
	checkBitCombox     *widgets.QComboBox
	stopBitCombox      *widgets.QComboBox
	asciiSendButton    *widgets.QRadioButton
	asciiReceiveButton *widgets.QRadioButton

	/// 按钮控件
	sendButton                *widgets.QPushButton
	openPortToolButton        *widgets.QToolButton
	autoNewLineReciveCheckBox *widgets.QCheckBox
	displayTimeCheckBox       *widgets.QCheckBox

	/// 数据显示
	receiveDataDisplay *widgets.QTextEdit
	sendDataDisplay    *widgets.QTextEdit

	/// 数据缓存
	receiveDataBuf   *core.QByteArray
	autoNewLineTimer *core.QTimer
}

func NewMainwindow() (mainWindow *MainWindow) {
	mainWindow = &MainWindow{portOpenFlag: false}
	mainWindow.QWidget = widgets.NewQWidget(nil, 0)
	mainWindow.SetMinimumSize2(300, 200)
	mainWindow.SetWindowTitle("串口调试工具")

	mainWindow.receiveDataBuf = core.NewQByteArray()
	mainWindow.serialPort = serialport.NewQSerialPort(nil)
	mainWindow.serialPort.ConnectReadyRead(mainWindow.readData)
	mainWindow.autoNewLineTimer = core.NewQTimer(nil)
	mainWindow.autoNewLineTimer.ConnectTimeout(mainWindow.receiveAutoNewLineTimeOut)

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
	mainWindow.portNameCombox = widgets.NewQComboBox(nil)
	mainWindow.buadCombox = widgets.NewQComboBox(nil)
	mainWindow.dataBitCombox = widgets.NewQComboBox(nil)
	mainWindow.checkBitCombox = widgets.NewQComboBox(nil)
	mainWindow.stopBitCombox = widgets.NewQComboBox(nil)
	portSettingLayout.AddWidget(portNameLabel, 0, 0, 0)
	portSettingLayout.AddWidget(buadLabel, 1, 0, 0)
	portSettingLayout.AddWidget(dataBitLabel, 2, 0, 0)
	portSettingLayout.AddWidget(checkBitLabel, 3, 0, 0)
	portSettingLayout.AddWidget(stopBitLabel, 4, 0, 0)
	portSettingLayout.AddWidget(mainWindow.portNameCombox, 0, 1, 0)
	portSettingLayout.AddWidget(mainWindow.buadCombox, 1, 1, 0)
	portSettingLayout.AddWidget(mainWindow.dataBitCombox, 2, 1, 0)
	portSettingLayout.AddWidget(mainWindow.checkBitCombox, 3, 1, 0)
	portSettingLayout.AddWidget(mainWindow.stopBitCombox, 4, 1, 0)
	/// 接收设置
	mainWindow.asciiReceiveButton = widgets.NewQRadioButton2("ASCII", nil)
	hexReceiveButton := widgets.NewQRadioButton2("Hex", nil)
	mainWindow.autoNewLineReciveCheckBox = widgets.NewQCheckBox2("自动换行", nil)
	mainWindow.displayTimeCheckBox = widgets.NewQCheckBox2("显示时间", nil)
	receiveSeetingLayout.AddWidget(mainWindow.asciiReceiveButton, 0, 0, 0)
	receiveSeetingLayout.AddWidget(hexReceiveButton, 0, 1, 0)
	receiveSeetingLayout.AddWidget(mainWindow.autoNewLineReciveCheckBox, 1, 0, 0)
	receiveSeetingLayout.AddWidget(mainWindow.displayTimeCheckBox, 2, 0, 0)
	/// 发送设置
	mainWindow.asciiSendButton = widgets.NewQRadioButton2("ASCII", nil)
	hexSendButton := widgets.NewQRadioButton2("Hex", nil)
	reSendCheckButton := widgets.NewQCheckBox2("重复发送:", nil)
	reSendSpinBox := widgets.NewQSpinBox(nil)
	reSendLabel := widgets.NewQLabel2("ms", nil, 0)
	sendSettingLayout.AddWidget(mainWindow.asciiSendButton, 0, 0, 0)
	sendSettingLayout.AddWidget(hexSendButton, 0, 1, 0)
	sendSettingLayout.AddWidget3(reSendCheckButton, 1, 0, 1, 2, 0)
	sendSettingLayout.AddWidget(reSendSpinBox, 2, 0, 0)
	sendSettingLayout.AddWidget(reSendLabel, 2, 1, 0)

	/// 发送数据显示
	sendDataDisplayLayout := widgets.NewQHBoxLayout()
	mainWindow.sendDataDisplay = widgets.NewQTextEdit(nil)
	sendButtonLayout := widgets.NewQVBoxLayout()
	mainWindow.sendButton = widgets.NewQPushButton2("打开串口", nil)
	advancedButton := widgets.NewQPushButton2("高级发送>>", nil)
	sendButtonLayout.AddWidget(mainWindow.sendButton, 0, 0)
	sendButtonLayout.AddWidget(advancedButton, 0, 0)
	sendDataDisplayLayout.AddWidget(mainWindow.sendDataDisplay, 0, 0)
	sendDataDisplayLayout.AddLayout(sendButtonLayout, 0)

	/// 数据显示
	mainWindow.receiveDataDisplay = widgets.NewQTextEdit(nil)
	dataDisplayLayout.AddWidget(mainWindow.receiveDataDisplay, 0, 0)
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
	mainWindow.openPortToolButton = widgets.NewQToolButton(nil)
	mainWindow.openPortToolButton.SetObjectName("openPortToolButton")
	mainWindow.openPortToolButton.SetToolTip("打开串口")
	mainWindow.openPortToolButton.SetCheckable(true)
	toolBar.AddWidget(mainWindow.openPortToolButton)

	/// 控件数据初始化
	serialsInfo := serialport.QSerialPortInfo{}
	for _, serialInfo := range serialsInfo.AvailablePorts() {
		mainWindow.portNameCombox.AddItem(serialInfo.PortName(), core.NewQVariant())
	}

	mainWindow.buadCombox.AddItems([]string{"115200", "57600", "38400", "19200", "9600"})
	mainWindow.dataBitCombox.AddItems([]string{"8", "7"})
	mainWindow.checkBitCombox.AddItems([]string{"None", "Even", "Odd", "Mark", "Space"})
	mainWindow.stopBitCombox.AddItems([]string{"1", "1.5", "2"})
	mainWindow.asciiReceiveButton.SetChecked(true)
	mainWindow.asciiSendButton.SetChecked(true)

	/// 控件功能绑定
	/// 发送按钮按下
	mainWindow.sendButton.ConnectClicked(func(checked bool) { ///< 发送按钮
		if !mainWindow.portOpenFlag {
			if mainWindow.openSerialPort() {
				mainWindow.openPortToolButton.SetChecked(true)
			}
		} else {
			mainWindow.sendData()
		}
	})
	/// 打开按钮按下
	mainWindow.openPortToolButton.ConnectClicked(func(checked bool) {
		fmt.Printf("clicked")
		if checked && !mainWindow.portOpenFlag {
			if !mainWindow.openSerialPort() {
				mainWindow.openPortToolButton.SetChecked(false)
			}
		} else if !checked && mainWindow.portOpenFlag {
			mainWindow.closeSerialPort()
		}
	})
	/// 高级发送按钮按下
	advancedButton.ConnectClicked(func(checked bool) { ///< 高级发送按钮
		if advancedWidget.IsHidden() {
			advancedWidget.Show()
		} else {
			advancedWidget.Hide()
		}
	})

	return
}

/// 打开串口
func (mainWindow *MainWindow) openSerialPort() (result bool) {
	result = true
	if !mainWindow.serialPort.IsOpen() {
		mainWindow.serialPort.SetPortName(mainWindow.portNameCombox.CurrentText())
		buadRate, _ := strconv.Atoi(mainWindow.buadCombox.CurrentText())
		mainWindow.serialPort.SetBaudRate(buadRate, serialport.QSerialPort__AllDirections)

		switch mainWindow.checkBitCombox.CurrentIndex() {
		case 0:
			mainWindow.serialPort.SetParity(serialport.QSerialPort__NoParity)
		case 1:
			mainWindow.serialPort.SetParity(serialport.QSerialPort__EvenParity)
		case 2:
			mainWindow.serialPort.SetParity(serialport.QSerialPort__OddParity)
		case 3:
			mainWindow.serialPort.SetParity(serialport.QSerialPort__MarkParity)
		case 4:
			mainWindow.serialPort.SetParity(serialport.QSerialPort__SpaceParity)
		}

		switch mainWindow.dataBitCombox.CurrentIndex() {
		case 0:
			mainWindow.serialPort.SetDataBits(serialport.QSerialPort__Data8)
		case 1:
			mainWindow.serialPort.SetDataBits(serialport.QSerialPort__Data7)
		}

		switch mainWindow.stopBitCombox.CurrentIndex() {
		case 0:
			mainWindow.serialPort.SetStopBits(serialport.QSerialPort__OneStop)
		case 1:
			mainWindow.serialPort.SetStopBits(serialport.QSerialPort__OneAndHalfStop)
		case 2:
			mainWindow.serialPort.SetStopBits(serialport.QSerialPort__TwoStop)
		}

		if !mainWindow.serialPort.Open(core.QIODevice__ReadWrite) {
			widgets.QMessageBox_Critical(mainWindow, "错误", "打开串口错误！！", widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
			return false
		}
	}

	mainWindow.portNameCombox.SetDisabled(true)
	mainWindow.buadCombox.SetDisabled(true)
	mainWindow.checkBitCombox.SetDisabled(true)
	mainWindow.dataBitCombox.SetDisabled(true)
	mainWindow.stopBitCombox.SetDisabled(true)
	mainWindow.sendButton.SetText("发  送")
	mainWindow.openPortToolButton.SetToolTip("关闭串口")
	mainWindow.portOpenFlag = true
	mainWindow.autoNewLineTimer.Start(100)

	return
}

/// 关闭串口
func (mainWindow *MainWindow) closeSerialPort() {
	if mainWindow.serialPort.IsOpen() {
		mainWindow.serialPort.Close()
	}

	mainWindow.portNameCombox.SetDisabled(false)
	mainWindow.buadCombox.SetDisabled(false)
	mainWindow.checkBitCombox.SetDisabled(false)
	mainWindow.dataBitCombox.SetDisabled(false)
	mainWindow.stopBitCombox.SetDisabled(false)
	mainWindow.sendButton.SetText("打开串口")
	mainWindow.openPortToolButton.SetToolTip("打开串口")
	mainWindow.portOpenFlag = false
	mainWindow.autoNewLineTimer.Stop()
}

/// 发送数据
func (mainWindow *MainWindow) sendData() {

}

/// 读取串口数据
func (mainWindow *MainWindow) readData() {
	if mainWindow.receiveDataBuf.Size() <= 0 {
		mainWindow.autoNewLineTimer.Start(100)
	}
	mainWindow.receiveDataBuf.Append(mainWindow.serialPort.ReadAll())
}

func (mainWindow *MainWindow) receiveAutoNewLineTimeOut() {
	if mainWindow.receiveDataBuf.Size() <= 0 {
		return
	}

	data := mainWindow.receiveDataBuf
	stringData := data.Data()
	stringData = stringData[:len(stringData)-1]

	if !mainWindow.asciiReceiveButton.IsChecked() {
		tempString := ""
		for _, byteData := range []byte(stringData) {
			byteString := strconv.FormatInt(int64(byteData), 16)
			if len(byteString) == 1 {
				tempString = tempString + "0" + byteString + " "
			} else {
				tempString = tempString + byteString + " "
			}
		}
		stringData = tempString
	}

	fmt.Println(stringData + "--")

	if mainWindow.displayTimeCheckBox.IsChecked() {
		currentTime := core.QDateTime_CurrentDateTime()
		stringData = currentTime.ToString("hh:mm:ss.zzz: ") + stringData
	}

	workCursor := mainWindow.receiveDataDisplay.TextCursor()
	workCursor.MovePosition(gui.QTextCursor__EndOfBlock, gui.QTextCursor__MoveAnchor, 1)
	if mainWindow.autoNewLineReciveCheckBox.IsChecked() {
		mainWindow.receiveDataDisplay.InsertHtml(stringData)
		mainWindow.receiveDataDisplay.InsertPlainText("\n")
	} else {
		mainWindow.receiveDataDisplay.InsertHtml(stringData)
	}
	mainWindow.receiveDataDisplay.VerticalScrollBar().SetValue(mainWindow.receiveDataDisplay.VerticalScrollBar().Maximum())

	mainWindow.receiveDataBuf.Clear()
}
