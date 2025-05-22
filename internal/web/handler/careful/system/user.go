/**
 * Description：
 * FileName：user.go
 * Author：CJiaの用心
 * Create：2025/5/18 21:22:50
 * Remark：
 */

package system

import (
	config "github.com/carefuly/carefuly-admin-go-gin/config/file"
	domainSystem "github.com/carefuly/carefuly-admin-go-gin/internal/domain/careful/system"
	serviceSystem "github.com/carefuly/carefuly-admin-go-gin/internal/service/careful/system"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// UserListPageResponse 用户列表分页响应
type UserListPageResponse struct {
	List     []domainSystem.User `json:"list"`     // 列表
	Total    int64               `json:"total"`    // 总数
	Page     int                 `json:"page"`     // 页码
	PageSize int                 `json:"pageSize"` // 每页数量
}

type UserHandler interface {
	RegisterRoutes(router *gin.RouterGroup)

	GetListPage(ctx *gin.Context)
}

type userHandler struct {
	rely config.RelyConfig
	svc  serviceSystem.UserService
}

func NewUserHandler(rely config.RelyConfig, svc serviceSystem.UserService) UserHandler {
	return &userHandler{
		rely: rely,
		svc:  svc,
	}
}

// RegisterRoutes 注册路由
func (h *userHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/user")
	base.GET("/listPage", h.GetListPage)
}

// GetListPage
// @Summary 获取用户分页列表
// @Description 获取用户分页列表
// @Tags 系统管理/用户管理
// @Accept application/json
// @Produce application/json
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "每页数量" default(10)
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param belongDept query string false "数据归属部门"
// @Param status query bool false "状态" default(true)
// @Param username query string false "用户名"
// @Success 200 {object} UserListPageResponse
// @Failure 400 {object} response.Response
// @Router /v1/system/user/listPage [get]
// @Security LoginToken
func (h *userHandler) GetListPage(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	creator := ctx.DefaultQuery("creator", "")
	modifier := ctx.DefaultQuery("modifier", "")
	belongDept := ctx.DefaultQuery("belongDept", "")
	status, _ := strconv.ParseBool(ctx.DefaultQuery("status", "true"))
	username := ctx.DefaultQuery("username", "")

	list, total, err := h.svc.GetListPage(ctx, domainSystem.UserFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: belongDept,
		},
		Pagination: filters.Pagination{
			Page:     page,
			PageSize: pageSize,
		},
		Status:   status,
		Username: username,
	})

	if err != nil {
		zap.L().Error("分页查询列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", UserListPageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}
