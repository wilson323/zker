# 测试生成 Skill

## 技能描述

为代码生成符合规范的测试用例，包括单元测试、集成测试和端到端测试。

## 适用场景

- 新功能开发完成后需要编写测试
- 现有代码缺少测试需要补充
- 提高测试覆盖率
- 测试驱动开发（TDD）

## 工作流程

### 1. 前端测试生成

#### 1.1 React 组件测试

**测试文件位置**: `{component-name}.test.tsx`

```typescript
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { UserProfile } from './user-profile';

// Mock 外部依赖
vi.mock('@coze-studio/api', () => ({
  userAPI: {
    getUserById: vi.fn(),
  },
}));

describe('UserProfile', () => {
  const defaultProps = {
    userId: 123,
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  // 1. 基础渲染测试
  describe('rendering', () => {
    it('should render component', () => {
      render(<UserProfile {...defaultProps} />);
      expect(screen.getByTestId('user-profile')).toBeInTheDocument();
    });

    it('should render user name when data is loaded', async () => {
      vi.mocked(userAPI.getUserById).mockResolvedValue({
        id: 123,
        name: 'John Doe',
      });

      render(<UserProfile {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText('John Doe')).toBeInTheDocument();
      });
    });
  });

  // 2. 加载状态测试
  describe('when loading', () => {
    it('should show spinner', () => {
      render(<UserProfile {...defaultProps} isLoading />);
      expect(screen.getByRole('progressbar')).toBeInTheDocument();
    });

    it('should not show user data', () => {
      render(<UserProfile {...defaultProps} isLoading />);
      expect(screen.queryByText('John Doe')).not.toBeInTheDocument();
    });
  });

  // 3. 错误状态测试
  describe('when error occurs', () => {
    it('should show error message', async () => {
      vi.mocked(userAPI.getUserById).mockRejectedValue(
        new Error('Failed to fetch')
      );

      render(<UserProfile {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByText(/failed to fetch/i)).toBeInTheDocument();
      });
    });
  });

  // 4. 用户交互测试
  describe('user interactions', () => {
    it('should call onEdit when edit button is clicked', async () => {
      const onEdit = vi.fn();
      vi.mocked(userAPI.getUserById).mockResolvedValue({
        id: 123,
        name: 'John Doe',
      });

      render(<UserProfile {...defaultProps} onEdit={onEdit} />);

      await waitFor(() => {
        fireEvent.click(screen.getByRole('button', { name: /edit/i }));
      });

      expect(onEdit).toHaveBeenCalledTimes(1);
    });
  });

  // 5. 条件渲染测试
  describe('conditional rendering', () => {
    it('should show edit button when user has permission', async () => {
      vi.mocked(userAPI.getUserById).mockResolvedValue({
        id: 123,
        name: 'John Doe',
      });

      render(<UserProfile {...defaultProps} hasPermission />);

      await waitFor(() => {
        expect(screen.getByRole('button', { name: /edit/i })).toBeInTheDocument();
      });
    });

    it('should not show edit button when user lacks permission', async () => {
      vi.mocked(userAPI.getUserById).mockResolvedValue({
        id: 123,
        name: 'John Doe',
      });

      render(<UserProfile {...defaultProps} hasPermission={false} />);

      await waitFor(() => {
        expect(screen.queryByRole('button', { name: /edit/i })).not.toBeInTheDocument();
      });
    });
  });

  // 6. Props 变化测试
  describe('when props change', () => {
    it('should refetch data when userId changes', async () => {
      const { rerender } = render(<UserProfile {...defaultProps} userId={123} />);

      await waitFor(() => {
        expect(userAPI.getUserById).toHaveBeenCalledWith(123);
      });

      rerender(<UserProfile {...defaultProps} userId={456} />);

      await waitFor(() => {
        expect(userAPI.getUserById).toHaveBeenCalledWith(456);
      });
    });
  });
});
```

