/**
 * Description：
 * FileName：dict_type.go
 * Author：CJiaの用心
 * Create：2025/4/17 10:34:29
 * Remark：
 */

package tools

import (
	"errors"
	"fmt"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/tools/dictType"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/xlsx"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"time"
)

type DictTypeHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	Create(ctx *gin.Context)
	Import(ctx *gin.Context)
	Delete(ctx *gin.Context)
	BatchDelete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetListPage(ctx *gin.Context)
	GetListAll(ctx *gin.Context)
}

type dictTypeHandler struct {
	rely config.RelyConfig
	svc  tools.DictTypeService
}

func NewDictTypeHandler(rely config.RelyConfig, svc tools.DictTypeService) DictTypeHandler {
	return &dictTypeHandler{
		rely: rely,
		svc:  svc,
	}
}

type DictTypeRequest struct {
	Name      string                `json:"name" binding:"required,max=50"`     // 字典信息名称
	StrValue  string                `json:"strValue"`                           // 字符串-字典信息值
	IntValue  int64                 `json:"intValue"`                           // 整型-字典信息值
	BoolValue bool                  `json:"boolValue"`                          // 布尔-字典信息值
	DictTag   dictType.DictTagConst `json:"dictTag" default:"primary"`          // 标签类型
	DictColor string                `json:"dictColor"`                          // 标签颜色
	DictId    string                `json:"dict_id" binding:"required,max=100"` // 字典ID
	Version   int                   `json:"version"`                            // 版本
	Remark    string                `json:"remark" binding:"max=255"`           // 备注
}

type ImportDictTypeRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (h *dictTypeHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/dictType")
	base.POST("/create", h.Create)
	base.POST("/import", h.Import)
	base.DELETE("/delete/:id", h.Delete)
	base.POST("/batchDelete", h.BatchDelete)
	base.PUT("/update/:id", h.Update)
	base.GET("/getById/:id", h.GetById)
	base.GET("/listPage", h.GetListPage)
	base.GET("/listAll", h.GetListAll)
}

// Create
// @Summary 创建字典信息
// @Description 创建字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param DictTypeRequest body DictTypeRequest true "字典信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/create [post]
// @Security LoginToken
func (h *dictTypeHandler) Create(ctx *gin.Context) {
	uid, ok := ctx.MustGet("userId").(string)
	if !ok {
		ctx.Set("internal", uid)
		zap.S().Error("用户ID获取失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	var req DictTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 校验参数
	_, err := dictType.ConvertDictTag(req.DictTag)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "dictTag:"+err.Error(), nil)
		return
	}

	err = h.svc.Create(ctx, domainTools.DictType{
		DictType: modelTools.DictType{
			CoreModels: models.CoreModels{
				Creator:  uid,
				Modifier: uid,
				Remark:   req.Remark,
			},
			Name:      req.Name,
			DictTag:   req.DictTag,
			DictColor: req.DictColor,
			DictId:    req.DictId,
		},
		StrValue:  req.StrValue,
		IntValue:  req.IntValue,
		BoolValue: req.BoolValue,
	})

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
	case errors.Is(err, tools.ErrDictNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典不存在", nil)
	case errors.Is(err, tools.ErrDictRecordNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典不存在", nil)
	case errors.Is(err, tools.ErrNotSupportedTypeValue):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "不支持的字典类型", nil)
	case errors.Is(err, tools.ErrDuplicateDictType):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典类型已存在", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("新增字典信息异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}

// Import
// @Summary 导入字典信息
// @Description 导入字典信息
// @Tags 系统工具/字典信息管理
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "文件(支持xlsx/csv格式)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/import [post]
// @Security LoginToken
func (h *dictTypeHandler) Import(ctx *gin.Context) {
	uid, ok := ctx.MustGet("userId").(string)
	if !ok {
		ctx.Set("internal", uid)
		zap.S().Error("用户ID获取失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	var req ImportDictRequest
	if err := ctx.ShouldBind(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 保存导入的文件信息
	format := time.Now().Format("2006-01-02")
	filePath := "./uploads/" + format + "/" + req.File.Filename
	if err := ctx.SaveUploadedFile(req.File, filePath); err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "保存文件失败", nil)
		return
	}

	// 读取Excel文件
	read, err := xlsx.NewXlsxFile(filePath).ReadBySheet("字典信息模板")
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, xlsx.ErrOpenFile, nil)
		return
	}

	result := h.svc.Import(ctx, uid, read)
	msg := fmt.Sprintf("导入成功【成功导入【%d】条数据, 失败【%d】条数据】", result.SuccessCount, result.FailCount)

	response.NewResponse().SuccessResponse(ctx, msg, result)
}

