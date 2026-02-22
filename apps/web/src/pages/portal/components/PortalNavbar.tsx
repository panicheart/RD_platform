import { useState, useEffect } from 'react'
import { Button, Space } from 'antd'
import { LoginOutlined } from '@ant-design/icons'
import styles from '../styles.module.css'

interface PortalNavbarProps {
  isAuthenticated: boolean
  onLogin: () => void
  onWorkbench: () => void
}

const navItems = [
  { id: 'about', label: '关于' },
  { id: 'services', label: '技术' },
  { id: 'projects', label: '产品' },
  { id: 'achievements', label: '荣誉' },
  { id: 'culture', label: '文化' },
]

export default function PortalNavbar({ isAuthenticated, onLogin, onWorkbench }: PortalNavbarProps) {
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 50)
    window.addEventListener('scroll', handleScroll)
    return () => window.removeEventListener('scroll', handleScroll)
  }, [])

  const scrollTo = (id: string) => {
    const el = document.getElementById(id)
    if (el) el.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <nav
      className={`${styles.navbar} ${scrolled ? styles.navbarScrolled : styles.navbarTransparent}`}
      style={{ height: scrolled ? '70px' : '90px' }}
    >
      <div className={styles.navInner}>
        <div className={styles.logo} onClick={() => scrollTo('home')}>
          <div className={styles.logoIcon}>M</div>
          <div>
            <div style={{ fontWeight: 'bold', fontSize: '18px', color: scrolled ? '#1e293b' : 'white' }}>微波室</div>
            <div style={{ fontSize: '12px', color: scrolled ? '#64748b' : 'rgba(255,255,255,0.7)' }}>航天科工集团</div>
          </div>
        </div>

        <div className={styles.navLinks}>
          {navItems.map((item) => (
            <a
              key={item.id}
              onClick={() => scrollTo(item.id)}
              className={styles.navLink}
              style={{ color: scrolled ? '#475569' : 'rgba(255,255,255,0.9)' }}
            >
              {item.label}
            </a>
          ))}
        </div>

        <Space>
          {isAuthenticated ? (
            <Button
              type="primary"
              onClick={onWorkbench}
              style={{
                borderRadius: '20px',
                background: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
                border: 'none',
              }}
            >
              进入工作台
            </Button>
          ) : (
            <>
              <Button
                onClick={onWorkbench}
                style={{
                  borderRadius: '20px',
                  borderColor: scrolled ? '#3b82f6' : 'rgba(255,255,255,0.5)',
                  color: scrolled ? '#3b82f6' : 'white',
                }}
              >
                工作台
              </Button>
              <Button
                type="primary"
                onClick={onLogin}
                icon={<LoginOutlined />}
                style={{
                  borderRadius: '20px',
                  background: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
                  border: 'none',
                }}
              >
                登录
              </Button>
            </>
          )}
        </Space>
      </div>
    </nav>
  )
}
