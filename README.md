# reptilelibs

## Logger Test
```
go test -bench=. -run=. logger.go logger_test.go
```

## Mysql Test
```
go test -bench=. -run=.  mysql.go mysql_test.go 
```

## Redis Test
```
go test -bench=. -run=. redis.go redis_test.go
```

## Table name Test
```
go test -bench=. -run=. table_name.go table_name_test.go
```

## Util Test
```
go test -bench=. -run=. util.go util_test.go
```

## hbase
```
go test -bench=. -run=. hbase.go hbase_test.go
```

## kafka
```
go test -bench=. -run=.  config.go kafka.go kafka_test.go util.go httpclient.go
```

## elasticsearch
```
go test -bench=. -run=. httpclient.go util.go elastic.go elastic_test.go 
```

# 可变参数的问题
```
func (this *SqlDB) ExecAction(statement string,callback func(SqlResult,error),args ...interface{}) (err error) {
	db := this.GetDB()
	err = db.Exec(statement,args).Error

    // ...
}
```
上面的调用Exec的可变参数会出现错误，在执行的时候会出现一个空的参数，在调用的使用需要如下调用:
```
func (this *SqlDB) ExecAction(statement string,callback func(SqlResult,error),args ...interface{}) (err error) {
	db := this.GetDB()
	err = db.Exec(statement,args...).Error

    // ...
}
```


# 接口
## NeedLogin
```
/**
* 权限中间件，主要是确认是否登陆成功
*
* @param callback 	主要是在每个服务中将用户ID传递出去
* @param timeout 	token的过期时间
*/ 
func NeedLogin(callback func(uid string) bool,timeout int64) gin.HandlerFunc{}
```

## LoadConfigFromFile
装载配置文件
```
/**
* 装载配置文件
*
* @param file 配置文件的路径
*/
func LoadConfigFromFile(file string)
```

## 编码Token
```
/**
* 编码token
* 
* @param src 	产生token的用户信息
* @param key 	产生token的秘钥
*
* @return token 产生以后的token信息
* @return err 	产生token是否出错
*/
func EnCookie(src TokenInfo, key string) (token string,err error)
```

## 解码Token
```
/**
* 解码token
*
* @param  src 	token信息
* @param  key 	解码token的秘钥
*
* @return token  解码后的token用户信息
* @return err    解码状态
*/
func DeCookie(src string, key string) (token TokenInfo,err error) 
```

## 装载错误码文件
```
/**
* 装载错误码文件
*
* @param file 	错误码文件路径
*/
func LoadCodeFromFile(file string) 
```

## 获取错误码
通过编码从错误码配置文件中获取错误码
```
func ResponseCode(enum string) float64
```

## 获取错误信息
通过编码从错误码配置文件中获取错误信息
```
func ResponseMsg(enum string) string
```

## 创建http操作对象
创建用于发起http请求的对象
```
func NewHttpClient() *HttpClient
```

## 创建http请求对象
```
// 创建一个请求
// @param url 		请求的地址
// @param data 		请求的数据，主要能发送的数据类型是：string,[]byte,map
// @return err
func (this *HttpClient) NewRequest(url string, data interface{}) (err error)
```

## 设置http请求头
```
// 设置请求的头部
// @param key
// @param value
func (this *HttpClient) Header(key, value string)
```

## 发起GET请求
发起HTTP的GET请求
```
// 发起Http的get请求
// @return data,err
func (this *HttpClient) Get() (r []byte, err error)

example:
	request := NewHttpClient()
	err := request.NewRequest(url, nil)
	if nil != err {
		fmt.Println("NewRequestFail, ", err)

		return
	}

	response, err := request.Get()
	if nil != err {
		fmt.Println("Http Get Fail, ", err)

		return
	}
```

## 发起PUT请求
```
// 发起Http的PUT请求
// @return data,err
func (this *HttpClient) Put() (r []byte, err error)

example:
	body := ...

	request := NewHttpClient()
	err := request.NewRequest(url, body)
	if nil != err {
		fmt.Println("NewRequest Fail, ", err)

		return
	}

	request.Header("Content-Type", "application/json")
	response, err := request.Put()
	if nil != err {
		fmt.Println("Http Put Fail, ", err)

		return
	}
```

