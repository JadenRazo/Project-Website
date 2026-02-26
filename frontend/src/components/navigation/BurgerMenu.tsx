// src/components/navigation/BurgerMenu.tsx
import React, { useState, useEffect, useCallback, useRef, memo } from 'react';
import styled, { css, keyframes } from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { useTheme } from '../../hooks/useTheme';
import type { ThemeMode } from '../../styles/theme.types';
import { lockScroll, unlockScroll } from '../../utils/scrollLock';

// Types and Interfaces
interface MenuItem {
  href: string;
  icon: React.JSX.Element;
  label: string;
  delay: number;
  isExternal?: boolean;
  isService?: boolean; // Flag for backend services
  description?: string; // Optional description for services
}

interface BurgerMenuProps {
  className?: string;
  activeSection: string;
  toggleTheme: () => void;
  themeMode: ThemeMode;
  onNavigate: (sectionId: string) => void;
}

interface MenuItemProps extends MenuItem {
  onClick: () => void;
  $isActive?: boolean;
  children: React.ReactNode;
}

interface BurgerLineProps {
  isOpen: boolean;
  delay: number;
}


// SVG Icons Component with error boundaries
const Icons = {
  Home: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
      <polyline points="9 22 9 12 15 12 15 22"></polyline>
    </svg>
  ),
  Projects: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path>
    </svg>
  ),
  About: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
      <circle cx="12" cy="7" r="4"></circle>
    </svg>
  ),
  Contact: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
      <polyline points="22,6 12,13 2,6"></polyline>
    </svg>
  ),
  Resume: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
      <polyline points="14 2 14 8 20 8"></polyline>
      <line x1="16" y1="13" x2="8" y2="13"></line>
      <line x1="16" y1="17" x2="8" y2="17"></line>
      <polyline points="10 9 9 9 8 9"></polyline>
    </svg>
  ),
  LinkedIn: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z"></path>
      <rect x="2" y="9" width="4" height="12"></rect>
      <circle cx="4" cy="4" r="2"></circle>
    </svg>
  ),
  Sun: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <circle cx="12" cy="12" r="5"></circle>
      <line x1="12" y1="1" x2="12" y2="3"></line>
      <line x1="12" y1="21" x2="12" y2="23"></line>
      <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
      <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
      <line x1="1" y1="12" x2="3" y2="12"></line>
      <line x1="21" y1="12" x2="23" y2="12"></line>
      <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
      <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
    </svg>
  ),
  Moon: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
    </svg>
  ),
  Skills: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"></path>
      <path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"></path>
    </svg>
  ),
  // New icons for backend services
  DevPanel: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
      <line x1="3" y1="9" x2="21" y2="9"></line>
      <line x1="9" y1="21" x2="9" y2="9"></line>
    </svg>
  ),
  URLShortener: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M13.5 10.5L21 3"></path>
      <path d="M16 3h5v5"></path>
      <path d="M21 14v5a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5"></path>
      <line x1="10" y1="14" x2="18" y2="14"></line>
      <line x1="7" y1="17" x2="18" y2="17"></line>
    </svg>
  ),
  Messaging: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
      <line x1="8" y1="10" x2="16" y2="10"></line>
      <line x1="8" y1="14" x2="14" y2="14"></line>
    </svg>
  ),
  Services: (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <rect x="3" y="3" width="7" height="7"></rect>
      <rect x="14" y="3" width="7" height="7"></rect>
      <rect x="14" y="14" width="7" height="7"></rect>
      <rect x="3" y="14" width="7" height="7"></rect>
    </svg>
  )
};

// Animation Variants
const menuVariants = {
  closed: {
    x: "100%",
    transition: { 
      type: "spring", 
      stiffness: 400, 
      damping: 40 
    }
  },
  open: {
    x: 0,
    transition: {
      type: "spring",
      stiffness: 400,
      damping: 40,
      staggerChildren: 0.1,
      delayChildren: 0.2
    }
  }
};

