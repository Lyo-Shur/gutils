package gutils

import (
	"github.com/Lyo-Shur/gutils/api"
	"github.com/Lyo-Shur/gutils/bean"
	"github.com/Lyo-Shur/gutils/cache"
	"github.com/Lyo-Shur/gutils/config/properties"
	"github.com/Lyo-Shur/gutils/config/read"
	"github.com/Lyo-Shur/gutils/convert"
	"github.com/Lyo-Shur/gutils/crypto"
	"github.com/Lyo-Shur/gutils/file"
	"github.com/Lyo-Shur/gutils/task"
	"github.com/Lyo-Shur/gutils/ticket"
	"github.com/Lyo-Shur/gutils/validator"
	"strings"
	"time"
)

// api
type CodeModeDTO = api.CodeModeDTO

func JsonCodeModeDTO(code int64, message string, data interface{}) string {
	return (&api.CodeModeDTO{
		Code:    code,
		Message: message,
		Data:    data,
	}).ToJson()
}

// bean
type Factory = bean.Factory

func GetBeanFactory() *Factory {
	return bean.GetFactory()
}

// cache
type Data = cache.Data
type Cache = cache.Cache
type CacheHolder = cache.Holder

func GetCacheHolder() *CacheHolder {
	return cache.GetHolder()
}

// config
func LoadConfig(path string) map[string]string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return LoadConfigFromHTTP(path)
	}
	return LoadConfigFromFile(path)
}
func LoadConfigFromFile(path string) map[string]string {
	content := read.File(path)
	return properties.Load(content)
}
func LoadConfigFromHTTP(path string) map[string]string {
	content := read.HTTP(path)
	return properties.Load(content)
}

// convert
func ToBigHump(s string) string {
	return convert.ToBigHump(s)
}
func ToSmallHump(s string) string {
	return convert.ToSmallHump(s)
}
func ToUnderline(s string) string {
	return convert.ToUnderline(s)
}
func MapBindToStruct(m map[string]string, v interface{}) error {
	return convert.MapBindToStruct(m, v)
}

// crypto
func EncodeBase64(data string) string {
	return crypto.EncodeBase64(data)
}
func DecodeBase64(data string) string {
	return crypto.DecodeBase64(data)
}
func EncodeMD5(data string) string {
	return crypto.EncodeMD5(data)
}

// file
func ReadFile(path string) string {
	return file.Read(path)
}

// task
func Run(d time.Duration, f func()) {
	task.Run(d, f)
}

// ticket
type TicketHolder = ticket.Holder
type TicketCacheHolder = ticket.CacheHolder

// validator
type ValidatorConfig = validator.Config
type ValidatorHelper = validator.Helper

func GetValidatorHelper(params map[string]string, configJson string) *ValidatorHelper {
	return validator.GetHelper(params, configJson)
}

type Rule = validator.Rule
type Required = validator.Required
type Ban = validator.Ban
type Length = validator.Length
type Range = validator.Range
type DateTime = validator.DateTime
type Regexp = validator.Regexp
