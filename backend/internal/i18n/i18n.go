package i18n

var Supported = []string{"en", "de"}

const Default = "en"

func IsSupported(code string) bool {
	for _, c := range Supported {
		if c == code {
			return true
		}
	}
	return false
}
