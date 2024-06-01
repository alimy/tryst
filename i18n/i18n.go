// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package i18n

var (
	_i18nAssets map[string]map[string]string
)

var (
	Zh = Alias("zh")
	En = Alias("en")
)

func init() {
	_i18nAssets = map[string]map[string]string{
		"zh": {},
		"en": {},
	}
}

type AliasFn func(string) string

func Register(assets map[string]map[string]string, fn ...func()) {
	defer aliasVar(fn...)

	for name, asset := range assets {
		add(name, asset)
	}
}

func Add(name string, kvs map[string]string, fn ...func()) {
	defer aliasVar(fn...)

	add(name, kvs)
}

func Alias(name string) AliasFn {
	res, exist := _i18nAssets[name]
	if !exist {
		res = make(map[string]string)
	}
	return func(key string) string {
		return res[key]
	}
}

func Get(name string, key string) (value string) {
	if kvs, exist := _i18nAssets[name]; exist {
		value = kvs[key]
	}
	return
}

func aliasVar(fn ...func()) {
	defer func() {
		if len(fn) > 0 && fn[0] != nil {
			fn[0]()
		}
	}()

	Zh = Alias("zh")
	En = Alias("en")
}

func add(name string, kvs map[string]string) {
	if len(name) == 0 || len(kvs) == 0 {
		return
	}

	if _, exist := _i18nAssets[name]; !exist {
		_i18nAssets[name] = make(map[string]string, 10)
	}
	for k, v := range kvs {
		if len(k) > 0 && len(v) > 0 {
			_i18nAssets[name][k] = v
		}
	}
}
