package model

// 扫描与单位相关的包级默认值。sweep 子命令没有在配置里写扫描区间时，
// 统一使用 C 波段附近的 1200–1700 nm 与 51 档密度。
const (
	// DefaultSweepStartNm 波长扫描起点（nm），约 1200 nm。
	DefaultSweepStartNm = 1200.0
	// DefaultSweepStopNm 波长扫描终点（nm），约 1700 nm。
	DefaultSweepStopNm = 1700.0
	// DefaultSweepSteps 扫描档数（含两端点）。
	DefaultSweepSteps = 51
)

// DefaultSweep 返回一组与 DefaultSweep* 常量一致的扫描参数。
func DefaultSweep() (startNm, stopNm float64, steps int) {
	return DefaultSweepStartNm, DefaultSweepStopNm, DefaultSweepSteps
}
