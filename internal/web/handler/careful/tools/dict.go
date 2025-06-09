/**
 * Description：
 * FileName：dict.go
 * Author：CJiaの用心
 * Create：2025/5/14 16:43:12
 * Remark：
 */

package tools

import (
	"errors"
	"fmt"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/tools"
	serviceTools "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/tools/dict"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/enumconv"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/xlsx"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// CreateDictRequest 创建
type CreateDictRequest struct {
	Name      string              `json:"name" binding:"required,max=100"`           // 字典名称
	Code      string              `json:"code" binding:"required,max=100"`           // 字典编码
	Type      dict.TypeConst      `json:"type" binding:"omitempty" default:"1"`      // 字典分类
	ValueType dict.TypeValueConst `json:"valueType" binding:"omitempty" default:"1"` // 字典值类型
	Sort      int                 `json:"sort" binding:"omitempty" default:"1"`      // 排序
	Remark    string              `json:"remark" binding:"omitempty,max=255"`        // 备注
}

// ImportDictRequest 导入
type ImportDictRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// UpdateDictRequest 更新
type UpdateDictRequest struct {
	Id      string `json:"id" binding:"required"`                // 主键ID
	Code    string `json:"code" binding:"required,max=100"`      // 字典编码
	Sort    int    `json:"sort" binding:"omitempty" default:"1"` // 排序
	Version int    `json:"version" binding:"omitempty"`          // 版本
	Remark  string `json:"remark" binding:"omitempty,max=255"`   // 备注
}

// DictListPageResponse 列表分页响应
type DictListPageResponse struct {
	List     []domainTools.Dict `json:"list"`     // 列表
	Total    int64              `json:"total"`    // 总数
	Page     int                `json:"page"`     // 页码
	PageSize int                `json:"pageSize"` // 每页数量
}

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
	rely    config.RelyConfig
	svc     tools.DictService
	userSvc system.UserService
}

func NewDictHandler(rely config.RelyConfig, svc tools.DictService, userSvc system.UserService) DictHandler {
	return &dictHandler{
		rely:    rely,
		svc:     svc,
		userSvc: userSvc,
	}
}

