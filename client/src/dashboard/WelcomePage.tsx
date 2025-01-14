import { useNavigate } from 'react-router-dom';
import './WelcomePage.css';

const WelcomePage = () => {
  const navigate = useNavigate();

  const handleLogin = () => {
    navigate('/login');
  };

  const handleSignup = () => {
    navigate('/signup');
  };

  return (
    <div className="welcome-container">
      <nav className="nav-bar">
        <a href="/" className="logo">
          MachineMart
        </a>
        <div className="nav-buttons">
          <button onClick={handleLogin} className="btn btn-outline">
            Login
          </button>
          <button onClick={handleSignup} className="btn btn-primary">
            Sign Up
          </button>
        </div>
      </nav>

      <section className="hero-section">
        <h1 className="hero-title">
          Welcome to MachineMart
        </h1>
        <p className="hero-subtitle">
          Your one-stop marketplace for industrial machinery. 
          Buy, sell, and discover quality equipment from trusted vendors.
        </p>
        <button onClick={handleSignup} className="btn btn-primary">
          Get Started
        </button>
      </section>

      <div className="features-grid">
        <div className="feature-card">
          <div className="feature-icon">🔍</div>
          <h3 className="feature-title">Easy Search</h3>
          <p className="feature-text">
            Find exactly what you need with our advanced search and filtering system.
          </p>
        </div>

        <div className="feature-card">
          <div className="feature-icon">🤝</div>
          <h3 className="feature-title">Secure Transactions</h3>
          <p className="feature-text">
            Trade with confidence using our secure payment and escrow services.
          </p>
        </div>

        <div className="feature-card">
          <div className="feature-icon">📊</div>
          <h3 className="feature-title">Market Insights</h3>
          <p className="feature-text">
            Access real-time market data and pricing trends for informed decisions.
          </p>
        </div>

        <div className="feature-card">
          <div className="feature-icon">💬</div>
          <h3 className="feature-title">Direct Communication</h3>
          <p className="feature-text">
            Connect directly with buyers and sellers through our messaging system.
          </p>
        </div>
      </div>
    </div>
  );
};

export default WelcomePage;