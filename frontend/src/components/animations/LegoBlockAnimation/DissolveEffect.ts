import type { LegoBlock, Particle, PerformanceConfig } from './types';
import { DISSOLVE } from './constants';

const createParticle = (
  block: LegoBlock,
  index: number,
  total: number
): Particle => {
  const angle = (index / total) * Math.PI * 2 + Math.random() * 0.5;
  const speed = 2 + Math.random() * 3;

  return {
    x: block.x + block.width / 2 + (Math.random() - 0.5) * block.width * 0.5,
    y: block.y + block.height / 2 + (Math.random() - 0.5) * block.height * 0.5,
    vx: Math.cos(angle) * speed * (0.5 + Math.random() * 0.5),
    vy: Math.sin(angle) * speed * (0.5 + Math.random() * 0.5) - 2,
    life: DISSOLVE.particleDuration,
    maxLife: DISSOLVE.particleDuration,
    size: 4 + Math.random() * 6,
    color: block.color.start,
    rotation: Math.random() * Math.PI * 2,
    rotationSpeed: (Math.random() - 0.5) * 0.3,
  };
};

export const startDissolve = (
  block: LegoBlock,
  config: PerformanceConfig
): void => {
  if (block.dissolving) return;

  block.dissolving = true;
  block.particles = [];

  if (!config.dissolveSimplified) {
    for (let i = 0; i < config.particlesPerBlock; i++) {
      block.particles.push(createParticle(block, i, config.particlesPerBlock));
    }
  }
};

export const updateParticles = (
  block: LegoBlock,
  deltaTime: number
): void => {
  for (const particle of block.particles) {
    if (particle.life <= 0) continue;

    particle.vy += DISSOLVE.particleGravity * (deltaTime / 16);

    particle.vx *= DISSOLVE.particleDrag;
    particle.vy *= DISSOLVE.particleDrag;

    particle.x += particle.vx * (deltaTime / 16);
    particle.y += particle.vy * (deltaTime / 16);

    particle.rotation += particle.rotationSpeed * (deltaTime / 16);

    particle.life -= deltaTime;

    particle.size *= 0.995;
  }

  block.particles = block.particles.filter(p => p.life > 0);
};

export const updateDissolve = (
  block: LegoBlock,
  deltaTime: number,
  config: PerformanceConfig
): void => {
  if (!block.dissolving) return;

  if (config.dissolveSimplified) {
    block.opacity -= 0.03 * (deltaTime / 16);
  } else {
    block.opacity -= 0.02 * (deltaTime / 16);
    updateParticles(block, deltaTime);
  }

  block.opacity = Math.max(0, block.opacity);
};

export const calculateDissolveProgress = (
  block: LegoBlock,
  canvasHeight: number,
  elapsedTime: number
): number => {
  const normalizedY = block.y / canvasHeight;
  const timeProgress = elapsedTime * DISSOLVE.sweepSpeed;
  return Math.max(0, timeProgress - normalizedY);
};

export const shouldStartDissolve = (
  block: LegoBlock,
  canvasHeight: number,
  elapsedTime: number
): boolean => {
  if (block.dissolving || !block.settled) return false;
  const progress = calculateDissolveProgress(block, canvasHeight, elapsedTime);
  return progress > 0;
};

export const isDissolveComplete = (blocks: LegoBlock[]): boolean => {
  return blocks.every(block => !block.spawned || block.opacity <= 0);
};

export const hasActiveParticles = (blocks: LegoBlock[]): boolean => {
  return blocks.some(block => block.particles.length > 0);
};
