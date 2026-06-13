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
  Grid
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import { DirectionsBus, EventSeat, Search } from '@mui/icons-material';

interface Bus {
  id: number;
  number: string;
  origin: string;
  dest: string;
  start_time: string;
  total_seat: number;
  left_seat: number;
}

const Home: React.FC = () => {
  const [buses, setBuses] = useState<Bus[]>([]);
  const [loading, setLoading] = useState(true);
  const [origin, setOrigin] = useState('');
  const [dest, setDest] = useState('');
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const navigate = useNavigate();

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

  useEffect(() => {
    fetchBuses();
  }, []);

  return (
    <Box sx={{ pb: 2 }}>
      <Paper sx={{ p: 2, mb: 3, borderRadius: 3 }}>
        <Stack spacing={2}>
          <Typography variant="subtitle2" color="primary" sx={{ fontWeight: 'bold' }}>
            班次查询
          </Typography>
          <Grid container spacing={2}>
            <Grid size={{ xs: 6 }}>
              <TextField
                select
                fullWidth
                label="出发地"
                size="small"
                value={origin}
                onChange={(e) => setOrigin(e.target.value)}
              >
                <MenuItem value="">全部</MenuItem>
                <MenuItem value="校区 A">校区 A</MenuItem>
                <MenuItem value="校区 B">校区 B</MenuItem>
                <MenuItem value="高铁站">高铁站</MenuItem>
              </TextField>
            </Grid>
            <Grid size={{ xs: 6 }}>
              <TextField
                select
                fullWidth
                label="目的地"
                size="small"
                value={dest}
                onChange={(e) => setDest(e.target.value)}
              >
                <MenuItem value="">全部</MenuItem>
                <MenuItem value="校区 A">校区 A</MenuItem>
                <MenuItem value="校区 B">校区 B</MenuItem>
                <MenuItem value="高铁站">高铁站</MenuItem>
              </TextField>
            </Grid>
            <Grid size={{ xs: 12 }}>
              <TextField
                fullWidth
                type="date"
                label="出发日期"
                size="small"
                value={date}
                onChange={(e) => setDate(e.target.value)}
                slotProps={{ inputLabel: { shrink: true } }}
              />
            </Grid>
          </Grid>
          <Button 
            variant="contained" 
            startIcon={<Search />} 
            onClick={fetchBuses}
            fullWidth
          >
            搜索
          </Button>
        </Stack>
      </Paper>

      <Typography variant="h6" sx={{ mb: 2, fontWeight: 'bold' }}>
        可用班次
      </Typography>

      {loading ? (
        <Stack spacing={2}>
          {[1, 2, 3].map(i => (
            <Skeleton key={i} variant="rectangular" height={120} sx={{ borderRadius: 3 }} />
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
                  <Grid container sx={{ alignItems: 'center' }}>
                    <Grid size={{ xs: 8 }}>
                      <Stack direction="row" spacing={1} sx={{ alignItems: 'center' }}>
                        <DirectionsBus color="primary" />
                        <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                          {bus.number}
                        </Typography>
                      </Stack>
                      <Typography color="text.secondary">
                        {bus.origin} ➔ {bus.dest}
                      </Typography>
                      <Typography variant="body2" sx={{ mt: 1 }}>
                        时间: {new Date(bus.start_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                      </Typography>
                    </Grid>
                    <Grid size={{ xs: 4 }} sx={{ textAlign: 'right' }}>
                      <Chip 
                        label={`余 ${bus.left_seat}`} 
                        color={bus.left_seat > 0 ? 'success' : 'error'}
                        size="small"
                        icon={<EventSeat />}
                        sx={{ mb: 1 }}
                      />
                      <Button variant="outlined" size="small" fullWidth>
                        详情
                      </Button>
                    </Grid>
                  </Grid>
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
