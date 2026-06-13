import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, CssBaseline } from '@mui/material';
import theme from './theme/theme';
import { AuthProvider, useAuth } from './context/AuthContext';
import Layout from './components/Layout';

// Pages
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Orders from './pages/Orders';
import Profile from './pages/Profile';
import BusDetail from './pages/BusDetail';
import SearchResults from './pages/SearchResults';
import ConfirmOrder from './pages/ConfirmOrder';

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
};

const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            
            <Route path="/" element={<Layout />}>
              <Route index element={<Home />} />
              <Route path="search-results" element={<SearchResults />} />
              <Route path="bus/:id" element={<BusDetail />} />
              <Route 
                path="confirm-order/:busId" 
                element={
                  <PrivateRoute>
                    <ConfirmOrder />
                  </PrivateRoute>
                } 
              />
              <Route 
                path="orders" 
                element={
                  <PrivateRoute>
                    <Orders />
                  </PrivateRoute>
                } 
              />
              <Route 
                path="profile" 
                element={
                  <PrivateRoute>
                    <Profile />
                  </PrivateRoute>
                } 
              />
            </Route>
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
};

export default App;
