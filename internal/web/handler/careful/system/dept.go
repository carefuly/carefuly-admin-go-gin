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

// CreateDeptRequest 创建
type CreateDeptRequest struct {
	Name     string `json:"name" binding:"required,max=100" example:"测试部门"`           // 部门名称
	Code     string `json:"code" binding:"required,max=100" example:"CARE_TEST"`      // 部门编码
	Owner    string `json:"owner" binding:"omitempty" example:"admin"`                // 负责人
	Phone    string `json:"phone" binding:"omitempty" example:"18566666666"`          // 联系电话
	Email    string `json:"email" binding:"omitempty,email" example:"admin@test.com"` // 邮箱
	Status   bool   `json:"status" binding:"omitempty" example:"true"`                // 状态
	ParentID string `json:"parent_id" binding:"omitempty" example:"1"`                // 上级部门
	Remark   string `json:"remark" binding:"omitempty,max=255" example:"测试部门"`        // 备注
}

type DeptHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	Create(ctx *gin.Context)
	GetDeptTree(ctx *gin.Context)
}

type deptHandler struct {
	rely config.RelyConfig
	svc  serviceSystem.DeptService
}

func NewDeptHandler(rely config.RelyConfig, svc serviceSystem.DeptService) DeptHandler {
	return &deptHandler{
		rely: rely,
		svc:  svc,
	}
}

// RegisterRoutes 注册路由
func (h *deptHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/dept")
	base.POST("/create", h.Create)
	base.GET("/listTree", h.GetDeptTree)
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

	var req CreateDeptRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainSystem.Dept{
		Dept: modelSystem.Dept{
			CoreModels: models.CoreModels{
				Creator:  uid,
				Modifier: uid,
				Remark:   req.Remark,
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
		if errors.Is(err, serviceSystem.ErrDeptNameDuplicate) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "部门名称已存在", nil)
			return
		} else if errors.Is(err, serviceSystem.ErrDeptCodeDuplicate) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "部门编码已存在", nil)
			return
		} else if errors.Is(err, serviceSystem.ErrDeptDuplicate) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "部门信息已存在", nil)
			return
		} else {
			zap.L().Error("创建部门失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
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
func (h *deptHandler) GetDeptTree(ctx *gin.Context) {
	creator := ctx.DefaultQuery("creator", "")
	modifier := ctx.DefaultQuery("modifier", "")
	belongDept := ctx.DefaultQuery("belongDept", "")
	status, _ := strconv.ParseBool(ctx.DefaultQuery("status", "true"))
	name := ctx.DefaultQuery("name", "")
	code := ctx.DefaultQuery("code", "")

	filter := domainSystem.DeptFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: belongDept,
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
