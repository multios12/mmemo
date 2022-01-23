import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import TopBar from './Components/TopBar'
import './App.css';
import ListView from './Views/ListView';
import EditView from './Views/EditView';

class App extends React.Component {
  render() {
    return (
      <div>
        <Router>
          <TopBar />
          <Routes>
            <Route path='/' element={<ListView />} />
            <Route path='/add' element={<EditView />} />
            <Route path='/memos/:id' element={<EditView />} />
          </Routes>
        </Router>
      </div>
    );
  }
}

//const darkTheme = createTheme({
//  palette: {
//    mode: 'dark',
//  },
//});

export default App;
