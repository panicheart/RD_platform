import { Row, Col, Typography } from 'antd'
import styles from '../styles.module.css'

const { Title, Text } = Typography

const projects = [
  { title: '相控阵雷达系统', category: '雷达', image: 'https://images.unsplash.com/photo-1569517282132-25d22f4573e6?w=600' },
  { title: '卫星通信地面站', category: '通信', image: 'https://images.unsplash.com/photo-1516849841032-87cbac4d88f7?w=600' },
  { title: '微波暗室测试', category: '测试', image: 'https://images.unsplash.com/photo-1581092160562-40aa08e78837?w=600' },
]

export default function ProjectsSection() {
  return (
    <section id="projects" className={styles.projectsSection}>
      <div className={styles.sectionInner}>
        <div className={styles.sectionHeader}>
          <Text className={styles.sectionLabel}>产品展示</Text>
          <Title level={2} style={{ marginTop: '8px' }}>已交付的关键产品与系统</Title>
        </div>

        <Row gutter={[24, 24]}>
          {projects.map((project, index) => (
            <Col xs={24} lg={8} key={index}>
              <div className={styles.projectCard}>
                <div className={styles.projectImage}>
                  <img src={project.image} alt={project.title} />
                </div>
                <div className={styles.projectInfo}>
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
