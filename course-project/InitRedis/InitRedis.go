package InitRedis

import (
	"github.com/go-redis/redis/v7"
)

var client *redis.Client = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "12345", // It's ok if password is "".
	DB:       0,       // use default DB
})

// 开启redis连接池
func InitRedisConnection() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "12345", // It's ok if password is "".
		DB:       0,       // use default DB
	})

	if _, err := FlushAll(); err != nil {
		println("Error when flushAll. " + err.Error())
	}
}

// 用于测试
func FlushAll() (string, error) {
	return client.FlushAll().Result()
}

func SetForever(key string, value interface{}) (string, error) {
	val, err := client.Set(key, value, 0).Result() // expiration表示无过期时间
	return val, err
}

func GetMap(key string, fields ...string) ([]interface{}, error) {
	return client.HMGet(key, fields...).Result()
}

// 确保redis加载lua脚本，若未加载则加载
func PrepareScript(script string) string {
	// sha := sha1.Sum([]byte(script))
	scriptsExists, err := client.ScriptExists(script).Result()
	if err != nil {
		panic("Failed to check if script exists: " + err.Error())
	}
	if !scriptsExists[0] {
		scriptSHA, err := client.ScriptLoad(script).Result()
		if err != nil {
			panic("Failed to load script " + script + " err: " + err.Error())
		}
		return scriptSHA
	}
	print("Script Exists.")
	return ""
}

// 执行lua脚本
func EvalSHA(sha string, args []string) (interface{}, error) {
	val, err := client.EvalSha(sha, args).Result()
	if err != nil {
		print("Error executing evalSHA... " + err.Error())
		return nil, err
	}
	return val, nil
}
