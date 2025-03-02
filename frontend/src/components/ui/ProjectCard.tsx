import React from 'react';
import styled from 'styled-components';
import { motion, useMotionValue, useTransform, useAnimation } from 'framer-motion';

const Card = styled(motion.div)`
  width: 300px;
  height: 400px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 20px;
  position: relative;
  overflow: hidden;
  backdrop-filter: blur(10px);
  cursor: pointer;
  user-select: none; /* Prevent text selection */
`;

const ProjectImage = styled.img`
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 10px;
`;

const ProjectContent = styled.div`
  margin-top: 20px;
  color: white;
`;

interface ProjectCardProps {
  title: string;
  description: string;
  image: string;
  link: string;
}

export const ProjectCard: React.FC<ProjectCardProps> = ({ title, description, image, link }) => {
  const x = useMotionValue(0);
  const y = useMotionValue(0);
  const controls = useAnimation();

  const rotateX = useTransform(y, [-100, 100], [10, -10]);
  const rotateY = useTransform(x, [-100, 100], [-10, 10]);

  const handleMouseMove = (event: React.MouseEvent<HTMLDivElement>) => {
    const rect = event.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    x.set(event.clientX - centerX);
    y.set(event.clientY - centerY);
  };

  const handleMouseLeave = () => {
    controls.start({ rotateX: 0, rotateY: 0, transition: { duration: 0.5 } });
  };

  const handleTouchMove = (event: React.TouchEvent<HTMLDivElement>) => {
    const rect = event.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    x.set(event.touches[0].clientX - centerX);
    y.set(event.touches[0].clientY - centerY);
  };

  const handleTouchEnd = () => {
    controls.start({ rotateX: 0, rotateY: 0, transition: { duration: 0.5 } });
  };

  return (
    <Card
      whileHover={{ scale: 1.05 }}
      style={{ rotateX, rotateY, perspective: 1000 }}
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
      animate={controls}
    >
      <ProjectImage src={image} alt={title} />
      <ProjectContent>
        <h3>{title}</h3>
        <p>{description}</p>
      </ProjectContent>
    </Card>
  );
};