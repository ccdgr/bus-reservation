import React, { useEffect, useState } from 'react';
import { 
  Box, 
  Card, 
  CardContent, 
  Typography, 
  Grid, 
  Button, 
  Skeleton,
  Stack,
  Chip
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import client from '../api/client';
import { DirectionsBus, EventSeat } from '@mui/icons-material';

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
  const navigate = useNavigate();

  useEffect(() => {
    client.get('/buses')
      .then(res => {
        setBuses(res.data);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);

  if (loading) {
    return (
      <Stack spacing={2}>
        {[1, 2, 3].map(i => (
          <Skeleton key={i} variant="rectangular" height={120} sx={{ borderRadius: 3 }} />
        ))}
      </Stack>
    );
  }

  return (
    <Box sx={{ pb: 2 }}>
      <Typography variant="h5" sx={{ mb: 2, fontWeight: 'bold' }}>
        今日班次
      </Typography>
      <Stack spacing={2}>
        {buses.map(bus => (
          <Card key={bus.id} onClick={() => navigate(`/bus/${bus.id}`)} sx={{ cursor: 'pointer' }}>
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
                    发车时间: {new Date(bus.start_time).toLocaleString()}
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
                  <Button variant="contained" size="small" fullWidth>
                    预定
                  </Button>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        ))}
      </Stack>
    </Box>
  );
};

export default Home;
