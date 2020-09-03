package tools

// StringSliceInclude Returns +true+ if the given +obj+ is present in +arr+ (that is, if any
// element == +obj+), otherwise returns +false+.
func (t *Tool) StringSliceInclude(arr []string, obj string) bool {
	for _, element := range arr {
		if element == obj {
			return true
		}
	}
	return false
}
