package waveguide

// ModeStatus 表示一根光纤在某波长下可承载的模式数量。
type ModeStatus int

const (
	// SingleMode 单模：V < 2.405，只承载基模 LP01。
	SingleMode ModeStatus = iota
	// MultiMode 多模：V >= 2.405，次低模式已截止进入导波。
	MultiMode
)

// Classify 按归一化频率 V 判定模式状态。边界取严格不等号，
// V = 2.405 这一档不判单模。
func Classify(v float64) ModeStatus {
	if v < CutoffV {
		return SingleMode
	}
	return MultiMode
}

// String 输出面向用户的状态文本。
func (m ModeStatus) String() string {
	switch m {
	case SingleMode:
		return "single-mode"
	default:
		return "multi-mode"
	}
}
