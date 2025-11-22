import React, { useState, FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authApi, isAxiosError } from '../global/api';
import './AuthPage.css';

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

const SignupPage = () => {
    const [formData, setFormData] = useState<SignupFormData>({
        Name: '',
        email: '',
        password: '',
        confirmPassword: ''
    });

    const [errors, setErrors] = useState<FormErrors>({});
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();

    const validateForm = (): boolean => {
        const newErrors: FormErrors = {};

        if (!formData.Name.trim()) {
            newErrors.Name = 'Name is required';
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
            await authApi.post('/api/v1/sign-up', signupData);
            navigate('/marketplace');
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
        <div className="auth-page">
            <div className="auth-container">
                <div className="auth-header">
                    <Link to="/" className="auth-logo">
                        MachineMart
                    </Link>
                    <h1 className="auth-title">Create your account</h1>
                    <p className="auth-subtitle">Join MachineMart to start trading machines</p>
                </div>

                <form onSubmit={handleSubmit} className="auth-form">
                    <div className="form-group">
                        <label htmlFor="Name" className="form-label">Full name</label>
                        <input
                            id="Name"
                            type="text"
                            name="Name"
                            value={formData.Name}
                            onChange={handleChange}
                            className={`form-input ${errors.Name ? 'form-input-error' : ''}`}
                            placeholder="Enter your full name"
                            autoComplete="name"
                        />
                        {errors.Name && <div className="form-error">{errors.Name}</div>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="email" className="form-label">Email address</label>
                        <input
                            id="email"
                            type="email"
                            name="email"
                            value={formData.email}
                            onChange={handleChange}
                            className={`form-input ${errors.email ? 'form-input-error' : ''}`}
                            placeholder="Enter your email"
                            autoComplete="email"
                        />
                        {errors.email && <div className="form-error">{errors.email}</div>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="password" className="form-label">Password</label>
                        <input
                            id="password"
                            type="password"
                            name="password"
                            value={formData.password}
                            onChange={handleChange}
                            className={`form-input ${errors.password ? 'form-input-error' : ''}`}
                            placeholder="Create a password"
                            autoComplete="new-password"
                        />
                        {errors.password && <div className="form-error">{errors.password}</div>}
                    </div>

                    <div className="form-group">
                        <label htmlFor="confirmPassword" className="form-label">Confirm password</label>
                        <input
                            id="confirmPassword"
                            type="password"
                            name="confirmPassword"
                            value={formData.confirmPassword}
                            onChange={handleChange}
                            className={`form-input ${errors.confirmPassword ? 'form-input-error' : ''}`}
                            placeholder="Confirm your password"
                            autoComplete="new-password"
                        />
                        {errors.confirmPassword && <div className="form-error">{errors.confirmPassword}</div>}
                    </div>

                    {errors.server && (
                        <div className="server-error">
                            {errors.server}
                        </div>
                    )}

                    <button
                        type="submit"
                        className="btn btn-primary btn-lg w-full"
                        disabled={isLoading}
                    >
                        {isLoading ? (
                            <>
                                <span className="loading-spinner"></span>
                                Creating account...
                            </>
                        ) : (
                            'Create account'
                        )}
                    </button>
                </form>

                <div className="auth-footer">
                    <p>
                        Already have an account?{' '}
                        <Link to="/login" className="auth-link">
                            Sign in
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default SignupPage;
