import { useEffect, useRef, useState } from 'react';
import { Star, Heart, Zap, Globe, Award, Users, Radio } from 'lucide-react';

const values = [
  {
    icon: Star,
    title: '科学严谨',
    description: '以麦克斯韦方程为理论基础，追求精确与完美',
  },
  {
    icon: Heart,
    title: ' passionate',
    description: '对微波技术的热爱是推动创新的源动力',
  },
  {
    icon: Zap,
    title: '勇于创新',
    description: '不断探索电磁波应用的新边界',
  },
  {
    icon: Globe,
    title: '服务航天',
    description: '以微波技术助力国家航天事业发展',
  },
  {
    icon: Award,
    title: '质量至上',
    description: '每一个组件都经过严格测试与验证',
  },
  {
    icon: Users,
    title: '协同攻关',
    description: '团队协作攻克技术难关',
  },
];

const testimonials = [
  {
    quote: '在微波室工作，让我有机会将麦克斯韦的理论转化为改变世界的技术。',
    author: '张工',
    role: '射频工程师',
  },
  {
    quote: '每一次成功的雷达测试，都是对电磁波理论的完美验证。',
    author: '李工',
    role: '测试工程师',
  },
  {
    quote: '从1GHz到40GHz，我们在每个频段都追求卓越。',
    author: '王工',
    role: '系统工程师',
  },
];

export default function Culture() {
  const sectionRef = useRef<HTMLElement>(null);
  const [isVisible, setIsVisible] = useState(false);
  const [activeTestimonial, setActiveTestimonial] = useState(0);

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

  // Auto-rotate testimonials
  useEffect(() => {
    const timer = setInterval(() => {
      setActiveTestimonial((prev) => (prev + 1) % testimonials.length);
    }, 5000);
    return () => clearInterval(timer);
  }, []);

  return (
    <section 
      id="culture" 
      ref={sectionRef}
      className="relative py-24 lg:py-32 bg-white overflow-hidden"
    >
      {/* 电磁波公式装饰 */}
      <div className="absolute top-40 left-10 opacity-5 text-space-black font-mono text-xl rotate-[-5deg]">
        P = |E|²/2η
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
              部门文化
            </span>
            <div className="w-8 h-[2px] bg-space-blue" />
          </div>
          
          <h2 
            className={`font-display text-3xl sm:text-4xl lg:text-5xl font-bold text-space-black mb-4 transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
            }`}
            style={{ transitionDelay: '100ms', transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
          >
            我们的<span className="text-space-blue">价值观</span>
          </h2>
          
          <p 
            className={`text-space-gray max-w-2xl mx-auto transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'
            }`}
            style={{ transitionDelay: '200ms', transitionTimingFunction: 'cubic-bezier(0.4, 0, 0.2, 1)' }}
          >
            微波室的文化建立在电磁波理论之上——严谨、精确、充满能量
          </p>
        </div>

        {/* Values Grid */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mb-20">
          {values.map((value, index) => (
            <div
              key={value.title}
              className={`group relative p-6 lg:p-8 rounded-2xl bg-gray-50 hover:bg-blue-50 transition-all duration-500 ${
                isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-12'
              }`}
              style={{ 
                transitionDelay: `${300 + index * 100}ms`,
                transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)'
              }}
            >
              {/* Icon */}
              <div className="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-white shadow-sm mb-4 group-hover:shadow-md group-hover:scale-110 transition-all duration-300">
                <value.icon className="w-6 h-6 text-space-blue" />
              </div>

              {/* Content */}
              <h3 className="font-display text-xl font-bold text-space-black mb-2">
                {value.title}
              </h3>
              <p className="text-space-gray text-sm leading-relaxed">
                {value.description}
              </p>

              {/* Hover Border */}
              <div className="absolute inset-0 rounded-2xl border-2 border-transparent group-hover:border-space-blue/20 transition-colors duration-300" />
            </div>
          ))}
        </div>

        {/* Testimonials */}
        <div 
          className={`relative bg-space-black rounded-3xl p-8 lg:p-12 overflow-hidden transition-all duration-700 ${
            isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-12'
          }`}
          style={{ transitionDelay: '600ms', transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
        >
          {/* Background Glow */}
          <div className="absolute top-0 right-0 w-96 h-96 bg-space-blue opacity-10 blur-3xl rounded-full -translate-y-1/2 translate-x-1/2" />
          
          {/* 麦克斯韦方程装饰 */}
          <div className="absolute bottom-4 left-4 opacity-10 text-white font-mono text-sm">
            ∇ · B = 0
          </div>
          
          <div className="relative z-10">
            <div className="text-center mb-8">
              <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-space-blue/20 mb-4">
                <Radio className="w-6 h-6 text-space-blue" />
              </div>
              <h3 className="font-display text-2xl font-bold text-white mb-2">
                员工心声
              </h3>
              <p className="text-white/60">
                听听微波人对我们文化的理解
              </p>
            </div>

            {/* Testimonial Content */}
            <div className="relative min-h-[160px]">
              {testimonials.map((testimonial, index) => (
                <div
                  key={index}
                  className={`absolute inset-0 flex flex-col items-center justify-center text-center transition-all duration-500 ${
                    activeTestimonial === index 
                      ? 'opacity-100 translate-y-0' 
                      : 'opacity-0 translate-y-4 pointer-events-none'
                  }`}
                >
                  <p className="text-white text-lg lg:text-xl italic mb-6 max-w-2xl">
                    "{testimonial.quote}"
                  </p>
                  <div>
                    <p className="text-white font-semibold">{testimonial.author}</p>
                    <p className="text-white/60 text-sm">{testimonial.role}</p>
                  </div>
                </div>
              ))}
            </div>

            {/* Dots */}
            <div className="flex justify-center gap-2 mt-6">
              {testimonials.map((_, index) => (
                <button
                  key={index}
                  onClick={() => setActiveTestimonial(index)}
                  className={`w-2 h-2 rounded-full transition-all duration-300 ${
                    activeTestimonial === index 
                      ? 'bg-space-blue w-6' 
                      : 'bg-white/30 hover:bg-white/50'
                  }`}
                />
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
