import React from 'react';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import LoginPage from './components/LoginPage.jsx';
import GroupsPage from './components/GroupsPage.jsx';
import ContestsPage from './components/ContestsPage.jsx';
import ContestDetails from './components/ContestDetails.jsx';
import WeightsDistrPage from "./components/WeightsDistrPage.jsx";

const App = () => {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<LoginPage/>}/>

                {localStorage.getItem('isAuthorized') &&
                    <Route path="/groups" element={<GroupsPage/>}/>
                }
                {localStorage.getItem('isAuthorized') &&
                    <Route path="/contests/:groupCode" element={<ContestsPage/>}/>
                }
                {localStorage.getItem('isAuthorized') &&
                    <Route path="/contest-details/:groupCode/:contestId" element={<ContestDetails/>}/>
                }
                {localStorage.getItem('isAuthorized') &&
                    <Route path="/weightsDistr/:GroupCode/:contestId" element={<WeightsDistrPage/>}/>
                }

            </Routes>
        </Router>
    );
};

export default App;
