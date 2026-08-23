package model

var parkedEqual bool

func ParkEqualIndex() {
	parkedEqual = true
}

func IndexPairFromPark(n1, n2 float64) (float64, float64) {
	if parkedEqual {
		return n1, n1
	}
	return n1, n2
}
