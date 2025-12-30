// backend/types/errno/storage_test.go
package errno

import (
	"net/http"
	"testing"
)

// TestStorageErrorCodes 测试存储错误码的完整性
func TestStorageErrorCodes(t *testing.T) {
	tests := []struct {
		name       string
		errCode    *BaseErrorCode
		wantCode   string
		wantStatus int
	}{
		{
			name:       "ErrInitStorageFailed",
			errCode:    ErrInitStorageFailed,
			wantCode:   "STORAGE61001",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "ErrUploadFileFailed",
			errCode:    ErrUploadFileFailed,
			wantCode:   "STORAGE61002",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "ErrDownloadFileFailed",
			errCode:    ErrDownloadFileFailed,
			wantCode:   "STORAGE61003",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.errCode.code != tt.wantCode {
				t.Errorf("code = %v, want %v", tt.errCode.code, tt.wantCode)
			}
			if tt.errCode.httpStatus != tt.wantStatus {
				t.Errorf("httpStatus = %v, want %v", tt.errCode.httpStatus, tt.wantStatus)
			}
			if tt.errCode.messageZH == "" {
				t.Error("messageZH is empty")
			}
			if tt.errCode.messageEN == "" {
				t.Error("messageEN is empty")
			}
		})
	}
}
