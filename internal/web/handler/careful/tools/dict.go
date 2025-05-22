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
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/tools"
	serviceTools "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/tools"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// CreateDictRequest 创建
type CreateDictRequest struct {
	Name      string `json:"name" binding:"required,max=100"`           // 字典名称
	Code      string `json:"code" binding:"required,max=100"`           // 字典编码
	Type      int    `json:"type" binding:"omitempty" default:"1"`      // 字典分类
	ValueType int    `json:"valueType" binding:"omitempty" default:"1"` // 字典值类型
	Sort      int    `json:"sort" binding:"omitempty" default:"1"`      // 排序
	Remark    string `json:"remark" binding:"omitempty,max=255"`        // 备注
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
	Delete(ctx *gin.Context)
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
	base.DELETE("/delete/:id", h.Delete)
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
		ctx.Set("internal", uid)
		zap.S().Error("获取用户失败", uid)
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	var req CreateDictRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
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
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典信息已存在", nil)
			return
		default:
			zap.L().Error("创建字典失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
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
	if id == "" {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, serviceTools.ErrDictNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典不存在", nil)
			return
		}
		zap.L().Error("删除字典失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
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
		ctx.Set("internal", uid)
		zap.S().Error("获取用户失败", uid)
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
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典信息已存在", nil)
			return
		default:
			zap.L().Error("创建字典失败", zap.Error(err))
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

	user, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceTools.ErrDictNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典不存在", nil)
			return
		}
		zap.L().Error("获取字典失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "获取成功", user)
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
		ctx.Set("internal", uid)
		zap.S().Error("获取用户失败", uid)
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
		Type:      dictType,
		ValueType: valueType,
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
		ctx.Set("internal", uid)
		zap.S().Error("获取用户失败", uid)
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
		Type:      dictType,
		ValueType: valueType,
	}

	list, err := h.svc.GetListAll(ctx, filter)
	if err != nil {
		zap.L().Error("获取列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}
