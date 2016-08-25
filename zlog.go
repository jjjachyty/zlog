package zlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//log var 非结构化日志
// var log = new(log.Logger)

//ConsoleLoger var 非结构化日志
//var ConsoleLoger = new(log.Logger)

//AppLogFile var 结构化日志
var AppLogFile *os.File

//Action type
type Action int

const (
	//SELECT action
	SELECT Action = 1 << iota
	//UPDATE action
	UPDATE
	//ADD action
	ADD
	//DELETE action
	DELETE
)

//SysLogFile var 日志文件
var SysLogFile *os.File

//AppLoger var 结构化日志
type AppLoger struct {
	OpTime   time.Time   //操作时间
	Method   string      //操作方法
	Action   string      //操作类型 select update add edit
	UserName string      //操作人
	OldData  interface{} //旧数据，此选项只用于update
	NewData  interface{} //目前最新数据
	Explan   string      //备注说明
}

//ErrorWithCode func 带错误号码的输出错误日志
//code 配置文件定义的编号 err 为系统错误，如为自定义则为nil arry 翻译数组填充值
func ErrorWithCode(code string, arry []string, err error) {
	recordLogsWithCode(4, code, arry, err)
}

//Error func 输出错误日志
//logs 输出日志 err 系统原始错误日志
func Error(logs string, err error) {
	recordLogf(4, logs, err)
}

//Errorf func 格式化输出错误日志
//logs 输出日志 err 系统原始错误日志
func Errorf(logs string, err error, vars ...interface{}) {
	recordLogf(4, logs, err, vars...)
}

//WarningWithCode func 带错误号码的输出错误日志
//code 配置文件定义的编号 err 为系统错误，如为自定义则为nil arry 翻译数组填充值
func WarningWithCode(code string, err error, arry []string) {
	recordLogsWithCode(3, code, arry, err)
}

//Warning func  输出错误日志
//logs 输出日志 err 系统原始错误日志
func Warning(logs string, err error) {
	recordLogf(3, logs, err)
}

//Warningf func 格式化输出错误日志
//logs 输出日志 err 系统原始错误日志
func Warningf(logs string, err error, vars ...interface{}) {
	recordLogf(3, logs, err, vars...)
}

//InfoWithCode func 带错误号码的输出错误日志
//code 配置文件定义的编号 err 为系统错误，如为自定义则为nil arry 翻译数组填充值
func InfoWithCode(code string, err error, arry []string) {
	recordLogsWithCode(2, code, arry, err)
}

//Info func  输出错误日志
//logs 输出日志 err 系统原始错误日志
func Info(logs string, err error) {
	recordLogf(2, logs, err)
}

//Infof func  格式化方式输出错误日志
//logs 输出日志 err 系统原始错误日志 vars替换的值
func Infof(logs string, err error, vars ...interface{}) {
	recordLogf(2, logs, err, vars...)
}

//Debug func 输出错误日志
//code 配置文件定义的编号 err 为系统错误，如为自定义则为nil arry 翻译数组填充值
func Debug(logs string, err error) {
	recordLogf(1, logs, err)
}

//Debugf func 输出错误日志
//code 配置文件定义的编号 err 为系统错误，如为自定义则为nil arry 翻译数组填充值
func Debugf(logs string, err error, vars ...interface{}) {
	recordLogf(1, logs, err, vars...)
}

// marshalCode func 翻译code
// code 编码 arry 填充对象
func marshalCode(code string, arry []string) string {
	var symbol bytes.Buffer
	valus := ErrorMaps[code]
	for i, v := range arry {
		symbol.WriteString("{")
		symbol.WriteString(strconv.Itoa(i))
		symbol.WriteString("}")
		valus = strings.Replace(valus, symbol.String(), v, -1)
	}
	return valus
}

//记录code方式日志
func recordLogsWithCode(leve int, code string, arry []string, err error) {
	var noteFileLogs bytes.Buffer
	var noteConsoleLogs bytes.Buffer

	logs := marshalCode(code, arry)
	noteFile, noteConsole := noteHandle(leve, code)
	noteFileLogs.WriteString(noteFile)
	noteFileLogs.WriteString(logs)
	noteConsoleLogs.WriteString(noteConsole)
	noteConsoleLogs.WriteString(logs)
	//控制台打印
	writeConsoleLog(leve, noteConsoleLogs.String(), err)
	if nil != SysLogFile {
		//写入日志文件
		writeSysFileLog(leve, noteFileLogs.String(), err)
	}

}

