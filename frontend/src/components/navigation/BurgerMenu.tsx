// src/components/navigation/BurgerMenu.tsx
import React, { useState, useEffect, useCallback, useRef, memo } from 'react';
import styled, { css, keyframes, useTheme as useStyledTheme } from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { useTheme } from '../../contexts/ThemeContext';
import { useClickOutside } from '../../hooks/useClickOutside';

// Types and Interfaces
interface MenuItem {
  href: string;
  icon: React.JSX.Element;
  label: string;
  delay: number;
  isExternal?: boolean;
}

interface BurgerMenuProps {
  className?: string;
}

interface MenuItemProps extends MenuItem {
  onClick: () => void;
  children: React.ReactNode;
}

interface BurgerLineProps {
  isOpen: boolean;
  delay: number;
}

// Custom Error Classes
class NavigationError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'NavigationError';
  }
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
  padding: 5rem 2rem 2rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  background: ${({ theme }) => `${theme.colors.backgroundAlt}F0`};
  box-shadow: -5px 0 25px rgba(0, 0, 0, 0.2);
  z-index: 995;
  overflow-y: auto;

  nav {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin-bottom: 2rem;
  }

  @media (max-width: 768px) {
    width: 80%;
  }
`;

const MenuLink = styled(motion.a)`
  display: flex;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  text-decoration: none;
  font-size: 1rem;
  font-weight: 500;
  background: ${({ theme }) => theme.colors.primaryLight};
  color: ${({ theme }) => theme.colors.primary};
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
      width: 0;
      height: 2px;
      bottom: -2px;
      left: 0;
      background: ${({ theme }) => theme.colors.primary};
      transition: width 0.3s ease;
    }
  }

  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover};
    transform: translateY(-2px);

    span::after {
      width: 100%;
    }
  }

  &:active {
    transform: translateY(1px);
  }
`;

const ThemeToggle = styled(motion.button)`
  display: flex;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  border: none;
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  cursor: pointer;
  background: ${({ theme }) => theme.colors.primaryLight};
  color: ${({ theme }) => theme.colors.primary};
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  margin-top: auto;

  svg {
    margin-right: 1rem;
  }

  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover};
    transform: translateY(-2px);
  }

  &:focus {
    outline: none;
    box-shadow: 0 0 0 2px ${({ theme }) => theme.colors.primary};
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
    href: '#home',
    icon: Icons.Home,
    label: 'Home',
    delay: 0.1,
    isExternal: false
  },
  { 
    href: '#projects',
    icon: Icons.Projects,
    label: 'Projects',
    delay: 0.2,
    isExternal: false
  },
  { 
    href: '#about',
    icon: Icons.About,
    label: 'About',
    delay: 0.3,
    isExternal: false
  },
  { 
    href: '#contact',
    icon: Icons.Contact,
    label: 'Contact',
    delay: 0.4,
    isExternal: false
  },
  { 
    href: '#resume',
    icon: Icons.Resume,
    label: 'Resume',
    delay: 0.5,
    isExternal: false
  },
  { 
    href: 'https://linkedin.com/in/jadenrazo',
    icon: Icons.LinkedIn,
    label: 'LinkedIn',
    delay: 0.6,
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
  isExternal 
}) => {
  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    try {
      if (isExternal) {
        return;
      }
      e.preventDefault();
      const element = document.querySelector(href);
      if (!element) {
        throw new NavigationError(`Element with selector "${href}" not found`);
      }
      element.scrollIntoView({ 
        behavior: 'smooth', 
        block: 'start' 
      });
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
    >
      {icon}
      <span>{children}</span>
    </MenuLink>
  );
};

// Memoize the MenuItem component to prevent unnecessary re-renders
const MenuItem = memo(MenuItemComponent);

// Manual implementation of click outside logic without the hook
const useManualClickOutside = <T extends HTMLElement>(
  ref: React.RefObject<T>,
  buttonRef: React.RefObject<HTMLElement>,
  isOpen: boolean,
  closeMenu: () => void
) => {
  useEffect(() => {
    // Return early if not open - no need to add listeners
    if (!isOpen) return;
    
    const handleClickOutside = (event: MouseEvent) => {
      // Check if the click was outside both the menu and the button
      if (
        ref.current && 
        !ref.current.contains(event.target as Node) && 
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
  }, [ref, buttonRef, isOpen, closeMenu]);
};

// Main BurgerMenu component with error boundary
export const BurgerMenu: React.FC<BurgerMenuProps> = ({ className }) => {
  const { themeMode, toggleTheme } = useTheme();
  const [isOpen, setIsOpen] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const menuRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);
  const isDarkMode = themeMode === 'dark';

  const closeMenu = useCallback(() => {
    setIsOpen(false);
  }, []);

  // Use manual click outside logic instead of the hook to avoid TypeScript errors
  useManualClickOutside(menuRef, buttonRef, isOpen, closeMenu);

  useEffect(() => {
    if (!isOpen) return;

    try {
      const handleEsc = (event: KeyboardEvent): void => {
        if (event.key === 'Escape') closeMenu();
      };
      
      document.body.style.overflow = 'hidden';
      window.addEventListener('keydown', handleEsc);
      
      return () => {
        document.body.style.overflow = '';
        window.removeEventListener('keydown', handleEsc);
      };
    } catch (err) {
      setError('Failed to initialize menu event handlers');
      console.error('Menu initialization error:', err);
    }
  }, [isOpen, closeMenu]);

  // Log to debug if the component is rendering
  useEffect(() => {
    console.log('BurgerMenu component rendered');
    
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
            >
              <motion.nav
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{ delay: 0.2 }}
              >
                {MENU_ITEMS.map((item) => (
                  <MenuItem
                    key={item.href}
                    {...item}
                    onClick={closeMenu}
                  >
                    {item.label}
                  </MenuItem>
                ))}
              </motion.nav>

              <ThemeToggle
                onClick={() => {
                  toggleTheme();
                  closeMenu();
                }}
                whileHover={{ scale: 1.05 }}
                whileTap={{ scale: 0.95 }}
                aria-label={isDarkMode ? "Switch to light mode" : "Switch to dark mode"}
              >
                {isDarkMode ? Icons.Sun : Icons.Moon}
                <span>{isDarkMode ? 'Light Mode' : 'Dark Mode'}</span>
              </ThemeToggle>
            </MenuContainer>
          </>
        )}
      </AnimatePresence>
    </MenuErrorBoundary>
  );
};

export default BurgerMenu;
