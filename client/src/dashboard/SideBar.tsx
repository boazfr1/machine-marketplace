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
    { icon: '🔧', label: 'Assign Machine', id: 'assign' },
    { icon: '💻', label: 'My Machines', id: 'my-machines' },
    { icon: '🖥️', label: 'Assigned Machines', id: 'assigned' }
  ];

  const eventHandler = (menuID: number) => {
    switch (menuID) {
      case 0:
        console.log("/feed = ", menuID);

        navigate('/feed');
        break;
      case 1:
        console.log("/assign-machine = ", menuID);

        navigate('/assign-machine');
        break;
      case 2:
        console.log("my-machines = ", menuID);
        
        navigate('/my-machines');
        break;
      case 4:
        console.log("/assigned-machines ", menuID);
        
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