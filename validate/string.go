package validate

func Capitalize(str string) string {
	vv := []rune(str)
	if len(vv) > 0 && (vv[0] >= 97 && vv[0] <= 122) {
		vv[0] = vv[0] - 32
	}
	return string(vv)
}