import React, {useEffect, useState} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import './styles.css'; // Import the provided CSS file

const ContestPage = () => {
    const {groupCode} = useParams();
    const navigate = useNavigate();
    const [contests, setContests] = useState([]);

    const queryParams = new URLSearchParams({
            userID: localStorage.getItem('userId'),
            groupCode: groupCode
        });

    useEffect(() => {
        const fetchContests = async () => {
            const response = await fetch(`/api/getContests?${queryParams}`);
            const data = await response.json();
            setContests(data.result);
            console.log(data);
        };

        fetchContests();
    }, [groupCode]);

    const handleContestClick = (contestId) => {
        navigate(`/contest-details/${groupCode}/${contestId}`);
    };

    return (
        <body>
        <div className="page-active">
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
                                        <a className='link' href={contest.ContestLink}>{contest.Name} on codeforces</a>
                                    </li>
                                ))}
                            </ul>
                        </nav>
                    </div>
                </div>
            </div>
        </div>
        </body>
    );
};

export default ContestPage;