## 发起DELETE请求
```
// 发起Http的delete请求
// @return data,err
func (this *HttpClient) Delete() (r []byte, err error)

example:
	request := NewHttpClient()
	err := request.NewRequest(url, nil)
	if nil != err {
		fmt.Println("NewRequest Fail, ", err)

		return
	}

	response, err := request.Delete()
	if nil != err {
		fmt.Println("Http Delete Fail, ", err)

		return
	}
```

## 发起POST请求
发起HTTP的POST请求
```
// 发起Http的post请求
// @return data,err
func (this *HttpClient) Post() (r []byte, err error) 

example:
	body := ...

	request := NewHttpClient()
	err := request.NewRequest(url, body)
	if nil != err {
		fmt.Println("NewRequest Fail, ", err)

		return
	}

	request.Header("Content-Type", "application/json")
	response, err := request.Post()
	if nil != err {
		fmt.Println("Http Post Fail, ", err)

		return
	}
```

## 发起base64 Post请求
发起HTTP的base64 POST请求
```
/**
* 发起BASE64 POST请求
*
* @param urlAddress 	请求的url
* @param data 			请求的请求体
*
* @return r 			请求成功以后的返回数据
* @return err 			请求是否成功
*/
func PostBase64(urlAddress string, data interface{}) (r []byte, err error) 
```

## 创建sql数据库操作实例
```
/**
* 创建sql数据库操作实例
*
* @param 	host 	 数据库主机
* @param 	port 	 数据库端口
* @param 	dbname 	 数据库名字
* @param 	username 数据库用户名
* @param 	password 数据库登录密码
*
* @return 	SqlDB	数据库操作实例
*/
func NewSqlDB(host,port,dbname,username,password string) *SqlDB
```

## 传教clickhouse数据库实例
```
// 创建clickhouse实例
//
// @param host 	clickhouse数据库连接地址
//
func NewClickHouse(host string) *ClickHouse
```

## 数据库查询接口
```
/**
* 查询接口
*
* @param statement 查询语句
* @param callback 查询回调函数   参数： 查询到的值 和  查询状态
* @param args 查询行的临时存储变量【主要用于查询的列,目前仅仅支持string和float64,也就是要查询的所有列的类型
*
* @return count 	查询到的行数
* @return err 		查询状态
*/
func (this *SqlDB) ExecRows(statement string,callback func([][]interface{}, error),args ...interface{}) (count int,err error){
```

## 数据库的动作执行（插入、更新、删除）
```
/**
* 执行数据操作动作，主要是插入数据和更新数据
*
* @param statement 	动作的语句
* @param callback 执行的回调函数
* @param args 动作的参数
*
* @return err 	执行的状态
*/
func (this *SqlDB) ExecAction(statement string,callback func(error),args ...interface{}) (err error) 
```

## 设置数据库的模式
```
/**
* 设置数据库为调试模式
*/
func (this *SqlDB) LogMode(mode bool)
```

## 设置数据库空闲连接数大小
```
/**
* 设置数据库空闲连接数大小
*/
func (this *SqlDB) SetMaxIdleConns(count int)
```

## 最大打开连接数
```
/**
* 最大打开连接数
*/
func (this *SqlDB) SetMaxOpenConns(count int)
```

## 创建redis实例
```
/**
* 创建新的redis实例
*
* @param host 	redis主机地址
* @param port 	redis主机端口
* @param passwd redis登录密码
* @param db 	redis数据库
* @param timeoutredis超时时间
*
* @return Redis 返回的redis操作实例
*/
func NewRedis(host,port,passwd string,db,timeout int) *Redis
```

## 打开redis连接
```
func (this *Redis) Open()
```

## 保存数据库到Redis中
```
/**
* 存储在线用户缓存数据
*
* @param key 	缓存数据key
* @param value	缓存的数据
* @param timeout 设置超时时间
*
* @return err 	错误信息
*/
func (this *Redis)SetValue(key string,value interface{},timeout int) (err error)
```

