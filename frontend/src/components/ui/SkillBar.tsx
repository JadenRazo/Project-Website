// src/components/ui/SkillBar.tsx
import React from 'react';
import styled from 'styled-components';
import { motion, useScroll, useTransform } from 'framer-motion';

const SkillContainer = styled.div`
  width: 100%;
  margin: 20px 0;
`;

const SkillName = styled.h4`
  margin: 0 0 5px 0;
  color: white;
`;

const ProgressBarContainer = styled.div`
  width: 100%;
  height: 10px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 5px;
  overflow: hidden;
`;

const Progress = styled(motion.div)`
  height: 100%;
  background: linear-gradient(90deg, var(--primary), var(--secondary));
  border-radius: 5px;
`;

interface SkillBarProps {
  skill: string;
  percentage: number;
}

export const SkillBar: React.FC<SkillBarProps> = ({ skill, percentage }) => {
  return (
    <SkillContainer>
      <SkillName>{skill}</SkillName>
      <ProgressBarContainer>
        <Progress
          initial={{ width: 0 }}
          whileInView={{ width: `${percentage}%` }}
          transition={{ duration: 1.5, ease: "easeOut" }}
          viewport={{ once: true }}
        />
      </ProgressBarContainer>
    </SkillContainer>
  );
};
