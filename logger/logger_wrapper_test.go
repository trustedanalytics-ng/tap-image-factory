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
package logger_wrapper

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInitLogger(t *testing.T) {
	Convey("Test InitLogger should initialize logger with defined module name", t, func() {
		logger := InitLogger("module")
		So(logger.Module, ShouldEqual, "module")
	})
}
