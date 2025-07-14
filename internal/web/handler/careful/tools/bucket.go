/**
 * Description：
 * FileName：bucket.go
 * Author：CJiaの用心
 * Create：2025/7/15 01:44:39
 * Remark：
 */

package tools

import (
	"errors"
	"fmt"
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainTools "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/tools"
	modelTools "github.com/carefuly/carefuly-admin-go-gin/internal/model/careful/tools"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	serviceTools "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/tools"
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

// CreateBucketRequest 创建
type CreateBucketRequest struct {
	Name   string `json:"name" binding:"required,max=50"`            // 存储桶名称
	Code   string `json:"code" binding:"required,max=50"`            // 存储桶编码
	Size   int    `json:"size" binding:"omitempty" default:"1"`      // 存储桶大小(GB)
	Sort   int    `json:"sort" binding:"omitempty" default:"1"`      // 排序
	Status bool   `json:"status" binding:"omitempty" default:"true"` // 状态【true-启用 false-停用】
	Remark string `json:"remark" binding:"omitempty,max=255"`        // 备注
}

// UpdateBucketRequest 更新
type UpdateBucketRequest struct {
	Id      string `json:"id" binding:"required"`                     // 主键ID
	Name    string `json:"name" binding:"required,max=50"`            // 存储桶名称
	Code    string `json:"code" binding:"required,max=50"`            // 存储桶大小(GB)
	Sort    int    `json:"sort" binding:"omitempty" default:"1"`      // 排序
	Status  bool   `json:"status" binding:"omitempty" default:"true"` // 状态【true-启用 false-停用】
	Version int    `json:"version" binding:"omitempty"`               // 版本
	Remark  string `json:"remark" binding:"omitempty,max=255"`        // 备注
}

// BucketListPageResponse 列表分页响应
type BucketListPageResponse struct {
	List     []domainTools.Bucket `json:"list"`     // 列表
	Total    int64                `json:"total"`    // 总数
	Page     int                  `json:"page"`     // 页码
	PageSize int                  `json:"pageSize"` // 每页数量
}

type BucketHandler interface {
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

type bucketHandler struct {
	rely    config.RelyConfig
	svc     serviceTools.BucketService
	userSvc serviceSystem.UserService
}

func NewBucketHandler(rely config.RelyConfig, svc serviceTools.BucketService, userSvc serviceSystem.UserService) BucketHandler {
	return &bucketHandler{
		rely:    rely,
		svc:     svc,
		userSvc: userSvc,
	}
}

// RegisterRoutes 注册路由
func (h *bucketHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/bucket")
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
// @Summary 创建存储桶
// @Description 创建存储桶
// @Tags 系统工具/存储桶管理
// @Accept application/json
// @Produce application/json
// @Param CreateBucketRequest body CreateBucketRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/bucket/create [post]
// @Security LoginToken
func (h *bucketHandler) Create(ctx *gin.Context) {
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

	var req CreateBucketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainTools.Bucket{
		Bucket: modelTools.Bucket{
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
			Size:   req.Size,
		},
	}

	if err := h.svc.Create(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceTools.ErrBucketNameDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "存储桶名称已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrBucketCodeDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "存储桶编码已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrBucketDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "存储桶已存在", nil)
			return
		default:
			zap.S().Error("创建存储桶失败", err.Error())
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
}

// Delete
// @Summary 删除存储桶
// @Description 删除指定id存储桶
// @Tags 系统工具/存储桶管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/bucket/delete/{id} [delete]
// @Security LoginToken
func (h *bucketHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		ctx.Set("internal", err.Error())
		zap.S().Error("删除存储桶失败", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
}

// BatchDelete
// @Summary 批量删除存储桶
// @Description 批量删除存储桶
// @Tags 系统工具/存储桶管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "id数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/bucket/delete/batchDelete [post]
// @Security LoginToken
func (h *bucketHandler) BatchDelete(ctx *gin.Context) {
	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := h.svc.BatchDelete(ctx, ids)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.S().Error("批量存储桶异常", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
}

// Update
// @Summary 更新存储桶
// @Description 更新存储桶信息
// @Tags 系统工具/存储桶管理
// @Accept application/json
// @Produce application/json
// @Param UpdateBucketRequest body UpdateBucketRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/bucket/update [put]
// @Security LoginToken
func (h *bucketHandler) Update(ctx *gin.Context) {
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

	var req UpdateBucketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainTools.Bucket{
		Bucket: modelTools.Bucket{
			CoreModels: models.CoreModels{
				Id:         req.Id,
				Sort:       req.Sort,
				Version:    req.Version,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status: req.Status,
			Name:   req.Name,
			Code:   req.Code,
		},
	}

	if err := h.svc.Update(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceTools.ErrBucketNameDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "存储桶名称已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrBucketCodeDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "存储桶编码已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrBucketDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "存储桶已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrBucketVersionInconsistency):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
			return
		default:
			zap.S().Error("更新存储桶失败", err.Error())
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
}

// GetById
// @Summary 获取存储桶
// @Description 获取指定id存储桶信息
// @Tags 系统工具/存储桶管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} domainTools.Bucket
// @Failure 400 {object} response.Response
// @Router /v1/tools/bucket/getById/{id} [get]
// @Security LoginToken
func (h *bucketHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	detail, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceTools.ErrBucketNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "存储桶不存在", nil)
			return
		}
		ctx.Set("internal", err.Error())
		zap.S().Error("获取存储桶失败", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
}

// GetListPage
// @Summary 获取存储桶分页列表
// @Description 获取存储桶分页列表
// @Tags 系统工具/存储桶管理
// @Accept application/json
// @Produce application/json
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "每页数量" default(10)
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "存储桶名称"
// @Param code query string false "存储桶编码"
// @Success 200 {object} BucketListPageResponse
// @Failure 400 {object} response.Response
// @Router /v1/tools/bucket/listPage [get]
// @Security LoginToken
func (h *bucketHandler) GetListPage(ctx *gin.Context) {
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

	filter := domainTools.BucketFilter{
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
		zap.S().Error("获取分页列表异常", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", BucketListPageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetListAll
// @Summary 获取所有存储桶
// @Description 获取所有存储桶列表
// @Tags 系统工具/存储桶管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "存储桶名称"
// @Param code query string false "存储桶编码"
// @Success 200 {array} []domainTools.Bucket
// @Failure 400 {object} response.Response
// @Router /v1/tools/bucket/listAll [get]
// @Security LoginToken
func (h *bucketHandler) GetListAll(ctx *gin.Context) {
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

	filter := domainTools.BucketFilter{
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
		zap.S().Error("获取列表异常", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}

// Export
// @Summary 导出存储桶数据
// @Description 导出存储桶数据到Excel文件
// @Tags 系统工具/存储桶管理
// @Accept application/json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "存储桶名称"
// @Param code query string false "存储桶编码"
// @Success 200 {file} file "Excel文件"
// @Failure 500 {object} response.Response
// @Router /v1/tools/bucket/export [get]
// @Security LoginToken
func (h *bucketHandler) Export(ctx *gin.Context) {
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

	filter := domainTools.BucketFilter{
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
		zap.S().Error("获取列表异常", err.Error())
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 准备导出配置
	filename := fmt.Sprintf("存储桶数据导出_%s.xlsx", time.Now().Format("20060102150405"))
	cfg := excelutil.ExcelExportConfig{
		SheetName:  "存储桶",
		FileName:   filename,
		StreamMode: true,
		Columns: []excelutil.ExcelColumn{
			{Title: "存储桶名称", Field: "Name", Width: 22},
			{Title: "存储桶编码", Field: "Code", Width: 17},
			{Title: "存储桶大小(GB)", Field: "Size", Width: 17},
			{
				Title: "状态",
				Field: "Status",
				Width: 10,
				Formatter: func(value interface{}) string {
					if status, ok := value.(bool); ok {
						if status {
							return "启用"
						}
						return "停用"
					}
					return fmt.Sprintf("%v", value)
				},
			},
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
		zap.S().Error("导出数据存储桶失败", err.Error())
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
