import React, { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { useInView } from 'react-intersection-observer';
import Badge from '../../components/common/Badge';
import { useScrollTo } from '../../hooks/useScrollTo';

// Import the image using require
const headshot = require('../../assets/images/headshot.webp');

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
  display: flex;
  flex-direction: column;
  gap: ${({ theme }) => theme.spacing.xxl};
`;

const ProfileSection = styled(motion.div)`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.xl};
  
  @media (min-width: ${({ theme }) => theme.breakpoints.tablet}) {
    flex-direction: row;
    align-items: flex-start;
    justify-content: center;
  }
`;

const ProfileImage = styled(motion.div)`
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: ${({ theme }) => theme.colors.primary};
  overflow: hidden;
  position: relative;
  flex-shrink: 0;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    position: absolute;
    top: 0;
    left: 0;
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    width: 150px;
    height: 150px;
  }
`;

const ProfileInfo = styled.div`
  text-align: center;
  
  @media (min-width: ${({ theme }) => theme.breakpoints.tablet}) {
    text-align: left;
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
  margin-bottom: ${({ theme }) => theme.spacing.md};
`;

const MainContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: ${({ theme }) => theme.spacing.xl};
`;

const ProfileSectionWrapper = styled.div`
  background: ${({ theme }) => theme.colors.surface};
  padding: ${({ theme }) => theme.spacing.xxl};
  border-radius: ${({ theme }) => theme.borderRadius.large};
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
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
  max-width: 100%;
  overflow-wrap: break-word;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 0.95rem;
    line-height: 1.5;
    padding-right: 10px; /* Add some padding to prevent text from touching the edge */
  }
`;

const TypedLine = styled(motion.div)`
  margin-bottom: 16px;
  position: relative;
  display: block;
  min-height: 24px;
  max-width: 100%;
  overflow-wrap: normal;
  word-break: normal;
  white-space: pre-wrap;
`;

const TypewriterText = styled(motion.div)`
  display: inline;
  color: ${({ theme }) => theme.colors.text};
  font-size: 1rem;
  line-height: 1.6;
  opacity: 1;
`;

const WordWrapper = styled.span`
  display: inline-block;
  white-space: nowrap;
  margin-right: 5px;
`;

const TypewriterCharacter = styled(motion.span)`
  display: inline-block;
  position: relative;
  color: ${({ theme }) => theme.colors.text};
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

const CertificationGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: ${({ theme }) => theme.spacing.lg};
`;

const CertificationCard = styled(motion.div)`
  background: ${({ theme }) => theme.colors.background};
  padding: ${({ theme }) => theme.spacing.lg};
  border-radius: ${({ theme }) => theme.borderRadius.medium};
  border: 1px solid ${({ theme }) => theme.colors.border};
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
    border-color: ${({ theme }) => theme.colors.primary};
  }
`;

const CertificationName = styled.h4`
  color: ${({ theme }) => theme.colors.text};
  font-size: 1.1rem;
  margin-bottom: ${({ theme }) => theme.spacing.xs};
`;

const CertificationIssuer = styled.p`
  color: ${({ theme }) => theme.colors.primary};
  font-size: 0.95rem;
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const CertificationDate = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.85rem;
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const CertificationLink = styled.a`
  color: ${({ theme }) => theme.colors.primary};
  text-decoration: none;
  font-size: 0.9rem;
  display: inline-flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.xs};
  
  &:hover {
    text-decoration: underline;
  }
  
  svg {
    width: 16px;
    height: 16px;
  }
`;

const CertificationBadgeContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: ${({ theme }) => theme.spacing.xs};
  margin-top: ${({ theme }) => theme.spacing.sm};
  align-items: center;
`;

const CertificationActions = styled.div`
  display: flex;
  flex-direction: column;
  gap: ${({ theme }) => theme.spacing.sm};
  margin-top: ${({ theme }) => theme.spacing.sm};
`;

// Define animation variants for experience section
const experienceSectionVariants = {
  initial: { opacity: 0 },
  animate: {
    opacity: 1,
    transition: {
      staggerChildren: 0.2
    }
  }
};

