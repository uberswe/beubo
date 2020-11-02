package component

func returnTIfNotEmpty(t string, d string) string {
	if t != "" {
		return t
	}
	return d
}
