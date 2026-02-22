import { Row, Col, Space, Typography } from 'antd'
import styles from '../styles.module.css'

const { Title, Text } = Typography

export default function PortalFooter() {
  return (
    <footer className={styles.footer}>
      <div className={styles.sectionInner}>
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
                <div style={{ fontWeight: 'bold', fontSize: '18px', color: 'white' }}>微波室</div>
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
        <div className={styles.footerDivider}>
          <Text style={{ color: 'rgba(255,255,255,0.4)' }}>
            © 2026 微波工程部 | 内部系统
          </Text>
        </div>
      </div>
    </footer>
  )
}
