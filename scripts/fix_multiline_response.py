#!/usr/bin/env python3
"""
修复bot_store_service.go中所有多行响应格式
"""

import re

file_path = "backend/api/handler/coze/bot_store_service.go"

with open(file_path, 'r', encoding='utf-8') as f:
    content = f.read()

# 定义所有替换规则
replacements = [
    # UpdateBotStoreItemResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.UpdateBotStoreItemResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*\}\)',
        'httputil.BuildSuccessResp(c, nil)'
    ),
    # ListBotStoreItemsResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.ListBotStoreItemsResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*Data: &botstoreAPI\.BotStoreItemListData\{\s*Items:.*?\}\s*\}\s*\})',
        'httputil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{\n\t\tItems:      items,\n\t\tTotal:      resp.Total,\n\t\tPage:       resp.Page,\n\t\tPageSize:   resp.PageSize,\n\t\tTotalPages: resp.TotalPages,\n\t})'
    ),
    # SearchBotStoreItemsResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.SearchBotStoreItemsResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*Data: &botstoreAPI\.BotStoreItemListData\{.*?\}\s*\}\s*\})',
        'httputil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{\n\t\tItems:      items,\n\t\tTotal:      resp.Total,\n\t\tPage:       resp.Page,\n\t\tPageSize:   resp.PageSize,\n\t\tTotalPages: resp.TotalPages,\n\t})'
    ),
    # GetBotStoreItemResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.GetBotStoreItemResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*Data: entityToDTO\(item\),\s*\}\)',
        'httputil.BuildSuccessResp(c, entityToDTO(item))'
    ),
    # GetBotCategoriesResponse - 注意这里有data（小写d）
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.GetBotCategoriesResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*data: dtos,\s*\}\)',
        'httputil.BuildSuccessResp(c, dtos)'
    ),
    # GetPendingReviewsResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.GetPendingReviewsResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*Data: &botstoreAPI\.BotStoreItemListData\{.*?\}\s*\}\s*\})',
        'httutil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{\n\t\tItems:      dtos,\n\t\tTotal:      total,\n\t\tPage:       req.Page,\n\t\tPageSize:   req.PageSize,\n\t\tTotalPages: totalPages,\n\t})'
    ),
    # ReviewBotStoreItemResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.ReviewBotStoreItemResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*\}\)',
        'httputil.BuildSuccessResp(c, nil)'
    ),
    # CreateReviewResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.CreateReviewResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*Data: reviewToDTO\(review\),\s*\}\)',
        'httputil.BuildSuccessResp(c, reviewToDTO(review))'
    ),
    # GetReviewsResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.GetReviewsResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*Data: &botstoreAPI\.ReviewListData\{.*?\}\s*\}\s*\})',
        'httputil.BuildSuccessResp(c, &botstoreAPI.ReviewListData{\n\t\tReviews:       reviews,\n\t\tTotal:         resp.Total,\n\t\tPage:          resp.Page,\n\t\tPageSize:      resp.PageSize,\n\t\tTotalPages:    resp.TotalPages,\n\t\tAverageRating: resp.AverageRating,\n\t\tRatingCount:   resp.RatingCount,\n\t})'
    ),
    # UpdateReviewResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.UpdateReviewResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*\}\)',
        'httputil.BuildSuccessResp(c, nil)'
    ),
    # DeleteReviewResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.DeleteReviewResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*\}\)',
        'httputil.BuildSuccessResp(c, nil)'
    ),
    # GetReviewStatisticsResponse
    (
        r'c\.JSON\(consts\.StatusOK, &botstoreAPI\.GetReviewStatisticsResponse\{\s*BaseResponse: botstoreAPI\.BaseResponse\{\s*Code: 0,\s*Msg:  "success",\s*\},\s*Data: &botstoreAPI\.ReviewStatisticsDTO\{.*?\}\s*\}\s*\})',
        'httputil.BuildSuccessResp(c, &botstoreAPI.ReviewStatisticsDTO{\n\t\tAverageRating: stats.AverageRating,\n\t\tRatingCount:   stats.RatingCount,\n\t\tRating1Count:  stats.Rating1Count,\n\t\tRating2Count:  stats.Rating2Count,\n\t\tRating3Count:  stats.Rating3Count,\n\t\tRating4Count:  stats.Rating4Count,\n\t\tRating5Count:  stats.Rating5Count,\n\t})'
    ),
    # handleReviewError中的错误响应
    (
        r'c\.JSON\(consts\.BadRequest, &botstoreAPI\.BaseResponse\{\s*Code: consts\.BadRequest,\s*Msg:  err\.Error\(\),\s*\}\)',
        'httputil.BuildErrorResp(c, errno.ErrInternalErrorCode, err.Error(), "处理失败", nil)'
    ),
    # invalidParamRequestResponse
    (
        r'c\.JSON\(consts\.BadRequest, &botstoreAPI\.BaseResponse\{\s*Code: consts\.BadRequest,\s*Msg:  msg,\s*\}\)',
        'httputil.BuildErrorResp(c, errno.ErrInvalidParamCode, msg, "参数验证失败", nil)'
    ),
    # internalServerErrorResponse
    (
        r'c\.JSON\(consts\.InternalServerError, &botstoreAPI\.BaseResponse\{\s*Code: consts\.InternalServerError,\s*Msg:  err\.Error\(\),\s*\}\)',
        'httputil.BuildErrorResp(c, errno.ErrInternalErrorCode, err.Error(), "内部服务器错误", nil)'
    ),
]

