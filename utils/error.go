package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// เป็น Extention Method โดยจะส่งจาก Business Logic ไปยัง Presentation Layer โดยจะ Conform ตาม Type Error
type Err_Handler struct {
	Code    int         `json:"code" bson:"code"`
	Status  bool        `json:"status" bson:"status"`
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data" bson:"data"`
}

func (e Err_Handler) Error() string {
	return e.Message
}

func CheckErrorMessage(responseData interface{}) (result string, err error) {

	for k, v := range responseData.(map[string]interface{}) {

		if k == "errorMessage" {
			// check type errorMessage is string
			if reflect.ValueOf(v).Kind() == reflect.String {
				return fmt.Sprintf("%v", v), nil
			}

			// check type errorMessage is map
			if reflect.ValueOf(v).Kind() == reflect.Map {
				listErrorData := []string{}
				for _, valList := range v.(map[string]interface{}) {
					for _, errorData := range valList.([]interface{}) {
						listErrorData = append(listErrorData, fmt.Sprintf("%v", errorData))
					}
				}

				return strings.Join(listErrorData[:], ","), nil
			}

			// check errorMessage is nil
			if v == nil {
				return "no error", nil
			}

			return "", fmt.Errorf("invalid type errorMessage")
		}
	}

	return "no error", nil
}
