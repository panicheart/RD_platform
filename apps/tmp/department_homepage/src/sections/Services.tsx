import { useState, useEffect, useRef } from 'react';
import { ArrowRight, Radio, Waves, Cpu } from 'lucide-react';

const services = [
  {
    id: 'anechoic',
    icon: Radio,
    title: '微波暗室测试',
    subtitle: '专业电磁兼容与天线测试',
    description: '配备国内领先的微波暗室设施，提供天线方向图测量、雷达散射截面(RCS)测试、电磁兼容性(EMC)测试等专业服务，测试频率覆盖1-40GHz全频段。',
    image: '/service-anechoic.jpg',
    features: ['天线方向图测试', 'RCS测量', 'EMC测试', '近场/远场转换'],
  },
  {
    id: 'radar',
    icon: Waves,
    title: '雷达系统工程',
    subtitle: '先进雷达技术研发',
    description: '专注于相控阵雷达、合成孔径雷达(SAR)、脉冲多普勒雷达等系统的研发与设计，为国防和民用领域提供高性能雷达解决方案。',
    image: '/service-radar.jpg',
    features: ['相控阵雷达', 'SAR成像', '信号处理', '目标识别'],
  },
  {
    id: 'circuit',
    icon: Cpu,
    title: '微波组件设计',
    subtitle: '高频电路与模块开发',
    description: '从事微波功率放大器、低噪声放大器、混频器、滤波器等核心组件的研制，产品广泛应用于航天、通信、雷达等领域。',
    image: '/service-circuit.jpg',
    features: ['功率放大器', '低噪声放大器', '混频器', '滤波器'],
  },
];

export default function Services() {
  const [activeTab, setActiveTab] = useState(0);
  const [isVisible, setIsVisible] = useState(false);
  const [isTransitioning, setIsTransitioning] = useState(false);
  const sectionRef = useRef<HTMLElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.unobserve(entry.target);
        }
      },
      { threshold: 0.2 }
    );

    if (sectionRef.current) {
      observer.observe(sectionRef.current);
    }

    return () => observer.disconnect();
  }, []);

  const handleTabChange = (index: number) => {
    if (index === activeTab || isTransitioning) return;
    
    setIsTransitioning(true);
    setTimeout(() => {
      setActiveTab(index);
      setTimeout(() => {
        setIsTransitioning(false);
      }, 100);
    }, 400);
  };

  const activeService = services[activeTab];

  return (
    <section 
      id="services" 
      ref={sectionRef}
      className="relative py-24 lg:py-32 bg-space-light overflow-hidden"
    >
      {/* 电磁波装饰 */}
      <div className="absolute top-10 left-10 opacity-5 text-space-black font-mono text-4xl">
        E = E₀e^(i(k·r-ωt))
      </div>
      
      <div className="max-w-7xl mx-auto px-6 sm:px-12">
        {/* Section Header */}
        <div className="text-center mb-16">
          <div 
            className={`inline-flex items-center gap-2 mb-4 transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'
            }`}
            style={{ transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
          >
            <div className="w-8 h-[2px] bg-space-blue" />
            <span className="text-space-blue font-semibold text-sm tracking-wider uppercase">
              核心技术
            </span>
            <div className="w-8 h-[2px] bg-space-blue" />
          </div>
          
          <h2 
            className={`font-display text-3xl sm:text-4xl lg:text-5xl font-bold text-space-black transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
            }`}
            style={{ transitionDelay: '100ms', transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
          >
            我们的<span className="text-space-blue">技术领域</span>
          </h2>
        </div>

        {/* Tab Navigation */}
        <div 
          className={`flex flex-wrap justify-center gap-2 sm:gap-4 mb-12 transition-all duration-700 ${
            isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
          }`}
          style={{ transitionDelay: '200ms', transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
        >
          {services.map((service, index) => (
            <button
              key={service.id}
              onClick={() => handleTabChange(index)}
              className={`relative px-6 py-3 rounded-full font-medium text-sm transition-all duration-300 ${
                activeTab === index 
                  ? 'bg-space-blue text-white shadow-glow' 
                  : 'bg-white text-space-gray hover:text-space-blue hover:bg-blue-50'
              }`}
            >
              <span className="flex items-center gap-2">
                <service.icon className="w-4 h-4" />
                {service.title}
              </span>
            </button>
          ))}
        </div>

        {/* Content Area */}
        <div 
          className={`grid grid-cols-1 lg:grid-cols-2 gap-8 lg:gap-12 items-center transition-all duration-700 ${
            isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-12'
          }`}
          style={{ 
            transitionDelay: '300ms', 
            transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)',
            perspective: '1200px'
          }}
        >
          {/* Image */}
          <div 
            className={`relative overflow-hidden rounded-2xl shadow-card transition-all duration-400 ${
              isTransitioning ? 'opacity-0 rotate-y-90' : 'opacity-100 rotate-y-0'
            }`}
            style={{ 
              transformStyle: 'preserve-3d',
              backfaceVisibility: 'hidden'
            }}
          >
            <img 
              src={activeService.image} 
              alt={activeService.title}
              className="w-full h-[300px] lg:h-[400px] object-cover transition-transform duration-700 hover:scale-105"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent" />
          </div>

          {/* Content */}
          <div 
            className={`space-y-6 transition-all duration-400 ${
              isTransitioning ? 'opacity-0 -rotate-y-90' : 'opacity-100 rotate-y-0'
            }`}
            style={{ 
              transformStyle: 'preserve-3d',
              backfaceVisibility: 'hidden'
            }}
          >
            <div>
              <h3 className="font-display text-2xl sm:text-3xl font-bold text-space-black mb-2">
                {activeService.subtitle}
              </h3>
              <p className="text-space-gray text-lg leading-relaxed">
                {activeService.description}
              </p>
            </div>

            {/* Features */}
            <div className="grid grid-cols-2 gap-4">
              {activeService.features.map((feature, index) => (
                <div 
                  key={feature}
                  className="flex items-center gap-3 p-3 bg-white rounded-lg shadow-sm"
                  style={{ 
                    animation: isTransitioning ? 'none' : `fadeInUp 0.5s ease ${index * 100}ms forwards`,
                    opacity: isTransitioning ? 0 : 1
                  }}
                >
                  <div className="w-2 h-2 rounded-full bg-space-blue" />
                  <span className="text-space-gray font-medium">{feature}</span>
                </div>
              ))}
            </div>

            {/* CTA */}
            <button className="group inline-flex items-center gap-2 text-space-blue font-semibold hover:gap-4 transition-all duration-300">
              了解更多
              <ArrowRight className="w-4 h-4 transition-transform duration-300 group-hover:translate-x-2" />
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}
