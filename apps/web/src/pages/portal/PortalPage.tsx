import { useNavigate } from 'react-router-dom'
import { useAuth } from '@/hooks/useAuth'
import PortalNavbar from './components/PortalNavbar'
import HeroSection from './components/HeroSection'
import AboutSection from './components/AboutSection'
import ServicesSection from './components/ServicesSection'
import ProjectsSection from './components/ProjectsSection'
import AchievementsSection from './components/AchievementsSection'
import CultureSection from './components/CultureSection'
import PortalFooter from './components/PortalFooter'

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
