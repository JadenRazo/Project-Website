import React from 'react';
import styled from 'styled-components';
import { useScrollTo } from '../../hooks/useScrollTo';

interface Section {
  id: string;
  title: string;
  isOpen: boolean;
}

interface DevPanelSectionNavProps {
  sections: Section[];
  activeSection?: string;
}

const NavContainer = styled.div`
  position: sticky;
  top: 80px;
  background: ${({ theme }) => theme.colors.card};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 2rem;
  z-index: 10;

  @media (max-width: 768px) {
    display: none;
  }
`;

const NavList = styled.ul`
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
`;

const NavItem = styled.li<{ isActive?: boolean; isOpen?: boolean }>`
  button {
    padding: 0.5rem 1rem;
    border: 1px solid ${({ theme }) => theme.colors.border};
    background: ${({ theme, isActive }) => 
      isActive ? theme.colors.primary : theme.colors.background};
    color: ${({ theme, isActive }) => 
      isActive ? 'white' : theme.colors.text};
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.875rem;
    position: relative;

    &:hover {
      background: ${({ theme, isActive }) => 
        isActive ? theme.colors.primary : theme.colors.border};
    }

    &::after {
      content: '';
      position: absolute;
      bottom: -6px;
      left: 50%;
      transform: translateX(-50%);
      width: 4px;
      height: 4px;
      background: ${({ theme, isOpen }) => 
        isOpen ? theme.colors.success : 'transparent'};
      border-radius: 50%;
    }
  }
`;

export const DevPanelSectionNav: React.FC<DevPanelSectionNavProps> = ({
  sections,
  activeSection
}) => {
  const { scrollToId } = useScrollTo();

  const handleNavClick = (sectionId: string) => {
    scrollToId(`section-${sectionId}`, {
      behavior: 'smooth',
      offset: 100
    });
  };

  return (
    <NavContainer>
      <NavList>
        {sections.map((section) => (
          <NavItem 
            key={section.id}
            isActive={activeSection === section.id}
            isOpen={section.isOpen}
          >
            <button onClick={() => handleNavClick(section.id)}>
              {section.title}
            </button>
          </NavItem>
        ))}
      </NavList>
    </NavContainer>
  );
};