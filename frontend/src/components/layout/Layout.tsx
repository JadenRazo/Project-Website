// src/components/layout/Layout.tsx
import styled from 'styled-components';

const LayoutWrapper = styled.div`
  position: relative;
  min-height: 100vh;
`;

export const Layout: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <LayoutWrapper>
    {children}
  </LayoutWrapper>
);
