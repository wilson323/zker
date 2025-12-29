// frontend/packages/arch/ui-components/src/components/Alert/stories/Alert.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { Alert } from '../Alert';

const meta: Meta<typeof Alert> = {
  title: 'Components/Alert',
  component: Alert,
  tags: ['autodocs'],
  argTypes: {
    type: {
      control: 'select',
      options: ['success', 'info', 'warning', 'error'],
      description: '类型',
    },
    message: {
      control: 'text',
      description: '提示内容',
    },
    description: {
      control: 'text',
      description: '辅助性文字',
    },
    closable: {
      control: 'boolean',
      description: '是否可关闭',
    },
    showIcon: {
      control: 'boolean',
      description: '是否显示图标',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Alert>;

export const Success: Story = {
  args: {
    type: 'success',
    message: '操作成功',
    description: '您的更改已保存',
    showIcon: true,
  },
};

export const Info: Story = {
  args: {
    type: 'info',
    message: '提示信息',
    description: '这是一条提示信息',
    showIcon: true,
  },
};

export const Warning: Story = {
  args: {
    type: 'warning',
    message: '警告信息',
    description: '请注意检查输入内容',
    showIcon: true,
  },
};

export const Error: Story = {
  args: {
    type: 'error',
    message: '错误信息',
    description: '操作失败，请重试',
    showIcon: true,
  },
};

export const Closable: Story = {
  args: {
    type: 'info',
    message: '可关闭的提示',
    closable: true,
    showIcon: true,
  },
};

export const WithoutIcon: Story = {
  args: {
    type: 'info',
    message: '不带图标的提示',
  },
};

export const WithLongDescription: Story = {
  args: {
    type: 'warning',
    message: '系统维护通知',
    description: '系统将于今晚 23:00 进行维护，预计持续 2 小时。请提前保存您的工作内容。',
    showIcon: true,
    closable: true,
  },
};

export const AllTypes: Story = {
  render: () => (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
      <Alert type="success" message="成功提示" showIcon />
      <Alert type="info" message="信息提示" showIcon />
      <Alert type="warning" message="警告提示" showIcon />
      <Alert type="error" message="错误提示" showIcon />
    </div>
  ),
};
