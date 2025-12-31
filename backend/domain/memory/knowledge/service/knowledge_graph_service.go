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

package service

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/memory/knowledge/entity"
)

// KnowledgeGraphService 知识图谱服务接口
type KnowledgeGraphService interface {
	// ExtractEntitiesAndRelations 从文本中抽取实体和关系
	ExtractEntitiesAndRelations(
		ctx context.Context,
		req *entity.ExtractEntitiesRequest,
	) (*entity.ExtractEntitiesResponse, error)

	// Query 查询知识图谱
	Query(
		ctx context.Context,
		req *entity.QueryKnowledgeGraphRequest,
	) (*entity.QueryKnowledgeGraphResponse, error)

	// Visualize 生成可视化数据
	Visualize(
		ctx context.Context,
		req *entity.VisualizeKnowledgeGraphRequest,
	) (*entity.VisualizeKnowledgeGraphResponse, error)

	// AddEntity 添加实体
	AddEntity(
		ctx context.Context,
		entity *entity.GraphEntity,
	) error

	// AddRelation 添加关系
	AddRelation(
		ctx context.Context,
		relation *entity.GraphRelation,
	) error

	// DeleteEntity 删除实体
	DeleteEntity(ctx context.Context, entityID string) error

	// DeleteRelation 删除关系
	DeleteRelation(ctx context.Context, relationID string) error

	// GetEntity 获取实体
	GetEntity(
		ctx context.Context,
		entityID string,
	) (*entity.GraphEntity, error)

	// GetEntityRelations 获取实体的所有关系
	GetEntityRelations(
		ctx context.Context,
		entityID string,
	) ([]entity.GraphRelation, error)

	// GetGraphStats 获取图谱统计
	GetGraphStats(
		ctx context.Context,
		conversationID string,
	) (*entity.GraphStats, error)
}
