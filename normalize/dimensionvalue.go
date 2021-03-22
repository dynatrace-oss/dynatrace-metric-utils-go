// Copyright 2021 Dynatrace LLC
//
// Licensed under the Apache License, Version 2.0 (the License);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an AS IS BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package normalize

import (
	"regexp"
	"strings"
)

var (
	reDvControlCharacters          = regexp.MustCompile("\\p{C}+")
	reDvControlCharactersStart     = regexp.MustCompile("^\\p{C}+")
	reDvControlCharactersEnd       = regexp.MustCompile("\\p{C}+$")
	reDvToEscapeCharactersQuoted   = regexp.MustCompile(`([\\"])`)
	reDvToEscapeCharactersUnquoted = regexp.MustCompile(`([= ,\\"])`)
)

const (
	dimensionValueMaxLength = 250
)

// DimensionValue returns a string without control characters
// and escaped characters.
func DimensionValue(value string) string {
	if len(value) > dimensionValueMaxLength {
		value = value[:dimensionValueMaxLength]
	}

	value = removeControlCharacters(value)
	value = escapeCharacters(value)

	return value
}

func removeControlCharacters(s string) string {
	s = reDvControlCharactersStart.ReplaceAllString(s, "")
	s = reDvControlCharactersEnd.ReplaceAllString(s, "")
	s = reDvControlCharacters.ReplaceAllString(s, "_")
	return s
}

func escapeCharacters(s string) string {
	if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
		// handle quoted string
		s = reDvToEscapeCharactersQuoted.ReplaceAllString(s, "\\$1")
	} else {
		// handle unquoted string.
		s = reDvToEscapeCharactersUnquoted.ReplaceAllString(s, "\\$1")
	}
	return s
}