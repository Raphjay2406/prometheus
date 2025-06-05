import React from 'react';
import { Routes, Route } from 'react-router-dom';
import LoginPage from './pages/auth/Login'; // Your Login component
import DashboardPage from './pages/dashboard/Dashboard'; // Placeholder for dashboard
import RegisterPage from './pages/auth/Register'; // Placeholder for registration
import { useAuth } from './contexts/AuthContext';
import { CircularProgress, Box } from '@mui/material';

const App: React.FC = () => {
  const { isLoading, isAuthenticated } = useAuth();

  if (isLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Routes>
      <Route path="/" element={<LoginPage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      {/* Protected route example */}
      {isAuthenticated && <Route path="/dashboard" element={<DashboardPage />} />}
      {/* Redirect unauthenticated users from protected routes, or show a 404 */}
      {!isAuthenticated && <Route path="/dashboard" element={<LoginPage />} />}
      {/* Add other routes here as you develop more pages */}
      {/* <Route path="*" element={<NotFoundPage />} /> */} {/* Optional 404 page */}
    </Routes>
  );
};

export default App;