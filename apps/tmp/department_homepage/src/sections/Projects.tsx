import { useEffect, useRef, useState } from 'react';
import { ArrowUpRight, Waves } from 'lucide-react';

const projects = [
  {
    id: 1,
    title: '相控阵雷达系统',
    category: '雷达技术',
    image: '/project-phasedarray.jpg',
    description: '高分辨率电子扫描相控阵雷达',
    frequency: 'X波段',
  },
  {
    id: 2,
    title: '卫星通信地面站',
    category: '通信系统',
    image: '/project-groundstation.jpg',
    description: '大口径抛物面天线地面接收系统',
    frequency: 'C/Ku波段',
  },
  {
    id: 3,
    title: '微波功率放大器',
    category: '微波组件',
    image: '/project-amplifier.jpg',
    description: '高功率固态微波放大模块',
    frequency: 'S/C波段',
  },
  {
    id: 4,
    title: '电磁兼容测试系统',
    category: '测试设备',
    image: '/project-emc.jpg',
    description: '全频段EMC测试与认证平台',
    frequency: '1-18GHz',
  },
  {
    id: 5,
    title: '毫米波雷达传感器',
    category: '传感器',
    image: '/project-mmwave.jpg',
    description: '77GHz车载毫米波雷达模块',
    frequency: '77GHz',
  },
  {
    id: 6,
    title: '天线测试场',
    category: '测试设施',
    image: '/project-antennatest.jpg',
    description: '室外远场天线测试与校准',
    frequency: '全频段',
  },
];

export default function Projects() {
  const sectionRef = useRef<HTMLElement>(null);
  const [isVisible, setIsVisible] = useState(false);
  const [hoveredId, setHoveredId] = useState<number | null>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.unobserve(entry.target);
        }
      },
      { threshold: 0.1 }
    );

    if (sectionRef.current) {
      observer.observe(sectionRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <section 
      id="projects" 
      ref={sectionRef}
      className="relative py-24 lg:py-32 bg-white overflow-hidden"
    >
      {/* 波动方程装饰 */}
      <div className="absolute bottom-20 right-10 opacity-5 text-space-black font-mono text-3xl rotate-[-10deg]">
        ∇²E = μ₀ε₀∂²E/∂t²
      </div>

      <div className="max-w-7xl mx-auto px-6 sm:px-12">
        {/* Section Header */}
        <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between mb-16">
          <div>
            <div 
              className={`inline-flex items-center gap-2 mb-4 transition-all duration-700 ${
                isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'
              }`}
              style={{ transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
            >
              <div className="w-8 h-[2px] bg-space-blue" />
              <span className="text-space-blue font-semibold text-sm tracking-wider uppercase">
                产品展示
              </span>
            </div>
            
            <h2 
              className={`font-display text-3xl sm:text-4xl lg:text-5xl font-bold text-space-black transition-all duration-700 ${
                isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
              }`}
              style={{ transitionDelay: '100ms', transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
            >
              探索我们的<span className="text-space-blue">产品</span>
            </h2>
          </div>
          
          <p 
            className={`mt-4 sm:mt-0 text-space-gray max-w-md transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
            }`}
            style={{ transitionDelay: '200ms', transitionTimingFunction: 'cubic-bezier(0.4, 0, 0.2, 1)' }}
          >
            从雷达系统到微波组件，我们的产品覆盖微波技术的各个领域
          </p>
        </div>

        {/* Projects Grid */}
        <div 
          className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
          style={{ perspective: '1500px' }}
        >
          {projects.map((project, index) => (
            <div
              key={project.id}
              className={`group relative transition-all duration-800 ${
                isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-20'
              }`}
              style={{ 
                transitionDelay: `${300 + index * 100}ms`,
                transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)',
                transform: isVisible ? `translateY(${index % 3 === 1 ? '40px' : '0'})` : 'translateY(80px)'
              }}
              onMouseEnter={() => setHoveredId(project.id)}
              onMouseLeave={() => setHoveredId(null)}
            >
              <div 
                className={`relative overflow-hidden rounded-2xl bg-gray-100 cursor-pointer card-3d ${
                  hoveredId === project.id ? 'shadow-card-hover' : 'shadow-card'
                }`}
              >
                {/* Image */}
                <div className="relative h-[400px] overflow-hidden">
                  <img 
                    src={project.image} 
                    alt={project.title}
                    className={`w-full h-full object-cover transition-transform duration-700 ${
                      hoveredId === project.id ? 'scale-105' : 'scale-100'
                    }`}
                  />
                  
                  {/* Overlay */}
                  <div 
                    className={`absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent transition-opacity duration-500 ${
                      hoveredId === project.id ? 'opacity-100' : 'opacity-70'
                    }`}
                  />
                  
                  {/* Content */}
                  <div className="absolute bottom-0 left-0 right-0 p-6">
                    <div className="flex items-center gap-2 mb-3">
                      <span 
                        className={`inline-block px-3 py-1 bg-space-blue/90 text-white text-xs font-medium rounded-full transition-all duration-500 ${
                          hoveredId === project.id ? 'translate-y-0 opacity-100' : 'translate-y-4 opacity-0'
                        }`}
                      >
                        {project.category}
                      </span>
                      <span 
                        className={`inline-flex items-center gap-1 px-2 py-1 bg-white/20 text-white text-xs rounded-full transition-all duration-500 ${
                          hoveredId === project.id ? 'translate-y-0 opacity-100' : 'translate-y-4 opacity-0'
                        }`}
                      >
                        <Waves className="w-3 h-3" />
                        {project.frequency}
                      </span>
                    </div>
                    
                    <h3 className="font-display text-xl font-bold text-white mb-2">
                      {project.title}
                    </h3>
                    
                    <p 
                      className={`text-white/80 text-sm transition-all duration-500 ${
                        hoveredId === project.id ? 'translate-y-0 opacity-100' : 'translate-y-4 opacity-0'
                      }`}
                    >
                      {project.description}
                    </p>
                  </div>
                  
                  {/* Arrow Icon */}
                  <div 
                    className={`absolute top-4 right-4 w-10 h-10 rounded-full bg-white flex items-center justify-center transition-all duration-500 ${
                      hoveredId === project.id ? 'scale-100 opacity-100' : 'scale-75 opacity-0'
                    }`}
                  >
                    <ArrowUpRight className="w-5 h-5 text-space-blue" />
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
