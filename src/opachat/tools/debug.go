package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RoomsDebugType struct {
	Cap  string           `json:"cap"`
	List []RoomsDebugType `json:"list,omitempty"`
}

// DebugJ pretty print structures
func DebugJ(v interface{}, echo bool, px, inde string) string {
	b, err := json.MarshalIndent(v, px, inde)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if echo {
		fmt.Printf(px+"%s\n", string(b))
		return ""
	}

	return string(b)
}

// DebugB print response body
func DebugB(r *http.Response) {
	b, _ := io.ReadAll(r.Body)
	fmt.Println(string(b))
}
