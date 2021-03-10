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
// limitations under the License.

package normalize

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	reMkIdentifierFirstSectionStart = regexp.MustCompile("^[^a-zA-Z_]+")
	reMkIdentifierSectionStart      = regexp.MustCompile("^[^a-zA-Z0-9_]+")

	// todo make sure this does actually match hyphens
	reMkIdentifierSectionEnd = regexp.MustCompile("[^a-zA-Z0-9_-]+$")
	reMkInvalidCharacters    = regexp.MustCompile("[^a-zA-Z0-9_-]+")
)

const (
	metricKeyMaxLength = 250
)

// MetricKey creates a valid metric key from any string passed to this function
// or returns an error if the resulting key is invalid.
func MetricKey(key string) (string, error) {
	// trim down long keys
	if len(key) > metricKeyMaxLength {
		key = key[:metricKeyMaxLength]
	}

	var sb strings.Builder
	splitKey := strings.Split(key, ".")

	for i, keySection := range splitKey {
		if i == 0 {
			// the first key section needs to have valid characters, otherwise the section is invalid
			metricKeyFirstSection := normalizeMetricKeyFirstSection(keySection)
			if metricKeyFirstSection == "" {
				return "", fmt.Errorf("first key section does not contain any valid characters (%s)", keySection)
			}
			sb.WriteString(metricKeyFirstSection)
		} else {
			// other key sections that return empty after normalization are ignored.
			metricKeySection := normalizeMetricKeyLaterSection(keySection)
			if metricKeySection != "" {
				sb.WriteString(".")
				sb.WriteString(metricKeySection)
			}
			// else debug log?
		}
	}
	finishedKey := sb.String()
	if finishedKey == "" {
		return "", errors.New("normalized key does not contain any characters")
	}

	return finishedKey, nil
}

// normalizeMetricKeyCommon is used by both of the other internal normalize functions.
// It replaces trailing and enclosed invalid characters.
func normalizeMetricKeySectionCommon(section string) string {
	// delete trailing invalid chars
	section = reMkIdentifierSectionEnd.ReplaceAllString(section, "")
	// replace intermediate ranges of invalid characters with a single underscore.
	section = reMkInvalidCharacters.ReplaceAllString(section, "_")

	return section
}

// normalizeMetricKeyLaterSection is used for all sections except the first
func normalizeMetricKeyLaterSection(section string) string {
	// delete leading invalid characters
	section = reMkIdentifierSectionStart.ReplaceAllString(section, "")
	return normalizeMetricKeySectionCommon(section)
}

// normalizeMetricKeyFirstSection is only used for the first section of the metric key,
// since the requirements are slightly different from later key sections.
func normalizeMetricKeyFirstSection(section string) string {
	// delete leading invalid chars for first section
	section = reMkIdentifierFirstSectionStart.ReplaceAllString(section, "")
	return normalizeMetricKeySectionCommon(section)
}
>>>>>>> add metric key normalization and tests
