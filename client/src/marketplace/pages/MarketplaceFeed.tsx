import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import MachineCard from "../components/MachineCard";
import MarketplaceLayout from "../components/MarketplaceLayout";
import { orderApi } from "../../global/api";
import { MachineType } from "../../global/types";
import "./MarketplaceFeed.css";

const MarketplaceFeed = () => {
    const [availableMachines, setAvailableMachines] = useState<MachineType[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [searchQuery, setSearchQuery] = useState("");
    const [sortBy, setSortBy] = useState<'name' | 'cpu' | 'ram' | 'memory'>('name');
    const [filterBy, setFilterBy] = useState<'all' | 'high-cpu' | 'high-ram' | 'high-storage' | 'gpu'>('all');
    const [minCpu, setMinCpu] = useState<number>(0);
    const [minRam, setMinRam] = useState<number>(0);
    const [minGpu, setMinGpu] = useState<number>(0);

    const navigate = useNavigate();

    const getAllAvailableMachines = async () => {
        try {
            const params = new URLSearchParams();
            if (minCpu > 0) params.append('cpu', minCpu.toString());
            if (minRam > 0) params.append('ram', minRam.toString());
            if (minGpu > 0) params.append('gpu', minGpu.toString());

            const url = `/api/v1/order${params.toString() ? `?${params.toString()}` : ''}`;
            const { data } = await orderApi.get<MachineType[]>(url);
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
    }, [minCpu, minRam, minGpu]);

    const navigateToMachineDetails = (machine: MachineType) => {
        navigate(`/machine/${machine.ID}`, {
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
                case 'gpu':
                    matchesFilter = machine.Gpu > 0;
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
                            <label className="filter-label">Min CPU:</label>
                            <input
                                type="range"
                                min="0"
                                max="32"
                                value={minCpu}
                                onChange={(e) => setMinCpu(Number(e.target.value))}
                                className="filter-range"
                            />
                            <span className="filter-value">{minCpu > 0 ? `${minCpu}+` : 'Any'}</span>
                        </div>

                        <div className="filter-group">
                            <label className="filter-label">Min RAM (GB):</label>
                            <input
                                type="range"
                                min="0"
                                max="128"
                                step="4"
                                value={minRam}
                                onChange={(e) => setMinRam(Number(e.target.value))}
                                className="filter-range"
                            />
                            <span className="filter-value">{minRam > 0 ? `${minRam}+` : 'Any'}</span>
                        </div>

                        <div className="filter-group">
                            <label className="filter-label">Min GPU:</label>
                            <input
                                type="range"
                                min="0"
                                max="8"
                                value={minGpu}
                                onChange={(e) => setMinGpu(Number(e.target.value))}
                                className="filter-range"
                            />
                            <span className="filter-value">{minGpu > 0 ? `${minGpu}+` : 'Any'}</span>
                        </div>

                        <div className="filter-group">
                            <label className="filter-label">Category:</label>
                            <select
                                value={filterBy}
                                onChange={(e) => setFilterBy(e.target.value as any)}
                                className="filter-select"
                            >
                                <option value="all">All Machines</option>
                                <option value="high-cpu">High CPU (8+ cores)</option>
                                <option value="high-ram">High RAM (16+ GB)</option>
                                <option value="gpu">With GPU</option>
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
                            {filteredAndSortedMachines.map((machine) => (
                                <MachineCard
                                    key={machine.ID}
                                    Name={machine.Name}
                                    Ram={machine.Ram}
                                    Cpu={machine.Cpu}
                                    Gpu={machine.Gpu}
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
                            {(searchQuery || filterBy !== 'all' || minCpu > 0 || minRam > 0 || minGpu > 0) && (
                                <button
                                    onClick={() => {
                                        setSearchQuery("");
                                        setFilterBy('all');
                                        setMinCpu(0);
                                        setMinRam(0);
                                        setMinGpu(0);
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
