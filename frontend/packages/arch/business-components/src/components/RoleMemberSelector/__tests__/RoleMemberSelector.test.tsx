// frontend/packages/arch/business-components/src/components/RoleMemberSelector/__tests__/RoleMemberSelector.test.tsx

import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { RoleMemberSelector } from '../RoleMemberSelector';
import { RoleMemberDTO } from '@coze-studio/api-client';

const mockMembers: RoleMemberDTO[] = [
  {
    user_id: '1',
    username: 'john',
    email: 'john@example.com',
    display_name: 'John Doe',
    joined_at: '2024-01-01T00:00:00Z',
  },
];

const mockAvailableUsers: RoleMemberDTO[] = [
  {
    user_id: '2',
    username: 'jane',
    email: 'jane@example.com',
    display_name: 'Jane Smith',
    joined_at: '2024-01-01T00:00:00Z',
  },
  {
    user_id: '3',
    username: 'bob',
    email: 'bob@example.com',
    display_name: 'Bob Johnson',
    joined_at: '2024-01-01T00:00:00Z',
  },
];

describe('RoleMemberSelector', () => {
  it('should render search input', () => {
    render(
      <RoleMemberSelector
        members={mockMembers}
        availableUsers={mockAvailableUsers}
      />
    );

    expect(screen.getByPlaceholderText('搜索用户...')).toBeInTheDocument();
  });

  it('should render current members list', () => {
    render(
      <RoleMemberSelector
        members={mockMembers}
        availableUsers={mockAvailableUsers}
      />
    );

    expect(screen.getByText('John Doe')).toBeInTheDocument();
    expect(screen.getByText('john@example.com')).toBeInTheDocument();
  });

  it('should show dropdown when input is focused', () => {
    render(
      <RoleMemberSelector
        members={[]}
        availableUsers={mockAvailableUsers}
      />
    );

    const input = screen.getByPlaceholderText('搜索用户...');
    fireEvent.focus(input);

    // Dropdown should be visible after focus
    expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    expect(screen.getByText('Bob Johnson')).toBeInTheDocument();
  });

  it('should filter users by search text', () => {
    render(
      <RoleMemberSelector
        members={[]}
        availableUsers={mockAvailableUsers}
      />
    );

    const input = screen.getByPlaceholderText('搜索用户...');
    fireEvent.focus(input);
    fireEvent.change(input, { target: { value: 'jane' } });

    expect(screen.getByText('Jane Smith')).toBeInTheDocument();
    expect(screen.queryByText('Bob Johnson')).not.toBeInTheDocument();
  });

  it('should add member when user is clicked', () => {
    const handleChange = jest.fn();
    render(
      <RoleMemberSelector
        members={[]}
        availableUsers={mockAvailableUsers}
        onChange={handleChange}
      />
    );

    const input = screen.getByPlaceholderText('搜索用户...');
    fireEvent.focus(input);

    const userItem = screen.getByText('Jane Smith');
    fireEvent.click(userItem);

    expect(handleChange).toHaveBeenCalledWith([
      expect.objectContaining({
        user_id: '2',
      }),
    ]);
  });

  it('should remove member when remove button is clicked', () => {
    const handleChange = jest.fn();
    render(
      <RoleMemberSelector
        members={mockMembers}
        availableUsers={mockAvailableUsers}
        onChange={handleChange}
      />
    );

    const removeButton = screen.getByText('✕');
    fireEvent.click(removeButton);

    expect(handleChange).toHaveBeenCalledWith([]);
  });

  it('should not remove members when disabled', () => {
    render(
      <RoleMemberSelector
        members={mockMembers}
        availableUsers={mockAvailableUsers}
        disabled={true}
      />
    );

    const removeButton = screen.queryByText('✕');
    expect(removeButton).not.toBeInTheDocument();
  });

  it('should show empty state when no users match search', () => {
    render(
      <RoleMemberSelector
        members={[]}
        availableUsers={mockAvailableUsers}
      />
    );

    const input = screen.getByPlaceholderText('搜索用户...');
    fireEvent.focus(input);
    fireEvent.change(input, { target: { value: 'nonexistent' } });

    expect(screen.getByText('未找到用户')).toBeInTheDocument();
  });
});