# 由于正则表达式难以处理复杂的多行结构，我们手动处理
manual_fixes = [
    # UpdateBotStoreItemResponse
    ('UpdateBotStoreItemResponse', '''
	httputil.BuildSuccessResp(c, nil)
'''),
    # ListBotStoreItemsResponse
    ('ListBotStoreItemsResponse', '''
	httputil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{
		Items:      items,
		Total:      resp.Total,
		Page:       resp.Page,
		PageSize:   resp.PageSize,
		TotalPages: resp.TotalPages,
	})
'''),
    # SearchBotStoreItemsResponse
    ('SearchBotStoreItemsResponse', '''
	httputil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{
		Items:      items,
		Total:      resp.Total,
		Page:       resp.Page,
		PageSize:   resp.PageSize,
		TotalPages: resp.TotalPages,
	})
'''),
    # GetBotStoreItemResponse
    ('GetBotStoreItemResponse', '''
	httputil.BuildSuccessResp(c, entityToDTO(item))
'''),
    # GetBotCategoriesResponse
    ('GetBotCategoriesResponse', '''
	httputil.BuildSuccessResp(c, dtos)
'''),
    # GetPendingReviewsResponse
    ('GetPendingReviewsResponse', '''
	httputil.BuildSuccessResp(c, &botstoreAPI.BotStoreItemListData{
		Items:      dtos,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	})
'''),
    # ReviewBotStoreItemResponse
    ('ReviewBotStoreItemResponse', '''
	httputil.BuildSuccessResp(c, nil)
'''),
    # CreateReviewResponse
    ('CreateReviewResponse', '''
	httputil.BuildSuccessResp(c, reviewToDTO(review))
'''),
    # GetReviewsResponse
    ('GetReviewsResponse', '''
	httputil.BuildSuccessResp(c, &botstoreAPI.ReviewListData{
		Reviews:       reviews,
		Total:         resp.Total,
		Page:          resp.Page,
		PageSize:      resp.PageSize,
		TotalPages:    resp.TotalPages,
		AverageRating: resp.AverageRating,
		RatingCount:   resp.RatingCount,
	})
'''),
    # UpdateReviewResponse
    ('UpdateReviewResponse', '''
	httputil.BuildSuccessResp(c, nil)
'''),
    # DeleteReviewResponse
    ('DeleteReviewResponse', '''
	httputil.BuildSuccessResp(c, nil)
'''),
    # GetReviewStatisticsResponse
    ('GetReviewStatisticsResponse', '''
	httputil.BuildSuccessResp(c, &botstoreAPI.ReviewStatisticsDTO{
		AverageRating: stats.AverageRating,
		RatingCount:   stats.RatingCount,
		Rating1Count:  stats.Rating1Count,
		Rating2Count:  stats.Rating2Count,
		Rating3Count:  stats.Rating3Count,
		Rating4Count:  stats.Rating4Count,
		Rating5Count:  stats.Rating5Count,
	})
'''),
]

# 手动替换（使用更简单的方法）
lines = content.split('\n')
new_lines = []
i = 0
while i < len(lines):
    line = lines[i]
    # 检查是否是需要替换的行
    if 'c.JSON(consts.StatusOK, &botstoreAPI.' in line:
        # 查找响应类型
        for response_type, replacement in manual_fixes:
            if response_type in line:
                # 跳过整个多行响应
                new_lines.append(replacement.strip())
                # 跳过后续行直到找到匹配的})
                brace_count = line.count('{') - line.count('}')
                i += 1
                while i < len(lines) and brace_count > 0:
                    brace_count += lines[i].count('{') - lines[i].count('}')
                    i += 1
                break
        else:
            new_lines.append(line)
            i += 1
    elif 'c.JSON(consts.StatusBadRequest, &botstoreAPI.BaseResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrInvalidParamCode, err.Error(), "参数验证失败", nil)')
        i += 2  # 跳过后续两行
    elif 'c.JSON(consts.StatusInternalServerError, &botstoreAPI.BaseResponse{' in line:
        new_lines.append('\thttputil.BuildErrorResp(c, errno.ErrInternalErrorCode, err.Error(), "内部服务器错误", nil)')
        i += 2  # 跳过后续两行
    else:
        new_lines.append(line)
        i += 1

content = '\n'.join(new_lines)

with open(file_path, 'w', encoding='utf-8') as f:
    f.write(content)

print("✅ 已修复 bot_store_service.go")
