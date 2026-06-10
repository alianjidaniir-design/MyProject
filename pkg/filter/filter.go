package filter

import (
	"fmt"
	"strings"
)

type Filter struct {
	Con   string
	Value interface{}
	Op    string
}

func Filtering(filter ...Filter) (string, []interface{}) {

	var conditions []string
	var args []interface{}

	for _, fil := range filter {
		if fil.Value == nil || isZeroValue(fil.Value) {
			continue
		}

		if fil.Op == "" {
			fil.Op = "="
		}

		conditions = append(conditions, fmt.Sprintf("%s %s ?", fil.Con, fil.Op))
		args = append(args, fil.Value)

	}
	if len(conditions) == 0 {
		return "1=1", []interface{}{}

	}
	
	return strings.Join(conditions, " AND "), args

}

func isZeroValue(v interface{}) bool {
	switch val := v.(type) {
	case int, int8, int16, int32, int64:
		return val == 0
	case uint, uint8, uint16, uint32, uint64:
		return val == 0
	case *int64:
		return val == nil || *val == 0
	case float32, float64:
		return val == 0.0
	case string:
		return val == ""
	case *string:
		return val == nil || *val == ""
	default:
		return false
	}
}
