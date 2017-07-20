// mainwindow
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	_ "reflect"
	"strconv"
	"strings"

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
	reSendCheckButton         *widgets.QCheckBox
	reSendSpinBox             *widgets.QSpinBox

	/// 数据显示
	receiveDataDisplay    *widgets.QTextEdit
	sendDataDisplay       *widgets.QTextEdit
	historySendListWidget *widgets.QTableWidget

	/// 数据缓存
	receiveDataBuf   *core.QByteArray
	autoNewLineTimer *core.QTimer
	reSendTimer      *core.QTimer /// 重复发送定时器
	sendStringBuf    string       /// 发送数据缓存区
}

type SettingType struct {
	SendHistorys []string
}

func NewMainwindow() (mainWindow *MainWindow) {
	mainWindow = &MainWindow{portOpenFlag: false}
	mainWindow.QWidget = widgets.NewQWidget(nil, 0)
	mainWindow.SetMinimumWidth(400)
	mainWindow.SetFixedHeight(600)
	mainWindow.SetWindowTitle("串口调试工具")

	mainWindow.receiveDataBuf = core.NewQByteArray()
	mainWindow.serialPort = serialport.NewQSerialPort(nil)
	mainWindow.serialPort.ConnectReadyRead(mainWindow.readData)
	mainWindow.autoNewLineTimer = core.NewQTimer(nil)
	mainWindow.autoNewLineTimer.ConnectTimeout(mainWindow.receiveAutoNewLineTimeOut)
	mainWindow.reSendTimer = core.NewQTimer(nil)
	mainWindow.reSendTimer.ConnectTimeout(mainWindow.reSendTimeOut)

	mainLayout := widgets.NewQVBoxLayout()
	widgetsLayout := widgets.NewQHBoxLayout()
	settingLayout := widgets.NewQVBoxLayout()
	dataDisplayLayout := widgets.NewQVBoxLayout()
	toolBar := widgets.NewQToolBar("工具栏", nil)
	mainWindow.historySendListWidget = widgets.NewQTableWidget(nil) ///< 发送历史LIST
	mainWindow.historySendListWidget.SetColumnCount(1)
	mainWindow.historySendListWidget.SetHorizontalHeaderLabels([]string{"历史数据"})
	mainWindow.historySendListWidget.HorizontalHeader().SetStretchLastSection(true)
	mainWindow.historySendListWidget.SetEditTriggers(widgets.QAbstractItemView__NoEditTriggers)

	mainLayout.AddWidget(toolBar, 0, 0)
	mainLayout.AddLayout(widgetsLayout, 0)
	widgetsLayout.AddLayout(settingLayout, 0)
	widgetsLayout.AddLayout(dataDisplayLayout, 0)
	widgetsLayout.AddWidget(mainWindow.historySendListWidget, 0, 0)

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
	mainWindow.reSendCheckButton = widgets.NewQCheckBox2("重复发送:", nil)
	mainWindow.reSendSpinBox = widgets.NewQSpinBox(nil)
	mainWindow.reSendSpinBox.SetMinimum(1)
	mainWindow.reSendSpinBox.SetMaximum(99999)
	mainWindow.reSendSpinBox.SetValue(1000)
	reSendLabel := widgets.NewQLabel2("ms", nil, 0)
	sendSettingLayout.AddWidget(mainWindow.asciiSendButton, 0, 0, 0)
	sendSettingLayout.AddWidget(hexSendButton, 0, 1, 0)
	sendSettingLayout.AddWidget3(mainWindow.reSendCheckButton, 1, 0, 1, 2, 0)
	sendSettingLayout.AddWidget(mainWindow.reSendSpinBox, 2, 0, 0)
	sendSettingLayout.AddWidget(reSendLabel, 2, 1, 0)

	/// 发送数据显示
	sendDataDisplayLayout := widgets.NewQHBoxLayout()
	mainWindow.sendDataDisplay = widgets.NewQTextEdit(nil)
	sendButtonLayout := widgets.NewQVBoxLayout()
	mainWindow.sendButton = widgets.NewQPushButton2("打开串口", nil)
	clearSendButton := widgets.NewQPushButton2("清除输入", nil)
	advancedButton := widgets.NewQPushButton2("高级发送>>", nil)
	sendButtonLayout.AddWidget(mainWindow.sendButton, 0, 0)
	sendButtonLayout.AddWidget(clearSendButton, 0, 0)
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
	advancedWidget.SetFixedHeight(100)
	advancedWidget.Hide()
	mainLayout.AddWidget(advancedWidget, 0, 0)

	/// 工具栏
	toolBar.SetObjectName("toolbar")
	toolBar.SetMinimumHeight(40)
	mainWindow.openPortToolButton = widgets.NewQToolButton(nil)
	mainWindow.openPortToolButton.SetObjectName("openPortToolButton")
	mainWindow.openPortToolButton.SetToolTip("打开串口")
	mainWindow.openPortToolButton.SetCheckable(true)
	mainWindow.openPortToolButton.SetFixedSize2(40, 40)
	clearReceiveToolButton := widgets.NewQToolButton(nil)
	clearReceiveToolButton.SetObjectName("clearReceiveToolButton")
	clearReceiveToolButton.SetToolTip("清除接收数据")
	clearReceiveToolButton.SetCheckable(true)
	clearReceiveToolButton.SetFixedSize2(40, 40)
	clearSendToolButton := widgets.NewQToolButton(nil)
	clearSendToolButton.SetObjectName("clearSendToolButton")
	clearSendToolButton.SetToolTip("清除发送数据")
	clearSendToolButton.SetCheckable(true)
	clearSendToolButton.SetFixedSize2(40, 40)
	clearHistoryToolButton := widgets.NewQToolButton(nil)
	clearHistoryToolButton.SetObjectName("clearHistoryToolButton")
	clearHistoryToolButton.SetToolTip("清除历史记录")
	clearHistoryToolButton.SetCheckable(true)
	clearHistoryToolButton.SetFixedSize2(40, 40)
	infoToolButton := widgets.NewQToolButton(nil)
	infoToolButton.SetObjectName("infoToolButton")
	infoToolButton.SetCheckable(true)
	infoToolButton.SetFixedSize2(40, 40)
	toolBar.AddWidget(mainWindow.openPortToolButton)
	toolBar.AddSeparator()
	toolBar.AddWidget(clearReceiveToolButton)
	toolBar.AddWidget(clearSendToolButton)
	toolBar.AddWidget(clearHistoryToolButton)
	toolBar.AddSeparator()
	toolBar.AddWidget(infoToolButton)

	/// 控件数据初始化
	serialsInfo := serialport.QSerialPortInfo{}
	for _, serialInfo := range serialsInfo.AvailablePorts() {
		mainWindow.portNameCombox.AddItem(serialInfo.PortName(), core.NewQVariant())
	}

	file, fileErr := os.OpenFile("setting.json", os.O_RDONLY, 0666)
	defer file.Close()
	settingData, _ := ioutil.ReadAll(file)

	settings := &SettingType{}
	jsonErr := json.Unmarshal(settingData, &settings)
	/// 历史数据界面数据初始化
	if jsonErr == nil && fileErr == nil {
		mainWindow.historySendListWidget.SetRowCount(len(settings.SendHistorys))
		for i := 0; i < len(settings.SendHistorys); i++ {
			mainWindow.historySendListWidget.SetItem(i, 0, widgets.NewQTableWidgetItem2(settings.SendHistorys[i], 0))
			mainWindow.historySendListWidget.Item(i, 0).SetToolTip(settings.SendHistorys[i])
		}
	} else {
		fmt.Errorf("Open file error Or json Unmarshal error")
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
			mainWindow.SetFixedHeight(mainWindow.Height() + advancedWidget.Height() + mainLayout.Spacing())
		} else {
			mainWindow.SetFixedHeight(mainWindow.Height() - advancedWidget.Height() - mainLayout.Spacing())
			advancedWidget.Hide()
		}
	})
	/// 清除输入按钮按下
	clearSendButton.ConnectClicked(func(checked bool) {
		mainWindow.sendDataDisplay.Clear()
	})
	/// 界面关闭
	mainWindow.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		mainWindow.closeDispose()
	})
	/// 历史发送列表单击
	mainWindow.historySendListWidget.ConnectCellClicked(func(row, column int) {
		mainWindow.sendDataDisplay.SetPlainText(mainWindow.historySendListWidget.Item(row, column).Text())
	})
	/// 历史发送列表双击
	mainWindow.historySendListWidget.ConnectCellDoubleClicked(func(row, column int) {
		text := mainWindow.historySendListWidget.Item(row, column).Text()
		mainWindow.sendDateWithDataToSerial(text)
	})
	/// 清除接收按钮
	clearReceiveToolButton.ConnectClicked(func(checked bool) {
		mainWindow.receiveDataDisplay.Clear()
	})
	/// 清除发送按钮按下
	clearSendToolButton.ConnectClicked(func(checked bool) {
		mainWindow.sendDataDisplay.Clear()
	})
	/// 清除发送历史按钮按下
	clearHistoryToolButton.ConnectClicked(func(checked bool) {
		mainWindow.historySendListWidget.SetRowCount(0)
	})
	/// 重复发送按钮按下
	mainWindow.reSendCheckButton.ConnectClicked(func(checked bool) {
		if !checked && (mainWindow.sendButton.Text() == "停止发送") {
			mainWindow.sendButton.SetText("发送")
			mainWindow.reSendTimer.Stop()
		}
	})
	/// 重复发送间隔发生变化
	mainWindow.reSendSpinBox.ConnectValueChanged(func(newVal int) {
		mainWindow.reSendTimer.Start(newVal)
	})
	/// 关于按钮按下
	infoToolButton.ConnectClicked(func(checked bool) {
		NewAboutWindow(mainWindow).Show()
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
	sendDataString := mainWindow.sendDataDisplay.ToPlainText()
	if len(sendDataString) <= 0 {
		widgets.QMessageBox_Warning(mainWindow, "警告", "请输入数据进行发送！！", widgets.QMessageBox__Ok, widgets.QMessageBox__NoButton)
		return
	}

	if mainWindow.sendButton.Text() == "停止发送" {
		mainWindow.sendButton.SetText("发送")
		mainWindow.reSendTimer.Stop()
		return
	}

	if mainWindow.reSendCheckButton.IsChecked() {
		mainWindow.sendButton.SetText("停止发送")
		mainWindow.reSendTimer.Start(mainWindow.reSendSpinBox.Value())
	} else {
		mainWindow.sendDateWithDataToSerial(sendDataString)
	}

	if sendDataString == "" {
		return
	}

	if sendDataString != mainWindow.historySendListWidget.Item(0, 0).Text() {
		mainWindow.historySendListWidget.InsertRow(0)
		tableItem := widgets.NewQTableWidgetItem2(sendDataString, 0)
		tableItem.SetToolTip(sendDataString)
		mainWindow.historySendListWidget.SetItem(0, 0, tableItem)
	}
}

/// 发送数据到串口
func (mainWindow *MainWindow) sendDateWithDataToSerial(data string) {
	if !mainWindow.asciiSendButton.IsChecked() {
		rx := core.NewQRegExp2("([a-fA-F0-9]{2}[ ]{0,1})*", core.Qt__CaseSensitive, core.QRegExp__RegExp)
		rx.IndexIn(data, 0, core.QRegExp__CaretAtZero)
		resultList := rx.CapturedTexts()
		if len(resultList) < 1 {
			return
		}
		resultString := resultList[0]
		var sendData []byte
		for _, byteString := range strings.Split(resultString, " ") {
			if byteString == "" {
				continue
			}

			if len(byteString) <= 2 {
				val, _ := strconv.ParseUint(byteString, 16, 8)
				sendData = append(sendData, byte(val))
			} else {
				byteStringSize := len(byteString)
				for i := 0; i < byteStringSize/2; i++ {
					tempString := byteString[i*2 : i*2+2]
					val, _ := strconv.ParseUint(tempString, 16, 8)
					sendData = append(sendData, byte(val))
				}
			}
		}
		fmt.Println("---Result:", sendData)
		data = string(sendData)
	}

	if mainWindow.serialPort.IsOpen() {
		mainWindow.serialPort.Write2(data)
	}
}

/// 读取串口数据
func (mainWindow *MainWindow) readData() {
	if mainWindow.receiveDataBuf.Size() <= 0 {
		mainWindow.autoNewLineTimer.Start(100)
	}
	mainWindow.receiveDataBuf.Append(mainWindow.serialPort.ReadAll())
}

/// 自动换行定时器回调
func (mainWindow *MainWindow) receiveAutoNewLineTimeOut() {
	if mainWindow.receiveDataBuf.Size() <= 0 {
		return
	}

	data := mainWindow.receiveDataBuf
	stringData := data.Data()
	//stringData = stringData[:len(stringData)-1]

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

	if mainWindow.displayTimeCheckBox.IsChecked() {
		currentTime := core.QDateTime_CurrentDateTime()
		stringData = currentTime.ToString("hh:mm:ss.zzz: ") + stringData
	}

	workCursor := mainWindow.receiveDataDisplay.TextCursor()
	workCursor.MovePosition(gui.QTextCursor__End, gui.QTextCursor__MoveAnchor, 1)
	if mainWindow.autoNewLineReciveCheckBox.IsChecked() {
		mainWindow.receiveDataDisplay.InsertHtml(stringData)
		mainWindow.receiveDataDisplay.InsertPlainText("\n")
	} else {
		mainWindow.receiveDataDisplay.InsertHtml(stringData)
	}
	mainWindow.receiveDataDisplay.VerticalScrollBar().SetValue(mainWindow.receiveDataDisplay.VerticalScrollBar().Maximum())

	mainWindow.receiveDataBuf.Clear()
}

/// 关闭时的处理
func (mainWindow *MainWindow) closeDispose() {
	settings := &SettingType{}
	/// 历史数据界面数据初始化
	for i := 0; i < mainWindow.historySendListWidget.RowCount(); i++ {
		settings.SendHistorys = append(settings.SendHistorys, mainWindow.historySendListWidget.Item(i, 0).Text())
	}
	fmt.Println(settings.SendHistorys)
	byteData, err := json.Marshal(&settings)
	fmt.Println(string(byteData))
	if err == nil {
		ioutil.WriteFile("setting.json", byteData, os.ModeCharDevice)
	}
}

/// 重复发送定时器超时
func (mainWindow *MainWindow) reSendTimeOut() {
	sendDataString := mainWindow.sendDataDisplay.ToPlainText()
	mainWindow.sendDateWithDataToSerial(sendDataString)
}
