package util

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"text/template"
)

const letterBytes1 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type stringUtil struct {
}

func (u *stringUtil) ArrayJoin(rows interface{}, fn func(index int) string) string {
	rv := reflect.ValueOf(rows)
	if rv.Len() == 0 {
		return ""
	}
	var result []string
	for i := 0; i < rv.Len(); i++ {
		r1 := fn(i)
		result = append(result, r1)
	}
	return strings.Join(result, ",")
}

func (u *stringUtil) GetMappingString(content string) map[string]string {
	var result map[string]string = make(map[string]string)
	if content == "" {
		return result
	}

	cs := strings.Split(content, ",")
	for _, c := range cs {
		result[c] = "1"
	}
	return result
}

func (u *stringUtil) IsFirstCharUpperCase(content string) bool {
	firstChar := content[0]
	if firstChar < 65 || firstChar > 90 {
		return false
	} else {
		return true
	}
}

func (u *stringUtil) RemoveEmptyRow(content string) string {
	return removeEmptyRowReg.ReplaceAllString(content, "")
}

func (u *stringUtil) GetRandomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (u *stringUtil) GetRandomStr1(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes1[rand.Int63()%int64(len(letterBytes1))]
	}
	return string(b)
}

func (u *stringUtil) FieldName(field string) string {
	fs := strings.Split(field, "_")
	var r1 []string
	for _, f := range fs {
		if f == "" {
			continue
		}
		r1 = append(r1, strings.Title(f))
	}
	return strings.Join(r1, "")
}

func (u *stringUtil) GetByTpl(tplStr string, target interface{}, funcMap template.FuncMap) (string, error) {
	var tpl *template.Template
	if funcMap == nil || len(funcMap) == 0 {
		tpl = template.Must(template.New(fmt.Sprintf("%s-%s", DateUtil.FormatNowByType(DatePattern2), StringUtil.GetRandomStr(5))).Parse(tplStr))
	} else {
		tpl = template.Must(template.New(fmt.Sprintf("%s-%s", DateUtil.FormatNowByType(DatePattern2), StringUtil.GetRandomStr(5))).Funcs(funcMap).Parse(tplStr))
	}

	buf := &bytes.Buffer{}

	err := tpl.Execute(buf, target)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (u *stringUtil) GetPhoneCode(phone string, area string) *countryCode {
	p := strings.TrimSpace(phone)
	if p[0:1] == "+" {
		l := strings.Index(phone, "-")
		if l > -1 {
			areaCode := p[0:l]
			phoneCode := p[l+1:]
			return &countryCode{
				AreaCode:  areaCode,
				PhoneCode: phoneCode,
			}
		} else {
			//+8526730
			//+852 9476
			//+8526730
			p = strings.ReplaceAll(p, " ", "")
			if strings.Index(p, "+852") > -1 {
				// 香港
				areaCode := p[0:4]
				phoneCode := p[4:]
				return &countryCode{
					AreaCode:  areaCode,
					PhoneCode: phoneCode,
				}
			} else if strings.Index(p, "+65") > -1 {
				// 新加坡
				areaCode := p[0:3]
				phoneCode := p[3:]
				return &countryCode{
					AreaCode:  areaCode,
					PhoneCode: phoneCode,
				}
			} else {
				// 其他都给前三位当区号
				areaCode := p[0:3]
				phoneCode := p[3:]
				return &countryCode{
					AreaCode:  areaCode,
					PhoneCode: phoneCode,
				}
			}
		}
	} else {
		var areaCode string
		if strings.Index(strings.ToLower(area), "hk") > -1 {
			areaCode = "+852"
		} else if strings.Index(strings.ToLower(area), "sin") > -1 {
			areaCode = "+65"
		}
		return &countryCode{
			AreaCode:  areaCode,
			PhoneCode: phone,
		}
	}
}

type countryCode struct {
	AreaCode  string
	PhoneCode string
}

var StringUtil stringUtil = stringUtil{}
