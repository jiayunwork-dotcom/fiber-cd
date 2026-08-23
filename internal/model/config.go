// Package model 定义阶跃光纤的输入配置、JSON 装载与非法输入校验。
// 所有 fiber-cd 子命令共享这份 schema：纤芯/包层折射率、芯径
// （直径，不是半径）、工作波长，以及可选的波长扫描区间。
package model

import (
	"encoding/json"
	"fmt"
)

// Config 描述一根圆阶跃光纤。单位约定：
//
//   - n1、n2：纤芯、包层折射率（无量纲），必须满足 n2 < n1；
//   - CoreDiameterUm：芯径，单位 μm，是直径；内部统一折算成半径使用；
//   - WavelengthNm：工作波长，单位 nm；
//   - SweepStartNm / SweepStopNm / SweepSteps：sweep 子命令的扫描区间
//     与步数，可缺省，缺省时用包级默认值。
type Config struct {
	Name string `json:"name"`

	N1 float64 `json:"n1"`
	N2 float64 `json:"n2"`

	CoreDiameterUm float64 `json:"core_diameter_um"`
	WavelengthNm   float64 `json:"wavelength_nm"`

	SweepStartNm float64 `json:"sweep_start_nm,omitempty"`
	SweepStopNm  float64 `json:"sweep_stop_nm,omitempty"`
	SweepSteps   int     `json:"sweep_steps,omitempty"`
}

// Clone 返回一份与接收者等价的独立拷贝（Config 全为值字段）。
func (c Config) Clone() Config {
	return c
}

// Description 返回一行可用于报告标题的描述文字。没有设置 name 时
// 退回用折射率差与芯径拼出的中性描述。
func (c Config) Description() string {
	if c.Name != "" {
		return c.Name
	}
	return "step-index fiber"
}

// String 把配置的主要参数压成一行便于调试打印。
func (c Config) String() string {
	return fmt.Sprintf("name=%q n1=%.6f n2=%.6f core_diameter_um=%.4f wavelength_nm=%.1f",
		c.Name, c.N1, c.N2, c.CoreDiameterUm, c.WavelengthNm)
}

// MarshalIndent 把配置序列化为缩进 JSON，供 dump 子命令回显，也方便
// 用户以当前配置为模板改字段。
func (c Config) MarshalIndent() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}
