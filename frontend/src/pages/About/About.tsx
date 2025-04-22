import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { motion, useAnimation } from 'framer-motion';
import { useInView } from 'react-intersection-observer';

// Import the image using require
const headshot = require('../../assets/images/headshot.jpg');

const AboutContainer = styled.div`
  min-height: calc(100vh - 200px);
  padding: ${({ theme }) => theme.spacing.xxl} ${({ theme }) => theme.spacing.xl};
  background: ${({ theme }) => theme.colors.background};
  margin-top: 60px;
  overflow-x: hidden; /* Prevent horizontal scrolling */

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: ${({ theme }) => theme.spacing.xl} ${({ theme }) => theme.spacing.md};
  }
`;

const ContentWrapper = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr;
  gap: ${({ theme }) => theme.spacing.xl};
  
  @media (min-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: 1fr 2fr;
  }
`;

const ProfileSection = styled(motion.div)`
  text-align: center;
  
  @media (min-width: ${({ theme }) => theme.breakpoints.tablet}) {
    text-align: left;
    position: sticky;
    top: 80px;
    height: fit-content;
  }
`;

const ProfileImage = styled(motion.div)`
  width: 200px;
  height: 200px;
  border-radius: 50%;
  margin: 0 auto;
  background: ${({ theme }) => theme.colors.primary};
  overflow: hidden;
  margin-bottom: ${({ theme }) => theme.spacing.lg};
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  
  @media (min-width: ${({ theme }) => theme.breakpoints.tablet}) {
    margin: 0 0 ${({ theme }) => theme.spacing.lg} 0;
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    width: 150px;
    height: 150px;
  }
`;

const Name = styled(motion.h1)`
  color: ${({ theme }) => theme.colors.primary};
  font-size: 2.5rem;
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const Title = styled(motion.h2)`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 1.2rem;
  margin-bottom: ${({ theme }) => theme.spacing.lg};
`;

const MainContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: ${({ theme }) => theme.spacing.xl};
`;

const Section = styled(motion.section)`
  background: ${({ theme }) => theme.colors.surface};
  padding: ${({ theme }) => theme.spacing.xl};
  border-radius: ${({ theme }) => theme.borderRadius.large};
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
`;

const SectionTitle = styled.h3`
  color: ${({ theme }) => theme.colors.primary};
  font-size: 1.5rem;
  margin-bottom: ${({ theme }) => theme.spacing.lg};
  border-bottom: 2px solid ${({ theme }) => theme.colors.border};
  padding-bottom: ${({ theme }) => theme.spacing.sm};
`;

const SectionContent = styled.div`
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.1rem;
  line-height: 1.6;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 0.95rem;
    word-break: break-word; /* Prevent text overflow on mobile */
  }
`;

const SkillsList = styled.ul`
  list-style: none;
  padding: 0;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: ${({ theme }) => theme.spacing.md};
`;

const SkillItem = styled(motion.li)`
  display: flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.sm};
  padding: ${({ theme }) => theme.spacing.sm};
  background: ${({ theme }) => theme.colors.background};
  border-radius: ${({ theme }) => theme.borderRadius.small};
  font-weight: 500;
  color: ${({ theme }) => theme.colors.primary}; /* Brightened skill text color */
  
  svg {
    width: 20px;
    height: 20px;
    color: ${({ theme }) => theme.colors.primary};
  }

  &:hover {
    background: ${({ theme }) => `${theme.colors.primary}15`};
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  }
`;

const ExperienceItem = styled(motion.div)`
  margin-bottom: ${({ theme }) => theme.spacing.lg};
  
  &:last-child {
    margin-bottom: 0;
  }
`;

const ExperienceTitle = styled.h4`
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.2rem;
  margin-bottom: ${({ theme }) => theme.spacing.xs};
`;

const ExperienceCompany = styled.h5`
  color: ${({ theme }) => theme.colors.primary};
  font-size: 1.1rem;
  margin-bottom: ${({ theme }) => theme.spacing.xs};
`;

const ExperienceDate = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.9rem;
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const ExperienceDescription = styled.div`
  color: ${({ theme }) => theme.colors.text};
  line-height: 1.6;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 0.95rem;
  }
