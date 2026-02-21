import type { Theme } from '../../../styles/theme.types';
import type { LegoBlock, PerformanceConfig, Particle } from './types';
import { GRID } from './constants';

const roundedRect = (
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  height: number,
  radius: number
): void => {
  ctx.beginPath();
  ctx.moveTo(x + radius, y);
  ctx.lineTo(x + width - radius, y);
  ctx.quadraticCurveTo(x + width, y, x + width, y + radius);
  ctx.lineTo(x + width, y + height - radius);
  ctx.quadraticCurveTo(x + width, y + height, x + width - radius, y + height);
  ctx.lineTo(x + radius, y + height);
  ctx.quadraticCurveTo(x, y + height, x, y + height - radius);
  ctx.lineTo(x, y + radius);
  ctx.quadraticCurveTo(x, y, x + radius, y);
  ctx.closePath();
};

const renderStuds = (
  ctx: CanvasRenderingContext2D,
  block: LegoBlock,
  theme: Theme
): void => {
  ctx.save();

  for (const stud of block.studs) {
    ctx.beginPath();
    ctx.arc(stud.x, stud.y, GRID.studRadius, 0, Math.PI * 2);

    const gradient = ctx.createRadialGradient(
      stud.x - 1,
      stud.y - 1,
      0,
      stud.x,
      stud.y,
      GRID.studRadius
    );
    gradient.addColorStop(0, 'rgba(255,255,255,0.4)');
    gradient.addColorStop(0.5, block.color.start);
    gradient.addColorStop(1, block.color.end);
    ctx.fillStyle = gradient;
    ctx.fill();

    ctx.strokeStyle = 'rgba(0,0,0,0.2)';
    ctx.lineWidth = 1;
    ctx.stroke();
  }

  ctx.restore();
};

const renderBlockContent = (
  ctx: CanvasRenderingContext2D,
  block: LegoBlock
): void => {
  if (block.content.type === 'empty' || !block.content.value) return;

  ctx.save();

  const centerX = block.width / 2;
  const centerY = block.height / 2 + 2;

  ctx.fillStyle = 'rgba(255,255,255,0.9)';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';

  let fontSize: number;
  const text = block.content.value;

  switch (block.content.type) {
    case 'icon':
      fontSize = Math.min(block.width, block.height) * 0.35;
      ctx.font = `bold ${fontSize}px "SF Mono", "Fira Code", monospace`;
      break;
    case 'code':
    case 'terminal':
      fontSize = Math.min(block.width, block.height) * 0.28;
      ctx.font = `${fontSize}px "SF Mono", "Fira Code", monospace`;
      break;
    case 'symbol':
      fontSize = Math.min(block.width, block.height) * 0.4;
      ctx.font = `bold ${fontSize}px "SF Mono", "Fira Code", monospace`;
      break;
    default:
      fontSize = 10;
      ctx.font = `${fontSize}px sans-serif`;
  }

  ctx.shadowColor = 'rgba(0,0,0,0.3)';
  ctx.shadowBlur = 2;
  ctx.shadowOffsetX = 1;
  ctx.shadowOffsetY = 1;

  ctx.fillText(text, centerX, centerY);

  ctx.restore();
};

const renderGlow = (
  ctx: CanvasRenderingContext2D,
  block: LegoBlock
): void => {
  if (block.glowIntensity <= 0) return;

  ctx.save();
  ctx.shadowColor = block.color.start;
  ctx.shadowBlur = 15 * block.glowIntensity;
  ctx.fillStyle = `rgba(255,255,255,${0.08 * block.glowIntensity})`;
  roundedRect(ctx, 0, 0, block.width, block.height, GRID.blockBorderRadius);
  ctx.fill();
  ctx.restore();
};

export const renderBlock = (
  ctx: CanvasRenderingContext2D,
  block: LegoBlock,
  config: PerformanceConfig,
  theme: Theme
): void => {
  if (!block.spawned || block.opacity <= 0) return;

  ctx.save();

  const centerX = block.x + block.width / 2;
  const centerY = block.y + block.height / 2;

  ctx.translate(centerX, centerY);
  ctx.rotate(block.rotation);
  ctx.scale(block.scaleX, block.scaleY);
  ctx.translate(-block.width / 2, -block.height / 2);

  ctx.globalAlpha = block.opacity;

  if (config.enableGlow && block.glowIntensity > 0) {
    renderGlow(ctx, block);
  }

  if (config.enableGradients) {
    const gradient = ctx.createLinearGradient(0, 0, block.width, block.height);
    gradient.addColorStop(0, block.color.start);
    gradient.addColorStop(1, block.color.end);
    ctx.fillStyle = gradient;
  } else {
    ctx.fillStyle = block.color.start;
  }

  roundedRect(ctx, 0, 0, block.width, block.height, GRID.blockBorderRadius);
  ctx.fill();

  ctx.strokeStyle = 'rgba(0,0,0,0.15)';
  ctx.lineWidth = 1;
  ctx.stroke();

  const highlightGradient = ctx.createLinearGradient(0, 0, 0, block.height * 0.4);
  highlightGradient.addColorStop(0, 'rgba(255,255,255,0.2)');
  highlightGradient.addColorStop(1, 'rgba(255,255,255,0)');
  ctx.fillStyle = highlightGradient;
  roundedRect(ctx, 1, 1, block.width - 2, block.height * 0.35, GRID.blockBorderRadius - 1);
  ctx.fill();

  if (config.enableStudDetail) {
    renderStuds(ctx, block, theme);
  }

  renderBlockContent(ctx, block);

  ctx.restore();
};

export const renderParticle = (
  ctx: CanvasRenderingContext2D,
  particle: Particle
): void => {
  if (particle.life <= 0) return;

  const alpha = particle.life / particle.maxLife;

  ctx.save();
  ctx.translate(particle.x, particle.y);
  ctx.rotate(particle.rotation);
  ctx.globalAlpha = alpha * 0.8;

  ctx.fillStyle = particle.color;
  ctx.fillRect(
    -particle.size / 2,
    -particle.size / 2,
    particle.size,
    particle.size
  );

  ctx.restore();
};

export const renderAllBlocks = (
  ctx: CanvasRenderingContext2D,
  blocks: LegoBlock[],
  config: PerformanceConfig,
  theme: Theme
): void => {
  const sortedBlocks = [...blocks].sort((a, b) => {
    if (a.settled !== b.settled) return a.settled ? -1 : 1;
    return a.targetY - b.targetY;
  });

  for (const block of sortedBlocks) {
    renderBlock(ctx, block, config, theme);

    for (const particle of block.particles) {
      renderParticle(ctx, particle);
    }
  }
};
