package libvalidate

import (
	"reflect"
	"shop-service/internal/util/liberr"
	"strings"
	"sync"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Global singleton state
var (
	once       sync.Once
	validate   *validator.Validate
	translator ut.Translator
)

func Validator() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())

		// Setup English Translator
		eng := en.New()
		uni := ut.New(eng, eng)

		translator, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(validate, translator)
	})
	return validate
}

// RegisterJSONTagField on Register JSON tag on validator error field
func RegisterJSONTagField(*validator.Validate) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
}

func ResolveError(err error, errCode string) error {
	if err == nil {
		return nil
	}

	if verr, ok := err.(validator.ValidationErrors); ok {
		vterr := verr.Translate(translator)

		errDetails := []*liberr.ErrorDetails{}
		for _, fieldErr := range verr {
			// Ignore first namespace / Struct name
			field := strings.SplitN(fieldErr.Namespace(), ".", 2)[1]

			errDetails = append(errDetails, &liberr.ErrorDetails{
				Code:    errCode,
				Message: vterr[fieldErr.Namespace()],
				Field:   field,
			})
		}

		return liberr.NewBaseError(errDetails...)
	}

	return liberr.NewTracerFromError(err)
}
