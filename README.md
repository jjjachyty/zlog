golang log 插件，分离console和文件纪录，支持文件覆盖，追加以及按日滚动纪录

1、 在项目根目录下面创建conf/zlog/目录
2、 配置zlog.xml
3、日志提示信息在code.ini按照系统格式填写,格式为日志级别-模块-编号，如错误级别的zlog日志为 E-ZLG-000,模块简写和编号统一为三位数
4、zlog提供Debug、Info、Warning、Error日志级别，分别对应
    Debug(logs string, err error)
    Info(logs string, err error)
    Warning(logs string, err error)
    Error(logs string, err error)
    logs为输出日志信息,err为系统原有的错误信息
5、如果使用code.ini配置文件输出错误信息，则选择对应WithCode方法
    如:ErrorWithCode(code string,arry []string,err error)
    code 为code.ini配置文件的编码，如E-ZLG-001
    arry 为替换占位符的数据，如{1} 替换成 机构
    err  为原始系统错误信息