## 从redis中获取值
```
/**
* 获取在线用户缓存数据
*
* @param key 	要获取缓存数据的key
* @param isset 是否重新更新缓存时间
* @param timeout 重新设置更新的时间
*
* @return (value,err)
*/
func (this *Redis)GetValue(key string,isset bool,timeout int) (value interface{},err error)
```

## 获取Redis所有的key
```
/**
* 获取在线缓存的所有keys
*
* @return (keys[],err)
*/
func (this *Redis)GetKeys() (keys []string,err error)
```

## 删除Redis中的值
```
/**
* 删除在线的用户缓存
*
* @param key 	要删除缓存的key
*/
func (this *Redis)DelValue(key string) 
```

## 删除Redis中所有的值
```
func (this *Redis)FlushAll()
```

## 创建MongoDB实例
```
/**
* 创建新的mongodb实例
*
* @param url 	mongodb的连接地址
*
* @return *MongoDB
*/
func NewMongoDB(url string) *MongoDB
```

## 打开MongoDB操作会话
```
/**
* 打开mongodb连接
*/
func (this *MongoDB) Open() (err error)
```

## 关闭mongoDB操作
```
/**
* 关闭集合
*/
func (this *MongoDB) Close() 
```

## 检查是否打开
```
/**
* 检查是否打开了数据库
*/
func (this *MongoDB) IsOpen() bool
```

## 插入数据
```
/**
* 插入数据到数据中
*
* @param db 		数据库
* @param collect	集合
* @param data 		插入的数据
*
* @return err 
*/
func (this *Session) Insert(db,collect string,data interface{}) (err error)
```

## 查询数据
```
/**
* 根据条件查询数据
*
* @param db 		数据库
* @param collect	集合
* @param param 		查询的参数
*
* @return query		返回索引对象
* @return err 		查询状态
*/
func (this *Mongo) Find(db,collect string,param interface{}) (query MongoQuery, err error)
```
注： query对象提供查询方法

## 更新或插入数据
```
/**
* 根据selector更新数据或插入数据
*
* @param db 		数据库
* @param collect	集合
* @param 	selector 	更新的条件
* @param 	update 		更新的数据
*
* @return err
*/
func (this *Mongo) UpdateOrInsert(db,collect string,selector interface{}, update interface{}) (err error)
```

## HBase响应结构
```
type HBaseResult struct {
	RowKey 		string 	// rowkey
	Family 		string  // 族名
	Qualifier 	string	// 列名
	Value 		[]byte  // 值
}
```

## 创建HBase实例
```
func NewHBase() *HBase ;
```

## 打开HBase客户操作
```
/**
* 打开hbase普通用户端
*
* @param host 	集群的主机地址列表[192.168.0.119:2181,192.169.0.120:2181]
*/
func (this *HBase) OpenClient(host string);
```

## 打开HBase管理操作
```
/**
* 打开hbase管理员端
*
* @param host 	集群的主机列表[192.168.0.119:2181,192.169.0.120:2181]
*/
func (this *HBase) OpenAdminClient(host string)
```

## 查寻HBase中一整行数据
```
/**
* 通过rowkey查询一整行数据
*
* @param table_name 	表名
* @param rowkey 		表的rowkey
*
* @return []HBaseResult,err
*/
func (this *HBase) GetRow(table_name string,rowkey string) (r []HBaseResult,err error)
```

## 在HBase中添加或更新数据
```
/**
* 插入或者更新一个值
*
* @param table_name 	表名
* @param value 			需要插入或者更新的值
* 
* @return err
*/
func (this *HBase) PutCell(table_name string,value HBaseResult) (err error)
```

## 从HBase中获取某写列的数据
```

/**
* 根据列名获取值
*
* @param  table_name 	表名
* @param  rowkey 		rowkey
* @param  family 		族名
* @param  columns 		列明
*
* @return map[string]HBaseResult,err
*/
func (this *HBase) GetCell(table_name string,rowkey string,family string,columns []string) (r map[string]HBaseResult,err error)
```

