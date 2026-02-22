import { Row, Col, Card, Typography } from 'antd'
import { ApiOutlined, RocketOutlined, ExperimentOutlined, TeamOutlined } from '@ant-design/icons'
import styles from '../styles.module.css'

const { Title, Text } = Typography

const services = [
  { icon: <ApiOutlined />, title: '雷达系统', desc: '相控阵雷达、合成孔径雷达系统研发', color: '#3b82f6' },
  { icon: <RocketOutlined />, title: '射频前端', desc: '功率放大器、混频器、滤波器设计', color: '#06b6d4' },
  { icon: <ExperimentOutlined />, title: '微波测试', desc: '暗室测试、天线测量、EMC检测', color: '#8b5cf6' },
  { icon: <TeamOutlined />, title: '卫星通信', desc: '星载微波、地面站系统集成', color: '#ec4899' },
]

export default function ServicesSection() {
  return (
    <section id="services" className={styles.servicesSection}>
      <div className={styles.sectionInner}>
        <div className={styles.sectionHeader}>
          <Text className={styles.sectionLabel}>技术服务</Text>
          <Title level={2} style={{ marginTop: '8px' }}>提供全面的微波技术解决方案</Title>
        </div>

        <Row gutter={[24, 24]}>
          {services.map((service, index) => (
            <Col xs={24} sm={12} lg={6} key={index}>
              <Card
                hoverable
                className={styles.serviceCard}
                styles={{ body: { padding: '32px 24px', textAlign: 'center' } }}
              >
                <div style={{ fontSize: '48px', color: service.color, marginBottom: '20px' }}>
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
