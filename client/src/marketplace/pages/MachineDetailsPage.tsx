import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import MarketplaceLayout from '../components/MarketplaceLayout';
import { MachineType, PurchaseRequest, PurchaseResponse } from '../../global/types';
import { orderApi } from '../../global/api';
import './MachineDetailsPage.css';

const MachineDetailsPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const machine = location.state?.machine as MachineType;
  const [selectedDuration, setSelectedDuration] = useState<number>(24);
  const [isPurchasing, setIsPurchasing] = useState(false);
  const [purchaseSuccess, setPurchaseSuccess] = useState(false);
  const [purchaseError, setPurchaseError] = useState<string | null>(null);

  if (!machine) {
    navigate('/marketplace');
    return null;
  }

  const handlePurchaseMachine = async () => {
    try {
      setIsPurchasing(true);
      setPurchaseError(null);

      const purchaseRequest: PurchaseRequest = {
        machine_id: machine.ID,
        deal_duration_hours: selectedDuration
      };

      const { data } = await orderApi.post<PurchaseResponse>('/api/v1/order/buy', purchaseRequest);
      console.log("Purchase successful:", data);

      setPurchaseSuccess(true);
      setTimeout(() => {
        navigate('/purchased-machines');
      }, 2000);
    } catch (error: any) {
      console.error("Purchase failed:", error);
      setPurchaseError(error.response?.data || "Failed to purchase machine. Please try again.");
    } finally {
      setIsPurchasing(false);
    }
  };

  const handleContactOwner = () => {
    console.log("Contact owner:", machine.OwnerID);
    alert("Contact feature coming soon!");
  };

  const specs = [
    { icon: '🧠', label: 'CPU Cores', value: machine.Cpu, unit: 'cores' },
    { icon: '💾', label: 'RAM', value: machine.Ram, unit: 'GB' },
    { icon: '🎮', label: 'GPU', value: machine.Gpu > 0 ? machine.Gpu : 'None', unit: machine.Gpu > 0 ? 'GPUs' : '' },
    { icon: '💿', label: 'Storage', value: machine.Memory, unit: 'GB' },
    { icon: '🌐', label: 'Network', value: '1 Gbps', unit: '' },
    { icon: '🔒', label: 'Security', value: 'Enterprise', unit: '' }
  ];

  const pricingTiers = [
    { duration: 'Daily', hours: 24, price: 5.50, popular: false },
    { duration: '3 Days', hours: 72, price: 15.00, popular: true, savings: '9%' },
    { duration: 'Weekly', hours: 168, price: 35.00, popular: false, savings: '17%' },
    { duration: 'Monthly', hours: 720, price: 120.00, popular: false, savings: '33%' }
  ];

  return (
    <MarketplaceLayout>
      <div className="machine-details-page">
        <div className="details-header">
          <button 
            onClick={() => navigate('/marketplace')} 
            className="back-button"
          >
            ← Back to Marketplace
          </button>
          
          <div className="machine-header-info">
            <div className="machine-title-section">
              <h1 className="machine-title">{machine.Name}</h1>
              <div className="machine-meta">
                <span className="owner-info">
                  <span className="owner-label">Owned by</span>
                  <span className="owner-id">#{machine.OwnerID}</span>
                </span>
                <span className="status-badge status-available">Available Now</span>
              </div>
            </div>
            
            <div className="quick-actions">
              <button
                onClick={handleContactOwner}
                className="btn btn-outline"
              >
                Contact Owner
              </button>
              {purchaseSuccess ? (
                <button
                  className="btn btn-success btn-lg"
                  disabled
                >
                  ✓ Purchase Successful!
                </button>
              ) : (
                <button
                  onClick={handlePurchaseMachine}
                  className="btn btn-primary btn-lg"
                  disabled={isPurchasing}
                >
                  {isPurchasing ? 'Processing...' : 'Purchase Machine'}
                </button>
              )}
            </div>
            {purchaseError && (
              <div className="error-message" style={{ color: 'red', marginTop: '10px' }}>
                {purchaseError}
              </div>
            )}
          </div>
        </div>

        <div className="details-content">
          <div className="details-main">
            <section className="specs-section">
              <h2 className="section-title">Technical Specifications</h2>
              <div className="specs-grid">
                {specs.map((spec, index) => (
                  <div key={index} className="spec-card">
                    <div className="spec-icon">{spec.icon}</div>
                    <div className="spec-info">
                      <span className="spec-label">{spec.label}</span>
                      <span className="spec-value">
                        {spec.value} {spec.unit}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </section>

            <section className="description-section">
              <h2 className="section-title">Machine Description</h2>
              <div className="description-content">
                <p>
                  This high-performance computing machine is perfect for demanding workloads 
                  including machine learning, data processing, and scientific computing. 
                  With {machine.Cpu} CPU cores and {machine.Ram}GB of RAM, it can handle 
                  intensive parallel processing tasks with ease.
                </p>
                <p>
                  The machine features enterprise-grade hardware with reliable performance 
                  and 99.9% uptime guarantee. All data is encrypted at rest and in transit, 
                  ensuring your workloads remain secure.
                </p>
              </div>
            </section>

            <section className="features-section">
              <h2 className="section-title">Key Features</h2>
              <div className="features-list">
                <div className="feature-item">
                  <span className="feature-icon">✅</span>
                  <span className="feature-text">Instant deployment and access</span>
                </div>
                <div className="feature-item">
                  <span className="feature-icon">✅</span>
                  <span className="feature-text">24/7 monitoring and support</span>
                </div>
                <div className="feature-item">
                  <span className="feature-icon">✅</span>
                  <span className="feature-text">Flexible scaling options</span>
                </div>
                <div className="feature-item">
                  <span className="feature-icon">✅</span>
                  <span className="feature-text">Enterprise-grade security</span>
                </div>
                <div className="feature-item">
                  <span className="feature-icon">✅</span>
                  <span className="feature-text">High-speed network connectivity</span>
                </div>
                <div className="feature-item">
                  <span className="feature-icon">✅</span>
                  <span className="feature-text">Backup and disaster recovery</span>
                </div>
              </div>
            </section>
          </div>

          <div className="details-sidebar">
            <div className="pricing-card">
              <h3 className="pricing-title">Select Rental Duration</h3>
              <div className="pricing-tiers">
                {pricingTiers.map((tier, index) => (
                  <div
                    key={index}
                    className={`pricing-tier ${tier.popular ? 'pricing-tier-popular' : ''} ${selectedDuration === tier.hours ? 'pricing-tier-selected' : ''}`}
                    onClick={() => setSelectedDuration(tier.hours)}
                    style={{ cursor: 'pointer' }}
                  >
                    <div className="tier-header">
                      <span className="tier-duration">{tier.duration}</span>
                      {tier.popular && <span className="popular-badge">Most Popular</span>}
                    </div>
                    <div className="tier-price">
                      <span className="price-amount">${tier.price}</span>
                      <span className="price-unit">/{tier.duration.toLowerCase()}</span>
                    </div>
                    {tier.savings && (
                      <div className="tier-savings">Save {tier.savings}</div>
                    )}
                  </div>
                ))}
              </div>

              {purchaseSuccess ? (
                <button
                  className="btn btn-success w-full btn-lg"
                  disabled
                >
                  ✓ Purchase Successful!
                </button>
              ) : (
                <button
                  onClick={handlePurchaseMachine}
                  className="btn btn-primary w-full btn-lg"
                  disabled={isPurchasing}
                >
                  {isPurchasing ? 'Processing Purchase...' : `Purchase for ${selectedDuration}h`}
                </button>
              )}

              <div className="pricing-note">
                <p>Duration: {selectedDuration} hours • Payment processed securely</p>
                {purchaseError && (
                  <p style={{ color: 'red', marginTop: '10px' }}>{purchaseError}</p>
                )}
              </div>
            </div>

            <div className="owner-card">
              <h3 className="owner-card-title">Machine Owner</h3>
              <div className="owner-info-detailed">
                <div className="owner-avatar">
                  <span className="avatar-text">#{machine.OwnerID}</span>
                </div>
                <div className="owner-details">
                  <span className="owner-name">Provider #{machine.OwnerID}</span>
                  <span className="owner-rating">
                    ⭐ 4.8 (127 reviews)
                  </span>
                  <span className="owner-status">Verified Provider</span>
                </div>
              </div>
              
              <button 
                onClick={handleContactOwner}
                className="btn btn-outline w-full"
              >
                Send Message
              </button>
            </div>
          </div>
        </div>
      </div>
    </MarketplaceLayout>
  );
};

export default MachineDetailsPage;
