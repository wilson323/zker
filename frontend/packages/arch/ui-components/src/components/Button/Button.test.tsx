// frontend/packages/arch/ui-components/src/components/Button/Button.test.tsx

import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Button } from './Button';

describe('Button', () => {
  it('renders children correctly', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });

  it('calls onClick when clicked', () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click me</Button>);

    fireEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('shows loading spinner when loading', () => {
    render(<Button loading>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
    expect(document.querySelector('.spinner')).toBeInTheDocument();
  });

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Click me</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('applies correct variant class', () => {
    const { container: container1 } = render(<Button variant="primary">Primary</Button>);
    const { container: container2 } = render(<Button variant="danger">Danger</Button>);

    expect(container1.querySelector('button')).toHaveStyle({ backgroundColor: '#1890FF' });
    expect(container2.querySelector('button')).toHaveStyle({ backgroundColor: '#F5222D' });
  });

  it('applies correct size class', () => {
    const { container: container1 } = render(<Button size="sm">Small</Button>);
    const { container: container2 } = render(<Button size="lg">Large</Button>);

    expect(container1.querySelector('button')).toHaveStyle({ padding: '4px 8px' });
    expect(container2.querySelector('button')).toHaveStyle({ padding: '16px 24px' });
  });
});
