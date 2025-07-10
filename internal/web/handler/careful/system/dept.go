/**
 * Description：
 * FileName：dept.go
 * Author：CJiaの用心
 * Create：2025/5/15 17:08:53
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

// CreateDeptRequest 创建
type CreateDeptRequest struct {
	Name     string `json:"name" binding:"required,max=100"`           // 部门名称
	Code     string `json:"code" binding:"required,max=100"`           // 部门编码
	Owner    string `json:"owner" binding:"omitempty"`                 // 负责人
	Phone    string `json:"phone" binding:"omitempty"`                 // 联系电话
	Email    string `json:"email" binding:"omitempty,email"`           // 邮箱
	ParentID string `json:"parent_id" binding:"omitempty"`             // 上级部门
	Sort     int    `json:"sort" binding:"omitempty" default:"1"`      // 排序
	Status   bool   `json:"status" binding:"omitempty" default:"true"` // 状态【true-启用 false-停用】
	Remark   string `json:"remark" binding:"omitempty,max=255"`        // 备注
}

// UpdateDeptRequest 更新
type UpdateDeptRequest struct {
	Id       string `json:"id" binding:"required"`                     // 主键ID
	Name     string `json:"name" binding:"required,max=100"`           // 部门名称
	Code     string `json:"code" binding:"required,max=100"`           // 部门编码
	Owner    string `json:"owner" binding:"omitempty"`                 // 负责人
	Phone    string `json:"phone" binding:"omitempty"`                 // 联系电话
	Email    string `json:"email" binding:"omitempty,email"`           // 邮箱
	ParentID string `json:"parent_id" binding:"omitempty"`             // 上级部门
	Sort     int    `json:"sort" binding:"omitempty" default:"1"`      // 排序
	Status   bool   `json:"status" binding:"omitempty" default:"true"` // 状态【true-启用 false-停用】
	Version  int    `json:"version" binding:"omitempty"`               // 版本
	Remark   string `json:"remark" binding:"omitempty,max=255"`        // 备注
}

// DeptListPageResponse 列表分页响应
type DeptListPageResponse struct {
	List     []domainSystem.Dept `json:"list"`     // 列表
	Total    int64               `json:"total"`    // 总数
	Page     int                 `json:"page"`     // 页码
	PageSize int                 `json:"pageSize"` // 每页数量
}

type DeptHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	Create(ctx *gin.Context)
	Delete(ctx *gin.Context)
	BatchDelete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetDeptTree(ctx *gin.Context)
	Export(ctx *gin.Context)
}

type deptHandler struct {
	rely    config.RelyConfig
	svc     serviceSystem.DeptService
	userSvc serviceSystem.UserService
}

func NewDeptHandler(rely config.RelyConfig, svc serviceSystem.DeptService, userSvc serviceSystem.UserService) DeptHandler {
	return &deptHandler{
		rely:    rely,
		svc:     svc,
		userSvc: userSvc,
	}
}

// RegisterRoutes 注册路由
func (h *deptHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/dept")
	base.POST("/create", h.Create)
	base.DELETE("/delete/:id", h.Delete)
	base.POST("/delete/batchDelete", h.BatchDelete)
	base.PUT("/update", h.Update)
	base.GET("/getById/:id", h.GetById)
	base.GET("/listTree", h.GetDeptTree)
	base.GET("/export", h.Export)
}

// Create
// @Summary 创建部门
// @Description 创建部门
// @Tags 系统管理/部门管理
// @Accept application/json
// @Produce application/json
// @Param CreateDeptRequest body CreateDeptRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/dept/create [post]
// @Security LoginToken
func (h *deptHandler) Create(ctx *gin.Context) {
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

	var req CreateDeptRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainSystem.Dept{
		Dept: modelSystem.Dept{
			CoreModels: models.CoreModels{
				Sort:       req.Sort,
				Creator:    uid,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status:   req.Status,
			Name:     req.Name,
			Code:     req.Code,
			Owner:    req.Owner,
			Phone:    req.Phone,
			Email:    req.Email,
			ParentID: req.ParentID,
		},
	}

	if err := h.svc.Create(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrDeptDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "部门信息已存在", nil)
			return
		default:
			zap.L().Error("创建部门失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
}

// Delete
// @Summary 删除部门
// @Description 删除指定id部门
// @Tags 系统管理/部门管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/dept/delete/{id} [delete]
// @Security LoginToken
func (h *deptHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrDeptNotFound):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "部门不存在", nil)
			return
		case errors.Is(err, serviceSystem.ErrDeptChildNodes):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "请先删除当前部门下的子部门", nil)
			return
		default:
			ctx.Set("internal", err.Error())
			zap.L().Error("删除部门失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
}

// BatchDelete
// @Summary 批量删除部门
// @Description 批量删除部门
// @Tags 系统管理/部门管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "id数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/dept/delete/batchDelete [post]
// @Security LoginToken
func (h *deptHandler) BatchDelete(ctx *gin.Context) {
	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := h.svc.BatchDelete(ctx, ids)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("批量删除部门异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
}

// Update
// @Summary 更新部门
// @Description 更新部门信息
// @Tags 系统管理/部门管理
// @Accept application/json
// @Produce application/json
// @Param UpdateDeptRequest body UpdateDeptRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/system/dept/update [put]
// @Security LoginToken
func (h *deptHandler) Update(ctx *gin.Context) {
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

	var req UpdateDeptRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainSystem.Dept{
		Dept: modelSystem.Dept{
			CoreModels: models.CoreModels{
				Id:         req.Id,
				Sort:       req.Sort,
				Version:    req.Version,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status:   req.Status,
			Name:     req.Name,
			Code:     req.Code,
			Owner:    req.Owner,
			Phone:    req.Phone,
			Email:    req.Email,
			ParentID: req.ParentID,
		},
	}

	if err := h.svc.Update(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceSystem.ErrDeptDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "部门信息已存在", nil)
			return
		case errors.Is(err, serviceSystem.ErrDeptVersionInconsistency):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
			return
		default:
			zap.L().Error("更新部门失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
}

// GetById
// @Summary 获取部门
// @Description 获取指定id部门信息
// @Tags 系统管理/部门管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} domainSystem.Dept
// @Failure 400 {object} response.Response
// @Router /v1/system/dept/getById/{id} [get]
// @Security LoginToken
func (h *deptHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	detail, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceSystem.ErrDeptNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "部门不存在", nil)
			return
		}
		zap.L().Error("获取部门失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
}

// GetDeptTree 获取部门树形结构
// @Summary 获取部门树形结构
// @Description 获取部门树形结构
// @Tags 系统管理/部门管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param belongDept query string false "数据归属部门"
// @Param status query bool false "状态" default(true)
// @Param name query string false "部门名称"
// @Param code query string false "部门编码"
// @Success 200 {object} serviceSystem.DeptTree
// @Failure 400 {object} response.Response
// @Router /v1/system/dept/listTree [get]
// @Security LoginToken
func (h *deptHandler) GetDeptTree(ctx *gin.Context) {
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

	filter := domainSystem.DeptFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status: status,
		Name:   name,
		Code:   code,
	}

	tree, err := h.svc.GetListTree(ctx, filter)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("获取部门树失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", tree)
}

// Export
// @Summary 导出部门信息
// @Description 导出部门信息到Excel文件
// @Tags 系统管理/部门管理
// @Accept application/json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param belongDept query string false "数据归属部门"
// @Param status query bool false "状态" default(true)
// @Param name query string false "部门名称"
// @Param code query string false "部门编码"
// @Success 200 {file} file "Excel文件"
// @Failure 500 {object} response.Response
// @Router /v1/system/dept/export [get]
// @Security LoginToken
func (h *deptHandler) Export(ctx *gin.Context) {
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

	filter := domainSystem.DeptFilter{
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
	filename := fmt.Sprintf("部门信息导出_%s.xlsx", time.Now().Format("20060102150405"))
	cfg := excelutil.ExcelExportConfig{
		SheetName:  "部门",
		FileName:   filename,
		StreamMode: true,
		Columns: []excelutil.ExcelColumn{
			{Title: "部门名称", Field: "Name", Width: 22},
			{Title: "部门编码", Field: "Code", Width: 18},
			{Title: "负责人", Field: "Owner", Width: 18},
			{Title: "联系电话", Field: "Phone", Width: 18},
			{Title: "邮箱", Field: "Email", Width: 22},
			{Title: "上级部门", Field: "ParentName", Width: 22},
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
		zap.L().Error("导出部门失败", zap.Error(err))
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
