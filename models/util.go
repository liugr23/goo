package models

import "reflect"

type UserSessionInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ConvertToInt64(number interface{}) int64 {
	if reflect.TypeOf(number).String() == "int" {
		return int64(number.(int))
	}
	return number.(int64)
}
