import { useNavigate } from 'react-router-dom';
import './SideBar.css';
import { useEffect } from 'react';

const Sidebar = () => {


  const navigate = useNavigate();

  useEffect(() => {
    console.log("MyMachinesPage");
  }, []);


  const menuItems = [
    { icon: '📊', label: 'Feed', id: 'feed' },
    { icon: '🔧', label: 'Add Machine', id: 'assign' },
    { icon: '💻', label: 'Owned Machines', id: 'my-machines' },
    { icon: '🖥️', label: 'Bought Machines', id: 'assigned' }
  ];

  const eventHandler = (menuID: number) => {
    switch (menuID) {
      case 0:
        navigate('/feed');
        break;
      case 1:
        navigate('/assign-machine');
        break;
      case 2:        
        navigate('/my-machines');
        break;
      case 3:        
        navigate('/assigned-machines');
        break;
      default:
        navigate('/')
    }
  }

  return (
    <div className="sidebar">
      <div className="sidebar-header">
        <h1>Dashboard</h1>
      </div>
      <nav>
        <ul className="menu-list">
          {menuItems.map((item, index) => (
            <li key={item.id} onClick={() => eventHandler(index)}>
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