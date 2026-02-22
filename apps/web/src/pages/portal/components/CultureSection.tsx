import { Row, Col, Typography } from 'antd'
import { StarOutlined, TeamOutlined, HeartOutlined } from '@ant-design/icons'
import styles from '../styles.module.css'

const { Title, Text } = Typography

const values = [
  { icon: <StarOutlined />, title: '创新进取', desc: '鼓励技术创新，勇于探索前沿' },
  { icon: <TeamOutlined />, title: '协同合作', desc: '跨团队协作，共克技术难题' },
  { icon: <HeartOutlined />, title: '质量至上', desc: '严把质量关，追求卓越品质' },
]

export default function CultureSection() {
  return (
    <section id="culture" className={styles.cultureSection}>
      <div className={styles.sectionInner} style={{ textAlign: 'center' }}>
        <Text className={styles.sectionLabel}>团队文化</Text>
        <Title level={2} style={{ marginTop: '8px', marginBottom: '48px' }}>打造卓越的微波研发团队</Title>

        <Row gutter={[48, 48]}>
          {values.map((item, index) => (
            <Col xs={24} md={8} key={index}>
              <div className={styles.cultureCard}>
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
}
