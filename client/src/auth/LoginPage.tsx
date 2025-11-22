import React, { useState, FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { authApi, isAxiosError } from '../global/api';
import './AuthPage.css';

interface LoginFormData {
    email: string;
    password: string;
}

interface FormErrors {
    email?: string;
    password?: string;
    server?: string;
}

const LoginPage = () => {
    const [formData, setFormData] = useState<LoginFormData>({
        email: '',
        password: ''
    });
    const [errors, setErrors] = useState<FormErrors>({});
    const [isLoading, setIsLoading] = useState(false);

    const navigate = useNavigate();

    const validateForm = (): boolean => {
        const newErrors: FormErrors = {};

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

        try {
            await authApi.post('/api/v1/login', formData);
            navigate('/marketplace');
        } catch (error) {
            if (isAxiosError(error)) {
                setErrors({
                    server: error.response?.data?.message || 'Invalid email or password'
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
                    <h1 className="auth-title">Welcome back</h1>
                    <p className="auth-subtitle">Sign in to your account to continue</p>
                </div>

                <form onSubmit={handleSubmit} className="auth-form">
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
                            placeholder="Enter your password"
                            autoComplete="current-password"
                        />
                        {errors.password && <div className="form-error">{errors.password}</div>}
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
                                Signing in...
                            </>
                        ) : (
                            'Sign in'
                        )}
                    </button>
                </form>

                <div className="auth-footer">
                    <p>
                        Don't have an account?{' '}
                        <Link to="/signup" className="auth-link">
                            Sign up
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    );
};

export default LoginPage;
