// frontend/packages/arch/ui-components/src/components/Dropdown/stories/Dropdown.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { Dropdown } from '../Dropdown';

const meta: Meta<typeof Dropdown> = {
  title: 'Components/Dropdown',
  component: Dropdown,
  tags: ['autodocs'],
  argTypes: {
    triggerType: {
      control: 'select',
      options: ['click', 'hover'],
      description: '触发方式',
    },
    disabled: {
      control: 'boolean',
      description: '是否禁用',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Dropdown>;

const menuItems = [
  { key: '1', label: '菜单项 1' },
  { key: '2', label: '菜单项 2' },
  { key: '3', label: '菜单项 3' },
  { key: '4', label: '菜单项 4', divided: true },
];

export const Default: Story = {
  args: {
    menu: menuItems,
    trigger: <button>点击打开菜单</button>,
  },
};

export const Hover: Story = {
  args: {
    menu: menuItems,
    trigger: <button>悬停打开菜单</button>,
    triggerType: 'hover',
  },
};

export const Disabled: Story = {
  args: {
    menu: menuItems,
    trigger: <button>禁用菜单</button>,
    disabled: true,
  },
};

export const WithDangerItem: Story = {
  args: {
    menu: [
      { key: '1', label: '查看' },
      { key: '2', label: '编辑' },
      { key: '3', label: '删除', danger: true },
    ],
    trigger: <button>操作菜单</button>,
  },
};

export const WithNestedMenu: Story = {
  args: {
    menu: [
      {
        key: '1',
        label: '文件操作',
        children: [
          { key: '1-1', label: '新建' },
          { key: '1-2', label: '打开' },
          { key: '1-3', label: '保存' },
        ],
      },
      {
        key: '2',
        label: '编辑操作',
        children: [
          { key: '2-1', label: '复制' },
          { key: '2-2', label: '粘贴' },
        ],
      },
    ],
    trigger: <button>嵌套菜单</button>,
  },
};

export const WithDisabledItems: Story = {
  args: {
    menu: [
      { key: '1', label: '可用项' },
      { key: '2', label: '禁用项', disabled: true },
      { key: '3', label: '另一可用项' },
    ],
    trigger: <button>包含禁用项</button>,
  },
};

export const AllVariants: Story = {
  render: () => (
    <div style={{ display: 'flex', gap: '20px' }}>
      <Dropdown
        menu={menuItems}
        trigger={<button>点击触发</button>}
      />
      <Dropdown
        menu={menuItems}
        trigger={<button>悬停触发</button>}
        triggerType="hover"
      />
      <Dropdown
        menu={[
          { key: '1', label: '选项1' },
          { key: '2', label: '删除', danger: true, divided: true },
        ]}
        trigger={<button>危险项</button>}
      />
    </div>
  ),
};
