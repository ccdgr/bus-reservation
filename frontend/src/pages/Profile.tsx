import React, { useEffect, useState } from 'react';
import { 
  Box, 
  Typography, 
  Paper, 
  Avatar, 
  List, 
  ListItem, 
  ListItemText, 
  ListItemIcon,
  Button,
  Divider
} from '@mui/material';
import { Person, Badge, Logout } from '@mui/icons-material';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

interface UserProfile {
  username: string;
  real_name: string;
  user_type: number;
}

const Profile: React.FC = () => {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const { logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    client.get('/users/profile')
      .then(res => setProfile(res.data));
  }, []);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  if (!profile) return <Typography sx={{ p: 2 }}>加载中...</Typography>;

  return (
    <Box>
      <Paper sx={{ p: 4, textAlign: 'center', mb: 3 }}>
        <Avatar 
          sx={{ width: 80, height: 80, mx: 'auto', mb: 2, bgcolor: 'primary.main' }}
        >
          {profile.real_name[0]}
        </Avatar>
        <Typography variant="h5" sx={{ fontWeight: 'bold' }}>{profile.real_name}</Typography>
        <Typography color="text.secondary">
          {profile.user_type === 0 ? '学生' : '教职工'}
        </Typography>
      </Paper>

      <Paper>
        <List>
          <ListItem>
            <ListItemIcon><Person /></ListItemIcon>
            <ListItemText primary="用户名" secondary={profile.username} />
          </ListItem>
          <Divider variant="inset" component="li" />
          <ListItem>
            <ListItemIcon><Badge /></ListItemIcon>
            <ListItemText primary="真实姓名" secondary={profile.real_name} />
          </ListItem>
        </List>
      </Paper>

      <Box sx={{ mt: 4 }}>
        <Button 
          variant="outlined" 
          color="error" 
          fullWidth 
          size="large"
          startIcon={<Logout />}
          onClick={handleLogout}
        >
          退出登录
        </Button>
      </Box>
      
      <Box sx={{ mt: 2, textAlign: 'center' }}>
        <Typography variant="caption" color="text.secondary">
          校车预定平台 v1.0.0
        </Typography>
      </Box>
    </Box>
  );
};

export default Profile;
