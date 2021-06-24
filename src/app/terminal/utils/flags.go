// Package utils - Utility Package
package utils

import (
	"reflect"
	"strings"

	"github.com/fatih/structs"
	"github.com/urfave/cli/v2"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type ()

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// GenerateFlags - 지정된 Options들을 Application과 LocalCommand 등을 실행하는데 필요한 Flag 값으로 전환 (Usage Help 등으로 활용)
func GenerateFlags(options ...interface{}) (flags []cli.Flag, mappings map[string]string, err error) {
	mappings = make(map[string]string)

	for _, targetStruct := range options {
		o := structs.New(targetStruct)
		for _, field := range o.Fields() {
			flagName := field.Tag("flagName")
			if flagName == "" {
				continue
			}
			envName := "K3WT_" + strings.ToUpper(strings.Join(strings.Split(flagName, "-"), "_"))
			mappings[flagName] = field.Name()

			flagShortName := field.Tag("flagSName")
			if flagShortName != "" {
				flagName += ", " + flagShortName
			}

			flagDescription := field.Tag("flagDescribe")

			switch field.Kind() {
			case reflect.String:
				flags = append(flags, &cli.StringFlag{
					Name:    flagName,
					Value:   field.Value().(string),
					Usage:   flagDescription,
					EnvVars: []string{envName},
				})
			case reflect.Bool:
				flags = append(flags, &cli.BoolFlag{
					Name:    flagName,
					Usage:   flagDescription,
					EnvVars: []string{envName},
				})
			case reflect.Int:
				flags = append(flags, &cli.IntFlag{
					Name:    flagName,
					Value:   field.Value().(int),
					Usage:   flagDescription,
					EnvVars: []string{envName},
				})
			}
		}
	}

	return
}

// ApplyFlags - Flag로 설정된 값을 지정한 옵션 구조로 복사
func ApplyFlags(flags []cli.Flag, flagsMappingHint map[string]string, c *cli.Context, options ...interface{}) {
	// 대상 옵션 구조 생성
	objects := make([]*structs.Struct, len(options))

	for i, targetStruct := range options {
		objects[i] = structs.New(targetStruct)
	}

	// Flag와 Option Fields 정보를 기준으로 입력 값 설정
	for flagName, fieldName := range flagsMappingHint {
		if !c.IsSet(flagName) {
			continue
		}
		var field *structs.Field
		var ok bool
		for _, o := range objects {
			field, ok = o.FieldOk(fieldName)
			if ok {
				break
			}
		}
		if field == nil {
			continue
		}
		var val interface{}
		switch field.Kind() {
		case reflect.String:
			val = c.String(flagName)
		case reflect.Bool:
			val = c.Bool(flagName)
		case reflect.Int:
			val = c.Int(flagName)
		}
		field.Set(val)
	}
}
