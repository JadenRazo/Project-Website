import React, { useState, useEffect, useCallback, useRef } from 'react';
import { motion, HTMLMotionProps, Variants, useAnimation } from 'framer-motion';
import styled from 'styled-components';

const SpaceContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
`;

const Rocket = styled(motion.div)<HTMLMotionProps<"div">>`
  position: absolute;
  width: 40px;
  height: 40px;
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%23FF6B6B'%3E%3Cpath d='M13.13 22.19L11.5 18.36C14.07 17.78 16.54 17 18.9 16.09L13.13 22.19M5.64 12.5L1.81 10.87L7.91 5.1C7 7.46 6.22 9.93 5.64 12.5M21.61 2.39C21.61 2.39 16.66 .269 11 5.93C8.81 8.12 7.5 10.53 6.65 12.64C6.37 13.39 6.56 14.21 7.11 14.77L9.24 16.89C9.79 17.45 10.61 17.63 11.36 17.35C13.5 16.53 15.88 15.19 18.07 13C23.73 7.34 21.61 2.39 21.61 2.39M14.54 9.46C13.76 8.68 13.76 7.41 14.54 6.63S17.59 5.85 18.37 6.63C19.14 7.41 19.15 8.68 18.37 9.46C17.59 10.24 16.32 10.24 15.54 9.46M8.88 16.53L7.47 15.12L8.88 16.53M6.24 22L9.24 19L8.53 17.76L9.47 16.82L10.71 17.53L13.71 14.53L12.47 13.29L13.41 12.35L14.65 13.06L17.65 10.06L14.65 7.06L6.24 22Z'/%3E%3C/svg%3E") no-repeat center;
  background-size: contain;
  filter: drop-shadow(0 0 8px rgba(255, 107, 107, 0.5));
  transform-origin: center;
  transform: rotate(40deg);
`;

const Thruster = styled(motion.div)<HTMLMotionProps<"div">>`
  position: absolute;
  width: 20px;
  height: 35px;
  left: -7px;
  top: 50%;
  transform: translateY(-50%);
  background: radial-gradient(
    ellipse at right,
    rgba(255,107,107,1) 0%,
    rgba(255,107,107,0.9) 15%,
    rgba(255,107,107,0.6) 30%,
    rgba(255,107,107,0.2) 60%,
    rgba(255,107,107,0) 80%
  );
  filter: blur(1px);
  transform-origin: right center;
  z-index: -1;
`;

const ThrusterCore = styled(motion.div)`
  position: absolute;
  width: 5px;
  height: -10px;
  right: -5px;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(255, 255, 255, 0.95);
  filter: blur(1px);
  border-radius: 50%;
`;

const Debris = styled(motion.div)<{ color: string }>`
  position: absolute;
  width: 4px;
  height: 4px;
  background: ${props => props.color};
  border-radius: 50%;
  filter: blur(1px);
  box-shadow: 0 0 6px ${props => props.color};
`;

interface DebrisPiece {
  id: number;
  color: string;
  x: number;
  y: number;
  rotate: number;
  scale: number;
  delay: number;
  velocity: number;
}

interface FlightPath {
  startX: number;
  startY: number;
  endX: number;
  endY: number;
  angle: number;
}

const SpaceAnimation: React.FC = () => {
  const [showExplosion, setShowExplosion] = useState(false);
  const [debris, setDebris] = useState<DebrisPiece[]>([]);
  const rocketControls = useAnimation();
  const explosionPoint = useRef({ x: 0, y: 0 });

  const generateFlightPath = useCallback((): FlightPath => {
    const windowWidth = window.innerWidth;
    const windowHeight = window.innerHeight;
    const rocketSize = 40;
    const margin = 100;

    const baseY = Math.random() * (windowHeight - margin * 2) + margin;
    
    const entry = {
      x: -rocketSize,
      y: baseY,
    };

    const exitY = baseY + (Math.random() * 40 - 20);
    const exitX = windowWidth - rocketSize;

    const angle = 30;

    return {
      startX: entry.x,
      startY: entry.y,
      endX: exitX,
      endY: exitY,
      angle: angle
    };
  }, []);

  const startRocketAnimation = useCallback(async () => {
    const path = generateFlightPath();

    try {
      await rocketControls.start({
        x: path.startX,
        y: path.startY,
        rotate: path.angle,
        scale: 0.8,
        opacity: 0,
        transition: { duration: 0 }
      });

      await rocketControls.start({
        opacity: 1,
        transition: {
          duration: 0.3
        }
      });

      await rocketControls.start({
        x: path.endX,
        y: path.endY,
        rotate: path.angle,
        scale: 1,
        transition: {
          duration: 5,
          ease: "linear"
        }
      });

      explosionPoint.current = { x: path.endX, y: path.endY };
      setShowExplosion(true);

      await rocketControls.start({
        opacity: 0,
        scale: [1, 1.5],
        transition: { 
          duration: 0.15,
          ease: "easeOut"
        }
      });

      setTimeout(() => {
        setShowExplosion(false);
        startRocketAnimation();
      }, 4000);
    } catch (error) {
      console.error("Animation error:", error);
    }
  }, [rocketControls, generateFlightPath]);

  const createDebrisPiece = useCallback((index: number): DebrisPiece => ({
    id: index,
    color: index % 5 === 0 ? '#FFA5A5' : '#FF6B6B',
    x: -(Math.random() * 400 + 100),
    y: (Math.random() - 0.5) * 400,
    rotate: Math.random() * 1440 - 720,
    scale: Math.random() * 0.7 + 0.3,
    delay: Math.random() * 0.1,
    velocity: Math.random() * 1.5 + 0.5,
  }), []);

  useEffect(() => {
    if (showExplosion) {
      const newDebris = Array.from({ length: 100 }, (_, i) => createDebrisPiece(i));
      setDebris(newDebris);
    }
  }, [showExplosion, createDebrisPiece]);

  useEffect(() => {
    startRocketAnimation();
  }, [startRocketAnimation]);

  return (
    <SpaceContainer>
      <Rocket
        animate={rocketControls}
        initial={{ 
          x: -100,
          y: window.innerHeight / 2,
          opacity: 0,
          rotate: 50,
          scale: 0.8
        }}
      >
        <Thruster
          animate={{
            scaleX: [1, 1.4, 1],
            opacity: [0.9, 1, 0.9],
          }}
          transition={{
            duration: 0.2,
            repeat: Infinity,
            repeatType: "reverse",
          }}
        >
          <ThrusterCore
            animate={{
              opacity: [0.8, 1, 0.8],
              scale: [0.9, 1.2, 0.9],
            }}
            transition={{
              duration: 0.15,
              repeat: Infinity,
              repeatType: "reverse",
            }}
          />
        </Thruster>
      </Rocket>
      
      {showExplosion && debris.map((piece) => (
        <Debris
          key={piece.id}
          color={piece.color}
          initial={{
            x: explosionPoint.current.x,
            y: explosionPoint.current.y,
            scale: 0,
            opacity: 0,
          }}
          animate={{
            x: explosionPoint.current.x + piece.x,
            y: explosionPoint.current.y + piece.y,
            scale: [0, piece.scale, 0],
            opacity: [0, 1, 0],
            rotate: piece.rotate,
          }}
          transition={{
            duration: 1.8 / piece.velocity,
            ease: "easeOut",
            delay: piece.delay,
          }}
        />
      ))}
    </SpaceContainer>
  );
};

export default SpaceAnimation; 