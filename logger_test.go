package youyoulibs

import (
	"testing"
)

/**

sudo docker run -it -d --rm  \
 -p 15672:15672 \
 -p 5672:5672 \
 -e RABBITMQ_DEFAULT_USER=qinxue \
 -e RABBITMQ_DEFAULT_PASS=qinxue \
 --name rabbitmq \
 rabbitmq:3-management
*/

func TestLog_init_rabbitmq(t *testing.T) {
	Log_init_rabbitmq("127.0.0.1", "5672", "qinxue", "qinxue", "qinxue")
}

func TestLog_I(t *testing.T) {
	var lg LogInfo
	Println("TestLog_W", lg)

	return
}

func TestLog_E(t *testing.T) {
	Error("TestLog_E")

	return
}

func TestLog_W(t *testing.T) {
	Log_IsRabbitmq(true)
	Log_LogMode(true)

	Warning("TestLog_W")

	return
}
