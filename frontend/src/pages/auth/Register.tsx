import React, { useState, FormEvent, useEffect, ChangeEvent } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import {
  Box,
  Button,
  TextField,
  Typography,
  Container,
  Stack,
  InputAdornment,
  IconButton,
  CircularProgress,
  Alert,
  Snackbar,
} from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material';

import { useAuth } from '../../contexts/AuthContext';

const RegisterPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [localValidationError, setLocalValidationError] = useState<string | null>(null);

  const { register, isLoading, error, isAuthenticated, clearError } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/dashboard');
    }
  }, [isAuthenticated, navigate]);

  useEffect(() => {
    return () => {
      clearError();
      setLocalValidationError(null);
    };
  }, [clearError]);

  const handleSnackbarClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
    if (reason === 'clickaway') {
      return;
    }
    setSnackbarOpen(false);
  };

  const handleUsernameChange = (event: ChangeEvent<HTMLInputElement>) => {
    setUsername(event.target.value);
    setLocalValidationError(null);
    if (error) clearError();
  };

  const handleEmailChange = (event: ChangeEvent<HTMLInputElement>) => {
    setEmail(event.target.value);
    setLocalValidationError(null);
    if (error) clearError();
  };

  const handlePasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    setPassword(event.target.value);
    setLocalValidationError(null);
    if (error) clearError();
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setLocalValidationError(null);
    if (error) clearError();

    if (!username.trim() || !email.trim() || !password.trim()) {
      setSnackbarOpen(false);
      setLocalValidationError('Username, email, and password cannot be empty.');
      return;
    }

    if (password.length < 6) {
      setSnackbarOpen(false);
      setLocalValidationError('Password must be at least 6 characters long.');
      return;
    }

    try {
      await register({ username, email, password });
      setSnackbarOpen(true);
      // After successful registration and auto-login, AuthContext will navigate to dashboard
    } catch (err) {
      setSnackbarOpen(false);
    }
  };

  return (
    <Container component="main" maxWidth="sm" sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', minHeight: '100vh', justifyContent: 'center', bgcolor: 'grey.100' }}>
      <Box
        sx={{
          padding: 4,
          border: 1,
          borderRadius: 'lg',
          boxShadow: 3,
          bgcolor: 'background.paper',
          width: '100%',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Stack spacing={3} sx={{ width: '100%' }}>
          <Typography component="h1" variant="h4" align="center" color="primary.main">
            Prometheus Register
          </Typography>
          <Typography variant="body1" align="center" color="text.secondary">
            Create your HR account.
          </Typography>

          <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1, width: '100%' }}>
            <Stack spacing={2} sx={{ width: '100%' }}>
              <TextField
                margin="normal"
                required
                fullWidth
                id="username"
                label="Username"
                name="username"
                autoComplete="username"
                value={username}
                onChange={handleUsernameChange}
                error={!!localValidationError || (!!error && error.toLowerCase().includes('username'))}
                helperText={
                  (!!localValidationError || (!!error && error.toLowerCase().includes('username')))
                    ? (localValidationError || error)
                    : ''
                }
              />
              <TextField
                margin="normal"
                required
                fullWidth
                id="email"
                label="Email Address"
                name="email"
                autoComplete="email"
                value={email}
                onChange={handleEmailChange}
                error={!!localValidationError || (!!error && error.toLowerCase().includes('email'))}
                helperText={
                  (!!localValidationError || (!!error && error.toLowerCase().includes('email')))
                    ? (localValidationError || error)
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
                autoComplete="new-password"
                value={password}
                onChange={handlePasswordChange}
                error={!!localValidationError || (!!error && error.toLowerCase().includes('password'))}
                helperText={
                  (!!localValidationError || (!!error && error.toLowerCase().includes('password')))
                    ? (localValidationError || error)
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

              {error && !localValidationError && !(error.toLowerCase().includes('username') || error.toLowerCase().includes('email') || error.toLowerCase().includes('password')) && (
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
                {isLoading ? 'Registering...' : 'Register'}
              </Button>
            </Stack>
          </Box>

          <Typography variant="body2" align="center" mt={2}>
            Already have an account?{' '}
            <RouterLink to="/login" style={{ textDecoration: 'none' }}>
              <Typography component="span" color="primary" sx={{ fontWeight: 'medium' }}>
                Sign In
              </Typography>
            </RouterLink>
          </Typography>
        </Stack>
      </Box>

      <Snackbar open={snackbarOpen} autoHideDuration={3000} onClose={handleSnackbarClose} anchorOrigin={{ vertical: 'top', horizontal: 'center' }}>
        <Alert onClose={handleSnackbarClose} severity="success" sx={{ width: '100%' }}>
          Registration Successful! Redirecting to dashboard...
        </Alert>
      </Snackbar>
    </Container>
  );
};

export default RegisterPage;