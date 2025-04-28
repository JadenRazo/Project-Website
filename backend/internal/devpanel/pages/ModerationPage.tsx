import React from 'react';
import { Container, Typography } from '@mui/material';
import { WordFilterManager } from '../components/WordFilterManager';

export const ModerationPage: React.FC = () => {
  return (
    <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
      <Typography variant="h4" component="h1" gutterBottom>
        Moderation Settings
      </Typography>
      <WordFilterManager />
    </Container>
  );
}; 