// backend/types/errno/storage.go
package errno

import "net/http"

// 存储模块错误码(STORAGE前缀)
// 遵循ZKER统一错误码定义规范

const (
	ErrInitStorageFailedCode         = 61001
	ErrUploadFileFailedCode          = 61002
	ErrDownloadFileFailedCode        = 61003
	ErrGenerateDownloadURLFailedCode = 61004
	ErrGenerateUploadURLFailedCode   = 61005
	ErrDeleteFileFailedCode          = 61006
	ErrDeleteFilePartialFailedCode   = 61007
	ErrGetFileInfoFailedCode         = 61008
	ErrCheckBucketFailedCode         = 61009
	ErrDeleteBucketFailedCode        = 61010
	ErrCreateBucketFailedCode        = 61011
)

var (
	// 存储初始化错误 (500)
	ErrInitStorageFailed = &BaseErrorCode{
		code:       "STORAGE61001",
		message:    "Failed to initialize storage",
		messageZH:  "存储初始化失败",
		messageEN:  "Failed to initialize storage",
		httpStatus: http.StatusInternalServerError,
	}

	// 文件上传错误 (500)
	ErrUploadFileFailed = &BaseErrorCode{
		code:       "STORAGE61002",
		message:    "Failed to upload file",
		messageZH:  "文件上传失败",
		messageEN:  "Failed to upload file",
		httpStatus: http.StatusInternalServerError,
	}

	// 文件下载错误 (500)
	ErrDownloadFileFailed = &BaseErrorCode{
		code:       "STORAGE61003",
		message:    "Failed to download file",
		messageZH:  "文件下载失败",
		messageEN:  "Failed to download file",
		httpStatus: http.StatusInternalServerError,
	}

	// 生成下载URL错误 (500)
	ErrGenerateDownloadURLFailed = &BaseErrorCode{
		code:       "STORAGE61004",
		message:    "Failed to generate download URL",
		messageZH:  "生成下载URL失败",
		messageEN:  "Failed to generate download URL",
		httpStatus: http.StatusInternalServerError,
	}

	// 生成上传URL错误 (500)
	ErrGenerateUploadURLFailed = &BaseErrorCode{
		code:       "STORAGE61005",
		message:    "Failed to generate upload URL",
		messageZH:  "生成上传URL失败",
		messageEN:  "Failed to generate upload URL",
		httpStatus: http.StatusInternalServerError,
	}

	// 删除文件错误 (500)
	ErrDeleteFileFailed = &BaseErrorCode{
		code:       "STORAGE61006",
		message:    "Failed to delete file",
		messageZH:  "删除文件失败",
		messageEN:  "Failed to delete file",
		httpStatus: http.StatusInternalServerError,
	}

	// 部分删除文件失败 (500)
	ErrDeleteFilePartialFailed = &BaseErrorCode{
		code:       "STORAGE61007",
		message:    "Partial files failed to delete",
		messageZH:  "部分文件删除失败",
		messageEN:  "Partial files failed to delete",
		httpStatus: http.StatusInternalServerError,
	}

	// 获取文件信息错误 (500)
	ErrGetFileInfoFailed = &BaseErrorCode{
		code:       "STORAGE61008",
		message:    "Failed to get file info",
		messageZH:  "获取文件信息失败",
		messageEN:  "Failed to get file info",
		httpStatus: http.StatusInternalServerError,
	}

	// 检查Bucket错误 (500)
	ErrCheckBucketFailed = &BaseErrorCode{
		code:       "STORAGE61009",
		message:    "Failed to check bucket",
		messageZH:  "检查Bucket失败",
		messageEN:  "Failed to check bucket",
		httpStatus: http.StatusInternalServerError,
	}

	// 删除Bucket错误 (500)
	ErrDeleteBucketFailed = &BaseErrorCode{
		code:       "STORAGE61010",
		message:    "Failed to delete bucket",
		messageZH:  "删除Bucket失败",
		messageEN:  "Failed to delete bucket",
		httpStatus: http.StatusInternalServerError,
	}

	// 创建Bucket错误 (500)
	ErrCreateBucketFailed = &BaseErrorCode{
		code:       "STORAGE61011",
		message:    "Failed to create bucket",
		messageZH:  "创建Bucket失败",
		messageEN:  "Failed to create bucket",
		httpStatus: http.StatusInternalServerError,
	}
)
