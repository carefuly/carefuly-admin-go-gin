/**
 * Description：
 * FileName：menu.go
 * Author：CJiaの用心
 * Create：2025/5/13 16:41:14
 * Remark：
 */

package system

import (
	"errors"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/menu"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/enumconv"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// CreateMenuRequest 创建
type CreateMenuRequest struct {
	Type        menu.TypeConst `json:"type" binding:"omitempty" default:"2"`                 // 菜单类型
	Icon        string         `json:"icon" binding:"omitempty,max=64" default:"HomeFilled"` // 菜单图标
	Title       string         `json:"title" binding:"required,max=64"`                      // 菜单标题
	Name        string         `json:"name" binding:"required,max=64"`                       // 组件名称
	Component   string         `json:"component" binding:"omitempty,max=128"`                // 组件名称
	Path        string         `json:"path" binding:"required,max=128"`                      // 路由地址
	Redirect    string         `json:"redirect" binding:"omitempty,max=128"`                 // 重定向地址
	IsHide      bool           `json:"isHide" binding:"omitempty" default:"false"`           // 是否隐藏
	IsLink      string         `json:"isLink" binding:"omitempty,max=255"`                   // 是否外链【不填写默认没有外链】
	IsKeepAlive bool           `json:"isKeepAlive" binding:"omitempty" default:"false"`      // 是否页面缓存
	IsFull      bool           `json:"isFull" binding:"omitempty" default:"false"`           // 是否缓存全屏
	IsAffix     bool           `json:"isAffix" binding:"omitempty" default:"false"`          // 是否缓存固定路由
	ParentID    string         `json:"parent_id" binding:"omitempty,max=100"`                // 上级菜单
	Sort        int            `json:"sort" binding:"omitempty" default:"1"`                 // 排序
	Status      bool           `json:"status" binding:"omitempty" default:"true"`            // 状态【true-启用 false-停用】
	Remark      string         `json:"remark" binding:"omitempty,max=255"`                   // 备注
}

// UpdateMenuRequest 更新
type UpdateMenuRequest struct {
	Id          string         `json:"id" binding:"required"`                                // 主键ID
	Type        menu.TypeConst `json:"type" binding:"omitempty" default:"2"`                 // 菜单类型
	Icon        string         `json:"icon" binding:"omitempty,max=64" default:"HomeFilled"` // 菜单图标
	Title       string         `json:"title" binding:"required,max=64"`                      // 菜单标题
	Name        string         `json:"name" binding:"required,max=64"`                       // 组件名称
	Component   string         `json:"component" binding:"omitempty,max=128"`                // 组件名称
	Path        string         `json:"path" binding:"required,max=128"`                      // 路由地址
	Redirect    string         `json:"redirect" binding:"omitempty,max=128"`                 // 重定向地址
	IsHide      bool           `json:"isHide" binding:"omitempty" default:"false"`           // 是否隐藏
	IsLink      string         `json:"isLink" binding:"omitempty,max=255"`                   // 是否外链【不填写默认没有外链】
	IsKeepAlive bool           `json:"isKeepAlive" binding:"omitempty" default:"false"`      // 是否页面缓存
	IsFull      bool           `json:"isFull" binding:"omitempty" default:"false"`           // 是否缓存全屏
	IsAffix     bool           `json:"isAffix" binding:"omitempty" default:"false"`          // 是否缓存固定路由
	ParentID    string         `json:"parent_id" binding:"omitempty,max=100"`                // 上级菜单
	Sort        int            `json:"sort" binding:"omitempty" default:"1"`                 // 排序
	Status      bool           `json:"status" binding:"omitempty" default:"true"`            // 状态【true-启用 false-停用】
	Version     int            `json:"version" binding:"omitempty"`                          // 版本
	Remark      string         `json:"remark" binding:"omitempty,max=255"`                   // 备注
}

// MenuListPageResponse 列表分页响应
type MenuListPageResponse struct {
	List     []domainSystem.Menu `json:"list"`     // 列表
	Total    int64               `json:"total"`    // 总数
	Page     int                 `json:"page"`     // 页码
	PageSize int                 `json:"pageSize"` // 每页数量
}

type MenuHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	BatchDelete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetListRouter(ctx *gin.Context)
}

type menuHandler struct {
	rely    config.RelyConfig
	svc     serviceSystem.MenuService
	userSvc serviceSystem.UserService
}

func NewMenuHandler(rely config.RelyConfig, svc serviceSystem.MenuService, userSvc serviceSystem.UserService) MenuHandler {
	return &menuHandler{
		rely:    rely,
		svc:     svc,
		userSvc: userSvc,
	}
}

