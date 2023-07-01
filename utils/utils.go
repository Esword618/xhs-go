package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GetUserAgent() string {
	uaList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.53 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return uaList[r.Intn(len(uaList))]
}

//	func ConvertInterfaceToMap(i interface{}) (map[string]string, error) {
//		if iMap, ok := i.(map[string]interface{}); ok {
//			result := make(map[string]string)
//			for key, value := range iMap {
//				if strValue, ok := value.(string); ok {
//					result[key] = strValue
//				} else {
//					return nil, fmt.Errorf("value for key %s is not a string", key)
//				}
//			}
//			return result, nil
//		}
//		return nil, fmt.Errorf("input is not a map[string]interface{}")
//	}
func ConvertInterfaceToMap(i interface{}) (map[string]string, error) {
	if iMap, ok := i.(map[string]string); ok {
		return iMap, nil
	}
	return nil, fmt.Errorf("input is not a map[string]string")
}
