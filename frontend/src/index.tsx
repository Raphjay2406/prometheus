import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css'; // Your global styles
import App from './App'; // Your main application component
import { ThemeProvider } from '@mui/material/styles'; // MUI Theme provider
import theme from './muiTheme'; // Your custom MUI theme configuration
import { BrowserRouter as Router } from 'react-router-dom'; // React Router for navigation
import { AuthProvider } from './contexts/AuthContext'; // Your authentication context
import reportWebVitals from './reportWebVitals'; // CRA's performance reporting utility

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

root.render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <Router>
        {/* AuthProvider MUST wrap App and any component that uses useAuth */}
        <AuthProvider>
          <App />
        </AuthProvider>
      </Router>
    </ThemeProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();