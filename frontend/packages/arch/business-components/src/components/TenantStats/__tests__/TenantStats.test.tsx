// frontend/packages/arch/business-components/src/components/TenantStats/__tests__/TenantStats.test.tsx

import React from 'react';
import { render, screen } from '@testing-library/react';
import { TenantStats } from '../TenantStats';
import { TenantStatsDTO } from '@coze-studio/api-client';

const mockStats: TenantStatsDTO = {
  tenant_id: 'tenant-1',
  total_bots: 100,
  active_bots: 75,
  total_users: 1000,
  active_users: 800,
  total_messages: 50000,
  total_tokens_used: 1000000,
  storage_used_mb: 5120,
};

describe('TenantStats', () => {
  it('should render all stat cards', () => {
    render(<TenantStats stats={mockStats} />);

    expect(screen.getByText('机器人总数')).toBeInTheDocument();
    expect(screen.getByText('用户总数')).toBeInTheDocument();
    expect(screen.getByText('消息总数')).toBeInTheDocument();
    expect(screen.getByText('Token使用量')).toBeInTheDocument();
    expect(screen.getByText('存储使用量')).toBeInTheDocument();
  });

  it('should display correct values', () => {
    render(<TenantStats stats={mockStats} />);

    expect(screen.getByText('100')).toBeInTheDocument(); // total_bots
    expect(screen.getByText('1000')).toBeInTheDocument(); // total_users
    expect(screen.getByText('50000')).toBeInTheDocument(); // total_messages
    expect(screen.getByText('1000000')).toBeInTheDocument(); // total_tokens_used
    expect(screen.getByText('5120 MB')).toBeInTheDocument(); // storage_used_mb
  });

  it('should display active counts and percentages', () => {
    render(<TenantStats stats={mockStats} />);

    expect(screen.getByText('75')).toBeInTheDocument(); // active_bots
    expect(screen.getByText('800')).toBeInTheDocument(); // active_users
    expect(screen.getByText('(75%)')).toBeInTheDocument(); // bots percentage
    expect(screen.getByText('(80%)')).toBeInTheDocument(); // users percentage
  });

  it('should render progress bars when showChart is true', () => {
    const { container } = render(<TenantStats stats={mockStats} showChart={true} />);

    const bars = container.querySelectorAll('.bar');
    expect(bars.length).toBeGreaterThan(0);
  });

  it('should not render progress bars when showChart is false', () => {
    const { container } = render(<TenantStats stats={mockStats} showChart={false} />);

    const bars = container.querySelectorAll('.bar');
    expect(bars.length).toBe(0);
  });

  it('should apply custom className', () => {
    const { container } = render(
      <TenantStats stats={mockStats} className="custom-class" />
    );

    const wrapper = container.firstChild as HTMLElement;
    expect(wrapper).toHaveClass('custom-class');
  });

  it('should render icons for each stat', () => {
    render(<TenantStats stats={mockStats} />);

    expect(screen.getByText('🤖')).toBeInTheDocument(); // bot icon
    expect(screen.getByText('👥')).toBeInTheDocument(); // user icon
    expect(screen.getByText('💬')).toBeInTheDocument(); // message icon
    expect(screen.getByText('📊')).toBeInTheDocument(); // token icon
    expect(screen.getByText('💾')).toBeInTheDocument(); // storage icon
  });
});
