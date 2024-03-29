// Copyright 2023 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package cfg

var (
	_features = newEmptyFeatures()

	// Use alias of Features.Use func
	Use = _features.Use

	// UseDeafult alias of Features.UseDefault func
	UseDefault = _features.UseDefault

	// As alias of Features.Cfg func
	Val = _features.Cfg

	// As alias of Features.CfgAs func
	As = _features.CfgAs

	// If alias of Features.CfgIf func
	If = _features.CfgIf

	// If alias of Features.CfgAll func
	All = _features.CfgAll

	// If alias of Features.CfgAny func
	Any = _features.CfgAny

	// In alias of Features.CfgIn func
	In = _features.CfgIn

	// On alias of Features.CfgOn func
	On = _features.CfgOn

	// Be alias of Feaures.CfgBe func
	Be = _features.CfgBe

	// Not alias of Features.CfgNot func
	Not = _features.CfgNot
)

// Initial initialize features in cfg pkg
func Initial(suites map[string][]string, kv map[string]string) {
	_features = NewFeatures(suites, kv)
	{
		// must re-assign variable below
		Use = _features.Use
		UseDefault = _features.UseDefault
		Val = _features.Cfg
		As = _features.CfgAs
		If = _features.CfgIf
		All = _features.CfgAll
		Any = _features.CfgAny
		In = _features.CfgIn
		On = _features.CfgOn
		Be = _features.CfgBe
		Not = _features.CfgNot
	}
}
