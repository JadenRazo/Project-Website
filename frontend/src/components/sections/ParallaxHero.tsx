// src/components/sections/ParallaxHero.tsx
import React from 'react';
import styled from 'styled-components';
import { motion, useScroll, useTransform } from 'framer-motion';

const ParallaxContainer = styled.div`
  height: 100vh;
  overflow: hidden;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const ParallaxLayer = styled(motion.div)`
  position: absolute;
  width: 100%;
  height: 100%;
`;

export const ParallaxHero = () => {
  const { scrollY } = useScroll();
  const y1 = useTransform(scrollY, [0, 500], [0, 100]);
  const y2 = useTransform(scrollY, [0, 500], [0, -100]);
  const opacity = useTransform(scrollY, [0, 300], [1, 0]);

  return (
    <ParallaxContainer>
      <ParallaxLayer style={{ y: y1, opacity }}>
        {/* Background elements */}
      </ParallaxLayer>
      <ParallaxLayer style={{ y: y2 }}>
        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
        >
          Your Name
        </motion.h1>
      </ParallaxLayer>
    </ParallaxContainer>
  );
};
