// frontend/packages/arch/ui-components/src/components/Select/Select.test.tsx

import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Select } from './Select';

const options = [
  { label: 'Option 1', value: '1' },
  { label: 'Option 2', value: '2' },
  { label: 'Option 3', value: '3' },
];

describe('Select', () => {
  it('renders placeholder', () => {
    render(
      <Select
        options={options}
        placeholder="Select an option"
      />
    );
    expect(screen.getByText('Select an option')).toBeInTheDocument();
  });

  it('opens dropdown on click', () => {
    render(<Select options={options} />);
    const trigger = screen.getByText('Select an option');

    fireEvent.click(trigger);

    expect(screen.getByText('Option 1')).toBeInTheDocument();
    expect(screen.getByText('Option 2')).toBeInTheDocument();
    expect(screen.getByText('Option 3')).toBeInTheDocument();
  });

  it('calls onChange when option selected', () => {
    const handleChange = vi.fn();
    render(
      <Select
        options={options}
        onChange={handleChange}
      />
    );

    const trigger = screen.getByText('Select an option');
    fireEvent.click(trigger);

    const option = screen.getByText('Option 1');
    fireEvent.click(option);

    expect(handleChange).toHaveBeenCalledWith('1');
  });

  it('displays selected value', () => {
    render(
      <Select
        options={options}
        value="2"
      />
    );
    expect(screen.getByText('Option 2')).toBeInTheDocument();
  });
});
