import React, { useEffect, useState } from 'react';
import { 
  Box, 
  Card, 
  CardContent, 
  Typography, 
  Button, 
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
  Search, 
  ChevronRight, 
  AccessTime,
  Person
} from '@mui/icons-material';
import { useAuth } from '../context/AuthContext';

interface Order {
  id: number;
  bus_id: number;
  status: number;
  created_at: string;
}

const Home: React.FC = () => {
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);
  const [origin, setOrigin] = useState('校区 A');
  const [dest, setDest] = useState('校区 B');
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  const handleSearch = () => {
    navigate(`/search-results?origin=${encodeURIComponent(origin)}&dest=${encodeURIComponent(dest)}&date=${date}`);
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
        bgcolor: 'background.paper', 
        border: '1px solid',
        borderColor: 'divider',
        boxShadow: '0 4px 12px rgba(0,0,0,0.05)'
      }}>
        <Stack spacing={3}>
          <Typography variant="h5" sx={{ fontWeight: 'bold', color: 'primary.main' }}>
            去哪里？
          </Typography>
          <Stack spacing={2}>
            <TextField
              select
              fullWidth
              label="出发地"
              value={origin}
              onChange={(e) => setOrigin(e.target.value)}
            >
              <MenuItem value="校区 A">校区 A</MenuItem>
              <MenuItem value="校区 B">校区 B</MenuItem>
              <MenuItem value="高铁站">高铁站</MenuItem>
            </TextField>
            <TextField
              select
              fullWidth
              label="目的地"
              value={dest}
              onChange={(e) => setDest(e.target.value)}
            >
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
                inputLabel: { shrink: true }
              }}
            />
          </Stack>
          <Button 
            variant="contained" 
            size="large"
            startIcon={<Search />} 
            onClick={handleSearch}
            fullWidth
            sx={{ 
              fontWeight: 'bold',
              py: 1.5,
              borderRadius: 2
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
    </Box>
  );
};

export default Home;
