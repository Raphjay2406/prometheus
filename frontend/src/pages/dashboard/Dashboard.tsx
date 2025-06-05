import React from 'react';
import { Box, Typography, Button, Container } from '@mui/material';
import { useAuth } from '../../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';

const DashboardPage: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login'); // Redirect to login page after logout
  };

  return (
    <Container component="main" maxWidth="md" sx={{ mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
      <Box sx={{ padding: 4, borderRadius: 'lg', boxShadow: 3, bgcolor: 'background.paper', width: '100%', textAlign: 'center' }}>
        <Typography component="h1" variant="h4" color="primary.main" gutterBottom>
          Welcome to Your Dashboard, {user?.username || user?.email}!
        </Typography>
        <Typography variant="body1" color="text.secondary" paragraph>
          Your role: {user?.role_name || 'N/A'}
        </Typography>
        <Typography variant="body2" color="text.secondary">
          This is a protected dashboard page. Only authenticated users can see this.
        </Typography>
        <Button variant="contained" color="secondary" onClick={handleLogout} sx={{ mt: 3 }}>
          Logout
        </Button>
      </Box>
    </Container>
  );
};

export default DashboardPage;