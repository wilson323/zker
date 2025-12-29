// frontend/packages/arch/ui-components/src/components/Progress/__tests__/Progress.test.tsx

import React from 'react';
import { render, screen } from '@testing-library/react';
import { Progress } from '../Progress';

describe('Progress', () => {
  it('should render progress bar with correct percent', () => {
    const { container } = render(<Progress percent={50} />);
    const bg = container.querySelector('.bg');
    expect(bg).toHaveStyle({ width: '50%' });
  });

  it('should clamp percent between 0 and 100', () => {
    const { container: container1 } = render(<Progress percent={-10} />);
    const { container: container2 } = render(<Progress percent={150} />);

    const bg1 = container1.querySelector('.bg');
    const bg2 = container2.querySelector('.bg');

    expect(bg1).toHaveStyle({ width: '0%' });
    expect(bg2).toHaveStyle({ width: '100%' });
  });

  it('should render info text by default', () => {
    render(<Progress percent={75} />);
    expect(screen.getByText('75%')).toBeInTheDocument();
  });

  it('should not render info when showInfo is false', () => {
    render(<Progress percent={75} showInfo={false} />);
    expect(screen.queryByText('75%')).not.toBeInTheDocument();
  });

  it('should render success icon when status is success', () => {
    render(<Progress percent={100} status="success" />);
    expect(screen.getByText('✓')).toBeInTheDocument();
  });

  it('should render exception icon when status is exception', () => {
    render(<Progress percent={50} status="exception" />);
    expect(screen.getByText('✗')).toBeInTheDocument();
  });

  it('should render custom format', () => {
    render(<Progress percent={30} format={(percent) => `${percent} of 100`} />);
    expect(screen.getByText('30 of 100')).toBeInTheDocument();
  });

  it('should render circle type', () => {
    const { container } = render(<Progress percent={75} type="circle" />);
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
  });

  it('should animate when status is active', () => {
    const { container } = render(<Progress percent={50} status="active" />);
    const activeBg = container.querySelector('.activeBg');
    expect(activeBg).toBeInTheDocument();
  });
});
