import './SideBar.css';

const Sidebar = () => {
  const menuItems = [
    { icon: '📊', label: 'Feed', id: 'feed' },
    { icon: '🔧', label: 'Assign Machine', id: 'assign' },
    { icon: '💻', label: 'My Machines', id: 'my-machines' },
    { icon: '🖥️', label: 'Assigned Machines', id: 'assigned' }
  ];

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <h1>Dashboard</h1>
      </div>
      <nav>
        <ul className="menu-list">
          {menuItems.map((item) => (
            <li key={item.id}>
              <button className="menu-button">
                <span className="icon">{item.icon}</span>
                <span className="label">{item.label}</span>
              </button>
            </li>
          ))}
        </ul>
      </nav>
    </div>
  );
};

export default Sidebar;