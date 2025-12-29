# ZKER 技术组件清单与使用指南 (完整版)

> **文档类型**: 技术规范 / 组件使用手册
> **版本**: v1.0
> **生成日期**: 2025-01-01
> **覆盖范围**: 前端组件、后端服务、基础设施、第三方集成
> **使用目标**: 确保所有开发可以精准执行

---

## 📋 目录

1. [文档说明](#文档说明)
2. [前端组件清单与使用指南](#前端组件清单与使用指南)
3. [后端服务组件清单](#后端服务组件清单)
4. [基础设施组件清单](#基础设施组件清单)
5. [第三方服务集成清单](#第三方服务集成清单)
6. [组件开发规范](#组件开发规范)

---

## 文档说明

### 1.1 文档目标

**核心目标**: 提供一份完整的、可执行的技术组件清单和使用指南,确保:

✅ **前端开发**: 精准使用组件,避免重复造轮子
✅ **后端开发**: 精准调用服务,理解服务边界
✅ **测试人员**: 理解组件功能,编写完整测试
✅ **运维人员**: 理解依赖关系,正确配置服务
✅ **新人上手**: 快速找到所需组件,理解使用方式

### 1.2 文档结构

本文档分为五大部分:
1. **前端组件**: 所有可复用的 UI 组件和业务组件
2. **后端服务**: 所有微服务和领域服务
3. **基础设施**: 数据库、缓存、消息队列等
4. **第三方集成**: 外部服务集成
5. **开发规范**: 组件开发和使用规范

### 1.3 组件分类标准

**按层次分类**:
- **Level 1 (Arch)**: 架构层,基础工具和类型
- **Level 2 (Common/Foundation)**: 通用层,共享组件
- **Level 3 (Domain)**: 领域层,业务组件
- **Level 4 (App)**: 应用层,应用入口

**按类型分类**:
- **UI 组件**: 纯展示组件
- **业务组件**: 包含业务逻辑的组件
- **布局组件**: 页面布局组件
- **工具组件**: 工具类组件

**按来源分类**:
- **自研组件**: 团队内部开发的组件
- **第三方组件**: 从外部库引入的组件
- **封装组件**: 对第三方组件的封装

---

## 前端组件清单与使用指南

### 2.1 通用 UI 组件

#### 2.1.1 Button 按钮

**组件路径**: `@coze-studio/ui/button`
**源文件**: `frontend/packages/components/bot-semi/src/components/ui-button/index.tsx`

**功能说明**:
基础按钮组件,封装 Semi Design Button,统一样式和行为

**Props API**:

```typescript
interface UIButtonProps extends SemiButtonProps {
  // 继承 Semi Design 所有 Button 属性
  // 常用属性:

  /**
   * 按钮类型
   * @default 'primary'
   */
  type?: 'primary' | 'secondary' | 'tertiary' | 'warning' | 'danger';

  /**
   * 按钮尺寸
   * @default 'medium'
   */
  size?: 'small' | 'medium' | 'large';

  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean;

  /**
   * 是否加载中
   * @default false
   */
  loading?: boolean;

  /**
   * 点击事件
   */
  onClick?: (event: React.MouseEvent) => void;

  /**
   * 按钮内容
   */
  children?: React.ReactNode;

  /**
   * 图标
   */
  icon?: React.ReactNode;

  /**
   * 块级按钮 (占满父容器宽度)
   * @default false
   */
  block?: boolean;

  /**
   * 幽灵按钮 (透明背景)
   * @default false
   */
  ghost?: boolean;

  /**
   * HTML 类型
   * @default 'button'
   */
  htmlType?: 'button' | 'submit' | 'reset';
}
```

**使用示例**:

```typescript
import { Button } from '@coze-studio/ui';

// 1. 基础用法
<Button type="primary" onClick={handleClick}>
  点击我
</Button>

// 2. 不同类型
<Button type="primary">主要按钮</Button>
<Button type="secondary">次要按钮</Button>
<Button type="warning">警告按钮</Button>
<Button type="danger">危险按钮</Button>

// 3. 不同尺寸
<Button size="small">小按钮</Button>
<Button size="medium">中按钮</Button>
<Button size="large">大按钮</Button>

// 4. 带图标
<Button type="primary" icon={<IconPlus />}>
  新建
</Button>

// 5. 加载状态
<Button type="primary" loading={isLoading}>
  提交中
</Button>

// 6. 禁用状态
<Button type="primary" disabled>
  禁用按钮
</Button>

// 7. 块级按钮
<Button type="primary" block>
  占满宽度
</Button>

// 8. 幽灵按钮
<Button type="primary" ghost>
  幽灵按钮
</Button>

// 9. 提交表单
<Button type="primary" htmlType="submit">
  提交表单
</Button>
```

**最佳实践**:

✅ **推荐做法**:
- 明确指定 `type` 属性,不要依赖默认值
- 使用 `loading` 状态提供更好的用户体验
- 为危险操作使用 `type="danger"`
- 为表单提交指定 `htmlType="submit"`

❌ **不推荐做法**:
- 不要嵌套按钮
- 不要在没有 `onClick` 的情况下使用按钮 (非提交场景)
- 不要过度使用 `type="primary"`,保持视觉层次

**注意事项**:
- 按钮文本应简洁明了,不超过 4-6 个字
- 同一页面中主要操作按钮不应超过 2 个
- 避免在按钮中使用过长文本

**常见问题**:

**Q: 按钮点击没有反应?**
A: 检查是否被禁用或有父元素阻止事件冒泡

**Q: 如何实现确认后再执行?**
A: 使用 Modal.confirm 或自定义确认逻辑

```typescript
const handleClick = async () => {
  const confirmed = await Modal.confirm({
    title: '确认删除?',
    content: '删除后无法恢复',
  });
  if (confirmed) {
    await handleDelete();
  }
};

<Button type="danger" onClick={handleClick}>
  删除
</Button>
```

---

#### 2.1.2 Input 输入框

**组件路径**: `@coze-studio/ui/input`
**源文件**: `frontend/packages/components/bot-semi/src/components/ui-input/index.tsx`

**功能说明**:
基础输入框组件,封装 Semi Design Input

**Props API**:

```typescript
interface UIInputProps extends InputProps {
  /**
   * 输入框类型
   * @default 'text'
   */
  type?: 'text' | 'password' | 'number' | 'email' | 'tel' | 'url';

  /**
   * 占位符文本
   */
  placeholder?: string;

  /**
   * 输入框尺寸
   * @default 'medium'
   */
  size?: 'small' | 'medium' | 'large';

  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean;

  /**
   * 值变化回调
   */
  onChange?: (value: string) => void;

  /**
   * 失去焦点回调
   */
  onBlur?: (value: string) => void;

  /**
   * 获得焦点回调
   */
  onFocus?: () => void;

  /**
   * 按下回车回调
   */
  onPressEnter?: () => void;

  /**
   * 最大长度
   */
  maxLength?: number;

  /**
   * 是否显示字数统计
   * @default false
   */
  showLengthLimit?: boolean;

  /**
   * 前缀图标
   */
  prefix?: React.ReactNode;

  /**
   * 后缀图标
   */
  suffix?: React.ReactNode;

  /**
   * 是否允许清除
   * @default false
   */
  allowClear?: boolean;

  /**
   * 输入框值 (受控)
   */
  value?: string;

  /**
   * 输入框默认值 (非受控)
   */
  defaultValue?: string;

  /**
   * 验证状态
   */
  validateStatus?: 'error' | 'warning';

  /**
   * 帮助文本
   */
  helpText?: string;
}
```

**使用示例**:

```typescript
import { Input } from '@coze-studio/ui';

// 1. 基础用法
<Input
  placeholder="请输入用户名"
  value={username}
  onChange={setUsername}
/>

// 2. 不同类型
<Input type="text" placeholder="文本输入" />
<Input type="password" placeholder="密码输入" />
<Input type="number" placeholder="数字输入" />
<Input type="email" placeholder="邮箱输入" />

// 3. 带前后缀
<Input
  prefix={<IconSearch />}
  placeholder="搜索..."
/>
<Input
  suffix={<IconUser />}
  placeholder="用户名"
/>

// 4. 限制长度
<Input
  maxLength={20}
  showLengthLimit
  placeholder="最多输入20个字符"
/>

// 5. 可清除
<Input
  allowClear
  placeholder="可清除输入"
  value={value}
  onChange={setValue}
/>

// 6. 验证状态
<Input
  validateStatus="error"
  helpText="用户名不能为空"
  value={username}
  onChange={setUsername}
/>

// 7. 受控与非受控
// 受控组件
<Input
  value={controlledValue}
  onChange={setControlledValue}
/>

// 非受控组件
<Input
  defaultValue="默认值"
  onBlur={handleBlur}
/>

// 8. 回车触发
<Input
  placeholder="输入后按回车搜索"
  onPressEnter={handleSearch}
/>

// 9. 数字输入限制
<Input
  type="number"
  placeholder="请输入年龄"
  onChange={(value) => {
    const num = parseInt(value);
    if (!isNaN(num) && num >= 0 && num <= 150) {
      setAge(num);
    }
  }}
/>
```

**最佳实践**:

✅ **推荐做法**:
- 表单场景优先使用受控组件
- 为敏感信息使用 `type="password"`
- 为有限制的输入添加 `maxLength`
- 为需要即时反馈的场景使用 `validateStatus`

❌ **不推荐做法**:
- 不要同时使用 `value` 和 `defaultValue`
- 不要在 `onChange` 中进行复杂计算或异步操作
- 不要过度使用 `suffix`,可能影响可读性

**注意事项**:
- 受控组件必须处理 `onChange` 事件
- 数字类型输入需要自行验证范围
- `maxLength` 只限制字符数,不限制字节

**常见问题**:

**Q: 输入框输入后值不变?**
A: 检查是否正确处理了 `onChange` 事件

```typescript
// ❌ 错误:没有更新状态
<Input value={value} />

// ✅ 正确
<Input value={value} onChange={setValue} />
```

**Q: 如何实现防抖搜索?**
A: 使用 `useDebounce` Hook

```typescript
import { useDebounce } from '@coze-arch/hooks';

const [keyword, setKeyword] = useState('');
const debouncedKeyword = useDebounce(keyword, 300);

useEffect(() => {
  if (debouncedKeyword) {
    performSearch(debouncedKeyword);
  }
}, [debouncedKeyword]);

<Input
  value={keyword}
  onChange={setKeyword}
  placeholder="搜索..."
/>
```

---

#### 2.1.3 Select 选择器

**组件路径**: `@coze-studio/ui/select`
**源文件**: `frontend/packages/components/bot-semi/src/components/ui-select/index.tsx`

**功能说明**:
下拉选择器组件,支持单选、多选、搜索、分组

**Props API**:

```typescript
interface UISelectProps {
  /**
   * 选择器数据源
   */
  options: Array<{
    label: string;
    value: string | number;
    disabled?: boolean;
    children?: Array<any>; // 分组子选项
  }>;

  /**
   * 占位符文本
   * @default '请选择'
   */
  placeholder?: string;

  /**
   * 选择器尺寸
   * @default 'medium'
   */
  size?: 'small' | 'medium' | 'large';

  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean;

  /**
   * 是否多选
   * @default false
   */
  multiple?: boolean;

  /**
   * 是否可搜索
   * @default false
   */
  filter?: boolean;

  /**
   * 是否允许清除
   * @default false
   */
  allowClear?: boolean;

  /**
   * 选中值
   */
  value?: string | number | Array<string | number>;

  /**
   * 默认选中值
   */
  defaultValue?: string | number | Array<string | number>;

  /**
   * 值变化回调
   */
  onChange?: (value: any) => void;

  /**
   * 选项渲染函数
   */
  optionRender?: (option: any) => React.ReactNode;

  /**
   * 最大显示标签数 (多选时)
   * @default undefined
   */
  maxTagCount?: number;

  /**
   * 是否全选 (多选时)
   * @default false
   */
  showSelectAll?: boolean;

  /**
   * 加载中
   * @default false
   */
  loading?: boolean;
}
```

**使用示例**:

```typescript
import { Select } from '@coze-studio/ui';

// 1. 基础单选
const options = [
  { label: '选项1', value: '1' },
  { label: '选项2', value: '2' },
  { label: '选项3', value: '3' },
];

<Select
  options={options}
  placeholder="请选择"
  value={value}
  onChange={setValue}
/>

// 2. 多选
<Select
  options={options}
  multiple
  placeholder="请选择多个"
  value={values}
  onChange={setValues}
/>

// 3. 可搜索
<Select
  options={options}
  filter
  placeholder="搜索选择"
  value={value}
  onChange={setValue}
/>

// 4. 分组
const groupedOptions = [
  {
    label: '分组1',
    children: [
      { label: '选项1-1', value: '1-1' },
      { label: '选项1-2', value: '1-2' },
    ],
  },
  {
    label: '分组2',
    children: [
      { label: '选项2-1', value: '2-1' },
      { label: '选项2-2', value: '2-2' },
    ],
  },
];

<Select
  options={groupedOptions}
  placeholder="请选择"
/>

// 5. 可清除
<Select
  options={options}
  allowClear
  placeholder="可清除选择"
  value={value}
  onChange={setValue}
/>

// 6. 异步加载
const [loading, setLoading] = useState(false);
const [asyncOptions, setAsyncOptions] = useState([]);

const handleSearch = async (keyword) => {
  setLoading(true);
  const data = await fetchOptions(keyword);
  setAsyncOptions(data);
  setLoading(false);
};

<Select
  options={asyncOptions}
  filter
  loading={loading}
  onSearch={handleSearch}
  placeholder="搜索加载"
/>

// 7. 自定义选项渲染
<Select
  options={userOptions}
  optionRender={(option) => (
    <div className="flex items-center">
      <Avatar src={option.avatar} size="small" />
      <span className="ml-2">{option.label}</span>
    </div>
  )}
  placeholder="选择用户"
/>

// 8. 多选限制
<Select
  options={options}
  multiple
  maxTagCount={2}
  placeholder="最多显示2个标签"
  value={values}
  onChange={setValues}
/>

// 9. 全选功能
<Select
  options={options}
  multiple
  showSelectAll
  placeholder="全选或多选"
  value={values}
  onChange={setValues}
/>

// 10. 禁用选项
const optionsWithDisabled = [
  { label: '选项1', value: '1' },
  { label: '选项2', value: '2', disabled: true },
  { label: '选项3', value: '3' },
];

<Select
  options={optionsWithDisabled}
  placeholder="部分选项禁用"
/>
```

**最佳实践**:

✅ **推荐做法**:
- 为选项提供有意义的 `label` 和 `value`
- 数据量大时使用 `filter` 可搜索
- 多选时使用 `maxTagCount` 避免占用过多空间
- 异步加载时显示 `loading` 状态

❌ **不推荐做法**:
- 不要在 `options` 中使用 `undefined` 作为 `value`
- 不要在 `onChange` 中进行复杂异步操作
- 不要过度使用分组,保持选项清晰

**注意事项**:
- `value` 必须与 `options` 中的某个 `value` 匹配
- 多选时 `value` 为数组
- 分组时 `children` 不能为空数组

**常见问题**:

**Q: 选择后显示的值不正确?**
A: 检查 `value` 是否与 `options` 中的 `value` 类型一致

```typescript
// ❌ 错误: 类型不一致
const value = '1'; // string
const options = [{ label: '选项', value: 1 }]; // number

// ✅ 正确
const value = 1; // number
const options = [{ label: '选项', value: 1 }]; // number
```

**Q: 如何实现远程搜索?**
A: 结合 `filter` 和 `onSearch`

```typescript
const [options, setOptions] = useState([]);
const [loading, setLoading] = useState(false);

const handleSearch = async (keyword) => {
  if (!keyword) {
    setOptions([]);
    return;
  }

  setLoading(true);
  try {
    const data = await api.search(keyword);
    setOptions(data.map(item => ({
      label: item.name,
      value: item.id,
    })));
  } finally {
    setLoading(false);
  }
};

<Select
  options={options}
  filter
  loading={loading}
  onSearch={handleSearch}
  placeholder="输入搜索关键词"
/>
```

---

#### 2.1.4 Modal 模态框

**组件路径**: `@coze-studio/ui/modal`
**源文件**: `frontend/packages/components/bot-semi/src/components/ui-modal/ui-modal.tsx`

**功能说明**:
模态对话框组件,支持确认框、表单对话框、自定义内容

**Props API**:

```typescript
interface UIModalProps {
  /**
   * 是否显示
   * @default false
   */
  visible: boolean;

  /**
   * 标题
   */
  title?: React.ReactNode;

  /**
   * 内容
   */
  children?: React.ReactNode;

  /**
   * 宽度
   * @default 520
   */
  width?: number | string;

  /**
   * 是否显示关闭图标
   * @default true
   */
  closable?: boolean;

  /**
   * 是否显示遮罩
   * @default true
   */
  maskClosable?: boolean;

  /**
   * 点击确定回调
   */
  onOk?: () => void | Promise<void>;

  /**
   * 点击取消回调
   */
  onCancel?: () => void;

  /**
   * 确定按钮文字
   * @default '确定'
   */
  okText?: string;

  /**
   * 取消按钮文字
   * @default '取消'
   */
  cancelText?: string;

  /**
   * 确定按钮加载状态
   * @default false
   */
  okLoading?: boolean;

  /**
   * 确定按钮禁用状态
   * @default false
   */
  okDisabled?: boolean;

  /**
   * 是否显示底部按钮
   * @default true
   */
  footer?: React.ReactNode | null;

  /**
   * 关闭后回调
   */
  afterClose?: () => void;

  /**
   * 居中位置
   * @default false
   */
  centered?: boolean;

  /**
   * 方向
   * @default 'ltr'
   */
  direction?: 'ltr' | 'rtl';
}
```

**使用示例**:

```typescript
import { Modal } from '@coze-studio/ui';

// 1. 基础用法
const [visible, setVisible] = useState(false);

<>
  <Button onClick={() => setVisible(true)}>打开对话框</Button>
  <Modal
    visible={visible}
    title="标题"
    onOk={() => setVisible(false)}
    onCancel={() => setVisible(false)}
  >
    <p>对话框内容</p>
  </Modal>
</>

// 2. 确认框
const handleConfirm = async () => {
  const confirmed = await Modal.confirm({
    title: '确认删除?',
    content: '删除后无法恢复,是否继续?',
    okText: '确定删除',
    cancelText: '取消',
    okType: 'danger',
  });

  if (confirmed) {
    await deleteItem();
  }
};

// 3. 异步操作
const [visible, setVisible] = useState(false);
const [loading, setLoading] = useState(false);

const handleOk = async () => {
  setLoading(true);
  try {
    await submitForm();
    setVisible(false);
  } finally {
    setLoading(false);
  }
};

<Modal
  visible={visible}
  title="提交表单"
  onOk={handleOk}
  onCancel={() => setVisible(false)}
  okLoading={loading}
>
  <Form>...</Form>
</Modal>

// 4. 自定义底部
<Modal
  visible={visible}
  title="自定义底部"
  footer={[
    <Button key="back" onClick={() => setVisible(false)}>
      返回
    </Button>,
    <Button key="submit" type="primary" onClick={handleSubmit}>
      提交
    </Button>,
  ]}
>
  <p>内容</p>
</Modal>

// 5. 无底部按钮
<Modal
  visible={visible}
  title="无底部"
  footer={null}
  onCancel={() => setVisible(false)}
>
  <p>只有关闭按钮</p>
</Modal>

// 6. 宽度调整
<Modal
  visible={visible}
  title="大尺寸对话框"
  width={800}
>
  <p>更宽的内容区域</p>
</Modal>

// 7. 居中显示
<Modal
  visible={visible}
  title="居中对话框"
  centered
>
  <p>垂直居中显示</p>
</Modal>

// 8. 点击遮罩不关闭
<Modal
  visible={visible}
  title="强制操作"
  maskClosable={false}
  onOk={handleOk}
  onCancel={handleCancel}
>
  <p>必须点击确定或取消才能关闭</p>
</Modal>

// 9. 信息提示框
Modal.info({
  title: '提示信息',
  content: '这是一条提示信息',
});

Modal.success({
  title: '操作成功',
  content: '您的操作已成功完成',
});

Modal.warning({
  title: '警告',
  content: '请注意可能存在的问题',
});

Modal.error({
  title: '错误',
  content: '操作失败,请稍后重试',
});

// 10. 销毁确认
const handleDelete = async () => {
  Modal.confirm({
    title: '确认删除?',
    content: '删除后无法恢复',
    okText: '确定删除',
    okType: 'danger',
    onOk: async () => {
      await deleteItem();
      message.success('删除成功');
    },
  });
};
```

**最佳实践**:

✅ **推荐做法**:
- 为危险操作使用 `Modal.confirm` 并指定 `okType="danger"`
- 异步操作时显示 `okLoading`
- 重要操作设置 `maskClosable=false` 防止误操作
- 关闭后清理状态使用 `afterClose`

❌ **不推荐做法**:
- 不要在 Modal 中嵌套另一个 Modal
- 不要过度使用 Modal,优先使用页面内提示
- 不要在 `onOk` 中抛出异常,应该捕获处理

**注意事项**:
- `visible` 控制显示隐藏
- 关闭操作包括: 点击取消、点击遮罩、点击关闭图标
- `afterClose` 在动画结束后触发,适合清理状态

**常见问题**:

**Q: Modal 关闭后状态没有清理?**
A: 使用 `afterClose` 清理状态

```typescript
<Modal
  visible={visible}
  afterClose={() => {
    setFormData(initialData);
    setErrors({});
  }}
>
  <Form data={formData} />
</Modal>
```

**Q: 如何在表单验证失败后阻止关闭?**
A: 控制 `okDisabled` 或在 `onOk` 中返回

```typescript
const handleOk = async () => {
  const valid = await form.validate();
  if (!valid) {
    return; // 阻止关闭
  }

  await submitForm();
  setVisible(false);
};
```

---

#### 2.1.5 Table 表格

**组件路径**: `@coze-studio/ui/table`
**源文件**: `frontend/packages/components/bot-semi/src/components/ui-table/index.tsx`

**功能说明**:
表格组件,支持排序、筛选、展开、选择等功能

**Props API**:

```typescript
interface UITableProps {
  /**
   * 列配置
   */
  columns: Array<{
    title: string; // 列标题
    dataIndex?: string; // 数据字段名
    key?: string; // 唯一键
    width?: number | string; // 列宽
    align?: 'left' | 'center' | 'right'; // 对齐方式
    fixed?: 'left' | 'right'; // 固定列
    sorter?: boolean | Function; // 排序
    filters?: Array<{ text: string; value: any }>; // 筛选
    render?: (value: any, record: any, index: number) => React.ReactNode; // 自定义渲染
  }>;

  /**
   * 数据源
   */
  dataSource?: Array<any>;

  /**
   * 行键
   */
  rowKey?: string | ((record: any) => string);

  /**
   * 是否加载中
   * @default false
   */
  loading?: boolean;

  /**
   * 是否显示边框
   * @default false
   */
  bordered?: boolean;

  /**
   * 表格大小
   */
  size?: 'small' | 'middle' | 'large';

  /**
   * 行选择配置
   */
  rowSelection?: {
    selectedRowKeys?: Array<any>;
    onChange?: (selectedRowKeys: Array<any>, selectedRows: Array<any>) => void;
    type?: 'checkbox' | 'radio';
  };

  /**
   * 分页配置
   */
  pagination?: {
    current?: number;
    pageSize?: number;
    total?: number;
    showSizeChanger?: boolean;
    showQuickJumper?: boolean;
    onChange?: (page: number, pageSize: number) => void;
  } | false;

  /**
   * 可滚动配置
   */
  scroll?: {
    x?: number | string;
    y?: number | string;
  };

  /**
   * 空数据时显示
   */
  emptyText?: React.ReactNode;

  /**
   * 行点击事件
   */
  onRow?: (record: any, index: number) => {
    onClick?: (event: React.MouseEvent) => void;
    onDoubleClick?: (event: React.MouseEvent) => void;
  };
}
```

**使用示例**:

```typescript
import { Table } from '@coze-studio/ui';

// 1. 基础用法
const columns = [
  {
    title: '姓名',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '年龄',
    dataIndex: 'age',
    key: 'age',
  },
  {
    title: '地址',
    dataIndex: 'address',
    key: 'address',
  },
];

const dataSource = [
  { key: '1', name: '张三', age: 32, address: '北京市' },
  { key: '2', name: '李四', age: 28, address: '上海市' },
];

<Table
  columns={columns}
  dataSource={dataSource}
  rowKey="key"
/>

// 2. 自定义渲染
const columns = [
  {
    title: '头像',
    dataIndex: 'avatar',
    render: (avatar) => <Avatar src={avatar} />,
  },
  {
    title: '姓名',
    dataIndex: 'name',
  },
  {
    title: '状态',
    dataIndex: 'status',
    render: (status) => {
      const color = status === 'active' ? 'green' : 'red';
      return <Tag color={color}>{status}</Tag>;
    },
  },
  {
    title: '操作',
    render: (_, record) => (
      <Space>
        <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
        <Button type="link" danger onClick={() => handleDelete(record)}>删除</Button>
      </Space>
    ),
  },
];

// 3. 排序
const columns = [
  {
    title: '姓名',
    dataIndex: 'name',
    sorter: (a, b) => a.name.localeCompare(b.name),
  },
  {
    title: '年龄',
    dataIndex: 'age',
    sorter: (a, b) => a.age - b.age,
  },
];

// 4. 行选择
const [selectedRowKeys, setSelectedRowKeys] = useState([]);

const rowSelection = {
  selectedRowKeys,
  onChange: (selectedKeys, selectedRows) => {
    setSelectedRowKeys(selectedKeys);
    console.log('选中的行:', selectedRows);
  },
};

<Table
  columns={columns}
  dataSource={dataSource}
  rowSelection={rowSelection}
/>

// 5. 分页
const [pagination, setPagination] = useState({
  current: 1,
  pageSize: 10,
  total: 100,
});

<Table
  columns={columns}
  dataSource={dataSource}
  pagination={{
    ...pagination,
    showSizeChanger: true,
    showQuickJumper: true,
    onChange: (page, pageSize) => {
      setPagination({ ...pagination, current: page, pageSize });
      fetchData(page, pageSize);
    },
  }}
/>

// 6. 加载状态
<Table
  columns={columns}
  dataSource={dataSource}
  loading={isLoading}
/>

// 7. 可滚动
<Table
  columns={columns}
  dataSource={dataSource}
  scroll={{ x: 1500, y: 400 }}
/>

// 8. 固定列
const columns = [
  {
    title: '姓名',
    dataIndex: 'name',
    fixed: 'left',
    width: 100,
  },
  // ... 其他列
  {
    title: '操作',
    fixed: 'right',
    width: 150,
    render: (_, record) => <Button>操作</Button>,
  },
];

// 9. 空数据
<Table
  columns={columns}
  dataSource={[]}
  emptyText={<Empty description="暂无数据" />}
/>

// 10. 行点击
const onRow = (record) => ({
  onClick: () => {
    console.log('点击了行:', record);
  },
  onDoubleClick: () => {
    console.log('双击了行:', record);
  },
});

<Table
  columns={columns}
  dataSource={dataSource}
  onRow={onRow}
/>
```

**最佳实践**:

✅ **推荐做法**:
- 为每列指定唯一的 `key`
- 使用 `rowKey` 指定唯一标识字段
- 大数据量时使用分页或虚拟滚动
- 为操作列固定在右侧

❌ **不推荐做法**:
- 不要在 `render` 函数中进行复杂计算或异步操作
- 不要在表格单元格中嵌套表格
- 不要使用 `index` 作为 `rowKey`

**注意事项**:
- `rowKey` 必须是唯一值,推荐使用 ID
- 自定义渲染时使用 `render` 而非 `dataIndex`
- 固定列需要设置 `width`

**常见问题**:

**Q: 警告 "Each child in a list should have a unique \"key\" prop"?**
A: 检查是否设置了 `rowKey`

```typescript
// ❌ 错误:没有设置 rowKey
<Table columns={columns} dataSource={dataSource} />

// ✅ 正确
<Table
  columns={columns}
  dataSource={dataSource}
  rowKey="id"
/>
```

**Q: 如何实现表格的刷新?**
A: 控制 `dataSource` 或调用刷新函数

```typescript
const [dataSource, setDataSource] = useState([]);
const [loading, setLoading] = useState(false);

const fetchData = async () => {
  setLoading(true);
  try {
    const data = await api.getUsers();
    setDataSource(data);
  } finally {
    setLoading(false);
  }
};

useEffect(() => {
  fetchData();
}, []);

<Table
  columns={columns}
  dataSource={dataSource}
  loading={loading}
/>
```

---

### 2.2 业务组件

#### 2.2.1 ChatInput 聊天输入框

**组件路径**: `@coze-studio/chat/chat-input`
**源文件**: `frontend/packages/common/chat-area/chat-uikit/src/components/chat/chat-input/index.tsx`

**功能说明**:
聊天输入组件,支持文本、图片、文件上传,支持语音输入

**Props API**:

```typescript
interface ChatInputProps {
  /**
   * 输入框占位符
   * @default '输入消息...'
   */
  placeholder?: string;

  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean;

  /**
   * 发送消息回调
   */
  onSend: (content: string, attachments?: Attachment[]) => void | Promise<void>;

  /**
   * 上传文件回调
   */
  onUpload?: (files: File[]) => void | Promise<void[]>;

  /**
   * 最大文件大小 (MB)
   * @default 10
   */
  maxFileSize?: number;

  /**
   * 允许的文件类型
   */
  acceptFileTypes?: string[];

  /**
   * 是否支持图片上传
   * @default true
   */
  supportImage?: boolean;

  /**
   * 是否支持文件上传
   * @default true
   */
  supportFile?: boolean;

  /**
   * 是否支持语音输入
   * @default true
   */
  supportVoice?: boolean;

  /**
   * 输入框最小行数
   * @default 1
   */
  minRows?: number;

  /**
   * 输入框最大行数
   * @default 6
   */
  maxRows?: number;

  /**
   * 快捷操作按钮
   */
  actions?: React.ReactNode[];

  /**
   * 发送中状态
   * @default false
   */
  sending?: boolean;
}
```

**使用示例**:

```typescript
import { ChatInput } from '@coze-studio/chat';

// 1. 基础用法
const [sending, setSending] = useState(false);

const handleSend = async (content) => {
  setSending(true);
  try {
    await api.sendMessage({ content });
  } finally {
    setSending(false);
  }
};

<ChatInput onSend={handleSend} sending={sending} />

// 2. 带附件上传
const handleUpload = async (files) => {
  const attachments = await Promise.all(
    files.map(async (file) => {
      const url = await api.uploadFile(file);
      return {
        type: file.type.startsWith('image/') ? 'image' : 'file',
        url,
        name: file.name,
      };
    })
  );
  return attachments;
};

<ChatInput
  onSend={handleSend}
  onUpload={handleUpload}
  supportImage
  supportFile
/>

// 3. 限制文件类型和大小
<ChatInput
  onSend={handleSend}
  acceptFileTypes={['image/png', 'image/jpeg', 'application/pdf']}
  maxFileSize={5} // 5MB
/>

// 4. 自定义操作按钮
const actions = [
  <Button icon={<IconPlus />} onClick={handleCustomAction}>
    自定义
  </Button>,
  <Button icon={<IconEmoji />} onClick={handleEmoji}>
    表情
  </Button>,
];

<ChatInput
  onSend={handleSend}
  actions={actions}
/>

// 5. 语音输入
<ChatInput
  onSend={handleSend}
  supportVoice
/>

// 6. 多行输入
<ChatInput
  onSend={handleSend}
  minRows={2}
  maxRows={10}
  placeholder="输入您的消息,支持换行..."
/>

// 7. 禁用状态
<ChatInput
  onSend={handleSend}
  disabled={!isConnected}
  placeholder="连接后可发送消息..."
/>

// 8. 快捷回复
const quickReplies = [
  { text: '你好', value: '你好' },
  { text: '谢谢', value: '谢谢' },
];

<ChatInput
  onSend={handleSend}
  quickReplies={quickReplies}
/>
```

**最佳实践**:

✅ **推荐做法**:
- 发送消息时显示 `sending` 状态防止重复提交
- 限制文件大小和类型,避免上传过大或不支持的文件
- 为语音输入提供清晰的引导
- 发送失败时恢复输入框内容

❌ **不推荐做法**:
- 不要在 `onSend` 中抛出未捕获的异常
- 不要同时显示多个快捷操作按钮,保持简洁
- 不要移除发送按钮,始终提供明确的发送入口

**注意事项**:
- 发送时清空输入框内容
- 上传的文件需要先转换服务器 URL
- 语音输入需要浏览器支持相应 API

**常见问题**:

**Q: 发送后输入框没有清空?**
A: 检查是否正确处理了 `onSend` 返回值

```typescript
const handleSend = async (content) => {
  const success = await api.sendMessage({ content });
  return success; // 返回 true 会清空输入框
};
```

---

#### 2.2.2 MessageBox 消息框

**组件路径**: `@coze-studio/chat/message-box`
**源文件**: `frontend/packages/common/chat-area/chat-uikit/src/components/common/message-box/index.tsx`

**功能说明**:
消息展示组件,支持文本、图片、文件、函数调用等多种消息类型

**Props API**:

```typescript
interface MessageBoxProps {
  /**
   * 消息数据
   */
  message: {
    id: string;
    role: 'user' | 'assistant' | 'system';
    content: string;
    contentType: 'text' | 'image' | 'file' | 'multimodal' | 'function_call';
    attachments?: Attachment[];
    createdAt: string;
    metadata?: Record<string, any>;
  };

  /**
   * 消息主题
   * @default 'light'
   */
  theme?: 'light' | 'dark';

  /**
   * 是否显示头像
   * @default true
   */
  showAvatar?: boolean;

  /**
   * 是否显示时间
   * @default true
   */
  showTime?: boolean;

  /**
   * 是否可复制
   * @default true
   */
  copyable?: boolean;

  /**
   * 头像 URL
   */
  avatarUrl?: string;

  /**
   * 用户名
   */
  username?: string;

  /**
   * 操作按钮
   */
  actions?: React.ReactNode[];

  /**
   * 点击事件
   */
  onClick?: () => void;
}
```

**使用示例**:

```typescript
import { MessageBox } from '@coze-studio/chat';

// 1. 基础用法 - 文本消息
const message = {
  id: '1',
  role: 'assistant',
  content: '您好,有什么可以帮助您的吗?',
  contentType: 'text',
  createdAt: '2025-01-01T12:00:00Z',
};

<MessageBox message={message} />

// 2. 用户消息
const userMessage = {
  id: '2',
  role: 'user',
  content: '帮我分析一下这个数据',
  contentType: 'text',
  createdAt: '2025-01-01T12:01:00Z',
};

<MessageBox
  message={userMessage}
  avatarUrl="/avatars/user.png"
  username="张三"
/>

// 3. 图片消息
const imageMessage = {
  id: '3',
  role: 'user',
  content: '![图片](/uploads/image.png)',
  contentType: 'image',
  attachments: [{
    type: 'image',
    url: '/uploads/image.png',
  }],
  createdAt: '2025-01-01T12:02:00Z',
};

<MessageBox message={imageMessage} />

// 4. 文件消息
const fileMessage = {
  id: '4',
  role: 'assistant',
  content: '这是您请求的文档',
  contentType: 'file',
  attachments: [{
    type: 'file',
    url: '/uploads/document.pdf',
    name: 'document.pdf',
    size: 1024000,
  }],
  createdAt: '2025-01-01T12:03:00Z',
};

<MessageBox message={fileMessage} />

// 5. 多模态消息
const multimodalMessage = {
  id: '5',
  role: 'assistant',
  content: '这是您要的图片和分析',
  contentType: 'multimodal',
  attachments: [
    { type: 'image', url: '/uploads/chart.png' },
    { type: 'text', content: '数据分析结果...' },
  ],
  createdAt: '2025-01-01T12:04:00Z',
};

<MessageBox message={multimodalMessage} />

// 6. 函数调用消息
const functionCallMessage = {
  id: '6',
  role: 'assistant',
  content: '正在为您查询天气...',
  contentType: 'function_call',
  metadata: {
    functionName: 'getWeather',
    arguments: { city: '北京' },
    result: { temperature: 25, weather: '晴' },
  },
  createdAt: '2025-01-01T12:05:00Z',
};

<MessageBox message={functionCallMessage} />

// 7. 自定义操作按钮
const actions = [
  <Button icon={<IconCopy />} size="small">
    复制
  </Button>,
  <Button icon={<IconRegenerate />} size="small">
    重新生成
  </Button>,
  <Button icon={<IconThumbUp />} size="small">
    点赞
  </Button>,
];

<MessageBox
  message={message}
  actions={actions}
/>

// 8. 深色主题
<MessageBox
  message={message}
  theme="dark"
/>

// 9. 点击事件
const handleClick = () => {
  console.log('点击了消息:', message);
};

<MessageBox
  message={message}
  onClick={handleClick}
/>

// 10. 自定义渲染
const CustomMessageBox = (props) => {
  return (
    <MessageBox {...props}>
      {(message) => (
        <div className="custom-content">
          {message.contentType === 'text' && (
            <Markdown>{message.content}</Markdown>
          )}
          {message.contentType === 'image' && (
            <Image src={message.attachments[0].url} />
          )}
        </div>
      )}
    </MessageBox>
  );
};
```

**最佳实践**:

✅ **推荐做法**:
- 为不同角色设置不同的头像和样式
- 长消息提供"展开/收起"功能
- 为函数调用消息展示清晰的调用过程和结果
- 提供复制、重新生成等操作按钮

❌ **不推荐做法**:
- 不要在消息中执行自动播放音频或视频
- 不要在消息中嵌入复杂的交互组件
- 不要过度使用自定义渲染,保持一致性

**注意事项**:
- Markdown 内容需要安全渲染,避免 XSS
- 图片和文件需要检查权限
- 时间显示应考虑本地化

**常见问题**:

**Q: Markdown 渲染不正确?**
A: 检查是否正确配置了 Markdown 渲染器

```typescript
import { Markdown } from '@coze-studio/ui';

<MessageBox message={message}>
  {(msg) => <Markdown>{msg.content}</Markdown>}
</MessageBox>
```

**Q: 如何实现消息流式输出?**
A: 使用特殊的流式消息组件或 Hook

```typescript
import { useStreamMessage } from '@coze-arch/hooks';

const StreamingMessage = () => {
  const { content, isComplete } = useStreamMessage(messageId);

  return (
    <MessageBox
      message={{
        ...message,
        content,
        metadata: { streaming: !isComplete },
      }}
    />
  );
};
```

---

#### 2.2.3 BotSelector Bot 选择器

**组件路径**: `@coze-studio/bot/bot-selector`
**源文件**: `frontend/packages/workflow/playground/src/components/bot-select/index.tsx`

**功能说明**:
Bot 选择器组件,支持搜索、分页、多选

**Props API**:

```typescript
interface BotSelectorProps {
  /**
   * 选中的 Bot ID
   */
  value?: string | string[];

  /**
   * 默认选中的 Bot ID
   */
  defaultValue?: string | string[];

  /**
   * 值变化回调
   */
  onChange?: (value: string | string[]) => void;

  /**
   * 是否多选
   * @default false
   */
  multiple?: boolean;

  /**
   * 是否可搜索
   * @default true
   */
  searchable?: boolean;

  /**
   * 占位符
   * @default '请选择 Bot'
   */
  placeholder?: string;

  /**
   * 是否禁用
   * @default false
   */
  disabled?: boolean;

  /**
   * Bot 筛选条件
   */
  filter?: {
    status?: string[];
    type?: string[];
    tags?: string[];
  };

  /**
   * 自定义 Bot 渲染
   */
  renderBot?: (bot: Bot) => React.ReactNode;

  /**
   * 加载 Bot 数据的函数
   */
  loadBots?: (params: any) => Promise<Bot[]>;
}
```

**使用示例**:

```typescript
import { BotSelector } from '@coze-studio/bot';

// 1. 基础用法 - 单选
const [botId, setBotId] = useState('');

<BotSelector
  value={botId}
  onChange={setBotId}
  placeholder="请选择一个 Bot"
/>

// 2. 多选
const [botIds, setBotIds] = useState([]);

<BotSelector
  multiple
  value={botIds}
  onChange={setBotIds}
  placeholder="请选择多个 Bot"
/>

// 3. 禁用搜索
<BotSelector
  searchable={false}
  value={botId}
  onChange={setBotId}
/>

// 4. 筛选条件
<BotSelector
  value={botId}
  onChange={setBotId}
  filter={{
    status: ['published'],
    type: ['chatbot'],
  }}
/>

// 5. 自定义渲染
<BotSelector
  value={botId}
  onChange={setBotId}
  renderBot={(bot) => (
    <div className="flex items-center">
      <Avatar src={bot.avatar} size="small" />
      <div className="ml-2">
        <div className="font-medium">{bot.name}</div>
        <div className="text-xs text-gray-500">{bot.description}</div>
      </div>
    </div>
  )}
/>

// 6. 自定义数据加载
const loadBots = async (params) => {
  const response = await api.bots.list({
    page: params.page,
    pageSize: params.pageSize,
    keyword: params.keyword,
  });
  return response.data;
};

<BotSelector
  value={botId}
  onChange={setBotId}
  loadBots={loadBots}
/>

// 7. 受控与非受控
// 受控
<BotSelector value={botId} onChange={setBotId} />

// 非受控
<BotSelector
  defaultValue="default-bot-id"
  onBlur={(value) => console.log('选中:', value)}
/>

// 8. 禁用状态
<BotSelector
  disabled={!canEdit}
  value={botId}
  onChange={setBotId}
/>

// 9. 标签筛选
<BotSelector
  value={botId}
  onChange={setBotId}
  filter={{
    tags: ['客服', '销售'],
  }}
/>

// 10. 异步加载
const [loading, setLoading] = useState(false);
const [bots, setBots] = useState([]);

const loadBots = async (keyword) => {
  setLoading(true);
  try {
    const data = await api.bots.search(keyword);
    setBots(data);
  } finally {
    setLoading(false);
  }
};

<BotSelector
  loading={loading}
  value={botId}
  onChange={setBotId}
  loadBots={loadBots}
/>
```

**最佳实践**:

✅ **推荐做法**:
- 大数据量时使用 `searchable` 提升用户体验
- 为 Bot 提供清晰的 `name` 和 `description`
- 使用 `filter` 缩小选择范围
- 自定义渲染时保持简洁

❌ **不推荐做法**:
- 不要在 `onChange` 中进行复杂异步操作
- 不要过度使用 `filter`,可能影响性能
- 不要在 `renderBot` 中渲染大量组件

**注意事项**:
- 多选时 `value` 为数组
- 搜索使用防抖,避免频繁请求
- 需要处理加载和错误状态

**常见问题**:

**Q: 选择后显示的值不正确?**
A: 检查 `value` 是否匹配 Bot 的 ID

```typescript
// 确保类型一致
const botId = 'bot-123'; // string
<BotSelector value={botId} /> // ✅

const botId = 123; // number
<BotSelector value={String(botId)} /> // ✅ 转换为字符串
```

---

### 2.3 布局组件

#### 2.3.1 GlobalLayout 全局布局

**组件路径**: `@coze-studio/layout/global-layout`
**源文件**: `frontend/packages/foundation/layout/src/components/global-layout/index.tsx`

**功能说明**:
应用全局布局,包含顶部导航、侧边栏、内容区域

**Props API**:

```typescript
interface GlobalLayoutProps {
  /**
   * 是否有侧边栏
   * @default false
   */
  hasSider?: boolean;

  /**
   * 横幅内容
   */
  banner?: React.ReactNode;

  /**
   * 子元素
   */
  children: React.ReactNode;
}

interface GlobalLayoutSlots {
  /**
   * 头部插槽
   */
  header?: React.ReactNode;

  /**
   * 头部上方插槽
   */
  headerTop?: React.ReactNode;

  /**
   * 头部下方插槽
   */
  headerBottom?: React.ReactNode;

  /**
   * 自定义 Provider
   */
  customProvider?: React.ComponentType<React.PropsWithChildren<any>>;
}
```

**使用示例**:

```typescript
import { GlobalLayout } from '@coze-studio/layout';

// 1. 基础用法
<GlobalLayout>
  <div>页面内容</div>
</GlobalLayout>

// 2. 带侧边栏
<GlobalLayout hasSider>
  <Sider>
    <Menu />
  </Sider>
  <Layout>
    <Content>页面内容</Content>
  </Layout>
</GlobalLayout>

// 3. 带横幅
<GlobalLayout banner={<Banner />}>
  <div>页面内容</div>
</GlobalLayout>

// 4. 使用插槽
<GlobalLayout>
  <GlobalLayout.Slot name="header">
    <Header />
  </GlobalLayout.Slot>

  <GlobalLayout.Slot name="headerTop">
    <TopBar />
  </GlobalLayout.Slot>

  <Content>页面内容</Content>
</GlobalLayout>

// 5. 自定义 Provider
<GlobalLayout>
  <GlobalLayout.Slot name="customProvider">
    {CustomProvider}
  </GlobalLayout.Slot>

  <Content>页面内容</Content>
</GlobalLayout>

// 6. 完整布局示例
<GlobalLayout hasSider>
  {/* 头部 */}
  <GlobalLayout.Slot name="header">
    <Header>
      <Logo />
      <Navigation />
      <UserMenu />
    </Header>
  </GlobalLayout.Slot>

  {/* 侧边栏 */}
  <Sider width={240}>
    <Menu mode="vertical" />
  </Sider>

  {/* 主内容 */}
  <Layout>
    <Breadcrumb />
    <Content className="p-4">
      {children}
    </Content>
    <Footer />
  </Layout>
</GlobalLayout>

// 7. 响应式布局
<GlobalLayout hasSider>
  <Sider
    width={240}
    breakpoint="lg"
    collapsedWidth="0"
  >
    <Menu />
  </Sider>

  <Layout>
    <Content>页面内容</Content>
  </Layout>
</GlobalLayout>
```

**最佳实践**:

✅ **推荐做法**:
- 使用插槽而非直接嵌套,保持灵活性
- 响应式设计使用 `breakpoint`
- 保持布局层次清晰

❌ **不推荐做法**:
- 不要在多个地方重复使用相同布局
- 不要过度嵌套布局组件
- 不要在布局组件中混入业务逻辑

**注意事项**:
- `hasSider` 必须在顶层声明
- 插槽内容需要是有效的 React 元素
- 横幅内容应轻量,避免影响性能

---

## 后端服务组件清单

### 3.1 微服务架构

ZKER 采用微服务架构,按业务领域拆分:

| 服务名 | 端口 | 职责 | 健康检查 |
|--------|------|------|---------|
| **user-service** | 8001 | 用户管理 | `/health` |
| **auth-service** | 8002 | 认证授权 | `/health` |
| **bot-service** | 8003 | Bot 管理 | `/health` |
| **conversation-service** | 8004 | 会话管理 | `/health` |
| **knowledge-service** | 8005 | 知识库管理 | `/health` |
| **workflow-service** | 8006 | 工作流引擎 | `/health` |
| **routing-service** | 8007 | 智能路由 | `/health` |
| **publishing-service** | 8008 | 发布渠道 | `/health` |
| **billing-service** | 8009 | 计费服务 | `/health` |
| **subscription-service** | 8010 | 订阅服务 | `/health` |
| **monitoring-service** | 8011 | 监控服务 | `/health` |
| **agent-service** | 8012 | Agent 执行 | `/health` |

**服务发现**: 使用 etcd
**API 网关**: Nginx / Kong
**消息队列**: NSQ

---

### 3.2 领域服务

#### 3.2.1 用户服务 (User Service)

**职责**:
- 用户注册、登录、注销
- 用户信息管理
- 用户画像管理

**主要接口**:
```
POST   /api/v1/users/register          # 用户注册
POST   /api/v1/users/login             # 用户登录
POST   /api/v1/users/logout            # 用户登出
GET    /api/v1/users/:id               # 获取用户详情
PUT    /api/v1/users/:id               # 更新用户信息
DELETE /api/v1/users/:id               # 删除用户
GET    /api/v1/users/:id/profile       # 获取用户画像
PUT    /api/v1/users/:id/profile       # 更新用户画像
```

**数据模型**:
```go
type User struct {
    UserID       string    `json:"user_id" gorm:"primaryKey"`
    TenantID     string    `json:"tenant_id" gorm:"index"`
    Username     string    `json:"username" gorm:"uniqueIndex"`
    Email        string    `json:"email" gorm:"uniqueIndex"`
    PasswordHash string    `json:"-"`
    Status       string    `json:"status"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}
```

---

#### 3.2.2 认证服务 (Auth Service)

**职责**:
- JWT Token 签发和验证
- 权限验证
- RBAC 权限管理

**主要接口**:
```
POST   /api/v1/auth/login              # 登录获取 Token
POST   /api/v1/auth/refresh            # 刷新 Token
POST   /api/v1/auth/logout             # 登出注销 Token
GET    /api/v1/auth/permissions        # 获取用户权限
POST   /api/v1/auth/check-permission   # 检查权限
```

**权限模型**:
```go
type Permission struct {
    PermissionID   string   `json:"permission_id"`
    Resource       string   `json:"resource"`       // 资源
    Action         string   `json:"action"`         // 操作
    Effect         string   `json:"effect"`         // 允许/拒绝
}

type Role struct {
    RoleID       string       `json:"role_id"`
    Name         string       `json:"name"`
    Permissions  []Permission `json:"permissions"`
}
```

---

#### 3.2.3 Bot 服务 (Bot Service)

**职责**:
- Bot 创建、编辑、删除
- Bot 版本管理
- Bot 发布管理

**主要接口**:
```
POST   /api/v1/bots                    # 创建 Bot
GET    /api/v1/bots                    # 获取 Bot 列表
GET    /api/v1/bots/:id                # 获取 Bot 详情
PUT    /api/v1/bots/:id                # 更新 Bot
DELETE /api/v1/bots/:id                # 删除 Bot
POST   /api/v1/bots/:id/publish        # 发布 Bot
POST   /api/v1/bots/:id/versions       # 创建 Bot 版本
GET    /api/v1/bots/:id/versions       # 获取 Bot 版本列表
```

---

## 基础设施组件清单

### 4.1 数据库

#### MySQL 8.4.5

**用途**: 主数据库,存储业务数据

**关键配置**:
```ini
[mysqld]
character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci
innodb_buffer_pool_size=4G
max_connections=500
```

**数据库**: `zker`

**表数量**: 157 张

**备份策略**:
- 全量备份: 每天 2:00 AM
- 增量备份: 每小时
- binlog 保留: 7 天

---

#### Milvus v2.5.10

**用途**: 向量数据库,存储文档向量和用户输入向量

**集合**:
- `knowledge_chunks` - 知识库分块向量 (1536 维)
- `user_inputs` - 用户输入向量 (1536 维)
- `intents` - 意图向量 (1536 维)

**检索参数**:
- `metric_type`: `COSINE` 或 `L2`
- `top_k`: 5-10
- `nprobe`: 16

---

### 4.2 缓存

#### Redis 8.0

**用途**: 缓存、会话、分布式锁

**数据结构**:
- String: Token 缓存、配置缓存
- Hash: 用户会话、Bot 状态
- List: 消息队列
- Set: 用户标签、Bot 标签
- ZSet: 排行榜、优先级队列

**缓存策略**:
- 热点数据: 1 小时过期
- 会话数据: 7 天过期
- 配置数据: 永久有效 (主动更新)

---

### 4.3 消息队列

#### NSQ

**用途**: 异步消息处理

**Topics**:
- `user-events` - 用户事件
- `bot-events` - Bot 事件
- `conversation-events` - 会话事件
- `workflow-events` - 工作流事件
- `notification-events` - 通知事件

**Channels**:
- `worker` - 工作进程
- `analytics` - 分析服务
- `audit` - 审计服务

---

### 4.4 搜索引擎

#### Elasticsearch 8.18.0

**用途**: 全文检索、日志检索

**索引**:
- `knowledge-documents` - 知识库文档
- `conversations` - 会话消息
- `logs-*` - 日志 (按日期滚动)

**分析器**:
- `smartcn` - 中文分词
- `standard` - 英文分词

---

### 4.5 对象存储

#### MinIO

**用途**: 文件存储

**Buckets**:
- `avatars` - 用户头像
- `documents` - 文档文件
- `images` - 图片文件
- `videos` - 视频文件
- `audio` - 音频文件
- `models` - 模型文件

**访问策略**:
- 私有桶: 需要签名 URL
- 公共桶: 直接访问 (7 天过期)

---

## 第三方服务集成清单

### 5.1 AI 模型服务

#### OpenAI API

**用途**: GPT-4, GPT-3.5, Embeddings

**配置**:
```yaml
openai:
  api_key: "${OPENAI_API_KEY}"
  base_url: "https://api.openai.com/v1"
  models:
    chat: "gpt-4-turbo-preview"
    embedding: "text-embedding-ada-002"
```

---

#### Anthropic Claude API

**用途**: Claude 3 Opus, Sonnet, Haiku

**配置**:
```yaml
anthropic:
  api_key: "${ANTHROPIC_API_KEY}"
  base_url: "https://api.anthropic.com"
  models:
    chat: "claude-3-opus-20240229"
```

---

### 5.2 语音识别

#### Azure Speech Service

**用途**: 语音转文字、文字转语音

**配置**:
```yaml
azure:
  speech:
    subscription_key: "${AZURE_SPEECH_KEY}"
    region: "eastasia"
```

---

### 5.3 短信/邮件

#### 阿里云短信

**用途**: 验证码、通知短信

**配置**:
```yaml
aliyun:
  sms:
    access_key_id: "${ALIYUN_ACCESS_KEY}"
    access_key_secret: "${ALIYUN_SECRET_KEY}"
    sign_name: "ZKER"
    template_code: "SMS_123456789"
```

---

#### SendGrid

**用途**: 邮件发送

**配置**:
```yaml
sendgrid:
  api_key: "${SENDGRID_API_KEY}"
  from_email: "noreply@zker.com"
  from_name: "ZKER"
```

---

## 组件开发规范

### 6.1 前端组件规范

#### 6.1.1 组件开发原则

**单一职责**:
- 每个组件只做一件事
- 复杂组件拆分为多个子组件

**可复用性**:
- 通过 Props 控制行为和样式
- 避免硬编码

**性能优化**:
- 使用 `React.memo` 避免不必要的重渲染
- 使用 `useMemo` 缓存计算结果
- 使用 `useCallback` 缓存回调函数

---

#### 6.1.2 组件命名规范

**组件名称**: PascalCase
```
ChatInput.tsx       ✅
chat-input.tsx      ❌
chat_input.tsx      ❌
```

**Hook 名称**: `use` 前缀
```
useBotStore.ts      ✅
botStore.ts         ❌
getBotStore.ts      ❌
```

**工具函数**: camelCase
```
formatDate.ts       ✅
FormatDate.ts       ❌
format_date.ts      ❌
```

---

#### 6.1.3 文件组织规范

```
component-name/
├── index.tsx              # 组件入口
├── component-name.tsx     # 组件实现
├── types.ts              # 类型定义
├── styles.module.less    # 样式文件
├── __tests__/            # 测试文件
│   ├── component-name.test.tsx
│   └── component-name.test.ts
└── README.md             # 组件文档
```

---

#### 6.1.4 TypeScript 类型定义

**Props 接口命名**: `组件名 + Props`
```typescript
interface ChatInputProps {
  placeholder: string;
  onSend: (content: string) => void;
}
```

**返回值类型**: 明确指定
```typescript
const formatMessage = (message: Message): FormattedMessage => {
  // ...
};
```

---

#### 6.1.5 组件文档模板

每个组件应包含:

```markdown
# 组件名称

## 功能说明
简要描述组件功能

## Props API
| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| ... | ... | ... | ... |

## 使用示例
\`\`\`typescript
// 示例代码
\`\`\`

## 注意事项
- 注意事项 1
- 注意事项 2

## 常见问题
**Q: 问题?**
A: 答案
```

---

### 6.2 后端服务规范

#### 6.2.1 服务分层

```
domain/         # 领域层 - 实体和值对象
application/    # 应用层 - 应用服务和用例
api/            # API 层 - HTTP 处理器和路由
infra/          # 基础设施层 - 技术实现
crossdomain/    # 跨领域 - 关注点横切
```

---

#### 6.2.2 接口设计规范

**RESTful API**:
```
GET    /api/v1/resources          # 列表
POST   /api/v1/resources          # 创建
GET    /api/v1/resources/:id      # 详情
PUT    /api/v1/resources/:id      # 更新
DELETE /api/v1/resources/:id      # 删除
```

**响应格式**:
```json
{
  "code": "SUCCESS",
  "message": "操作成功",
  "data": { ... },
  "request_id": "req_1234567890",
  "timestamp": "2025-01-01T12:00:00Z"
}
```

---

## 附录

### A. 组件快速索引

**前端组件**:
- Button → `@coze-studio/ui/button`
- Input → `@coze-studio/ui/input`
- Select → `@coze-studio/ui/select`
- Modal → `@coze-studio/ui/modal`
- Table → `@coze-studio/ui/table`
- ChatInput → `@coze-studio/chat/chat-input`
- MessageBox → `@coze-studio/chat/message-box`
- BotSelector → `@coze-studio/bot/bot-selector`
- GlobalLayout → `@coze-studio/layout/global-layout`

**后端服务**:
- User Service → `user-service:8001`
- Auth Service → `auth-service:8002`
- Bot Service → `bot-service:8003`
- Conversation Service → `conversation-service:8004`
- Knowledge Service → `knowledge-service:8005`
- Workflow Service → `workflow-service:8006`
- Routing Service → `routing-service:8007`

---

### B. 相关文档

- [ZKER 统一错误码定义规范](./ZKER-统一错误码定义规范.md)
- [可执行性文档缺失分析与补充方案](./可执行性文档缺失分析与补充方案_v1.0.md)
- [数据库设计完整交付清单](./数据库设计完整交付清单.md)
- [数据库DDL脚本使用指南](./数据库DDL脚本使用指南.md)

---

**文档状态**: ✅ 初版完成
**维护责任**: ZKER 架构团队
**下次更新**: 根据组件变更定期更新
