package global

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

func ApiInitValidator() {
	validate := validator.New()

	// 定义不同的语言翻译
	enT := en.New()
	cn := zh_Hans_CN.New()

	uni := ut.New(enT, cn)

	enTrans, _ := uni.GetTranslator("en")
	zhTrans, _ := uni.GetTranslator("zh_Hans_CN")

	err := zh_translations.RegisterDefaultTranslations(validate, zhTrans)
	if err != nil {
		panic(err)
	}
	err = en_translations.RegisterDefaultTranslations(validate, enTrans)
	if err != nil {
		panic(err)
	}

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	Validator.Validate = validate
	Validator.UT = uni // 存储 Universal Translator
	Validator.VTrans = zhTrans

	Validator.ValidStruct = func(ctx *gin.Context, i interface{}) []string {
		err := Validator.Validate.Struct(i)
		lang := ctx.GetHeader("Accept-Language")
		if lang == "" {
			lang = Config.Lang
		}
		trans := getTranslatorForLang(lang)
		errList := make([]string, 0, 10)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				errList = append(errList, err.Error())
				return errList
			}
			for _, err2 := range err.(validator.ValidationErrors) {
				errList = append(errList, err2.Translate(trans))
			}
		}
		return errList
	}
	Validator.ValidVar = func(ctx *gin.Context, field interface{}, tag string) []string {
		err := Validator.Validate.Var(field, tag)
		lang := ctx.GetHeader("Accept-Language")
		if lang == "" {
			lang = Config.Lang
		}
		trans := getTranslatorForLang(lang)
		errList := make([]string, 0, 10)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				errList = append(errList, err.Error())
				return errList
			}
			for _, err2 := range err.(validator.ValidationErrors) {
				errList = append(errList, err2.Translate(trans))
			}
		}
		return errList
	}

}
func getTranslatorForLang(lang string) ut.Translator {
	switch lang {
	case "zh_CN":
		fallthrough
	case "zh-CN":
		fallthrough
	case "zh":
		trans, _ := Validator.UT.GetTranslator("zh_Hans_CN")
		return trans
	case "en":
		fallthrough
	default:
		trans, _ := Validator.UT.GetTranslator("en")
		return trans
	}
}