import { ReactNode } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import './MarketplaceLayout.css';

interface MarketplaceLayoutProps {
    children: ReactNode;
}

const MarketplaceLayout = ({ children }: MarketplaceLayoutProps) => {
    const navigate = useNavigate();
    const location = useLocation();

    const navigationItems = [
        {
            id: 'marketplace',
            icon: '🏪',
            label: 'Marketplace',
            path: '/marketplace',
            description: 'Browse available machines'
        },
        {
            id: 'add-machine',
            icon: '➕',
            label: 'List Machine',
            path: '/add-machine',
            description: 'Add your machine to marketplace'
        },
        {
            id: 'my-machines',
            icon: '🖥️',
            label: 'My Machines',
            path: '/my-machines',
            description: 'Manage your listed machines'
        },
        {
            id: 'purchased-machines',
            icon: '💰',
            label: 'Purchased',
            path: '/purchased-machines',
            description: 'View your purchased machines'
        }
    ];

    const handleNavigation = (path: string) => {
        navigate(path);
    };

    const handleLogout = () => {
        // Add logout logic here
        navigate('/');
    };

    return (
        <div className="marketplace-layout">
            <aside className="marketplace-sidebar">
                <div className="sidebar-header">
                    <div className="sidebar-brand">
                        <span className="brand-logo">⚙️</span>
                        <span className="brand-name">MachineMart</span>
                    </div>
                </div>

                <nav className="sidebar-nav">
                    <ul className="nav-list">
                        {navigationItems.map((item) => (
                            <li key={item.id}>
                                <button
                                    onClick={() => handleNavigation(item.path)}
                                    className={`nav-item ${location.pathname === item.path ? 'nav-item-active' : ''}`}
                                    title={item.description}
                                >
                                    <span className="nav-icon">{item.icon}</span>
                                    <div className="nav-content">
                                        <span className="nav-label">{item.label}</span>
                                        <span className="nav-description">{item.description}</span>
                                    </div>
                                </button>
                            </li>
                        ))}
                    </ul>
                </nav>

                <div className="sidebar-footer">
                    <button onClick={handleLogout} className="logout-btn">
                        <span className="nav-icon">🚪</span>
                        <span className="nav-label">Sign out</span>
                    </button>
                </div>
            </aside>

            <main className="marketplace-main">
                <div className="main-content">
                    {children}
                </div>
            </main>
        </div>
    );
};

export default MarketplaceLayout;
