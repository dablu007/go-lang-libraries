package utility

type MapUtil struct {
}

func (m MapUtil) GetKeyListFromKeyValueMap(keyMap map[int]bool) []int {
	keys := []int{}
	for k, _ := range keyMap {
		keys = append(keys, k)
	}
	return keys
}
