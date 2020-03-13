package middleware

import (
	"mioqq"
)

// Config 跨域设置
type Config struct {
}

// Generate 生成跨域中间件
func (c *Config) Generate() mioqq.HandlerFunc {
	return func(c *mioqq.Context) {
		
	}
}
