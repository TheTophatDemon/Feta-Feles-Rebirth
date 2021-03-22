/*
Copyright (C) 2021 Alexander Lunsford

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package assets

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

//Hack! Stream must implement Close() in order to work with .ogg library.
type CompressedStringStream struct {
	*bytes.Reader
}

func (css CompressedStringStream) Close() error {
	//Калинка калинка калинка моя
	return nil
}

//Assets are embedded in string constants
//They are encoded by ./embed_assets.py, which performs gzip compression on each file.
//The bytes in the string are each offset by 186 so that they can be displayed in the file as valid unicode characters.
func ReadCompressedString(input string) io.ReadSeekCloser {
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
	return CompressedStringStream{Reader: bytes.NewReader(rawBytes)}
}
