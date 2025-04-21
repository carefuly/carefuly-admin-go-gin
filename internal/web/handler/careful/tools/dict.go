/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/4/15 14:05:55
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
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/tools/dict"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/query/filters"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/xlsx"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type DictHandler interface {
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

type dictHandler struct {
	rely config.RelyConfig
	svc  tools.DictService
}

func NewDictHandler(rely config.RelyConfig, svc tools.DictService) DictHandler {
	return &dictHandler{
		rely: rely,
		svc:  svc,
	}
}

type DictRequest struct {
	Name      string              `json:"name" binding:"required,max=100"` // 字典名称
	Code      string              `json:"code" binding:"required,max=100"` // 字典编码
	Type      dict.TypeConst      `json:"type" default:"0"`                // 字典类型
	TypeValue dict.TypeValueConst `json:"typeValue" default:"0"`           // 字典类型值
	Version   int                 `json:"version"`                         // 版本
	Remark    string              `json:"remark" binding:"max=255"`        // 备注
}

type ImportDictRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (h *dictHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/dict")
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
// @Summary 创建字典
// @Description 创建字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param DictRequest body DictRequest true "字典信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/create [post]
// @Security LoginToken
func (h *dictHandler) Create(ctx *gin.Context) {
	uid, ok := ctx.MustGet("userId").(string)
	if !ok {
		ctx.Set("internal", uid)
		zap.S().Error("用户ID获取失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	var req DictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 校验参数
	_, err := dict.ConvertDictType(req.Type)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "type:"+err.Error(), nil)
		return
	}
	_, err = dict.ConvertDictTypeValue(req.TypeValue)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "typeValue:"+err.Error(), nil)
		return
	}

	err = h.svc.Create(ctx, domainTools.Dict{
		Dict: modelTools.Dict{
			CoreModels: models.CoreModels{
				Creator:  uid,
				Modifier: uid,
				Remark:   req.Remark,
			},
			Name:      req.Name,
			Code:      req.Code,
			Type:      req.Type,
			TypeValue: req.TypeValue,
		},
	})

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
	case errors.Is(err, tools.ErrDuplicateDictName):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典名称已存在", nil)
	case errors.Is(err, tools.ErrDuplicateDictCode):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典编码已存在", nil)
	case errors.Is(err, tools.ErrDuplicateDict):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典已存在", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("新增数据字典异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}

// Import
// @Summary 导入字典
// @Description 导入字典
// @Tags 系统工具/字典管理
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "文件(支持xlsx/csv格式)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/import [post]
// @Security LoginToken
func (h *dictHandler) Import(ctx *gin.Context) {
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
	read, err := xlsx.NewXlsxFile(filePath).ReadBySheet("字典模板")
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, xlsx.ErrOpenFile, nil)
		return
	}

	result := h.svc.Import(ctx, uid, read)
	msg := fmt.Sprintf("导入成功【成功导入【%d】条数据, 失败【%d】条数据】", result.SuccessCount, result.FailCount)

	response.NewResponse().SuccessResponse(ctx, msg, result)
}

