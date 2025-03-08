import React, { useRef, useEffect } from 'react';
import { motion, useInView, useAnimation, Variants } from 'framer-motion';
import styled from 'styled-components';
import { useTheme } from '../../contexts/ThemeContext';
import { useZIndex } from '../../hooks/useZIndex';

interface Milestone {
  year: string;
  title: string;
  description: string;
}

interface SkillGroup {
  category: string;
  skills: string[];
}

const AboutContainer = styled.div`
  position: relative;
  width: 100%;
  padding: 4rem 2rem;
  
  @media (min-width: 768px) {
    padding: 8rem 4rem;
  }
`;

const SectionHeading = styled(motion.h2)`
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 3rem;
  position: relative;
  display: inline-block;
  
  &::after {
    content: '';
    position: absolute;
    left: 0;
    bottom: -10px;
    height: 3px;
    width: 60%;
    background: ${({ theme }) => theme.colors.accent};
  }
`;

const AboutGrid = styled.div`
  display: grid;
  grid-template-columns: 1fr;
  gap: 4rem;
  margin-top: 2rem;
  
  @media (min-width: 992px) {
    grid-template-columns: 1fr 1fr;
  }
`;

const BioSection = styled(motion.div)`
  display: flex;
  flex-direction: column;
  gap: 2rem;
`;

const BioText = styled(motion.p)`
  font-size: 1.125rem;
  line-height: 1.7;
  margin-bottom: 1.5rem;
  max-width: 650px;
`;

const TimelineSection = styled(motion.div)`
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 2px;
    height: 100%;
    background: ${({ theme }) => theme.colors.accent};
    opacity: 0.3;
  }
`;

const MilestoneItem = styled(motion.div)`
  position: relative;
  padding-left: 2rem;
  margin-bottom: 3rem;
  
  &::before {
    content: '';
    position: absolute;
    left: -6px;
    top: 8px;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: ${({ theme }) => theme.colors.accent};
  }
`;

const MilestoneYear = styled.span`
  font-weight: 600;
  color: ${({ theme }) => theme.colors.accent};
  display: block;
  margin-bottom: 0.5rem;
`;

const MilestoneTitle = styled.h3`
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.75rem;
`;

const MilestoneDescription = styled.p`
  font-size: 1rem;
  line-height: 1.6;
  color: ${({ theme }) => theme.colors.text};
  opacity: 0.7;
`;

const SkillsContainer = styled(motion.div)`
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
  margin-top: 3rem;
  
  @media (min-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }
`;

const SkillGroupContainer = styled(motion.div)`
  margin-bottom: 2rem;
`;

const SkillCategory = styled.h3`
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: ${({ theme }) => theme.colors.text};
`;

const SkillList = styled.ul`
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  padding: 0;
  list-style: none;
`;

const SkillItem = styled(motion.li)`
  padding: 0.5rem 1rem;
  background: ${({ theme }) => theme.colors.backgroundAlt};
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s ease;
  
  &:hover {
    background: ${({ theme }) => theme.colors.accent};
    color: #fff;
    transform: translateY(-2px);
  }
`;

const animationVariants: Record<string, Variants> = {
  heading: {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0, 
      transition: { duration: 0.6, ease: [0.23, 1, 0.32, 1] } 
    }
  },
  bio: {
    hidden: { opacity: 0, y: 30 },
    visible: { 
      opacity: 1, 
      y: 0, 
      transition: { duration: 0.7, ease: [0.23, 1, 0.32, 1] } 
    }
  },
  timeline: {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1, 
      transition: { staggerChildren: 0.2 } 
    }
  },
  milestone: {
    hidden: { opacity: 0, x: -20 },
    visible: { 
      opacity: 1, 
      x: 0, 
      transition: { duration: 0.5, ease: [0.23, 1, 0.32, 1] } 
    }
  },
  skills: {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1, 
      transition: { staggerChildren: 0.1 } 
    }
  },
  skillGroup: {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0, 
      transition: { duration: 0.5, ease: [0.23, 1, 0.32, 1] } 
    }
  },
  skill: {
    hidden: { opacity: 0, scale: 0.9 },
    visible: { 
      opacity: 1, 
      scale: 1,
      transition: { duration: 0.3, ease: [0.23, 1, 0.32, 1] } 
    }
  }
};

