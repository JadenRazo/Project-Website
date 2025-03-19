import React, { useEffect } from 'react';

const Home: React.FC = () => {
  useEffect(() => {
    window.location.href = 'https://jadenrazo.dev';
  }, []);

  return (
    <div style={{ 
      height: '100vh', 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center' 
    }}>
      Redirecting...
    </div>
  );
};

export default Home; 