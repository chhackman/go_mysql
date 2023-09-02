//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		//本地连接
		DSN: "indigo:indigotest@tcp(10.1.90.122:3306)/go_test",
	},
	Redis: RedisConfig{
		Addr: "10.1.90.235:6379",
	},
}