#### 1.2 Hooks 测试

**测试文件位置**: `{hook-name}.test.ts`

```typescript
import { renderHook, act, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { useAuthForm } from './use-auth-form';

describe('useAuthForm', () => {
  it('should initialize with default values', () => {
    const { result } = renderHook(() => useAuthForm());

    expect(result.current.email).toBe('');
    expect(result.current.password).toBe('');
    expect(result.current.loading).toBe(false);
  });

  it('should update email when setEmail is called', () => {
    const { result } = renderHook(() => useAuthForm());

    act(() => {
      result.current.setEmail('test@example.com');
    });

    expect(result.current.email).toBe('test@example.com');
  });

  it('should call login API when handleSubmit is called', async () => {
    const mockLogin = vi.fn().mockResolvedValue({ token: 'abc123' });
    const { result } = renderHook(() => useAuthForm(mockLogin));

    act(() => {
      result.current.setEmail('test@example.com');
      result.current.setPassword('password123');
    });

    await act(async () => {
      await result.current.handleSubmit();
    });

    expect(mockLogin).toHaveBeenCalledWith({
      email: 'test@example.com',
      password: 'password123',
    });
  });

  it('should set loading to true during submission', async () => {
    const mockLogin = vi.fn();
    const { result } = renderHook(() => useAuthForm(mockLogin));

    act(() => {
      result.current.handleSubmit();
    });

    expect(result.current.loading).toBe(true);
  });
});
```

#### 1.3 API 测试

**测试文件位置**: `{api-name}.test.ts`

```typescript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { UserAPI } from './user-api';
import { HTTPClient } from '@coze-studio/http-client';

describe('UserAPI', () => {
  let api: UserAPI;
  let mockHTTP: HTTPClient;

  beforeEach(() => {
    mockHTTP = {
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn(),
    } as unknown as HTTPClient;

    api = new UserAPI(mockHTTP);
  });

  describe('getUserById', () => {
    it('should fetch user by ID', async () => {
      const mockUser = { id: 123, name: 'John Doe' };
      vi.mocked(mockHTTP.get).mockResolvedValue({ data: mockUser });

      const result = await api.getUserById(123);

      expect(mockHTTP.get).toHaveBeenCalledWith('/api/v1/users/123');
      expect(result).toEqual(mockUser);
    });

    it('should throw error when fetch fails', async () => {
      vi.mocked(mockHTTP.get).mockRejectedValue(new Error('Network error'));

      await expect(api.getUserById(123)).rejects.toThrow('Network error');
    });
  });

  describe('createUser', () => {
    it('should create new user', async () => {
      const newUser = { name: 'John', email: 'john@example.com' };
      const createdUser = { id: 123, ...newUser };
      vi.mocked(mockHTTP.post).mockResolvedValue({ data: createdUser });

      const result = await api.createUser(newUser);

      expect(mockHTTP.post).toHaveBeenCalledWith('/api/v1/users', newUser);
      expect(result).toEqual(createdUser);
    });
  });
});
```

### 2. 后端测试生成

#### 2.1 服务测试

**测试文件位置**: `service_test.go`

