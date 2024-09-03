import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Dashboard from './components/Dashboard';

function App() {
  return(
    <main>
        <Router>
            <Routes>
                <Route exact path='/' element={<Dashboard />} />
                <Route exact path='/register' element={<Dashboard />} />
                <Route exact path='/library' element={<Dashboard />} />
                <Route exact path='/preferences' element={<Dashboard />} />
                <Route exact path='/recommend' element={<Dashboard />} />
                <Route exact path='/books/search' element={<Dashboard />} />
                <Route exact path='/trending' element={<Dashboard />} />
                <Route exact path='/admin/book/new' element={<Dashboard />}/>
            </Routes>
        </Router>
    </main>
)
}

export default App;
