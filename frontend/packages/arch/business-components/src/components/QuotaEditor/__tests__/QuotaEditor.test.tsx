// frontend/packages/arch/business-components/src/components/QuotaEditor/__tests__/QuotaEditor.test.tsx

import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { QuotaEditor } from '../QuotaEditor';
import { QuotaResourceType, QuotaUnit } from '@coze-studio/api-client';

const mockQuotas = [
  {
    resource_type: QuotaResourceType.BOT,
    limit: 100,
    unit: QuotaUnit.COUNT,
    is_soft_limit: false,
    overage_fee: 0,
  },
  {
    resource_type: QuotaResourceType.USER,
    limit: 1000,
    unit: QuotaUnit.COUNT,
    is_soft_limit: true,
    overage_fee: 0.01,
  },
];

describe('QuotaEditor', () => {
  it('should render quota items', () => {
    render(<QuotaEditor quotas={mockQuotas} />);

    expect(screen.getByText('机器人数量')).toBeInTheDocument();
    expect(screen.getByText('用户数量')).toBeInTheDocument();
  });

  it('should render slider and input for each quota', () => {
    const { container } = render(<QuotaEditor quotas={mockQuotas} />);

    const sliders = container.querySelectorAll('.control input[type="range"]');
    const inputs = container.querySelectorAll('.control input[type="number"]');

    expect(sliders.length).toBe(2);
    expect(inputs.length).toBe(2);
  });

  it('should call onChange when quota limit changes', () => {
    const handleChange = jest.fn();
    const { container } = render(
      <QuotaEditor quotas={mockQuotas} onChange={handleChange} />
    );

    const firstInput = container.querySelectorAll('.control input[type="number"]')[0];
    fireEvent.change(firstInput, { target: { value: '200' } });

    expect(handleChange).toHaveBeenCalled();
  });

  it('should toggle soft limit checkbox', () => {
    const handleChange = jest.fn();
    const { container } = render(
      <QuotaEditor quotas={mockQuotas} onChange={handleChange} />
    );

    const checkboxes = container.querySelectorAll('input[type="checkbox"]');
    const firstCheckbox = checkboxes[0];

    fireEvent.click(firstCheckbox);

    expect(handleChange).toHaveBeenCalled();
  });

  it('should show overage fee input when soft limit is enabled', () => {
    const { container } = render(<QuotaEditor quotas={mockQuotas} />);

    const overageFeeInputs = container.querySelectorAll('.feeInput');
    expect(overageFeeInputs.length).toBeGreaterThan(0);
  });

  it('should calculate and display price when showPrice is true', () => {
    render(<QuotaEditor quotas={mockQuotas} showPrice={true} />);

    expect(screen.getByText(/预计费用/)).toBeInTheDocument();
  });

  it('should display total price', () => {
    render(<QuotaEditor quotas={mockQuotas} showPrice={true} />);

    expect(screen.getByText(/总计/)).toBeInTheDocument();
  });

  it('should not display price when showPrice is false', () => {
    render(<QuotaEditor quotas={mockQuotas} showPrice={false} />);

    expect(screen.queryByText(/预计费用/)).not.toBeInTheDocument();
  });
});