```go
package service_test

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
    "gorm.io/gorm"

    "github.com/coze-dev/coze-studio/backend/domain/user/entity"
    "github.com/coze-dev/coze-studio/backend/domain/user/repository"
    "github.com/coze-dev/coze-studio/backend/domain/user/service"
)

// Mock Repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) (int64, error) {
    args := m.Called(ctx, user)
    return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (*entity.User, bool, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, false, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Bool(1), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    args := m.Called(ctx, email)
    return args.Bool(0), args.Error(1)
}

func TestCreateUser(t *testing.T) {
    tests := []struct {
        name    string
        req     *service.CreateUserRequest
        setup   func(*MockUserRepository)
        want    *entity.User
        wantErr bool
        errCode int
    }{
        {
            name: "successful creation",
            req: &service.CreateUserRequest{
                Name:    "John Doe",
                Email:   "john@example.com",
                SpaceID: 1,
                OwnerID: 1,
            },
            setup: func(m *MockUserRepository) {
                m.On("ExistsByEmail", mock.Anything, "john@example.com").Return(false, nil)
                m.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(int64(123), nil)
            },
            want: &entity.User{
                ID:    123,
                Name:  "John Doe",
                Email: "john@example.com",
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            req: &service.CreateUserRequest{
                Name:    "John Doe",
                Email:   "john@example.com",
                SpaceID: 1,
                OwnerID: 1,
            },
            setup: func(m *MockUserRepository) {
                m.On("ExistsByEmail", mock.Anything, "john@example.com").Return(true, nil)
            },
            wantErr: true,
            errCode: errno.ErrUserAlreadyExists,
        },
        {
            name: "repository error",
            req: &service.CreateUserRequest{
                Name:    "John Doe",
                Email:   "john@example.com",
                SpaceID: 1,
                OwnerID: 1,
            },
            setup: func(m *MockUserRepository) {
                m.On("ExistsByEmail", mock.Anything, "john@example.com").Return(false, nil)
                m.On("Create", mock.Anything, mock.Anything).Return(int64(0), assert.AnError)
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            ctx := context.Background()
            mockRepo := new(MockUserRepository)
            tt.setup(mockRepo)

            svc := service.NewService(&service.Components{
                UserRepo: mockRepo,
            })

            // Act
            got, err := svc.CreateUser(ctx, tt.req)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
                if tt.errCode != 0 {
                    assert.Equal(t, tt.errCode, errorx.Code(err))
                }
            } else {
                require.NoError(t, err)
                assert.NotNil(t, got)
                assert.Equal(t, tt.want.Name, got.Name)
                assert.Equal(t, tt.want.Email, got.Email)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}

func TestGetUserByID(t *testing.T) {
    t.Run("user exists", func(t *testing.T) {
        // Arrange
        ctx := context.Background()
        mockRepo := new(MockUserRepository)
        mockUser := &entity.User{ID: 123, Name: "John Doe"}
        mockRepo.On("FindByID", mock.Anything, int64(123)).Return(mockUser, true, nil)

        svc := service.NewService(&service.Components{
            UserRepo: mockRepo,
        })

        // Act
        got, err := svc.GetUserByID(ctx, 123)

        // Assert
        require.NoError(t, err)
        assert.NotNil(t, got)
        assert.Equal(t, mockUser, got)
        mockRepo.AssertExpectations(t)
    })

    t.Run("user not found", func(t *testing.T) {
        // Arrange
        ctx := context.Background()
        mockRepo := new(MockUserRepository)
        mockRepo.On("FindByID", mock.Anything, int64(123)).Return(nil, false, nil)

        svc := service.NewService(&service.Components{
            UserRepo: mockRepo,
        })

        // Act
        got, err := svc.GetUserByID(ctx, 123)

        // Assert
        require.Error(t, err)
        assert.Nil(t, got)
        assert.Equal(t, errno.ErrUserNotFound, errorx.Code(err))
        mockRepo.AssertExpectations(t)
    })
}
```

#### 2.2 HTTP Handler 测试

