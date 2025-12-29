// frontend/packages/arch/ui-components/src/index.ts

/**
 * ZKER 基础UI组件库
 *
 * 提供统一的UI组件，确保视觉一致性
 * 所有组件使用 @coze-studio/common/themes 设计令牌系统
 */

// 基础组件
export { Button } from './components/Button';
export type { ButtonProps } from './components/Button';

export { Input } from './components/Input';
export type { InputProps } from './components/Input';

export { Modal } from './components/Modal';
export type { ModalProps } from './components/Modal';

export { Table } from './components/Table';
export type { TableProps, Column, PaginationConfig } from './components/Table';

export { Card } from './components/Card';
export type { CardProps } from './components/Card';

// 表单组件
export { Form, FormItem } from './components/Form';
export type { FormProps, FormContextValue, FormItemProps } from './components/Form';

export { TextArea } from './components/TextArea';
export type { TextAreaProps } from './components/TextArea';

export { Checkbox } from './components/Checkbox';
export type { CheckboxProps } from './components/Checkbox';

export { Radio } from './components/Radio';
export type { RadioProps } from './components/Radio';

export { Select } from './components/Select';
export type { SelectProps, SelectOption } from './components/Select';

export { Slider } from './components/Slider';
export type { SliderProps } from './components/Slider';

export { Switch } from './components/Switch';
export type { SwitchProps } from './components/Switch';

// 数据展示组件
export { Badge } from './components/Badge';
export type { BadgeProps } from './components/Badge';

export { Tag } from './components/Tag';
export type { TagProps } from './components/Tag';

export { Progress } from './components/Progress';
export type { ProgressProps, ProgressType } from './components/Progress';

// 反馈组件
export { Spin } from './components/Spin';
export type { SpinProps } from './components/Spin';

export { Alert } from './components/Alert';
export type { AlertProps, AlertType } from './components/Alert';

export { Message } from './components/Message';
export type { MessageProps, MessageType } from './components/Message';

// 交互组件
export { Tooltip } from './components/Tooltip';
export type { TooltipProps, TooltipPlacement } from './components/Tooltip';

export { Dropdown } from './components/Dropdown';
export type { DropdownProps, DropdownOption } from './components/Dropdown';

export { DatePicker } from './components/DatePicker';
export type { DatePickerProps, DatePickerType } from './components/DatePicker';

export { Upload } from './components/Upload';
export type { UploadProps, UploadFile } from './components/Upload';
