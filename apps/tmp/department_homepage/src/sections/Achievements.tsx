import { useEffect, useRef, useState } from 'react';
import { Radio, Waves, Zap, Award } from 'lucide-react';

const stats = [
  {
    icon: Radio,
    value: 40,
    suffix: 'GHz',
    label: '最高工作频率',
  },
  {
    icon: Waves,
    value: 500,
    suffix: '+',
    label: '微波组件型号',
  },
  {
    icon: Zap,
    value: 10,
    suffix: 'kW',
    label: '峰值输出功率',
  },
  {
    icon: Award,
    value: 100,
    suffix: '+',
    label: '技术专利',
  },
];

function AnimatedCounter({ value, suffix, isVisible }: { value: number; suffix: string; isVisible: boolean }) {
  const [count, setCount] = useState(0);

  useEffect(() => {
    if (!isVisible) return;

    const duration = 1500;
    const steps = 60;
    const increment = value / steps;
    let current = 0;

    const timer = setInterval(() => {
      current += increment;
      if (current >= value) {
        setCount(value);
        clearInterval(timer);
      } else {
        setCount(Math.floor(current));
      }
    }, duration / steps);

    return () => clearInterval(timer);
  }, [value, isVisible]);

  return (
    <span className="font-display text-4xl sm:text-5xl lg:text-6xl font-bold text-space-black">
      {count}{suffix}
    </span>
  );
}

export default function Achievements() {
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
      { threshold: 0.3 }
    );

    if (sectionRef.current) {
      observer.observe(sectionRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <section 
      ref={sectionRef}
      className="relative py-24 lg:py-32 bg-space-light overflow-hidden"
    >
      {/* 麦克斯韦方程背景 */}
      <div className="absolute inset-0 opacity-5">
        <div className="absolute top-10 left-1/4 text-space-black font-mono text-2xl rotate-12">
          c = 1/√(μ₀ε₀)
        </div>
        <div className="absolute bottom-20 right-1/4 text-space-black font-mono text-2xl rotate-[-8deg]">
          λ = c/f
        </div>
      </div>

      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-5">
        <div className="absolute inset-0" style={{
          backgroundImage: `radial-gradient(circle at 2px 2px, #3b82f6 1px, transparent 0)`,
          backgroundSize: '40px 40px'
        }} />
      </div>

      <div className="relative z-10 max-w-7xl mx-auto px-6 sm:px-12">
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
              技术实力
            </span>
            <div className="w-8 h-[2px] bg-space-blue" />
          </div>
          
          <h2 
            className={`font-display text-3xl sm:text-4xl lg:text-5xl font-bold text-space-black transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
            }`}
            style={{ transitionDelay: '100ms', transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
          >
            数字见证<span className="text-space-blue">实力</span>
          </h2>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-6 lg:gap-8">
          {stats.map((stat, index) => (
            <div
              key={stat.label}
              className={`group relative bg-white rounded-2xl p-6 lg:p-8 text-center shadow-card hover:shadow-card-hover transition-all duration-600 ${
                isVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-80'
              }`}
              style={{ 
                transitionDelay: `${index * 150}ms`,
                transitionTimingFunction: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)'
              }}
            >
              {/* Icon */}
              <div 
                className={`inline-flex items-center justify-center w-14 h-14 rounded-full bg-blue-50 mb-4 transition-all duration-600 ${
                  isVisible ? 'rotate-0 scale-100' : 'rotate-180 scale-0'
                }`}
                style={{ 
                  transitionDelay: `${index * 150}ms`,
                  transitionTimingFunction: 'cubic-bezier(0.34, 1.56, 0.64, 1)'
                }}
              >
                <stat.icon className="w-7 h-7 text-space-blue group-hover:scale-110 transition-transform duration-300" />
              </div>

              {/* Value */}
              <div className="mb-2">
                <AnimatedCounter 
                  value={stat.value} 
                  suffix={stat.suffix}
                  isVisible={isVisible}
                />
              </div>

              {/* Label */}
              <p className="text-space-gray font-medium">
                {stat.label}
              </p>

              {/* Hover Glow */}
              <div className="absolute inset-0 rounded-2xl bg-space-blue opacity-0 group-hover:opacity-5 transition-opacity duration-300" />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
