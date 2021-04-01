// Copyright 2021 Dynatrace LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package oneagentenrichment

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dynatrace-oss/dynatrace-metric-utils-go/metric/dimensions"
)

const indirectionFilename = "dt_metadata_e617c525669e072eebe3d0f08212e8f2.properties"

// readIndirectionFile reads the first line from the Reader and returns it
// or an error if there was a problem while reading.
func readIndirectionFile(reader io.Reader) (string, error) {
	if reader == nil {
		return "", errors.New("reader cannot be nil")
	}

	scanner := bufio.NewScanner(reader)

	scanner.Scan()
	indirectionFilename := strings.TrimSpace(scanner.Text())

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return indirectionFilename, nil
}

// readMetadataFile reads the actual OneAgent metadata file from the passed Reader, discarding empty lines.
func readMetadataFile(reader io.Reader) ([]string, error) {
	if reader == nil {
		return nil, errors.New("reader cannot be nil")
	}

	scanner := bufio.NewScanner(reader)
	lines := []string{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// readOneAgentMetadata takes the name of the properties file. It then reads the indirection file
// to get the name of the actual metadata file. That file is then read and parsed into an array of strings,
// which represent the lines of that file. Errors from function calls inside this function are passed on to the caller.
func readOneAgentMetadata(indirectionFileName string) ([]string, error) {
	// Currently, this only works on Windows hosts, since the indirection on Linux
	// is based on libc. As Go does not use libc to open files, this doesnt currently
	// work on Linux hosts.
	indirection, err := os.Open(indirectionFileName)
	if err != nil {
		// an error occurred during opening of the file
		return nil, err
	}
	defer indirection.Close()

	filename, err := readIndirectionFile(indirection)
	if err != nil {
		// an error occurred during reading of the file
		return nil, err
	}

	if filename == "" {
		return nil, errors.New("metadata file name is empty")
	}

	metadataFile, err := os.Open(filename)
	if err != nil {
		// an error occurred during opening of the file
		return nil, err
	}

	content, err := readMetadataFile(metadataFile)

	if err != nil {
		// an error occurred during reading of the file
		return nil, err
	}

	return content, nil
}

// parseOneAgentMetadata transforms lines into key-value pairs and discards
// pairs that do not conform to the 'key=value' (trailing additional equal signs are added to
// the value)
func parseOneAgentMetadata(lines []string) []dimensions.Dimension {
	result := []dimensions.Dimension{}
	for _, line := range lines {
		split := strings.SplitN(line, "=", 2)
		if len(split) != 2 {
			log.Println(fmt.Sprintf("Could not parse OneAgent metadata line '%s'", line))
			continue
		}
		key, value := split[0], split[1]

		if key == "" || value == "" {
			log.Println(fmt.Sprintf("Could not parse OneAgent metadata line '%s'", line))
			continue
		}

		result = append(result, dimensions.NewDimension(key, value))
	}
	return result
}

func asNormalizedDimensionList(lines []string) dimensions.NormalizedDimensionList {
	if len(lines) == 0 {
		return dimensions.NewNormalizedDimensionList()
	}

	dims := []dimensions.Dimension{}
	for _, dim := range parseOneAgentMetadata(lines) {
		dims = append(dims, dim)
	}
	return dimensions.NewNormalizedDimensionList(dims...)
}

// GetOneAgentMetadata reads the metadata and returns them as NormalizedDimensionList
func GetOneAgentMetadata() dimensions.NormalizedDimensionList {
	lines, err := readOneAgentMetadata(indirectionFilename)
	if err != nil {
		log.Println("Could not read OneAgent metadata. This is normal if no OneAgent is installed, or if you are running this on Linux.")
		return dimensions.NewNormalizedDimensionList()
	}

	return asNormalizedDimensionList(lines)
}