// Delete
// @Summary 删除字典信息
// @Description 删除字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/delete/{id} [delete]
// @Security LoginToken
func (h *dictTypeHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "参数错误", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "参数校验", id)
	// err := h.svc.Delete(ctx, id)
	//
	// switch {
	// case err == nil:
	// 	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
	// case errors.Is(err, tools.ErrDictNotFound):
	// 	response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "记录不存在", nil)
	// default:
	// 	ctx.Set("internal", err.Error())
	// 	zap.L().Error("删除字典异常", zap.Error(err))
	// 	response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	// }
}

// BatchDelete
// @Summary 批量删除字典信息
// @Description 批量删除字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "ID数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/batchDelete [post]
// @Security LoginToken
func (h *dictTypeHandler) BatchDelete(ctx *gin.Context) {
	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "参数校验", ids)

	// err := h.svc.BatchDelete(ctx, ids)
	// if err != nil {
	// 	ctx.Set("internal", err.Error())
	// 	zap.L().Error("批量删除字典异常", zap.Error(err))
	// 	response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	// 	return
	// }
	//
	// response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
}

// Update
// @Summary 更新字典信息
// @Description 更新字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "ID"
// @Param DictTypeRequest body DictTypeRequest true "字典信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/update/{id} [put]
// @Security LoginToken
func (h *dictTypeHandler) Update(ctx *gin.Context) {
	uid, ok := ctx.MustGet("userId").(string)
	if !ok {
		ctx.Set("internal", uid)
		zap.S().Error("用户ID获取失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "参数错误", nil)
		return
	}

	var req DictTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 校验参数
	_, err := dictType.ConvertDictTag(req.DictTag)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "dictTag:"+err.Error(), nil)
		return
	}

	err = h.svc.Update(ctx, id, domainTools.DictType{
		DictType: modelTools.DictType{
			CoreModels: models.CoreModels{
				Version:  req.Version,
				Modifier: uid,
				Remark:   req.Remark,
			},
			Name:      req.Name,
			DictTag:   req.DictTag,
			DictColor: req.DictColor,
			DictId:    req.DictId,
		},
		StrValue:  req.StrValue,
		IntValue:  req.IntValue,
		BoolValue: req.BoolValue,
	})

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
	case errors.Is(err, tools.ErrDictNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典不存在", nil)
	case errors.Is(err, tools.ErrDictRecordNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典不存在", nil)
	case errors.Is(err, tools.ErrNotSupportedTypeValue):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "不支持的字典类型", nil)
	case errors.Is(err, tools.ErrDuplicateDictType):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典类型已存在", nil)
	case errors.Is(err, tools.ErrDictTypeNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "记录不存在", nil)
	case errors.Is(err, tools.ErrDictTypeVersionInconsistency):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("新增字典信息异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}

// GetById
// @Summary 根据ID获取字典信息
// @Description 根据ID获取字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/getById/{id} [get]
// @Security LoginToken
func (h *dictTypeHandler) GetById(ctx *gin.Context) {
	uid, ok := ctx.MustGet("userId").(string)
	if !ok {
		ctx.Set("internal", uid)
		zap.S().Error("用户ID获取失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "参数错误", nil)
		return
	}

	detail, err := h.svc.GetById(ctx, id)

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
	case errors.Is(err, tools.ErrDictTypeNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "记录不存在", nil)
	case errors.Is(err, tools.ErrDictTypeRecordNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "记录不存在", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.S().Error("根据Id查询字典信息异常", err)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}

// GetListPage
// @Summary 分页获取字典信息
// @Description 分页获取字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "每页数量" default(10)
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "字典信息名称"
// @Param dictTag query string false "标签类型" default(primary)
// @Param dictId query string false "字典ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/listPage [get]
// @Security LoginToken
func (h *dictTypeHandler) GetListPage(ctx *gin.Context) {

}

// GetListAll
// @Summary 获取所有字典信息
// @Description 获取所有字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "字典信息名称"
// @Param dictTag query string false "标签类型" default(primary)
// @Param dictId query string false "字典ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/listAll [get]
// @Security LoginToken
func (h *dictTypeHandler) GetListAll(ctx *gin.Context) {

}
