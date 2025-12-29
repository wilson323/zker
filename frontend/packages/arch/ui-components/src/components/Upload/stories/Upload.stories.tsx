// frontend/packages/arch/ui-components/src/components/Upload/stories/Upload.stories.tsx

import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { Upload } from '../Upload';

const meta: Meta<typeof Upload> = {
  title: 'Components/Upload',
  component: Upload,
  tags: ['autodocs'],
  argTypes: {
    accept: {
      control: 'text',
      description: '接受的文件类型',
    },
    multiple: {
      control: 'boolean',
      description: '是否支持多文件上传',
    },
    drag: {
      control: 'boolean',
      description: '是否支持拖拽上传',
    },
    disabled: {
      control: 'boolean',
      description: '是否禁用',
    },
    maxCount: {
      control: 'number',
      description: '最大上传数量',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Upload>;

function UploadWrapper(props: any) {
  const [fileList, setFileList] = useState<any[]>([]);

  return <Upload {...props} fileList={fileList} onChange={setFileList} />;
}

export const Default: Story = {
  render: () => <UploadWrapper action="/api/upload" />,
};

export const DragUpload: Story = {
  render: () => <UploadWrapper action="/api/upload" drag />,
};

export const Multiple: Story = {
  render: () => <UploadWrapper action="/api/upload" multiple />,
};

export const WithLimit: Story = {
  render: () => <UploadWrapper action="/api/upload" maxCount={3} />,
};

export const ImageOnly: Story = {
  render: () => <UploadWrapper action="/api/upload" accept="image/*" />,
};

export const Disabled: Story = {
  render: () => <UploadWrapper action="/api/upload" disabled />,
};

export const CustomTrigger: Story = {
  render: () => (
    <UploadWrapper action="/api/upload">
      <button style={{ padding: '8px 16px' }}>自定义上传按钮</button>
    </UploadWrapper>
  ),
};

export const AllVariants: Story = {
  render: () => (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '40px' }}>
      <div>
        <h4>默认上传</h4>
        <UploadWrapper action="/api/upload" />
      </div>
      <div>
        <h4>拖拽上传</h4>
        <UploadWrapper action="/api/upload" drag />
      </div>
      <div>
        <h4>多文件上传</h4>
        <UploadWrapper action="/api/upload" multiple />
      </div>
      <div>
        <h4>仅图片</h4>
        <UploadWrapper action="/api/upload" accept="image/*" />
      </div>
    </div>
  ),
};
