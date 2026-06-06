// Copyright (c) 2016-2026 AtomAI, All rights reserved.
//
// See the COPYRIGHT file at the top-level directory of this distribution and at
// https://github.com/go-vgo/robotgo/blob/master/LICENSE
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
//
// This file may not be copied, modified, or distributed
// except according to those terms.

//go:build ocr
// +build ocr

package robotgo

import (
	"github.com/otiai10/gosseract/v2"
)

// GetText get the image text by tesseract ocr
func GetText(imgPath string, args ...string) (string, error) {
	var lang = "eng"

	if len(args) > 0 {
		lang = args[0]
		if lang == "zh" {
			lang = "chi_sim"
		}
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(imgPath)
	client.SetLanguage(lang)
	return client.Text()
}
