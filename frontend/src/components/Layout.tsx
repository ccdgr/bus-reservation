import React from 'react';
import { Outlet } from 'react-router-dom';
import { 
  Box, 
  AppBar, 
  Toolbar, 
  Typography,
  Container
} from '@mui/material';

const Layout: React.FC = () => {
  return (
    <Box sx={{ bgcolor: 'background.default', minHeight: '100vh' }}>
      <AppBar position="sticky" elevation={0} sx={{ borderBottom: '1px solid', borderColor: 'divider' }}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1, textAlign: 'center', fontWeight: 'bold', color: 'primary.main' }}>
            校车预定
          </Typography>
        </Toolbar>
      </AppBar>

      <Container maxWidth="sm" sx={{ mt: 2, pb: 4 }}>
        <Outlet />
      </Container>
    </Box>
  );
};

export default Layout;
