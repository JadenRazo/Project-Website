import React, { useRef, useState, useCallback, useEffect } from 'react';
import { motion, useInView, useAnimation } from 'framer-motion';
import styled from 'styled-components';
import { useTheme } from '../../contexts/ThemeContext';
import { useZIndex } from '../../hooks/useZIndex';

interface Milestone {
  year: string;
  title: string;
  description: string;
}

// Main container with stable positioning
const AboutContainer = styled.section`
  position: relative;
  width: 100%;
  padding: 6rem 2rem;
  min-height: 100vh;
  overflow-x: hidden;
  
  @media (min-width: 768px) {
    padding: 8rem 4rem;
  }
`;

const ContentWrapper = styled.div`
  max-width: 1400px;
  margin: 0 auto;
`;

// Section title with underline animation
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
    transform-origin: left center;
  }
`;

// Grid layout for sections
const AboutGrid = styled.div`
  display: grid;
  grid-template-columns: 1fr;
  gap: 4rem;
  margin-top: 2rem;
  
  @media (min-width: 992px) {
    grid-template-columns: 1fr 1fr;
  }
`;

// Bio section with staggered text
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

// Timeline section with visual line
const TimelineSection = styled(motion.div)`
  position: relative;
  padding-left: 2rem;
  
  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    width: 2px;
    height: 100%;
    background: ${({ theme }) => theme.colors.accent};
    opacity: 0.3;
    transform-origin: top center;
  }
`;

const MilestoneItem = styled(motion.div)`
  position: relative;
  margin-bottom: 3rem;
  
  &::before {
    content: '';
    position: absolute;
    left: -2rem;
    top: 8px;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: ${({ theme }) => theme.colors.accent};
    transform-origin: center;
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

// Animation variants with simpler transitions
const animationVariants = {
  // Heading animation
  heading: {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0, 
      transition: { duration: 0.6, ease: "easeOut" } 
    }
  },
  // Bio text animation with stagger
  container: {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1, 
      transition: { staggerChildren: 0.15, delayChildren: 0.1 } 
    }
  },
  // Individual bio paragraph
  item: {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0, 
      transition: { duration: 0.5, ease: "easeOut" } 
    }
  },
  // Timeline wrapper
  timeline: {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1, 
      transition: { staggerChildren: 0.15, delayChildren: 0.2 } 
    }
  },
  // Individual milestone
  milestone: {
    hidden: { opacity: 0, x: -20 },
    visible: { 
      opacity: 1, 
      x: 0, 
      transition: { duration: 0.5, ease: "easeOut" } 
    }
  }
};

export const About: React.FC = () => {
  const { theme } = useTheme();
  const { zIndex } = useZIndex();
  const [isFirstRender, setIsFirstRender] = useState(true);
  
  // Animation controls for better performance
  const sectionControls = useAnimation();
  const bioControls = useAnimation();
  const timelineControls = useAnimation();
  
  // Refs for sections with stable thresholds
  const sectionRef = useRef<HTMLDivElement>(null);
  const bioRef = useRef<HTMLDivElement>(null);
  const timelineRef = useRef<HTMLDivElement>(null);
  
  // InView with optimized settings - removed 'once: true' to allow re-animation
  const isSectionInView = useInView(sectionRef, { amount: 0.2, once: false });
  const isBioInView = useInView(bioRef, { amount: 0.2, once: false });
  const isTimelineInView = useInView(timelineRef, { amount: 0.2, once: false });
  
  // Handle scroll-based animations with performance optimizations
  useEffect(() => {
    // Skip animations on first render to avoid initial jank
    if (isFirstRender) {
      setIsFirstRender(false);
      sectionControls.set("visible");
      bioControls.set("visible");
      timelineControls.set("visible");
      return;
    }
    
    // Animate section heading
    if (isSectionInView) {
      sectionControls.start("visible");
    } else {
      sectionControls.start("hidden");
    }
    
    // Animate bio section
    if (isBioInView) {
      bioControls.start("visible");
    } else {
      bioControls.start("hidden");
    }
    
    // Animate timeline
    if (isTimelineInView) {
      timelineControls.start("visible");
    } else {
      timelineControls.start("hidden");
    }
  }, [
    isFirstRender,
    isSectionInView, 
    isBioInView, 
    isTimelineInView,
    sectionControls,
    bioControls,
    timelineControls
  ]);

  // Define milestone data
  const milestones: Milestone[] = [
    {
      year: '2025',
      title: 'DevOps Trainee',
      description: 'Studying at Western Governors University to become a DevOps Engineer. B.S in Cloud Computing.'
    },
    {
      year: '2023',
      title: 'Linux Server Admin',
      description: 'Developed persistent critical server solutions for clients using Linux and Tmux to achieve high availability and scalability.'
    },
    {
      year: '2021',
      title: 'API Specialist',
      description: 'Specialized in creating responsive, accessible web applications. Specifically accessing APIs and integrating them into the application.'
    },
    {
      year: '2019',
      title: 'Backend Developer',
      description: 'Expanded skills to cover backend technologies. Specifically focused on database management and server-side logic.'
    },
    {
      year: '2017',
      title: 'Junior Developer',
      description: 'Started as a hobbyist as a mod developer focused on Java.'
    }
  ];

  return (
    <AboutContainer ref={sectionRef} id="about">
      <ContentWrapper>
        <SectionHeading
          variants={animationVariants.heading}
          initial="hidden"
          animate={sectionControls}
          aria-label="About Me Section"
        >
          About Me
        </SectionHeading>
        
        <AboutGrid>
          <BioSection
            ref={bioRef}
            variants={animationVariants.container}
            initial="hidden"
            animate={bioControls}
          >
            <BioText variants={animationVariants.item}>
              I'm a passionate web developer and designer with a focus on creating
              engaging digital experiences. With expertise in modern frameworks and
              a keen eye for design, I bridge the gap between functionality and aesthetics.
            </BioText>
            
            <BioText variants={animationVariants.item}>
              My approach combines technical knowledge with creative problem-solving,
              resulting in performant and visually appealing solutions. I'm constantly
              exploring new technologies and techniques to enhance my craft.
            </BioText>
          </BioSection>
          
          <TimelineSection
            ref={timelineRef}
            variants={animationVariants.timeline}
            initial="hidden"
            animate={timelineControls}
          >
            {milestones.map((milestone, index) => (
              <MilestoneItem 
                key={`${milestone.year}-${index}`}
                variants={animationVariants.milestone}
              >
                <MilestoneYear>{milestone.year}</MilestoneYear>
                <MilestoneTitle>{milestone.title}</MilestoneTitle>
                <MilestoneDescription>{milestone.description}</MilestoneDescription>
              </MilestoneItem>
            ))}
          </TimelineSection>
        </AboutGrid>
      </ContentWrapper>
    </AboutContainer>
  );
};

export default About; 