// RegisterRoutes 注册路由
func (h *dictHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/dict")
	base.POST("/create", h.Create)
	base.POST("/import", h.Import)
	base.DELETE("/delete/:id", h.Delete)
	base.POST("/delete/batchDelete", h.BatchDelete)
	base.PUT("/update", h.Update)
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
// @Param CreateDictRequest body CreateDictRequest true "请求"
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

	user, err := h.userSvc.GetById(ctx, uid)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.S().Error("获取用户失败", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	var req CreateDictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 校验参数
	typeValidValues := []string{"普通字典", "系统字典", "枚举字典"}
	converter := enumconv.NewEnumConverter(dict.TypeMapping, dict.TypeImportMapping, typeValidValues, "字典分类")
	_, err = converter.FromEnum(req.Type)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}
	typeValueValidValues := []string{"字符串", "整型", "布尔"}
	typeValueConverter := enumconv.NewEnumConverter(dict.TypeValueMapping, dict.TypeValueImportMapping, typeValueValidValues, "字典值类型")
	_, err = typeValueConverter.FromEnum(req.ValueType)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 转换为领域模型
	domain := domainTools.Dict{
		Dict: modelTools.Dict{
			CoreModels: models.CoreModels{
				Sort:       req.Sort,
				Creator:    uid,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Name:      req.Name,
			Code:      req.Code,
			Type:      req.Type,
			ValueType: req.ValueType,
		},
	}

	if err := h.svc.Create(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceTools.ErrDictNameDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典名称已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrDictCodeDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典编码已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrDictDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据字典已存在", nil)
			return
		default:
			zap.L().Error("创建数据字典失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
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

	user, err := h.userSvc.GetById(ctx, uid)
	if err != nil {
		ctx.Set("internal", err)
		zap.S().Error("获取用户失败", err)
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

	result := h.svc.Import(ctx, uid, user.DeptId, read)
	msg := fmt.Sprintf("导入成功【成功导入【%d】条数据, 失败【%d】条数据】", result.SuccessCount, result.FailCount)

	response.NewResponse().SuccessResponse(ctx, msg, result)
}

// Delete
// @Summary 删除字典
// @Description 删除指定id字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/delete/{id} [delete]
// @Security LoginToken
func (h *dictHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, serviceTools.ErrDictNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典不存在", nil)
			return
		}
		ctx.Set("internal", err.Error())
		zap.L().Error("删除字典失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
}

// BatchDelete
// @Summary 批量删除字典
// @Description 批量删除字典
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "id数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/delete/batchDelete [post]
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
// @Description 更新字典信息
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param UpdateDictRequest body UpdateDictRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/update [put]
// @Security LoginToken
func (h *dictHandler) Update(ctx *gin.Context) {
	uid, ok := ctx.MustGet("userId").(string)
	if !ok {
		ctx.Set("internal", uid)
		zap.S().Error("用户ID获取失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	user, err := h.userSvc.GetById(ctx, uid)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.S().Error("获取用户失败", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	var req UpdateDictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainTools.Dict{
		Dict: modelTools.Dict{
			CoreModels: models.CoreModels{
				Id:         req.Id,
				Sort:       req.Sort,
				Version:    req.Version,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Code: req.Code,
		},
	}

	if err := h.svc.Update(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceTools.ErrDictNameDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典名称已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrDictCodeDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典编码已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrDictDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据字典已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrDictVersionInconsistency):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
			return
		default:
			zap.L().Error("更新数据字典失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
}

// GetById
// @Summary 获取字典
// @Description 获取指定id字典信息
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} domainTools.Dict
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/getById/{id} [get]
// @Security LoginToken
func (h *dictHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户ID不能为空", nil)
		return
	}

	detail, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceTools.ErrDictNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典不存在", nil)
			return
		}
		zap.L().Error("获取字典失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
}

// GetListPage
// @Summary 获取字典分页列表
// @Description 获取字典分页列表
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
// @Param type query int true "字典分类" default(0)
// @Param valueType query int true "字典值类型" default(0)
// @Success 200 {object} DictListPageResponse
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/listPage [get]
// @Security LoginToken
func (h *dictHandler) GetListPage(ctx *gin.Context) {
	uid, ok := ctx.MustGet("userId").(string)
	if !ok {
		ctx.Set("internal", uid)
		zap.S().Error("用户ID获取失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	user, err := h.userSvc.GetById(ctx, uid)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.S().Error("获取用户失败", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	creator := ctx.DefaultQuery("creator", "")
	modifier := ctx.DefaultQuery("modifier", "")
	status, _ := strconv.ParseBool(ctx.DefaultQuery("status", "true"))

	name := ctx.DefaultQuery("name", "")
	code := ctx.DefaultQuery("code", "")
	dictType, _ := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	valueType, _ := strconv.Atoi(ctx.DefaultQuery("valueType", "0"))

	filter := domainTools.DictFilter{
		Pagination: filters.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status:    status,
		Name:      name,
		Code:      code,
		Type:      dict.TypeConst(dictType),
		ValueType: dict.TypeValueConst(valueType),
	}

	list, total, err := h.svc.GetListPage(ctx, filter)
	if err != nil {
		zap.L().Error("获取分页列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", DictListPageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetListAll
// @Summary 获取所有字典
// @Description 获取所有字典列表
// @Tags 系统工具/字典管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "字典名称"
// @Param code query string false "字典编码"
// @Param type query int true "字典分类" default(0)
// @Param valueType query int true "字典值类型" default(0)
// @Success 200 {array} []domainTools.Dict
// @Failure 400 {object} response.Response
// @Router /v1/tools/dict/listAll [get]
// @Security LoginToken
func (h *dictHandler) GetListAll(ctx *gin.Context) {
	uid, ok := ctx.MustGet("userId").(string)
	if !ok {
		ctx.Set("internal", uid)
		zap.S().Error("用户ID获取失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	user, err := h.userSvc.GetById(ctx, uid)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.S().Error("获取用户失败", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	creator := ctx.DefaultQuery("creator", "")
	modifier := ctx.DefaultQuery("modifier", "")
	status, _ := strconv.ParseBool(ctx.DefaultQuery("status", "true"))
	name := ctx.DefaultQuery("name", "")
	code := ctx.DefaultQuery("code", "")
	dictType, _ := strconv.Atoi(ctx.DefaultQuery("type", "0"))
	valueType, _ := strconv.Atoi(ctx.DefaultQuery("valueType", "0"))

	filter := domainTools.DictFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status:    status,
		Name:      name,
		Code:      code,
		Type:      dict.TypeConst(dictType),
		ValueType: dict.TypeValueConst(valueType),
	}

	list, err := h.svc.GetListAll(ctx, filter)
	if err != nil {
		zap.L().Error("获取列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}