```go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/coze-dev/coze-studio/backend/api/handler"
    "github.com/coze-dev/coze-studio/backend/domain/user/service"
)

func TestCreateUserHandler(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    interface{}
        setup          func(*MockUserService)
        expectedStatus int
        expectedBody   interface{}
    }{
        {
            name: "successful creation",
            requestBody: map[string]interface{}{
                "name":  "John Doe",
                "email": "john@example.com",
            },
            setup: func(m *MockUserService) {
                m.On("CreateUser", mock.Anything, mock.AnythingOfType("*service.CreateUserRequest")).
                    Return(&entity.User{ID: 123, Name: "John Doe"}, nil)
            },
            expectedStatus: http.StatusOK,
            expectedBody: map[string]interface{}{
                "code":    0,
                "message": "success",
                "data": map[string]interface{}{
                    "user": map[string]interface{}{
                        "id":   float64(123),
                        "name": "John Doe",
                    },
                },
            },
        },
        {
            name: "invalid request",
            requestBody: map[string]interface{}{
                "name": "",
            },
            setup:          func(m *MockUserService) {},
            expectedStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockService := new(MockUserService)
            tt.setup(mockService)

            h := &handler.UserHandler{
                UserService: mockService,
            }

            body, _ := json.Marshal(tt.requestBody)
            req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()

            // Act
            h.CreateUser(w, req)

            // Assert
            assert.Equal(t, tt.expectedStatus, w.Code)

            var response map[string]interface{}
            json.Unmarshal(w.Body.Bytes(), &response)
            assert.Equal(t, tt.expectedBody, response)

            mockService.AssertExpectations(t)
        })
    }
}
```

### 3. 测试覆盖要求

#### 3.1 前端覆盖率要求

| 包层级 | 覆盖率要求 | 增量覆盖率 |
|-------|----------|-----------|
| Level 1 (arch) | 80% | 90% |
| Level 2 (common) | 30% | 60% |
| Level 3-4 | 0% | 灵活 |

#### 3.2 测试场景要求

**必须测试的场景**：
1. ✅ 正常流程（Happy Path）
2. ✅ 边界情况（Boundary Cases）
3. ✅ 错误处理（Error Handling）
4. ✅ 加载状态（Loading States）
5. ✅ 空数据状态（Empty States）
6. ✅ 用户交互（User Interactions）
7. ✅ Props 变化（Props Changes）
8. ✅ 条件渲染（Conditional Rendering）

## 测试模板

### 前端组件测试模板

```typescript
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ComponentName } from './component-name';

describe('ComponentName', () => {
  const defaultProps = {};

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('rendering', () => {
    it('should render component', () => {
      render(<ComponentName {...defaultProps} />);
      expect(screen.getByTestId('component-name')).toBeInTheDocument();
    });
  });

  describe('when loading', () => {
    it('should show loading state', () => {
      render(<ComponentName {...defaultProps} isLoading />);
      // assertions
    });
  });

  describe('when error occurs', () => {
    it('should show error message', () => {
      // test setup and assertions
    });
  });

  describe('user interactions', () => {
    it('should handle user action', () => {
      // test setup and assertions
    });
  });
});
```

### 后端服务测试模板

```go
package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/require"
)

func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   *InputType
        setup   func(*MockRepository)
        want    *OutputType
        wantErr bool
    }{
        {
            name:  "successful case",
            input: &InputType{...},
            setup: func(m *MockRepository) {
                m.On("Method", mock.Anything, mock.Anything).Return(..., nil)
            },
            want:    &OutputType{...},
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := new(MockRepository)
            tt.setup(mockRepo)
            svc := NewService(mockRepo)

            // Act
            got, err := svc.Method(context.Background(), tt.input)

            // Assert
            if tt.wantErr {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

## 输出检查清单

生成测试后，必须确保：

- [ ] 测试文件命名正确
- [ ] 使用 describe/it 或 t.Run 结构
- [ ] 测试描述清晰明确
- [ ] 包含正常流程测试
- [ ] 包含错误处理测试
- [ ] 包含边界情况测试
- [ ] Mock 外部依赖
- [ ] 断言完整准确
- [ ] 测试独立可重复
- [ ] 覆盖率符合要求

## 注意事项

1. **测试独立性**: 每个测试应该独立运行
2. **Mock 使用**: 合理使用 Mock 隔离依赖
3. **断言完整**: 验证所有重要的行为
4. **测试速度**: 单元测试应该快速运行
5. **可维护性**: 测试代码也应该清晰易懂
