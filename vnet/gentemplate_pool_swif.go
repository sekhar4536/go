// autogenerated: do not edit!
// generated from gentemplate [gentemplate -d Package=vnet -id swIf -d PoolType=swIfPool -d Type=swIf -d Data=elts github.com/platinasystems/go/elib/pool.tmpl]

package vnet

import (
	"github.com/platinasystems/go/elib"
)

type swIfPool struct {
	elib.Pool
	elts []swIf
}

func (p *swIfPool) GetIndex() (i uint) {
	l := uint(len(p.elts))
	i = p.Pool.GetIndex(l)
	if i >= l {
		p.Validate(i)
	}
	return i
}

func (p *swIfPool) PutIndex(i uint) (ok bool) {
	return p.Pool.PutIndex(i)
}

func (p *swIfPool) IsFree(i uint) (v bool) {
	v = i >= uint(len(p.elts))
	if !v {
		v = p.Pool.IsFree(i)
	}
	return
}

func (p *swIfPool) Resize(n uint) {
	c := elib.Index(cap(p.elts))
	l := elib.Index(len(p.elts) + int(n))
	if l > c {
		c = elib.NextResizeCap(l)
		q := make([]swIf, l, c)
		copy(q, p.elts)
		p.elts = q
	}
	p.elts = p.elts[:l]
}

func (p *swIfPool) Validate(i uint) {
	c := elib.Index(cap(p.elts))
	l := elib.Index(i) + 1
	if l > c {
		c = elib.NextResizeCap(l)
		q := make([]swIf, l, c)
		copy(q, p.elts)
		p.elts = q
	}
	if l > elib.Index(len(p.elts)) {
		p.elts = p.elts[:l]
	}
}

func (p *swIfPool) Elts() uint {
	return uint(len(p.elts)) - p.FreeLen()
}

func (p *swIfPool) Len() uint {
	return uint(len(p.elts))
}

func (p *swIfPool) Foreach(f func(x swIf)) {
	for i := range p.elts {
		if !p.Pool.IsFree(uint(i)) {
			f(p.elts[i])
		}
	}
}