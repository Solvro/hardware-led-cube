package main

import (
	"bytes"
	"encoding/json"
	"gotest.tools/v3/assert"
	"os"
	"testing"
)

func TestJSONFileDecoding(t *testing.T) {
	file, err := os.Open("parsing_test_files/json_file_input.json")
	assert.NilError(t, err)
	fs, ec := NewJSONFileAnimation(file)

	checkParsingError := errChanCheckerRecover(t, fs, ec)
	checkParsingError()
	for range fs.frameChan {
		checkParsingError()
	}
	// all this does is check if decoding into frames and encoding back works
	jsonBytes, err := json.Marshal(fs.frames)
	assert.NilError(t, err)
	tmp := bytes.NewBuffer(nil)
	err = json.Compact(tmp, jsonBytes)
	assert.NilError(t, err)

	out := tmp.String()

	jsonBytes, err = os.ReadFile("parsing_test_files/json_file_input.json")
	assert.NilError(t, err)
	tmp.Reset()
	err = json.Compact(tmp, jsonBytes)
	assert.NilError(t, err)

	in := tmp.String()

	assert.DeepEqual(t, out, in)
}

// TODO: test the normalization function
// TODO: more and better test files

// recovers the panic on errors from ec and fails through t instead
func errChanCheckerRecover(t *testing.T, fs FrameSource, ec <-chan error) func() {
	checkError := errChanChecker(fs, ec)
	return func() {
		checkError()
		defer func() {
			if r := recover(); r != nil {
				t.Fatal(recover())
			}
		}()
	}
}
