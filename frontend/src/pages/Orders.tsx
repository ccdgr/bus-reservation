import React, { useEffect, useState } from 'react';
import { 
  Box, 
  Card, 
  CardContent, 
  Typography, 
  Stack, 
  Chip,
  Button,
  Divider,
  Alert,
  Snackbar,
  IconButton,
  Tabs,
  Tab,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions
} from '@mui/material';
import client from '../api/client';
import { 
  ShoppingBag, 
  AccessTime, 
  Cancel, 
  ArrowBack, 
  Payments, 
  FactCheck,
  QrCode2,
  CheckCircleOutlined
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

interface Order {
  id: number;
  bus_id: number;
  status: number;
  created_at: string;
  bus?: {
    number: string;
    origin: string;
    dest: string;
    start_time: string;
  };
}

const Orders: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [filteredOrders, setFilteredOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [msg, setMsg] = useState('');
  const [tabValue, setTabValue] = useState(0); 
  const navigate = useNavigate();

  // Verification dialog states
  const [verifyDialogOpen, setVerifyDialogOpen] = useState(false);
  const [selectedOrderID, setSelectedOrderID] = useState<number | null>(null);
  const [verifying, setVerifying] = useState(false);

  useEffect(() => {
    fetchOrders();
  }, []);

  useEffect(() => {
    filterOrders();
  }, [tabValue, orders]);

  const fetchOrders = () => {
    setLoading(true);
    client.get('/orders')
      .then(res => setOrders(res.data))
      .finally(() => setLoading(false));
  };

  const filterOrders = () => {
    if (tabValue === 0) {
      setFilteredOrders(orders);
    } else {
      const statusMap = [ -1, 0, 1, 2, 3, 4 ];
      const targetStatus = statusMap[tabValue];
      setFilteredOrders(orders.filter(o => o.status === targetStatus));
    }
  };

  const handleAction = async (id: number, action: 'cancel' | 'pay' | 'verify') => {
    try {
      await client.post(`/orders/${id}/${action}`);
      setMsg(`操作成功`);
      setOpen(true);
      fetchOrders();
    } catch (err: any) {
      setMsg(err.response?.data?.error || '操作失败');
      setOpen(true);
    }
  };

  const openVerifyDialog = (id: number) => {
    setSelectedOrderID(id);
    setVerifyDialogOpen(true);
  };

  const handleSimulateVerify = async () => {
    if (!selectedOrderID) return;
    setVerifying(true);
    try {
      await client.post(`/orders/${selectedOrderID}/verify`);
      setMsg('核验成功，欢迎乘车！');
      setOpen(true);
      setVerifyDialogOpen(false);
      fetchOrders();
    } catch (err: any) {
      setMsg(err.response?.data?.error || '核验失败');
      setOpen(true);
    } finally {
      setVerifying(false);
    }
  };

  const getStatusChip = (status: number) => {
    switch (status) {
      case 0: return <Chip label="待支付" size="small" color="warning" variant="outlined" />;
      case 1: return <Chip label="待核验" size="small" color="info" variant="outlined" />;
      case 2: return <Chip label="已取消" size="small" color="default" variant="outlined" />;
      case 3: return <Chip label="已过期" size="small" color="error" variant="outlined" />;
      case 4: return <Chip label="已核验" size="small" color="success" variant="outlined" />;
      default: return <Chip label="未知" size="small" />;
    }
  };

  if (loading) return <Typography sx={{ p: 2 }}>加载中...</Typography>;

  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <IconButton onClick={() => navigate('/')} sx={{ mr: 1 }}>
          <ArrowBack />
        </IconButton>
        <Typography variant="h6" sx={{ fontWeight: 'bold' }}>我的订单</Typography>
      </Box>

      <Tabs 
        value={tabValue} 
        onChange={(_, v) => setTabValue(v)} 
        variant="scrollable" 
        scrollButtons="auto"
        sx={{ mb: 3, borderBottom: 1, borderColor: 'divider' }}
      >
        <Tab label="全部" />
        <Tab label="待支付" />
        <Tab label="待核验" />
        <Tab label="已取消" />
        <Tab label="已过期" />
        <Tab label="已核验" />
      </Tabs>
      
      {filteredOrders.length === 0 && (
        <Alert severity="info" sx={{ borderRadius: 3 }}>暂无相关订单</Alert>
      )}

      <Stack spacing={2}>
        {filteredOrders.map(order => (
          <Card key={order.id} sx={{ borderRadius: 3 }}>
            <CardContent>
              <Stack direction="row" sx={{ justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
                  <ShoppingBag color="action" sx={{ fontSize: 20 }} />
                  <Typography variant="subtitle2" color="text.secondary">
                    订单号: {order.id}
                  </Typography>
                </Stack>
                {getStatusChip(order.status)}
              </Stack>
              
              <Divider sx={{ my: 1.5 }} />
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="body1" sx={{ fontWeight: 'bold' }}>
                  {order.bus ? `${order.bus.number} (${order.bus.origin} ➔ ${order.bus.dest})` : `班次 ID: ${order.bus_id}`}
                </Typography>
                <Stack direction="row" spacing={0.5} sx={{ alignItems: 'center', mt: 0.5 }}>
                  <AccessTime sx={{ fontSize: 14, color: 'text.secondary' }} />
                  <Typography variant="caption" color="text.secondary">
                    下单时间: {new Date(order.created_at).toLocaleString()}
                  </Typography>
                </Stack>
              </Box>

              <Stack direction="row" spacing={1} sx={{ justifyContent: 'flex-end' }}>
                {order.status === 0 && (
                  <>
                    <Button 
                      variant="outlined" 
                      color="error" 
                      size="small" 
                      startIcon={<Cancel />}
                      onClick={() => handleAction(order.id, 'cancel')}
                    >
                      取消
                    </Button>
                    <Button 
                      variant="contained" 
                      color="primary" 
                      size="small" 
                      startIcon={<Payments />}
                      onClick={() => handleAction(order.id, 'pay')}
                    >
                      立即支付
                    </Button>
                  </>
                )}
                {order.status === 1 && (
                  <>
                    <Button 
                      variant="outlined" 
                      color="error" 
                      size="small" 
                      onClick={() => handleAction(order.id, 'cancel')}
                    >
                      退票
                    </Button>
                    <Button 
                      variant="contained" 
                      color="success" 
                      size="small" 
                      startIcon={<FactCheck />}
                      onClick={() => openVerifyDialog(order.id)}
                    >
                      核验上车
                    </Button>
                  </>
                )}
                {order.status === 4 && (
                  <Chip icon={<CheckCircleOutlined />} label="核验成功" color="success" variant="outlined" size="small" />
                )}
              </Stack>
            </CardContent>
          </Card>
        ))}
      </Stack>

      {/* Verification Dialog */}
      <Dialog open={verifyDialogOpen} onClose={() => !verifying && setVerifyDialogOpen(false)} fullWidth maxWidth="xs">
        <DialogTitle sx={{ textAlign: 'center', fontWeight: 'bold' }}>乘车核验</DialogTitle>
        <DialogContent sx={{ textAlign: 'center', py: 3 }}>
          <Box sx={{ 
            p: 2, 
            bgcolor: '#f8f9fa', 
            borderRadius: 4, 
            display: 'inline-block',
            border: '2px solid #e9ecef',
            mb: 2
          }}>
            <QrCode2 sx={{ fontSize: 200, color: '#333' }} />
          </Box>
          <Typography variant="body2" color="text.secondary">
            请将二维码对准校车扫码器进行核验
          </Typography>
        </DialogContent>
        <DialogActions sx={{ p: 2 }}>
          <Button 
            fullWidth 
            variant="contained" 
            size="large" 
            onClick={handleSimulateVerify} 
            disabled={verifying}
            startIcon={<FactCheck />}
            sx={{ py: 1.5, fontWeight: 'bold' }}
          >
            {verifying ? '核验中...' : '模拟核验成功'}
          </Button>
        </DialogActions>
      </Dialog>

      <Snackbar open={open} autoHideDuration={3000} onClose={() => setOpen(false)}>
        <Alert severity="info" sx={{ width: '100%' }}>
          {msg}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Orders;
