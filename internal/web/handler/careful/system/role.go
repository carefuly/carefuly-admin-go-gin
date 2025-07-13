/**
 * Description：
 * FileName：role.go
 * Author：CJiaの用心
 * Create：2025/6/12 14:23:45
 * Remark：
 */

package system

import (
	"errors"
	"fmt"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	modelSystem "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/system"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/system/role"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/excelutil"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// CreateRoleRequest 创建
type CreateRoleRequest struct {
	Name   string `json:"name" binding:"required,max=100"`           // 角色名称
	Code   string `json:"code" binding:"required,max=100"`           // 角色编码
	Sort   int    `json:"sort" binding:"omitempty" default:"1"`      // 排序
	Status bool   `json:"status" binding:"omitempty" default:"true"` // 状态【true-启用 false-停用】
	Remark string `json:"remark" binding:"omitempty,max=255"`        // 备注
}

// UpdateRoleRequest 更新
type UpdateRoleRequest struct {
	Id            string              `json:"id" binding:"required"`                     // 主键ID
	Name          string              `json:"name" binding:"required,max=100"`           // 角色名称
	Code          string              `json:"code" binding:"required,max=100"`           // 角色编码
	DataRange     role.DataRangeConst `json:"data_range" binding:"omitempty"`            // 数据权限范围
	DeptIDs       []string            `json:"dept_ids" binding:"omitempty"`              // 部门ID数组
	MenuIDs       []string            `json:"menu_ids" binding:"omitempty"`              // 菜单ID数组
	MenuButtonIDs []string            `json:"menu_button_ids" binding:"omitempty"`       // 按钮ID数组
	MenuColumnIDs []string            `json:"menu_column_ids" binding:"omitempty"`       // 列权限ID数组
	Sort          int                 `json:"sort" binding:"omitempty" default:"1"`      // 排序
	Status        bool                `json:"status" binding:"omitempty" default:"true"` // 状态【true-启用 false-停用】
	Version       int                 `json:"version" binding:"omitempty"`               // 版本
	Remark        string              `json:"remark" binding:"omitempty,max=255"`        // 备注
}

// RoleListPageResponse 列表分页响应
type RoleListPageResponse struct {
	List     []domainSystem.Role `json:"list"`     // 列表
	Total    int64               `json:"total"`    // 总数
	Page     int                 `json:"page"`     // 页码
	PageSize int                 `json:"pageSize"` // 每页数量
}

type RoleHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	BatchDelete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetListPage(ctx *gin.Context)
	GetListAll(ctx *gin.Context)
	Export(ctx *gin.Context)
}

type roleHandler struct {
	rely    config.RelyConfig
	svc     serviceSystem.RoleService
	userSvc serviceSystem.UserService
}

func NewRoleHandler(rely config.RelyConfig, svc serviceSystem.RoleService, userSvc serviceSystem.UserService) RoleHandler {
	return &roleHandler{
		rely:    rely,
		svc:     svc,
		userSvc: userSvc,
	}
}

// RegisterRoutes 注册路由
func (h *roleHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/role")
	base.POST("/create", h.Create)
	base.DELETE("/delete/:id", h.Delete)
	base.POST("/delete/batchDelete", h.BatchDelete)
	base.PUT("/update", h.Update)
	base.GET("/getById/:id", h.GetById)
	base.GET("/listPage", h.GetListPage)
	base.GET("/listAll", h.GetListAll)
	base.GET("/export", h.Export)
}

