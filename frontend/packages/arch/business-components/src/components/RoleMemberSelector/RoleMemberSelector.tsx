// frontend/packages/arch/business-components/src/components/RoleMemberSelector/RoleMemberSelector.tsx

import React, { useState, CSSProperties } from 'react';
import { useStyles } from './RoleMemberSelector.styles';
import { Input } from '@coze-studio/ui-components';
import { RoleMemberDTO } from '@coze-studio/api-client';

export interface RoleMemberSelectorProps {
  /**
   * 当前成员列表
   */
  members: RoleMemberDTO[];

  /**
   * 可选择的用户列表
   */
  availableUsers: RoleMemberDTO[];

  /**
   * 是否多选
   */
  multiple?: boolean;

  /**
   * 是否禁用
   */
  disabled?: boolean;

  /**
   * 成员变更回调
   */
  onChange?: (members: RoleMemberDTO[]) => void;

  /**
   * 成员移除回调
   */
  onRemove?: (userId: string) => void;

  /**
   * 自定义样式类名
   */
  className?: string;

  /**
   * 自定义样式
   */
  style?: CSSProperties;
}

export const RoleMemberSelector: React.FC<RoleMemberSelectorProps> = ({
  members,
  availableUsers,
  multiple = true,
  disabled = false,
  onChange,
  onRemove,
  className,
  style,
}) => {
  const classes = useStyles();
  const [searchText, setSearchText] = useState('');
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = React.useRef<HTMLDivElement>(null);

  // 点击外部关闭下拉列表
  React.useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // 过滤用户列表
  const filteredUsers = availableUsers.filter(
    (user) =>
      !members.some((m) => m.user_id === user.user_id) &&
      (searchText === '' ||
        user.username.toLowerCase().includes(searchText.toLowerCase()) ||
        user.email.toLowerCase().includes(searchText.toLowerCase()) ||
        user.display_name?.toLowerCase().includes(searchText.toLowerCase()))
  );

  // 添加成员
  const handleAddMember = (user: RoleMemberDTO) => {
    if (!multiple) {
      onChange?.([user]);
    } else {
      onChange?.([...members, user]);
    }
    setIsOpen(false);
  };

  // 移除成员
  const handleRemoveMember = (userId: string) => {
    const updated = members.filter((m) => m.user_id !== userId);
    onChange?.(updated);
    onRemove?.(userId);
  };

  return (
    <div
      ref={containerRef}
      className={`${classes.container} ${className || ''}`}
      style={style}
    >
      <div className={classes.searchBox}>
        <Input
          placeholder="搜索用户..."
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          onFocus={() => !disabled && setIsOpen(true)}
          disabled={disabled}
        />
      </div>

      {isOpen && !disabled && (
        <div className={classes.dropdown}>
          {filteredUsers.length === 0 ? (
            <div className={classes.empty}>未找到用户</div>
          ) : (
            filteredUsers.map((user) => (
              <div
                key={user.user_id}
                className={classes.userItem}
                onClick={() => handleAddMember(user)}
              >
                {user.avatar_url && (
                  <img
                    src={user.avatar_url}
                    alt={user.username}
                    className={classes.avatar}
                  />
                )}
                <div className={classes.userInfo}>
                  <div className={classes.username}>
                    {user.display_name || user.username}
                  </div>
                  <div className={classes.email}>{user.email}</div>
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {members.length > 0 && (
        <div className={classes.memberList}>
          {members.map((member) => (
            <div key={member.user_id} className={classes.memberItem}>
              {member.avatar_url && (
                <img
                  src={member.avatar_url}
                  alt={member.username}
                  className={classes.avatar}
                />
              )}
              <div className={classes.userInfo}>
                <div className={classes.username}>
                  {member.display_name || member.username}
                </div>
                <div className={classes.email}>{member.email}</div>
              </div>
              {!disabled && (
                <button
                  className={classes.removeButton}
                  onClick={() => handleRemoveMember(member.user_id)}
                >
                  ✕
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default RoleMemberSelector;