//记录非code方式日志
// func recordLog(leve int, logs string, err error) {
// 	var fileLog bytes.Buffer
// 	var consoleLog bytes.Buffer
// 	noteFile, noteConsole := noteHandle(leve, "")
// 	consoleLog.WriteString(noteConsole)
// 	consoleLog.WriteString(logs)

// 	fileLog.WriteString(noteFile)
// 	fileLog.WriteString(logs)
// 	//控制台打印
// 	writeConsoleLog(leve, consoleLog.String(), err)

// 	if nil != SysLogFile {
// 		//写入日志文件
// 		writeSysFileLog(leve, fileLog.String(), err)
// 	}

// }

//记录非code 格式化方式日志
func recordLogf(leve int, logs string, err error, vars ...interface{}) {
	var fileLog bytes.Buffer
	var consoleLog bytes.Buffer
	noteFile, noteConsole := noteHandle(leve, "")
	consoleLog.WriteString(noteConsole)
	consoleLog.WriteString(logs)

	fileLog.WriteString(noteFile)
	fileLog.WriteString(logs)
	//控制台打印
	writeConsoleLogf(leve, consoleLog.String(), err, vars...)

	if nil != SysLogFile {
		//写入日志文件
		writeFileLogf(leve, fileLog.String(), err, vars...)
	}

}

//日志提示处理 [info] ［warning］[error]
func noteHandle(leve int, code string) (string, string) {
	var colour string
	var normalColour = "\033[0m" //黑色
	var codeDef string
	var noteFile bytes.Buffer
	var noteConsole bytes.Buffer
	switch leve {
	case 2: //info
		colour = "\033[36m" //青色
		codeDef = "I-ZLG-000"
	case 3: //warning
		colour = "\033[33m" //黄色
		codeDef = "W-ZLG-000"
	case 4: //error
		colour = "\033[31m" //红色
		codeDef = "E-ZLG-000"
	default:
		codeDef = "D-ZLG-000"
	}
	if 0 == len(code) {
		code = codeDef
	}
	noteFile.WriteString("[")
	noteFile.WriteString(code)
	noteFile.WriteString("]")

	noteConsole.WriteString("[")
	noteConsole.WriteString(colour)
	noteConsole.WriteString(code)
	noteConsole.WriteString(normalColour)
	noteConsole.WriteString("]")

	return noteFile.String(), noteConsole.String()
}

//writeLog func 日志输出处理
func writeSysFileLog(level int, logs string, errs error) {
	switch ZlogConfig.SysFile.Level {
	case "error": //日志级别为错误，只打印错误级别的日志
		if 4 == level {
			outputFileLogf(logs, errs)
		}
	case "warning":
		if 3 == level || 4 == level {
			outputFileLogf(logs, errs)
		}
	case "info":
		if 3 == level || 4 == level || 2 == level {
			outputFileLogf(logs, errs)
		}
	default:
		outputFileLogf(logs, errs)
	}

}

//writeLog func 日志输出处理
func writeConsoleLog(level int, logs string, errs error) {
	switch ZlogConfig.Console.Level {
	case "error": //日志级别为错误，只打印错误级别的日志
		if 4 == level {
			outputConsoleLogf(logs, errs)
		}
	case "warning":
		if 3 == level || 4 == level {
			outputConsoleLogf(logs, errs)
		}
	case "info":
		if 3 == level || 4 == level || 2 == level {
			outputConsoleLogf(logs, errs)
		}
	default:
		outputConsoleLogf(logs, errs)
	}

}

//writeLog func 日志输出处理
func writeFileLogf(level int, logs string, errs error, vars ...interface{}) {
	switch ZlogConfig.SysFile.Level {
	case "error": //日志级别为错误，只打印错误级别的日志
		if 4 == level {
			outputFileLogf(logs, errs, vars...)
		}
	case "warning":
		if 3 == level || 4 == level {
			outputFileLogf(logs, errs, vars...)
		}
	case "info":
		if 3 == level || 4 == level || 2 == level {
			outputFileLogf(logs, errs, vars...)
		}
	default:
		outputFileLogf(logs, errs, vars...)
	}

}

