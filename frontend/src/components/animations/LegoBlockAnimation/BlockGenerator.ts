import type { Theme } from '../../../styles/theme.types';
import type {
  LegoBlock,
  BlockShape,
  BlockContent,
  GradientColor,
  StudPosition,
  BlockGeneratorConfig,
} from './types';
import {
  BLOCK_SHAPES,
  SHAPE_WEIGHTS,
  BLOCK_CONTENT,
  GRID,
  ANIMATION,
  getBlockGradients,
} from './constants';

export class BlockGenerator {
  private seed: number;
  private gradients: GradientColor[];
  private availableTechnologies: string[];

  constructor(theme: Theme) {
    this.seed = Date.now() + Math.random() * 10000;
    this.gradients = getBlockGradients(theme);
    this.availableTechnologies = [...BLOCK_CONTENT.technologies];
    this.shuffleArray(this.availableTechnologies);
  }

  private shuffleArray<T>(array: T[]): void {
    for (let i = array.length - 1; i > 0; i--) {
      const j = Math.floor(this.seededRandom() * (i + 1));
      [array[i], array[j]] = [array[j], array[i]];
    }
  }

  private seededRandom(): number {
    this.seed = (this.seed + 0x6D2B79F5) | 0;
    let t = this.seed;
    t = Math.imul(t ^ (t >>> 15), t | 1);
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  }

  private weightedRandom<T extends string>(weights: Record<T, number>): T {
    const entries = Object.entries(weights) as [T, number][];
    const total = entries.reduce((sum, [, weight]) => sum + weight, 0);
    let random = this.seededRandom() * total;

    for (const [key, weight] of entries) {
      random -= weight;
      if (random <= 0) return key;
    }

    return entries[entries.length - 1][0];
  }

  private selectShape(): BlockShape {
    return this.weightedRandom(SHAPE_WEIGHTS);
  }

  private selectShapeForWidth(remainingWidth: number): BlockShape {
    const shapes: BlockShape[] = ['1x4', '1x3', '1x2', '1x1'];

    for (const shape of shapes) {
      if (BLOCK_SHAPES[shape].cols <= remainingWidth) {
        if (this.seededRandom() > 0.3 || BLOCK_SHAPES[shape].cols === remainingWidth) {
          return shape;
        }
      }
    }

    return '1x1';
  }

  private selectContent(): BlockContent {
    if (this.availableTechnologies.length === 0) {
      return { type: 'empty' };
    }

    const value = this.availableTechnologies.pop()!;
    return { type: 'icon', value };
  }

  private selectColor(): GradientColor {
    const index = Math.floor(this.seededRandom() * this.gradients.length);
    return this.gradients[index];
  }

  private generateStuds(width: number): StudPosition[] {
    const studs: StudPosition[] = [];
    const cols = Math.max(1, Math.floor(width / GRID.studSpacing));
    const offsetX = (width - (cols - 1) * GRID.studSpacing) / 2;
    const offsetY = GRID.studRadius + 2;

    for (let col = 0; col < cols; col++) {
      studs.push({
        x: offsetX + col * GRID.studSpacing,
        y: offsetY,
      });
    }

    return studs;
  }

  private createBlock(
    id: string,
    col: number,
    row: number,
    shape: BlockShape,
    spawnIndex: number,
    gridCols: number,
    gridRows: number
  ): LegoBlock {
    const shapeData = BLOCK_SHAPES[shape];
    const width = shapeData.cols * GRID.cellSize;
    const height = GRID.cellSize;

    const targetX = GRID.padding + col * GRID.cellSize;
    const targetY = GRID.padding + row * GRID.cellSize;

    const startY = -height - 20 - this.seededRandom() * 30;

    return {
      id,
      shape,
      x: targetX,
      y: startY,
      targetX,
      targetY,
      width,
      height,
      rotation: 0,
      color: this.selectColor(),
      studs: this.generateStuds(width),
      content: this.selectContent(),
      opacity: 1,
      velocity: 0,
      settled: false,
      dissolving: false,
      particles: [],
      spawnDelay: ANIMATION.initialDelay + spawnIndex * ANIMATION.spawnInterval,
      spawned: false,
      glowIntensity: 0,
      scaleX: 1,
      scaleY: 1,
      bouncePhase: 0,
      settleStartTime: 0,
    };
  }

  generatePattern(config: BlockGeneratorConfig): LegoBlock[] {
    const { viewportWidth, viewportHeight, blockCount } = config;

    const gridCols = Math.floor((viewportWidth - GRID.padding * 2) / GRID.cellSize);
    const gridRows = Math.floor((viewportHeight - GRID.padding * 2) / GRID.cellSize);

    const placements: { col: number; row: number; shape: BlockShape }[] = [];

    for (let row = gridRows - 1; row >= 0 && placements.length < blockCount; row--) {
      let col = 0;

      while (col < gridCols && placements.length < blockCount) {
        const remainingWidth = gridCols - col;
        const shape = this.selectShapeForWidth(remainingWidth);
        const shapeCols = BLOCK_SHAPES[shape].cols;

        placements.push({ col, row, shape });
        col += shapeCols;
      }
    }

    placements.sort((a, b) => {
      if (a.row !== b.row) return b.row - a.row;
      return a.col - b.col;
    });

    const blocks: LegoBlock[] = [];
    for (let i = 0; i < Math.min(placements.length, blockCount); i++) {
      const { col, row, shape } = placements[i];
      const block = this.createBlock(
        `block-${i}-${Date.now()}`,
        col,
        row,
        shape,
        i,
        gridCols,
        gridRows
      );
      blocks.push(block);
    }

    return blocks;
  }

  updateTheme(theme: Theme): void {
    this.gradients = getBlockGradients(theme);
  }
}

export const createBlockGenerator = (theme: Theme): BlockGenerator => {
  return new BlockGenerator(theme);
};