export const About: React.FC = () => {
  const { theme } = useTheme();
  const { zIndex } = useZIndex();
  
  const sectionRef = useRef<HTMLDivElement>(null);
  const isInView = useInView(sectionRef, { once: false, amount: 0.2 });
  const controls = useAnimation();
  
  useEffect(() => {
    if (isInView) {
      controls.start('visible');
    }
  }, [isInView, controls]);

  const milestones: Milestone[] = [
    {
      year: '2023',
      title: 'Senior Developer',
      description: 'Led development teams and architected enterprise-scale applications with a focus on performance optimization.'
    },
    {
      year: '2021',
      title: 'Full Stack Developer',
      description: 'Created modern web applications utilizing React, Node.js, and cloud infrastructure.'
    },
    {
      year: '2019',
      title: 'Frontend Specialist',
      description: 'Developed responsive user interfaces and interactive web experiences with cutting-edge technologies.'
    },
    {
      year: '2017',
      title: 'Started Coding Journey',
      description: 'Began learning programming fundamentals and building small projects to develop skills.'
    }
  ];

  const skillGroups: SkillGroup[] = [
    {
      category: 'Frontend',
      skills: ['React', 'TypeScript', 'Next.js', 'CSS/SCSS', 'Framer Motion', 'Redux', 'Styled Components']
    },
    {
      category: 'Backend',
      skills: ['Node.js', 'Express', 'PostgreSQL', 'MongoDB', 'GraphQL', 'REST API Design']
    },
    {
      category: 'DevOps & Tools',
      skills: ['Git', 'Docker', 'AWS', 'CI/CD', 'Jest', 'Webpack', 'Performance Optimization']
    },
    {
      category: 'Design & UX',
      skills: ['Figma', 'UI Design', 'Accessibility', 'User Research', 'Animation', 'Design Systems']
    }
  ];

  return (
    <AboutContainer ref={sectionRef}>
      <SectionHeading
        variants={animationVariants.heading}
        initial="hidden"
        animate={controls}
      >
        About Me
      </SectionHeading>
      
      <AboutGrid>
        <BioSection
          variants={animationVariants.bio}
          initial="hidden"
          animate={controls}
        >
          <BioText>
            I'm a passionate developer with a keen eye for detail and a love for creating beautiful, functional web experiences. My journey in technology began with a curiosity about how digital products shape our daily lives.
          </BioText>
          
          <BioText>
            With expertise in both frontend and backend technologies, I bring a holistic approach to development. I believe in writing clean, maintainable code that scales with your business needs while providing exceptional user experiences.
          </BioText>
          
          <BioText>
            When I'm not coding, you can find me exploring new technologies, contributing to open source projects, or sharing knowledge with the developer community.
          </BioText>
          
          <TimelineSection
            variants={animationVariants.timeline}
            initial="hidden"
            animate={controls}
          >
            {milestones.map((milestone, index) => (
              <MilestoneItem 
                key={index}
                variants={animationVariants.milestone}
              >
                <MilestoneYear>{milestone.year}</MilestoneYear>
                <MilestoneTitle>{milestone.title}</MilestoneTitle>
                <MilestoneDescription>{milestone.description}</MilestoneDescription>
              </MilestoneItem>
            ))}
          </TimelineSection>
        </BioSection>
        
        <SkillsContainer
          variants={animationVariants.skills}
          initial="hidden"
          animate={controls}
        >
          {skillGroups.map((group, groupIndex) => (
            <SkillGroupContainer 
              key={groupIndex}
              variants={animationVariants.skillGroup}
            >
              <SkillCategory>{group.category}</SkillCategory>
              <SkillList>
                {group.skills.map((skill, skillIndex) => (
                  <SkillItem 
                    key={skillIndex}
                    variants={animationVariants.skill}
                  >
                    {skill}
                  </SkillItem>
                ))}
              </SkillList>
            </SkillGroupContainer>
          ))}
        </SkillsContainer>
      </AboutGrid>
    </AboutContainer>
  );
};

export default About; 