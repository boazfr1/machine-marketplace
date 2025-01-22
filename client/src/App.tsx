import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LoginForm from './identification/LoginForm';
import SignupForm from './identification/SignupForm';
import Dashboard from './dashboard/WelcomePage';
import Feed from './tradeAndOparate/Feed';
import MachinePage from './tradeAndOparate/MachinePage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/login" element={<LoginForm />} />
        <Route path="/signup" element={<SignupForm />} />
        <Route path="/feed" element={<Feed />} />
        <Route path="/machine" element={<MachinePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;