package util

func Exn(e error) {
	if e != nil {
		panic("Exception: " + e.Error())
	}
}
