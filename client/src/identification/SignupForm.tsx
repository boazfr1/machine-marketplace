import React, { useState, FormEvent } from 'react';
import './SignupForm.css';
import api, { isAxiosError } from '../api';

interface SignupFormData {
  Name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

interface FormErrors {
  Name?: string;
  email?: string;
  password?: string;
  confirmPassword?: string;
  server?: string;
}

const SignupForm = () => {
  const [formData, setFormData] = useState<SignupFormData>({
    Name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });

  const [errors, setErrors] = useState<FormErrors>({});
  const [isLoading, setIsLoading] = useState(false);

  const validateForm = (): boolean => {
    const newErrors: FormErrors = {};
    
    if (!formData.Name.trim()) {
      newErrors.Name = 'First name is required';
    }
    
    if (!formData.email) {
      newErrors.email = 'Email is required';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'Please enter a valid email';
    }
    
    if (!formData.password) {
      newErrors.password = 'Password is required';
    } else if (formData.password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters';
    }

    if (!formData.confirmPassword) {
      newErrors.confirmPassword = 'Please confirm your password';
    } else if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    setErrors({});

    const signupData = {
      Name: formData.Name,
      email: formData.email,
      password: formData.password
    };

    try {
      await api.post('/api/v1/sign-up', signupData);      
    } catch (error) {
      if (isAxiosError(error)) {
        setErrors({
          server: error.response?.data?.message || 'An error occurred during sign up'
        });
      } else {
        setErrors({
          server: 'An unexpected error occurred'
        });
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    if (errors[name as keyof FormErrors]) {
      setErrors(prev => ({
        ...prev,
        [name]: undefined
      }));
    }
  };

  return (
    <div className="signup-container">
      <form onSubmit={handleSubmit} className="signup-form">
        <h2 className="form-title">Create Account</h2>
        
        <div className="name-group">
          <div className="form-group">
            <label htmlFor="Name">Name</label>
            <input
              id="Name"
              type="text"
              name="Name"
              value={formData.Name}
              onChange={handleChange}
              autoComplete="given-name"
            />
            {errors.Name && <div className="error-message">{errors.Name}</div>}
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            autoComplete="email"
          />
          {errors.email && <div className="error-message">{errors.email}</div>}
        </div>

        <div className="form-group">
          <label htmlFor="password">Password</label>
          <input
            id="password"
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            autoComplete="new-password"
          />
          {errors.password && <div className="error-message">{errors.password}</div>}
        </div>

        <div className="form-group">
          <label htmlFor="confirmPassword">Confirm Password</label>
          <input
            id="confirmPassword"
            type="password"
            name="confirmPassword"
            value={formData.confirmPassword}
            onChange={handleChange}
            autoComplete="new-password"
          />
          {errors.confirmPassword && <div className="error-message">{errors.confirmPassword}</div>}
        </div>

        {errors.server && <div className="form-error">{errors.server}</div>}

        <button 
          type="submit" 
          className="submit-button"
          disabled={isLoading}
        >
          {isLoading ? 'Creating Account...' : 'Sign Up'}
        </button>
      </form>
    </div>
  );
};

export default SignupForm;