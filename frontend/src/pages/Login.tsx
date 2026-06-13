import React, { useState } from 'react';
import { 
  Box, 
  TextField, 
  Button, 
  Typography, 
  Paper, 
  Link,
  Alert,
  Snackbar,
  CircularProgress
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';

const Login: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (loading) return;
    
    setLoading(true);
    try {
      const res = await client.post('/users/login', { username, password });
      login(res.data.token);
      navigate('/');
    } catch (err: any) {
      console.error('Login error:', err);
      const msg = err.response?.data?.error || err.message || '登录失败，请检查网络或账号密码';
      setError(msg);
      setOpen(true);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ 
      display: 'flex', 
      flexDirection: 'column', 
      alignItems: 'center', 
      justifyContent: 'center',
      minHeight: '100vh',
      px: 2,
      bgcolor: 'background.default'
    }}>
      <Paper sx={{ p: 4, width: '100%', maxWidth: 400, borderRadius: 4 }}>
        <Typography variant="h5" align="center" gutterBottom sx={{ fontWeight: 'bold', color: 'primary.main', mb: 3 }}>
          欢迎回来
        </Typography>
        <form onSubmit={handleSubmit}>
          <TextField
            label="用户名"
            fullWidth
            margin="normal"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            disabled={loading}
          />
          <TextField
            label="密码"
            type="password"
            fullWidth
            margin="normal"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            disabled={loading}
          />
          <Button 
            type="submit" 
            variant="contained" 
            fullWidth 
            size="large"
            disabled={loading}
            sx={{ mt: 3, mb: 2, py: 1.5, fontWeight: 'bold' }}
          >
            {loading ? <CircularProgress size={24} color="inherit" /> : '登录'}
          </Button>
          <Box sx={{ textAlign: 'center', mt: 1 }}>
            <Link 
              component="button"
              type="button"
              variant="body2"
              onClick={() => navigate('/register')} 
              sx={{ cursor: 'pointer', textDecoration: 'none' }}
              disabled={loading}
            >
              没有账号？立即注册
            </Link>
          </Box>
        </form>
      </Paper>
      
      <Snackbar 
        open={open} 
        autoHideDuration={4000} 
        onClose={() => setOpen(false)}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert onClose={() => setOpen(false)} severity="error" variant="filled" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Login;
