import React from 'react';
import { motion } from 'framer-motion';
import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';

const NotFoundContainer = styled.div`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: ${({ theme }) => theme.colors.background};
  overflow: hidden;
  position: relative;
`;

const Content = styled.div`
  text-align: center;
  z-index: 2;
  padding: ${({ theme }) => theme.spacing.xl};
`;

const Title = styled(motion.h1)`
  font-size: 8rem;
  color: ${({ theme }) => theme.colors.primary};
  margin: 0;
  line-height: 1;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    font-size: 6rem;
  }
`;

const Subtitle = styled(motion.p)`
  font-size: 1.5rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  margin: ${({ theme }) => theme.spacing.md} 0;
`;

const Button = styled(motion.button)`
  background: ${({ theme }) => theme.colors.primary};
  color: ${({ theme }) => theme.colors.background};
  border: none;
  padding: ${({ theme }) => `${theme.spacing.sm} ${theme.spacing.lg}`};
  border-radius: ${({ theme }) => theme.borderRadius.medium};
  font-size: 1.1rem;
  cursor: pointer;
  transition: transform 0.2s ease;

  &:hover {
    transform: scale(1.05);
  }
`;

const Star = styled(motion.div)`
  position: absolute;
  width: 2px;
  height: 2px;
  background: ${({ theme }) => theme.colors.primary};
  border-radius: 50%;
`;

const NotFound: React.FC = () => {
  const navigate = useNavigate();
  const numberOfStars = 50;

  const starVariants = {
    animate: (i: number) => ({
      y: [0, -2000],
      opacity: [0, 1, 0],
      transition: {
        duration: Math.random() * 2 + 3,
        repeat: Infinity,
        delay: i * 0.1,
      },
    }),
  };

  return (
    <NotFoundContainer>
      {[...Array(numberOfStars)].map((_, i) => (
        <Star
          key={i}
          variants={starVariants}
          custom={i}
          animate="animate"
          style={{
            left: `${Math.random() * 100}%`,
            top: '100%',
          }}
        />
      ))}
      
      <Content>
        <Title
          initial={{ y: -100, opacity: 0 }}
          animate={{ y: 0, opacity: 1 }}
          transition={{ duration: 0.8, type: "spring" }}
        >
          404
        </Title>
        <Subtitle
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          Looks like you've ventured into deep space
        </Subtitle>
        <Button
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ delay: 0.8 }}
          onClick={() => navigate('/')}
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
        >
          Return to Earth
        </Button>
      </Content>
    </NotFoundContainer>
  );
};

export default NotFound; 