import React, { useRef, useEffect, useCallback } from 'react';
import styled from 'styled-components';
import type { LegoBlock, LegoCanvasProps } from './types';
import { PHASE_TIMINGS, SETTLE } from './constants';
import { BlockGenerator } from './BlockGenerator';
import { PhysicsEngine } from './PhysicsEngine';
import { renderAllBlocks } from './BlockRenderer';
import {
  startDissolve,
  updateDissolve,
  shouldStartDissolve,
  isDissolveComplete,
  hasActiveParticles,
} from './DissolveEffect';
import { useAnimationLoop } from './useLegoAnimation';

const Canvas = styled.canvas`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
`;

export const LegoCanvas: React.FC<LegoCanvasProps> = ({
  config,
  theme,
  isVisible,
  phase,
  onPhaseChange,
  onDissolveStart,
  width,
  height,
}) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const blocksRef = useRef<LegoBlock[]>([]);
  const generatorRef = useRef<BlockGenerator | null>(null);
  const physicsRef = useRef<PhysicsEngine | null>(null);
  const phaseStartTimeRef = useRef<number>(0);
  const animationStartTimeRef = useRef<number>(0);
  const settleStartTimeRef = useRef<number>(0);

  useEffect(() => {
    if (!theme || width <= 0 || height <= 0) return;

    generatorRef.current = new BlockGenerator(theme);
    physicsRef.current = new PhysicsEngine();

    blocksRef.current = generatorRef.current.generatePattern({
      viewportWidth: width,
      viewportHeight: height,
      blockCount: config.blockCount,
      gridCellSize: 32,
    });

    animationStartTimeRef.current = 0;
    phaseStartTimeRef.current = 0;
    settleStartTimeRef.current = 0;

    return () => {
      blocksRef.current = [];
    };
  }, [theme, width, height, config.blockCount]);

  const updateBuildPhase = useCallback((timestamp: number, deltaTime: number) => {
    const elapsedSinceStart = timestamp - animationStartTimeRef.current;

    for (const block of blocksRef.current) {
      if (!block.spawned && elapsedSinceStart >= block.spawnDelay) {
        block.spawned = true;
      }

      if (block.spawned && physicsRef.current) {
        if (!block.settled) {
          physicsRef.current.updateBlock(block, deltaTime, timestamp);
        } else {
          physicsRef.current.updateSettledBlock(block, timestamp);
        }
      }

      for (const particle of block.particles) {
        if (particle.life > 0) {
          particle.x += particle.vx;
          particle.y += particle.vy;
          particle.vy += 0.12;
          particle.vx *= 0.98;
          particle.rotation += particle.rotationSpeed;
          particle.life -= deltaTime;
        }
      }

      block.particles = block.particles.filter(p => p.life > 0);
    }

    const allSpawned = blocksRef.current.every(b => b.spawned);
    const allSettled = physicsRef.current?.isAllSettled(blocksRef.current) ?? false;

    if (allSpawned && allSettled) {
      onPhaseChange('settle');
      settleStartTimeRef.current = timestamp;
    }
  }, [onPhaseChange]);

  const updateSettlePhase = useCallback((timestamp: number) => {
    const settleElapsed = timestamp - settleStartTimeRef.current;
    const pulseProgress = (settleElapsed % SETTLE.glowPulseDuration) / SETTLE.glowPulseDuration;
    const glowIntensity = Math.sin(pulseProgress * Math.PI) * SETTLE.glowMaxIntensity;

    for (const block of blocksRef.current) {
      if (block.settled) {
        block.glowIntensity = glowIntensity;
      }
    }

    if (settleElapsed >= PHASE_TIMINGS.settle) {
      for (const block of blocksRef.current) {
        block.glowIntensity = 0;
      }
      onPhaseChange('dissolve');
      phaseStartTimeRef.current = timestamp;
      onDissolveStart();
    }
  }, [onPhaseChange, onDissolveStart]);

  const updateDissolvePhase = useCallback((timestamp: number, deltaTime: number) => {
    const dissolveElapsed = timestamp - phaseStartTimeRef.current;

    for (const block of blocksRef.current) {
      if (shouldStartDissolve(block, height, dissolveElapsed)) {
        startDissolve(block, config);
      }

      if (block.dissolving) {
        updateDissolve(block, deltaTime, config);
      }
    }

    if (isDissolveComplete(blocksRef.current) && !hasActiveParticles(blocksRef.current)) {
      onPhaseChange('revealed');
    }
  }, [height, config, onPhaseChange]);

  const animate = useCallback((timestamp: number, deltaTime: number) => {
    const canvas = canvasRef.current;
    if (!canvas || !theme) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    if (animationStartTimeRef.current === 0) {
      animationStartTimeRef.current = timestamp;
    }

    const dpr = Math.min(window.devicePixelRatio, 2);
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.clearRect(0, 0, canvas.width / dpr, canvas.height / dpr);

    switch (phase) {
      case 'build':
        updateBuildPhase(timestamp, deltaTime);
        break;
      case 'settle':
        updateSettlePhase(timestamp);
        break;
      case 'dissolve':
        updateDissolvePhase(timestamp, deltaTime);
        break;
      case 'revealed':
        return;
    }

    renderAllBlocks(ctx, blocksRef.current, config, theme);
  }, [phase, theme, config, updateBuildPhase, updateSettlePhase, updateDissolvePhase]);

  useAnimationLoop(
    animate,
    isVisible && phase !== 'revealed',
    config.frameTargetMs
  );

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas || width <= 0 || height <= 0) return;

    const dpr = Math.min(window.devicePixelRatio, 2);
    canvas.width = width * dpr;
    canvas.height = height * dpr;
  }, [width, height]);

  if (width <= 0 || height <= 0) return null;

  return (
    <Canvas
      ref={canvasRef}
      style={{ width: `${width}px`, height: `${height}px` }}
    />
  );
};

export default LegoCanvas;
