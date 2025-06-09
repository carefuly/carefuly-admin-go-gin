/**
 * Description：
 * FileName：menu_column.go
 * Author：CJiaの用心
 * Create：2025/6/9 14:43:31
 * Remark：
 */

package system

import (
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// CreateMenuColumnRequest 创建
type CreateMenuColumnRequest struct {
	Title  string `json:"title" binding:"required,max=64"`        // 标题
	Field  string `json:"field" binding:"required,max=64"`        // 字段名
	Width  int    `json:"width" binding:"required" default:"150"` // 宽度
	MenuId string `json:"menuId" binding:"required,max=100"`      // 菜单ID
	Sort   int    `json:"sort" binding:"omitempty" default:"1"`   // 排序
	Remark string `json:"remark" binding:"omitempty,max=255"`     // 备注
}

// UpdateMenuColumnRequest 更新
type UpdateMenuColumnRequest struct {
	Id      string `json:"id" binding:"required"`                  // 主键ID
	Title   string `json:"title" binding:"required,max=64"`        // 标题
	Field   string `json:"field" binding:"required,max=64"`        // 字段名
	Width   int    `json:"width" binding:"required" default:"150"` // 宽度
	Sort    int    `json:"sort" binding:"omitempty" default:"1"`   // 排序
	Version int    `json:"version" binding:"omitempty"`            // 版本
	Remark  string `json:"remark" binding:"omitempty,max=255"`     // 备注
}

// MenuColumnListPageResponse 列表分页响应
type MenuColumnListPageResponse struct {
	List     []domainSystem.MenuColumn `json:"list"`     // 列表
	Total    int64                     `json:"total"`    // 总数
	Page     int                       `json:"page"`     // 页码
	PageSize int                       `json:"pageSize"` // 每页数量
}

type MenuColumnHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	BatchDelete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetListPage(ctx *gin.Context)
	GetListAll(ctx *gin.Context)
}

type menuColumnHandler struct {
	rely    config.RelyConfig
	svc     serviceSystem.MenuColumnService
	userSvc serviceSystem.UserService
}

func NewMenuColumnHandler(rely config.RelyConfig, svc serviceSystem.MenuColumnService, userSvc serviceSystem.UserService) MenuColumnHandler {
	return &menuColumnHandler{
		rely:    rely,
		svc:     svc,
		userSvc: userSvc,
	}
}

// RegisterRoutes 注册路由
func (h *menuColumnHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/menuColumn")
	base.POST("/create", h.Create)
	base.DELETE("/delete/:id", h.Delete)
	base.POST("/delete/batchDelete", h.BatchDelete)
	base.PUT("/update", h.Update)
	base.GET("/getById/:id", h.GetById)
	base.GET("/listPage", h.GetListPage)
	base.GET("/listAll", h.GetListAll)
}

