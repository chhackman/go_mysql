//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		//k8s标签
		DSN: "root:root@tcp(webook-live-mysql:3308)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-live-redis:6380",
	},
}
