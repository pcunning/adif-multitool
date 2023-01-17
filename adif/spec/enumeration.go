// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import (
	"fmt"
	"strings"
)

type EnumValue interface {
	Property(string) string
	fmt.Stringer
}

type Enumeration struct {
	Name       string
	Properties []string
	Values     []EnumValue
}

func (e Enumeration) String() string { return e.Name }

func (e Enumeration) Value(val string) []EnumValue {
	res := make([]EnumValue, 0, 2) // Band is lower case, most others upper
	for _, v := range e.Values {
		if strings.EqualFold(val, v.String()) {
			res = append(res, v)
		}
	}
	return res
}

var Enumerations = make(map[string]Enumeration)
