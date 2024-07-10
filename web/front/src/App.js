import React from 'react';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import LoginPage from './components/LoginPage.jsx';
import ContestDetails from './components/ContestDetails.jsx';
import WeightsDistrPage from "./components/WeightsDistrPage.jsx";
import LinkPage from "./components/LinkPage.jsx";
import LateSubmissionPage from "./components/LateSubmissionPage.jsx";
import UploadCsvFilePage from "./components/UploadCsvFilePage.jsx";


const App = () => {
    return (
        <Router>
            <Routes>
                {/* Default Route for Unauthorized Users */}
                <Route path="/" element={<LoginPage/>}/>

                {/* Protected Routes (with Redirection) */}
                <Route path="/link"
                       element={
                           // localStorage.getItem('isAuthorized') ?
                           (<LinkPage/>)
                       }
                />
                <Route path="/contest-details/:groupCode/:contestId"
                       element={
                           // localStorage.getItem('isAuthorized') ?
                           (<ContestDetails/>)
                           // : (<Navigate to='/'/>)
                       }
                />
                <Route path="/weights-distribution/:groupCode/:contestId"
                       element={
                           // localStorage.getItem('isAuthorized') ?
                           (<WeightsDistrPage/>)
                           // : (<Navigate to='/'/>)
                       }
                />
                <Route path="/late-submissions/:groupCode/:contestId"
                       element={
                           // localStorage.getItem('isAuthorized') ?
                           (<LateSubmissionPage/>)
                           // : (<Navigate to='/'/>)
                       }
                />
                <Route path="/upload-csv/:groupCode/:contestId"
                       element={
                           // localStorage.getItem('isAuthorized') ?
                           (<UploadCsvFilePage/>)
                           // : (<Navigate to='/'/>)
                       }
                />
                {/*<Route path="*" element={<Navigate to="/" replace/>}/>*/}
            </Routes>
        </Router>
    );
};

export default App;
