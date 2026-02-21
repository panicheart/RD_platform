import { useEffect, useRef, useState } from 'react';
import { ArrowRight, Radio, Zap, Target } from 'lucide-react';
import { Button } from '@/components/ui/button';

const features = [
  {
    icon: Radio,
    title: '射频 expertise',
    description: '覆盖1-40GHz全频段微波技术',
  },
  {
    icon: Zap,
    title: '高功率技术',
    description: '千瓦级微波功率放大与传输',
  },
  {
    icon: Target,
    title: '精密测量',
    description: '微波参数精确测试与校准',
  },
];

export default function About() {
  const sectionRef = useRef<HTMLElement>(null);
  const [isVisible, setIsVisible] = useState(false);

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

  return (
    <section 
      id="about" 
      ref={sectionRef}
      className="relative py-24 lg:py-32 bg-white overflow-hidden"
    >
      {/* Background Decoration */}
      <div className="absolute top-0 right-0 w-1/2 h-full bg-gradient-to-l from-gray-50 to-transparent opacity-50" />
      
      {/* 麦克斯韦方程装饰 */}
      <div className="absolute top-20 right-20 opacity-5 text-space-black font-mono text-6xl rotate-12">
        ∇ × E
      </div>
      
      <div className="relative z-10 max-w-7xl mx-auto px-6 sm:px-12">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-center">
          {/* Left Content */}
          <div className="space-y-8">
            {/* Section Label */}
            <div 
              className={`inline-flex items-center gap-2 transition-all duration-700 ${
                isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'
              }`}
              style={{ transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
            >
              <div className="w-8 h-[2px] bg-space-blue" />
              <span className="text-space-blue font-semibold text-sm tracking-wider uppercase">
                关于微波室
              </span>
            </div>

            {/* Title */}
            <h2 
              className={`font-display text-3xl sm:text-4xl lg:text-5xl font-bold text-space-black leading-tight transition-all duration-700 ${
                isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
              }`}
              style={{ transitionDelay: '100ms', transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
            >
              以麦克斯韦方程为基石，
              <br />
              <span className="text-space-blue">开拓微波技术新境界</span>
            </h2>

            {/* Description */}
            <p 
              className={`text-space-gray text-lg leading-relaxed max-w-xl transition-all duration-700 ${
                isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
              }`}
              style={{ transitionDelay: '200ms', transitionTimingFunction: 'cubic-bezier(0.4, 0, 0.2, 1)' }}
            >
              微波室隶属于航天科工集团，是国内领先的微波技术研发中心。
              我们深耕微波领域三十余年，在雷达系统、通信设备、微波组件等方面积累了深厚的技术底蕴，
              为国家航天事业和国防建设提供了大量关键技术与产品。
            </p>

            {/* Features */}
            <div 
              className={`grid grid-cols-1 sm:grid-cols-3 gap-6 pt-4 transition-all duration-700 ${
                isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
              }`}
              style={{ transitionDelay: '300ms', transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
            >
              {features.map((feature, index) => (
                <div 
                  key={feature.title}
                  className="group p-4 rounded-xl bg-gray-50 hover:bg-blue-50 transition-all duration-300"
                  style={{ transitionDelay: `${400 + index * 100}ms` }}
                >
                  <feature.icon className="w-8 h-8 text-space-blue mb-3 transition-transform duration-300 group-hover:scale-110" />
                  <h3 className="font-display font-bold text-space-black mb-1">
                    {feature.title}
                  </h3>
                  <p className="text-sm text-space-gray">
                    {feature.description}
                  </p>
                </div>
              ))}
            </div>

            {/* CTA Button */}
            <div 
              className={`pt-4 transition-all duration-700 ${
                isVisible ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-8'
              }`}
              style={{ transitionDelay: '600ms', transitionTimingFunction: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)' }}
            >
              <Button 
                variant="outline"
                size="lg"
                className="border-space-black text-space-black hover:bg-space-black hover:text-white rounded-full px-8 group transition-all duration-300"
              >
                了解更多
                <ArrowRight className="ml-2 w-4 h-4 transition-transform duration-300 group-hover:translate-x-2" />
              </Button>
            </div>
          </div>

          {/* Right Content - Image */}
          <div 
            className={`relative transition-all duration-1000 ${
              isVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-95'
            }`}
            style={{ transitionDelay: '200ms', transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
          >
            {/* Orbit Circle Decoration */}
            <div 
              className={`absolute -top-10 -right-10 w-32 h-32 rounded-full bg-gradient-to-br from-blue-200 to-transparent opacity-60 transition-all duration-800 ${
                isVisible ? 'scale-100 rotate-180' : 'scale-0 rotate-0'
              }`}
              style={{ 
                transitionDelay: '800ms',
                transitionTimingFunction: 'cubic-bezier(0.34, 1.56, 0.64, 1)'
              }}
            />
            
            {/* Main Image */}
            <div className="relative overflow-hidden rounded-2xl shadow-2xl">
              <div 
                className={`absolute inset-0 bg-space-blue transition-all duration-1000 ${
                  isVisible ? 'clip-path-full' : 'clip-path-zero'
                }`}
                style={{ 
                  clipPath: isVisible ? 'circle(100% at 50% 50%)' : 'circle(0% at 50% 50%)',
                  transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)'
                }}
              />
              <img 
                src="/about-antenna.jpg" 
                alt="微波天线阵列" 
                className="w-full h-auto object-cover transition-transform duration-700 hover:scale-105"
              />
              
              {/* Overlay Stats */}
              <div className="absolute bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-black/80 to-transparent">
                <div className="flex items-center gap-8">
                  <div>
                    <div className="text-3xl font-display font-bold text-white">30+</div>
                    <div className="text-sm text-white/70">年技术积累</div>
                  </div>
                  <div className="w-px h-12 bg-white/30" />
                  <div>
                    <div className="text-3xl font-display font-bold text-white">200+</div>
                    <div className="text-sm text-white/70">专业研发人员</div>
                  </div>
                  <div className="w-px h-12 bg-white/30" />
                  <div>
                    <div className="text-3xl font-display font-bold text-white">40GHz</div>
                    <div className="text-sm text-white/70">最高工作频率</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
