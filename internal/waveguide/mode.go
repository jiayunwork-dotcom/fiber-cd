package waveguide

type ModeStatus int

const (
	SingleMode ModeStatus = iota
	MultiMode
)

func Classify(v float64) ModeStatus {
	if v < CutoffV {
		return SingleMode
	}
	return MultiMode
}

func (m ModeStatus) String() string {
	switch m {
	case SingleMode:
		return "single-mode"
	default:
		return "multi-mode"
	}
}
