// frontend/packages/studio/src/pages/settings/UserProfile/UserProfile.tsx

import React, { useState } from 'react';
import { Button, Input, Upload, Alert } from '@coze-studio/ui-components';
import { useAuthStore } from '@coze-studio/stores';
import './UserProfile.styles.css';

export const UserProfile: React.FC = () => {
  const { userInfo, updateUserInfo } = useAuthStore();
  const [isEditing, setIsEditing] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const [formData, setFormData] = useState({
    display_name: userInfo?.display_name || '',
    email: userInfo?.email || '',
    phone: '',
    bio: '',
  });

  const [passwordData, setPasswordData] = useState({
    current: '',
    new: '',
    confirm: '',
  });

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handlePasswordChange = (field: string, value: string) => {
    setPasswordData((prev) => ({ ...prev, [field]: value }));
  };

  const handleSaveProfile = () => {
    // TODO: 调用 API 保存个人信息
    updateUserInfo(formData);
    setMessage({ type: 'success', text: '个人信息已更新' });
    setIsEditing(false);
    setTimeout(() => setMessage(null), 3000);
  };

  const handleUploadAvatar = (fileList: any[]) => {
    if (fileList.length > 0) {
      const file = fileList[0];
      // TODO: 上传头像到服务器
      console.log('Upload avatar:', file);
    }
  };

  const handleChangePassword = () => {
    if (passwordData.new !== passwordData.confirm) {
      setMessage({ type: 'error', text: '两次输入的密码不一致' });
      setTimeout(() => setMessage(null), 3000);
      return;
    }

    // TODO: 调用 API 修改密码
    console.log('Change password:', passwordData);
    setMessage({ type: 'success', text: '密码已修改' });
    setPasswordData({ current: '', new: '', confirm: '' });
    setTimeout(() => setMessage(null), 3000);
  };

  if (!userInfo) {
    return <div className="user-profile-loading">加载中...</div>;
  }

  return (
    <div className="user-profile">
      <div className="profile-header">
        <h1>个人设置</h1>
        <p>管理您的个人信息和偏好设置</p>
      </div>

      {message && (
        <Alert
          type={message.type}
          message={message.text}
          closable
          onClose={() => setMessage(null)}
          className="profile-alert"
        />
      )}

      {/* 头像设置 */}
      <div className="profile-section">
        <h2>头像</h2>
        <div className="avatar-section">
          <div className="avatar-preview">
            {userInfo.avatar_url ? (
              <img src={userInfo.avatar_url} alt="头像" />
            ) : (
              <div className="avatar-placeholder">
                {userInfo.display_name?.[0] || userInfo.username[0]}
              </div>
            )}
          </div>
          <Upload
            accept="image/*"
            fileList={[]}
            onChange={handleUploadAvatar}
          >
            <Button>更换头像</Button>
          </Upload>
          <p className="upload-hint">支持 JPG、PNG 格式，大小不超过 2MB</p>
        </div>
      </div>

      {/* 基本信息 */}
      <div className="profile-section">
        <div className="section-header">
          <h2>基本信息</h2>
          {!isEditing && (
            <Button onClick={() => setIsEditing(true)}>编辑</Button>
          )}
        </div>

        <div className="profile-form">
          <div className="form-group">
            <label>用户名</label>
            <Input value={userInfo.username} disabled />
          </div>

          <div className="form-group">
            <label>显示名称</label>
            <Input
              value={formData.display_name}
              onChange={(e) => handleInputChange('display_name', e.target.value)}
              disabled={!isEditing}
              placeholder="请输入显示名称"
            />
          </div>

          <div className="form-group">
            <label>邮箱地址</label>
            <Input
              type="email"
              value={formData.email}
              onChange={(e) => handleInputChange('email', e.target.value)}
              disabled={!isEditing}
              placeholder="请输入邮箱地址"
            />
          </div>

          <div className="form-group">
            <label>手机号码</label>
            <Input
              type="tel"
              value={formData.phone}
              onChange={(e) => handleInputChange('phone', e.target.value)}
              disabled={!isEditing}
              placeholder="请输入手机号码"
            />
          </div>

          <div className="form-group">
            <label>个人简介</label>
            <textarea
              value={formData.bio}
              onChange={(e) => handleInputChange('bio', e.target.value)}
              disabled={!isEditing}
              placeholder="介绍一下自己..."
              rows={4}
              className="bio-textarea"
            />
          </div>

          {isEditing && (
            <div className="form-actions">
              <Button onClick={() => setIsEditing(false)}>取消</Button>
              <Button onClick={handleSaveProfile}>保存</Button>
            </div>
          )}
        </div>
      </div>

      {/* 修改密码 */}
      <div className="profile-section">
        <h2>修改密码</h2>
        <div className="password-form">
          <div className="form-group">
            <label>当前密码</label>
            <Input
              type="password"
              value={passwordData.current}
              onChange={(e) => handlePasswordChange('current', e.target.value)}
              placeholder="请输入当前密码"
            />
          </div>

          <div className="form-group">
            <label>新密码</label>
            <Input
              type="password"
              value={passwordData.new}
              onChange={(e) => handlePasswordChange('new', e.target.value)}
              placeholder="请输入新密码（至少8位）"
            />
          </div>

          <div className="form-group">
            <label>确认新密码</label>
            <Input
              type="password"
              value={passwordData.confirm}
              onChange={(e) => handlePasswordChange('confirm', e.target.value)}
              placeholder="请再次输入新密码"
            />
          </div>

          <Button onClick={handleChangePassword}>修改密码</Button>
        </div>
      </div>

      {/* 账号信息 */}
      <div className="profile-section">
        <h2>账号信息</h2>
        <div className="account-info">
          <div className="info-item">
            <span className="label">用户ID：</span>
            <span className="value">{userInfo.user_id}</span>
          </div>
          <div className="info-item">
            <span className="label">账号状态：</span>
            <span className={`value status ${userInfo.is_active ? 'active' : 'inactive'}`}>
              {userInfo.is_active ? '正常' : '已禁用'}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UserProfile;
