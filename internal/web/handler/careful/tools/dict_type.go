/**
 * Description：
 * FileName：dict_type.go
 * Author：CJiaの用心
 * Create：2025/6/6 11:57:30
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
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/tools/dict"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/constants/careful/tools/dictType"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/filters"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/ginx/response"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/models"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/enumconv"
	"github.com/carefuly/carefuly-admin-go-gin/pkg/utils/excelutil"
	validate "github.com/carefuly/carefuly-admin-go-gin/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// CreateDictTypeRequest 创建
type CreateDictTypeRequest struct {
	Name      string                `json:"name" binding:"required,max=50"`             // 字典信息名称
	StrValue  string                `json:"strValue" binding:"max=50"`                  // 字符串-字典信息值
	IntValue  int64                 `json:"intValue"`                                   // 整型-字典信息值
	BoolValue bool                  `json:"boolValue"`                                  // 布尔-字典信息值
	DictTag   dictType.DictTagConst `json:"dictTag" binding:"max=10" default:"primary"` // 标签类型
	DictColor string                `json:"dictColor" binding:"max=50"`                 // 标签颜色
	DictId    string                `json:"dict_id" binding:"required,max=100"`         // 字典ID
	Sort      int                   `json:"sort" binding:"omitempty" default:"1"`       // 排序
	Status    bool                  `json:"status" binding:"omitempty" default:"true"`  // 状态【true-启用 false-停用】
	Remark    string                `json:"remark" binding:"omitempty,max=255"`         // 备注
}

// ImportDictTypeRequest 导入
type ImportDictTypeRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// UpdateDictTypeRequest 更新
type UpdateDictTypeRequest struct {
	Id        string                `json:"id" binding:"required"`                      // 主键ID
	Name      string                `json:"name" binding:"required,max=50"`             // 字典信息名称
	DictTag   dictType.DictTagConst `json:"dictTag" binding:"max=10" default:"primary"` // 标签类型
	DictColor string                `json:"dictColor" binding:"max=50"`                 // 标签颜色
	DictId    string                `json:"dict_id" binding:"required,max=100"`         // 字典ID
	Sort      int                   `json:"sort" binding:"omitempty" default:"1"`       // 排序
	Status    bool                  `json:"status" binding:"omitempty" default:"true"`  // 状态【true-启用 false-停用】
	Version   int                   `json:"version" binding:"omitempty"`                // 版本
	Remark    string                `json:"remark" binding:"omitempty,max=255"`         // 备注
}

type ListByDictNamesRequest struct {
	DictNames []string `json:"dictNames"` // 数组参数格式: ?dictNames=性别&dictNames=计量单位
}

// DictTypeListPageResponse 列表分页响应
type DictTypeListPageResponse struct {
	List     []domainTools.DictType `json:"list"`     // 列表
	Total    int64                  `json:"total"`    // 总数
	Page     int                    `json:"page"`     // 页码
	PageSize int                    `json:"pageSize"` // 每页数量
}

type DictTypeHandler interface {
	RegisterRoutes(router *gin.RouterGroup)
	Create(ctx *gin.Context)
	Import(ctx *gin.Context)
	Delete(ctx *gin.Context)
	BatchDelete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetById(ctx *gin.Context)
	GetListByDictNames(ctx *gin.Context)
	GetListPage(ctx *gin.Context)
	GetListAll(ctx *gin.Context)
	Export(ctx *gin.Context)
}

type dictTypeHandler struct {
	rely    config.RelyConfig
	svc     serviceTools.DictTypeService
	userSvc serviceSystem.UserService
}

func NewDictTypeHandler(rely config.RelyConfig, svc serviceTools.DictTypeService, userSvc serviceSystem.UserService) DictTypeHandler {
	return &dictTypeHandler{
		rely:    rely,
		svc:     svc,
		userSvc: userSvc,
	}
}

func (h *dictTypeHandler) RegisterRoutes(router *gin.RouterGroup) {
	base := router.Group("/dictType")
	base.POST("/create", h.Create)
	base.POST("/import", h.Import)
	base.DELETE("/delete/:id", h.Delete)
	base.POST("/delete/batchDelete", h.BatchDelete)
	base.PUT("/update", h.Update)
	base.GET("/getById/:id", h.GetById)
	base.POST("/listByDictNames", h.GetListByDictNames)
	base.GET("/listPage", h.GetListPage)
	base.GET("/listAll", h.GetListAll)
	base.GET("/export", h.Export)
}

// Create
// @Summary 创建字典信息
// @Description 创建字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param CreateDictTypeRequest body CreateDictTypeRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/create [post]
// @Security LoginToken
func (h *dictTypeHandler) Create(ctx *gin.Context) {
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

	var req CreateDictTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 校验参数
	dictTagValues := []string{"primary", "success", "warning", "danger", "info"}
	converter := enumconv.NewEnumConverter(dictType.DictTagMapping, dictType.DictTagImportMapping, dictTagValues, "标签类型")
	_, err = converter.FromEnum(req.DictTag)
	if err != nil {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// 转换为领域模型
	domain := domainTools.DictType{
		DictType: modelTools.DictType{
			CoreModels: models.CoreModels{
				Sort:       req.Sort,
				Creator:    uid,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status:    req.Status,
			Name:      req.Name,
			DictTag:   req.DictTag,
			DictColor: req.DictColor,
			DictId:    req.DictId,
		},
		StrValue:  req.StrValue,
		IntValue:  req.IntValue,
		BoolValue: req.BoolValue,
	}

	if err := h.svc.Create(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceTools.ErrDictTypeDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "同一字典下存在相同的字典项/值", nil)
			return
		case errors.Is(err, serviceTools.ErrDictTypeInvalidDictValueType):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "不支持的字典类型", nil)
		default:
			ctx.Set("internal", err.Error())
			zap.L().Error("创建字典信息失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
}

// Import
// @Summary 导入字典信息
// @Description 导入字典信息
// @Tags 系统工具/字典信息管理
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "文件(支持xlsx/csv格式)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/import [post]
// @Security LoginToken
func (h *dictTypeHandler) Import(ctx *gin.Context) {

}

// Delete
// @Summary 删除字典信息
// @Description 删除指定id信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/delete/{id} [delete]
// @Security LoginToken
func (h *dictTypeHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
		return
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		if errors.Is(err, serviceTools.ErrDictTypeNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典信息不存在", nil)
			return
		}
		ctx.Set("internal", err.Error())
		zap.L().Error("删除字典信息异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
	}

	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
}

// BatchDelete
// @Summary 批量删除字典信息
// @Description 批量删除字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param ids body []string true "id数组"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/batchDelete [post]
// @Security LoginToken
func (h *dictTypeHandler) BatchDelete(ctx *gin.Context) {
	var ids []string
	if err := ctx.ShouldBindJSON(&ids); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	err := h.svc.BatchDelete(ctx, ids)
	if err != nil {
		ctx.Set("internal", err.Error())
		zap.L().Error("批量删除字典信息异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
}

// Update
// @Summary 更新字典信息
// @Description 更新字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param UpdateDictTypeRequest body UpdateDictTypeRequest true "请求"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/update [put]
// @Security LoginToken
func (h *dictTypeHandler) Update(ctx *gin.Context) {
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

	var req UpdateDictTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	// 转换为领域模型
	domain := domainTools.DictType{
		DictType: modelTools.DictType{
			CoreModels: models.CoreModels{
				Id:         req.Id,
				Sort:       req.Sort,
				Version:    req.Version,
				Modifier:   uid,
				BelongDept: user.DeptId,
				Remark:     req.Remark,
			},
			Status:    req.Status,
			Name:      req.Name,
			DictTag:   req.DictTag,
			DictColor: req.DictColor,
			DictId:    req.DictId,
		},
	}

	if err := h.svc.Update(ctx, domain); err != nil {
		switch {
		case errors.Is(err, serviceTools.ErrDictTypeDuplicate):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典信息已存在", nil)
			return
		case errors.Is(err, serviceTools.ErrDictTypeVersionInconsistency):
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
			return
		default:
			ctx.Set("internal", err.Error())
			zap.L().Error("更新字典信息失败", zap.Error(err))
			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
			return
		}
	}

	response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
}

// GetById
// @Summary 获取字典信息
// @Description 获取指定id字典信息
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param id path string true "id"
// @Success 200 {object} domainTools.DictType
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/getById/{id} [get]
// @Security LoginToken
func (h *dictTypeHandler) GetById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" || len(id) == 0 {
		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "id不能为空", nil)
		return
	}

	detail, err := h.svc.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, serviceTools.ErrDictTypeNotFound) {
			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "字典信息不存在", nil)
			return
		}
		ctx.Set("internal", err.Error())
		zap.L().Error("获取字典信息失败", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
}

// GetListByDictNames
// @Summary 根据字典名称批量查询字典项
// @Description 返回分层结构的字典项映射
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param dictNames body []string true "字典名称数组"
// @Success 200 {object} map[string][]domainTools.DictType
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/listByDictNames [post]
// @Security LoginToken
func (h *dictTypeHandler) GetListByDictNames(ctx *gin.Context) {
	var req []string
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
		return
	}

	list, err := h.svc.GetByDictNames(ctx, req)
	if err != nil {
		zap.L().Error("获取字典名称批量查询字典项异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
}

// GetListPage
// @Summary 获取字典信息分页列表
// @Description 获取字典信息分页列表
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param page query int true "页码" default(1)
// @Param pageSize query int true "每页数量" default(10)
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "字典信息名称"
// @Param dictTag query string false "标签类型" default(primary)
// @Param dictName query string false "数据字典名称"
// @Param valueType query int true "数据类型" default(1)
// @Param dict_id query string false "数据字典id"
// @Success 200 {object} DictTypeListPageResponse
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/listPage [get]
// @Security LoginToken
func (h *dictTypeHandler) GetListPage(ctx *gin.Context) {
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
	dictTag := ctx.DefaultQuery("dictTag", "")
	dictName := ctx.DefaultQuery("dictName", "")
	valueType, _ := strconv.Atoi(ctx.DefaultQuery("valueType", "0"))
	dictId := ctx.DefaultQuery("dict_id", "")

	filter := domainTools.DictTypeFilter{
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
		DictTag:   dictTag,
		DictName:  dictName,
		ValueType: valueType,
		DictId:    dictId,
	}

	list, total, err := h.svc.GetListPage(ctx, filter)
	if err != nil {
		zap.L().Error("获取字典信息分页列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	response.NewResponse().SuccessResponse(ctx, "查询成功", DictTypeListPageResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// GetListAll
// @Summary 获取所有字典信息
// @Description 获取所有字典信息列表
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/json
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "字典信息名称"
// @Param dictTag query string false "标签类型" default(primary)
// @Param dictName query string false "数据字典名称"
// @Param valueType query int true "数据类型" default(1)
// @Param dict_id query string false "数据字典id"
// @Success 200 {array} []domainTools.DictType
// @Failure 400 {object} response.Response
// @Router /v1/tools/dictType/listAll [get]
// @Security LoginToken
func (h *dictTypeHandler) GetListAll(ctx *gin.Context) {
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
	dictTag := ctx.DefaultQuery("dictTag", "")
	dictName := ctx.DefaultQuery("dictName", "")
	valueType, _ := strconv.Atoi(ctx.DefaultQuery("valueType", "0"))
	dictId := ctx.DefaultQuery("dict_id", "")

	filter := domainTools.DictTypeFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status:    status,
		Name:      name,
		DictTag:   dictTag,
		DictName:  dictName,
		ValueType: valueType,
		DictId:    dictId,
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
// @Summary 导出字典信息
// @Description 导出字典信息到Excel文件
// @Tags 系统工具/字典信息管理
// @Accept application/json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param creator query string false "创建人"
// @Param modifier query string false "修改人"
// @Param status query bool false "状态" default(true)
// @Param name query string false "字典信息名称"
// @Param dictTag query string false "标签类型" default(primary)
// @Param dictName query string false "数据字典名称"
// @Param valueType query int true "数据类型" default(1)
// @Param dict_id query string false "数据字典id"
// @Success 200 {file} file "Excel文件"
// @Failure 500 {object} response.Response
// @Router /v1/tools/dictType/export [get]
// @Security LoginToken
func (h *dictTypeHandler) Export(ctx *gin.Context) {
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
	dictTag := ctx.DefaultQuery("dictTag", "")
	dictName := ctx.DefaultQuery("dictName", "")
	valueType, _ := strconv.Atoi(ctx.DefaultQuery("valueType", "0"))
	dictId := ctx.DefaultQuery("dict_id", "")

	filter := domainTools.DictTypeFilter{
		Filters: filters.Filters{
			Creator:    creator,
			Modifier:   modifier,
			BelongDept: user.DeptId,
		},
		Status:    status,
		Name:      name,
		DictTag:   dictTag,
		DictName:  dictName,
		ValueType: valueType,
		DictId:    dictId,
	}

	list, err := h.svc.GetListAll(ctx, filter)
	if err != nil {
		zap.L().Error("获取列表异常", zap.Error(err))
		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
		return
	}

	// 准备导出配置
	filename := fmt.Sprintf("字典信息导出_%s.xlsx", time.Now().Format("20060102150405"))
	cfg := excelutil.ExcelExportConfig{
		SheetName:  "字典信息",
		FileName:   filename,
		StreamMode: true,
		Columns: []excelutil.ExcelColumn{
			{Title: "字典项名称", Field: "Name", Width: 22},
			{Title: "字符串-值", Field: "StrValue", Width: 17},
			{Title: "整型-值", Field: "IntValue", Width: 17},
			{Title: "布尔-值", Field: "BoolValue", Width: 17},
			{
				Title: "标签类型",
				Field: "DictTag",
				Width: 15,
				Formatter: func(value interface{}) string {
					typeValidValues := []string{"primary", "success", "warning", "danger", "info"}
					converter := enumconv.NewEnumConverter(dictType.DictTagMapping, dictType.DictTagImportMapping, typeValidValues, "标签类型")
					str, _ := converter.FromEnum(value.(dictType.DictTagConst))
					return str
				},
			},
			{Title: "标签颜色", Field: "DictColor", Width: 17},
			{Title: "字典名称", Field: "DictName", Width: 17},
			{
				Title: "数据类型",
				Field: "ValueType",
				Width: 15,
				Formatter: func(value interface{}) string {
					typeValueValidValues := []string{"字符串", "整型", "布尔"}
					typeValueConverter := enumconv.NewEnumConverter(dict.TypeValueMapping, dict.TypeValueImportMapping, typeValueValidValues, "数据类型")
					str, _ := typeValueConverter.FromEnum(value.(dict.TypeValueConst))
					return str
				},
			},
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
		zap.L().Error("导出数据字典失败", zap.Error(err))
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
