import React, {useEffect, useState} from 'react';
import {useNavigate} from "react-router-dom"; // Import the provided CSS file
import './styles.css';

const GroupsPage = () => {
    const [groups, setGroups] = useState([]);
    // const navigate = useNavigate();

    useEffect(() => {
        const queryParams = new URLSearchParams({
            userID: localStorage.getItem('userId'),
        });
        console.log(`** userID for groups ${localStorage.getItem('userId')}`);
        const fetchGroups = async () => {
            const response = await fetch(`/api/getGroups?${queryParams}`);
            const data = await response.json();
            // for (const group of data.result) {
            //     console.log(group.AccessLevel);
            // }
            setGroups(data.result.filter(g => (g.AccessLevel === 'Manager' || g.AccessLevel === 'Менеджер')));
            console.log(data);
        };

        fetchGroups();
    }, []);

    return (
        <body>
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
                                    )
                                )
                                }
                            </ul>
                        </nav>
                    </div>
                </div>
            </div>
        </div>
        </body>
    );
};

export default GroupsPage;
