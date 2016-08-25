package zlog

import (
	"encoding/xml"
	"io/ioutil"
)

//LogXMLSRT struct
type LogXMLSRT struct {
	Console CosXMLSRT     `xml:"Cos"`
	SysFile SysFileXMLSRT `xml:"sysFile"`
	AppFile AppFileXMLSRT `xml:"appFile"`
}

//SysFileXMLSRT struct
type SysFileXMLSRT struct {
	Layout string `xml:"layout"`
	Path   string `xml:"path"`
	Mode   string `xml:"mode"`
	Level  string `xml:"level"`
}

//AppFileXMLSRT ileXMLSRT struct
type AppFileXMLSRT struct {
	Path string `xml:"path"`
	Mode string `xml:"mode"`
}

//CosXMLSRT struct
type CosXMLSRT struct {
	Layout string `xml:"layout"`
	Level  string `xml:"level"`
}

//ZlogConfig var
var ZlogConfig = &LogXMLSRT{CosXMLSRT{"date|time|SysileName", "debug"}, SysFileXMLSRT{"date|time|SysileName", "zlogSys.log", "cover", "debug"}, AppFileXMLSRT{"zlogApp.log", "cover"}}

//readXML读取XML配置
func readXML() {

	content, err := ioutil.ReadFile("conf/zlog/zlog.xml")
	if err == nil {
		err = xml.Unmarshal(content, &ZlogConfig)
		Debug("zlog.xml 文件解析成功", nil)
		Infof("zlog Cos 日志输出格式为:%s", nil, ZlogConfig.Console.Layout)
		Infof("zlog Sys 日志输出格式为:%s", nil, ZlogConfig.SysFile.Layout)
		Infof("zlog Sys 日志输出地址为:%s", nil, ZlogConfig.SysFile.Path)

		Infof("zlog Sys 日志输出模式为:%s", nil, ZlogConfig.SysFile.Mode)
		Infof("zlog Sys 日志输出级别为:%s", nil, ZlogConfig.SysFile.Level)
		Infof("zlog App 日志输出地址为:%s", nil, ZlogConfig.AppFile.Path)
		Infof("zlog App 日志输出模式为:%s", nil, ZlogConfig.AppFile.Mode)
		if err != nil {
			Error("zlog.xml文件解析失败，请按照说明检查配置项！", err)
		}
	} else { //启用默认配置
		Warning("无法在项目目录中发现conf/zlog/zlog.xml文件,zlog将启动默认配置。", nil)
		Infof("zlog Cos 日志输出格式为:%s", nil, ZlogConfig.Console.Layout)
		Infof("zlog Sys 日志输出格式为:%s", nil, ZlogConfig.SysFile.Layout)
		Infof("zlog Sys 日志输出地址为:%s", nil, ZlogConfig.SysFile.Path)
		Infof("zlog Sys 日志输出模式为:%s", nil, ZlogConfig.SysFile.Mode)
		Infof("zlog Sys 日志输出级别为:%s", nil, ZlogConfig.SysFile.Level)
		Infof("zlog App 日志输出地址为:%s", nil, ZlogConfig.AppFile.Path)
		Infof("zlog App 日志输出模式为:%s", nil, ZlogConfig.AppFile.Mode)
	}
}