// Create
// @Summary 创建角色
// @Description 创建角色
// @Tags 系统管理/角色管理
// @Accept application/json
// @Produce application/json
// @Param CreateRoleRequest body CreateRoleRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/role/create [post]
// @Security LoginToken
func (h *roleHandler) Create(ctx *gin.Context) {
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

	var req CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainSystem.Role{
		Role: modelSystem.Role{
			CoreModels: models.CoreModels{
				Sort:       req.Sort,
				Creator:    uid,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status: req.Status,
			Name:   req.Name,
			Code:   req.Code,
		},
	}

	if err := h.svc.Create(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrRoleCodeDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "角色编码已存在", nil)
			return
		default:
			zap.L().Error("创建角色失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
}

// Delete
// @Summary 删除角色
// @Description 删除指定id角色
// @Tags 系统管理/角色管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/role/delete/{id} [delete]
// @Security LoginToken
func (h *roleHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrRoleNotFound):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "角色不存在", nil)
			return
		default:
			ctx.Set("internal", err.Error())
			zap.L().Error("删除角色失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
}

// BatchDelete
// @Summary 批量删除角色
// @Description 批量删除角色
// @Tags 系统管理/角色管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "id数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/role/delete/batchDelete [post]
// @Security LoginToken
func (h *roleHandler) BatchDelete(ctx *gin.Context) {
	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := h.svc.BatchDelete(ctx, ids)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("批量删除角色异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
}

// Update
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 系统管理/角色管理
// @Accept application/json
// @Produce application/json
// @Param UpdateRoleRequest body UpdateRoleRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/role/update [put]
// @Security LoginToken
func (h *roleHandler) Update(ctx *gin.Context) {
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

	var req UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainSystem.Role{
		Role: modelSystem.Role{
			CoreModels: models.CoreModels{
				Id:         req.Id,
				Sort:       req.Sort,
				Version:    req.Version,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status:        req.Status,
			Name:          req.Name,
			Code:          req.Code,
			DataRange:     req.DataRange,
			DeptIDs:       req.DeptIDs,
			MenuIDs:       req.MenuIDs,
			MenuButtonIDs: req.MenuButtonIDs,
			MenuColumnIDs: req.MenuColumnIDs,
		},
	}

	if err := h.svc.Update(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrRoleCodeDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "角色编码已存在", nil)
			return
		case errors.Is(err, serviceSystem.ErrRoleVersionInconsistency):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
			return
		default:
			zap.L().Error("更新角色失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
}

// GetById
// @Summary 获取角色
// @Description 获取指定id角色信息
// @Tags 系统管理/角色管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} domainSystem.Role
// @Failure 400 {object} response.Response
// @Router /v1/system/role/getById/{id} [get]
// @Security LoginToken
func (h *roleHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	detail, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceSystem.ErrRoleNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "角色不存在", nil)
			return
		}
		zap.L().Error("获取角色失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
}

// GetListPage
// @Summary 获取角色分页列表
// @Description 获取角色分页列表
// @Tags 系统管理/角色管理
// @Accept application/json
// @Produce application/json
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "每页数量" default(10)
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "角色名称"
// @Param code query string false "角色编码"
// @Success 200 {object} RoleListPageResponse
// @Failure 400 {object} response.Response
// @Router /v1/system/role/listPage [get]
// @Security LoginToken
func (h *roleHandler) GetListPage(ctx *gin.Context) {
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

	filter := domainSystem.RoleFilter{
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
		Name:   name,
		Code:   code,
	}

	list, total, err := h.svc.GetListPage(ctx, filter)
	if err != nil {
		zap.L().Error("获取分页列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", RoleListPageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetListAll
// @Summary 获取所有角色
// @Description 获取所有角色列表
// @Tags 系统管理/角色管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "角色名称"
// @Param code query string false "角色编码"
// @Success 200 {array} []domainSystem.Role
// @Failure 400 {object} response.Response
// @Router /v1/system/role/listAll [get]
// @Security LoginToken
func (h *roleHandler) GetListAll(ctx *gin.Context) {
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

	filter := domainSystem.RoleFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status: status,
		Name:   name,
		Code:   code,
	}

	list, err := h.svc.GetListAll(ctx, filter)
	if err != nil {
		zap.L().Error("获取列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}

// Export
// @Summary 导出角色信息
// @Description 导出角色信息到Excel文件
// @Tags 系统管理/角色管理
// @Accept application/json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "角色名称"
// @Param code query string false "角色编码"
// @Success 200 {file} file "Excel文件"
// @Failure 500 {object} response.Response
// @Router /v1/system/role/export [get]
// @Security LoginToken
func (h *roleHandler) Export(ctx *gin.Context) {
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

	filter := domainSystem.RoleFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status: status,
		Name:   name,
		Code:   code,
	}

	list, err := h.svc.GetListAll(ctx, filter)
	if err != nil {
		zap.L().Error("获取列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 准备导出配置
	filename := fmt.Sprintf("角色信息导出_%s.xlsx", time.Now().Format("20060102150405"))
	cfg := excelutil.ExcelExportConfig{
		SheetName:  "角色",
		FileName:   filename,
		StreamMode: true,
		Columns: []excelutil.ExcelColumn{
			{Title: "角色名称", Field: "Name", Width: 22},
			{Title: "角色编码", Field: "Code", Width: 17},
			{Title: "排序", Field: "Sort", Width: 8},
			{Title: "创建时间", Field: "CreateTime", Width: 22},
			{Title: "更新时间", Field: "UpdateTime", Width: 22},
			{Title: "备注", Field: "Remark", Width: 40},
		},
		Data: list,
	}

	// 创建并执行导出器
	exporter := excelutil.NewExcelExporter(&cfg)
	f, err := exporter.Export()
	if err != nil {
		zap.L().Error("导出角色失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 设置响应头
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename=export.xlsx")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Cache-Control", "no-store")

	// 流式写入响应
	if _, err := f.WriteTo(ctx.Writer); err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "生成Excel失败", nil)
	}
}
