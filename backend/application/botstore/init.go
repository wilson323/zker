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

package botstore

import (
	"github.com/coze-dev/coze-studio/backend/domain/botstore/repository"
	botstoreService "github.com/coze-dev/coze-studio/backend/domain/botstore/service"
	"gorm.io/gorm"
)

// BotStoreService Bot商店服务容器
type BotStoreService struct {
	Publisher botstoreService.BotStorePublisher
	Browser   botstoreService.BotStoreBrowser
	Reviewer  botstoreService.BotStoreReviewer
}

var (
	botStoreService *BotStoreService
)

// InitBotStoreService 初始化Bot商店服务
func InitBotStoreService(db *gorm.DB) error {
	// 初始化Repository (使用repository包的构造函数)
	itemRepo := repository.NewBotStoreRepository(db)
	categoryRepo := repository.NewBotStoreCategoryRepository(db)

	// 初始化Service
	botStoreService = &BotStoreService{
		Publisher: botstoreService.NewBotStorePublisher(itemRepo, categoryRepo),
		Browser:   botstoreService.NewBotStoreBrowser(itemRepo, categoryRepo),
		Reviewer:  botstoreService.NewBotStoreReviewer(itemRepo),
	}

	return nil
}

// GetBotStoreService 获取Bot商店服务
func GetBotStoreService() *BotStoreService {
	return botStoreService
}

// GetPublisher 获取发布服务
func GetPublisher() botstoreService.BotStorePublisher {
	if botStoreService != nil {
		return botStoreService.Publisher
	}
	return nil
}

// GetBrowser 获取浏览服务
func GetBrowser() botstoreService.BotStoreBrowser {
	if botStoreService != nil {
		return botStoreService.Browser
	}
	return nil
}

// GetReviewer 获取审核服务
func GetReviewer() botstoreService.BotStoreReviewer {
	if botStoreService != nil {
		return botStoreService.Reviewer
	}
	return nil
}
