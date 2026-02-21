import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import type { User, Role } from '../types';

interface AuthState {
  // 用户信息
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  
  // 登录/登出方法
  setAuth: (user: User, token: string, refreshToken?: string) => void;
  clearAuth: () => void;
  updateUser: (user: Partial<User>) => void;
  
  // Token管理
  setToken: (token: string) => void;
  setRefreshToken: (token: string) => void;
  
  // 权限检查
  hasPermission: (permission: string) => boolean;
  hasRole: (roleCode: string) => boolean;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      // 初始状态
      user: null,
      token: null,
      refreshToken: null,
      isAuthenticated: false,

      // 设置认证信息（登录成功时调用）
      setAuth: (user, token, refreshToken) => {
        localStorage.setItem('access_token', token);
        if (refreshToken) {
          localStorage.setItem('refresh_token', refreshToken);
        }
        set({ 
          user, 
          token, 
          refreshToken: refreshToken || null,
          isAuthenticated: true 
        });
      },

      // 清除认证信息（登出时调用）
      clearAuth: () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        set({ 
          user: null, 
          token: null, 
          refreshToken: null,
          isAuthenticated: false 
        });
      },

      // 更新用户信息
      updateUser: (userData) => {
        const currentUser = get().user;
        if (currentUser) {
          set({ user: { ...currentUser, ...userData } });
        }
      },

      // 设置Token
      setToken: (token) => {
        localStorage.setItem('access_token', token);
        set({ token });
      },

      // 设置刷新Token
      setRefreshToken: (refreshToken) => {
        localStorage.setItem('refresh_token', refreshToken);
        set({ refreshToken });
      },

      // 检查是否有指定权限
      hasPermission: (permission: string) => {
        const user = get().user;
        if (!user) return false;
        
        // 超级管理员拥有所有权限
        if (user.roles.some((role: Role) => role.code === 'super_admin')) {
          return true;
        }
        
        // 检查用户角色中的权限
        return user.roles.some((role: Role) => 
          role.permissions.includes(permission) || role.permissions.includes('*')
        );
      },

      // 检查是否有指定角色
      hasRole: (roleCode: string) => {
        const user = get().user;
        if (!user) return false;
        return user.roles.some((role: Role) => role.code === roleCode);
      },
    }),
    {
      name: 'rdp-auth-storage',
      storage: createJSONStorage(() => localStorage),
      partialize: (state) => ({ 
        user: state.user, 
        token: state.token,
        refreshToken: state.refreshToken,
        isAuthenticated: state.isAuthenticated 
      }),
    }
  )
);

// 导出类型
export type { AuthState };
