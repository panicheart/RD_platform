import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Button, Card, Row, Col, Typography, Space } from 'antd'
import { 
  ArrowRightOutlined, 
  TeamOutlined, 
  RocketOutlined,
  ExperimentOutlined,
  ApiOutlined,
  TrophyOutlined,
  HeartOutlined,
  StarOutlined,
  LoginOutlined,
} from '@ant-design/icons'
import { useAuth } from '@/hooks/useAuth'

const { Title, Text, Paragraph } = Typography

// ============ 导航栏 ============
const PortalNavbar = ({ isAuthenticated, onLogin, onWorkbench }: { 
  isAuthenticated: boolean, 
  onLogin: () => void, 
  onWorkbench: () => void 
}) => {
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
    <nav style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      zIndex: 1000,
      height: scrolled ? '70px' : '90px',
      background: scrolled ? 'rgba(255,255,255,0.95)' : 'transparent',
      backdropFilter: scrolled ? 'blur(10px)' : 'none',
      boxShadow: scrolled ? '0 2px 10px rgba(0,0,0,0.1)' : 'none',
      transition: 'all 0.3s ease',
    }}>
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 48px', height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        {/* Logo */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px', cursor: 'pointer' }} onClick={() => scrollTo('home')}>
          <div style={{ 
            width: '40px', 
            height: '40px', 
            borderRadius: '50%', 
            background: 'linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}>
            <span style={{ color: 'white', fontWeight: 'bold', fontSize: '18px' }}>M</span>
          </div>
          <div>
            <div style={{ fontWeight: 'bold', fontSize: '18px', color: scrolled ? '#1e293b' : 'white' }}>微波室</div>
            <div style={{ fontSize: '12px', color: scrolled ? '#64748b' : 'rgba(255,255,255,0.7)' }}>航天科工集团</div>
          </div>
        </div>

        {/* Navigation Links */}
        <div style={{ display: 'flex', gap: '32px' }}>
          {['about', 'services', 'projects', 'achievements', 'culture'].map((item) => (
            <a 
              key={item}
              onClick={() => scrollTo(item)}
              style={{ 
                color: scrolled ? '#475569' : 'rgba(255,255,255,0.9)',
                cursor: 'pointer',
                fontSize: '14px',
                fontWeight: 500,
                textDecoration: 'none',
                transition: 'color 0.3s',
              }}
            >
              {item === 'about' && '关于'}
              {item === 'services' && '技术'}
              {item === 'projects' && '产品'}
              {item === 'achievements' && '荣誉'}
              {item === 'culture' && '文化'}
            </a>
          ))}
        </div>

        {/* Auth Buttons */}
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

// ============ Hero 区 ============
const HeroSection = () => (
  <section id="home" style={{ 
    minHeight: '100vh', 
    background: 'linear-gradient(135deg, #0f172a 0%, #1e3a5f 50%, #0f172a 100%)',
    position: 'relative',
    overflow: 'hidden',
  }}>
    {/* 背景装饰 */}
    <div style={{ 
      position: 'absolute', 
      top: '10%', 
      right: '10%', 
      fontSize: '200px', 
      opacity: 0.05, 
      color: 'white',
      fontFamily: 'serif',
      transform: 'rotate(12deg)',
    }}>
      ∇ × E
    </div>
    
    <div style={{ 
      maxWidth: '1280px', 
      margin: '0 auto', 
      padding: '0 48px', 
      minHeight: '100vh', 
      display: 'flex', 
      alignItems: 'center',
      position: 'relative',
      zIndex: 1,
    }}>
      <div style={{ maxWidth: '600px' }}>
        {/* Badge */}
        <div style={{ 
          display: 'inline-flex', 
          alignItems: 'center', 
          gap: '8px',
          padding: '8px 16px',
          borderRadius: '20px',
          background: 'rgba(255,255,255,0.1)',
          marginBottom: '24px',
        }}>
          <RocketOutlined style={{ color: '#3b82f6' }} />
          <Text style={{ color: 'rgba(255,255,255,0.8)', fontSize: '14px' }}>
            航天科工集团 · 微波技术
          </Text>
        </div>
        
        <Title level={1} style={{ color: 'white', fontSize: '64px', fontWeight: 700, marginBottom: '16px', lineHeight: 1.2 }}>
          驾驭电磁波
          <br />
          <span style={{ color: '#3b82f6' }}>连接天地</span>
        </Title>
        
        <Paragraph style={{ color: 'rgba(255,255,255,0.7)', fontSize: '18px', marginBottom: '32px' }}>
          以麦克斯韦方程组为理论基础，专注于微波技术研发与应用，
          为航天事业提供先进的射频与微波解决方案。
        </Paragraph>
        
        <Space size="middle">
          <Button 
            type="primary"
            size="large"
            icon={<ArrowRightOutlined />}
            style={{ 
              borderRadius: '24px',
              height: '48px',
              paddingLeft: '24px',
              paddingRight: '24px',
              fontSize: '16px',
              background: 'linear-gradient(135deg, #3b82f6 0%, #2563eb 100%)',
              border: 'none',
            }}
          >
            探索技术
          </Button>
        </Space>
      </div>
    </div>
    
    {/* 底部渐变 */}
    <div style={{ 
      position: 'absolute', 
      bottom: 0, 
      left: 0, 
      right: 0, 
      height: '150px', 
      background: 'linear-gradient(to top, white, transparent)',
    }} />
  </section>
)

// ============ About 区 ============
const AboutSection = () => (
  <section id="about" style={{ padding: '120px 0', background: 'white', position: 'relative' }}>
    <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 48px' }}>
      <Row gutter={64} align="middle">
        <Col span={12}>
          <div style={{ marginBottom: '24px' }}>
            <div style={{ width: '32px', height: '2px', background: '#3b82f6', marginBottom: '8px' }} />
            <Text style={{ color: '#3b82f6', fontWeight: 600, fontSize: '14px', letterSpacing: '2px' }}>关于微波室</Text>
          </div>
          
          <Title level={2} style={{ fontSize: '42px', marginBottom: '24px' }}>
            以麦克斯韦方程为基石，
            <br />
            <span style={{ color: '#3b82f6' }}>开拓微波技术新境界</span>
          </Title>
          
          <Paragraph style={{ fontSize: '16px', color: '#64748b', lineHeight: 1.8, marginBottom: '32px' }}>
            微波室隶属于航天科工集团，是国内领先的微波技术研发中心。
            我们深耕微波领域三十余年，在雷达系统、通信设备、微波组件等方面积累了深厚的技术底蕴，
            为国家航天事业和国防建设提供了大量关键技术与产品。
          </Paragraph>
          
          {/* 特性卡片 */}
          <Row gutter={16}>
            {[
              { icon: <ApiOutlined />, title: '射频 expertise', desc: '覆盖1-40GHz全频段' },
              { icon: <RocketOutlined />, title: '高功率技术', desc: '千瓦级功率放大' },
              { icon: <ExperimentOutlined />, title: '精密测量', desc: '参数精确测试' },
            ].map((item, i) => (
              <Col span={8} key={i}>
                <div style={{ 
                  padding: '20px', 
                  borderRadius: '12px', 
                  background: '#f8fafc',
                  textAlign: 'center',
                }}>
                  <div style={{ fontSize: '32px', color: '#3b82f6', marginBottom: '12px' }}>{item.icon}</div>
                  <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>{item.title}</div>
                  <Text type="secondary" style={{ fontSize: '12px' }}>{item.desc}</Text>
                </div>
              </Col>
            ))}
          </Row>
        </Col>
        
        <Col span={12}>
          <div style={{ 
            borderRadius: '16px', 
            overflow: 'hidden',
            boxShadow: '0 25px 50px -12px rgba(0,0,0,0.25)',
            position: 'relative',
          }}>
            <img 
              src="https://images.unsplash.com/photo-1517976487492-5750f3195933?w=800" 
              alt="微波技术" 
              style={{ width: '100%', display: 'block' }}
            />
            <div style={{ 
              position: 'absolute', 
              bottom: 0, 
              left: 0, 
              right: 0, 
              padding: '24px',
              background: 'linear-gradient(to top, rgba(0,0,0,0.8), transparent)',
            }}>
              <Row gutter={24}>
                <Col><div style={{ color: 'white', fontSize: '28px', fontWeight: 'bold' }}>30+</div><Text style={{ color: 'rgba(255,255,255,0.7)', fontSize: '12px' }}>年技术积累</Text></Col>
                <Col><div style={{ color: 'white', fontSize: '28px', fontWeight: 'bold' }}>200+</div><Text style={{ color: 'rgba(255,255,255,0.7)', fontSize: '12px' }}>研发人员</Text></Col>
                <Col><div style={{ color: 'white', fontSize: '28px', fontWeight: 'bold' }}>40GHz</div><Text style={{ color: 'rgba(255,255,255,0.7)', fontSize: '12px' }}>最高频率</Text></Col>
              </Row>
            </div>
          </div>
        </Col>
      </Row>
    </div>
  </section>
)

// ============ Services 区 ============
const ServicesSection = () => {
  const services = [
    { icon: <ApiOutlined />, title: '雷达系统', desc: '相控阵雷达、合成孔径雷达系统研发', color: '#3b82f6' },
    { icon: <RocketOutlined />, title: '射频前端', desc: '功率放大器、混频器、滤波器设计', color: '#06b6d4' },
    { icon: <ExperimentOutlined />, title: '微波测试', desc: '暗室测试、天线测量、EMC检测', color: '#8b5cf6' },
    { icon: <TeamOutlined />, title: '卫星通信', desc: '星载微波、地面站系统集成', color: '#ec4899' },
  ]

  return (
    <section id="services" style={{ padding: '100px 0', background: '#f8fafc' }}>
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 48px' }}>
        <div style={{ textAlign: 'center', marginBottom: '48px' }}>
          <Text style={{ color: '#3b82f6', fontWeight: 600, fontSize: '14px', letterSpacing: '2px' }}>技术服务</Text>
          <Title level={2} style={{ marginTop: '8px' }}>提供全面的微波技术解决方案</Title>
        </div>
        
        <Row gutter={[24, 24]}>
          {services.map((service, index) => (
            <Col xs={24} sm={12} lg={6} key={index}>
              <Card 
                hoverable
                style={{ borderRadius: '16px', height: '100%' }}
                styles={{ body: { padding: '32px 24px', textAlign: 'center' } }}
              >
                <div style={{ 
                  fontSize: '48px', 
                  color: service.color, 
                  marginBottom: '20px',
                }}>
                  {service.icon}
                </div>
                <Title level={4} style={{ marginBottom: '8px' }}>{service.title}</Title>
                <Text type="secondary">{service.desc}</Text>
              </Card>
            </Col>
          ))}
        </Row>
      </div>
    </section>
  )
}

// ============ Projects 区 ============
const ProjectsSection = () => {
  const projects = [
    { title: '相控阵雷达系统', category: '雷达', image: 'https://images.unsplash.com/photo-1569517282132-25d22f4573e6?w=600' },
    { title: '卫星通信地面站', category: '通信', image: 'https://images.unsplash.com/photo-1516849841032-87cbac4d88f7?w=600' },
    { title: '微波暗室测试', category: '测试', image: 'https://images.unsplash.com/photo-1581092160562-40aa08e78837?w=600' },
  ]

  return (
    <section id="projects" style={{ padding: '100px 0', background: 'white' }}>
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 48px' }}>
        <div style={{ textAlign: 'center', marginBottom: '48px' }}>
          <Text style={{ color: '#3b82f6', fontWeight: 600, fontSize: '14px', letterSpacing: '2px' }}>产品展示</Text>
          <Title level={2} style={{ marginTop: '8px' }}>已交付的关键产品与系统</Title>
        </div>
        
        <Row gutter={[24, 24]}>
          {projects.map((project, index) => (
            <Col xs={24} lg={8} key={index}>
              <div style={{ 
                borderRadius: '16px', 
                overflow: 'hidden',
                boxShadow: '0 10px 30px rgba(0,0,0,0.1)',
              }}>
                <div style={{ height: '240px', overflow: 'hidden' }}>
                  <img 
                    src={project.image} 
                    alt={project.title}
                    style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                  />
                </div>
                <div style={{ padding: '20px' }}>
                  <Text type="secondary" style={{ fontSize: '12px' }}>{project.category}</Text>
                  <Title level={4} style={{ marginTop: '4px', marginBottom: 0 }}>{project.title}</Title>
                </div>
              </div>
            </Col>
          ))}
        </Row>
      </div>
    </section>
  )
}

// ============ Achievements 区 ============
const AchievementsSection = () => {
  const achievements = [
    { title: '国家科技进步奖', year: '2024', level: '一等奖' },
    { title: '集团技术创新奖', year: '2023', level: '一等奖' },
    { title: '军用微波技术突破', year: '2023', level: '重大项目' },
  ]

  return (
    <section id="achievements" style={{ padding: '100px 0', background: 'linear-gradient(135deg, #0f172a 0%, #1e293b 100%)' }}>
      <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 48px' }}>
        <div style={{ textAlign: 'center', marginBottom: '48px' }}>
          <Text style={{ color: '#3b82f6', fontWeight: 600, fontSize: '14px', letterSpacing: '2px' }}>荣誉成就</Text>
          <Title level={2} style={{ marginTop: '8px', color: 'white' }}>见证我们的技术实力</Title>
        </div>
        
        <Row gutter={[24, 24]}>
          {achievements.map((item, index) => (
            <Col xs={24} lg={8} key={index}>
              <Card 
                style={{ 
                  borderRadius: '16px', 
                  background: 'rgba(255,255,255,0.05)',
                  border: '1px solid rgba(255,255,255,0.1)',
                }}
              >
                <div style={{ textAlign: 'center' }}>
                  <TrophyOutlined style={{ fontSize: '48px', color: '#f59e0b', marginBottom: '16px' }} />
                  <Title level={4} style={{ color: 'white', marginBottom: '8px' }}>{item.title}</Title>
                  <Space>
                    <span style={{ color: '#f59e0b', fontSize: '12px' }}>{item.year}</span>
                    <span style={{ color: '#3b82f6', fontSize: '12px' }}>{item.level}</span>
                  </Space>
                </div>
              </Card>
            </Col>
          ))}
        </Row>
      </div>
    </section>
  )
}

// ============ Culture 区 ============
const CultureSection = () => (
  <section id="culture" style={{ padding: '100px 0', background: '#f8fafc' }}>
    <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 48px', textAlign: 'center' }}>
      <Text style={{ color: '#3b82f6', fontWeight: 600, fontSize: '14px', letterSpacing: '2px' }}>团队文化</Text>
      <Title level={2} style={{ marginTop: '8px', marginBottom: '48px' }}>打造卓越的微波研发团队</Title>
      
      <Row gutter={[48, 48]}>
        {[
          { icon: <StarOutlined />, title: '创新进取', desc: '鼓励技术创新，勇于探索前沿' },
          { icon: <TeamOutlined />, title: '协同合作', desc: '跨团队协作，共克技术难题' },
          { icon: <HeartOutlined />, title: '质量至上', desc: '严把质量关，追求卓越品质' },
        ].map((item, index) => (
          <Col xs={24} md={8} key={index}>
            <div style={{ padding: '32px', borderRadius: '16px', background: 'white', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
              <div style={{ fontSize: '40px', color: '#3b82f6', marginBottom: '16px' }}>{item.icon}</div>
              <Title level={4}>{item.title}</Title>
              <Text type="secondary">{item.desc}</Text>
            </div>
          </Col>
        ))}
      </Row>
    </div>
  </section>
)

// ============ Footer ============
const PortalFooter = () => (
  <footer style={{ background: '#0f172a', padding: '60px 0 30px', color: 'white' }}>
    <div style={{ maxWidth: '1280px', margin: '0 auto', padding: '0 48px' }}>
      <Row gutter={[48, 32]}>
        <Col xs={24} md={8}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '12px', marginBottom: '16px' }}>
            <div style={{ 
              width: '40px', 
              height: '40px', 
              borderRadius: '50%', 
              background: 'linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}>
              <span style={{ color: 'white', fontWeight: 'bold', fontSize: '18px' }}>M</span>
            </div>
            <div>
              <div style={{ fontWeight: 'bold', fontSize: '18px' }}>微波室</div>
              <div style={{ fontSize: '12px', color: 'rgba(255,255,255,0.6)' }}>航天科工集团</div>
            </div>
          </div>
          <Text style={{ color: 'rgba(255,255,255,0.6)' }}>
            致力于微波技术研发与应用，为航天事业提供先进的射频与微波解决方案。
          </Text>
        </Col>
        <Col xs={24} md={8}>
          <Title level={5} style={{ color: 'white', marginBottom: '16px' }}>快速链接</Title>
          <Space direction="vertical">
            <a style={{ color: 'rgba(255,255,255,0.6)' }}>关于我们</a>
            <a style={{ color: 'rgba(255,255,255,0.6)' }}>技术产品</a>
            <a style={{ color: 'rgba(255,255,255,0.6)' }}>团队文化</a>
          </Space>
        </Col>
        <Col xs={24} md={8}>
          <Title level={5} style={{ color: 'white', marginBottom: '16px' }}>联系方式</Title>
          <Text style={{ color: 'rgba(255,255,255,0.6)' }}>
            航天科工集团 · 微波工程部<br />
            internal@microwave.rd
          </Text>
        </Col>
      </Row>
      <div style={{ borderTop: '1px solid rgba(255,255,255,0.1)', marginTop: '40px', paddingTop: '24px', textAlign: 'center' }}>
        <Text style={{ color: 'rgba(255,255,255,0.4)' }}>
          © 2026 微波工程部 | 内部系统
        </Text>
      </div>
    </div>
  </footer>
)

// ============ 主页面 ============
export default function PortalPage() {
  const navigate = useNavigate()
  const { isAuthenticated } = useAuth()

  const handleLogin = () => {
    navigate('/login', { state: { from: '/workbench' } })
  }

  const handleWorkbench = () => {
    if (isAuthenticated) {
      navigate('/workbench')
    } else {
      navigate('/login', { state: { from: '/workbench' } })
    }
  }

  return (
    <div style={{ background: 'white' }}>
      <PortalNavbar 
        isAuthenticated={isAuthenticated} 
        onLogin={handleLogin} 
        onWorkbench={handleWorkbench} 
      />
      <HeroSection />
      <AboutSection />
      <ServicesSection />
      <ProjectsSection />
      <AchievementsSection />
      <CultureSection />
      <PortalFooter />
    </div>
  )
}
