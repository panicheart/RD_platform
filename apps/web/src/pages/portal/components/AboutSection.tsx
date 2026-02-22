import { Row, Col, Typography } from 'antd'
import { ApiOutlined, RocketOutlined, ExperimentOutlined } from '@ant-design/icons'
import styles from '../styles.module.css'

const { Title, Text, Paragraph } = Typography

const features = [
  { icon: <ApiOutlined />, title: '射频 expertise', desc: '覆盖1-40GHz全频段' },
  { icon: <RocketOutlined />, title: '高功率技术', desc: '千瓦级功率放大' },
  { icon: <ExperimentOutlined />, title: '精密测量', desc: '参数精确测试' },
]

const stats = [
  { value: '30+', label: '年技术积累' },
  { value: '200+', label: '研发人员' },
  { value: '40GHz', label: '最高频率' },
]

export default function AboutSection() {
  return (
    <section id="about" className={styles.aboutSection}>
      <div className={styles.sectionInner}>
        <Row gutter={64} align="middle">
          <Col span={12}>
            <div style={{ marginBottom: '24px' }}>
              <div className={styles.sectionDivider} />
              <Text className={styles.sectionLabel}>关于微波室</Text>
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

            <Row gutter={16}>
              {features.map((item, i) => (
                <Col span={8} key={i}>
                  <div className={styles.featureCard}>
                    <div style={{ fontSize: '32px', color: '#3b82f6', marginBottom: '12px' }}>{item.icon}</div>
                    <div style={{ fontWeight: 'bold', marginBottom: '4px' }}>{item.title}</div>
                    <Text type="secondary" style={{ fontSize: '12px' }}>{item.desc}</Text>
                  </div>
                </Col>
              ))}
            </Row>
          </Col>

          <Col span={12}>
            <div className={styles.aboutImageWrapper}>
              <img
                src="https://images.unsplash.com/photo-1517976487492-5750f3195933?w=800"
                alt="微波技术"
                className={styles.aboutImage}
              />
              <div className={styles.aboutImageOverlay}>
                <Row gutter={24}>
                  {stats.map((item, i) => (
                    <Col key={i}>
                      <div style={{ color: 'white', fontSize: '28px', fontWeight: 'bold' }}>{item.value}</div>
                      <Text style={{ color: 'rgba(255,255,255,0.7)', fontSize: '12px' }}>{item.label}</Text>
                    </Col>
                  ))}
                </Row>
              </div>
            </div>
          </Col>
        </Row>
      </div>
    </section>
  )
}