## 从HBase中删除族或列
```
///
// 删除Cell操作
//
// @param table_name 	表名
// @param value 			删除的值
//	注： 如果删除的是族，只需要将value中的Qualifier设置为""
// 	注： 如果删除的是某个列，只需要将value中的value设置为nil
// @return err	
//
func (this *HBase) DelCell(table_name string,value HBaseResult) (err error)
```


## 请求的响应函数
### 成功
```
/**
* @param data 	响应的数据
*/
func HandleOK(c *gin.Context, data interface{})

响应格式是：
	{
		"code": 200,
		"message": "success",
		"data": data
	}
```
### 失败
```
/**
*
* @param code 	响应的状态码
* @param msg 	响应的错误信息
*/
func HandleErr(c *gin.Context, code float64, msg string, err error)

响应格式是：
	{
		"code": code,
		"message", msg,
		"error": err,
	}
```

## 装载数据库表文件
所有的数据库表通过配置文件设置
```
/**
* @param 	file 	配置文件路径
*/
func LoadTableFromFile(file string)
```

## 获取数据表名
```
/**
* 获取表的key
*/
func TableName(key string) string
```

## 产生错误对象
```
func GenError(msg string) error
```

## 获取时间戳
```
func Timestamp() int64 
```

## MD5加密
```
func MD5(str string) string
```

## 产生32 UID
```
func GenUUID32() (string, error)
```

## 得到年月日
```
func GetNowYmd() string
```

## 得到当前日期和时间
```
func GetNowYmdHms() string
```

## 格式化日期格式
也就是将日期中的T和Z替换
```
func FormatDateString(str string) string
```

## 得到n为随机数
```
func GetRandomNumberString(n int) string
```

## 请求验证码
```
// 发送短信验证码
// url 		请求的utl接口
// appkey 	请求发送短信的appkey
// template 请求发送短信的template
// secret 	请求发送短信的secret
// mobile 	接收短信的号码
// 
// code,err
func SendSmsCode(url,appkey,template,secret, mobile string) (code string, err error)
```

## 如果日志需要输出到消息队列，进行初始化
```
func Log_init_rabbitmq(host,port,username,passwd,name string)
```

## 控制日志是否输出到控制台
```
/**
* 控制日志是否输出到控制台
*/
func Log_IsConsole(is bool) 
```

## 控制日志输出到消息队列中
```
func Log_IsRabbitmq(is bool)
```

## 设置日志是否是调试模式
```
func Log_LogMode(is bool)
```

## 输出日志信息
```
func Println(msg ...interface{})

func Error(msg ...interface{})

func Warning(msg ...interface{})
```

## 创建kafka实例
```
// 创建kafka实例
// @param c 配置结构
func NewKafka(c *KafkaConfig) *KafKa
```

## kafka消息的类型
```
// 消息的类型
type MsgType int32
const (
	MsgType_Login MsgType = 0
)
```

## 创建生产者
```
/ 创建生产者
// @param c 配置结构
func newKafkaPub(c *KafkaConfig) kafka.SyncProducer
```

## 创建消费者
```
// 创建消费者
// @param c 配置结构
func newKafkaSub(c *KafkaConfig) *cluster.Consumer
```

## 生产消息
```
// 生产消息
// @param c
// @param key 消息的key
// @param m_type 消息类型
// @param msg 	消息
func (d *KafKa) Product(c context.Context, key string, m_type MsgType, msg []byte) (err error) 
```

## 循环消费消息
```
// 循环消费消息
// @param errf 	收到错误回调函数
// @param notif 收到通知回调函数
// @param msgf  收到消息毁掉函数
func (j *KafKa) Consume(errf func(error), notif func(interface{}), msgf func(*KafkaMsg)) 
```

## 关闭kafka消费
```
// 关闭消息队列监听
func (j *KafKa) Close() 
```

## 创建elasticsearch操作对象
```
func NewElasticSearch(url string) *ElasticSearch
```

