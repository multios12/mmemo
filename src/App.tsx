import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import TopBar from './Components/TopBar'
import './App.css';
import ListView from './Views/ListView';
import EditView from './Views/EditView';
import { createTheme, ThemeProvider } from '@mui/material/styles';

export default function App() {
  const theme = createTheme({ palette: { mode: 'dark' } })
  return (
    <ThemeProvider theme={theme}>
      <Router>
        <TopBar />
        <Routes>
          <Route path='/' element={<ListView />} />
          <Route path='/add' element={<EditView />} />
          <Route path='/memos/:id' element={<EditView />} />
        </Routes>
      </Router>
    </ThemeProvider>
  );
}
