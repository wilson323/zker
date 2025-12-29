// frontend/packages/arch/ui-components/src/components/Spin/__tests__/Spin.test.tsx

import React from 'react';
import { render, screen } from '@testing-library/react';
import { Spin } from '../Spin';

describe('Spin', () => {
  it('should render spinning indicator', () => {
    render(<Spin />);
    expect(screen.getByRole('status')).toBeInTheDocument();
  });

  it('should not render when spinning is false and no children', () => {
    const { container } = render(<Spin spinning={false} />);
    expect(container.firstChild).toBeNull();
  });

  it('should render custom tip', () => {
    render(<Spin tip="Loading..." />);
    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('should render with custom size', () => {
    const { container } = render(<Spin size="large" />);
    const spinner = container.querySelector('.spinner');
    expect(spinner).toBeInTheDocument();
  });

  it('should render children when not spinning', () => {
    render(
      <Spin spinning={false}>
        <div>Content</div>
      </Spin>
    );
    expect(screen.getByText('Content')).toBeInTheDocument();
  });

  it('should render overlay when spinning with children', () => {
    const { container } = render(
      <Spin spinning>
        <div>Content</div>
      </Spin>
    );
    const overlay = container.querySelector('.overlay');
    expect(overlay).toBeInTheDocument();
  });

  it('should apply delay before showing', () => {
    jest.useFakeTimers();
    render(<Spin delay={500} />);
    expect(screen.queryByRole('status')).not.toBeInTheDocument();

    jest.advanceTimersByTime(500);
    expect(screen.getByRole('status')).toBeInTheDocument();

    jest.useRealTimers();
  });
});
