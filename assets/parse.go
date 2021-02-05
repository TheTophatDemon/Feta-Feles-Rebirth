package assets

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

//Assets are embedded in string constants
//They are encoded by ./embed_assets.py, which performs gzip compression on each file.
//The bytes in the string are each offset by 186 so that they can be displayed in the file as valid unicode characters.
func ReadCompressedString(input string) io.ReadSeeker {
	zipBytes := make([]byte, len(input))
	i := 0
	for _, r := range input {
		v := int(r) - 186
		zipBytes[i] = byte(v)
		i++
	}
	//Error checking is skipped for now. Python's gzip library doesn't seem to write valid headers.
	zipReader, _ := gzip.NewReader(bytes.NewReader(zipBytes))
	defer zipReader.Close()
	rawBytes, _ := ioutil.ReadAll(zipReader)
	return bytes.NewReader(rawBytes)
}
