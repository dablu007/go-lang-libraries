package utility

import (
	. "flow/utility"
	"reflect"
	"testing"
)

func TestMapUtil_GetKeyListFromKeyValueMap(t *testing.T) {
	type args struct {
		keyMap map[int]bool
	}
	tests := []struct {
		name string
		m    MapUtil
		args args
		want []int
	}{
		{name: "TestWithEmptyMap", args: args{keyMap: make(map[int]bool)}, want: []int{}},
		{name: "TestWithNonEmptyMap", args: args{map[int]bool{1: true, 2: true, 3: true}}, want: []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MapUtil{}
			if got := m.GetKeyListFromKeyValueMap(tt.args.keyMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapUtil.GetKeyListFromKeyValueMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
