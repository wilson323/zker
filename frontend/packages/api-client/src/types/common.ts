// frontend/packages/api-client/src/types/common.ts

/**
 * 通用分页参数
 */
export interface PageRequest {
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

/**
 * 通用分页响应
 */
export interface PageResponse<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
  has_next: boolean;
  has_prev: boolean;
}

/**
 * 通用API响应
 */
export interface APIResponse<T> {
  code: string;
  message: string;
  message_en?: string;
  data: T;
  request_id?: string;
  timestamp: string;
}

/**
 * 时间戳
 */
export type Timestamp = string;

/**
 * 实体ID
 */
export type EntityID = string;

/**
 * 租户ID
 */
export type TenantID = string;

/**
 * 用户ID
 */
export type UserID = string;
