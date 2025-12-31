// frontend/apps/coze-studio/src/router/types.ts

import { RouteObject } from 'react-router-dom';

/**
 * 路由元信息
 */
export interface RouteMeta {
  /** 路由标题 */
  title?: string;
  /** 是否需要登录 */
  requireAuth?: boolean;
  /** 需要的权限(格式: resource:action) */
  permission?: string;
  /** 是否需要选择租户 */
  requireTenant?: boolean;
  /** 布局组件 */
  layout?: 'default' | 'empty' | 'fullscreen';
  /** 是否隐藏(在菜单中) */
  hidden?: boolean;
  /** 图标 */
  icon?: string;
  /** 排序 */
  order?: number;
}

/**
 * 扩展路由对象
 */
export interface ExtRouteObject extends Omit<RouteObject, 'children'> {
  /** 路由元信息 */
  meta?: RouteMeta;
  /** 子路由 */
  children?: ExtRouteObject[];
  /** 路由路径(相对于父路由) */
  path?: string;
}

/**
 * 路由菜单项
 */
export interface MenuItem {
  /** 菜单key */
  key: string;
  /** 菜单标题 */
  label: string;
  /** 菜单图标 */
  icon?: string;
  /** 路由路径 */
  path: string;
  /** 子菜单 */
  children?: MenuItem[];
  /** 排序 */
  order?: number;
}
