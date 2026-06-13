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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
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
  LocationOn, 
  AccessTime, 
  EventSeat,
  CreditCard,
  ChatBubble,
  AccountBalanceWallet,
  CheckCircleOutlined
} from '@mui/icons-material';
import client from '../api/client';
import { useAuth } from '../context/AuthContext';

interface Bus {
  id: number;
  number: string;
  origin: string;
  dest: string;
  start_time: string;
  total_seat: number;
  left_seat: number;
}

const BusDetail: React.FC = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  
  const [bus, setBus] = useState<Bus | null>(null);
  const [loading, setLoading] = useState(true);
  const [status, setStatus] = useState({ type: 'info', msg: '' });
  const [open, setOpen] = useState(false);

  // Payment states
  const [payDialogOpen, setPayDialogOpen] = useState(false);
  const [processing, setProcessing] = useState(false);
  const [paySuccessOpen, setPaySuccessOpen] = useState(false);
  const [selectedMethod, setSelectedMethod] = useState('');

  useEffect(() => {
    client.get(`/buses/${id}`)
      .then(res => setBus(res.data))
      .finally(() => setLoading(false));
  }, [id]);

  const handleOpenPayment = () => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    setPayDialogOpen(true);
  };

  const handleFinalPay = async () => {
    if (!selectedMethod) {
      setStatus({ type: 'error', msg: '请选择支付方式' });
      setOpen(true);
      return;
    }

    setProcessing(true);
    try {
      // 1. 创建订单 (Status 0)
      const orderRes = await client.post('/orders', { bus_id: Number(id) });
      const orderID = orderRes.data.id;

      // 2. 模拟支付延迟
      await new Promise(resolve => setTimeout(resolve, 1500));

      // 3. 调用支付接口 (Status 1)
      if (orderID) {
        let payRes;
        if (selectedMethod === 'alipay') {
          payRes = await client.post(`/orders/${orderID}/pay`);
        } else {
          // If not Alipay, we simulate the internal mock
          payRes = await client.post(`/orders/${orderID}/pay`);
        }

        if (payRes.data.payment_url) {
          // Real Alipay flow: redirect to Alipay Sandbox
          window.location.href = payRes.data.payment_url;
          return;
        }
      }

      setPayDialogOpen(false);
      setPaySuccessOpen(true);
      
      // 4. 跳转
      setTimeout(() => navigate('/orders'), 2000);
    } catch (err: any) {
      console.error('Reservation failed:', err);
      setStatus({ type: 'error', msg: err.response?.data?.error || '预定或支付失败' });
      setOpen(true);
    } finally {
      setProcessing(false);
    }
  };

  if (loading || !bus) return <Typography sx={{ p: 4 }}>加载中...</Typography>;

  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <IconButton onClick={() => navigate(-1)} sx={{ mr: 1 }}>
          <ArrowBack />
        </IconButton>
        <Typography variant="h6" sx={{ fontWeight: 'bold' }}>班次详情</Typography>
      </Box>

      <Paper sx={{ p: 3, mb: 3, borderRadius: 3 }}>
        <Stack spacing={3}>
          <Box>
            <Typography color="text.secondary" variant="caption" sx={{ display: 'block' }}>班次号</Typography>
            <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
              <DirectionsBus color="primary" />
              <Typography variant="h5" sx={{ fontWeight: 'bold' }}>{bus.number}</Typography>
            </Stack>
          </Box>

          <Divider />

          <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
            <Box>
              <Typography color="text.secondary" variant="caption" sx={{ display: 'block' }}>起点</Typography>
              <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
                <LocationOn sx={{ color: 'success.main' }} />
                <Typography variant="body1">{bus.origin}</Typography>
              </Stack>
            </Box>
            <Box sx={{ textAlign: 'right' }}>
              <Typography color="text.secondary" variant="caption" sx={{ display: 'block' }}>终点</Typography>
              <Stack direction="row" spacing={1} sx={{ alignItems: 'center', justifyContent: 'flex-end' }}>
                <Typography variant="body1">{bus.dest}</Typography>
                <LocationOn sx={{ color: 'error.main' }} />
              </Stack>
            </Box>
          </Box>

          <Box>
            <Typography color="text.secondary" variant="caption" sx={{ display: 'block' }}>出发时间</Typography>
            <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
              <AccessTime color="action" />
              <Typography variant="body1">{new Date(bus.start_time).toLocaleString()}</Typography>
            </Stack>
          </Box>

          <Box>
            <Typography color="text.secondary" variant="caption" sx={{ display: 'block' }}>座位情况</Typography>
            <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
              <EventSeat color="action" />
              <Typography variant="body1">剩余 {bus.left_seat} / 总计 {bus.total_seat}</Typography>
            </Stack>
          </Box>
        </Stack>
      </Paper>

      <Box sx={{ 
        position: 'fixed', 
        bottom: 32, 
        left: 0, 
        right: 0, 
        px: 2,
        zIndex: 100
      }}>
        <Button 
          variant="contained" 
          fullWidth 
          size="large" 
          disabled={bus.left_seat <= 0}
          onClick={handleOpenPayment}
          sx={{ py: 1.5, fontSize: '1.1rem', boxShadow: 3, borderRadius: 3 }}
        >
          {bus.left_seat > 0 ? '立即预定' : '已售罄'}
        </Button>
      </Box>

      {/* Payment Selection Dialog */}
      <Dialog open={payDialogOpen} onClose={() => !processing && setPayDialogOpen(false)} fullWidth maxWidth="xs">
        <DialogTitle sx={{ fontWeight: 'bold', textAlign: 'center' }}>选择支付方式</DialogTitle>
        <DialogContent sx={{ p: 0 }}>
          <List sx={{ pt: 0 }}>
            {[
              { id: 'card', name: '一卡通支付', icon: <CreditCard sx={{ color: '#ff9800' }} /> },
              { id: 'wechat', name: '微信支付', icon: <ChatBubble sx={{ color: '#4caf50' }} /> },
              { id: 'alipay', name: '支付宝支付', icon: <AccountBalanceWallet sx={{ color: '#2196f3' }} /> }
            ].map((method) => (
              <ListItem key={method.id} disablePadding>
                <ListItemButton 
                  selected={selectedMethod === method.id}
                  onClick={() => setSelectedMethod(method.id)}
                  disabled={processing}
                >
                  <ListItemIcon>{method.icon}</ListItemIcon>
                  <ListItemText primary={method.name} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button fullWidth variant="contained" size="large" onClick={handleFinalPay} disabled={processing || !selectedMethod}>
            {processing ? <CircularProgress size={24} color="inherit" /> : '立即支付'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Success Simulation Popup */}
      <Dialog open={paySuccessOpen}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', py: 4, px: 4, textAlign: 'center' }}>
          <CheckCircleOutlined sx={{ fontSize: 64, color: 'success.main', mb: 2 }} />
          <Typography variant="h5" sx={{ fontWeight: 'bold', mb: 1 }}>支付成功</Typography>
          <Typography color="text.secondary">正在为您出票，请稍后...</Typography>
        </Box>
      </Dialog>

      <Snackbar open={open} autoHideDuration={3000} onClose={() => setOpen(false)}>
        <Alert severity={status.type as any} sx={{ width: '100%' }}>
          {status.msg}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default BusDetail;
