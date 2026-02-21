import { useState, useEffect } from 'react';
import { Menu, X, Radio } from 'lucide-react';
import { Button } from '@/components/ui/button';

const navLinks = [
  { name: '首页', href: '#home' },
  { name: '关于', href: '#about' },
  { name: '技术', href: '#services' },
  { name: '产品', href: '#projects' },
  { name: '文化', href: '#culture' },
];

export default function Navbar() {
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isLoaded, setIsLoaded] = useState(false);

  useEffect(() => {
    setIsLoaded(true);
    
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 50);
    };
    
    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const scrollToSection = (href: string) => {
    const id = href.replace('#', '');
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
    setIsMobileMenuOpen(false);
  };

  return (
    <nav 
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-500 ${
        isScrolled 
          ? 'bg-white/95 backdrop-blur-lg shadow-sm h-[70px]' 
          : 'bg-transparent h-[90px]'
      }`}
      style={{ transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
    >
      <div className="max-w-7xl mx-auto px-6 sm:px-12 h-full">
        <div className="flex items-center justify-between h-full">
          {/* Logo */}
          <a 
            href="#home"
            onClick={(e) => { e.preventDefault(); scrollToSection('#home'); }}
            className={`flex items-center gap-2 transition-all duration-600 ${
              isLoaded ? 'opacity-100 scale-100' : 'opacity-0 scale-80'
            }`}
            style={{ transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
          >
            <Radio 
              className={`w-8 h-8 transition-colors duration-300 ${
                isScrolled ? 'text-space-blue' : 'text-space-blue'
              }`} 
            />
            <div className="flex flex-col">
              <span 
                className={`font-display text-xl font-bold transition-colors duration-300 ${
                  isScrolled ? 'text-space-black' : 'text-white'
                }`}
              >
                微波室
              </span>
              <span 
                className={`text-xs transition-colors duration-300 ${
                  isScrolled ? 'text-space-gray' : 'text-white/70'
                }`}
              >
                航天科工集团
              </span>
            </div>
          </a>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center gap-8">
            {navLinks.map((link, index) => (
              <a
                key={link.name}
                href={link.href}
                onClick={(e) => { e.preventDefault(); scrollToSection(link.href); }}
                className={`relative text-sm font-medium transition-all duration-500 link-underline ${
                  isScrolled ? 'text-space-gray hover:text-space-blue' : 'text-white/90 hover:text-white'
                } ${isLoaded ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-5'}`}
                style={{ 
                  transitionDelay: `${100 + index * 100}ms`,
                  transitionTimingFunction: 'cubic-bezier(0.16, 1, 0.3, 1)'
                }}
              >
                {link.name}
              </a>
            ))}
          </div>

          {/* CTA Button */}
          <div 
            className={`hidden md:block transition-all duration-500 ${
              isLoaded ? 'opacity-100 scale-100' : 'opacity-0 scale-90'
            }`}
            style={{ 
              transitionDelay: '500ms',
              transitionTimingFunction: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)'
            }}
          >
            <Button 
              className="bg-space-blue hover:bg-space-blue-dark text-white rounded-full px-6 transition-all duration-300 hover:shadow-glow"
              onClick={() => scrollToSection('#contact')}
            >
              联系我们
            </Button>
          </div>

          {/* Mobile Menu Button */}
          <button
            className="md:hidden p-2"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? (
              <X className={`w-6 h-6 ${isScrolled ? 'text-space-black' : 'text-white'}`} />
            ) : (
              <Menu className={`w-6 h-6 ${isScrolled ? 'text-space-black' : 'text-white'}`} />
            )}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      <div 
        className={`md:hidden absolute top-full left-0 right-0 bg-white shadow-lg transition-all duration-500 overflow-hidden ${
          isMobileMenuOpen ? 'max-h-96 opacity-100' : 'max-h-0 opacity-0'
        }`}
        style={{ transitionTimingFunction: 'cubic-bezier(0.23, 1, 0.32, 1)' }}
      >
        <div className="px-6 py-4 space-y-4">
          {navLinks.map((link) => (
            <a
              key={link.name}
              href={link.href}
              onClick={(e) => { e.preventDefault(); scrollToSection(link.href); }}
              className="block text-space-gray hover:text-space-blue font-medium transition-colors duration-300"
            >
              {link.name}
            </a>
          ))}
          <Button 
            className="w-full bg-space-blue hover:bg-space-blue-dark text-white rounded-full"
            onClick={() => scrollToSection('#contact')}
          >
            联系我们
          </Button>
        </div>
      </div>
    </nav>
  );
}
