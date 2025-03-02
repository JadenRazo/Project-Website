// src/components/sections/Projects.tsx
import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';

const ProjectsSection = styled.section`
  padding: 100px 150px;
  max-width: 1600px;
  margin: 0 auto;
`;

const ProjectGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 50px;
`;

export const Projects = () => {
  return (
    <ProjectsSection>
      <h2>Selected Work</h2>
      <ProjectGrid>
        {/* Add your ProjectCard components here */}
      </ProjectGrid>
    </ProjectsSection>
  );
};