// RegisterRoutes 注册路由
func (h *menuHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/menu")
	base.POST("/create", h.Create)
	base.DELETE("/delete/:id", h.Delete)
	base.POST("/delete/batchDelete", h.BatchDelete)
	base.PUT("/update", h.Update)
	base.GET("/getById/:id", h.GetById)
	base.GET("/listRouter", h.GetListRouter)
}

// Create
// @Summary 创建菜单
// @Description 创建菜单
// @Tags 系统管理/菜单管理
// @Accept application/json
// @Produce application/json
// @Param CreateMenuRequest body CreateMenuRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/menu/create [post]
// @Security LoginToken
func (h *menuHandler) Create(ctx *gin.Context) {
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

	var req CreateMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 校验参数
	typeValidValues := []string{"目录", "菜单"}
	converter := enumconv.NewEnumConverter(menu.TypeMapping, menu.TypeImportMapping, typeValidValues, "菜单类型")
	_, err = converter.FromEnum(req.Type)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 转换为领域模型
	domain := domainSystem.Menu{
		Menu: modelSystem.Menu{
			CoreModels: models.CoreModels{
				Sort:       req.Sort,
				Creator:    uid,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status:      req.Status,
			Type:        req.Type,
			Icon:        req.Icon,
			Title:       req.Title,
			Name:        req.Name,
			Component:   req.Component,
			Path:        req.Path,
			Redirect:    req.Redirect,
			IsHide:      req.IsHide,
			IsLink:      req.IsLink,
			IsKeepAlive: req.IsKeepAlive,
			IsFull:      req.IsFull,
			IsAffix:     req.IsAffix,
			ParentID:    req.ParentID,
		},
	}

	if err := h.svc.Create(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrMenuNameDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "菜单已存在", nil)
			return
		default:
			zap.L().Error("创建菜单失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
}

// Delete
// @Summary 删除菜单
// @Description 删除指定id菜单
// @Tags 系统管理/菜单管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/menu/delete/{id} [delete]
// @Security LoginToken
func (h *menuHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, serviceSystem.ErrMenuNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "菜单不存在", nil)
			return
		}
		ctx.Set("internal", err.Error())
		zap.L().Error("删除菜单失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
}

// BatchDelete
// @Summary 批量删除菜单
// @Description 批量删除菜单
// @Tags 系统管理/菜单管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "id数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/menu/delete/batchDelete [post]
// @Security LoginToken
func (h *menuHandler) BatchDelete(ctx *gin.Context) {
	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := h.svc.BatchDelete(ctx, ids)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("批量删除菜单异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
}

// Update
// @Summary 更新菜单
// @Description 更新菜单信息
// @Tags 系统管理/菜单管理
// @Accept application/json
// @Produce application/json
// @Param UpdateMenuRequest body UpdateMenuRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/menu/update [put]
// @Security LoginToken
func (h *menuHandler) Update(ctx *gin.Context) {
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

	var req UpdateMenuRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainSystem.Menu{
		Menu: modelSystem.Menu{
			CoreModels: models.CoreModels{
				Id:         req.Id,
				Sort:       req.Sort,
				Version:    req.Version,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status:      req.Status,
			Type:        req.Type,
			Icon:        req.Icon,
			Title:       req.Title,
			Name:        req.Name,
			Component:   req.Component,
			Path:        req.Path,
			Redirect:    req.Redirect,
			IsHide:      req.IsHide,
			IsLink:      req.IsLink,
			IsKeepAlive: req.IsKeepAlive,
			IsFull:      req.IsFull,
			IsAffix:     req.IsAffix,
			ParentID:    req.ParentID,
		},
	}

	if err := h.svc.Update(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrMenuNameDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "菜单已存在", nil)
			return
		case errors.Is(err, serviceSystem.ErrMenuVersionInconsistency):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
			return
		default:
			zap.L().Error("更新菜单失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
}

// GetById
// @Summary 获取菜单
// @Description 获取指定id菜单信息
// @Tags 系统管理/菜单管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} domainSystem.Menu
// @Failure 400 {object} response.Response
// @Router /v1/system/menu/getById/{id} [get]
// @Security LoginToken
func (h *menuHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	detail, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceSystem.ErrMenuNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "菜单不存在", nil)
			return
		}
		zap.L().Error("获取菜单失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
}

// GetListRouter
// @Summary 获取所有菜单
// @Description 获取所有菜单列表
// @Tags 系统管理/菜单管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param title query string false "菜单标题"
// @Success 200 {array} []domainSystem.Menu
// @Failure 400 {object} response.Response
// @Router /v1/system/menu/listRouter [get]
// @Security LoginToken
func (h *menuHandler) GetListRouter(ctx *gin.Context) {
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

	filter := domainSystem.MenuFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status: status,
		Title:  title,
	}

	list, err := h.svc.GetListAll(ctx, filter)
	if err != nil {
		zap.L().Error("获取列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}
