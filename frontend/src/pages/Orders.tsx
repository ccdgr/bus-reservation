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
  Snackbar
} from '@mui/material';
import client from '../api/client';
import { ShoppingBag, AccessTime, Cancel } from '@mui/icons-material';

interface Order {
  id: number;
  bus_id: number;
  status: number; // 0: Pending, 1: Paid, 2: Cancelled
  created_at: string;
}

const Orders: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [msg, setMsg] = useState('');

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = () => {
    client.get('/orders')
      .then(res => setOrders(res.data))
      .finally(() => setLoading(false));
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
      <Typography variant="h5" sx={{ mb: 2, fontWeight: 'bold' }}>
        我的订单
      </Typography>
      
      {orders.length === 0 && (
        <Alert severity="info">暂无订单信息</Alert>
      )}

      <Stack spacing={2}>
        {orders.map(order => (
          <Card key={order.id}>
            <CardContent>
              <Stack direction="row" sx={{ justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
                  <ShoppingBag color="action" fontSize="small" />
                  <Typography variant="body2" color="text.secondary">
                    订单号: {order.id}
                  </Typography>
                </Stack>
                {getStatusChip(order.status)}
              </Stack>
              
              <Divider sx={{ my: 1 }} />
              
              <Typography variant="body1" sx={{ mt: 1 }}>
                班次 ID: {order.bus_id}
              </Typography>
              
              <Stack direction="row" spacing={1} sx={{ alignItems: 'center', mt: 1 }}>
                <AccessTime fontSize="inherit" color="disabled" />
                <Typography variant="caption" color="text.secondary">
                  预定时间: {new Date(order.created_at).toLocaleString()}
                </Typography>
              </Stack>

              {order.status === 0 && (
                <Box sx={{ mt: 2, textAlign: 'right' }}>
                  <Button 
                    variant="outlined" 
                    color="error" 
                    size="small" 
                    startIcon={<Cancel />}
                    onClick={() => handleCancel(order.id)}
                  >
                    取消订单
                  </Button>
                </Box>
              )}
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
