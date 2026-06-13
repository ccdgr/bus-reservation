import React, { useEffect, useState } from 'react';
import { 
  Box, 
  Card, 
  CardContent, 
  Typography, 
  Button, 
  Skeleton,
  Stack,
  Chip,
  TextField,
  MenuItem,
  Paper,
  IconButton
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import { 
  DirectionsBus, 
  EventSeat, 
  Search, 
  ChevronRight, 
  AccessTime,
  Person
} from '@mui/icons-material';
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

interface Order {
  id: number;
  bus_id: number;
  status: number;
  created_at: string;
}

const Home: React.FC = () => {
  const [buses, setBuses] = useState<Bus[]>([]);
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [origin, setOrigin] = useState('');
  const [dest, setDest] = useState('');
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const fetchBuses = () => {
    setLoading(true);
    client.get('/buses', {
      params: { origin, dest, date }
    })
      .then(res => {
        setBuses(res.data);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  const fetchRecentOrders = () => {
    if (!isAuthenticated) return;
    client.get('/orders')
      .then(res => {
        const ongoing = res.data.filter((o: Order) => o.status === 0 || o.status === 1).slice(0, 2);
        setRecentOrders(ongoing);
      });
  };

  useEffect(() => {
    fetchBuses();
    fetchRecentOrders();
  }, [isAuthenticated]);

  return (
    <Box sx={{ pb: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'flex-end', mb: 1 }}>
        <IconButton onClick={() => navigate('/profile')} color="primary">
          <Person />
        </IconButton>
      </Box>

      <Paper elevation={0} sx={{ 
        p: 3, 
        mb: 4, 
        borderRadius: 4, 
        bgcolor: 'primary.main', 
        color: 'white',
        boxShadow: '0 8px 24px rgba(25, 118, 210, 0.2)'
      }}>
        <Stack spacing={3}>
          <Typography variant="h5" sx={{ fontWeight: 'bold' }}>
            去哪里？
          </Typography>
          <Stack spacing={2}>
            <TextField
              select
              fullWidth
              label="出发地"
              value={origin}
              onChange={(e) => setOrigin(e.target.value)}
              sx={{ 
                bgcolor: 'rgba(255,255,255,0.1)', 
                borderRadius: 2,
                '& .MuiOutlinedInput-root': { color: 'white' },
                '& .MuiInputLabel-root': { color: 'rgba(255,255,255,0.7)' },
                '& .MuiOutlinedInput-notchedOutline': { borderColor: 'rgba(255,255,255,0.3)' }
              }}
            >
              <MenuItem value="">全部</MenuItem>
              <MenuItem value="校区 A">校区 A</MenuItem>
              <MenuItem value="校校区 B">校区 B</MenuItem>
              <MenuItem value="高铁站">高铁站</MenuItem>
            </TextField>
            <TextField
              select
              fullWidth
              label="目的地"
              value={dest}
              onChange={(e) => setDest(e.target.value)}
              sx={{ 
                bgcolor: 'rgba(255,255,255,0.1)', 
                borderRadius: 2,
                '& .MuiOutlinedInput-root': { color: 'white' },
                '& .MuiInputLabel-root': { color: 'rgba(255,255,255,0.7)' },
                '& .MuiOutlinedInput-notchedOutline': { borderColor: 'rgba(255,255,255,0.3)' }
              }}
            >
              <MenuItem value="">全部</MenuItem>
              <MenuItem value="校区 A">校区 A</MenuItem>
              <MenuItem value="校区 B">校区 B</MenuItem>
              <MenuItem value="高铁站">高铁站</MenuItem>
            </TextField>
            <TextField
              fullWidth
              type="date"
              label="出发日期"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              slotProps={{ 
                inputLabel: { shrink: true, sx: { color: 'rgba(255,255,255,0.7)' } },
                htmlInput: { sx: { color: 'white' } }
              }}
              sx={{ 
                bgcolor: 'rgba(255,255,255,0.1)', 
                borderRadius: 2,
                '& .MuiOutlinedInput-notchedOutline': { borderColor: 'rgba(255,255,255,0.3)' }
              }}
            />
          </Stack>
          <Button 
            variant="contained" 
            size="large"
            startIcon={<Search />} 
            onClick={fetchBuses}
            fullWidth
            sx={{ 
              bgcolor: 'white', 
              color: 'primary.main',
              fontWeight: 'bold',
              py: 1.5,
              '&:hover': { bgcolor: 'rgba(255,255,255,0.9)' }
            }}
          >
            立即查询
          </Button>
        </Stack>
      </Paper>

      {isAuthenticated && (
        <Box sx={{ mb: 4 }}>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
            <Typography variant="h6" sx={{ fontWeight: 'bold' }}>进行中的订单</Typography>
            <Button 
              size="small" 
              endIcon={<ChevronRight />} 
              onClick={() => navigate('/orders')}
            >
              全部
            </Button>
          </Box>
          
          {recentOrders.length === 0 ? (
            <Typography variant="body2" color="text.secondary">暂无进行中的订单</Typography>
          ) : (
            <Stack spacing={2}>
              {recentOrders.map(order => (
                <Card key={order.id} variant="outlined" sx={{ borderRadius: 3 }}>
                  <CardContent sx={{ py: '12px !important' }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <Box>
                        <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>订单 #{order.id}</Typography>
                        <Stack direction="row" spacing={0.5} sx={{ alignItems: 'center' }}>
                          <AccessTime sx={{ fontSize: 14, color: 'text.secondary' }} />
                          <Typography variant="caption" color="text.secondary">
                            {new Date(order.created_at).toLocaleString()}
                          </Typography>
                        </Stack>
                      </Box>
                      <Chip 
                        label={order.status === 0 ? '处理中' : '预定成功'} 
                        color={order.status === 0 ? 'warning' : 'success'} 
                        size="small" 
                      />
                    </Box>
                  </CardContent>
                </Card>
              ))}
            </Stack>
          )}
        </Box>
      )}

      <Typography variant="h6" sx={{ mb: 2, fontWeight: 'bold' }}>
        可用班次
      </Typography>

      {loading ? (
        <Stack spacing={2}>
          {[1, 2, 3].map(i => (
            <Skeleton key={i} variant="rectangular" height={100} sx={{ borderRadius: 3 }} />
          ))}
        </Stack>
      ) : (
        <Stack spacing={2}>
          {buses.length === 0 ? (
            <Typography align="center" color="text.secondary" sx={{ py: 4 }}>
              未找到相关班次
            </Typography>
          ) : (
            buses.map(bus => (
              <Card key={bus.id} onClick={() => navigate(`/bus/${bus.id}`)} sx={{ cursor: 'pointer', borderRadius: 3 }}>
                <CardContent>
                  <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                    <Box sx={{ flexGrow: 1 }}>
                      <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
                        <DirectionsBus color="primary" />
                        <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                          {bus.number}
                        </Typography>
                      </Stack>
                      <Typography color="text.secondary">
                        {bus.origin} ➔ {bus.dest}
                      </Typography>
                      <Typography variant="body2" sx={{ mt: 0.5 }}>
                        时间: {new Date(bus.start_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                      </Typography>
                    </Box>
                    <Box sx={{ textAlign: 'right', ml: 2 }}>
                      <Chip 
                        label={`余 ${bus.left_seat}`} 
                        color={bus.left_seat > 0 ? 'success' : 'error'}
                        size="small"
                        icon={<EventSeat sx={{ fontSize: 16 }} />}
                        sx={{ mb: 1 }}
                      />
                      <Button variant="outlined" size="small" fullWidth>
                        详情
                      </Button>
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            ))
          )}
        </Stack>
      )}
    </Box>
  );
};

export default Home;