const experienceItemVariants = {
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

// Experience section component
const ExperienceSection = () => {
  // Experience lines data
  const experienceLines = useMemo(() => [
    "• Developed and maintained multiple full-stack SaaS applications using React, HTML, CSS, TypeScript, Python, and Go",
    "• Implemented microservices architecture for scalable backend solutions",
    "• Created efficient CI/CD pipelines using GitHub Actions and Docker",
    "• Designed and implemented RESTful APIs and WebSocket services"
  ], []);
  
  // State for typing animation
  const [activeLineIndex, setActiveLineIndex] = useState(0);
  const [typedCharacters, setTypedCharacters] = useState<string[]>(["", "", "", ""]);
  const [showCursor, setShowCursor] = useState(true);
  const animationRef = useRef<NodeJS.Timeout | null>(null);
  
  // Reference for the section element
  const [sectionRef, inView] = useInView({
    triggerOnce: false,
    threshold: 0.2,
    rootMargin: "-50px 0px",
  });
  
  // Clear all timeouts on unmount
  useEffect(() => {
    return () => {
      if (animationRef.current) {
        clearTimeout(animationRef.current);
      }
    };
  }, []);
  
  // Type each character one by one
  const typeCharacter = useCallback((lineIndex: number, charIndex: number) => {
    if (lineIndex >= experienceLines.length) {
      setTimeout(() => setShowCursor(false), 800);
      return;
    }
    
    setActiveLineIndex(lineIndex);
    
    // Get the current line
    const currentLine = experienceLines[lineIndex];
    
    if (charIndex < currentLine.length) {
      // Add next character
      setTypedCharacters(prev => {
        const updated = [...prev];
        updated[lineIndex] = currentLine.substring(0, charIndex + 1);
        return updated;
      });
      
      // Schedule next character with a typing speed of 60 WPM
      const typingDelay = Math.random() * 10 + 40; // 40-50ms per character for ~60 WPM
      animationRef.current = setTimeout(() => {
        typeCharacter(lineIndex, charIndex + 1);
      }, typingDelay);
    } else {
      // Line is complete, move to next line after a pause
      animationRef.current = setTimeout(() => {
        typeCharacter(lineIndex + 1, 0);
      }, 700);
    }
  }, [experienceLines]);
  
  // Start or reset animation when section comes into view
  useEffect(() => {
    if (inView) {
      // Reset animation state
      setActiveLineIndex(0);
      setTypedCharacters(["", "", "", ""]);
      setShowCursor(true);
      
      // Clear any existing timeouts
      if (animationRef.current) {
        clearTimeout(animationRef.current);
      }
      
      // Start typing after a small delay
      animationRef.current = setTimeout(() => {
        typeCharacter(0, 0);
      }, 500);
    }
  }, [inView, typeCharacter]);
  
  return (
    <Section
      variants={experienceSectionVariants}
      initial="initial"
      animate="animate"
      ref={sectionRef}
    >
      <SectionTitle>Experience</SectionTitle>
      <SectionContent>
        <ExperienceItem variants={experienceItemVariants}>
          <ExperienceTitle>Full Stack Developer</ExperienceTitle>
          <ExperienceCompany>Personal Projects & Freelance</ExperienceCompany>
          <ExperienceDate>2020 - Present</ExperienceDate>
          <ExperienceDescription>
            {experienceLines.map((line, lineIndex) => (
              <TypedLine key={`line-${lineIndex}`}>
                <TypewriterText>
                  {/* Process each word as a non-breaking unit */}
                  {typedCharacters[lineIndex].split(' ').map((word, wordIndex) => {
                    // Skip empty words
                    if (word === '') return null;
                    
                    return (
                      <WordWrapper key={`${lineIndex}-word-${wordIndex}`}>
                        {word.split('').map((char, charIndex) => (
                          <TypewriterCharacter
                            key={`${lineIndex}-${wordIndex}-${charIndex}`}
                            initial={{ opacity: 0 }}
                            animate={{ opacity: 1 }}
                            transition={{ duration: 0.1 }}
                          >
                            {char}
                          </TypewriterCharacter>
                        ))}
                      </WordWrapper>
                    );
                  })}
                  {activeLineIndex === lineIndex && showCursor && (
                    <Cursor
                      animate={{ opacity: [1, 0] }}
                      transition={{
                        duration: 0.5,
                        repeat: Infinity,
                        repeatType: "reverse"
                      }}
                    />
                  )}
                </TypewriterText>
              </TypedLine>
            ))}
          </ExperienceDescription>
        </ExperienceItem>
      </SectionContent>
    </Section>
  );
};

const About: React.FC = () => {
  const [certifications, setCertifications] = useState<any[]>([]);
  const [certificationsLoading, setCertificationsLoading] = useState(true);

  // Scroll to top on component mount (if not using the ScrollToTop component)
  const { scrollToTop } = useScrollTo();

  useEffect(() => {
    scrollToTop({ behavior: 'auto' });
  }, [scrollToTop]);

  // Fetch certifications
  useEffect(() => {
    const fetchCertifications = async () => {
      try {
        const response = await fetch('/api/v1/devpanel/public/certifications');
        if (response.ok) {
          const data = await response.json();
          setCertifications(data);
        }
      } catch (error) {
        console.error('Failed to fetch certifications:', error);
      } finally {
        setCertificationsLoading(false);
      }
    };

    fetchCertifications();
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

  return (
    <AboutContainer>
      <ContentWrapper>
        {/* Profile and Certifications Section */}
        <ProfileSectionWrapper>
          <ProfileSection
            initial={{ y: -30, opacity: 0 }}
            animate={{ y: 0, opacity: 1 }}
            transition={{ type: "spring", stiffness: 100 }}
          >
            <ProfileImage
              whileHover={{ scale: 1.05 }}
              transition={{ type: "spring", stiffness: 300 }}
            >
              <img 
                src={headshot} 
                alt="Jaden Razo - Full Stack Developer"
                loading="eager"
              />
            </ProfileImage>
            <ProfileInfo>
              <Name>Jaden Razo</Name>
              <Title>Full Stack Developer</Title>
              <motion.p 
                style={{ color: 'var(--text)', lineHeight: 1.6, marginTop: '1rem' }}
                variants={itemVariants}
              >
                I am a passionate Full Stack Developer with expertise in building scalable web applications
                and microservices. With a strong foundation in both frontend and backend development,
                I specialize in creating efficient, user-friendly solutions that solve real-world problems.
              </motion.p>
            </ProfileInfo>
          </ProfileSection>

          {/* Certifications directly below profile */}
          {!certificationsLoading && certifications.length > 0 && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 }}
              style={{ marginTop: '2rem' }}
            >
              <SectionTitle style={{ marginBottom: '1.5rem' }}>Certifications</SectionTitle>
              <CertificationGrid>
                {certifications.map((cert, index) => (
                  <CertificationCard
                    key={cert.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: 0.1 * index }}
                    whileHover={{ scale: 1.02 }}
                  >
                    <CertificationName>{cert.name}</CertificationName>
                    <CertificationIssuer>{cert.issuer}</CertificationIssuer>
                    <CertificationDate>
                      Issued: {new Date(cert.issue_date).toLocaleDateString('en-US', { 
                        year: 'numeric', 
                        month: 'long' 
                      })}
                      {cert.expiry_date && (
                        <>
                          <br />
                          Expires: {new Date(cert.expiry_date).toLocaleDateString('en-US', { 
                            year: 'numeric', 
                            month: 'long' 
                          })}
                        </>
                      )}
                    </CertificationDate>
                    {cert.credential_id && (
                      <p style={{ fontSize: '0.85rem', color: 'var(--text-secondary)', marginBottom: '0.5rem' }}>
                        Credential ID: {cert.credential_id}
                      </p>
                    )}
                    <CertificationActions>
                      {cert.verification_url && (
                        <CertificationLink 
                          href={cert.verification_url} 
                          target="_blank" 
                          rel="noopener noreferrer"
                        >
                          {cert.verification_text || 'Verify Certificate'}
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6m4-3h6v6m-11 5L21 3" />
                          </svg>
                        </CertificationLink>
                      )}
                      
                      <CertificationBadgeContainer>
                        {cert.is_featured && (
                          <Badge variant="featured">
                            Featured
                          </Badge>
                        )}
                        {cert.category && cert.category.name && (
                          <Badge variant="category">
                            {cert.category.name}
                          </Badge>
                        )}
                        {cert.expiry_date && (() => {
                          const expiryDate = new Date(cert.expiry_date);
                          const now = new Date();
                          const monthsUntilExpiry = (expiryDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24 * 30);
                          
                          if (monthsUntilExpiry < 0) {
                            return (
                              <Badge variant="expired">
                                Expired
                              </Badge>
                            );
                          } else if (monthsUntilExpiry < 3) {
                            return (
                              <Badge variant="expiry">
                                Expires Soon
                              </Badge>
                            );
                          }
                          return null;
                        })()}
                      </CertificationBadgeContainer>
                    </CertificationActions>
                  </CertificationCard>
                ))}
              </CertificationGrid>
            </motion.div>
          )}
        </ProfileSectionWrapper>

        <MainContent>
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
          <ExperienceSection />

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
                  Focused on cloud computing, cybersecurity and distributed systems. Studying for A+, Security+, Network+.
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