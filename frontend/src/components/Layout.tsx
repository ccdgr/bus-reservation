import React from 'react';
import { Outlet } from 'react-router-dom';
import { 
  Box, 
  Container
} from '@mui/material';

const Layout: React.FC = () => {
  return (
    <Box sx={{ bgcolor: 'background.default', minHeight: '100vh' }}>
      <Container maxWidth="sm" sx={{ mt: 2, pb: 4 }}>
        <Outlet />
      </Container>
    </Box>
  );
};

export default Layout;
