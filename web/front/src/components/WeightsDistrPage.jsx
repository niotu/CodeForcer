import React, {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import './styles.css'; // Import the provided CSS file

const WeightsDistrPage = () => {

    const navigate = useNavigate();
    const [tasks, setTasks] = useState([]);

    useEffect(() => {
        const fetchTasks = async () => {
            const response = await fetch('/api/getGroups');
            const data = await response.json();
            setTasks(data);
            console.log(data);
        };

        fetchTasks()
    }, [])

    return (
        <body>
        <div className='page-active'>
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Set up tasks weights</h1>
                        <h4>Note: summary points distribution must be 15.</h4>
                    </div>
                    <div className="right-part">
                        <nav className="distribution-form">
                            <ul>
                                {tasks.map(task => (
                                    <li key={task.taskCode}><p>{task.taskName }</p>
                                        <label>
                                            <input type="number" name="task1" value="1"/>
                                        </label>
                                    </li>
                                ))}
                                <li>

                                    <p>task1</p>
                                    <label>
                                        <input type="number" name="task1" value="1"/>
                                    </label>
                                </li>
                            </ul>
                        </nav>
                    </div>
                </div>
            </div>
            <button className={'logout'} onSubmit={(e) => logout()}>Logout</button>
        </div>
        </body>
    );
};

export default WeightsDistrPage;