`;

const TypedLine = styled(motion.div)`
  margin-bottom: 16px;
  position: relative;
  display: flex;
  min-height: 24px;
`;

const TypewriterText = styled(motion.span)`
  display: inline-block;
  white-space: pre-wrap;
  word-break: break-word;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1rem;
  line-height: 1.6;
  opacity: 1;
`;

const Cursor = styled(motion.span)`
  display: inline-block;
  width: 2px;
  height: 1em;
  background-color: ${({ theme }) => theme.colors.primary};
  margin-left: 2px;
  align-self: center;
  border-radius: 1px;
`;

const About: React.FC = () => {
  // Scroll to top on component mount (if not using the ScrollToTop component)
  useEffect(() => {
    if (typeof window !== 'undefined') {
      window.scrollTo(0, 0);
    }
  }, []);

  // Animation variants for staggered animations
  const containerVariants = {
    initial: { opacity: 0 },
    animate: {
      opacity: 1,
      transition: {
        staggerChildren: 0.2
      }
    }
  };

  const itemVariants = {
    initial: { y: 20, opacity: 0 },
    animate: {
      y: 0,
      opacity: 1,
      transition: {
        type: "spring",
        stiffness: 100
      }
    }
  };

  // Improved typing animation with simpler and more reliable approach
  const [experienceRef, experienceInView] = useInView({
    triggerOnce: false,
    threshold: 0.3,
    rootMargin: "-50px 0px",
  });

  // Experience lines data
  const experienceLines = [
    "• Developed and maintained multiple full-stack SaaS applications using React, HTML, CSS, TypeScript, Python, and Go",
    "• Implemented microservices architecture for scalable backend solutions",
    "• Created efficient CI/CD pipelines using GitHub Actions and Docker",
    "• Designed and implemented RESTful APIs and WebSocket services"
  ];
  
  // State for letter counts in each line
  const [letterCounts, setLetterCounts] = useState<number[]>([0, 0, 0, 0]);
  const [activeLineIndex, setActiveLineIndex] = useState(0);
  const [isTypingComplete, setIsTypingComplete] = useState(false);
  const animationRef = useRef<NodeJS.Timeout | null>(null);
  
  // Clear all timeouts on unmount
  useEffect(() => {
    return () => {
      if (animationRef.current) {
        clearTimeout(animationRef.current);
      }
    };
  }, []);

  // Reset and start animation when section comes into view
  useEffect(() => {
    // When component comes into view
    if (experienceInView) {
      // Reset animation state
      setLetterCounts([0, 0, 0, 0]);
      setActiveLineIndex(0);
      setIsTypingComplete(false);
      
      // Start typing animation for first line
      if (animationRef.current) {
        clearTimeout(animationRef.current);
      }
      
      // Begin typing first line
      animationRef.current = setTimeout(() => {
        typeNextLetter(0);
      }, 400);
    } else {
      // Reset when out of view
      setLetterCounts([0, 0, 0, 0]);
      setActiveLineIndex(0);
      setIsTypingComplete(false);
      
      if (animationRef.current) {
        clearTimeout(animationRef.current);
        animationRef.current = null;
      }
    }
  }, [experienceInView]);
  
  // Type out letters one by one
  const typeNextLetter = (lineIndex: number) => {
    if (lineIndex >= experienceLines.length) {
      setIsTypingComplete(true);
      return;
    }
    
    setActiveLineIndex(lineIndex);
    
    setLetterCounts(prev => {
      const newCounts = [...prev];
      
      if (newCounts[lineIndex] < experienceLines[lineIndex].length) {
        // Type next letter
        newCounts[lineIndex]++;
        
        // Schedule next letter
        animationRef.current = setTimeout(
          () => typeNextLetter(lineIndex),
          Math.random() * 30 + 20 // Random typing speed for realism
        );
      } else {
        // Current line complete, move to next after a pause
        animationRef.current = setTimeout(
          () => typeNextLetter(lineIndex + 1),
          700 // Pause between lines
        );
      }
      
      return newCounts;
    });
  };

  return (
    <AboutContainer>
      <ContentWrapper>
        <ProfileSection
          initial={{ x: -50, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          transition={{ type: "spring", stiffness: 100 }}
        >
          <ProfileImage
            whileHover={{ scale: 1.05 }}
            transition={{ type: "spring", stiffness: 300 }}
          >
            <img src={headshot} alt="Jaden Razo" />
          </ProfileImage>
          <Name>Jaden Razo</Name>
          <Title>Full Stack Developer</Title>
        </ProfileSection>

        <MainContent>
          {/* About Me Section */}
          <Section
            variants={containerVariants}
            initial="initial"
            animate="animate"
          >
            <SectionTitle>About Me</SectionTitle>
            <SectionContent>
              <motion.p variants={itemVariants}>
                I am a passionate Full Stack Developer with expertise in building scalable web applications
                and microservices. With a strong foundation in both frontend and backend development,
                I specialize in creating efficient, user-friendly solutions that solve real-world problems.
              </motion.p>
            </SectionContent>
          </Section>

          {/* Technical Skills Section */}
          <Section
            variants={containerVariants}
            initial="initial"
            animate="animate"
          >
            <SectionTitle>Technical Skills</SectionTitle>
            <SkillsList>
              {[
                "JavaScript/TypeScript",
                "React.js",
                "Node.js",
                "Go",
                "Python",
                "Docker",
                "Kubernetes",
                "AWS",
                "MongoDB",
                "PostgreSQL",
                "Redis",
                "GraphQL"
              ].map((skill, index) => (
                <SkillItem
                  key={skill}
                  variants={itemVariants}
                  whileHover={{ scale: 1.05 }}
                >
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
                  </svg>
                  {skill}
                </SkillItem>
              ))}
            </SkillsList>
          </Section>

          {/* Experience Section with typing animation */}
          <Section
            variants={containerVariants}
            initial="initial"
            animate="animate"
            ref={experienceRef}
          >
            <SectionTitle>Experience</SectionTitle>
            <SectionContent>
              <ExperienceItem variants={itemVariants}>
                <ExperienceTitle>Full Stack Developer</ExperienceTitle>
                <ExperienceCompany>Personal Projects & Freelance</ExperienceCompany>
                <ExperienceDate>2020 - Present</ExperienceDate>
                <ExperienceDescription>
                  {experienceLines.map((line, index) => (
                    <TypedLine key={index}>
                      <TypewriterText>
                        {line.substring(0, letterCounts[index])}
                      </TypewriterText>
                      {activeLineIndex === index && letterCounts[index] < line.length && (
                        <Cursor
                          animate={{ opacity: [1, 0] }}
                          transition={{
                            duration: 0.5,
                            repeat: Infinity,
                            repeatType: "reverse"
                          }}
                        />
                      )}
                    </TypedLine>
                  ))}
                </ExperienceDescription>
              </ExperienceItem>
            </SectionContent>
          </Section>

          {/* Education Section */}
          <Section
            variants={containerVariants}
            initial="initial"
            animate="animate"
          >
            <SectionTitle>Education</SectionTitle>
            <SectionContent>
              <ExperienceItem variants={itemVariants}>
                <ExperienceTitle>B.S Cloud Computing</ExperienceTitle>
                <ExperienceCompany>Western Governors University</ExperienceCompany>
                <ExperienceDate>Present</ExperienceDate>
                <ExperienceDescription>
                  Focused on software engineering, cloud computing, and distributed systems. Focused on DevOps and earning my certifications.
                </ExperienceDescription>
              </ExperienceItem> 
              <ExperienceItem variants={itemVariants}>
                <ExperienceTitle>High School Diploma</ExperienceTitle>
                <ExperienceCompany>Sky Mountain High School</ExperienceCompany>
                <ExperienceDate>2017 - 2021</ExperienceDate>
                <ExperienceDescription>
                  Graduated with honors while participating in coding camps and volunteering projects. Developed strong foundation in programming fundamentals and problem-solving skills.
                </ExperienceDescription>
              </ExperienceItem>
            </SectionContent>
          </Section>
        </MainContent>
      </ContentWrapper>
    </AboutContainer>
  );
};

export default About; 