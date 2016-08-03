package zlog

import (
	"bytes"
	"fmt"
	"log"
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

//StructLoger var 结构化日志
var StructLoger = new(log.Logger)

//ErrorWithCode func 带错误号码的输出错误日志
//code 配置文件定义的编号 err 为系统错误，如为自定义则为nil arry 翻译数组填充值
func ErrorWithCode(code string, arry []string, err error) {
	recordLogsWithCode(4, code, arry, err)
}

//Error func 输出错误日志
//logs 输出日志 err 系统原始错误日志
func Error(logs string, err error) {
	recordLog(4, logs, err)
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
	recordLog(3, logs, err)
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
	recordLog(2, logs, err)
}

//Infof func  格式化方式输出错误日志
//logs 输出日志 err 系统原始错误日志 vars替换的值
func Infof(logs string, err error, vars ...interface{}) {
	recordLogf(2, logs, err, vars...)
}

//Debug func 输出错误日志
//code 配置文件定义的编号 err 为系统错误，如为自定义则为nil arry 翻译数组填充值
func Debug(logs string, err error) {
	recordLog(1, logs, err)
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
	//写入日志文件
	writeFileLog(leve, noteFileLogs.String(), err)
}

//记录非code方式日志
func recordLog(leve int, logs string, err error) {
	var fileLog bytes.Buffer
	var consoleLog bytes.Buffer
	noteFile, noteConsole := noteHandle(leve, "")
	consoleLog.WriteString(noteConsole)
	consoleLog.WriteString(logs)

	fileLog.WriteString(noteFile)
	fileLog.WriteString(logs)
	//控制台打印
	writeConsoleLog(leve, consoleLog.String(), err)
	//写入日志文件
	writeFileLog(leve, fileLog.String(), err)

}

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
	//写入日志文件
	writeFileLogf(leve, fileLog.String(), err, vars...)

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
func writeFileLog(level int, logs string, errs error) {
	switch ZlogConfig.File.Level {
	case "error": //日志级别为错误，只打印错误级别的日志
		if 4 == level {
			outputFileLog(logs, errs)
		}
	case "warning":
		if 3 == level || 4 == level {
			outputFileLog(logs, errs)
		}
	case "info":
		if 3 == level || 4 == level || 2 == level {
			outputFileLog(logs, errs)
		}
	default:
		outputFileLog(logs, errs)
	}

}

//writeLog func 日志输出处理
func writeConsoleLog(level int, logs string, errs error) {
	switch ZlogConfig.Console.Level {
	case "error": //日志级别为错误，只打印错误级别的日志
		if 4 == level {
			outputConsoleLog(logs, errs)
		}
	case "warning":
		if 3 == level || 4 == level {
			outputConsoleLog(logs, errs)
		}
	case "info":
		if 3 == level || 4 == level || 2 == level {
			outputConsoleLog(logs, errs)
		}
	default:
		outputConsoleLog(logs, errs)
	}

}

//writeLog func 日志输出处理
func writeFileLogf(level int, logs string, errs error, vars ...interface{}) {
	switch ZlogConfig.File.Level {
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

	_, err := fmt.Printf(logs, vars...)
	if err == nil {
		fmt.Println(err)
	}
	if errs != nil {
		fmt.Println(errs.Error())
	}
}
func outputConsoleLog(logs string, errs error) {

	fmt.Println(getFormatHeader(), logs)
	if errs != nil {
		fmt.Println(errs.Error())
	}
}
func getFormatHeader() string {
	var header bytes.Buffer

	var times = time.Now()

	layouts := strings.Split(ZlogConfig.Console.Layout, "|")
	_, file, line, _ := runtime.Caller(2)

	for _, v := range layouts {
		switch v {
		case "date":
			header.WriteString(times.Format("2006/01/02"))
		case "time":
			header.WriteString(" ")
			header.WriteString(times.Format("15:04:05"))
		case "utc":
			header.WriteString(times.UTC().Format("2006/01/02 15:04:05"))
		case "fileName":
			header.WriteString(" ")
			header.WriteString(strings.Split(file, "/")[strings.Count(file, "/")])
		case "fullpath":
			header.WriteString(file)

		}

	}
	header.WriteString(":")
	header.WriteString(strconv.Itoa(line))
	return header.String()
}
func outputFileLogf(logs string, errs error, vars ...interface{}) {
	log.Printf(logs, vars...)
	if errs != nil {
		log.Println(errs.Error())
	}
}
func outputFileLog(logs string, errs error) {
	log.Output(0, logs)
	if errs != nil {
		log.Println(errs.Error())
	}
}

func fileHandle() {

	if 0 != len(ZlogConfig.File.Path) {
		createLogFile()
		if "dailyRolling" == ZlogConfig.File.Mode {
			DailyRollingLogs(func() {

				nowTime := time.Now().Format("2006-01-02")
				_, err := os.Stat(ZlogConfig.File.Path) //判断文件是否存在

				if err == nil {
					err := os.Rename(ZlogConfig.File.Path, strings.Replace(ZlogConfig.File.Path, ".", nowTime+".", -1))
					if err == nil {
						createLogFile()
					} else {
						Error("日志备份失败", err)
					}

				} else {
					createLogFile()
					Error("日志文件丢失,系统自动创建日志", nil)
				}

			})
		}

	}
}

func createLogFile() {
	var logFlag = log.LstdFlags
	var fileMode = os.O_APPEND
	var erros bytes.Buffer
	if "cover" == ZlogConfig.File.Mode { //如果为覆盖
		fileMode = os.O_TRUNC
	}
	logfile, err := os.OpenFile(ZlogConfig.File.Path, os.O_CREATE|os.O_RDWR|fileMode, 0666)
	if err != nil {
		erros.WriteString(" zlog日志文件 ")
		erros.WriteString(logfile.Name())
		erros.WriteString(" 创建失败,无法使用文件记录")
		Error(erros.String(), err)
		//os.Exit(1)
	}
	log.SetOutput(logfile)
	layouts := strings.Split(ZlogConfig.File.Layout, "|")
	for _, v := range layouts {
		switch v {
		case "date":
			logFlag = log.Ldate | logFlag
		case "time":
			logFlag = log.Ltime | logFlag
		case "utc":
			logFlag = log.LUTC | logFlag
		case "fileName":
			logFlag = log.Lshortfile | logFlag
		case "fullpath":
			logFlag = log.Llongfile | logFlag

		}

	}

	log.SetFlags(logFlag)

}

func init() {
	readXML()
	fileHandle()
	readCode()
	Info("zlog 启动成功，如有任何疑问请提交问题至http://git.pccqcpa.com.cn/components/zlog.git", nil)
}
