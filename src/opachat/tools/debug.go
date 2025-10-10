package tools

import (
	"encoding/json"
	"fmt"
)

func ShowJSON(v any, echos ...bool) string {
	echo := true
	if len(echos) > 0 {
		echo = echos[0]
	}

	b, err := json.MarshalIndent(v, "", "  ")

	if err == nil {
		if echo {
			fmt.Println(string(b))
			return ""
		}

		return string(b)
	}

	fmt.Println(err)
	return ""
}
