package lib

import "encoding/json"

func STM(a any) map[string]interface{} {
	var inInterface map[string]interface{}
	ms, _ := json.Marshal(a)
	json.Unmarshal(ms, &inInterface)
	return inInterface
}
