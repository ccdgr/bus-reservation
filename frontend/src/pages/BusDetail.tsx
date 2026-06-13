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
  IconButton
} from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowBack, DirectionsBus, LocationOn, AccessTime, EventSeat } from '@mui/icons-material';
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
  const [status, setStatus] = useState({ type: '', msg: '' });
  const [open, setOpen] = useState(false);

  useEffect(() => {
    client.get(`/buses/${id}`)
      .then(res => setBus(res.data))
      .finally(() => setLoading(false));
  }, [id]);

  const handleReserve = async () => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }

    try {
      await client.post('/orders', { bus_id: Number(id) });
      setStatus({ type: 'success', msg: '预定成功！请前往订单查看结果' });
      setOpen(true);
      setTimeout(() => navigate('/orders'), 1500);
    } catch (err: any) {
      setStatus({ type: 'error', msg: err.response?.data?.error || '预定失败' });
      setOpen(true);
    }
  };

  if (loading || !bus) return <Typography sx={{ p: 4 }}>加载中...</Typography>;

  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <IconButton onClick={() => navigate(-1)} sx={{ mr: 1 }}>
          <ArrowBack />
        </IconButton>
        <Typography variant="h6">班次详情</Typography>
      </Box>

      <Paper sx={{ p: 3, mb: 3 }}>
        <Stack spacing={3}>
          <Box>
            <Typography color="text.secondary" variant="caption" sx={{ display: 'block' }}>班次号</Typography>
            <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
              <DirectionsBus color="primary" />
              <Typography variant="h5" sx={{ fontWeight: 'bold' }}>{bus.number}</Typography>
            </Stack>
          </Box>

          <Divider />

          <Stack direction="row" sx={{ justifyContent: 'space-between' }}>
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
          </Stack>

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
        bottom: 72, 
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
          onClick={handleReserve}
          sx={{ py: 1.5, fontSize: '1.1rem', boxShadow: 3 }}
        >
          {bus.left_seat > 0 ? '立即预定' : '已售罄'}
        </Button>
      </Box>

      <Snackbar open={open} autoHideDuration={3000} onClose={() => setOpen(false)}>
        <Alert severity={status.type as any} sx={{ width: '100%' }}>
          {status.msg}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default BusDetail;
