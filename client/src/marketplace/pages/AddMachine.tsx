// CreateMachine.tsx
import { useState } from 'react';
import {
  TextField,
  Button,
  Paper,
  Typography,
  Container,
  Alert,
  CircularProgress
} from '@mui/material';
import { orderApi } from '../../global/api';
import '../style/AddMachine.css';
import Sidebar from '../../dashboard/SideBar';

interface MachineFormData {
  name: string;
  ram: number;
  cpu: number;
  memory: number;
  key: string;
  host: string;
  ssh_user: string;
}

interface FormErrors {
  [key: string]: string;
}

const AssignMachine = () => {
  const [formData, setFormData] = useState<MachineFormData>({
    name: '',
    ram: 0,
    cpu: 0,
    memory: 0,
    key: '',
    host: '',
    ssh_user: ''
  });

  const formConfig = [
    { name: 'name', label: 'Machine Name', type: 'text' },
    { name: 'ram', label: 'RAM (GB)', type: 'number' },
    { name: 'cpu', label: 'CPU Cores', type: 'number' },
    { name: 'memory', label: 'Memory (GB)', type: 'number' },
    { name: 'key', label: 'SSH Key', type: 'text', multiline: true, rows: 3 },
    { name: 'host', label: 'Host', type: 'text' },
    { name: 'ssh_user', label: 'SSH User', type: 'text' }
  ];

  const [errors, setErrors] = useState<FormErrors>({});
  const [isLoading, setIsLoading] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');

  const validateForm = () => {
    const newErrors: FormErrors = {};

    if (!formData.name) newErrors.name = 'Name is required';
    if (!formData.ram || formData.ram <= 0) newErrors.ram = 'RAM must be greater than 0';
    if (!formData.cpu || formData.cpu <= 0) newErrors.cpu = 'CPU must be greater than 0';
    if (!formData.memory || formData.memory <= 0) newErrors.memory = 'Memory must be greater than 0';
    if (!formData.key) newErrors.key = 'SSH key is required';
    if (!formData.host) newErrors.host = 'Host is required';
    if (!formData.ssh_user) newErrors.ssh_user = 'SSH user is required';

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setIsLoading(true);
    setSuccessMessage('');
    setErrors({});

    try {
      await orderApi.post('/api/v1/order/create', formData);
      setSuccessMessage('Machine created successfully!');
      setFormData({
        name: '',
        ram: 0,
        cpu: 0,
        memory: 0,
        key: '',
        host: '',
        ssh_user: ''
      });
    } catch (error: any) {
      setErrors({
        server: error.response?.data?.message || 'Failed to create machine'
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'ram' || name === 'cpu' || name === 'memory'
        ? parseInt(value) || 0
        : value
    }));
  };

  return (
    <Container maxWidth="sm" className="create-machine-container">
      <Sidebar />
      <Paper elevation={3} className="form-paper">
        <Typography variant="h5" component="h1" gutterBottom>
          Create New Machine
        </Typography>

        {successMessage && (
          <Alert severity="success" className="alert-message">
            {successMessage}
          </Alert>
        )}

        {errors.server && (
          <Alert severity="error" className="alert-message">
            {errors.server}
          </Alert>
        )}

        <form onSubmit={handleSubmit}>
          {formConfig.map((field: any) => (
            <TextField
              key={field.name}
              fullWidth
              label={field.label}
              name={field.name}
              type={field.type === 'number' ? 'number' : undefined}
              multiline={field.multiline}
              rows={field.rows}
              value={field.type === 'number'
                ? (formData[field.name as keyof MachineFormData] as number) || ''
                : formData[field.name as keyof MachineFormData]
              }
              onChange={handleChange}
              error={!!errors[field.name]}
              helperText={errors[field.name]}
              margin="normal"
            />
          ))}

          <Button
            type="submit"
            variant="contained"
            color="primary"
            fullWidth
            disabled={isLoading}
            className="submit-button"
          >
            {isLoading ? <CircularProgress size={24} /> : 'Create Machine'}
          </Button>
        </form>
      </Paper>
    </Container>
  );
};

export default AssignMachine;