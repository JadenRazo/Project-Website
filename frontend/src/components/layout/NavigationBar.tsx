// src/components/layout/NavigationBar.tsx
import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';

const Nav = styled(motion.nav)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
`;

const Logo = styled.div`
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--primary);
`;

const NavItems = styled.div`
  display: flex;
  align-items: center;
  gap: 2rem;
`;

interface NavigationBarProps {
  isDarkMode: boolean;
  toggleTheme: () => void;
}


export const NavigationBar: React.FC<NavigationBarProps> = ({ isDarkMode, toggleTheme }) => {
  return (
    <Nav
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      transition={{ duration: 0.5 }}
    >
      <Logo></Logo>
      <NavItems>
        {/* Your other nav items */}
      </NavItems>
    </Nav>
  );
};
