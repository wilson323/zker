// frontend/packages/arch/ui-components/src/components/Modal/Modal.tsx

import React, { useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import { useStyles } from './Modal.styles';

/**
 * Modal组件Props接口
 */
export interface ModalProps {
  /** 是否显示Modal */
  visible: boolean;
  /** 关闭回调 */
  onClose: () => void;
  /** 标题 */
  title?: string;
  /** 子元素 */
  children: React.ReactNode;
  /** 底部操作区 */
  footer?: React.ReactNode;
  /** 宽度 */
  width?: number | string;
  /** 点击遮罩层是否关闭 */
  maskClosable?: boolean;
  /** 是否显示确认按钮 */
  showOk?: boolean;
  /** 确认按钮文字 */
  okText?: string;
  /** 是否显示取消按钮 */
  showCancel?: boolean;
  /** 取消按钮文字 */
  cancelText?: string;
  /** 确认回调 */
  onOk?: () => void;
  /** 取消回调 */
  onCancel?: () => void;
  /** 确认按钮loading状态 */
  confirmLoading?: boolean;
  /** 自定义className */
  className?: string;
}

/**
 * Modal 对话框
 *
 * 模态对话框组件，用于需要用户处理事务，又不希望跳转页面时
 *
 * @example
 * ```tsx
 * <Modal
 *   visible={isVisible}
 *   onClose={handleClose}
 *   title="确认删除"
 *   onOk={handleOk}
 *   onCancel={handleCancel}
 * >
 *   确定要删除这条数据吗？
 * </Modal>
 * ```
 */
export const Modal: React.FC<ModalProps> = ({
  visible,
  onClose,
  title,
  children,
  footer,
  width = 520,
  maskClosable = true,
  showOk = true,
  okText = '确定',
  showCancel = true,
  cancelText = '取消',
  onOk,
  onCancel,
  confirmLoading = false,
  className,
}) => {
  const classes = useStyles({ width });
  const modalRef = useRef<HTMLDivElement>(null);

  // 处理ESC键关闭
  useEffect(() => {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && visible) {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEsc);
    return () => document.removeEventListener('keydown', handleEsc);
  }, [visible, onClose]);

  // 处理点击遮罩层关闭
  const handleMaskClick = (e: React.MouseEvent) => {
    if (maskClosable && e.target === e.currentTarget) {
      onClose();
    }
  };

  // 处理确认
  const handleOk = () => {
    onOk?.();
  };

  // 处理取消
  const handleCancel = () => {
    onCancel?.();
    onClose();
  };

  // 默认footer
  const defaultFooter = (
    <div className={classes.footer}>
      {showCancel && (
        <button
          className={`${classes.footerButton} ${classes.cancelButton}`}
          onClick={handleCancel}
        >
          {cancelText}
        </button>
      )}
      {showOk && (
        <button
          className={`${classes.footerButton} ${classes.okButton}`}
          onClick={handleOk}
          disabled={confirmLoading}
        >
          {confirmLoading && <span className={classes.spinner} />}
          {okText}
        </button>
      )}
    </div>
  );

  if (!visible) return null;

  return createPortal(
    <div className={`${classes.mask} ${className || ''}`} onClick={handleMaskClick}>
      <div ref={modalRef} className={classes.modal}>
        {title && (
          <div className={classes.header}>
            <h3 className={classes.title}>{title}</h3>
            <button
              className={classes.closeButton}
              onClick={onClose}
              aria-label="关闭"
            >
              ×
            </button>
          </div>
        )}
        <div className={classes.body}>{children}</div>
        {footer !== null && footer !== undefined ? footer : defaultFooter}
      </div>
    </div>,
    document.body
  );
};

export default Modal;
