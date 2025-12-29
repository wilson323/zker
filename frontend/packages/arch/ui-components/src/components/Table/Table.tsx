// frontend/packages/arch/ui-components/src/components/Table/Table.tsx

import React from 'react';
import { useStyles } from './Table.styles';

/**
 * 列定义接口
 */
export interface Column<T = any> {
  /** 列标题 */
  title: string;
  /** 数据字段名 */
  dataIndex?: keyof T;
  /** 列key */
  key: string;
  /** 自定义渲染函数 */
  render?: (value: any, record: T, index: number) => React.ReactNode;
  /** 列宽度 */
  width?: number | string;
  /** 是否可排序 */
  sorter?: boolean | ((a: T, b: T) => number);
  /** 对齐方式 */
  align?: 'left' | 'center' | 'right';
}

/**
 * 分页配置接口
 */
export interface PaginationConfig {
  /** 当前页 */
  current: number;
  /** 每页条数 */
  pageSize: number;
  /** 总数 */
  total: number;
  /** 页码改变回调 */
  onChange?: (page: number, pageSize: number) => void;
}

/**
 * Table组件Props接口
 */
export interface TableProps<T = any> {
  /** 列定义 */
  columns: Column<T>[];
  /** 数据源 */
  dataSource: T[];
  /** 加载状态 */
  loading?: boolean;
  /** 分页配置 */
  pagination?: PaginationConfig | false;
  /** 是否显示边框 */
  bordered?: boolean;
  /** 行key */
  rowKey?: keyof T | ((record: T) => string);
  /** 自定义className */
  className?: string;
  /** 空数据时的文本 */
  emptyText?: React.ReactNode;
}

/**
 * Table 表格
 *
 * 表格组件，展示行列数据
 *
 * @example
 * ```tsx
 * const columns = [
 *   { title: '姓名', dataIndex: 'name', key: 'name' },
 *   { title: '年龄', dataIndex: 'age', key: 'age' },
 * ];
 *
 * <Table
 *   columns={columns}
 *   dataSource={data}
 *   pagination={{ current: 1, pageSize: 10, total: 100 }}
 * />
 * ```
 */
export const Table = <T extends Record<string, any>>({
  columns,
  dataSource,
  loading = false,
  pagination,
  bordered = false,
  rowKey = 'id',
  className,
  emptyText = '暂无数据',
}: TableProps<T>) => {
  const classes = useStyles({ bordered });

  // 获取行的key
  const getRowKey = (record: T, index: number): string => {
    if (typeof rowKey === 'function') {
      return rowKey(record);
    }
    return String(record[rowKey] ?? index);
  };

  // 渲染表头
  const renderHeader = () => {
    return (
      <thead className={classes.thead}>
        <tr>
          {columns.map((column) => (
            <th
              key={column.key}
              className={classes.th}
              style={{ width: column.width, textAlign: column.align || 'left' }}
            >
              {column.title}
            </th>
          ))}
        </tr>
      </thead>
    );
  };

  // 渲染表格内容
  const renderBody = () => {
    if (loading) {
      return (
        <tbody>
          <tr>
            <td colSpan={columns.length} className={classes.loadingCell}>
              <div className={classes.spinner} />
              <span className={classes.loadingText}>加载中...</span>
            </td>
          </tr>
        </tbody>
      );
    }

    if (dataSource.length === 0) {
      return (
        <tbody>
          <tr>
            <td colSpan={columns.length} className={classes.emptyCell}>
              {emptyText}
            </td>
          </tr>
        </tbody>
      );
    }

    return (
      <tbody>
        {dataSource.map((record, rowIndex) => {
          const key = getRowKey(record, rowIndex);
          return (
            <tr key={key} className={classes.tr}>
              {columns.map((column) => {
                const value = column.dataIndex ? record[column.dataIndex] : undefined;
                const content = column.render
                  ? column.render(value, record, rowIndex)
                  : value;

                return (
                  <td
                    key={column.key}
                    className={classes.td}
                    style={{ textAlign: column.align || 'left' }}
                  >
                    {content}
                  </td>
                );
              })}
            </tr>
          );
        })}
      </tbody>
    );
  };

  // 渲染分页
  const renderPagination = () => {
    if (!pagination) return null;

    const { current, pageSize, total, onChange } = pagination;
    const totalPages = Math.ceil(total / pageSize);
    const pageNumbers = [];

    // 生成页码（简化版，实际项目可能需要更复杂的逻辑）
    let startPage = Math.max(1, current - 2);
    let endPage = Math.min(totalPages, current + 2);

    if (startPage > 1) {
      pageNumbers.push(1);
      if (startPage > 2) {
        pageNumbers.push('...');
      }
    }

    for (let i = startPage; i <= endPage; i++) {
      pageNumbers.push(i);
    }

    if (endPage < totalPages) {
      if (endPage < totalPages - 1) {
        pageNumbers.push('...');
      }
      pageNumbers.push(totalPages);
    }

    return (
      <div className={classes.pagination}>
        <button
          className={classes.pageButton}
          disabled={current === 1}
          onClick={() => onChange?.(current - 1, pageSize)}
        >
          上一页
        </button>

        {pageNumbers.map((page, index) => {
          if (page === '...') {
            return (
              <span key={`ellipsis-${index}`} className={classes.ellipsis}>
                ...
              </span>
            );
          }

          return (
            <button
              key={page}
              className={`${classes.pageButton} ${current === page ? classes.activePage : ''}`}
              onClick={() => onChange?.(Number(page), pageSize)}
            >
              {page}
            </button>
          );
        })}

        <button
          className={classes.pageButton}
          disabled={current === totalPages}
          onClick={() => onChange?.(current + 1, pageSize)}
        >
          下一页
        </button>
      </div>
    );
  };

  return (
    <div className={`${classes.container} ${className || ''}`}>
      <div className={classes.tableWrapper}>
        <table className={classes.table}>
          {renderHeader()}
          {renderBody()}
        </table>
      </div>
      {renderPagination()}
    </div>
  );
};

export default Table;
