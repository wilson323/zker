// frontend/packages/arch/ui-components/src/components/DatePicker/stories/DatePicker.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { DatePicker } from '../DatePicker';

const meta: Meta<typeof DatePicker> = {
  title: 'Components/DatePicker',
  component: DatePicker,
  tags: ['autodocs'],
  argTypes: {
    type: {
      control: 'select',
      options: ['date', 'datetime', 'dateRange', 'month', 'year'],
      description: '选择器类型',
    },
    placeholder: {
      control: 'text',
      description: '占位文本',
    },
    format: {
      control: 'text',
      description: '日期格式',
    },
    disabled: {
      control: 'boolean',
      description: '是否禁用',
    },
  },
};

export default meta;
type Story = StoryObj<typeof DatePicker>;

export const Default: Story = {
  args: {
    placeholder: '请选择日期',
  },
};

export const WithValue: Story = {
  args: {
    placeholder: '请选择日期',
    value: '2024-01-15',
  },
};

export const Disabled: Story = {
  args: {
    placeholder: '禁用日期选择器',
    disabled: true,
  },
};

export const CustomFormat: Story = {
  args: {
    placeholder: '请选择日期',
    format: 'YYYY/MM/DD',
  },
};

export const Multiple: Story = {
  render: () => (
    <div style={{ display: 'flex', gap: '20px' }}>
      <DatePicker placeholder="请选择开始日期" />
      <DatePicker placeholder="请选择结束日期" />
    </div>
  ),
};

export const AllVariants: Story = {
  render: () => (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div>
        <h4>默认</h4>
        <DatePicker placeholder="请选择日期" />
      </div>
      <div>
        <h4>禁用</h4>
        <DatePicker placeholder="禁用状态" disabled />
      </div>
      <div>
        <h4>自定义格式</h4>
        <DatePicker placeholder="YYYY/MM/DD" format="YYYY/MM/DD" />
      </div>
    </div>
  ),
};
