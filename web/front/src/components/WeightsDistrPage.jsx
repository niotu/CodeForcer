import React, {useEffect, useState} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import './styles.css'; // Import the provided CSS file

const WeightsDistrPage = () => {
    let comment = 'the comment'
    const {groupCode, contestId} = useParams(); // Extracting groupCode and contestId from URL parameters
    const navigate = useNavigate();
    const [tasks, setTasks] = useState([]);
    const [taskWeights, setTaskWeights] = useState([]);

    function logout() {
        localStorage.setItem('isAuthorized', 'false');
        localStorage.setItem('userId', null);
    }

    const handleWeights = async (e) => {
        e.preventDefault();
        const weights = Array.from(e.target).reduce((acc, curr) => {
            if (curr.name.startsWith('task-')) {
                const taskId = parseInt(curr.name.split('-')[1]);
                const weight = parseInt(curr.value);
                acc[taskId] = weight;
            }
            return acc;
        }, {});

        const response = await fetch(`/api/setTaskWeights`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                groupCode,
                contestId,
                userID: localStorage.getItem('userId'),
                weights,
            }),
        });
    };

    useEffect(() => {
        const queryParams = new URLSearchParams({
            groupCode: groupCode,
            contestId: contestId,
            userID: localStorage.getItem('userId'),
        });

        const fetchTasks = async () => {
            const response = await fetch(`/api/getTasks?${queryParams}`);
            const data = await response.json();
            setTasks(data.result);
            console.log(data);
        };

        fetchTasks()


    }, [])
    return (

        <div className='page-active'>
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Set up tasks weights</h1>
                    </div>
                    <div className="right-part">
                        <form onSubmit={handleWeights} autoComplete='on'>
                            <nav className="list-view" id='distribution-form'>
                                <ul>
                                    {tasks.map(task => (
                                        <li key={task.Index}>
                                            <label className='task'> {task.Name}</label>
                                            <input type='number' onChange={(e) => {
                                            setTaskWeights(e.target.value, task.Index)}} min={0} max={task.MaxPoints}/>
                                        </li>
                                    ))}
                                </ul>
                            </nav>
                            <button type='submit'>Submit</button>
                        </form>
                    </div>
                </div>
            </div>
            <div className="navigation">
                <div className="left-navigation-part">
                    <a href="/link">
                        <button className="previous-page">Back</button>
                    </a>
                </div>
                <p>{comment}</p>
                <div className="right-navigation-part">
                    <a href="/">
                        <button className={'logout'} onClick={(e) => logout(e)}>Logout</button>
                    </a>
                </div>
            </div>
        </div>
    );
};

export default WeightsDistrPage;