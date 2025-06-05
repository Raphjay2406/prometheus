// frontend/src/pages/auth/Login.tsx
import React, { useState, FormEvent, useEffect, ChangeEvent } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import {
  Box,
  Button,
  TextField, // MUI's input component
  Typography, // MUI's text component (replaces Heading, Text)
  Container,  // For centering and max-width
  Stack,      // Replaces VStack for vertical stacking
  InputAdornment, // For input suffixes/prefixes like password visibility toggle
  IconButton, // For icon buttons
  CircularProgress, // For loading spinner inside button
  Alert,      // For displaying error messages
  Snackbar,   // For transient notifications (like Chakra's Toast)
} from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material'; // MUI icons for password visibility

import { useAuth } from '../../contexts/AuthContext'; // Your AuthContext

const LoginPage: React.FC = () => {
  const [identifier, setIdentifier] = useState(''); // Can be username or email
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [snackbarOpen, setSnackbarOpen] = useState(false); // For Snackbar success message
  const [localValidationError, setLocalValidationError] = useState<string | null>(null); // NEW: For client-side validation errors

  const { login, isLoading, error, isAuthenticated, clearError } = useAuth(); // 'error' from useAuth is for backend errors
  const navigate = useNavigate();

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard'); // Or your desired authenticated route
    }
  }, [isAuthenticated, navigate]);

  // Clear errors when the component unmounts or input fields change
  useEffect(() => {
    return () => {
      clearError(); // Clear any existing auth errors from context on unmount
      setLocalValidationError(null); // Clear any local errors on unmount
    };
  }, [clearError]); // Depend on clearError which is useCallback-wrapped

  // Snackbar close handler
  const handleSnackbarClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackbarOpen(false);
  };

  const handleIdentifierChange = (event: ChangeEvent<HTMLInputElement>) => {
    setIdentifier(event.target.value);
    setLocalValidationError(null); // Clear local error on input change
    if (error) clearError(); // Clear auth context error on input change
  };

  const handlePasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    setPassword(event.target.value);
    setLocalValidationError(null); // Clear local error on input change
    if (error) clearError(); // Clear auth context error on input change
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLocalValidationError(null); // Clear previous local error before a new attempt
    if (error) clearError(); // Clear previous auth context error

    if (!identifier.trim() || !password.trim()) {
      setSnackbarOpen(false); // Close any existing success snackbar
      setLocalValidationError('Username/Email and password cannot be empty.'); // SETTING LOCAL ERROR
      return;
    }

    try {
      await login({ username: identifier, password }); // AuthContext's login handles backend call and its own 'error' state
      setSnackbarOpen(true); // Open success snackbar if login (and subsequent navigation) is successful
    } catch (err) {
      // Error is already handled by AuthContext, its 'error' state will be set.
      setSnackbarOpen(false); // Ensure success snackbar is not open
    }
  };

  return (
    <Container component="main" maxWidth="sm" sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', minHeight: '100vh', justifyContent: 'center', bgcolor: 'grey.100' }}>
      <Box
        sx={{
          padding: 4, // Equivalent to Chakra's p={8} (2*8=16px)
          border: 1, // borderWidth={1}
          borderRadius: 'lg', // borderRadius="lg"
          boxShadow: 3, // boxShadow="xl" (MUI shadows are 0-24)
          bgcolor: 'background.paper', // bg="white"
          width: '100%',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Stack spacing={3} sx={{ width: '100%' }}>
          <Typography component="h1" variant="h4" align="center" color="primary.main">
            Prometheus Login
          </Typography>
          <Typography variant="body1" align="center" color="text.secondary">
            Sign in to access your HR dashboard.
          </Typography>

          <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1, width: '100%' }}>
            <Stack spacing={2} sx={{ width: '100%' }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="identifier"
                label="Username or Email"
                name="identifier"
                autoComplete="username"
                value={identifier}
                onChange={handleIdentifierChange}
                // Check for local validation error or a backend error related to credentials/username
                error={!!localValidationError || (!!error && (error.toLowerCase().includes('username') || error.toLowerCase().includes('credentials') || error.toLowerCase().includes('empty')))}
                helperText={
                  (!!localValidationError || (!!error && (error.toLowerCase().includes('username') || error.toLowerCase().includes('credentials') || error.toLowerCase().includes('empty'))))
                    ? (localValidationError || error) // Display local error if present, otherwise backend error
                    : ''
                }
              />

              <TextField
                margin="normal"
                required
                fullWidth
                name="password"
                label="Password"
                type={showPassword ? 'text' : 'password'}
                id="password"
                autoComplete="current-password"
                value={password}
                onChange={handlePasswordChange}
                // Check for local validation error or a backend error related to credentials/password
                error={!!localValidationError || (!!error && (error.toLowerCase().includes('password') || error.toLowerCase().includes('credentials') || error.toLowerCase().includes('empty')))}
                helperText={
                  (!!localValidationError || (!!error && (error.toLowerCase().includes('password') || error.toLowerCase().includes('credentials') || error.toLowerCase().includes('empty'))))
                    ? (localValidationError || error) // Display local error if present, otherwise backend error
                    : ''
                }
                InputProps={{
                  endAdornment: (
                    <InputAdornment position="end">
                      <IconButton
                        aria-label={showPassword ? 'hide password' : 'show password'}
                        onClick={() => setShowPassword(!showPassword)}
                        edge="end"
                      >
                        {showPassword ? <VisibilityOff /> : <Visibility />}
                      </IconButton>
                    </InputAdornment>
                  ),
                }}
              />

              {/* Display a general backend error message if it's not tied to specific fields or covered by local validation */}
              {error && !localValidationError && !(error.toLowerCase().includes('username') || error.toLowerCase().includes('password') || error.toLowerCase().includes('credentials') || error.toLowerCase().includes('empty')) && (
                <Alert severity="error" sx={{ width: '100%', mt: 2 }}>
                  {error}
                </Alert>
              )}

              <Button
                type="submit"
                fullWidth
                variant="contained"
                color="primary"
                sx={{ mt: 3, mb: 2 }}
                disabled={isLoading}
                startIcon={isLoading ? <CircularProgress size={20} color="inherit" /> : null}
              >
                {isLoading ? 'Signing In...' : 'Sign In'}
              </Button>
            </Stack>
          </Box>

          <Typography variant="body2" align="center" mt={2}>
            Don't have an account?{' '}
            <RouterLink to="/register" style={{ textDecoration: 'none' }}>
              <Typography component="span" color="primary" sx={{ fontWeight: 'medium' }}>
                Sign Up
              </Typography>
            </RouterLink>
          </Typography>
        </Stack>
      </Box>

      {/* Snackbar for success notifications */}
      <Snackbar open={snackbarOpen} autoHideDuration={3000} onClose={handleSnackbarClose} anchorOrigin={{ vertical: 'top', horizontal: 'center' }}>
        <Alert onClose={handleSnackbarClose} severity="success" sx={{ width: '100%' }}>
          Login Successful!
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default LoginPage;