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
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	"github.com/coze-dev/coze-studio/backend/domain/knowledgegraph/entity"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// GraphQueryService 图谱查询服务
var GraphQueryService *service.GraphQueryService

// QueryEntity 查询单个实体
// @router /api/knowledge_graph/query_entity [POST]
func QueryEntity(ctx context.Context, c *app.RequestContext) {
	var err error
	var req QueryEntityRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// 查询实体
	ent, err := GraphQueryService.QueryEntity(ctx, req.EntityID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &QueryEntityResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Entity: &EntityInfo{
			ID:         ent.ID,
			EntityType: string(ent.EntityType),
			EntityName: ent.EntityName,
			Properties: ent.Properties,
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// QueryNeighbors 查询实体邻居
// @router /api/knowledge_graph/query_neighbors [POST]
func QueryNeighbors(ctx context.Context, c *app.RequestContext) {
	var err error
	var req QueryNeighborsRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// 查询邻居
	subgraph, err := GraphQueryService.QueryNeighbors(ctx, req.EntityID, req.Depth)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为响应格式
	neighbors := make([]*EntityInfo, 0, len(subgraph.Neighbors))
	for _, node := range subgraph.Neighbors {
		neighbors = append(neighbors, &EntityInfo{
			ID:         node.ID,
			EntityType: string(node.EntityType),
			EntityName: node.EntityName,
			Properties: node.Properties,
		})
	}

	relationships := make([]*RelationshipInfo, 0, len(subgraph.Relationships))
	for _, rel := range subgraph.Relationships {
		relationships = append(relationships, &RelationshipInfo{
			ID:           rel.ID,
			SourceID:     rel.SourceID,
			TargetID:     rel.TargetID,
			RelationType: string(rel.RelationType),
			Properties:   rel.Properties,
		})
	}

	resp := &QueryNeighborsResponse{
		BaseResponse:   baseModel.NewBaseResponse(),
		CenterNode: &EntityInfo{
			ID:         subgraph.CenterNode.ID,
			EntityType: string(subgraph.CenterNode.EntityType),
			EntityName: subgraph.CenterNode.EntityName,
			Properties: subgraph.CenterNode.Properties,
		},
		Neighbors:     neighbors,
		Relationships: relationships,
		Depth:         subgraph.Depth,
	}

	c.JSON(consts.StatusOK, resp)
}

// QueryPath 查询两个实体之间的路径
// @router /api/knowledge_graph/query_path [POST]
func QueryPath(ctx context.Context, c *app.RequestContext) {
	var err error
	var req QueryPathRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// 查询路径
	paths, err := GraphQueryService.QueryPath(ctx, req.SourceEntityID, req.TargetEntityID, req.MaxDepth)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	// 转换为响应格式
	pathList := make([]*PathInfo, 0, len(paths))
	for _, path := range paths {
		pathList = append(pathList, &PathInfo{
			Nodes:         path.Nodes,
			Relationships: path.Relationships,
			Length:        path.Length,
			Weight:        path.Weight,
		})
	}

	resp := &QueryPathResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Paths:        pathList,
		Count:        len(pathList),
	}

	c.JSON(consts.StatusOK, resp)
}

// SearchEntities 搜索实体
// @router /api/knowledge_graph/search_entities [POST]
func SearchEntities(ctx context.Context, c *app.RequestContext) {
	var err error
	var req SearchEntitiesRequest
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

	// 搜索实体
	entities, err := GraphQueryService.SearchEntities(ctx, tenantID, req.Keyword, entity.EntityType(req.EntityType), req.Limit)
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

	resp := &SearchEntitiesResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Entities:     entityList,
		Count:        len(entityList),
	}

	c.JSON(consts.StatusOK, resp)
}

// GetGraphStatistics 获取图谱统计信息
// @router /api/knowledge_graph/statistics [POST]
func GetGraphStatistics(ctx context.Context, c *app.RequestContext) {
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

	// 获取统计信息
	stats, err := GraphQueryService.GetGraphStatistics(ctx, tenantID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &GetGraphStatisticsResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Statistics: &GraphStatisticsInfo{
			NodeCount:           stats.NodeCount,
			RelationshipCount:   stats.RelationshipCount,
			AvgDegree:           stats.AvgDegree,
			Density:             stats.Density,
			ConnectedComponents: stats.ConnectedComponents,
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// FindShortestPath 查找最短路径
// @router /api/knowledge_graph/shortest_path [POST]
func FindShortestPath(ctx context.Context, c *app.RequestContext) {
	var err error
	var req QueryPathRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// 查找最短路径
	path, err := GraphQueryService.FindShortestPath(ctx, req.SourceEntityID, req.TargetEntityID)
	if err != nil {
		coze.InternalServerErrorResponse(ctx, c, err)
		return
	}

	resp := &FindShortestPathResponse{
		BaseResponse: baseModel.NewBaseResponse(),
		Path: &PathInfo{
			Nodes:         path.Nodes,
			Relationships: path.Relationships,
			Length:        path.Length,
			Weight:        path.Weight,
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// ========================================
// 请求和响应结构
// ========================================

// QueryEntityRequest 查询实体请求
type QueryEntityRequest struct {
	EntityID string `json:"entity_id" validate:"required"`
}

// QueryEntityResponse 查询实体响应
type QueryEntityResponse struct {
	baseModel.BaseResponse
	Entity *EntityInfo `json:"entity"`
}

// QueryNeighborsRequest 查询邻居请求
type QueryNeighborsRequest struct {
	EntityID string `json:"entity_id" validate:"required"`
	Depth    int    `json:"depth" validate:"min=1,max=5"`
}

// QueryNeighborsResponse 查询邻居响应
type QueryNeighborsResponse struct {
	baseModel.BaseResponse
	CenterNode    *EntityInfo       `json:"center_node"`
	Neighbors     []*EntityInfo     `json:"neighbors"`
	Relationships []*RelationshipInfo `json:"relationships"`
	Depth         int               `json:"depth"`
}

// RelationshipInfo 关系信息（用于查询响应）
type RelationshipInfo struct {
	ID           string                       `json:"id"`
	SourceID     string                       `json:"source_id"`
	TargetID     string                       `json:"target_id"`
	RelationType string                       `json:"relation_type"`
	Properties   entity.RelationshipProperties `json:"properties"`
}

// QueryPathRequest 查询路径请求
type QueryPathRequest struct {
	SourceEntityID string `json:"source_entity_id" validate:"required"`
	TargetEntityID string `json:"target_entity_id" validate:"required"`
	MaxDepth       int    `json:"max_depth" validate:"min=1,max=10"`
}

// QueryPathResponse 查询路径响应
type QueryPathResponse struct {
	baseModel.BaseResponse
	Paths []*PathInfo `json:"paths"`
	Count int         `json:"count"`
}

// PathInfo 路径信息
type PathInfo struct {
	Nodes         []string `json:"nodes"`
	Relationships []string `json:"relationships"`
	Length        int      `json:"length"`
	Weight        float32  `json:"weight"`
}

// SearchEntitiesRequest 搜索实体请求
type SearchEntitiesRequest struct {
	Keyword    string `json:"keyword"`
	EntityType string `json:"entity_type"`
	Limit      int    `json:"limit" validate:"min=1,max=100"`
}

// SearchEntitiesResponse 搜索实体响应
type SearchEntitiesResponse struct {
	baseModel.BaseResponse
	Entities []*EntityInfo `json:"entities"`
	Count    int           `json:"count"`
}

// GetGraphStatisticsResponse 获取图谱统计响应
type GetGraphStatisticsResponse struct {
	baseModel.BaseResponse
	Statistics *GraphStatisticsInfo `json:"statistics"`
}

// GraphStatisticsInfo 图谱统计信息
type GraphStatisticsInfo struct {
	NodeCount           int     `json:"node_count"`
	RelationshipCount   int     `json:"relationship_count"`
	AvgDegree           float32 `json:"avg_degree"`
	Density             float32 `json:"density"`
	ConnectedComponents int     `json:"connected_components"`
}

// FindShortestPathResponse 查找最短路径响应
type FindShortestPathResponse struct {
	baseModel.BaseResponse
	Path *PathInfo `json:"path"`
}
