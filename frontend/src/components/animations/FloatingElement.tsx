// src/components/animations/FloatingElement.tsx
import React from 'react';
import styled from 'styled-components';
import { motion, useMotionValue, useTransform } from 'framer-motion';

const FloatingContainer = styled(motion.div)`
  position: absolute;
  pointer-events: none;
`;

interface FloatingElementProps {
  children: React.ReactNode;
  x?: number;
  y?: number;
}

export const FloatingElement: React.FC<FloatingElementProps> = ({ children, x = 0, y = 0 }) => {
  const mouseX = useMotionValue(0);
  const mouseY = useMotionValue(0);

  const translateX = useTransform(mouseX, [-500, 500], [-15, 15]);
  const translateY = useTransform(mouseY, [-500, 500], [-15, 15]);

  const handleMouseMove = (event: MouseEvent) => {
    mouseX.set(event.clientX - window.innerWidth / 2);
    mouseY.set(event.clientY - window.innerHeight / 2);
  };

  React.useEffect(() => {
    window.addEventListener('mousemove', handleMouseMove);
    return () => window.removeEventListener('mousemove', handleMouseMove);
  }, []);

  return (
    <FloatingContainer
      style={{
        translateX,
        translateY,
        left: x,
        top: y,
      }}
    >
      {children}
    </FloatingContainer>
  );
};