const itemVariants = {
  closed: { 
    x: 50, 
    opacity: 0 
  },
  open: { 
    x: 0, 
    opacity: 1,
    transition: { 
      type: "spring", 
      stiffness: 300, 
      damping: 24 
    }
  }
};

const overlayVariants = {
  closed: { 
    opacity: 0,
    transition: {
      duration: 0.2
    }
  },
  open: { 
    opacity: 1,
    transition: {
      duration: 0.2
    }
  }
};

// Cool pulse animation for icon glow
const pulseAnimation = keyframes`
  0% {
    box-shadow: 0 0 0 0 rgba(var(--primary-rgb, 0, 120, 255), 0.7);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(var(--primary-rgb, 0, 120, 255), 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(var(--primary-rgb, 0, 120, 255), 0);
  }
`;

// Styled Components with complete implementations
const BurgerButton = styled(motion.button)<{ isOpen: boolean }>`
  position: fixed;
  top: 1.5rem;
  right: 1.5rem;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  width: 2.5rem;
  height: 2.5rem;
  padding: 0.5rem;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  background: ${({ theme }) => theme.colors.primaryLight};

  @media (max-width: 768px) {
    right: 1rem;
    top: 1rem;
    width: 2.75rem;
    height: 2.75rem;
  }
  border: 2px solid ${({ theme }) => theme.colors.primaryHover};
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;

  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover};
    transform: scale(1.05);
  }

  &:focus {
    outline: none;
    box-shadow: 0 0 0 3px ${({ theme }) => theme.colors.primary}80;
  }

  /* Pulse animation when open */
  ${({ isOpen }) => isOpen && css`
    animation: ${pulseAnimation} 2s infinite;
  `}
`;

const BurgerLine = styled(motion.div)<BurgerLineProps>`
  width: 100%;
  height: 2px;
  border-radius: 10px;
  background: ${({ theme }) => theme.colors.primary};
  position: relative;
  transform-origin: center;
  transition: all 0.3s ease;
`;

const Overlay = styled(motion.div)`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(3px);
  z-index: 990;
`;

const MenuContainer = styled(motion.div)`
  position: fixed;
  top: 0;
  right: 0;
  width: 300px;
  height: 100vh;
  background: ${({ theme }) => theme.colors.background};
  box-shadow: -2px 0 10px rgba(0, 0, 0, 0.1);
  z-index: 995;
  display: flex;
  flex-direction: column;
  padding: 5rem 2rem 2rem;
  overflow-y: auto;

  nav {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 1rem;
    flex: 1 0 auto;
  }

  @media (max-width: 768px) {
    width: 100%;
    padding: 5rem 1.5rem 1.5rem;
    
    nav {
      flex: 0 0 auto;
      margin-bottom: 0;
    }
  }
`;

const MenuLink = styled(motion.a)<{ $isActive?: boolean; $isService?: boolean }>`
  display: flex;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  text-decoration: none;
  font-size: 1rem;
  font-weight: 500;
  background: ${({ theme, $isActive, $isService }) => 
    $isActive ? theme.colors.primary : 
    $isService ? `linear-gradient(135deg, ${theme.colors.primary}, ${theme.colors.primaryHover})` : 
    theme.colors.primaryLight};
  color: ${({ theme, $isActive, $isService }) => 
    $isActive || $isService ? '#ffffff' : theme.colors.primary};
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;

  svg {
    margin-right: 1rem;
  }

  span {
    position: relative;
    
    &::after {
      content: '';
      position: absolute;
      width: ${({ $isActive }) => $isActive ? '100%' : '0'};
      height: 2px;
      bottom: -2px;
      left: 0;
      background: ${({ theme, $isActive }) => 
        $isActive ? '#ffffff' : theme.colors.primary};
      transition: width 0.3s ease;
    }
  }

  &:hover {
    background: ${({ theme, $isActive, $isService }) => 
      $isActive ? theme.colors.primary : 
      $isService ? `linear-gradient(135deg, ${theme.colors.primaryHover}, ${theme.colors.primary})` : 
      theme.colors.primaryHover};
    transform: translateY(-2px);

    span::after {
      width: 100%;
    }
  }

  &:active {
    transform: translateY(1px);
  }
  
  @media (max-width: 768px) {
    padding: 0.9rem;
    font-size: 0.95rem;
  }
  
  @media (max-width: 480px) {
    padding: 0.8rem;
    font-size: 0.9rem;
    
    svg {
      width: 18px;
      height: 18px;
      margin-right: 0.8rem;
    }
  }
`;

