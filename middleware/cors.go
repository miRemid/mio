package middleware

import (
	"github.com/miRemid/mio"
)

// Config 跨域设置
type Config struct {
}

// Generate 生成跨域中间件
func (c *Config) Generate() mio.HandlerFunc {
	return func(c *mio.Context) {

	}
}
