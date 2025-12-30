// backend/types/errno/cache_test.go
package errno

import (
	"net/http"
	"testing"
)

// TestCacheErrorCodes 测试缓存错误码的完整性
func TestCacheErrorCodes(t *testing.T) {
	tests := []struct {
		name       string
		errCode    *BaseErrorCode
		wantCode   string
		wantStatus int
	}{
		{
			name:       "ErrCacheMiss",
			errCode:    ErrCacheMiss,
			wantCode:   "CACHE404001",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "ErrCacheGetFailed",
			errCode:    ErrCacheGetFailed,
			wantCode:   "CACHE500002",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "ErrCacheDecodeFailed",
			errCode:    ErrCacheDecodeFailed,
			wantCode:   "CACHE500007",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "ErrCacheEncodeFailed",
			errCode:    ErrCacheEncodeFailed,
			wantCode:   "CACHE500006",
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
