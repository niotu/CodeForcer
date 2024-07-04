import React from 'react';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import LoginPage from './components/LoginPage';
import GroupsPage from './components/GroupsPage';
import ContestsPage from './components/ContestsPage';
import ContestDetails from './components/ContestDetails';
import WeightsDistrPage from "./components/WeightsDistrPage.jsx";

const App = () => {
    return (
        <Router>
            <Routes>
                {this.localStorage.getItem("isAuthorized") === false &&
                    <Route path="/" element={<LoginPage/>}/>
                }
                {this.localStorage.getItem('isAuthorized') &&
                    <Route path="/groups" element={<GroupsPage/>}/>
                }
                {this.localStorage.getItem('isAuthorized') &&
                    <Route path="/contests/:groupCode" element={<ContestsPage/>}/>
                }
                {this.localStorage.getItem('isAuthorized') &&
                    <Route path="/contest-details/:groupCode/:contestId" element={<ContestDetails/>}/>
                }
                {this.localStorage.getItem('isAuthorized') &&
                    <Route path="/weightsDistr/:GroupCode/:contestId" element={<WeightsDistrPage/>}/>
                }

            </Routes>
        </Router>
    );
};

export default App;
