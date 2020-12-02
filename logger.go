package youyoulibs

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"runtime"
	"time"
)

type LogInfo struct {
	File      string        `json:"file"`
	Line      int           `json:"line"`
	Level     string        `json:"level"`
	Log       []interface{} `json:"log"`
	Timestamp int64         `json:"timestamp"`
}

var Console bool
var Rabbitmq bool
var LogMode bool
var channel *amqp.Channel
var queue amqp.Queue

func init() {
	Console = true
	Rabbitmq = false
	channel = nil
	LogMode = false
}

// gin框架的Logger中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if raw != "" {
			path = path + "?" + raw
		}

		Println(path)

		c.Next()
	}
}

// 控制日志是否输出到控制台
func Log_IsConsole(is bool) {
	Console = is
}

func Log_IsRabbitmq(is bool) {
	Rabbitmq = is
}

func Log_LogMode(is bool) {
	LogMode = is
}

// 如果要将日志输出到rabbitmq中，需要调用这个方法进行初始化
//
// @param 	host	主机地址
// @param 	port	主机端口
// @param 	username 用户名
// @param 	passwd	用户密码
// @param 	name	对列名
func Log_init_rabbitmq(host, port, username, passwd, name string) {
	user := fmt.Sprintf("%s:%s", username, passwd)
	url := fmt.Sprintf("amqp://%s@%s:%s/", user, host, port)
	conn, err := amqp.Dial(url)
	if nil != err {
		log.Fatalf("%s: %s", "Open rabbitmq err", err)
	}

	ch, err := conn.Channel()
	if nil != err {
		log.Fatalf("%s: %s", "Channel error", err)
	}

	q, err := ch.QueueDeclare(name, false,
		false, false, false, nil)
	if nil != err {
		log.Fatalf("%s: %s", "QueueDeclare error ", err)
	}

	channel = ch
	queue = q

	return
}

func rabbitmq(msg interface{}) {
	b, err := json.Marshal(msg)
	if nil != err {
		return
	}

	err = channel.Publish("", queue.Name, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})

	if nil != err {
		log.Println("Publish error ", err)

		return
	}

	return
}

func Println(msg ...interface{}) {
	_, file, line, _ := runtime.Caller(1)

	info := LogInfo{
		File:      file,
		Line:      line,
		Level:     "Info",
		Timestamp: time.Now().Unix(),
	}

	info.Log = make([]interface{}, 0)
	for _, v := range msg {
		info.Log = append(info.Log, v)
	}

	var logs interface{}
	if !LogMode {
		logs = info.Log
	} else {
		logs = info
	}

	if Console {
		b, err := json.Marshal(logs)
		if nil != err {
			log.Println(logs, err)
		} else {
			log.Println(string(b))
		}
	}

	if nil != channel && Rabbitmq {
		rabbitmq(logs)
	}

	return
}

func Error(msg ...interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if !LogMode {
		file = ""
		line = 0
	}

	info := LogInfo{
		File:      file,
		Line:      line,
		Level:     "Error",
		Timestamp: time.Now().Unix(),
	}

	info.Log = make([]interface{}, 0)
	for _, v := range msg {
		info.Log = append(info.Log, v)
	}

	var logs interface{}
	if !LogMode {
		logs = info.Log
	} else {
		logs = info
	}

	if Console {
		b, err := json.Marshal(logs)
		if nil != err {
			log.Error(logs, err)
		} else {
			log.Error(string(b))
		}
	}

	if nil != channel && Rabbitmq {
		rabbitmq(logs)
	}

	return
}

func Warning(msg ...interface{}) {
	_, file, line, _ := runtime.Caller(1)

	if !LogMode {
		file = ""
		line = 0
	}

	info := LogInfo{
		File:      file,
		Line:      line,
		Level:     "Warning",
		Timestamp: time.Now().Unix(),
	}

	info.Log = make([]interface{}, 0)
	for _, v := range msg {
		info.Log = append(info.Log, v)
	}

	var logs interface{}
	if !LogMode {
		logs = info.Log
	} else {
		logs = info
	}
	if Console {
		b, err := json.Marshal(logs)
		if nil != err {
			log.Warning(logs, err)
		} else {
			log.Warning(string(b))
		}
	}

	if nil != channel && Rabbitmq {
		rabbitmq(logs)
	}

	return
}
