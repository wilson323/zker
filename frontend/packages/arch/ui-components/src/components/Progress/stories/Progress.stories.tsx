// frontend/packages/arch/ui-components/src/components/Progress/stories/Progress.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { Progress } from '../Progress';

const meta: Meta<typeof Progress> = {
  title: 'Components/Progress',
  component: Progress,
  tags: ['autodocs'],
  argTypes: {
    percent: {
      control: 'number',
      description: '百分比（0-100）',
    },
    type: {
      control: 'select',
      options: ['line', 'circle', 'dashboard'],
      description: '类型',
    },
    status: {
      control: 'select',
      options: ['normal', 'active', 'success', 'exception'],
      description: '状态',
    },
    showInfo: {
      control: 'boolean',
      description: '是否显示百分比',
    },
    strokeWidth: {
      control: 'number',
      description: '进度条线的宽度',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Progress>;

export const Default: Story = {
  args: {
    percent: 50,
  },
};

export const Active: Story = {
  args: {
    percent: 60,
    status: 'active',
  },
};

export const Success: Story = {
  args: {
    percent: 100,
    status: 'success',
  },
};

export const Exception: Story = {
  args: {
    percent: 50,
    status: 'exception',
  },
};

export const WithoutInfo: Story = {
  args: {
    percent: 75,
    showInfo: false,
  },
};

export const Circle: Story = {
  args: {
    percent: 75,
    type: 'circle',
  },
};

export const CustomFormat: Story = {
  args: {
    percent: 30,
    format: (percent: number) => `${percent} of 100`,
  },
};

export const AllSizes: Story = {
  render: () => (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div>
        <h4>Small</h4>
        <Progress percent={50} strokeWidth={6} />
      </div>
      <div>
        <h4>Default</h4>
        <Progress percent={50} strokeWidth={10} />
      </div>
      <div>
        <h4>Large</h4>
        <Progress percent={50} strokeWidth={15} />
      </div>
    </div>
  ),
};

export const AllStatuses: Story = {
  render: () => (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
      <div>
        <h4>Normal</h4>
        <Progress percent={50} status="normal" />
      </div>
      <div>
        <h4>Active</h4>
        <Progress percent={50} status="active" />
      </div>
      <div>
        <h4>Success</h4>
        <Progress percent={100} status="success" />
      </div>
      <div>
        <h4>Exception</h4>
        <Progress percent={50} status="exception" />
      </div>
    </div>
  ),
};

export const AllTypes: Story = {
  render: () => (
    <div style={{ display: 'flex', gap: '40px' }}>
      <div>
        <h4>Line</h4>
        <Progress percent={50} type="line" />
      </div>
      <div>
        <h4>Circle</h4>
        <Progress percent={75} type="circle" width={120} />
      </div>
    </div>
  ),
};