//writeLog func 日志输出处理
func writeConsoleLogf(level int, logs string, errs error, vars ...interface{}) {
	switch ZlogConfig.Console.Level {
	case "error": //日志级别为错误，只打印错误级别的日志
		if 4 == level {
			outputConsoleLogf(logs, errs, vars...)
		}
	case "warning":
		if 3 == level || 4 == level {
			outputConsoleLogf(logs, errs, vars...)
		}
	case "info":
		if 3 == level || 4 == level || 2 == level {
			outputConsoleLogf(logs, errs, vars...)
		}
	default:
		outputConsoleLogf(logs, errs, vars...)
	}

}

func outputConsoleLogf(logs string, errs error, vars ...interface{}) {
	var logsLn bytes.Buffer
	logsLn.WriteString(getFormatHeader(ZlogConfig.Console.Layout))
	logsLn.WriteString(logs)
	logsLn.WriteString("\n")

	_, err := fmt.Printf(logsLn.String(), vars...)
	if err != nil {
		fmt.Println(err)
	}
	if errs != nil {
		fmt.Println(errs.Error())
	}
}

// func outputConsoleLog(logs string, errs error) {

// 	fmt.Println(getFormatHeader(ZlogConfig.Console.Layout), logs)
// 	if errs != nil {
// 		fmt.Println(errs.Error())
// 	}
// }
func getFormatHeader(layout string) string {
	var header bytes.Buffer

	var times = time.Now()

	layouts := strings.Split(layout, "|")
	_, file, line, _ := runtime.Caller(5)

	for _, v := range layouts {
		switch v {
		case "date":
			header.WriteString(times.Format("2006/01/02"))
			header.WriteString(" ")
		case "time":
			header.WriteString(times.Format("15:04:05"))
			header.WriteString(" ")
		case "utc":
			header.WriteString(times.UTC().Format("2006/01/02 15:04:05"))
			header.WriteString(" ")
		case "fileName":

			header.WriteString(strings.Split(file, "/")[strings.Count(file, "/")])
			header.WriteString(":")
			header.WriteString(strconv.Itoa(line))
			header.WriteString(" ")
		case "fullpath":
			header.WriteString(file)
			header.WriteString(":")
			header.WriteString(strconv.Itoa(line))
		}

	}

	return header.String()
}
func outputFileLogf(logs string, errs error, vars ...interface{}) {
	var logsBuff bytes.Buffer
	logsBuff.WriteString(getFormatHeader(ZlogConfig.SysFile.Layout))
	logsBuff.WriteString(logs)
	logsBuff.WriteString("\n")
	fmt.Fprintf(SysLogFile, logsBuff.String(), vars...)

	if errs != nil {
		fmt.Fprintln(SysLogFile, errs.Error())
	}
}

// func outputFileLog(logs string, errs error) {
// 	fmt.Sprintln(logs)
// 	if errs != nil {
// 		log.Println(errs.Error())
// 	}
// }

func fileHandle() {

	if 0 != len(ZlogConfig.SysFile.Path) {
		SysLogFile = createFile("Sys", ZlogConfig.SysFile.Mode, ZlogConfig.SysFile.Path)
		if "dailyRolling" == ZlogConfig.SysFile.Mode {
			DailyRollingLogs(func() {

				nowTime := time.Now().Format("2006-01-02")
				_, err := os.Stat(ZlogConfig.SysFile.Path) //判断文件是否存在

				if err == nil {
					err := os.Rename(ZlogConfig.SysFile.Path, strings.Replace(ZlogConfig.SysFile.Path, ".", nowTime+".", -1))
					if err == nil {
						SysLogFile = createFile("Sys", ZlogConfig.SysFile.Mode, ZlogConfig.SysFile.Path)
					} else {
						Error("日志备份失败", err)
					}

				} else {
					SysLogFile = createFile("Sys", ZlogConfig.SysFile.Mode, ZlogConfig.SysFile.Path)
					Error("zlog 日志文件丢失,系统自动创建日志", nil)
				}

			})
		}
	}
	if 0 != len(ZlogConfig.AppFile.Path) {
		AppLogFile = createFile("App", ZlogConfig.AppFile.Mode, ZlogConfig.AppFile.Path)
		if "dailyRolling" == ZlogConfig.AppFile.Mode {
			DailyRollingLogs(func() {

				nowTime := time.Now().Format("2006-01-02")
				_, err := os.Stat(ZlogConfig.AppFile.Path) //判断文件是否存在

				if err == nil {
					err := os.Rename(ZlogConfig.AppFile.Path, strings.Replace(ZlogConfig.AppFile.Path, ".", nowTime+".", -1))
					if err == nil {
						AppLogFile = createFile("App", ZlogConfig.AppFile.Mode, ZlogConfig.AppFile.Path)
					} else {
						Error("zlog App 日志备份失败", err)
					}

				} else {
					AppLogFile = createFile("App", ZlogConfig.AppFile.Mode, ZlogConfig.AppFile.Path)
					Error("zlog App 日志文件丢失,系统自动创建日志", nil)
				}

			})
		}
	}
}

