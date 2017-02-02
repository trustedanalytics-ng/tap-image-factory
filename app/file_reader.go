/**
 * Copyright (c) 2016 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"archive/tar"
	"compress/gzip"
	"io"
)

type FileReader interface {
	NewGzipReader(reader io.Reader) (*gzip.Reader, error)
	NewTarReader(reader io.Reader) *tar.Reader
}

type ArchiveReader struct{}

func (c *ArchiveReader) NewGzipReader(reader io.Reader) (*gzip.Reader, error) {
	return gzip.NewReader(reader)
}

func (c *ArchiveReader) NewTarReader(reader io.Reader) *tar.Reader {
	return tar.NewReader(reader)
}
