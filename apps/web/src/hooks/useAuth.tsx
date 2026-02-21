import { useState, useEffect, createContext, useContext } from 'react'
import { api } from '@/utils/api'

interface User {
  id: string
  username: string
  display_name: string
  email?: string
  avatar_url?: string
  role: string
  team?: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (username: string, password: string) => Promise<void>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    checkAuth()
  }, [])

  const checkAuth = async () => {
    const token = localStorage.getItem('token')
    if (!token) {
      setIsLoading(false)
      return
    }

    try {
      const response = await api.get('/users/me')
      setUser(response.data.data)
    } catch (error) {
      localStorage.removeItem('token')
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
      display_name: username === 'admin' ? '管理员' : username,
      role: username === 'admin' ? 'admin' : 'designer',
      team: '软件',
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
      delete api.defaults.headers.common['Authorization']
      setUser(null)
    }
  }

  const refreshUser = async () => {
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
