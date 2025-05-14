/**
 * Description：
 * FileName：validator.go
 * Author：CJiaの用心
 * Create：2025/2/19 16:06:49
 * Remark：
 */

package validate

import (
	"errors"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"strings"
)

type ValidatorError struct {
	trans ut.Translator
}

func NewValidatorError(trans ut.Translator) *ValidatorError {
	return &ValidatorError{
		trans: trans,
	}
}

func (v *ValidatorError) HandleValidatorError(ctx *gin.Context, err error) {
	// 检查是否是因为EOF导致的错误
	if err == io.EOF {
		// 参数为空
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "请求参数为空", nil)
		return
	}
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)
	if !ok {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, v.removeTopStruct(errs.Translate(v.trans)), nil)
	return
}

func (v *ValidatorError) removeTopStruct(filed map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range filed {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}
