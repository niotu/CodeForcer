import React, { useEffect, useState } from 'react';
import './styles.css'; // Import the provided CSS file

const GroupsPage = () => {
    const [groups, setGroups] = useState([]);
    const [error, setError] = useState(null); // State for error message

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch('/api/getGroups');
                if (!response.ok) {
                    throw new Error('Failed to fetch groups');
                }
                const data = await response.json();
                setGroups(data);
            } catch (error) {
                setError(error.message);
            }
        };

        fetchGroups();
    }, []);

    return (
        <div className="wizard">
            <div className="panel">
                <div className="left-part">
                    <h1>Select Group</h1>
                </div>
                <div className="right-part">
                    <nav className="list-view">
                        <ul>
                            {groups.map(group => (
                                <li key={group.GroupCode}><a href={`/contests/${group.GroupCode}`}>{group.GroupName}</a></li>
                            ))}
                        </ul>
                    </nav>
                </div>
            </div>
            {error && (
                <div className="error-box">
                    <p>{error}</p>
                </div>
            )}
        </div>
    );
};

export default GroupsPage;
