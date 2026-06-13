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
  MenuItem
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';

const Register: React.FC = () => {
  const [formData, setFormData] = useState({
    username: '',
    password: '',
    real_name: '',
    user_type: 0
  });
  const [status, setStatus] = useState({ type: '', msg: '' });
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await client.post('/users/register', formData);
      setStatus({ type: 'success', msg: '注册成功，正在跳转登录...' });
      setOpen(true);
      setTimeout(() => navigate('/login'), 2000);
    } catch (err: any) {
      setStatus({ type: 'error', msg: err.response?.data?.error || '注册失败' });
      setOpen(true);
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
      <Paper sx={{ p: 4, width: '100%', maxWidth: 400 }}>
        <Typography variant="h5" align="center" gutterBottom sx={{ fontWeight: 'bold' }}>
          创建账号
        </Typography>
        <form onSubmit={handleSubmit}>
          <TextField
            label="用户名"
            fullWidth
            margin="normal"
            required
            onChange={(e) => setFormData({...formData, username: e.target.value})}
          />
          <TextField
            label="真实姓名"
            fullWidth
            margin="normal"
            required
            onChange={(e) => setFormData({...formData, real_name: e.target.value})}
          />
          <TextField
            label="密码"
            type="password"
            fullWidth
            margin="normal"
            required
            onChange={(e) => setFormData({...formData, password: e.target.value})}
          />
          <TextField
            select
            label="用户类型"
            fullWidth
            margin="normal"
            value={formData.user_type}
            onChange={(e) => setFormData({...formData, user_type: Number(e.target.value)})}
          >
            <MenuItem value={0}>学生</MenuItem>
            <MenuItem value={1}>教职工</MenuItem>
          </TextField>
          <Button 
            type="submit" 
            variant="contained" 
            fullWidth 
            size="large"
            sx={{ mt: 3, mb: 2 }}
          >
            注册
          </Button>
          <Box sx={{ textAlign: 'center' }}>
            <Link onClick={() => navigate('/login')} sx={{ cursor: 'pointer' }}>
              已有账号？返回登录
            </Link>
          </Box>
        </form>
      </Paper>
      <Snackbar open={open} autoHideDuration={3000} onClose={() => setOpen(false)}>
        <Alert severity={status.type as any} sx={{ width: '100%' }}>
          {status.msg}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Register;