## 创建索引
```
// 创建索引
// @param index 	索引名
// @param  attr		索引参数[包含索引属性以及映射等信息]
// @body 	json
func (this *ElasticSearch) CreateIndex(index string, attr interface{}) (r []byte, err error)

example:
	body := `
		{
			"settings": {
				"number_of_shards": 3,
				"number_of_replicas": 2
			},
			"mappings": {
				"_source": {
					"enabled": true
				},
				"dynamic": true,
				"properties": {
					"message": {
						"type": "text",
						"analyzer": "standard"
					},
					"uid": {
						"type": "text"
					}
				}
				
			}
		}
	`

	elastis := NewElasticSearch("http://127.0.0.1:9200")
	r, err := elastis.CreateIndex("test", body)
	if nil != err {
		fmt.Println("TestCreateIndex error: ", err)
		return
	}

	fmt.Println("CreateIndex success! data: ", string(r))
	return
```

## 更新索引
```
// 更新索引，主要用于修改索引的属性以及添加索引的映射
// @param index 	索引名
// @param attr 		索引属性[包含索引属性以及映射等信息]
// @body json
// 注： 有些信息是不能被修改的，如:number_of_shards
func (this *ElasticSearch) UpdateIndex(index string, field string, attr interface{}) (r []byte, err error)

example:
	body := `
		{
			"number_of_replicas": 2
		}
	`

	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.UpdateIndex("test", "_settings", body)
	if nil != err {
		fmt.Println("UpdateIndex error: ", err)
		return
	}

	fmt.Println("UpdateIndex success! data: ", string(response))
	return
```

## 打开索引
```
// 打开索引
// @param index
func (this *ElasticSearch) OpenIndex(index string) (r []byte, err error)

example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.OpenIndex("test")
	if nil != err {
		fmt.Println("TestOpenIndex error: ", err)
		return
	}

	fmt.Println("TestOpenIndex success! data: ", string(response))
	return
```

## 关闭索引
```
// 关闭索引
// @param index
func (this *ElasticSearch) CloseIndex(index string) (r []byte, err error)


example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.CloseIndex("test")
	if nil != err {
		fmt.Println("TestCloseIndex error: ", err)
		return
	}

	fmt.Println("TestCloseIndex success! data: ", string(response))
	return
```

## 删除索引
```
// 删除索引
// @param index 索引名
func (this *ElasticSearch) DeleteIndex(index string) (r []byte, err error)

example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.DeleteIndex("test")
	if nil != err {
		fmt.Println("TestDeleteIndex error: ", err)
		return
	}

	fmt.Println("TestDeleteIndex success! data: ", string(response))
	return
```

## 设置索引别名
```
// 添加索引别名
// @param 	index 	索引名
// @param 	alias 	别名
func (this *ElasticSearch) AddIndexAlias(index string, alias string) (r []byte, err error)

example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.AddIndexAlias("test", "myalias")
	if nil != err {
		fmt.Println("TestAddIndexAlias error: ", err)
		return
	}

	fmt.Println("TestAddIndexAlias success! data: ", string(response))
	return
```


## 获取索引信息
```
// 获取索引信息
// @param 	index 	索引名
func (this *ElasticSearch) GetIndex(index string) (r []byte, err error) 

example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.GetIndex("test")
	if nil != err {
		fmt.Println("TestGetIndex error: ", err)
		return
	}

	fmt.Println("TestGetIndex success! data: ", string(response))
	return
```

## 清除索引缓存
```
// 清除索引缓存
// @param 	index 	索引名
func (this *ElasticSearch) ClearCache(index string) (r []byte, err error)

example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.ClearCache("test")
	if nil != err {
		fmt.Println("TestClearCache error: ", err)
		return
	}

	fmt.Println("TestClearCache success! data: ", string(response))
	return
```

## 刷新索引
```
// 刷新索引
// @param  	index 	索引名
func (this *ElasticSearch) RefreshIndex(index string) (r []byte, err error)

example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.RefreshIndex("test")
	if nil != err {
		fmt.Println("TestRefreshIndex error: ", err)
		return
	}

	fmt.Println("TestRefreshIndex success! data: ", string(response))
	return
```

## 冲洗索引
```
// 冲洗索引
// 索引主要是通过执行冲洗将数据保存到索引存储并清除内部事物日志
// @param  	index 	索引名
func (this *ElasticSearch) FlushIndex(index string) (r []byte, err error)

example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.FlushIndex("test")
	if nil != err {
		fmt.Println("TestFlushIndex error: ", err)
		return
	}

	fmt.Println("TestFlushIndex success! data: ", string(response))
	return
```


