import { useState, useEffect, createContext, useContext } from 'react'
import { api } from '@/utils/api'
import type { User } from '@/types'

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
  refreshUser: () => Promise<User>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    const token = localStorage.getItem('token') || localStorage.getItem('access_token')
    if (!token) {
      setIsLoading(false)
      return
    }

    try {
      const response = await api.get('/users/me')
      const userData = response.data?.data || response.data
      setUser(userData)
    } catch (error) {
      localStorage.removeItem('token')
      localStorage.removeItem('access_token')
      api.defaults.headers.common['Authorization'] = ''
    } finally {
      setIsLoading(false)
    }
  }

  const login = async (username: string, password: string) => {
    // Demo mode: simulate login without backend
    // In production, this would call: await api.post('/auth/login', { username, password })
    
    if (!username || !password) {
      throw new Error('请输入用户名和密码')
    }
    
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // Demo credentials
    const demoUser: User = {
      id: '1',
      username: username,
      email: `${username}@example.com`,
      display_name: username === 'admin' ? '管理员' : username,
      displayName: username === 'admin' ? '管理员' : username,
      role: username === 'admin' ? 'admin' : 'designer',
      roles: [{ id: '1', name: username === 'admin' ? '管理员' : '工程师', code: username === 'admin' ? 'admin' : 'designer', permissions: [] }],
      status: 'active',
      isActive: true,
      title: '工程师',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    }
    
    const demoToken = 'demo-token-' + Date.now()
    
    localStorage.setItem('token', demoToken)
    api.defaults.headers.common['Authorization'] = `Bearer ${demoToken}`
    setUser(demoUser)
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout')
    } catch {
      // Ignore errors when backend is not available
    } finally {
      localStorage.removeItem('token')
      localStorage.removeItem('access_token')
      delete api.defaults.headers.common['Authorization']
      setUser(null)
    }
  }

  const refreshUser = async (): Promise<User> => {
    try {
      const response = await api.get('/users/me')
      const userData = response.data?.data || response.data
      setUser(userData)
      return userData
    } catch (error) {
      console.error('Failed to refresh user:', error)
      throw error
    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
