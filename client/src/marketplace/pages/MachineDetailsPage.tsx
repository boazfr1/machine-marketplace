import { useLocation, useNavigate } from 'react-router-dom';
import MarketplaceLayout from '../components/MarketplaceLayout';
import { MachineType } from '../../global/types';
import './MachineDetailsPage.css';

const MachineDetailsPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const machine = location.state?.machine as MachineType;

  if (!machine) {
    navigate('/marketplace');
    return null;
  }

  const handleRentMachine = () => {
    console.log("Rent machine:", machine.Name);
    // Add rent logic here
  };

  const handleContactOwner = () => {
    console.log("Contact owner:", machine.OwnerID);
    // Add contact logic here
  };

  const specs = [
    { icon: '🧠', label: 'CPU Cores', value: machine.Cpu, unit: 'cores' },
    { icon: '💾', label: 'RAM', value: machine.Ram, unit: 'GB' },
    { icon: '💿', label: 'Storage', value: machine.Memory, unit: 'GB' },
    { icon: '⚡', label: 'Performance', value: 'High', unit: '' },
    { icon: '🌐', label: 'Network', value: '1 Gbps', unit: '' },
    { icon: '🔒', label: 'Security', value: 'Enterprise', unit: '' }
  ];

  const pricingTiers = [
    { duration: 'Hourly', price: 0.25, popular: false },
    { duration: 'Daily', price: 5.50, popular: true, savings: '8%' },
    { duration: 'Weekly', price: 35.00, popular: false, savings: '17%' },
    { duration: 'Monthly', price: 120.00, popular: false, savings: '33%' }
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
              <button 
                onClick={handleRentMachine}
                className="btn btn-primary btn-lg"
              >
                Rent This Machine
              </button>
            </div>
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
              <h3 className="pricing-title">Pricing Options</h3>
              <div className="pricing-tiers">
                {pricingTiers.map((tier, index) => (
                  <div 
                    key={index} 
                    className={`pricing-tier ${tier.popular ? 'pricing-tier-popular' : ''}`}
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
              
              <button 
                onClick={handleRentMachine}
                className="btn btn-primary w-full btn-lg"
              >
                Start Rental
              </button>
              
              <div className="pricing-note">
                <p>All prices include infrastructure costs and 24/7 support</p>
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