// Create
// @Summary 创建菜单数据列
// @Description 创建菜单数据列
// @Tags 系统管理/菜单数据列管理
// @Accept application/json
// @Produce application/json
// @Param CreateMenuColumnRequest body CreateMenuColumnRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/menuColumn/create [post]
// @Security LoginToken
func (h *menuColumnHandler) Create(ctx *gin.Context) {
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

	var req CreateMenuColumnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainSystem.MenuColumn{
		MenuColumn: modelSystem.MenuColumn{
			CoreModels: models.CoreModels{
				Sort:       req.Sort,
				Creator:    uid,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Title:  req.Title,
			Field:  req.Field,
			MenuId: req.MenuId,
		},
	}

	if err := h.svc.Create(ctx, domain); err != nil {
		zap.L().Error("创建菜单数据列异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
}

// Delete
// @Summary 删除菜单数据列
// @Description 删除指定id菜单数据列
// @Tags 系统管理/菜单数据列管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/menuColumn/delete/{id} [delete]
// @Security LoginToken
func (h *menuColumnHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, serviceSystem.ErrMenuColumnNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "菜单数据列不存在", nil)
			return
		}
		ctx.Set("internal", err.Error())
		zap.L().Error("删除菜单数据列失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
}

// BatchDelete
// @Summary 批量删除菜单数据列
// @Description 批量删除菜单数据列
// @Tags 系统管理/菜单数据列管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "id数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/menuColumn/delete/batchDelete [post]
// @Security LoginToken
func (h *menuColumnHandler) BatchDelete(ctx *gin.Context) {
	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := h.svc.BatchDelete(ctx, ids)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("批量删除菜单数据列异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
}

// Update
// @Summary 更新菜单数据列
// @Description 更新菜单数据列信息
// @Tags 系统管理/菜单数据列管理
// @Accept application/json
// @Produce application/json
// @Param UpdateMenuColumnRequest body UpdateMenuColumnRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/menuColumn/update [put]
// @Security LoginToken
func (h *menuColumnHandler) Update(ctx *gin.Context) {
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

	var req UpdateMenuColumnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainSystem.MenuColumn{
		MenuColumn: modelSystem.MenuColumn{
			CoreModels: models.CoreModels{
				Id:         req.Id,
				Sort:       req.Sort,
				Version:    req.Version,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Title: req.Title,
			Field: req.Field,
			Width: req.Width,
		},
	}

	if err := h.svc.Update(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrMenuColumnVersionInconsistency):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
			return
		default:
			zap.L().Error("更新菜单数据列失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
}

// GetById
// @Summary 获取菜单数据列
// @Description 获取指定id菜单数据列信息
// @Tags 系统管理/菜单数据列管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} domainSystem.MenuColumn
// @Failure 400 {object} response.Response
// @Router /v1/system/menuColumn/getById/{id} [get]
// @Security LoginToken
func (h *menuColumnHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户ID不能为空", nil)
		return
	}

	detail, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceSystem.ErrMenuColumnNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "菜单数据列不存在", nil)
			return
		}
		zap.L().Error("获取菜单数据列失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
}

// GetListPage
// @Summary 获取菜单数据列分页列表
// @Description 获取菜单数据列分页列表
// @Tags 系统管理/菜单数据列管理
// @Accept application/json
// @Produce application/json
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "每页数量" default(10)
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param title query string false "标题"
// @Param field query string false "字段名"
// @Param menu_id query string false "菜单ID"
// @Success 200 {object} MenuColumnListPageResponse
// @Failure 400 {object} response.Response
// @Router /v1/system/menuColumn/listPage [get]
// @Security LoginToken
func (h *menuColumnHandler) GetListPage(ctx *gin.Context) {
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

	title := ctx.DefaultQuery("title", "")
	field := ctx.DefaultQuery("field", "")
	menuId := ctx.DefaultQuery("menu_id", "")

	filter := domainSystem.MenuColumnFilter{
		Pagination: filters.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status: status,
		Title:  title,
		Field:  field,
		MenuId: menuId,
	}

	list, total, err := h.svc.GetListPage(ctx, filter)
	if err != nil {
		zap.L().Error("获取分页列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", MenuColumnListPageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetListAll
// @Summary 获取所有菜单数据列
// @Description 获取所有菜单数据列列表
// @Tags 系统管理/菜单数据列管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param title query string false "标题"
// @Param field query string false "字段名"
// @Param menu_id query string false "菜单ID"
// @Success 200 {array} []domainSystem.MenuColumn
// @Failure 400 {object} response.Response
// @Router /v1/system/menuColumn/listAll [get]
// @Security LoginToken
func (h *menuColumnHandler) GetListAll(ctx *gin.Context) {
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

	title := ctx.DefaultQuery("title", "")
	field := ctx.DefaultQuery("field", "")
	menuId := ctx.DefaultQuery("menu_id", "")

	filter := domainSystem.MenuColumnFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status: status,
		Title:  title,
		Field:  field,
		MenuId: menuId,
	}

	list, err := h.svc.GetListAll(ctx, filter)
	if err != nil {
		zap.L().Error("获取列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}
