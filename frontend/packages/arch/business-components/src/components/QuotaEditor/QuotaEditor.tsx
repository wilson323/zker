// frontend/packages/arch/business-components/src/components/QuotaEditor/QuotaEditor.tsx

import React, { useState, CSSProperties } from 'react';
import { useStyles } from './QuotaEditor.styles';
import { Input } from '@coze-studio/ui-components';
import { Slider } from '@coze-studio/ui-components';
import { QuotaResourceType, QuotaUnit } from '@coze-studio/api-client';

export interface QuotaLimit {
  resource_type: QuotaResourceType;
  limit: number;
  unit: QuotaUnit;
  is_soft_limit: boolean;
  overage_fee?: number;
}

export interface QuotaEditorProps {
  /**
   * 配额限制列表
   */
  quotas: QuotaLimit[];

  /**
   * 变更回调
   */
  onChange?: (quotas: QuotaLimit[]) => void;

  /**
   * 是否显示价格计算
   */
  showPrice?: boolean;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;
}

const resourceTypeLabels: Record<QuotaResourceType, string> = {
  [QuotaResourceType.BOT]: '机器人数量',
  [QuotaResourceType.USER]: '用户数量',
  [QuotaResourceType.API_CALL]: 'API调用次数',
  [QuotaResourceType.MESSAGE]: '消息数量',
  [QuotaResourceType.TOKEN]: 'Token使用量',
  [QuotaResourceType.STORAGE]: '存储空间(MB)',
  [QuotaResourceType.KNOWLEDGE_BASE]: '知识库数量',
  [QuotaResourceType.WORKFLOW]: '工作流数量',
};

const unitLabels: Record<QuotaUnit, string> = {
  [QuotaUnit.COUNT]: '个',
  [QuotaUnit.MB]: 'MB',
  [QuotaUnit.GB]: 'GB',
  [QuotaUnit.TIMES]: '次',
};

const defaultLimits: Record<QuotaResourceType, { max: number; step: number }> = {
  [QuotaResourceType.BOT]: { max: 1000, step: 10 },
  [QuotaResourceType.USER]: { max: 10000, step: 100 },
  [QuotaResourceType.API_CALL]: { max: 10000000, step: 100000 },
  [QuotaResourceType.MESSAGE]: { max: 10000000, step: 100000 },
  [QuotaResourceType.TOKEN]: { max: 100000000, step: 1000000 },
  [QuotaResourceType.STORAGE]: { max: 10240, step: 100 },
  [QuotaResourceType.KNOWLEDGE_BASE]: { max: 1000, step: 10 },
  [QuotaResourceType.WORKFLOW]: { max: 1000, step: 10 },
};

export const QuotaEditor: React.FC<QuotaEditorProps> = ({
  quotas,
  onChange,
  showPrice = true,
  className,
  style,
}) => {
  const classes = useStyles();
  const [localQuotas, setLocalQuotas] = useState<QuotaLimit[]>(quotas);

  React.useEffect(() => {
    setLocalQuotas(quotas);
  }, [quotas]);

  const handleChange = (index: number, field: keyof QuotaLimit, value: any) => {
    const updated = [...localQuotas];
    updated[index] = { ...updated[index], [field]: value };
    setLocalQuotas(updated);
    onChange?.(updated);
  };

  const calculatePrice = (quota: QuotaLimit): number => {
    // 简化的价格计算逻辑
    const basePrices: Record<QuotaResourceType, number> = {
      [QuotaResourceType.BOT]: 10,
      [QuotaResourceType.USER]: 5,
      [QuotaResourceType.API_CALL]: 0.0001,
      [QuotaResourceType.MESSAGE]: 0.0001,
      [QuotaResourceType.TOKEN]: 0.00001,
      [QuotaResourceType.STORAGE]: 0.01,
      [QuotaResourceType.KNOWLEDGE_BASE]: 1,
      [QuotaResourceType.WORKFLOW]: 1,
    };

    const basePrice = basePrices[quota.resource_type] || 0;
    return quota.limit * basePrice;
  };

  const getTotalPrice = (): number => {
    return localQuotas.reduce((sum, quota) => sum + calculatePrice(quota), 0);
  };

  return (
    <div className={`${classes.container} ${className || ''}`} style={style}>
      {localQuotas.map((quota, index) => {
        const { max, step } = defaultLimits[quota.resource_type];
        const price = calculatePrice(quota);

        return (
          <div key={quota.resource_type} className={classes.quotaItem}>
            <div className={classes.header}>
              <label className={classes.label}>
                {resourceTypeLabels[quota.resource_type]}
              </label>
              <span className={classes.unit}>{unitLabels[quota.unit]}</span>
            </div>

            <div className={classes.control}>
              <Slider
                value={quota.limit}
                min={0}
                max={max}
                step={step}
                onChange={(value) => handleChange(index, 'limit', value)}
                style={{ flex: 1 }}
              />
              <Input
                type="number"
                value={quota.limit}
                min={0}
                max={max}
                step={step}
                onChange={(e) =>
                  handleChange(index, 'limit', Number(e.target.value))
                }
                className={classes.input}
              />
            </div>

            <div className={classes.options}>
              <label className={classes.checkboxLabel}>
                <input
                  type="checkbox"
                  checked={quota.is_soft_limit}
                  onChange={(e) =>
                    handleChange(index, 'is_soft_limit', e.target.checked)
                  }
                />
                允许超限
              </label>

              {quota.is_soft_limit && (
                <div className={classes.overageFee}>
                  <span>超量单价：</span>
                  <Input
                    type="number"
                    value={quota.overage_fee || 0}
                    min={0}
                    step={0.01}
                    onChange={(e) =>
                      handleChange(index, 'overage_fee', Number(e.target.value))
                    }
                    className={classes.feeInput}
                  />
                  <span>元/{unitLabels[quota.unit]}</span>
                </div>
              )}
            </div>

            {showPrice && (
              <div className={classes.price}>预计费用：¥{price.toFixed(2)}</div>
            )}
          </div>
        );
      })}

      {showPrice && (
        <div className={classes.total}>
          <strong>总计：¥{getTotalPrice().toFixed(2)}/月</strong>
        </div>
      )}
    </div>
  );
};

export default QuotaEditor;
