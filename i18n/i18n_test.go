// Copyright 2024 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package i18n_test

import (
	"testing"

	"github.com/alimy/tryst/i18n"
)

var (
	cn   i18n.AliasFn
	zhTw = i18n.Alias("zh_tw")

	assets = map[string]map[string]string{
		"zh": {
			"hello":  "你好",
			"golang": "go语言",
		},
		"en": {
			"hello":  "hello",
			"golang": "golang",
		},
	}

	zhTwKvs = map[string]string{
		"hello":  "你好",
		"golang": "go语言",
	}
)

func TestRegister(t *testing.T) {
	i18n.Register(assets, func() {
		cn = i18n.Alias("cn")
	})

	if value := i18n.Get("zh", "hello"); value != "你好" {
		t.Errorf(`call i18n.Get("zh", "hello") want "你好" but got %s`, value)
	}

	if value := i18n.Get("zh", "nothing"); value != "" {
		t.Errorf(`call i18n.Get("zh", "nothing") want "" but got %s`, value)
	}

	if value := i18n.Get("en", "hello"); value != "hello" {
		t.Errorf(`call i18n.Get("en", "hello") want "hello" but got %s`, value)
	}

	if value := i18n.Get("en", "nothing"); value != "" {
		t.Errorf(`call i18n.Get("en", "nothing") want "" but got %s`, value)
	}

	if value := i18n.Get("zh", "hello"); value != "你好" {
		t.Errorf(`call i18n.Get("zh", "hello") want "你好" but got %s`, value)
	}

	if value := i18n.Get("zh", "nothing"); value != "" {
		t.Errorf(`call i18n.Get("zh", "nothing") want "" but got %s`, value)
	}

	if value := i18n.Zh("hello"); value != "你好" {
		t.Errorf(`call i18n.Zh("hello") want "你好" but got %s`, value)
	}

	if value := i18n.Zh("nothing"); value != "" {
		t.Errorf(`call i18n.Get("nothing") want "" but got %s`, value)
	}

	if value := i18n.En("hello"); value != "hello" {
		t.Errorf(`call i18n.En("hello") want "hello" but got %s`, value)
	}

	if value := i18n.En("nothing"); value != "" {
		t.Errorf(`call i18.En("nothing") want "" but got %s`, value)
	}

	if value := cn("hello"); value != "" {
		t.Errorf(`call i18n.En("hello") want "" but got %s`, value)
	}

	if value := i18n.Get("cn", "nothing"); value != "" {
		t.Errorf(`call i18.Get("cn", "nothing") want "" but got %s`, value)
	}

	if value := i18n.Get("cn", "nothing", "abc"); value != "abc" {
		t.Errorf(`call i18.Get("cn", "nothing", "abc") want "abc" but got %s`, value)
	}
}

func TestAdd(t *testing.T) {
	i18n.Add("zh_tw", zhTwKvs, func() {
		zhTw = i18n.Alias("zh_tw")
	})

	if value := i18n.Get("zh_tw", "hello"); value != "你好" {
		t.Errorf(`call i18n.Get("zh_tw", "hello") want "你好" but got %s`, value)
	}

	if value := i18n.Get("zh_tw", "nothing"); value != "" {
		t.Errorf(`call i18n.Get("zh_tw", "nothing") want "" but got %s`, value)
	}

	if value := zhTw("hello"); value != "你好" {
		t.Errorf(`call zhTw("hello") want "你好" but got %s`, value)
	}

	if value := zhTw("nothing"); value != "" {
		t.Errorf(`call zhTw("nothing") want "" but got %s`, value)
	}
}
