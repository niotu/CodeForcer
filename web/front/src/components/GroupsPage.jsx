import React, {useEffect, useState} from 'react';
import './styles.css'; // Import the provided CSS file

const GroupsPage = () => {
    const [groups, setGroups] = useState([]);

    useEffect(() => {
        const fetchGroups = async () => {
            const response = await fetch('/api/getGroups');
            const data = await response.json();
            setGroups(data);
            console.log(data);
        };

        fetchGroups();
    }, []);

    return (

        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Select Group</h1>
                    </div>
                    <div className="right-part">
                        <nav className="list-view">
                            <ul>
                                {groups.map(group => (
                                    <li key={group.GroupCode}><a
                                        href={`/contests/${group.GroupCode}`}>{group.GroupName}</a></li>
                                ))}
                            </ul>
                        </nav>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default GroupsPage;
