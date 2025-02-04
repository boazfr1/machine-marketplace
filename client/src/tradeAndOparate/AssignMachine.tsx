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
import api from '../api';
import './AssignMachine.css';
import Sidebar from '../dashboard/SideBar';

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
      await api.post('/api/v1/machine/create', formData);
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
          <TextField
            fullWidth
            label="Machine Name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            error={!!errors.name}
            helperText={errors.name}
            margin="normal"
          />

          <TextField
            fullWidth
            label="RAM (GB)"
            name="ram"
            type="number"
            value={formData.ram || ''}
            onChange={handleChange}
            error={!!errors.ram}
            helperText={errors.ram}
            margin="normal"
          />

          <TextField
            fullWidth
            label="CPU Cores"
            name="cpu"
            type="number"
            value={formData.cpu || ''}
            onChange={handleChange}
            error={!!errors.cpu}
            helperText={errors.cpu}
            margin="normal"
          />

          <TextField
            fullWidth
            label="Memory (GB)"
            name="memory"
            type="number"
            value={formData.memory || ''}
            onChange={handleChange}
            error={!!errors.memory}
            helperText={errors.memory}
            margin="normal"
          />

          <TextField
            fullWidth
            label="SSH Key"
            name="key"
            multiline
            rows={3}
            value={formData.key}
            onChange={handleChange}
            error={!!errors.key}
            helperText={errors.key}
            margin="normal"
          />

          <TextField
            fullWidth
            label="Host"
            name="host"
            value={formData.host}
            onChange={handleChange}
            error={!!errors.host}
            helperText={errors.host}
            margin="normal"
          />

          <TextField
            fullWidth
            label="SSH User"
            name="ssh_user"
            value={formData.ssh_user}
            onChange={handleChange}
            error={!!errors.ssh_user}
            helperText={errors.ssh_user}
            margin="normal"
          />

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