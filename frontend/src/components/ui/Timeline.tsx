// src/components/sections/Timeline.tsx
import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';

const TimelineContainer = styled.div`
  position: relative;
  max-width: 1200px;
  margin: 100px auto;
  padding: 20px;
`;

const TimelineLine = styled.div`
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 2px;
  height: 100%;
  background: var(--primary);
`;

const TimelineItem = styled(motion.div)`
  display: flex;
  justify-content: flex-end;
  padding-right: 50%;
  position: relative;
  margin: 50px 0;

  &:nth-child(odd) {
    justify-content: flex-start;
    padding-right: 0;
    padding-left: 50%;
  }
`;

const TimelineContent = styled(motion.div)`
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 20px;
  border-radius: 10px;
  width: 400px;
  position: relative;
`;

const TimelineDot = styled(motion.div)`
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 20px;
  background: var(--primary);
  border-radius: 50%;
`;

interface TimelineEventProps {
  date: string;
  title: string;
  description: string;
}

export const Timeline: React.FC<{ events: TimelineEventProps[] }> = ({ events }) => {
  return (
    <TimelineContainer>
      <TimelineLine />
      {events.map((event, index) => (
        <TimelineItem
          key={index}
          initial={{ opacity: 0, y: 50 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: index * 0.2 }}
          viewport={{ once: true }}
        >
          <TimelineContent>
            <h3>{event.date}</h3>
            <h4>{event.title}</h4>
            <p>{event.description}</p>
          </TimelineContent>
          <TimelineDot
            initial={{ scale: 0 }}
            whileInView={{ scale: 1 }}
            transition={{ duration: 0.5 }}
            viewport={{ once: true }}
          />
        </TimelineItem>
      ))}
    </TimelineContainer>
  );
};
