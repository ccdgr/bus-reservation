import React from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { 
  Box, 
  BottomNavigation, 
  BottomNavigationAction, 
  AppBar, 
  Toolbar, 
  Typography,
  Container
} from '@mui/material';
import { 
  Home as HomeIcon, 
  ListAlt as OrdersIcon, 
  Person as ProfileIcon 
} from '@mui/icons-material';

const Layout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();

  // Determine active tab based on path
  const getValue = () => {
    if (location.pathname === '/') return 0;
    if (location.pathname === '/orders') return 1;
    if (location.pathname === '/profile') return 2;
    return 0;
  };

  return (
    <Box sx={{ pb: 7, bgcolor: 'background.default', minHeight: '100vh' }}>
      <AppBar position="sticky" elevation={0}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1, textAlign: 'center' }}>
            校车预定
          </Typography>
        </Toolbar>
      </AppBar>

      <Container maxWidth="sm" sx={{ mt: 2 }}>
        <Outlet />
      </Container>

      <BottomNavigation
        value={getValue()}
        onChange={(_, newValue) => {
          if (newValue === 0) navigate('/');
          if (newValue === 1) navigate('/orders');
          if (newValue === 2) navigate('/profile');
        }}
        showLabels
      >
        <BottomNavigationAction label="首页" icon={<HomeIcon />} />
        <BottomNavigationAction label="订单" icon={<OrdersIcon />} />
        <BottomNavigationAction label="我的" icon={<ProfileIcon />} />
      </BottomNavigation>
    </Box>
  );
};

export default Layout;
