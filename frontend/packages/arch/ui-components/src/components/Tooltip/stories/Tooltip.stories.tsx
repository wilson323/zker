// frontend/packages/arch/ui-components/src/components/Tooltip/stories/Tooltip.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { Tooltip } from '../Tooltip';

const meta: Meta<typeof Tooltip> = {
  title: 'Components/Tooltip',
  component: Tooltip,
  tags: ['autodocs'],
  argTypes: {
    title: {
      control: 'text',
      description: '提示内容',
    },
    placement: {
      control: 'select',
      options: ['top', 'bottom', 'left', 'right'],
      description: '位置',
    },
    delay: {
      control: 'number',
      description: '延迟显示（毫秒）',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Tooltip>;

export const Top: Story = {
  args: {
    title: '顶部提示',
    placement: 'top',
    children: <button>鼠标悬停查看上方提示</button>,
  },
};

export const Bottom: Story = {
  args: {
    title: '底部提示',
    placement: 'bottom',
    children: <button>鼠标悬停查看下方提示</button>,
  },
};

export const Left: Story = {
  args: {
    title: '左侧提示',
    placement: 'left',
    children: <button>鼠标悬停查看左侧提示</button>,
  },
};

export const Right: Story = {
  args: {
    title: '右侧提示',
    placement: 'right',
    children: <button>鼠标悬停查看右侧提示</button>,
  },
};

export const WithDelay: Story = {
  args: {
    title: '延迟显示',
    placement: 'top',
    delay: 500,
    children: <button>悬停500ms后显示</button>,
  },
};

export const WithLongText: Story = {
  args: {
    title: '这是一段很长的提示文本，用于测试提示框的换行和宽度适应情况',
    placement: 'top',
    children: <button>长文本提示</button>,
  },
};

export const AllPlacements: Story = {
  render: () => (
    <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '20px' }}>
      <Tooltip title="顶部提示" placement="top">
        <button>Top</button>
      </Tooltip>
      <Tooltip title="底部提示" placement="bottom">
        <button>Bottom</button>
      </Tooltip>
      <Tooltip title="左侧提示" placement="left">
        <button>Left</button>
      </Tooltip>
      <Tooltip title="右侧提示" placement="right">
        <button>Right</button>
      </Tooltip>
    </div>
  ),
};
