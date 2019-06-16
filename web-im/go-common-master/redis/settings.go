package redis

import (
	"fmt"
	"strconv"

	myLog "../log"

	"github.com/gomodule/redigo/redis"
)

var (
	ErrorNoFoundValue      = fmt.Errorf("no found value")
	ErrorValueNoArrayByte  = fmt.Errorf("value no []byte")
	ErrorValueNoArrayUint8 = fmt.Errorf("value no []uint8")

	maxIdle      = 3
	maxActive    = 8
	network      = "tcp"
	ipv4         = "localhost"
	dbPort       = "6379"
	dialPassword = ""
)

func GetPool() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:   maxIdle,
		MaxActive: maxActive,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(network, ipv4+":"+dbPort, redis.DialPassword(dialPassword))
			if err != nil {
				return nil, err
			}
			return conn, nil
		}}
	return pool
}

func GetConn() redis.Conn {
	return GetPool().Get()
}

//获取string值
func GetString(key string) (string, error) {
	result, err := GetResult(key)
	if err != nil {
		return "", err
	}
	arrayByteResult, ok := result.([]byte)
	if !ok {
		return "", ErrorValueNoArrayByte
	}
	resultString := string(arrayByteResult)
	return resultString, nil
}

//获取int值
func GetInt(key string) (int, error) {
	result, err := GetResult(key)
	if err != nil {
		return 0, err
	}
	arrayUint8Result, ok := result.([]uint8)
	if !ok {
		return 0, ErrorValueNoArrayUint8
	}
	resultInt, _ := strconv.Atoi(string(arrayUint8Result))
	return resultInt, nil
}

//获取结果
func GetResult(key string) (interface{}, error) {
	conn := GetConn()
	defer conn.Close()
	result, err := conn.Do("GET", key)
	if result == nil {
		return "", ErrorNoFoundValue
	}
	return result, err
}

//没有效果
// func ConnDo(commandName string, args ...interface{}) (interface{}, error) {
// 	conn := GetConn()
// 	defer conn.Close()
// 	reply, err := conn.Do(commandName, args)
// 	return reply, err
// }

//并集两个集合，并简历新集合
func ZunionStoreTwo(destination string, key1 string, key2 string) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("ZUNIONSTORE", destination, 2,
		key1, key2,
		"AGGREGATE", "SUM")
	return err
}

//按分值，从高至低返回集合数据
func ZrevRangeKey(key string) (interface{}, error) {
	conn := GetConn()
	defer conn.Close()
	reply, err := conn.Do("ZREVRANGE", key, 0, -1)
	return reply, err
}

//延长key时间
func ExpireKey(key string, seconds int) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", key, seconds)
	return err
}

//按正则公式，返回满足条件的key
func Keys(pattern string) (interface{}, error) {
	conn := GetConn()
	defer conn.Close()
	reply, err := conn.Do("KEYS", pattern)
	return reply, err
}

//删除key
func DelKey(key string) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("DEL", key)
	return err
}

//延长key到某个时间点
func ExpireAtKey(key string, timestamp int64) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("EXPIREAT", key, timestamp)
	return err
}

//设置key的值和过期时间
func SetExpireKey(key string, value interface{}, seconds int) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("SET", key, value, "EX", seconds)
	return err
}

//若key不存在，设置key
func SetNXExpireKey(key string, value interface{}, seconds int) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("SET", key, value, "EX", seconds, "NX")
	return err
}

//列表插入头部数据
func LpushKey(key string, value interface{}) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("LPUSH", key, value)
	return err
}

//列表删除尾部数据
func RpopKey(key string) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("RPOP", key)
	return err
}

//获取列表头部数据
func LheadKey(key string) (string, error) {
	conn := GetConn()
	defer conn.Close()
	reply, err := conn.Do("LRANGE", key, 0, 1)
	if err != nil {
		return "", nil
	}
	for _, dbMessage := range reply.([]interface{}) {
		messageString := string(dbMessage.([]uint8))
		return messageString, nil
	}
	return "", nil
}

//获取列表大小
func LlenKey(key string) (int, error) {
	conn := GetConn()
	defer conn.Close()
	reply, err := conn.Do("LLEN", key)
	return int(reply.(int64)), err
}

//返回列表所有数据
func LrangeKey(key string) (interface{}, error) {
	conn := GetConn()
	defer conn.Close()
	reply, err := conn.Do("LRANGE", key, 0, -1)
	return reply, err
}

//有序集合新增数据
func ZaddKey(key string, score int, member interface{}) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("ZADD", key, score, member)
	return err
}

//值++
func INCRKey(key string) error {
	conn := GetConn()
	defer conn.Close()
	_, err := conn.Do("INCR", key)
	return err
}

func PrintErrorMsg(errString string) error {
	return PrintError(fmt.Errorf(errString), errString)
}

func PrintError(err error, errString string) error {
	fmt.Println(err)
	myLog.LogErrorInfo("redis", "err:"+err.Error(), "errString:"+errString)
	return fmt.Errorf(errString)
}
