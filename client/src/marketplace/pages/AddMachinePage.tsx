import { useState, FormEvent } from "react";
import MarketplaceLayout from "../components/MarketplaceLayout";
import "./AddMachinePage.css";

interface MachineFormData {
  name: string;
  cpu: number;
  ram: number;
  storage: number;
  description: string;
  hourlyRate: number;
}

const AddMachinePage = () => {
  const [formData, setFormData] = useState<MachineFormData>({
    name: '',
    cpu: 1,
    ram: 1,
    storage: 10,
    description: '',
    hourlyRate: 0.25
  });
  
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      // Add API call to submit machine
      console.log("Adding machine:", formData);
      await new Promise(resolve => setTimeout(resolve, 2000)); // Simulate API call
      
      // Reset form or redirect
      alert("Machine added successfully!");
    } catch (error) {
      console.error("Failed to add machine:", error);
      alert("Failed to add machine. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'number' ? Number(value) : value
    }));
  };

  return (
    <MarketplaceLayout>
      <div className="add-machine-page">
        <div className="page-header">
          <h1 className="page-title">Add New Machine</h1>
          <p className="page-subtitle">
            List your machine on the marketplace and start earning
          </p>
        </div>

        <div className="form-container">
          <form onSubmit={handleSubmit} className="machine-form">
            <div className="form-section">
              <h2 className="section-title">Basic Information</h2>
              
              <div className="form-group">
                <label htmlFor="name" className="form-label">Machine Name</label>
                <input
                  id="name"
                  name="name"
                  type="text"
                  value={formData.name}
                  onChange={handleChange}
                  className="form-input"
                  placeholder="e.g., High-Performance ML Server"
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="description" className="form-label">Description</label>
                <textarea
                  id="description"
                  name="description"
                  value={formData.description}
                  onChange={handleChange}
                  className="form-textarea"
                  placeholder="Describe your machine's capabilities and ideal use cases..."
                  rows={4}
                  required
                />
              </div>
            </div>

            <div className="form-section">
              <h2 className="section-title">Technical Specifications</h2>
              
              <div className="specs-grid">
                <div className="form-group">
                  <label htmlFor="cpu" className="form-label">CPU Cores</label>
                  <input
                    id="cpu"
                    name="cpu"
                    type="number"
                    min="1"
                    max="128"
                    value={formData.cpu}
                    onChange={handleChange}
                    className="form-input"
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="ram" className="form-label">RAM (GB)</label>
                  <input
                    id="ram"
                    name="ram"
                    type="number"
                    min="1"
                    max="1024"
                    value={formData.ram}
                    onChange={handleChange}
                    className="form-input"
                    required
                  />
                </div>

                <div className="form-group">
                  <label htmlFor="storage" className="form-label">Storage (GB)</label>
                  <input
                    id="storage"
                    name="storage"
                    type="number"
                    min="10"
                    max="10000"
                    value={formData.storage}
                    onChange={handleChange}
                    className="form-input"
                    required
                  />
                </div>
              </div>
            </div>

            <div className="form-section">
              <h2 className="section-title">Pricing</h2>
              
              <div className="form-group">
                <label htmlFor="hourlyRate" className="form-label">Hourly Rate ($)</label>
                <input
                  id="hourlyRate"
                  name="hourlyRate"
                  type="number"
                  min="0.01"
                  max="100"
                  step="0.01"
                  value={formData.hourlyRate}
                  onChange={handleChange}
                  className="form-input"
                  required
                />
                <div className="form-help">
                  Recommended rate: $0.25 - $2.00 per hour based on specifications
                </div>
              </div>
            </div>

            <div className="form-actions">
              <button 
                type="button" 
                className="btn btn-secondary"
                disabled={isSubmitting}
              >
                Save as Draft
              </button>
              <button 
                type="submit" 
                className="btn btn-primary btn-lg"
                disabled={isSubmitting}
              >
                {isSubmitting ? (
                  <>
                    <span className="loading-spinner"></span>
                    Adding Machine...
                  </>
                ) : (
                  'List Machine'
                )}
              </button>
            </div>
          </form>

          <div className="preview-section">
            <h3 className="preview-title">Preview</h3>
            <div className="machine-preview">
              <div className="preview-header">
                <h4 className="preview-name">
                  {formData.name || 'Machine Name'}
                </h4>
                <span className="preview-rate">
                  ${formData.hourlyRate}/hour
                </span>
              </div>
              
              <div className="preview-specs">
                <div className="preview-spec">
                  <span className="spec-label">CPU:</span>
                  <span className="spec-value">{formData.cpu} cores</span>
                </div>
                <div className="preview-spec">
                  <span className="spec-label">RAM:</span>
                  <span className="spec-value">{formData.ram} GB</span>
                </div>
                <div className="preview-spec">
                  <span className="spec-label">Storage:</span>
                  <span className="spec-value">{formData.storage} GB</span>
                </div>
              </div>
              
              {formData.description && (
                <div className="preview-description">
                  {formData.description}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </MarketplaceLayout>
  );
};

export default AddMachinePage;
