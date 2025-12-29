// frontend/packages/arch/ui-components/src/components/Message/stories/Message.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { Message } from '../Message';

const meta: Meta<typeof Message> = {
  title: 'Components/Message',
  component: Message,
  tags: ['autodocs'],
  argTypes: {
    type: {
      control: 'select',
      options: ['success', 'info', 'warning', 'error'],
      description: '类型',
    },
    content: {
      control: 'text',
      description: '提示内容',
    },
    duration: {
      control: 'number',
      description: '持续时间（毫秒）',
    },
    icon: {
      control: 'text',
      description: '自定义图标',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Message>;

export const Success: Story = {
  args: {
    type: 'success',
    content: '操作成功',
  },
};

export const Info: Story = {
  args: {
    type: 'info',
    content: '这是一条信息',
  },
};

export const Warning: Story = {
  args: {
    type: 'warning',
    content: '请注意',
  },
};

export const Error: Story = {
  args: {
    type: 'error',
    content: '操作失败',
  },
};

export const WithDuration: Story = {
  args: {
    type: 'info',
    content: '5秒后关闭',
    duration: 5000,
  },
};

export const Persistent: Story = {
  args: {
    type: 'info',
    content: '不自动关闭',
    duration: 0,
  },
};

export const CustomIcon: Story = {
  args: {
    type: 'success',
    content: '自定义图标',
    icon: '🎉',
  },
};

export const AllTypes: Story = {
  render: () => (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
      <Message type="success" content="成功消息" />
      <Message type="info" content="信息消息" />
      <Message type="warning" content="警告消息" />
      <Message type="error" content="错误消息" />
    </div>
  ),
};
