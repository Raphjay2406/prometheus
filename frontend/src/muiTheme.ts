// frontend/src/muiTheme.ts
import { createTheme, responsiveFontSizes } from '@mui/material/styles';
import { red } from '@mui/material/colors'; // Example color import

// A custom theme for this app
let theme = createTheme({
  palette: {
    primary: {
      main: '#1890ff', // Your primary brand color (similar to previous 'brand.500')
    },
    secondary: {
      main: '#096dd9', // A secondary color
    },
    error: {
      main: red.A400, // Standard error color from MUI
    },
    // You can also define mode-specific palettes for light/dark mode
    // background: {
    //   default: '#fff',
    // },
  },
  typography: {
    fontFamily: 'Roboto, Arial, sans-serif',
    h1: {
      fontSize: '2.5rem',
      fontWeight: 500,
    },
    // Add other typography variants as needed
  },
  // You can customize components globally here
  // components: {
  //   MuiButton: {
  //     styleOverrides: {
  //       root: {
  //         borderRadius: 8,
  //       },
  //     },
  //   },
  //   MuiTextField: {
  //     defaultProps: {
  //       variant: 'outlined',
  //       margin: 'normal',
  //     }
  //   }
  // },
});

// Make typography responsive
theme = responsiveFontSizes(theme);

export default theme;