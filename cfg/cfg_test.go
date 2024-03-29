// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package cfg

import (
	"testing"
)

func TestCfg(t *testing.T) {
	suites := map[string][]string{
		"default": {"Sms", "Alipay", "Zinc", "MySQL", "Redis", "AliOSS", "LogZinc"},
		"develop": {"Zinc", "MySQL", "AliOSS", "LogFile"},
		"slim":    {"Zinc", "MySQL", "Redis", "AliOSS", "LogFile"},
	}
	kv := map[string]string{
		"sms": "SmsJuhe",
	}

	Initial(suites, kv)
	UseDefault()

	if !If("Sms") {
		t.Error(`want If("Sms") == true but not`)
	}

	if All("Sms", "Alipay", "Zinc", "LogFile") {
		t.Error(`want All("Sms", "Alipay", "Zinc", "LogFile") == false but not`)
	}

	if !Any("Sms", "Alipay", "Zinc") {
		t.Error(`want Any("Sms", "Alipay", "Zinc", "LogFile") == true but not`)
	}

	if !Any("SmsNo", "Alipays", "Zinc", "LogFile") {
		t.Error(`want Any("SmsNo", "Alipays", "Zinc", "LogFile") == true but not`)
	}

	if v, exist := Val("Sms"); exist && v != "SmsJuhe" {
		t.Errorf(`want Val("Sms") == "SmsJuhe", true but got: "%s", "%t"`, v, exist)
	}

	As("sms", func(v string) {
		if v != "SmsJuhe" {
			t.Errorf(`want As("Sms") == "SmsJuhe", true but got: "%s"`, v)
		}
	})

	matched := false
	Be("Alipay", func() {
		matched = true
	})
	if !matched {
		t.Error(`want Be("Alipay", ...) matched but not`)
	}

	matched = false
	Not("LogFile", func() {
		matched = true
	})
	if !matched {
		t.Error(`want Not("LogFile", ...) matched but not`)
	}

	var m1, m2, m3, m4 bool
	In(Actions{
		"Sms": func() {
			m1 = true
		},
		"Alipay": func() {
			m2 = true
		},
		"Meili": func() {
			m4 = true
		},
	}, func() {
		m3 = true
	})
	if !m1 || !m2 || m3 || m4 {
		t.Errorf(`In("Sms", "Alipay", "Meili", ...) not correct -> m1: %t m2:%t m3:%t m4:%t`, m1, m2, m3, m4)
	}

	m1 = false
	m2 = false
	m3 = false
	In(Actions{
		"LogFile": func() {
			m1 = true
		},
		"Meili": func() {
			m2 = true
		},
	}, func() {
		m3 = true
	})
	if m1 || m2 || !m3 {
		t.Errorf(`In("Zinc", "MySQL", ...) not correct -> m1: %t m2:%t m3:%t`, m1, m2, m3)
	}

	m1 = false
	m2 = false
	m3 = false
	On(Actions{
		"AliOSS": func() {
			m1 = true
		},
		"Localoss": func() {
			m2 = true
		},
	}, func() {
		m3 = true
	})
	if !m1 || m2 || m3 {
		t.Errorf(`On("AliOSS", "Localoss", ...) not correct -> m1: %t m2:%t m3:%t`, m1, m2, m3)
	}

	m1 = false
	m2 = false
	m3 = false
	On(Actions{
		"COS": func() {
			m1 = true
		},
		"Localoss": func() {
			m2 = true
		},
	}, func() {
		m3 = true
	})
	if m1 || m2 || !m3 {
		t.Errorf(`On("AliOSS", "Localoss", ...) not correct -> m1: %t m2:%t m3:%t`, m1, m2, m3)
	}
}
