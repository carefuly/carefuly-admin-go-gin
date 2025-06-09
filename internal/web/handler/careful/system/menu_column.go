/**
 * Description：
 * FileName：menu_column.go
 * Author：CJiaの用心
 * Create：2025/6/9 14:43:31
 * Remark：
 */

package system

// // CreateMenuColumnRequest 创建
// type CreateMenuColumnRequest struct {
// 	Title  string `json:"title" binding:"required,max=64"`        // 标题
// 	Field  string `json:"field" binding:"required,max=64"`        // 字段名
// 	Width  int    `json:"width" binding:"required" default:"150"` // 宽度
// 	MenuId string `json:"menuId" binding:"required,max=100"`      // 菜单ID
// 	Sort   int    `json:"sort" binding:"omitempty" default:"1"`   // 排序
// 	Remark string `json:"remark" binding:"omitempty,max=255"`     // 备注
// }
//
// // UpdateMenuColumnRequest 更新
// type UpdateMenuColumnRequest struct {
// 	Id      string `json:"id" binding:"required"`                  // 主键ID
// 	Title   string `json:"title" binding:"required,max=64"`        // 标题
// 	Field   string `json:"field" binding:"required,max=64"`        // 字段名
// 	Width   int    `json:"width" binding:"required" default:"150"` // 宽度
// 	Sort    int    `json:"sort" binding:"omitempty" default:"1"`   // 排序
// 	Version int    `json:"version" binding:"omitempty"`            // 版本
// 	Remark  string `json:"remark" binding:"omitempty,max=255"`     // 备注
// }
//
// // MenuColumnListPageResponse 列表分页响应
// type MenuColumnListPageResponse struct {
// 	List     []domainSystem.MenuColumn `json:"list"`     // 列表
// 	Total    int64                     `json:"total"`    // 总数
// 	Page     int                       `json:"page"`     // 页码
// 	PageSize int                       `json:"pageSize"` // 每页数量
// }
//
// type MenuColumnHandler interface {
// }
//
// type menuColumnHandler struct {
// 	rely    config.RelyConfig
// 	svc     system.MenuColumnService
// 	userSvc system.UserService
// }
//
// func NewMenuColumnHandler(rely config.RelyConfig, svc system.MenuColumnService, userSvc system.UserService) MenuColumnHandler {
// 	return &menuColumnHandler{
// 		rely:    rely,
// 		svc:     svc,
// 		userSvc: userSvc,
// 	}
// }
//
// // RegisterRoutes 注册路由
// func (h *menuColumnHandler) RegisterRoutes(router *gin.RouterGroup) {
// 	// base := router.Group("/menuButton")
// 	// base.POST("/create", h.Create)
// 	// base.POST("/import", h.Import)
// 	// base.DELETE("/delete/:id", h.Delete)
// 	// base.POST("/delete/batchDelete", h.BatchDelete)
// 	// base.PUT("/update", h.Update)
// 	// base.GET("/getById/:id", h.GetById)
// 	// base.GET("/listPage", h.GetListPage)
// 	// base.GET("/listAll", h.GetListAll)
// }
//
// // Create
// // @Summary 创建菜单数据列
// // @Description 创建菜单数据列
// // @Tags 系统管理/菜单数据列管理
// // @Accept application/json
// // @Produce application/json
// // @Param CreateMenuColumnRequest body CreateMenuColumnRequest true "请求"
// // @Success 200 {object} response.Response
// // @Failure 400 {object} response.Response
// // @Router /v1/system/menuColumn/create [post]
// // @Security LoginToken
// func (h *menuColumnHandler) Create(ctx *gin.Context) {
// 	uid, ok := ctx.MustGet("userId").(string)
// 	if !ok {
// 		ctx.Set("internal", uid)
// 		zap.S().Error("用户ID获取失败", uid)
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	user, err := h.userSvc.GetById(ctx, uid)
// 	if err != nil {
// 		ctx.Set("internal", err)
// 		zap.S().Error("获取用户失败", err)
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	var req CreateMenuColumnRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
// 		return
// 	}
//
// 	// 转换为领域模型
// 	domain := domainSystem.MenuColumn{
// 		MenuColumn: modelSystem.MenuColumn{
// 			CoreModels: models.CoreModels{
// 				Sort:       req.Sort,
// 				Creator:    uid,
// 				Modifier:   uid,
// 				BelongDept: user.DeptId,
// 				Remark:     req.Remark,
// 			},
// 			Title:   req.Title,
// 			Field:   req.Field,
// 			// MenuId: req.MenuId,
// 		},
// 	}
//
// 	if err := h.svc.Create(ctx, domain); err != nil {
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	response.NewResponse().SuccessResponse(ctx, "新增成功", nil)
// }
//
// // Delete
// // @Summary 删除菜单数据列
// // @Description 删除指定id菜单数据列
// // @Tags 系统管理/菜单数据列管理
// // @Accept application/json
// // @Produce application/json
// // @Param id path string true "id"
// // @Success 200 {object} response.Response
// // @Failure 400 {object} response.Response
// // @Router /v1/system/menuColumn/delete/{id} [delete]
// // @Security LoginToken
// func (h *menuColumnHandler) Delete(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	if id == "" || len(id) == 0 {
// 		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "ID不能为空", nil)
// 		return
// 	}
//
// 	if err := h.svc.Delete(ctx, id); err != nil {
// 		if errors.Is(err, serviceSystem.ErrMenuColumnNotFound) {
// 			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "菜单数据列不存在", nil)
// 			return
// 		}
// 		ctx.Set("internal", err.Error())
// 		zap.L().Error("删除菜单数据列失败", zap.Error(err))
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	response.NewResponse().SuccessResponse(ctx, "删除成功", nil)
// }
//
// // BatchDelete
// // @Summary 批量删除菜单权限
// // @Description 批量删除菜单权限
// // @Tags 系统管理/菜单权限管理
// // @Accept application/json
// // @Produce application/json
// // @Param ids body []string true "id数组"
// // @Success 200 {object} response.Response
// // @Failure 400 {object} response.Response
// // @Router /v1/system/menuColumn/delete/batchDelete [post]
// // @Security LoginToken
// func (h *menuColumnHandler) BatchDelete(ctx *gin.Context) {
// 	var ids []string
// 	if err := ctx.ShouldBindJSON(&ids); err != nil {
// 		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
// 		return
// 	}
//
// 	err := h.svc.BatchDelete(ctx, ids)
// 	if err != nil {
// 		ctx.Set("internal", err.Error())
// 		zap.L().Error("批量删除菜单权限异常", zap.Error(err))
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	response.NewResponse().SuccessResponse(ctx, "批量删除成功", nil)
// }
//
// // Update
// // @Summary 更新菜单权限
// // @Description 更新菜单权限信息
// // @Tags 系统管理/菜单权限管理
// // @Accept application/json
// // @Produce application/json
// // @Param UpdateMenuButtonRequest body UpdateMenuButtonRequest true "请求"
// // @Success 200 {object} response.Response
// // @Failure 400 {object} response.Response
// // @Router /v1/system/menuButton/update [put]
// // @Security LoginToken
// func (h *menuColumnHandler) Update(ctx *gin.Context) {
// 	uid, ok := ctx.MustGet("userId").(string)
// 	if !ok {
// 		ctx.Set("internal", uid)
// 		zap.S().Error("用户ID获取失败", uid)
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	user, err := h.userSvc.GetById(ctx, uid)
// 	if err != nil {
// 		ctx.Set("internal", err)
// 		zap.S().Error("获取用户失败", err)
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	var req UpdateMenuButtonRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		validate.NewValidatorError(h.rely.Trans).HandleValidatorError(ctx, err)
// 		return
// 	}
//
// 	// 转换为领域模型
// 	domain := domainSystem.MenuButton{
// 		MenuButton: modelSystem.MenuButton{
// 			CoreModels: models.CoreModels{
// 				Id:         req.Id,
// 				Sort:       req.Sort,
// 				Version:    req.Version,
// 				Modifier:   uid,
// 				BelongDept: user.DeptId,
// 				Remark:     req.Remark,
// 			},
// 			Name:   req.Name,
// 			Code:   req.Code,
// 			Api:    req.Api,
// 			Method: req.Method,
// 		},
// 	}
//
// 	if err := h.svc.Update(ctx, domain); err != nil {
// 		switch {
// 		case errors.Is(err, serviceSystem.ErrMenuButtonVersionInconsistency):
// 			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "数据版本不一致，取消修改，请刷新后重试", nil)
// 			return
// 		default:
// 			zap.L().Error("更新菜单权限失败", zap.Error(err))
// 			response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 			return
// 		}
// 	}
//
// 	response.NewResponse().SuccessResponse(ctx, "更新成功", nil)
// }
//
// // GetById
// // @Summary 获取菜单权限
// // @Description 获取指定id菜单权限信息
// // @Tags 系统管理/菜单权限管理
// // @Accept application/json
// // @Produce application/json
// // @Param id path string true "id"
// // @Success 200 {object} domainSystem.MenuButton
// // @Failure 400 {object} response.Response
// // @Router /v1/system/menuButton/getById/{id} [get]
// // @Security LoginToken
// func (h *menuColumnHandler) GetById(ctx *gin.Context) {
// 	id := ctx.Param("id")
// 	if id == "" {
// 		response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "用户ID不能为空", nil)
// 		return
// 	}
//
// 	detail, err := h.svc.GetById(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, serviceSystem.ErrMenuButtonNotFound) {
// 			response.NewResponse().ErrorResponse(ctx, http.StatusBadRequest, "菜单权限不存在", nil)
// 			return
// 		}
// 		zap.L().Error("获取菜单权限失败", zap.Error(err))
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	response.NewResponse().SuccessResponse(ctx, "获取成功", detail)
// }
//
// // GetListPage
// // @Summary 获取菜单权限分页列表
// // @Description 获取菜单权限分页列表
// // @Tags 系统管理/菜单权限管理
// // @Accept application/json
// // @Produce application/json
// // @Param page query int true "页码" default(1)
// // @Param pageSize query int true "每页数量" default(10)
// // @Param creator query string false "创建人"
// // @Param modifier query string false "修改人"
// // @Param status query bool false "状态" default(true)
// // @Param name query string false "名称"
// // @Param code query string false "权限值"
// // @Param menu_id query string false "菜单ID"
// // @Success 200 {object} MenuButtonListPageResponse
// // @Failure 400 {object} response.Response
// // @Router /v1/system/menuButton/listPage [get]
// // @Security LoginToken
// func (h *menuColumnHandler) GetListPage(ctx *gin.Context) {
// 	uid, ok := ctx.MustGet("userId").(string)
// 	if !ok {
// 		ctx.Set("internal", uid)
// 		zap.S().Error("用户ID获取失败", uid)
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	user, err := h.userSvc.GetById(ctx, uid)
// 	if err != nil {
// 		ctx.Set("internal", err)
// 		zap.S().Error("获取用户失败", err)
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
// 	creator := ctx.DefaultQuery("creator", "")
// 	modifier := ctx.DefaultQuery("modifier", "")
// 	status, _ := strconv.ParseBool(ctx.DefaultQuery("status", "true"))
//
// 	name := ctx.DefaultQuery("name", "")
// 	code := ctx.DefaultQuery("code", "")
// 	menuId := ctx.DefaultQuery("menu_id", "")
//
// 	filter := domainSystem.MenuButtonFilter{
// 		Pagination: filters.Pagination{
// 			Page:     page,
// 			PageSize: pageSize,
// 		},
// 		Filters: filters.Filters{
// 			Creator:    creator,
// 			Modifier:   modifier,
// 			BelongDept: user.DeptId,
// 		},
// 		Status: status,
// 		Name:   name,
// 		Code:   code,
// 		MenuId: menuId,
// 	}
//
// 	list, total, err := h.svc.GetListPage(ctx, filter)
// 	if err != nil {
// 		zap.L().Error("获取分页列表异常", zap.Error(err))
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	response.NewResponse().SuccessResponse(ctx, "查询成功", MenuButtonListPageResponse{
// 		List:     list,
// 		Total:    total,
// 		Page:     page,
// 		PageSize: pageSize,
// 	})
// }
//
// // GetListAll
// // @Summary 获取所有菜单权限
// // @Description 获取所有菜单权限列表
// // @Tags 系统管理/菜单权限管理
// // @Accept application/json
// // @Produce application/json
// // @Param creator query string false "创建人"
// // @Param modifier query string false "修改人"
// // @Param status query bool false "状态" default(true)
// // @Param name query string false "名称"
// // @Param code query string false "权限值"
// // @Param menu_id query string false "菜单ID"
// // @Success 200 {array} []domainSystem.MenuButton
// // @Failure 400 {object} response.Response
// // @Router /v1/system/menuButton/listAll [get]
// // @Security LoginToken
// func (h *menuColumnHandler) GetListAll(ctx *gin.Context) {
// 	uid, ok := ctx.MustGet("userId").(string)
// 	if !ok {
// 		ctx.Set("internal", uid)
// 		zap.S().Error("用户ID获取失败", uid)
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	user, err := h.userSvc.GetById(ctx, uid)
// 	if err != nil {
// 		ctx.Set("internal", err)
// 		zap.S().Error("获取用户失败", err)
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	creator := ctx.DefaultQuery("creator", "")
// 	modifier := ctx.DefaultQuery("modifier", "")
// 	status, _ := strconv.ParseBool(ctx.DefaultQuery("status", "true"))
//
// 	name := ctx.DefaultQuery("name", "")
// 	code := ctx.DefaultQuery("code", "")
// 	menuId := ctx.DefaultQuery("menu_id", "")
//
// 	filter := domainSystem.MenuButtonFilter{
// 		Filters: filters.Filters{
// 			Creator:    creator,
// 			Modifier:   modifier,
// 			BelongDept: user.DeptId,
// 		},
// 		Status: status,
// 		Name:   name,
// 		Code:   code,
// 		MenuId: menuId,
// 	}
//
// 	list, err := h.svc.GetListAll(ctx, filter)
// 	if err != nil {
// 		zap.L().Error("获取列表异常", zap.Error(err))
// 		response.NewResponse().ErrorResponse(ctx, http.StatusInternalServerError, "服务器异常", nil)
// 		return
// 	}
//
// 	response.NewResponse().SuccessResponse(ctx, "查询成功", list)
// }

