/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package knowledgegraph

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/handler/coze"
	baseModel "github.com/coze-dev/coze-studio/backend/api/model/base"
	"github.com/coze-dev/coze-studio/backend/api/internal/httputil"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	"github.com/coze-dev/coze-studio/backend/domain/knowledgegraph/entity"
	"github.com/coze-dev/coze-studio/backend/domain/knowledgegraph/service"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// GraphConstructionService 图谱构建服务
var GraphConstructionService *service.GraphConstructionService

// BuildGraphFromText 从文本构建知识图谱
// @router /api/knowledge_graph/build_from_text [POST]
func BuildGraphFromText(ctx context.Context, c *app.RequestContext) {
	var err error
	var req BuildGraphFromTextRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// 获取租户ID和用户ID
	session := ctxutil.GetSession(ctx)
	if session == nil {
		baseModel.HttpError(c, errno.ErrUserAuthenticationFailed)
		return
	}

	tenantID := session.TenantID

	// 执行图谱构建
	err = GraphConstructionService.BuildGraphFromText(ctx, tenantID, req.Text)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &BuildGraphFromTextResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Message:      "知识图谱构建成功",
	}

	httputil.BuildSuccessResp(c, resp)
}

// ExtractEntities 从文本中抽取实体
// @router /api/knowledge_graph/extract_entities [POST]
func ExtractEntities(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ExtractEntitiesRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	session := ctxutil.GetSession(ctx)
	if session == nil {
		baseModel.HttpError(c, errno.ErrUserAuthenticationFailed)
		return
	}

	tenantID := session.TenantID

	// 抽取实体
	entities, err := GraphConstructionService.ExtractEntities(ctx, tenantID, req.Text)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为响应格式
	entityList := make([]*EntityInfo, 0, len(entities))
	for _, ent := range entities {
		entityList = append(entityList, &EntityInfo{
			ID:         ent.ID,
			EntityType: string(ent.EntityType),
			EntityName: ent.EntityName,
			Properties: ent.Properties,
		})
	}

	resp := &ExtractEntitiesResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Entities:     entityList,
		Count:        len(entityList),
	}

	httputil.BuildSuccessResp(c, resp)
}

// ExtractRelationships 从文本中抽取关系
// @router /api/knowledge_graph/extract_relationships [POST]
func ExtractRelationships(ctx context.Context, c *app.RequestContext) {
	var err error
	var req ExtractRelationshipsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	session := ctxutil.GetSession(ctx)
	if session == nil {
		baseModel.HttpError(c, errno.ErrUserAuthenticationFailed)
		return
	}

	tenantID := session.TenantID

	// 构建实体列表
	entities := make([]*entity.GraphEntity, 0, len(req.Entities))
	for _, entReq := range req.Entities {
		entities = append(entities, &entity.GraphEntity{
			ID:         entReq.ID,
			TenantID:   tenantID,
			EntityType: entity.EntityType(entReq.EntityType),
			EntityName: entReq.EntityName,
			Properties: entReq.Properties,
		})
	}

	// 抽取关系
	relationships, err := GraphConstructionService.ExtractRelationships(ctx, tenantID, req.Text, entities)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为响应格式
	relList := make([]*RelationshipInfo, 0, len(relationships))
	for _, rel := range relationships {
		relList = append(relList, &RelationshipInfo{
			ID:             rel.ID,
			SourceEntityID: rel.SourceEntityID,
			TargetEntityID: rel.TargetEntityID,
			RelationType:   string(rel.RelationType),
			Properties:     rel.Properties,
			Weight:         rel.Weight,
		})
	}

	resp := &ExtractRelationshipsResponse{
		BaseResponse:   baseModel.NewBaseResponse(),
		Relationships:  relList,
		Count:          len(relList),
	}

	httputil.BuildSuccessResp(c, resp)
}

// AddEntity 添加实体
// @router /api/knowledge_graph/add_entity [POST]
func AddEntity(ctx context.Context, c *app.RequestContext) {
	var err error
	var req AddEntityRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	session := ctxutil.GetSession(ctx)
	if session == nil {
		baseModel.HttpError(c, errno.ErrUserAuthenticationFailed)
		return
	}

	tenantID := session.TenantID

	// 构建实体
	ent := &entity.GraphEntity{
		TenantID:   tenantID,
		EntityType: entity.EntityType(req.EntityType),
		EntityName: req.EntityName,
		Properties: req.Properties,
	}

	// 添加实体
	if err := GraphConstructionService.AddEntity(ctx, ent); err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &AddEntityResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		EntityID:     ent.ID,
		Message:      "实体添加成功",
	}

	httputil.BuildSuccessResp(c, resp)
}

// AddRelationship 添加关系
// @router /api/knowledge_graph/add_relationship [POST]
func AddRelationship(ctx context.Context, c *app.RequestContext) {
	var err error
	var req AddRelationshipRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	session := ctxutil.GetSession(ctx)
	if session == nil {
		baseModel.HttpError(c, errno.ErrUserAuthenticationFailed)
		return
	}

	tenantID := session.TenantID

	// 构建关系
	rel := &entity.GraphRelationship{
		TenantID:       tenantID,
		SourceEntityID: req.SourceEntityID,
		TargetEntityID: req.TargetEntityID,
		RelationType:   entity.RelationType(req.RelationType),
		Properties:     req.Properties,
		Weight:         req.Weight,
	}

	// 添加关系
	if err := GraphConstructionService.AddRelationship(ctx, rel); err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &AddRelationshipResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		RelationshipID: rel.ID,
		Message:      "关系添加成功",
	}

	httputil.BuildSuccessResp(c, resp)
}

