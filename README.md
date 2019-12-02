# gutils
Go语言常用工具。

api
    提供接口CodeModeDTO返回对象格式
    
bean
    提供bean注册、获取、自动注入
    
cache
    提供数据缓存、获取、定时销毁功能
    
    * 依赖于task
    
config
    提供读取Json配置文件信息
    
    * 依赖于file
    
controller
    提供controller缓存、注册功能。提供了简易API控制器
    
    * 依赖于github.com/kataras/iris
    * 依赖于github.com/Lyo-Shur/gorm
    * 依赖于api
    
file
    提供txt类文件读取方法
    
form
    提供表单字段、上传文件提取。保存文件到本地功能
    
    * 依赖于github.com/kataras/iris
    
sso
    提供根据信息生成票据、使用票据反查信息功能
    
    * 依赖于github.com/satori/go.uuid
    * 依赖于cache
    
task
    提供间隔性循环执行方法功能
    
validator
    提供结构体字段类型格式校验功能
