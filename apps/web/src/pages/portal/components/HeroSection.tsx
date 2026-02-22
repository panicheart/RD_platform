import { Button, Space, Typography } from 'antd'
import { RocketOutlined, ArrowRightOutlined } from '@ant-design/icons'
import styles from '../styles.module.css'

const { Title, Paragraph, Text } = Typography

export default function HeroSection() {
  return (
    <section id="home" className={styles.hero}>
      <div className={styles.heroBgDecor}>∇ × E</div>

      <div className={styles.heroInner}>
        <div className={styles.heroContent}>
          <div className={styles.heroBadge}>
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

      <div className={styles.heroGradientBottom} />
    </section>
  )
}
