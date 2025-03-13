/**
 * @deprecated This component has been replaced with CreativeShaderBackground
 * for improved visual quality and performance. Please use CreativeShaderBackground instead.
 * 
 * This file is kept for reference but will be removed in a future update.
 */

import React, { useRef, useEffect, useState, useCallback, useMemo } from 'react';
import styled from 'styled-components';
import { useTheme } from '../../contexts/ThemeContext';
import useDeviceCapabilities from '../../hooks/useDeviceCapabilities';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

// Animation container
const AnimationContainer = styled.div<{ opacity: number }>`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
  opacity: ${props => props.opacity};
  transition: opacity 0.5s ease;
  will-change: transform;
  transform: translateZ(0);
`;

const Canvas = styled.canvas`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: block;
`;

interface Pixel {
  x: number;
  y: number;
  size: number;
  color: string;
  alpha: number;
  targetAlpha: number;
  velocity: number;
  state: 'idle' | 'active' | 'fading';
  activationTime: number;
  lifetime: number;
}

interface Wave {
  x: number;
  y: number;
  radius: number;
  maxRadius: number;
  speed: number;
  opacity: number;
}

interface PixelGridAnimationProps {
  cellSize?: number;
  density?: number;
  highlightColor?: string;
  baseColor?: string;
  animationSpeed?: number;
  interactive?: boolean;
  showWaves?: boolean;
  pulseEffect?: boolean;
}

const PixelGridAnimation: React.FC<PixelGridAnimationProps> = ({
  cellSize = 12,
  density = 0.8,
  highlightColor,
  baseColor,
  animationSpeed = 1,
  interactive = true,
  showWaves = true,
  pulseEffect = true,
}) => {
  console.warn('PixelGridAnimation is deprecated. Please use CreativeShaderBackground instead.');
  return null;
};

export default PixelGridAnimation; 