// Service description for the menu items
const ServiceDescription = styled.small`
  display: block;
  font-size: 0.8rem;
  margin-top: 4px;
  opacity: 0.85;
  font-weight: normal;
`;

// Section divider for menu categories
const MenuDivider = styled.div`
  display: flex;
  align-items: center;
  margin: 1rem 0 0.5rem 0;
  color: ${({ theme }) => theme.colors.text}80;
  font-size: 0.8rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 1px;
  
  &::before, &::after {
    content: '';
    flex: 1;
    height: 1px;
    background-color: ${({ theme }) => theme.colors.text}20;
  }
  
  &::before {
    margin-right: 0.5rem;
  }
  
  &::after {
    margin-left: 0.5rem;
  }
`;

// Ensure ThemeToggle is properly styled and visible
const ThemeToggle = styled(motion.button)`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: auto; 
  margin-bottom: 1rem;
  padding: 12px 16px;
  background-color: ${({ theme }) => theme.colors.primaryLight || `${theme.colors.background}80`};
  border: 1px solid ${({ theme }) => theme.colors.primary}30;
  border-radius: 8px;
  color: ${({ theme }) => theme.colors.primary};
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  width: 100%;
  z-index: 5;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  &:hover, &:focus {
    background-color: ${({ theme }) => theme.colors.primaryHover || `${theme.colors.primary}10`};
    transform: translateY(-2px);
  }

  &:active {
    transform: translateY(0);
  }

  @media (max-width: 768px) {
    margin-top: 16px; /* Fixed spacing on mobile */
    margin-bottom: 1.5rem;
    width: 100%;
    padding: 14px;
    font-size: 16px;
  }
`;

// Error Boundary Component
class MenuErrorBoundary extends React.Component<{ children: React.ReactNode }, { hasError: boolean }> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError() {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Menu Error:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return <div>Something went wrong with the menu. Please refresh the page.</div>;
    }
    return this.props.children;
  }
}

// Define menu items with proper typing
const MENU_ITEMS: MenuItem[] = [
  { 
    href: '#hero',
    icon: Icons.Home,
    label: 'Home',
    delay: 0.1
  },
  { 
    href: '#skills',
    icon: Icons.Skills,
    label: 'Skills',
    delay: 0.2
  },
  { 
    href: '#projects',
    icon: Icons.Projects,
    label: 'Projects',
    delay: 0.3
  },
  { 
    href: '#about',
    icon: Icons.About,
    label: 'About',
    delay: 0.4
  }
];

// Backend services menu items
const SERVICES_ITEMS: MenuItem[] = [
  {
    href: '/devpanel',
    icon: Icons.DevPanel,
    label: 'Developer Panel',
    delay: 0.5,
    isService: true,
    description: 'Manage your projects and website'
  },
  {
    href: '/urlshortener',
    icon: Icons.URLShortener,
    label: 'URL Shortener',
    delay: 0.6,
    isService: true,
    description: 'Create and manage short URLs'
  },
  {
    href: '/messaging',
    icon: Icons.Messaging,
    label: 'Messaging',
    delay: 0.7,
    isService: true,
    description: 'Real-time messaging platform'
  }
];

// External links
const EXTERNAL_ITEMS: MenuItem[] = [
  { 
    href: 'https://jadenrazo.dev/s/linkedin',
    icon: Icons.LinkedIn,
    label: 'LinkedIn',
    delay: 0.8,
    isExternal: true
  },
  { 
    href: '/resume.pdf',
    icon: Icons.Resume,
    label: 'Resume',
    delay: 0.9,
    isExternal: true
  }
];

