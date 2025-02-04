import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LoginForm from './identification/LoginForm';
import SignupForm from './identification/SignupForm';
import Dashboard from './dashboard/WelcomePage';
import Feed from './marketplace/pages/Feed';
import MachineDetails from './marketplace/pages/MachineDetails';
import MyMachinesPage from './marketplace/pages/OwnedMachines';
import AssignMachine from './marketplace/pages/AddMachine';
import BoughtMachines from './marketplace/pages/BoughtMachines';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/login" element={<LoginForm />} />
        <Route path="/signup" element={<SignupForm />} />
        <Route path="/feed" element={<Feed />} />
        <Route path="/machine" element={<MachineDetails />} />
        <Route path="/my-machines" element={<MyMachinesPage />} />
        <Route path="/assign-machine" element={<AssignMachine />} />
        <Route path="/assigned-machines" element={<BoughtMachines />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;