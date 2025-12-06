import { useNavigate } from 'react-router-dom';
import './LandingPage.css';

const LandingPage = () => {
    const navigate = useNavigate();

    const handleLogin = () => {
        navigate('/login');
    };

    const handleSignup = () => {
        navigate('/signup');
    };

    const features = [
        {
            icon: '🔍',
            title: 'Smart Search',
            description: 'Find exactly what you need with our advanced search and intelligent filtering system.'
        },
        {
            icon: '🛡️',
            title: 'Secure Trading',
            description: 'Trade with confidence using our secure payment system and verified seller network.'
        },
        {
            icon: '📊',
            title: 'Market Analytics',
            description: 'Access real-time market data, pricing trends, and performance insights for informed decisions.'
        },
        {
            icon: '💬',
            title: 'Direct Communication',
            description: 'Connect instantly with buyers and sellers through our integrated messaging platform.'
        },
        {
            icon: '⚡',
            title: 'Instant Deployment',
            description: 'Deploy and access your machines instantly with our cloud-based infrastructure.'
        },
        {
            icon: '🔧',
            title: 'Easy Management',
            description: 'Manage your entire machine portfolio from a single, intuitive dashboard.'
        }
    ];

    return (
        <div className="landing-page">
            <header className="landing-header">
                <nav className="landing-nav">
                    <div className="nav-brand">
                        <span className="brand-logo">⚙️</span>
                        <span className="brand-name">MachineMart</span>
                    </div>
                    <div className="nav-actions">
                        <button onClick={handleLogin} className="btn btn-outline">
                            Sign in
                        </button>
                        <button onClick={handleSignup} className="btn btn-primary">
                            Get started
                        </button>
                    </div>
                </nav>
            </header>

            <main className="landing-main">
                <section className="hero-section">
                    <div className="hero-content">
                        <h1 className="hero-title">
                            The Future of
                            <span className="hero-highlight"> Machine Trading</span>
                        </h1>
                        <p className="hero-description">
                            Connect with a global network of machine providers and users.
                            Rent, buy, or sell computing power with unprecedented ease and security.
                        </p>
                        <div className="hero-actions">
                            <button onClick={handleSignup} className="btn btn-primary btn-lg">
                                Start Trading
                                <span className="btn-icon">→</span>
                            </button>
                            <button className="btn btn-secondary btn-lg">
                                Learn More
                            </button>
                        </div>
                    </div>
                    <div className="hero-visual">
                        <div className="hero-card">
                            <div className="hero-card-header">
                                <div className="hero-card-title">High-Performance Server</div>
                                <div className="hero-card-status">Available</div>
                            </div>
                            <div className="hero-card-specs">
                                <div className="spec-item">
                                    <span className="spec-label">CPU</span>
                                    <span className="spec-value">32 cores</span>
                                </div>
                                <div className="spec-item">
                                    <span className="spec-label">RAM</span>
                                    <span className="spec-value">128 GB</span>
                                </div>
                                <div className="spec-item">
                                    <span className="spec-label">Storage</span>
                                    <span className="spec-value">2 TB SSD</span>
                                </div>
                            </div>
                            <div className="hero-card-price">$0.50/hour</div>
                        </div>
                    </div>
                </section>

                <section className="features-section">
                    <div className="features-header">
                        <h2 className="features-title">Everything you need to succeed</h2>
                        <p className="features-description">
                            Powerful tools and features designed to make machine trading simple, secure, and profitable.
                        </p>
                    </div>
                    <div className="features-grid">
                        {features.map((feature, index) => (
                            <div key={index} className="feature-card">
                                <div className="feature-icon">{feature.icon}</div>
                                <h3 className="feature-title">{feature.title}</h3>
                                <p className="feature-description">{feature.description}</p>
                            </div>
                        ))}
                    </div>
                </section>

                <section className="cta-section">
                    <div className="cta-content">
                        <h2 className="cta-title">Ready to get started?</h2>
                        <p className="cta-description">
                            Join thousands of users who are already trading machines on MachineMart.
                        </p>
                        <button onClick={handleSignup} className="btn btn-primary btn-lg">
                            Create your account
                        </button>
                    </div>
                </section>
            </main>

            <footer className="landing-footer">
                <div className="footer-content">
                    <div className="footer-brand">
                        <span className="brand-logo">⚙️</span>
                        <span className="brand-name">MachineMart</span>
                    </div>
                    <p className="footer-text">
                        The world's leading marketplace for computing resources.
                    </p>
                </div>
            </footer>
        </div>
    );
};

export default LandingPage;