// Memoized MenuItem component with error handling
const MenuItemComponent: React.FC<MenuItemProps> = ({ 
  delay, 
  href, 
  icon, 
  children, 
  onClick, 
  isExternal,
  $isActive,
  isService,
  description
}) => {
  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    try {
      if (isExternal) {
        return; // Let browser handle external links
      }
      e.preventDefault();
      onClick();
    } catch (error) {
      console.error('Navigation error:', error);
      // Fallback to default behavior if something goes wrong
      if (!isExternal) {
        window.location.href = href;
      }
    }
  };

  return (
    <MenuLink
      href={href}
      onClick={handleClick}
      variants={itemVariants}
      custom={delay}
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
      target={isExternal ? "_blank" : undefined}
      rel={isExternal ? "noopener noreferrer" : undefined}
      $isActive={$isActive}
      $isService={isService}
    >
      {icon}
      <div>
        <span>{children}</span>
        {description && <ServiceDescription>{description}</ServiceDescription>}
      </div>
    </MenuLink>
  );
};

// Memoize the MenuItem component to prevent unnecessary re-renders
const MemoizedMenuItem = memo(MenuItemComponent);

// Main BurgerMenu component with error boundary
export const BurgerMenu: React.FC<BurgerMenuProps> = ({ 
  className,
  activeSection,
  toggleTheme: externalToggleTheme,
  themeMode: externalThemeMode,
  onNavigate
}) => {
  // Use provided props with context as fallback
  const themeContext = useTheme();
  const finalThemeMode = externalThemeMode || themeContext.themeMode;
  const finalToggleTheme = externalToggleTheme || themeContext.toggleTheme;
  
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [isMobile, setIsMobile] = useState<boolean>(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const isDarkMode = finalThemeMode === 'dark';

  // Check if mobile on initial render and window resize
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth <= 768);
    };
    
    // Initial check
    checkMobile();
    
    // Add listener for window resize
    window.addEventListener('resize', checkMobile);
    
    // Cleanup
    return () => {
      window.removeEventListener('resize', checkMobile);
    };
  }, []);

  const closeMenu = useCallback(() => {
    setIsOpen(false);
  }, []);

  // Handle navigation for menu items
  const handleItemClick = useCallback((href: string, isExternal: boolean = false, isService: boolean = false) => {
    if (isExternal) return;
    
    try {
      // If it's a service route, handle with normal browser navigation
      if (isService || href.startsWith('/')) {
        // For service links, navigate using window.location
        window.location.href = href;
        closeMenu();
        return;
      }
      
      // For section navigation within the page
      const sectionId = href.replace('#', '');
      onNavigate(sectionId);
      closeMenu();
    } catch (error) {
      console.error('Navigation error:', error);
    }
  }, [onNavigate, closeMenu]);

  // Manual click outside logic
  useEffect(() => {
    if (!isOpen) return;
    
    const handleClickOutside = (event: MouseEvent) => {
      if (
        menuRef.current && 
        !menuRef.current.contains(event.target as Node) && 
        buttonRef.current && 
        !buttonRef.current.contains(event.target as Node)
      ) {
        closeMenu();
      }
    };
    
    document.addEventListener('mousedown', handleClickOutside);
    
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [isOpen, closeMenu]);

  // Keyboard handler for ESC key and prevent scrolling when menu is open
  useEffect(() => {
    if (!isOpen) return;

    try {
      const handleEsc = (event: KeyboardEvent): void => {
        if (event.key === 'Escape') closeMenu();
      };

      lockScroll();
      window.addEventListener('keydown', handleEsc);

      return () => {
        unlockScroll();
        window.removeEventListener('keydown', handleEsc);
      };
    } catch (err) {
      setError('Failed to initialize menu event handlers');
      console.error('Menu initialization error:', err);
    }
  }, [isOpen, closeMenu]);

  // Initialize pulse animation CSS variables
  useEffect(() => {
    // Additional setup for the pulse animation
    const root = document.documentElement;
    const primaryColor = getComputedStyle(document.body).getPropertyValue('--primary-color') || '#0078ff';
    
    // Convert hex to RGB if needed
    const hexToRgb = (hex: string): string => {
      // Remove # if present
      hex = hex.replace('#', '');
      
      // Convert 3-digit hex to 6-digits
      if (hex.length === 3) {
        hex = hex.split('').map(h => h + h).join('');
      }
      
      // Parse the hex values
      const r = parseInt(hex.substring(0, 2), 16);
      const g = parseInt(hex.substring(2, 4), 16);
      const b = parseInt(hex.substring(4, 6), 16);
      
      return `${r}, ${g}, ${b}`;
    };
    
    // Set the RGB values for the animation
    if (primaryColor.startsWith('#')) {
      root.style.setProperty('--primary-rgb', hexToRgb(primaryColor));
    }
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <MenuErrorBoundary>
      <BurgerButton
        ref={buttonRef}
        isOpen={isOpen}
        onClick={() => setIsOpen(prev => !prev)}
        whileHover={{ scale: 1.05 }}
        whileTap={{ scale: 0.95 }}
        aria-label={isOpen ? "Close menu" : "Open menu"}
        aria-expanded={isOpen}
        role="button"
        className={className}
      >
        {Array(3).fill(null).map((_, i) => (
          <BurgerLine
            key={i}
            isOpen={isOpen}
            delay={i * 0.1}
            animate={{
              rotate: isOpen && i !== 1 ? (i === 0 ? 45 : -45) : 0,
              y: isOpen ? (i === 0 ? 8 : i === 2 ? -8 : 0) : 0,
              opacity: isOpen && i === 1 ? 0 : 1
            }}
          />
        ))}
      </BurgerButton>

      <AnimatePresence>
        {isOpen && (
          <>
            <Overlay
              initial="closed"
              animate="open"
              exit="closed"
              variants={overlayVariants}
              onClick={closeMenu}
            />
            <MenuContainer
              ref={menuRef}
              className="menu-container"
              initial="closed"
              animate="open"
              exit="closed"
              variants={menuVariants}
              data-lenis-prevent
              data-scroll-lock-scrollable
            >
              <motion.nav
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ delay: 0.2 }}
              >
                {/* Main navigation items */}
                {MENU_ITEMS.map((item) => (
                  <MemoizedMenuItem
                    key={item.href}
                    {...item}
                    onClick={() => handleItemClick(item.href, item.isExternal, item.isService)}
                    $isActive={!item.isExternal && !item.isService && item.href.replace('#', '') === activeSection}
                  >
                    {item.label}
                  </MemoizedMenuItem>
                ))}

                {/* Services section */}
                <MenuDivider>Services</MenuDivider>
                {SERVICES_ITEMS.map((item) => (
                  <MemoizedMenuItem
                    key={item.href}
                    {...item}
                    onClick={() => handleItemClick(item.href, item.isExternal, item.isService)}
                    description={item.description}
                  >
                    {item.label}
                  </MemoizedMenuItem>
                ))}
                
                {/* External links */}
                <MenuDivider>Links</MenuDivider>
                {EXTERNAL_ITEMS.map((item) => (
                  <MemoizedMenuItem
                    key={item.href}
                    {...item}
                    onClick={() => handleItemClick(item.href, item.isExternal)}
                  >
                    {item.label}
                  </MemoizedMenuItem>
                ))}
              </motion.nav>

              {/* Fixed positioning for theme toggle */}
              <motion.div 
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ delay: 0.3 }}
                style={{ 
                  width: '100%', 
                  display: 'block', 
                  marginTop: isMobile ? '16px' : 'auto' 
                }}
              >
                <ThemeToggle
                  onClick={() => {
                    finalToggleTheme();
                    // Don't close menu on theme toggle to improve UX
                  }}
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  aria-label={isDarkMode ? "Switch to light mode" : "Switch to dark mode"}
                >
                  {isDarkMode ? Icons.Sun : Icons.Moon}
                  <span>{isDarkMode ? 'Light Mode' : 'Dark Mode'}</span>
                </ThemeToggle>
              </motion.div>
            </MenuContainer>
          </>
        )}
      </AnimatePresence>
    </MenuErrorBoundary>
  );
};

export default BurgerMenu;