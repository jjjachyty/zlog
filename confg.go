package zlog

import (
	"encoding/xml"
	"io/ioutil"
)

//LogXMLSRT struct
type LogXMLSRT struct {
	Console ConsoleXMLSRT `xml:"console"`
	File    FileXMLSRT    `xml:"file"`
}

//FileXMLSRT struct
type FileXMLSRT struct {
	Layout string `xml:"layout"`
	Path   string `xml:"path"`
	Mode   string `xml:"mode"`
	Level  string `xml:"level"`
}

//ConsoleXMLSRT struct
type ConsoleXMLSRT struct {
	Layout string `xml:"layout"`
	Level  string `xml:"level"`
}

//ZlogConfig var
var ZlogConfig = LogXMLSRT{ConsoleXMLSRT{"date|time|fileName", "debug"}, FileXMLSRT{"date|time|fileName", "zlog.log", "cover", "debug"}}

//readXML读取XML配置
func readXML() {

	content, err := ioutil.ReadFile("conf/zlog/zlog.xml")
	if err == nil {
		err = xml.Unmarshal(content, &ZlogConfig)
		Infof(" zlog－console 日志输出格式为:%s", nil, ZlogConfig.Console.Layout)
		Infof(" zlog－file 日志输出格式为:%s", nil, ZlogConfig.File.Layout)
		Infof(" zlog－file 日志输出地址为:%s", nil, ZlogConfig.File.Path)
		Infof(" zlog－file 日志输出模式为:%s", nil, ZlogConfig.File.Mode)
		Infof(" zlog－file 日志输出级别为:%s", nil, ZlogConfig.File.Level)
		if err != nil {
			Error("zlog.xml文件解析失败，请按照说明检查配置项！", err)
		}
	} else { //启用默认配置
		fileHandle()
		Warning(" 无法在项目目录中发现conf/zlog/zlog.xml文件,zlog将启动默认配置。", err)
		Infof(" zlog－console 日志输出格式为:%s", nil, ZlogConfig.Console.Layout)
		Infof(" zlog－file 日志输出格式为:%s", nil, ZlogConfig.File.Layout)
		Infof(" zlog－file 日志输出地址为:%s", nil, ZlogConfig.File.Path)
		Infof(" zlog－file 日志输出模式为:%s", nil, ZlogConfig.File.Mode)
		Infof(" zlog－file 日志输出级别为:%s", nil, ZlogConfig.File.Level)

	}
}
