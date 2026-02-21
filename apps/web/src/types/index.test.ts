import { describe, it, expect } from 'vitest';
import type { User, Project, APIResponse } from '@types';

describe('Type Definitions', () => {
  it('should define User type correctly', () => {
    const user: User = {
      id: '01H',
      username: 'testuser',
      email: 'test@example.com',
      displayName: 'Test User',
      status: 'active',
      roles: [],
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    };

    expect(user.id).toBeDefined();
    expect(user.status).toBe('active');
  });

  it('should define Project type correctly', () => {
    const project: Project = {
      id: '01H',
      name: 'Test Project',
      code: 'TEST-001',
      category: 'new_product',
      status: 'in_progress',
      priority: 'high',
      classification: 'internal',
      ownerId: '01H',
      teamMembers: [],
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    };

    expect(project.category).toBe('new_product');
    expect(project.priority).toBe('high');
  });

  it('should define APIResponse type correctly', () => {
    const response: APIResponse = {
      code: 200,
      message: 'success',
      data: { id: '1' },
    };

    expect(response.code).toBe(200);
    expect(response.message).toBe('success');
  });
});
