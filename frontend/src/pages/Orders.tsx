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
  Tab
} from '@mui/material';
import client from '../api/client';
import { ShoppingBag, AccessTime, Cancel, ArrowBack } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

interface Order {
  id: number;
  bus_id: number;
  status: number; // 0: Pending, 1: Paid, 2: Cancelled
  created_at: string;
}

const Orders: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [filteredOrders, setFilteredOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [msg, setMsg] = useState('');
  const [tabValue, setTabValue] = useState(0); // 0: All, 1: Pending, 2: Success, 3: Cancelled
  const navigate = useNavigate();

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
    } else if (tabValue === 1) {
      setFilteredOrders(orders.filter(o => o.status === 0));
    } else if (tabValue === 2) {
      setFilteredOrders(orders.filter(o => o.status === 1));
    } else if (tabValue === 3) {
      setFilteredOrders(orders.filter(o => o.status === 2));
    }
  };

  const handleCancel = async (id: number) => {
    try {
      await client.post(`/orders/${id}/cancel`);
      setMsg('订单已取消');
      setOpen(true);
      fetchOrders();
    } catch (err: any) {
      setMsg(err.response?.data?.error || '取消失败');
      setOpen(true);
    }
  };

  const getStatusChip = (status: number) => {
    switch (status) {
      case 0: return <Chip label="处理中" size="small" color="warning" />;
      case 1: return <Chip label="预定成功" size="small" color="success" />;
      case 2: return <Chip label="已取消" size="small" color="default" />;
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
        variant="fullWidth" 
        sx={{ mb: 3, borderBottom: 1, borderColor: 'divider' }}
      >
        <Tab label="全部" />
        <Tab label="处理中" />
        <Tab label="成功" />
        <Tab label="已取消" />
      </Tabs>
      
      {filteredOrders.length === 0 && (
        <Alert severity="info" sx={{ borderRadius: 3 }}>暂无订单信息</Alert>
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
              
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Box>
                  <Typography variant="body1" sx={{ fontWeight: 'bold' }}>
                    班次 ID: {order.bus_id}
                  </Typography>
                  <Stack direction="row" spacing={0.5} sx={{ alignItems: 'center', mt: 0.5 }}>
                    <AccessTime sx={{ fontSize: 14, color: 'text.secondary' }} />
                    <Typography variant="caption" color="text.secondary">
                      {new Date(order.created_at).toLocaleString()}
                    </Typography>
                  </Stack>
                </Box>
                
                {order.status === 0 && (
                  <Button 
                    variant="outlined" 
                    color="error" 
                    size="small" 
                    startIcon={<Cancel />}
                    onClick={() => handleCancel(order.id)}
                    sx={{ borderRadius: 2 }}
                  >
                    取消
                  </Button>
                )}
              </Box>
            </CardContent>
          </Card>
        ))}
      </Stack>

      <Snackbar open={open} autoHideDuration={3000} onClose={() => setOpen(false)}>
        <Alert severity="info" sx={{ width: '100%' }}>
          {msg}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default Orders;
