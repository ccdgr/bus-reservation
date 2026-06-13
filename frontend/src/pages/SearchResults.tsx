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
  IconButton
} from '@mui/material';
import { useNavigate, useSearchParams } from 'react-router-dom';
import client from '../api/client';
import { DirectionsBus, EventSeat, ArrowBack } from '@mui/icons-material';

interface Bus {
  id: number;
  number: string;
  origin: string;
  dest: string;
  start_time: string;
  total_seat: number;
  left_seat: number;
}

const SearchResults: React.FC = () => {
  const [buses, setBuses] = useState<Bus[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const origin = searchParams.get('origin') || '';
  const dest = searchParams.get('dest') || '';
  const date = searchParams.get('date') || '';

  useEffect(() => {
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
  }, [origin, dest, date]);

  return (
    <Box sx={{ pb: 4 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 3 }}>
        <IconButton onClick={() => navigate(-1)} sx={{ mr: 1 }}>
          <ArrowBack />
        </IconButton>
        <Box>
          <Typography variant="h6" sx={{ fontWeight: 'bold' }}>班次查询结果</Typography>
          <Typography variant="caption" color="text.secondary">
            {origin} ➔ {dest} | {date}
          </Typography>
        </Box>
      </Box>

      {loading ? (
        <Stack spacing={2}>
          {[1, 2, 3, 4].map(i => (
            <Skeleton key={i} variant="rectangular" height={100} sx={{ borderRadius: 3 }} />
          ))}
        </Stack>
      ) : (
        <Stack spacing={2}>
          {buses.length === 0 ? (
            <Box sx={{ textAlign: 'center', py: 8 }}>
              <Typography color="text.secondary" sx={{ mb: 2 }}>
                未找到相关班次
              </Typography>
              <Button variant="outlined" onClick={() => navigate(-1)}>
                返回重新搜索
              </Button>
            </Box>
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
                      <Typography color="text.secondary" variant="body2">
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
                      <Button variant="outlined" size="small" fullWidth sx={{ borderRadius: 2 }}>
                        预定
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

export default SearchResults;
