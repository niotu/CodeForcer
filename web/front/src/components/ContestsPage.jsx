import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './styles.css'; // Import the provided CSS file

const ContestPage = () => {
    const { groupCode } = useParams();
    const navigate = useNavigate();
    const [contests, setContests] = useState([]);

    useEffect(() => {
        const fetchContests = async () => {
            const response = await fetch(`/api/getContests?groupCode=${groupCode}`);
            const data = await response.json();
            setContests(data);
            console.log(data);
        };

        fetchContests();
    }, [groupCode]);

    const handleContestClick = (contestId) => {
        navigate(`/contest-details/${groupCode}/${contestId}`);
    };

    return (
        <div className="wizard">
            <div className="panel">
                <div className="left-part">
                    <h1>Select Contest</h1>
                </div>
                <div className="right-part">
                    <nav className="list-view">
                        <ul>
                            {contests.map(contest => (
                                <li key={contest.Id}>
                                    <a href="#" onClick={() => handleContestClick(contest.Id)}>
                                        {contest.Name}
                                    </a><br/>
                                    <a href={contest.ContestLink}>{contest.Name} on codeforces</a>
                                </li>
                            ))}
                        </ul>
                    </nav>
                </div>
            </div>
            <div className="navigation">
                <div className="left-navigation-part">
                    <button onClick={() => navigate(-1)} className="previous-page">Previous Page</button>
                </div>
                <div className="right-navigation-part">
                    <button onClick={() => navigate('/assignment')} className="next-page">Next Page</button>
                </div>
            </div>
        </div>
    );
};

export default ContestPage;