// Delete
// @Summary 删除字典
// @Description 删除字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/delete/{id} [delete]
// @Security LoginToken
func (h *dictHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "参数错误", nil)
		return
	}

	err := h.svc.Delete(ctx, id)

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
	case errors.Is(err, tools.ErrDictNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "记录不存在", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("删除字典异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}

// BatchDelete
// @Summary 批量删除字典
// @Description 批量删除字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "ID数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/batchDelete [post]
// @Security LoginToken
func (h *dictHandler) BatchDelete(ctx *gin.Context) {
	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := h.svc.BatchDelete(ctx, ids)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("批量删除字典异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
}

// Update
// @Summary 更新字典
// @Description 更新字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "ID"
// @Param DictRequest body DictRequest true "字典信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/update/{id} [put]
// @Security LoginToken
func (h *dictHandler) Update(ctx *gin.Context) {
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

	var req DictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 校验参数
	_, err := dict.ConvertDictType(req.Type)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "type:"+err.Error(), nil)
		return
	}
	_, err = dict.ConvertDictTypeValue(req.TypeValue)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "typeValue:"+err.Error(), nil)
		return
	}

	err = h.svc.Update(ctx, id, domainTools.Dict{
		Dict: modelTools.Dict{
			CoreModels: models.CoreModels{
				Version:  req.Version,
				Modifier: uid,
				Remark:   req.Remark,
			},
			Name:      req.Name,
			Code:      req.Code,
			Type:      req.Type,
			TypeValue: req.TypeValue,
		},
	})

	switch {
	case err == nil:
		response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
	case errors.Is(err, tools.ErrDuplicateDictName):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典名称已存在", nil)
	case errors.Is(err, tools.ErrDuplicateDictCode):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典编码已存在", nil)
	case errors.Is(err, tools.ErrDuplicateDict):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典已存在", nil)
	case errors.Is(err, tools.ErrDictNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "记录不存在", nil)
	case errors.Is(err, tools.ErrDictVersionInconsistency):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.L().Error("新增数据字典异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}

// GetById
// @Summary 根据ID获取字典
// @Description 根据ID获取字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/getById/{id} [get]
// @Security LoginToken
func (h *dictHandler) GetById(ctx *gin.Context) {
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
	case errors.Is(err, tools.ErrDictNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "记录不存在", nil)
	case errors.Is(err, tools.ErrDictRecordNotFound):
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "记录不存在", nil)
	default:
		ctx.Set("internal", err.Error())
		zap.S().Error("根据Id查询字典异常", err)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}
}

// GetListPage
// @Summary 分页获取字典
// @Description 分页获取字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "每页数量" default(10)
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "字典名称"
// @Param code query string false "字典编码"
// @Param type query int false "字典类型" default(-1)
// @Param typeValue query int false "字典数据类型" default(-1)
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/listPage [get]
// @Security LoginToken
func (h *dictHandler) GetListPage(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	creator := ctx.DefaultQuery("creator", "")
	modifier := ctx.DefaultQuery("modifier", "")
	status, _ := strconv.ParseBool(ctx.DefaultQuery("status", "true"))
	name := ctx.DefaultQuery("name", "")
	code := ctx.DefaultQuery("code", "")
	dictType, _ := strconv.Atoi(ctx.Query("type"))
	typeValue, _ := strconv.Atoi(ctx.Query("typeValue"))

	list, total, err := h.svc.GetListPage(ctx, domainTools.DictFilter{
		Filters: filters.Filters{
			Creator:  creator,
			Modifier: modifier,
			Status:   status,
		},
		Pagination: filters.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		Name:      name,
		Code:      code,
		Type:      dictType,
		TypeValue: typeValue,
	})

	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("分页查询列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", gin.H{
		"list":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GetListAll
// @Summary 获取所有字典
// @Description 获取所有字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "字典名称"
// @Param code query string false "字典编码"
// @Param type query int false "字典类型" default(-1)
// @Param typeValue query int false "字典数据类型" default(-1)
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/listAll [get]
// @Security LoginToken
func (h *dictHandler) GetListAll(ctx *gin.Context) {
	creator := ctx.DefaultQuery("creator", "")
	modifier := ctx.DefaultQuery("modifier", "")
	status, _ := strconv.ParseBool(ctx.DefaultQuery("status", "true"))
	name := ctx.DefaultQuery("name", "")
	code := ctx.DefaultQuery("code", "")
	dictType, _ := strconv.Atoi(ctx.DefaultQuery("type", "-1"))
	typeValue, _ := strconv.Atoi(ctx.DefaultQuery("typeValue", "-1"))

	list, err := h.svc.GetListAll(ctx, domainTools.DictFilter{
		Filters: filters.Filters{
			Creator:  creator,
			Modifier: modifier,
			Status:   status,
		},
		Name:      name,
		Code:      code,
		Type:      dictType,
		TypeValue: typeValue,
	})

	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("查询所有列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}