// DeleteEntity 删除实体
// @router /api/knowledge_graph/delete_entity [POST]
func DeleteEntity(ctx context.Context, c *app.RequestContext) {
	var err error
	var req DeleteEntityRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	session := ctxutil.GetSession(ctx)
	if session == nil {
		baseModel.HttpError(c, errno.ErrUserAuthenticationFailed)
		return
	}

	// 删除实体
	if err := GraphConstructionService.DeleteEntity(ctx, req.EntityID); err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &DeleteEntityResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Message:      "实体删除成功",
	}

	httputil.BuildSuccessResp(c, resp)
}

// GetEntityStats 获取实体统计
// @router /api/knowledge_graph/entity_stats [POST]
func GetEntityStats(ctx context.Context, c *app.RequestContext) {
	var err error
	var req baseModel.BaseRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	session := ctxutil.GetSession(ctx)
	if session == nil {
		baseModel.HttpError(c, errno.ErrUserAuthenticationFailed)
		return
	}

	tenantID := session.TenantID

	// 获取统计
	stats, err := GraphConstructionService.GetEntityStats(ctx, tenantID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为响应格式
	statsMap := make(map[string]int64)
	for k, v := range stats {
		statsMap[string(k)] = v
	}

	resp := &GetEntityStatsResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Stats:        statsMap,
	}

	httputil.BuildSuccessResp(c, resp)
}

// ========================================
// 请求和响应结构
// ========================================

// BuildGraphFromTextRequest 从文本构建图谱请求
type BuildGraphFromTextRequest struct {
	Text string `json:"text" validate:"required"`
}

// BuildGraphFromTextResponse 从文本构建图谱响应
type BuildGraphFromTextResponse struct {
	baseModel.BaseResponse
	Message string `json:"message"`
}

// EntityInfo 实体信息
type EntityInfo struct {
	ID         string                 `json:"id"`
	EntityType string                 `json:"entity_type"`
	EntityName string                 `json:"entity_name"`
	Properties entity.EntityProperties `json:"properties"`
}

// ExtractEntitiesRequest 抽取实体请求
type ExtractEntitiesRequest struct {
	Text string `json:"text" validate:"required"`
}

// ExtractEntitiesResponse 抽取实体响应
type ExtractEntitiesResponse struct {
	baseModel.BaseResponse
	Entities []*EntityInfo `json:"entities"`
	Count    int           `json:"count"`
}

// RelationshipInfo 关系信息
type RelationshipInfo struct {
	ID             string                       `json:"id"`
	SourceEntityID string                       `json:"source_entity_id"`
	TargetEntityID string                       `json:"target_entity_id"`
	RelationType   string                       `json:"relation_type"`
	Properties     entity.RelationshipProperties `json:"properties"`
	Weight         float32                      `json:"weight"`
}

// ExtractRelationshipsRequest 抽取关系请求
type ExtractRelationshipsRequest struct {
	Text     string       `json:"text" validate:"required"`
	Entities []*EntityInfo `json:"entities"`
}

// ExtractRelationshipsResponse 抽取关系响应
type ExtractRelationshipsResponse struct {
	baseModel.BaseResponse
	Relationships []*RelationshipInfo `json:"relationships"`
	Count         int                 `json:"count"`
}

// AddEntityRequest 添加实体请求
type AddEntityRequest struct {
	EntityName string                 `json:"entity_name" validate:"required"`
	EntityType string                 `json:"entity_type" validate:"required"`
	Properties entity.EntityProperties `json:"properties"`
}

// AddEntityResponse 添加实体响应
type AddEntityResponse struct {
	baseModel.BaseResponse
	EntityID string `json:"entity_id"`
	Message  string `json:"message"`
}

// AddRelationshipRequest 添加关系请求
type AddRelationshipRequest struct {
	SourceEntityID string                       `json:"source_entity_id" validate:"required"`
	TargetEntityID string                       `json:"target_entity_id" validate:"required"`
	RelationType   string                       `json:"relation_type" validate:"required"`
	Properties     entity.RelationshipProperties `json:"properties"`
	Weight         float32                      `json:"weight"`
}

// AddRelationshipResponse 添加关系响应
type AddRelationshipResponse struct {
	baseModel.BaseResponse
	RelationshipID string `json:"relationship_id"`
	Message        string `json:"message"`
}

// DeleteEntityRequest 删除实体请求
type DeleteEntityRequest struct {
	EntityID string `json:"entity_id" validate:"required"`
}

// DeleteEntityResponse 删除实体响应
type DeleteEntityResponse struct {
	baseModel.BaseResponse
	Message string `json:"message"`
}

// GetEntityStatsResponse 获取实体统计响应
type GetEntityStatsResponse struct {
	baseModel.BaseResponse
	Stats map[string]int64 `json:"stats"`
}
