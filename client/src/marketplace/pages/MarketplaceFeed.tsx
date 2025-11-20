import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import MachineCard from "../components/MachineCard";
import MarketplaceLayout from "../components/MarketplaceLayout";
import api from "../../global/api";
import { MachineType } from "../../global/types";
import "./MarketplaceFeed.css";

const MarketplaceFeed = () => {
    const [availableMachines, setAvailableMachines] = useState<MachineType[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");
    const [sortBy, setSortBy] = useState<'name' | 'cpu' | 'ram' | 'memory'>('name');
    const [filterBy, setFilterBy] = useState<'all' | 'high-cpu' | 'high-ram' | 'high-storage'>('all');

    const navigate = useNavigate();

    const getAllAvailableMachines = async () => {
        try {
            const { data } = await api<MachineType[]>('/api/v1/machine');
            console.log("Available machines:", data);
            setAvailableMachines(data);
        } catch (error) {
            console.error("Failed to fetch machines:", error);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        getAllAvailableMachines();
    }, []);

    const navigateToMachineDetails = (machine: MachineType) => {
        navigate('/machine/details', {
            state: { machine }
        });
    };

    const filteredAndSortedMachines = availableMachines
        .filter(machine => {
            // Search filter
            const matchesSearch = machine.Name.toLowerCase().includes(searchQuery.toLowerCase());

            // Category filter
            let matchesFilter = true;
            switch (filterBy) {
                case 'high-cpu':
                    matchesFilter = machine.Cpu >= 8;
                    break;
                case 'high-ram':
                    matchesFilter = machine.Ram >= 16;
                    break;
                case 'high-storage':
                    matchesFilter = machine.Memory >= 500;
                    break;
                default:
                    matchesFilter = true;
            }

            return matchesSearch && matchesFilter;
        })
        .sort((a, b) => {
            switch (sortBy) {
                case 'cpu':
                    return b.Cpu - a.Cpu;
                case 'ram':
                    return b.Ram - a.Ram;
                case 'memory':
                    return b.Memory - a.Memory;
                default:
                    return a.Name.localeCompare(b.Name);
            }
        });

    return (
        <MarketplaceLayout>
            <div className="marketplace-feed">
                <div className="feed-header">
                    <div className="header-content">
                        <h1 className="feed-title">Machine Marketplace</h1>
                        <p className="feed-subtitle">
                            Discover and rent high-performance computing resources from verified providers
                        </p>
                    </div>

                    <div className="feed-stats">
                        <div className="stat-item">
                            <span className="stat-value">{availableMachines.length}</span>
                            <span className="stat-label">Available Machines</span>
                        </div>
                        <div className="stat-item">
                            <span className="stat-value">{filteredAndSortedMachines.length}</span>
                            <span className="stat-label">Matching Results</span>
                        </div>
                    </div>
                </div>

                <div className="feed-controls">
                    <div className="search-section">
                        <div className="search-input-wrapper">
                            <span className="search-icon">🔍</span>
                            <input
                                type="text"
                                placeholder="Search machines by name..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                className="search-input"
                            />
                        </div>
                    </div>

                    <div className="filter-section">
                        <div className="filter-group">
                            <label className="filter-label">Filter by:</label>
                            <select
                                value={filterBy}
                                onChange={(e) => setFilterBy(e.target.value as any)}
                                className="filter-select"
                            >
                                <option value="all">All Machines</option>
                                <option value="high-cpu">High CPU (8+ cores)</option>
                                <option value="high-ram">High RAM (16+ GB)</option>
                                <option value="high-storage">High Storage (500+ GB)</option>
                            </select>
                        </div>

                        <div className="filter-group">
                            <label className="filter-label">Sort by:</label>
                            <select
                                value={sortBy}
                                onChange={(e) => setSortBy(e.target.value as any)}
                                className="filter-select"
                            >
                                <option value="name">Name</option>
                                <option value="cpu">CPU Cores</option>
                                <option value="ram">RAM</option>
                                <option value="memory">Storage</option>
                            </select>
                        </div>
                    </div>
                </div>

                <div className="feed-content">
                    {isLoading ? (
                        <div className="loading-state">
                            <div className="loading-spinner-large"></div>
                            <p>Loading available machines...</p>
                        </div>
                    ) : filteredAndSortedMachines.length > 0 ? (
                        <div className="machines-grid">
                            {filteredAndSortedMachines.map((machine, index) => (
                                <MachineCard
                                    key={`${machine.Name}-${machine.OwnerID}-${index}`}
                                    Name={machine.Name}
                                    Ram={machine.Ram}
                                    Cpu={machine.Cpu}
                                    Memory={machine.Memory}
                                    OwnerID={machine.OwnerID}
                                    onClick={() => navigateToMachineDetails(machine)}
                                />
                            ))}
                        </div>
                    ) : (
                        <div className="empty-state">
                            <div className="empty-icon">🔍</div>
                            <h3 className="empty-title">No machines found</h3>
                            <p className="empty-description">
                                {searchQuery || filterBy !== 'all'
                                    ? "Try adjusting your search or filter criteria"
                                    : "There are no machines available at the moment"
                                }
                            </p>
                            {(searchQuery || filterBy !== 'all') && (
                                <button
                                    onClick={() => {
                                        setSearchQuery("");
                                        setFilterBy('all');
                                    }}
                                    className="btn btn-outline"
                                >
                                    Clear filters
                                </button>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </MarketplaceLayout>
    );
};

export default MarketplaceFeed;
