// frontend/packages/arch/business-components/src/components/SubscriptionSelector/SubscriptionSelector.tsx

import React from 'react';
import { Select } from '@coze-studio/ui-components';
import { useStyles } from './SubscriptionSelector.styles';

export type SubscriptionTier = 'free' | 'pro' | 'enterprise';

export interface SubscriptionTierOption {
  value: SubscriptionTier;
  label: string;
  description: string;
  features: string[];
  price?: string;
}

export interface SubscriptionSelectorProps {
  value?: SubscriptionTier;
  onChange?: (tier: SubscriptionTier) => void;
  disabled?: boolean;
  className?: string;
}

const defaultOptions: SubscriptionTierOption[] = [
  {
    value: 'free',
    label: '免费版',
    description: '适合个人用户和小团队',
    features: ['5个机器人', '1000次API调用/月', '基础技术支持'],
    price: '¥0/月',
  },
  {
    value: 'pro',
    label: '专业版',
    description: '适合成长型团队',
    features: ['50个机器人', '50000次API调用/月', '优先技术支持', '高级分析功能'],
    price: '¥999/月',
  },
  {
    value: 'enterprise',
    label: '企业版',
    description: '适合大型企业',
    features: ['无限机器人', '无限API调用', '24/7专属支持', '定制开发', '私有化部署'],
    price: '联系销售',
  },
];

export const SubscriptionSelector: React.FC<SubscriptionSelectorProps> = ({
  value,
  onChange,
  disabled = false,
  className,
}) => {
  const classes = useStyles();

  const options = defaultOptions.map(option => ({
    label: option.label,
    value: option.value,
  }));

  return (
    <div className={`${classes.container} ${className || ''}`}>
      <label className={classes.label}>订阅等级</label>
      <Select
        value={value}
        onChange={onChange}
        options={options}
        disabled={disabled}
        placeholder="请选择订阅等级"
      />

      {value && (
        <div className={classes.details}>
          {(() => {
            const selected = defaultOptions.find(opt => opt.value === value);
            if (!selected) return null;

            return (
              <>
                <div className={classes.header}>
                  <h3 className={classes.title}>{selected.label}</h3>
                  {selected.price && (
                    <span className={classes.price}>{selected.price}</span>
                  )}
                </div>
                <p className={classes.description}>{selected.description}</p>
                <ul className={classes.features}>
                  {selected.features.map((feature, index) => (
                    <li key={index}>{feature}</li>
                  ))}
                </ul>
              </>
            );
          })()}
        </div>
      )}
    </div>
  );
};
