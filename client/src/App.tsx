import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LoginForm from './identification/LoginForm';
import SignupForm from './identification/SignupForm';
import Dashboard from './dashboard/WelcomePage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/login" element={<LoginForm />} />
        <Route path="/signup" element={<SignupForm />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;