package utils

import (
	"github.com/tidwall/gjson"
)

func GetOperators() []gjson.Result {
	operatorsJson := RedisGet("data_source")
	return gjson.Parse(operatorsJson).Array()
}
