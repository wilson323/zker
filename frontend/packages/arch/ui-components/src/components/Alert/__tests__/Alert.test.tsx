// frontend/packages/arch/ui-components/src/components/Alert/__tests__/Alert.test.tsx

import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { Alert } from '../Alert';

describe('Alert', () => {
  it('should render alert with message', () => {
    render(<Alert message="This is an alert" />);
    expect(screen.getByText('This is an alert')).toBeInTheDocument();
  });

  it('should render with correct type styling', () => {
    const { container } = render(<Alert message="Success" type="success" />);
    const alert = container.querySelector('.container');
    expect(alert).toHaveStyle({ backgroundColor: '#f6ffed' });
  });

  it('should render description', () => {
    render(
      <Alert message="Warning" description="This is a warning message" type="warning" />
    );
    expect(screen.getByText('This is a warning message')).toBeInTheDocument();
  });

  it('should render close button when closable', () => {
    const { container } = render(<Alert message="Alert" closable />);
    const closeButton = container.querySelector('.closeButton');
    expect(closeButton).toBeInTheDocument();
  });

  it('should hide alert when close button clicked', () => {
    const { container } = render(<Alert message="Alert" closable />);
    const closeButton = container.querySelector('.closeButton') as HTMLButtonElement;

    fireEvent.click(closeButton);

    expect(container.firstChild).toBeNull();
  });

  it('should call onClose when close button clicked', () => {
    const handleClose = jest.fn();
    const { container } = render(<Alert message="Alert" closable onClose={handleClose} />);
    const closeButton = container.querySelector('.closeButton') as HTMLButtonElement;

    fireEvent.click(closeButton);

    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it('should render icon when showIcon is true', () => {
    const { container } = render(<Alert message="Info" type="info" showIcon />);
    const icon = container.querySelector('.icon');
    expect(icon).toBeInTheDocument();
  });

  it('should render custom icon', () => {
    render(<Alert message="Custom" icon="🎉" showIcon />);
    expect(screen.getByText('🎉')).toBeInTheDocument();
  });

  it('should apply custom className', () => {
    const { container } = render(<Alert message="Alert" className="custom-class" />);
    const alert = container.firstChild as HTMLElement;
    expect(alert).toHaveClass('custom-class');
  });
});
