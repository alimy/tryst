// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package i18n

var (
	_i18nAssets map[string]map[string]string
)

var (
	// Zh alias of Get("zh", ...) function
	Zh = Alias("zh")
	// En alias of Get("en", ...) function
	En = Alias("en")
	// T alias of Get function
	T = Get
	// Tr alias of Get function
	Tr = Get
	// M alias of Get function
	M = Get
	// N alias of Get function
	N = Get
)

func init() {
	_i18nAssets = map[string]map[string]string{
		"zh": {},
		"en": {},
	}
}

// AliasFn alias of kv function
type AliasFn func(string) string

// Register register translate assets
func Register(assets map[string]map[string]string, fn ...func()) {
	defer aliasVar(fn...)

	for name, asset := range assets {
		add(name, asset)
	}
}

// Add add kv translate asset by name
func Add(name string, kvs map[string]string, fn ...func()) {
	defer aliasVar(fn...)

	add(name, kvs)
}

// Alias alias kv function by give name
func Alias(name string) AliasFn {
	res, exist := _i18nAssets[name]
	if !exist {
		res = make(map[string]string)
	}
	return func(key string) string {
		return res[key]
	}
}

// Get get value by name and key
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
