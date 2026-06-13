import React, { useEffect, useState } from 'react';
import { 
  Box, 
  Typography, 
  Button, 
  Paper, 
  Stack, 
  Divider,
  Alert,
  Snackbar,
  IconButton,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  CircularProgress
} from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  ArrowBack, 
  DirectionsBus, 
  AccessTime, 
  Person,
  CreditCard,
  ChatBubble,
  AccountBalanceWallet,
  ReceiptLong
} from '@mui/icons-material';
import client from '../api/client';

interface Bus {
  id: number;
  number: string;
  origin: string;
  dest: string;
  start_time: string;
}

interface UserProfile {
  real_name: string;
}

const ConfirmOrder: React.FC = () => {
  const { busId } = useParams();
  const navigate = useNavigate();
  
  const [bus, setBus] = useState<Bus | null>(null);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState({ type: 'info', msg: '' });
  const [open, setOpen] = useState(false);

  // Payment Drawer states
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [processing, setProcessing] = useState(false);
  const [selectedMethod, setSelectedMethod] = useState('alipay'); // Only Alipay allowed

  useEffect(() => {
    // Fetch bus and user profile simultaneously
    Promise.all([
      client.get(`/buses/${busId}`),
      client.get('/users/profile')
    ])
    .then(([busRes, profileRes]) => {
      setBus(busRes.data);
      setProfile(profileRes.data);
    })
    .catch(() => {
      setStatus({ type: 'error', msg: '获取数据失败，请检查网络或登录状态' });
      setOpen(true);
    })
    .finally(() => setLoading(false));
  }, [busId]);

  const handlePay = async () => {
    if (selectedMethod !== 'alipay') return;

    setProcessing(true);
    try {
      // 1. Create order
      const orderRes = await client.post('/orders', { bus_id: Number(busId) });
      const orderID = orderRes.data.id;

      // 2. Call payment API
      if (orderID) {
        const payRes = await client.post(`/orders/${orderID}/pay`);
        if (payRes.data.payment_url) {
          // Redirect to Alipay Sandbox
          window.location.href = payRes.data.payment_url;
          return;
        } else {
          // Fallback if no URL returned (e.g. backend missing alipay config)
          setStatus({ type: 'success', msg: '支付成功' });
          setOpen(true);
          setTimeout(() => navigate('/orders'), 1500);
        }
      }
    } catch (err: any) {
      console.error('Payment failed:', err);
      setStatus({ type: 'error', msg: err.response?.data?.error || '支付失败' });
      setOpen(true);
    } finally {
      setProcessing(false);
      setDrawerOpen(false);
    }
  };

  if (loading || !bus || !profile) return <Typography sx={{ p: 4 }}>加载中...</Typography>;

  return (
    <Box sx={{ pb: 12 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <IconButton onClick={() => navigate(-1)} sx={{ mr: 1 }}>
          <ArrowBack />
        </IconButton>
        <Typography variant="h6" sx={{ fontWeight: 'bold' }}>确认预约</Typography>
      </Box>

      <Paper sx={{ p: 3, mb: 3, borderRadius: 3 }}>
        <Stack spacing={3}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="subtitle1" color="text.secondary">乘车路线</Typography>
            <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
              <DirectionsBus color="primary" fontSize="small" />
              <Typography variant="body1" sx={{ fontWeight: 'bold' }}>
                {bus.origin} ➔ {bus.dest}
              </Typography>
            </Stack>
          </Box>

          <Divider />

          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="subtitle1" color="text.secondary">发车时间</Typography>
            <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
              <AccessTime color="action" fontSize="small" />
              <Typography variant="body1">
                {new Date(bus.start_time).toLocaleString()}
              </Typography>
            </Stack>
          </Box>

          <Divider />

          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="subtitle1" color="text.secondary">预约人</Typography>
            <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
              <Person color="action" fontSize="small" />
              <Typography variant="body1">
                {profile.real_name}
              </Typography>
            </Stack>
          </Box>

          <Divider />

          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Typography variant="subtitle1" color="text.secondary">乘车费用</Typography>
            <Stack direction="row" spacing={0.5} sx={{ alignItems: 'center', color: 'error.main' }}>
              <Typography variant="body2">¥</Typography>
              <Typography variant="h6" sx={{ fontWeight: 'bold' }}>0.01</Typography>
            </Stack>
          </Box>
        </Stack>
      </Paper>

      <Alert severity="info" sx={{ borderRadius: 3, mb: 3 }} icon={<ReceiptLong />}>
        支付成功后请在“我的订单”中进行核验上车。
      </Alert>

      {/* Bottom Fixed Action Bar */}
      <Paper 
        elevation={4} 
        sx={{ 
          position: 'fixed', 
          bottom: 0, 
          left: 0, 
          right: 0, 
          p: 2,
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          borderTopLeftRadius: 16,
          borderTopRightRadius: 16,
          zIndex: 100
        }}
      >
        <Box>
          <Typography variant="caption" color="text.secondary">总计</Typography>
          <Typography variant="h5" color="error.main" sx={{ fontWeight: 'bold' }}>¥ 0.01</Typography>
        </Box>
        <Button 
          variant="contained" 
          size="large" 
          onClick={() => setDrawerOpen(true)}
          sx={{ px: 4, py: 1.5, borderRadius: 8, fontWeight: 'bold' }}
        >
          确定预约
        </Button>
      </Paper>

      {/* Payment Drawer (Bottom Sheet) */}
      <Drawer
        anchor="bottom"
        open={drawerOpen}
        onClose={() => !processing && setDrawerOpen(false)}
        sx={{ '& .MuiDrawer-paper': { borderTopLeftRadius: 16, borderTopRightRadius: 16, pb: 4 } }}
      >
        <Box sx={{ p: 3 }}>
          <Typography variant="h6" sx={{ fontWeight: 'bold', mb: 2, textAlign: 'center' }}>
            选择支付方式
          </Typography>
          <List sx={{ pt: 0 }}>
            <ListItem disablePadding sx={{ mb: 1 }}>
              <ListItemButton 
                selected={selectedMethod === 'alipay'}
                onClick={() => setSelectedMethod('alipay')}
                sx={{ borderRadius: 2, border: selectedMethod === 'alipay' ? '1px solid #1677FF' : '1px solid transparent' }}
                disabled={processing}
              >
                <ListItemIcon><AccountBalanceWallet sx={{ color: '#1677FF', fontSize: 32 }} /></ListItemIcon>
                <ListItemText 
                  primary={<Typography sx={{ fontWeight: 'bold', color: '#1677FF' }}>支付宝支付 (沙箱环境)</Typography>} 
                  secondary="推荐使用" 
                />
              </ListItemButton>
            </ListItem>
            
            <ListItem disablePadding sx={{ mb: 1 }}>
              <ListItemButton disabled>
                <ListItemIcon><ChatBubble sx={{ color: '#ccc', fontSize: 32 }} /></ListItemIcon>
                <ListItemText primary="微信支付" secondary="暂不支持" />
              </ListItemButton>
            </ListItem>

            <ListItem disablePadding>
              <ListItemButton disabled>
                <ListItemIcon><CreditCard sx={{ color: '#ccc', fontSize: 32 }} /></ListItemIcon>
                <ListItemText primary="校园一卡通" secondary="维护中" />
              </ListItemButton>
            </ListItem>
          </List>

          <Button 
            fullWidth 
            variant="contained" 
            size="large" 
            onClick={handlePay} 
            disabled={processing || selectedMethod !== 'alipay'}
            sx={{ mt: 3, py: 1.5, borderRadius: 2, fontWeight: 'bold', bgcolor: '#1677FF' }}
          >
            {processing ? <CircularProgress size={24} color="inherit" /> : '立即支付 ¥ 0.01'}
          </Button>
        </Box>
      </Drawer>

      <Snackbar open={open} autoHideDuration={3000} onClose={() => setOpen(false)}>
        <Alert severity={status.type as any} sx={{ width: '100%' }}>
          {status.msg}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default ConfirmOrder;
