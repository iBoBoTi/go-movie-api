package pkg

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iBoBoTi/go-movie-api/errors"
	"net/http"
	"strings"
)

func Decode(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindJSON(v); err != nil {
		e := &errors.Error{
			Status: http.StatusBadRequest,
		}
		if verr, ok := err.(validator.ValidationErrors); ok {
			errs := []string{}
			for _, fieldErr := range verr {
				errs = append(errs, fmt.Sprintf("%s is invalid: '%s'", fieldErr.Field(), fieldErr.Value()))
			}
			e.Message = strings.Join(errs, ";")
			return e
		}
		e.Message = err.Error()
		return e
	}
	return nil
}
