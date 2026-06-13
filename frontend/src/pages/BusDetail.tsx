import React, { useEffect, useState } from 'react';
import { 
  Box, 
  Typography, 
  Button, 
  Paper, 
  Stack, 
  Divider,
  IconButton
} from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  ArrowBack, 
  DirectionsBus, 
  LocationOn, 
  AccessTime, 
  EventSeat
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
    navigate(`/confirm-order/${id}`);
  };

  if (loading || !bus) return <Typography sx={{ p: 4 }}>加载中...</Typography>;

  return (
    <Box sx={{ pb: 10 }}>
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
    </Box>
  );
};

export default BusDetail;
