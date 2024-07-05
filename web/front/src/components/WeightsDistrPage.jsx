import React, {useEffect, useState} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import './styles.css'; // Import the provided CSS file

const WeightsDistrPage = () => {
    let comment = 'The task weights must be in percentage (0-100%)'
    const {groupCode, contestId} = useParams(); // Extracting groupCode and contestId from URL parameters
    const navigate = useNavigate();
    const [tasks, setTasks] = useState([]);
    const [weights, setWeights] = useState([]);

    function logout() {
        localStorage.setItem('isAuthorized', 'false');
        localStorage.setItem('userId', null);
    }

    const handleWeights = async (e) => {
        e.preventDefault();

        console.log(`Weights: [ ${weights.join('-')} ]`)

        localStorage.setItem('weights', weights);

        navigate(`/contest-details/${groupCode}/${contestId}`);
    };

    const setTaskWeights = (value, index) => {
        console.log(`Setting weight for ${index} to ${value}`);
        weights[index] = value;
        setWeights([...weights]);
    }

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
                                    {
                                        tasks.map((task, index) => ( // Add the 'index' parameter here
                                            <li key={task.Index}>
                                                <label className='task'> {task.Name}</label>
                                                <input
                                                    type='number'
                                                    onChange={(e) => {
                                                        setTaskWeights(e.target.value, index); // Pass the index
                                                    }}
                                                    min={0}
                                                    max={task.MaxPoints}
                                                />
                                            </li>
                                        ))
                                    }
                                </ul>
                            </nav>

                            <button type='submit'>Submit</button>
                        </form>
                    </div>
                </div>
            </div>
            <div className="navigation">
                <div className="left-navigation-part">
                    <a href="">
                        <button className="previous-page" onClick={(e) => {
                            e.preventDefault();
                            history.go(-1);
                        }}>Back
                        </button>
                    </a>
                </div>
                <p>{comment}</p>
                <div className="right-navigation-part">
                    <a href="/">
                        <button className={'logout'} onClick={(e) => {localStorage.clear()}}>Logout</button>
                    </a>
                </div>
            </div>
        </div>
    );
};

export default WeightsDistrPage;