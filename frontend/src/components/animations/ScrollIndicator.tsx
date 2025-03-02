import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';

const ScrollContainer = styled(motion.div)`
  position: absolute; /* Changed from fixed to absolute */
  bottom: 40px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  z-index: 1000; /* Ensure it stays on top of other content */
  user-select: none; /* Disable text selection */
`;

const ScrollText = styled(motion.span)`
  font-size: 14px;
  color: var(--primary);
  opacity: 0.8;
  font-weight: 500;
  letter-spacing: 0.5px;
  pointer-events: none; /* Disable pointer events */
`;

const ArrowContainer = styled(motion.div)`
  width: 28px;
  height: 44px;
  border: 2px solid var(--primary);
  border-radius: 14px;
  position: relative;
  display: flex;
  justify-content: center;
  padding-top: 8px;
  pointer-events: none; /* Disable pointer events */
`;

const ScrollDot = styled(motion.div)`
  width: 6px;
  height: 6px;
  background-color: var(--primary);
  border-radius: 50%;
  pointer-events: none; /* Disable pointer events */
`;

export const ScrollIndicator = () => {
  const scrollToContent = () => {
    window.scrollTo({
      top: window.innerHeight,
      behavior: 'smooth'
    });
  };

  return (
    <ScrollContainer
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 2, duration: 0.5 }}
      onClick={scrollToContent}
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
    >
      <ScrollText
        animate={{ opacity: [0.4, 0.8, 0.4] }}
        transition={{
          duration: 2,
          repeat: Infinity,
          ease: "easeInOut"
        }}
      >
        Scroll to explore
      </ScrollText>
      <ArrowContainer>
        <ScrollDot
          animate={{
            y: [0, 24, 0],
          }}
          transition={{
            duration: 1.5,
            repeat: Infinity,
            ease: "easeInOut",
          }}
        />
      </ArrowContainer>
    </ScrollContainer>
  );
};