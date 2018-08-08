// +build ocr

package robotgo

import (
	"github.com/otiai10/gosseract"
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