## 添加文档
```
// 添加文档{添加数据}
// 注： 数据格式和索引中的mapping格式样的
// @param  	index 		索引名
// @param 	doc[json] 	文档数据
func (this *ElasticSearch) PostDoc(index string, doc string) (r []byte, err error) 

example:
	doc := `
		{
			"log": "This is a test"
		}
	`

	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.PostDoc("test", doc)
	if nil != err {
		fmt.Println("TestPostDoc error: ", err)
		return
	}

	fmt.Println("TestPostDoc success! data: ", string(response))
	return
```


## 获取文档
```
// 通过id获取文档
// @param 	index 	索引名
// @param 	id 		文档id
func (this *ElasticSearch) GetDoc(index string, id string) (r []byte, err error) 


example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.GetDoc("test", "1")
	if nil != err {
		fmt.Println("TestGetDoc error: ", err)
		return
	}

	fmt.Println("TestGetDoc success! data: ", string(response))
	return
```


## 获取文档的源数据
```
// 仅仅是获取文档的内容
// @param 	index 	索引名
func (this *ElasticSearch) DocSource(index string, id string) (r []byte, err error)


example:
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.DocSource("index", "ktpDYXMBskf7Xeau3ZnL")
	if nil != err {
		fmt.Println("TestDocSource error: ", err)
		return
	}

	fmt.Println("TestDocSource success! data: ", string(response))
	return
```


## 通过query进行搜索
```
// 通过url中添加搜索信息的方式搜索
// @param 	index 	索引名
// @param 	query 	索引参数
func (this *ElasticSearch) SearchByQuery(index string, query string) (r []byte, err error) 

```

## 通过body进行搜索
```
// 通过body进行搜索
// @param index		索引名
// @param body 		索引参数体
// 注： 可以进行任意的查询
//		包含全文查询、短语、前缀、多字段等等
func (this *ElasticSearch) SearchByBody(index string, body string) (r []byte, err error)


example:
	body := `
		{
			"from": 0,
			"size": 20,
			"sort": [
				{
					"message": {
						"order": "asc"
					}
				}
			],
			"query": {
				"term": {
					"data.Timu": "搜索"
				}
			},
			"highlight": {
				"fields": {
				  "message": {} 
				}
			  },
			"_source": true
		}
	`
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.SearchByBody("topic", body)
	if nil != err {
		fmt.Println("TestSearchByBody error: ", err)
		return
	}

	fmt.Println("TestSearchByBody success! data: ", string(response))
	return
```


## 验证搜索条件是否可执行
```
// 验证查询语句是否正确
// @param 	index 	索引名字
// @param 	body 	查询体
// 注： 查询体只能包含query域
// 正确返回:{"_shards":{"total":1,"successful":1,"failed":0},"valid":true}
// 不能查询返回: {"valid": false}
func (this *ElasticSearch) SearchValidate(index string, body string) (r []byte, err error)


example:
	body := `
		{
			"query": {
				"term": {
					"message": "this"
				}
			}
		}
	`
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.SearchValidate("index", body)
	if nil != err {
		fmt.Println("TestSearchValidate error: ", err)
		return
	}

	fmt.Println("TestSearchValidate success! data: ", string(response))
	return
```

## 获取数据的条数
```
// 查询满足条件的数据条数
// @param index 	索引名
// @param body 		查询参数体
// 注： 查询体仅仅包含查询参数中的query域
func (this *ElasticSearch) CountByBody(index string, body string) (r []byte, err error) 


example：
	body := `
		{
			"query": {
				"multi_match": {
					"query": "资源管理 目录索引",
					"type": "most_fields",
					"fields": ["data.Timu","data.Daan.Name"]
				}
			}
		}
	`
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.CountByBody("topic", body)
	if nil != err {
		fmt.Println("TestCountByBody error: ", err)
		return
	}

	fmt.Println("TestCountByBody success! data: ", string(response))
	return
```