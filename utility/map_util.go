package utility

type MapUtility interface {
	GetKeyListFromKeyValueMap(keyMap map[int]bool) []int
}

type MapUtil struct {
}

func NewMapUtil() *MapUtil {
	return &MapUtil{}
}

func (m MapUtil) GetKeyListFromKeyValueMap(keyMap map[int]bool) []int {
	keys := []int{}
	for k, _ := range keyMap {
		keys = append(keys, k)
	}
	return keys
}
