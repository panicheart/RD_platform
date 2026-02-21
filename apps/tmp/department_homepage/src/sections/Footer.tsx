import { useEffect, useRef, useState } from 'react';
import { Radio, Mail, Phone, MapPin, Send, Twitter, Linkedin, Youtube, Instagram } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

const quickLinks = [
  { name: '首页', href: '#home' },
  { name: '关于', href: '#about' },
  { name: '技术', href: '#services' },
  { name: '产品', href: '#projects' },
  { name: '文化', href: '#culture' },
];

const socialLinks = [
  { icon: Twitter, href: '#', label: 'Twitter' },
  { icon: Linkedin, href: '#', label: 'LinkedIn' },
  { icon: Youtube, href: '#', label: 'YouTube' },
  { icon: Instagram, href: '#', label: 'Instagram' },
];

export default function Footer() {
  const footerRef = useRef<HTMLElement>(null);
  const [isVisible, setIsVisible] = useState(false);
  const [email, setEmail] = useState('');

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

    if (footerRef.current) {
      observer.observe(footerRef.current);
    }

    return () => observer.disconnect();
  }, []);

  const scrollToSection = (href: string) => {
    const id = href.replace('#', '');
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };

  const handleSubscribe = (e: React.FormEvent) => {
    e.preventDefault();
    alert('感谢订阅！我们会定期发送最新微波技术资讯给您。');
    setEmail('');
  };

  return (
    <footer 
      id="contact"
      ref={footerRef}
      className="relative bg-space-black pt-20 pb-8 overflow-hidden"
    >
      {/* 麦克斯韦方程背景 */}
      <div className="absolute top-20 right-20 opacity-5 text-white font-mono text-lg">
        <div>∇ × E = -∂B/∂t</div>
        <div className="mt-2">∇ × B = μ₀J + μ₀ε₀∂E/∂t</div>
      </div>

      {/* Particle Effect (subtle) */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-10 left-10 w-1 h-1 bg-space-blue rounded-full animate-pulse" />
        <div className="absolute top-20 right-20 w-1 h-1 bg-space-blue rounded-full animate-pulse" style={{ animationDelay: '1s' }} />
        <div className="absolute bottom-20 left-1/4 w-1 h-1 bg-space-blue rounded-full animate-pulse" style={{ animationDelay: '2s' }} />
        <div className="absolute top-1/2 right-1/3 w-1 h-1 bg-space-blue rounded-full animate-pulse" style={{ animationDelay: '0.5s' }} />
      </div>

      <div className="relative z-10 max-w-7xl mx-auto px-6 sm:px-12">
        {/* Main Footer Content */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 mb-16">
          {/* Brand Column */}
          <div 
            className={`transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'
            }`}
            style={{ transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
          >
            <a href="#home" className="flex items-center gap-2 mb-4">
              <Radio className="w-8 h-8 text-space-blue" />
              <div className="flex flex-col">
                <span className="font-display text-xl font-bold text-white">
                  微波室
                </span>
                <span className="text-xs text-white/50">
                  航天科工集团
                </span>
              </div>
            </a>
            <p className="text-white/60 text-sm leading-relaxed mb-6">
              以麦克斯韦方程为理论基础，专注于微波技术研发与应用，
              为航天事业提供先进的射频与微波解决方案。
            </p>
            {/* Social Links */}
            <div className="flex gap-3">
              {socialLinks.map((social, index) => (
                <a
                  key={social.label}
                  href={social.href}
                  className={`w-10 h-10 rounded-full bg-white/10 flex items-center justify-center text-white/60 hover:bg-space-blue hover:text-white transition-all duration-300 ${
                    isVisible ? 'opacity-100 scale-100' : 'opacity-0 scale-0'
                  }`}
                  style={{ 
                    transitionDelay: `${200 + index * 50}ms`,
                    transitionTimingFunction: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)'
                  }}
                  aria-label={social.label}
                >
                  <social.icon className="w-4 h-4" />
                </a>
              ))}
            </div>
          </div>

          {/* Quick Links */}
          <div 
            className={`transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'
            }`}
            style={{ transitionDelay: '100ms', transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
          >
            <h3 className="font-display text-lg font-bold text-white mb-6">
              快速链接
            </h3>
            <ul className="space-y-3">
              {quickLinks.map((link) => (
                <li key={link.name}>
                  <a
                    href={link.href}
                    onClick={(e) => { e.preventDefault(); scrollToSection(link.href); }}
                    className="text-white/60 hover:text-space-blue transition-colors duration-300 text-sm link-underline"
                  >
                    {link.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          {/* Contact Info */}
          <div 
            className={`transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'
            }`}
            style={{ transitionDelay: '200ms', transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
          >
            <h3 className="font-display text-lg font-bold text-white mb-6">
              联系我们
            </h3>
            <ul className="space-y-4">
              <li className="flex items-start gap-3">
                <MapPin className="w-5 h-5 text-space-blue flex-shrink-0 mt-0.5" />
                <span className="text-white/60 text-sm">
                  北京市海淀区航天路88号<br />
                  航天科工集团微波室
                </span>
              </li>
              <li className="flex items-center gap-3">
                <Phone className="w-5 h-5 text-space-blue flex-shrink-0" />
                <span className="text-white/60 text-sm">
                  +86 10 8888 8888
                </span>
              </li>
              <li className="flex items-center gap-3">
                <Mail className="w-5 h-5 text-space-blue flex-shrink-0" />
                <span className="text-white/60 text-sm">
                  microwave@casic.cn
                </span>
              </li>
            </ul>
          </div>

          {/* Newsletter */}
          <div 
            className={`transition-all duration-700 ${
              isVisible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'
            }`}
            style={{ transitionDelay: '300ms', transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)' }}
          >
            <h3 className="font-display text-lg font-bold text-white mb-6">
              订阅资讯
            </h3>
            <p className="text-white/60 text-sm mb-4">
              获取最新的微波技术动态
            </p>
            <form onSubmit={handleSubscribe} className="flex gap-2">
              <Input
                type="email"
                placeholder="您的邮箱"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="bg-white/10 border-white/20 text-white placeholder:text-white/40 focus:border-space-blue"
                required
              />
              <Button 
                type="submit"
                className="bg-space-blue hover:bg-space-blue-dark text-white px-4"
              >
                <Send className="w-4 h-4" />
              </Button>
            </form>
          </div>
        </div>

        {/* Divider */}
        <div className="border-t border-white/10 pt-8">
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            <p className="text-white/40 text-sm text-center sm:text-left">
              © 2024 航天科工集团微波室. 保留所有权利.
            </p>
            <div className="flex gap-6">
              <a href="#" className="text-white/40 hover:text-white text-sm transition-colors duration-300">
                隐私政策
              </a>
              <a href="#" className="text-white/40 hover:text-white text-sm transition-colors duration-300">
                使用条款
              </a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}
