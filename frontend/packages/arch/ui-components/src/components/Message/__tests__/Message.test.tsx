// frontend/packages/arch/ui-components/src/components/Message/__tests__/Message.test.tsx

import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { Message } from '../Message';

describe('Message', () => {
  beforeEach(() => {
    jest.useFakeTimers();
  });

  afterEach(() => {
    jest.useRealTimers();
  });

  it('should render message with content', () => {
    render(<Message content="Success message" type="success" />);
    expect(screen.getByText('Success message')).toBeInTheDocument();
  });

  it('should render with correct type styling', () => {
    const { container } = render(<Message content="Info" type="info" />);
    const message = container.querySelector('.container');
    expect(message).toHaveStyle({ backgroundColor: '#1890ff' });
  });

  it('should render icon', () => {
    render(<Message content="Error" type="error" />);
    expect(screen.getByText('✕')).toBeInTheDocument();
  });

  it('should auto close after duration', async () => {
    const { container } = render(<Message content="Auto close" duration={3000} />);
    expect(container.firstChild).toBeInTheDocument();

    jest.advanceTimersByTime(3000);

    await waitFor(() => {
      expect(container.firstChild).toBeNull();
    });
  });

  it('should not auto close when duration is 0', () => {
    const { container } = render(<Message content="Persistent" duration={0} />);

    jest.advanceTimersByTime(10000);

    expect(container.firstChild).toBeInTheDocument();
  });

  it('should call onClose when auto closed', () => {
    const handleClose = jest.fn();
    render(<Message content="Message" duration={1000} onClose={handleClose} />);

    jest.advanceTimersByTime(1000);

    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it('should render custom icon', () => {
    render(<Message content="Custom" icon="🎉" />);
    expect(screen.getByText('🎉')).toBeInTheDocument();
  });

  it('should apply custom className', () => {
    const { container } = render(<Message content="Message" className="custom-class" />);
    const message = container.firstChild as HTMLElement;
    expect(message).toHaveClass('custom-class');
  });
});
