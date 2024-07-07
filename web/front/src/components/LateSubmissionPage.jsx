import React, {useState} from 'react';
import {useNavigate, useParams} from "react-router-dom"; // Import the provided CSS file
import './styles.css';
import logout from "./globalFunctions.js";

const LateSubmissionPage = () => {
    // ... other state variables ...
    const {groupCode, contestId} = useParams(); // Extracting groupCode and contestId from URL parameters
    const [date, setDate] = useState(new Date()); // Initialize with current date
    const navigate = useNavigate();
    const [comment, setComment] = useState('');
    const [isCorrect, setIsCorrect] = useState(true);

    const lateSubmit = async (e) => {
        e.preventDefault();
        console.log(' processing...')

        // ... (Your logic to get groupCode, contestId, and userID) ...

        // Convert the date to a string in ISO format for sending to the API
        const formattedDate = date.toISOString();
        console.log(formattedDate);
        sessionStorage.setItem('date', formattedDate);

        navigate(`/contest-details/${groupCode}/${contestId}`);

    }

    // ... (Rest of your component logic) ...

    return (
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Set up late submission date</h1>
                    </div>
                    <div className="right-part">
                        <form onSubmit={lateSubmit} autoComplete='on'>

                            <label htmlFor="date">Choose the date:</label>
                            <input
                                type="datetime-local"
                                id="date"
                                value={date.toISOString().slice(0, 16)} // Format for datetime-local input
                                onChange={(e) => setDate(new Date(e.target.value))}
                                required
                                className={isCorrect ? 'correct' : 'incorrect'}
                            /><br/><br/>
                            <button type="submit">Submit</button>
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

export default LateSubmissionPage;
