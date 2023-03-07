// Copyright (c) 2017 Yandex LLC. All rights reserved.
// Use of this source code is governed by a MPL 2.0
// license that can be found in the LICENSE file.
// Author: Vladimir Skipor <skipor@yandex-team.ru>

package simple

import (
	"net/http"

	phttp "github.com/yandex/pandora/components/guns/http"
	"github.com/yandex/pandora/core/aggregator/netsample"
)

type Ammo struct {
	// OPTIMIZE(skipor): reuse *http.Request.
	// Need to research is it possible. http.Transport can hold reference to http.Request.
	req       *http.Request
	tag       string
	id        uint64
	isInvalid bool
}

func (a *Ammo) Request() (*http.Request, *netsample.Sample) {
	sample := netsample.Acquire(a.tag)
	sample.SetID(a.id)
	return a.req, sample
}

func (a *Ammo) Reset(req *http.Request, tag string) {
	*a = Ammo{req, tag, 0, false}
}

func (a *Ammo) SetID(id uint64) {
	a.id = id
}

func (a *Ammo) ID() uint64 {
	return a.id
}

func (a *Ammo) Invalidate() {
	a.isInvalid = true
}

func (a *Ammo) IsInvalid() bool {
	return a.isInvalid
}

func (a *Ammo) IsValid() bool {
	return !a.isInvalid
}

func (a *Ammo) Tag() string {
	return a.tag
}

var _ phttp.Ammo = (*Ammo)(nil)
