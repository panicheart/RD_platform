import { Row, Col, Card, Space, Typography } from 'antd'
import { TrophyOutlined } from '@ant-design/icons'
import styles from '../styles.module.css'

const { Title, Text } = Typography

const achievements = [
  { title: '国家科技进步奖', year: '2024', level: '一等奖' },
  { title: '集团技术创新奖', year: '2023', level: '一等奖' },
  { title: '军用微波技术突破', year: '2023', level: '重大项目' },
]

export default function AchievementsSection() {
  return (
    <section id="achievements" className={styles.achievementsSection}>
      <div className={styles.sectionInner}>
        <div className={styles.sectionHeader}>
          <Text className={styles.sectionLabel}>荣誉成就</Text>
          <Title level={2} style={{ marginTop: '8px', color: 'white' }}>见证我们的技术实力</Title>
        </div>

        <Row gutter={[24, 24]}>
          {achievements.map((item, index) => (
            <Col xs={24} lg={8} key={index}>
              <Card className={styles.achievementCard}>
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