//AppOperateLog func
//UserName 操作人
//method  方法名
//action   方法类型
//newData  新数据
//oldData  旧数据
//explan   说明
func AppOperateLog(userName string, method string, action Action, newData interface{}, oldData interface{}, explan string) {
	var conslogsBuff bytes.Buffer
	var logsBuff bytes.Buffer
	logsBuff.WriteString("{\"OpTime\":\"")
	logsBuff.WriteString(time.Now().String())
	logsBuff.WriteString("\",\"UserName\":\"")
	logsBuff.WriteString(userName)
	logsBuff.WriteString("\",\"Action\":\"")
	switch action {
	case 1:
		logsBuff.WriteString("select\",")
	case 2:
		logsBuff.WriteString("update\",")
	case 3:
		logsBuff.WriteString("add\",")
	case 4:
		logsBuff.WriteString("delete\",")
	}
	if new, err := json.Marshal(newData); err == nil {
		logsBuff.WriteString("\"NewData\":")
		logsBuff.WriteString(string(new))
	} else {
		fmt.Fprintln(SysLogFile, err.Error())
	}
	if old, err := json.Marshal(oldData); err == nil {
		logsBuff.WriteString(",\"OldData\":")
		logsBuff.WriteString(string(old))
	} else {
		fmt.Fprintln(SysLogFile, err.Error())
	}
	logsBuff.WriteString(",\"Explan\":\"")
	logsBuff.WriteString(explan)
	logsBuff.WriteString("\"}")
	logsBuff.WriteString("\n")
	fmt.Fprintf(AppLogFile, logsBuff.String())
	conslogsBuff.WriteString("zlog App 操作日志:")
	conslogsBuff.Write(logsBuff.Bytes())
	fmt.Print(conslogsBuff.String())
}

// //AppOperateLog func产品日志
// func AppOperateLog(applog AppLoger) {
// 	var conslogsBuff bytes.Buffer
// 	var logsBuff bytes.Buffer
// 	applog.OpTime = time.Now()
// 	if b, err := json.Marshal(applog); err == nil {
// 		logsBuff.WriteString(string(b))
// 	} else {
// 		fmt.Fprintln(SysLogFile, err.Error())
// 	}
// 	logsBuff.WriteString("\n")
// 	fmt.Fprintf(AppLogFile, logsBuff.String())
// 	conslogsBuff.WriteString("zlog App 操作日志:")
// 	conslogsBuff.Write(logsBuff.Bytes())
// 	fmt.Print(conslogsBuff.String())

// }

func createFile(typ string, model string, path string) *os.File {
	var fileMode = os.O_APPEND
	var erros bytes.Buffer

	if "cover" == model { //如果为覆盖
		fileMode = os.O_TRUNC
	}

	//Infof("zlog %s 日志文件地址为%q", nil, typ, path)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|fileMode, 0666)
	if err != nil {
		erros.WriteString("zlog %s 日志文件 ")
		erros.WriteString(path)
		erros.WriteString(" 创建失败,无法使用此文件记录")
		Errorf(erros.String(), nil, typ)
		if "App" == typ {
			path = "zlogApp.log"
		} else {
			path = "zlogSys.log"
		}
		Infof("zlog 将在项目根路径下创建 %s 日志文件用于纪录日志信息", nil, typ)

		file, err = os.OpenFile(path, os.O_CREATE|os.O_RDWR|fileMode, 0666)

		if err != nil {
			Errorf("zlog 无法创建 %s 日志文件用于纪录日志信息", nil, typ)
			//os.Exit(1)
		}

		//Infof("zlog %s 日志文件地址为 %q", nil, typ, path)
	}
	return file
}

func init() {
	readXML()
	fileHandle()
	readCode()
	Info("zlog 启动成功，如有任何疑问请提交问题至http://git.pccqcpa.com.cn/components/zlog.git", nil)
}

func main() {
	ErrorWithCode("xxx", []string{"Z", "9"}, nil)
}
