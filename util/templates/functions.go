package templates

// ActionUsed build template function which checks if action is used
func ActionUsed(actions []string) interface{} {
	return func(input string) bool {
		for _, item := range actions {
			if item == input {
				return true
			}
		}

		return false
	}
}
