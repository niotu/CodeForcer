import React, {useEffect, useState} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";


const WeightsDistrPage = () => {
    let comment = 'The task weights minimum value is 0'
    const {groupCode, contestId} = useParams(); // Extracting groupCode and contestId from URL parameters
    const navigate = useNavigate();

    const [tasks, setTasks] = useState([]);

    const [weights, setWeights] = useState([]);
    const [mode, setMode] = useState('best');

    if (!localStorage.getItem('isAuthorized')) {
        return show404page();
    }

    const handleWeights = async (e) => {
        e.preventDefault();

        console.log(`Weights: [ ${weights.join('-')} ]`)

        sessionStorage.setItem('weights', weights);
        sessionStorage.setItem('mode', mode)

        navigate(`/late-submissions/${groupCode}/${contestId}`);
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

            let url = process.env.REACT_APP_BACKEND_URL +
                '/api/getTasks?' + queryParams;

            const response = await fetch(url);
            const data = await response.json();
            setTasks(data.result.Problems);
            console.log(data);
        };

        fetchTasks()


    }, [])
    // const setMode = (value) => {
    //     mode = useState(value);
    // };

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
                                                    required
                                                    min={0}
                                                />
                                            </li>
                                        ))
                                    }
                                    <label className='task'>Mode: </label>
                                    <select id='mode'
                                            onChange={(e) =>
                                                setMode(e.target.value)}
                                            defaultValue={mode}>
                                        <option value='last'>Last</option>
                                        <option value='best'>Best</option>
                                    </select>
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
                        <button className={'logout'} onClick={() => logout()}>Logout
                        </button>
                    </a>
                </div>
            </div>
        </div>
    );
};

export default WeightsDistrPage;