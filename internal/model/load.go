package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// LoadFile 从文件系统读取一份 JSON 配置。文件不存在、不可读或 JSON
// 非法时返回带路径的错误，CLI 会把它原样打到 stderr。
func LoadFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read config %q: %w", path, err)
	}
	cfg, err := Parse(data)
	if err != nil {
		return Config{}, fmt.Errorf("config %q: %w", path, err)
	}
	return cfg, nil
}

// Parse 把一段 JSON 字节解码为 Config。缺省字段保持零值，由调用方
// 决定何时走 Validate 或 SweepRange 的默认值逻辑。
func Parse(data []byte) (Config, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return cfg, nil
}

// Decode 从任意 io.Reader 解码配置，供测试与管道输入复用。
func Decode(r io.Reader) (Config, error) {
	dec := json.NewDecoder(r)
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("invalid JSON: %w", err)
	}
	return cfg, nil
}
