#!/usr/bin/env node

const fs = require('fs');
const path = require('path');

console.log('🧹 Clearing React app caches...\n');

// Clear node_modules/.cache
const cacheDir = path.join(__dirname, 'node_modules', '.cache');
if (fs.existsSync(cacheDir)) {
  console.log('Removing node_modules/.cache...');
  fs.rmSync(cacheDir, { recursive: true, force: true });
  console.log('✅ node_modules/.cache cleared');
} else {
  console.log('ℹ️  No node_modules/.cache found');
}

// Clear build directory
const buildDir = path.join(__dirname, 'build');
if (fs.existsSync(buildDir)) {
  console.log('\nRemoving build directory...');
  fs.rmSync(buildDir, { recursive: true, force: true });
  console.log('✅ Build directory cleared');
} else {
  console.log('\nℹ️  No build directory found');
}

// Clear .eslintcache if it exists
const eslintCache = path.join(__dirname, '.eslintcache');
if (fs.existsSync(eslintCache)) {
  console.log('\nRemoving .eslintcache...');
  fs.unlinkSync(eslintCache);
  console.log('✅ .eslintcache cleared');
}

// Clear tsconfig.tsbuildinfo if it exists
const tsBuildInfo = path.join(__dirname, 'tsconfig.tsbuildinfo');
if (fs.existsSync(tsBuildInfo)) {
  console.log('\nRemoving tsconfig.tsbuildinfo...');
  fs.unlinkSync(tsBuildInfo);
  console.log('✅ tsconfig.tsbuildinfo cleared');
}

console.log('\n✨ All caches cleared! You can now run npm start with fresh caches.');