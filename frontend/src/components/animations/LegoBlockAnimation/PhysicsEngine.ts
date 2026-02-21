import type { LegoBlock, Particle } from './types';
import { ANIMATION } from './constants';

const createLandingParticle = (block: LegoBlock, index: number): Particle => {
  const angle = Math.PI * 0.3 + Math.random() * Math.PI * 0.4;
  const speed = 1.5 + Math.random() * 2;
  const side = index % 2 === 0 ? -1 : 1;

  return {
    x: block.x + block.width / 2 + side * (block.width / 4) * Math.random(),
    y: block.targetY + block.height - 2,
    vx: Math.cos(angle) * speed * side,
    vy: -Math.sin(angle) * speed * 0.6,
    life: 150 + Math.random() * 100,
    maxLife: 250,
    size: 2 + Math.random() * 3,
    color: block.color.start,
    rotation: Math.random() * Math.PI * 2,
    rotationSpeed: (Math.random() - 0.5) * 0.2,
  };
};

export class PhysicsEngine {
  updateBlock(block: LegoBlock, deltaTime: number, timestamp: number): void {
    if (block.settled || !block.spawned) return;

    const dt = Math.min(deltaTime / 16, 2);

    if (block.bouncePhase > 0) {
      this.updateBounce(block, timestamp);
      return;
    }

    block.velocity += ANIMATION.gravity * dt;
    block.velocity = Math.min(block.velocity, ANIMATION.maxFallSpeed);

    block.y += block.velocity * dt;

    if (block.y >= block.targetY) {
      block.y = block.targetY;

      if (block.velocity > 3) {
        block.bouncePhase = timestamp;
        block.scaleY = 0.85;
        block.scaleX = 1.1;

        for (let i = 0; i < 3; i++) {
          block.particles.push(createLandingParticle(block, i));
        }
      } else {
        this.settleBlock(block, timestamp);
      }
    }
  }

  private updateBounce(block: LegoBlock, timestamp: number): void {
    const elapsed = timestamp - block.bouncePhase;
    const progress = Math.min(elapsed / ANIMATION.bounceDuration, 1);

    const bounceHeight = Math.sin(progress * Math.PI) * 4 * ANIMATION.bounceStrength;
    block.y = block.targetY - bounceHeight;

    const squashProgress = Math.sin(progress * Math.PI);
    block.scaleY = 1 - squashProgress * 0.1;
    block.scaleX = 1 + squashProgress * 0.05;

    if (progress >= 1) {
      this.settleBlock(block, timestamp);
    }
  }

  private settleBlock(block: LegoBlock, timestamp: number): void {
    block.y = block.targetY;
    block.velocity = 0;
    block.settled = true;
    block.settleStartTime = timestamp;
    block.scaleX = 1;
    block.scaleY = 1;
    block.bouncePhase = 0;
  }

  updateSettledBlock(block: LegoBlock, timestamp: number): void {
    if (!block.settled || block.settleStartTime === 0) return;

    const elapsed = timestamp - block.settleStartTime;
    if (elapsed < 100) {
      const progress = elapsed / 100;
      block.scaleX = 1 + (1 - progress) * 0.02;
      block.scaleY = 1 - (1 - progress) * 0.02;
    } else {
      block.scaleX = 1;
      block.scaleY = 1;
    }
  }

  isAllSettled(blocks: LegoBlock[]): boolean {
    return blocks.every(block => !block.spawned || block.settled);
  }
}

export const createPhysicsEngine = (): PhysicsEngine => {
  return new PhysicsEngine();
};
