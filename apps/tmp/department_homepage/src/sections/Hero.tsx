import { useEffect, useRef, useState } from 'react';
import { ArrowRight, Waves } from 'lucide-react';
import { Button } from '@/components/ui/button';

// 麦克斯韦方程组组件
const MaxwellEquations = () => {
  return (
    <div className="absolute top-20 left-10 opacity-20 text-white font-mono text-sm md:text-base lg:text-lg">
      <div className="space-y-2">
        <div className="flex items-center gap-2">
          <span className="text-space-blue">∇</span> · <span className="text-blue-400">E</span> = ρ/ε₀
        </div>
        <div className="flex items-center gap-2">
          <span className="text-space-blue">∇</span> · <span className="text-cyan-400">B</span> = 0
        </div>
        <div className="flex items-center gap-2">
          <span className="text-space-blue">∇</span> × <span className="text-blue-400">E</span> = -∂<span className="text-cyan-400">B</span>/∂t
        </div>
        <div className="flex items-center gap-2">
          <span className="text-space-blue">∇</span> × <span className="text-cyan-400">B</span> = μ₀<span className="text-yellow-400">J</span> + μ₀ε₀∂<span className="text-blue-400">E</span>/∂t
        </div>
      </div>
    </div>
  );
};

// 电磁波可视化Canvas
const ElectromagneticWave = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const resize = () => {
      canvas.width = canvas.offsetWidth * 2;
      canvas.height = canvas.offsetHeight * 2;
    };
    resize();
    window.addEventListener('resize', resize);

    let time = 0;
    const animate = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      
      const width = canvas.width;
      const height = canvas.height;
      const centerY = height / 2;
      const amplitude = height / 6;
      const frequency = 0.02;
      const speed = 0.05;

      // 绘制电场 (蓝色)
      ctx.beginPath();
      ctx.strokeStyle = 'rgba(59, 130, 246, 0.8)';
      ctx.lineWidth = 3;
      for (let x = 0; x < width; x++) {
        const y = centerY + amplitude * Math.sin(frequency * x - time * speed);
        if (x === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
      }
      ctx.stroke();

      // 绘制磁场 (青色)
      ctx.beginPath();
      ctx.strokeStyle = 'rgba(34, 211, 238, 0.8)';
      ctx.lineWidth = 3;
      for (let x = 0; x < width; x++) {
        const y = centerY + amplitude * Math.cos(frequency * x - time * speed);
        if (x === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
      }
      ctx.stroke();

      // 绘制传播方向箭头
      ctx.beginPath();
      ctx.strokeStyle = 'rgba(255, 255, 255, 0.5)';
      ctx.lineWidth = 2;
      ctx.setLineDash([10, 10]);
      ctx.moveTo(0, centerY);
      ctx.lineTo(width, centerY);
      ctx.stroke();
      ctx.setLineDash([]);

      time++;
      requestAnimationFrame(animate);
    };

    animate();
    return () => window.removeEventListener('resize', resize);
  }, []);

  return (
    <canvas 
      ref={canvasRef} 
      className="absolute bottom-0 left-0 right-0 h-32 md:h-48 w-full opacity-60"
    />
  );
};

export default function Hero() {
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    setIsLoaded(true);
  }, []);

  const scrollToSection = (id: string) => {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <section id="home" className="relative min-h-screen w-full overflow-hidden bg-black">
      {/* Background Gradient */}
      <div className="absolute inset-0 bg-gradient-to-br from-black via-gray-900 to-black opacity-80" />
      
      {/* 麦克斯韦方程组背景 */}
      <MaxwellEquations />
      
      {/* 电磁波可视化 */}
      <ElectromagneticWave />
      
      {/* Content */}
      <div className="relative z-10 grid min-h-screen grid-cols-1 lg:grid-cols-[1fr_0.8fr] items-center px-6 sm:px-12 lg:px-20">
        {/* Left Content */}
        <div className="flex flex-col justify-center py-20 lg:py-0">
          <div className="space-y-6">
            {/* Badge */}
            <div 
              className={`inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/10 backdrop-blur-sm transition-all duration-700 ${
                isLoaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-5'
              }`}
              style={{ transitionDelay: '200ms' }}
            >
              <Waves className="w-4 h-4 text-space-blue" />
              <span className="text-white/80 text-sm">航天科工集团 · 微波技术</span>
            </div>

            {/* Title */}
            <h1 
              className={`font-display text-4xl sm:text-5xl lg:text-6xl xl:text-7xl font-bold text-white leading-tight transition-all duration-1000 ${
                isLoaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'
              }`}
              style={{ transitionDelay: '300ms', transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
            >
              <span className="block overflow-hidden">
                <span className="inline-block">驾驭电磁波</span>
              </span>
              <span className="block overflow-hidden mt-2">
                <span className="inline-block text-space-blue">连接天地</span>
              </span>
            </h1>
            
            {/* Subtitle */}
            <p 
              className={`text-lg sm:text-xl text-gray-300 max-w-xl leading-relaxed transition-all duration-700 ${
                isLoaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'
              }`}
              style={{ transitionDelay: '600ms', transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
            >
              以麦克斯韦方程组为理论基础，专注于微波技术研发与应用，
              为航天事业提供先进的射频与微波解决方案。
            </p>
            
            {/* CTA Button */}
            <div 
              className={`pt-4 transition-all duration-600 ${
                isLoaded ? 'opacity-100 scale-100' : 'opacity-0 scale-90'
              }`}
              style={{ transitionDelay: '900ms', transitionTimingFunction: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)' }}
            >
              <Button 
                size="lg"
                className="bg-space-blue hover:bg-space-blue-dark text-white px-8 py-6 text-lg font-semibold rounded-full animate-pulse-glow group transition-all duration-300"
                onClick={() => scrollToSection('about')}
              >
                探索技术
                <ArrowRight className="ml-2 w-5 h-5 transition-transform duration-300 group-hover:translate-x-2 group-hover:rotate-[-45deg]" />
              </Button>
            </div>
          </div>
        </div>
        
        {/* Right Content - Microwave Image */}
        <div 
          className={`relative hidden lg:flex items-center justify-center transition-all duration-1200 ${
            isLoaded ? 'opacity-100 translate-x-0' : 'opacity-0 translate-x-24'
          }`}
          style={{ 
            transitionDelay: '400ms', 
            transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)',
            perspective: '1000px'
          }}
        >
          <div className="relative animate-float">
            {/* Glow Effect */}
            <div className="absolute inset-0 bg-gradient-to-r from-blue-500 via-cyan-500 to-blue-500 opacity-30 blur-3xl rounded-full scale-75" />
            
            {/* Microwave Image */}
            <img 
              src="/hero-microwave.jpg" 
              alt="电磁波可视化" 
              className="relative z-10 w-full max-w-md xl:max-w-lg rounded-2xl shadow-2xl"
              style={{ 
                transform: 'rotateY(-5deg) rotateX(3deg)',
                transformStyle: 'preserve-3d'
              }}
            />
            
            {/* 频率标签 */}
            <div className="absolute -bottom-4 -right-4 bg-black/80 backdrop-blur-sm px-4 py-2 rounded-lg border border-space-blue/30">
              <span className="text-space-blue font-mono text-sm">f = 1-40 GHz</span>
            </div>
          </div>
        </div>
      </div>
      
      {/* Bottom Gradient Fade */}
      <div className="absolute bottom-0 left-0 right-0 h-32 bg-gradient-to-t from-white to-transparent z-20" />
    </section>
  );
}
