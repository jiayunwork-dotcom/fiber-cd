package model

const (
	MetersPerUm = 1e-6
	MetersPerNm = 1e-9
	NmPerUm     = 1000.0
)

func NmToM(nm float64) float64 { return nm * MetersPerNm }

func UmToM(um float64) float64 { return um * MetersPerUm }

func MToNm(m float64) float64 { return m / MetersPerNm }

func MToUm(m float64) float64 { return m / MetersPerUm }

func UmToNm(um float64) float64 { return um * NmPerUm }

func NmToUm(nm float64) float64 { return nm / NmPerUm }

func PercentToFraction(percent float64) float64 {
	return percent / 100
}
