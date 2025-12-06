import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LoginPage from './auth/LoginPage';
import SignupPage from './auth/SignupPage';
import LandingPage from './landing/LandingPage';
import MarketplaceFeed from './marketplace/pages/MarketplaceFeed';
import MachineDetailsPage from './marketplace/pages/MachineDetailsPage';
import MyMachinesPage from './marketplace/pages/MyMachinesPage';
import AddMachinePage from './marketplace/pages/AddMachinePage';
import PurchasedMachinesPage from './marketplace/pages/PurchasedMachinesPage';

function MachineMarketplaceApp() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/signup" element={<SignupPage />} />
        <Route path="/marketplace" element={<MarketplaceFeed />} />
        <Route path="/machine/:id" element={<MachineDetailsPage />} />
        <Route path="/my-machines" element={<MyMachinesPage />} />
        <Route path="/add-machine" element={<AddMachinePage />} />
        <Route path="/purchased-machines" element={<PurchasedMachinesPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default MachineMarketplaceApp;