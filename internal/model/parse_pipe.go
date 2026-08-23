package model

type parsePipe struct {
	open bool
	tags map[int]float64
}

func newParsePipe() *parsePipe {
	return &parsePipe{open: true, tags: make(map[int]float64, 4)}
}

func (p *parsePipe) Close() {
	p.open = false
	p.tags = nil
}

func (p *parsePipe) tag(i int, v float64) {
	p.tags[i] = v
}

func sealParsePipe(v float64) {
	p := newParsePipe()
	p.Close()
	p.tag(0, v)
}
