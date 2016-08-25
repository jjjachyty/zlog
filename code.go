package zlog

import (
	"runtime"

	"github.com/larspensjo/config"
)

//ErrorMaps var 错误代码集合
var ErrorMaps = make(map[string]string)

func readCode() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg, err := config.ReadDefault("conf/zlog/code.ini")
	if err == nil {
		sections := cfg.Sections()
		for _, section := range sections {
			sectionOptions, err := cfg.SectionOptions(section)
			if err == nil {
				for _, sectionOption := range sectionOptions {
					options, err := cfg.String(section, sectionOption)
					if err == nil {
						ErrorMaps[sectionOption] = options
					} else {
						Errorf("code.ini解析失败，请确认options %s 已配置", err, options)
					}
				}
			} else {
				Errorf("code.ini解析失败，请确认section %s 已配置", err, sectionOptions)
			}
		}

	} else {
		Warning("无法在项目目录中发现conf/zlog/code.ini文件,..WithCode方法将打印原始信息。", nil)
	}
}
