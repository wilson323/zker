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

package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

//go:embed translations/*.json
var translationFiles embed.FS

// Language 语言类型
type Language string

const (
	LanguageZH   Language = "zh" // 中文
	LanguageEN   Language = "en" // 英文
	LanguageDefault Language = LanguageEN
)

// Translator 翻译器
type Translator struct {
	mu         sync.RWMutex
	translations map[Language]map[int32]string // language -> error_code -> message
}

var (
	globalTranslator *Translator
	translatorOnce   sync.Once
)

// GetGlobalTranslator 获取全局翻译器实例
func GetGlobalTranslator() *Translator {
	translatorOnce.Do(func() {
		globalTranslator = NewTranslator()
		if err := globalTranslator.LoadTranslations(); err != nil {
			logs.Errorf("Failed to load translations: %v", err)
		}
	})
	return globalTranslator
}

// NewTranslator 创建新的翻译器
func NewTranslator() *Translator {
	return &Translator{
		translations: make(map[Language]map[int32]string),
	}
}

// LoadTranslations 加载所有翻译文件
func (t *Translator) LoadTranslations() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 加载中文翻译
	if err := t.loadTranslationFile(LanguageZH); err != nil {
		return fmt.Errorf("failed to load Chinese translations: %w", err)
	}

	// 加载英文翻译
	if err := t.loadTranslationFile(LanguageEN); err != nil {
		return fmt.Errorf("failed to load English translations: %w", err)
	}

	logs.Infof("Loaded %d languages with translations", len(t.translations))
	return nil
}

// loadTranslationFile 加载单个语言翻译文件
func (t *Translator) loadTranslationFile(lang Language) error {
	filename := fmt.Sprintf("translations/%s.json", lang)

	data, err := translationFiles.ReadFile(filename)
	if err != nil {
		return err
	}

	var translations map[int32]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return fmt.Errorf("failed to parse translation file %s: %w", filename, err)
	}

	t.translations[lang] = translations
	logs.Infof("Loaded %d translations for language: %s", len(translations), lang)
	return nil
}

// Translate 翻译错误码到指定语言
func (t *Translator) Translate(code int32, lang Language) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 如果请求的语言不存在，回退到默认语言
	if _, exists := t.translations[lang]; !exists {
		lang = LanguageDefault
	}

	messages, exists := t.translations[lang]
	if !exists {
		return "", false
	}

	message, exists := messages[code]
	return message, exists
}

// TranslateOrDefault 翻译错误码，如果不存在则返回默认消息
func (t *Translator) TranslateOrDefault(code int32, lang Language, defaultMsg string) string {
	msg, exists := t.Translate(code, lang)
	if !exists {
		return defaultMsg
	}
	return msg
}

// TranslateWithFallback 翻译错误码，如果指定语言不存在则回退到英文，再回退到中文，最后回退到默认消息
func (t *Translator) TranslateWithFallback(code int32, lang Language, defaultMsg string) string {
	// 尝试指定语言
	if msg, exists := t.Translate(code, lang); exists {
		return msg
	}

	// 回退到英文
	if lang != LanguageEN {
		if msg, exists := t.Translate(code, LanguageEN); exists {
			return msg
		}
	}

	// 回退到中文
	if lang != LanguageZH {
		if msg, exists := t.Translate(code, LanguageZH); exists {
			return msg
		}
	}

	// 最后回退到默认消息
	return defaultMsg
}

// GetSupportedLanguages 获取支持的语言列表
func (t *Translator) GetSupportedLanguages() []Language {
	t.mu.RLock()
	defer t.mu.RUnlock()

	languages := make([]Language, 0, len(t.translations))
	for lang := range t.translations {
		languages = append(languages, lang)
	}
	return languages
}

// AddTranslation 动态添加翻译（用于测试或运行时更新）
func (t *Translator) AddTranslation(lang Language, code int32, message string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.translations[lang]; !exists {
		t.translations[lang] = make(map[int32]string)
	}

	t.translations[lang][code] = message
}

// 便捷函数：使用全局翻译器

// Translate 翻译错误码
func Translate(code int32, lang Language) (string, bool) {
	return GetGlobalTranslator().Translate(code, lang)
}

// TranslateOrDefault 翻译错误码，提供默认消息
func TranslateOrDefault(code int32, lang Language, defaultMsg string) string {
	return GetGlobalTranslator().TranslateOrDefault(code, lang, defaultMsg)
}

// TranslateWithFallback 翻译错误码，支持多语言回退
func TranslateWithFallback(code int32, lang Language, defaultMsg string) string {
	return GetGlobalTranslator().TranslateWithFallback(code, lang, defaultMsg)
}
