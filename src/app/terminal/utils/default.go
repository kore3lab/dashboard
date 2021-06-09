// Package utils - Utility Package
package utils

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/fatih/structs"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type ()

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// ApplyDefaultValues - 지정된 구조체 형식에 Structure Tags로 지정된 기본 값 설정
func ApplyDefaultValues(targetStructure interface{}) (err error) {
	// 구조체 인스턴스 생성
	o := structs.New(targetStructure)

	// 필드들에 대한 Structure Tag를 통한 기본 값 설정
	for _, field := range o.Fields() {
		defaultValue := field.Tag("default")
		if "" == defaultValue {
			continue
		}

		var val interface{}
		switch field.Kind() {
		case reflect.String:
			val = defaultValue
		case reflect.Bool:
			if defaultValue == "true" {
				val = true
			} else if defaultValue == "false" {
				val = false
			} else {
				return fmt.Errorf("invalid bool expression: %v, use true/false", defaultValue)
			}
		case reflect.Int:
			val, err = strconv.Atoi(defaultValue)
			if err != nil {
				return err
			}
		default:
			val = field.Value()
		}
		field.Set(val)
	}
	return nil
}
