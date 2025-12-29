// frontend/packages/arch/ui-components/src/components/Spin/stories/Spin.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { Spin } from '../Spin';

const meta: Meta<typeof Spin> = {
  title: 'Components/Spin',
  component: Spin,
  tags: ['autodocs'],
  argTypes: {
    spinning: {
      control: 'boolean',
      description: '是否加载中',
    },
    size: {
      control: 'select',
      options: ['small', 'default', 'large'],
      description: '尺寸',
    },
    tip: {
      control: 'text',
      description: '提示文本',
    },
    delay: {
      control: 'number',
      description: '延迟显示（毫秒）',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Spin>;

export const Default: Story = {
  args: {
    spinning: true,
  },
};

export const Small: Story = {
  args: {
    spinning: true,
    size: 'small',
  },
};

export const Large: Story = {
  args: {
    spinning: true,
    size: 'large',
  },
};

export const WithTip: Story = {
  args: {
    spinning: true,
    tip: '加载中...',
  },
};

export const WithDelay: Story = {
  args: {
    spinning: true,
    delay: 500,
    tip: '延迟显示',
  },
};

export const AsWrapper: Story = {
  args: {
    spinning: true,
    tip: '加载中...',
    children: (
      <div style={{ padding: '20px', background: '#f5f5f5' }}>
        <p>这是包装的内容</p>
        <p>加载时会显示遮罩层</p>
      </div>
    ),
  },
};

export const NotSpinning: Story = {
  args: {
    spinning: false,
    children: <div>内容区域</div>,
  },
};

export const AllVariants: Story = {
  render: () => (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div>
        <h4>Small</h4>
        <Spin size="small" />
      </div>
      <div>
        <h4>Default</h4>
        <Spin />
      </div>
      <div>
        <h4>Large</h4>
        <Spin size="large" />
      </div>
      <div>
        <h4>With Tip</h4>
        <Spin tip="加载中..." />
      </div>
    </div>
  ),
};
